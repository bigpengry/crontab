package worker

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/coreos/etcd/mvcc/mvccpb"

	"github.com/bigpengry/crontab/common"

	"github.com/coreos/etcd/clientv3"
)

// JobManager 任务管理器
type JobManager struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

// JobMgr 是一个JobManager对象的全局化实例
var (
	JobMgr *JobManager
)

// watchJobs 监听etcd中的任务变化
func (m *JobManager) watchJobs() (err error) {
	job := new(common.Job)
	jobEvent := new(common.JobEvent)
	// 获取当前etcd中的任务
	gresp, err := m.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix())
	if err != nil {
		return
	}

	for _, kvPair := range gresp.Kvs {
		job := new(common.Job)
		if err = json.Unmarshal(kvPair.Value, job); err != nil {
			continue
		}
		jobEvent = common.NewJobEvent(common.JOB_EVENT_SAVE, job)
		Queue.jobEventChan <- jobEvent
	}

	// 监听协程
	go func() {
		startRevision := gresp.Header.Revision + 1
		ch := m.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithRev(startRevision), clientv3.WithPrefix())

		for resp := range ch {
			for _, event := range resp.Events {

				switch event.Type {
				case mvccpb.PUT:
					if err = json.Unmarshal(event.Kv.Value, job); err != nil {
						continue
					}
					jobEvent = common.NewJobEvent(common.JOB_EVENT_SAVE, job)

					// 同步给调度协程(更新)
					Queue.jobEventChan <- jobEvent
				case mvccpb.DELETE:
					jobName := strings.TrimPrefix(string(event.Kv.Key), common.JOB_SAVE_DIR)
					job.Name = jobName
					jobEvent = common.NewJobEvent(common.JOB_EVENT_DELETE, job)

					// 同步给调度协程(删除)
					Queue.jobEventChan <- jobEvent
				}
			}
		}
	}()

	return
}

// watchKiller 监听强杀
func (m *JobManager) watchKiller() (err error) {
	job := new(common.Job)
	// 监听协程
	go func() {
		ch := m.watcher.Watch(context.TODO(), common.JOB_KILLER_DIR, clientv3.WithPrefix())

		for resp := range ch {
			for _, event := range resp.Events {
				switch event.Type {
				// 杀死任务
				case mvccpb.PUT:
					jobName := strings.TrimPrefix(string(event.Kv.Key), common.JOB_KILLER_DIR)
					job.Name = jobName
					jobEvent := common.NewJobEvent(common.JOB_EVENT_KILL, job)
					Queue.jobEventChan <- jobEvent
				case mvccpb.DELETE:

				}
			}
		}
	}()
	return
}

// InitJobManager 初始化任务管理器
func InitJobManager() (err error) {
	// 设置配置文件
	conf := clientv3.Config{
		Endpoints:   Conf.ETCDEndPoints,
		DialTimeout: time.Duration(Conf.ETCDDialTimeOut) * time.Millisecond,
	}

	// 创建客户端
	cli, err := clientv3.New(conf)
	if err != nil {
		return
	}

	// 创建kv
	kv := clientv3.NewKV(cli)

	// 创建租约
	lease := clientv3.NewLease(cli)

	// 创建watcher
	watcher := clientv3.NewWatcher(cli)

	JobMgr = &JobManager{
		client:  cli,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}

	//　启动监听
	JobMgr.watchJobs()
	JobMgr.watchKiller()
	return
}
