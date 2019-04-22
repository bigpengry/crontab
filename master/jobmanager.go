package master

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bigpengry/crontab/common"

	"github.com/coreos/etcd/clientv3"
)

// JobMgr 是一个JobManager对象的全局化实例
var (
	JobMgr *JobManager
)

// JobManager 任务管理器
type JobManager struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

// SaveJob 接受一个Job对象，把任务写入到etcd的“/cron/jobs/"目录下，返回上次job的信息和错误信息
func (m *JobManager) SaveJob(j *common.Job) (oldJob *common.Job, err error) {
	oldJob = new(common.Job)
	// 任务的key
	jobKey := common.JOB_SAVE_DIR + j.Name
	// 序列化
	jobValue, err := json.Marshal(j)
	if err != nil {
		return
	}
	// 保存到ETCD
	presp, err := m.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV())
	if err != nil {
		return
	}

	if presp.PrevKv != nil {
		if err = json.Unmarshal(presp.PrevKv.Value, oldJob); err != nil {
			return
		}
		return
	}
	return
}

// DeleteJob 输入一个任务的名称，从etcd中删除这个任务，返回被删除任务的信息和错误信息
func (m *JobManager) DeleteJob(jobName string) (oldJob *common.Job, err error) {
	oldJob = new(common.Job)

	// 得到key
	jobKey := common.JOB_SAVE_DIR + jobName

	// 从etcd中删除key
	dresp, err := m.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV())
	if err != nil {
		return
	}

	// 返回被删除的任务信息
	if len(dresp.PrevKvs) != 0 {
		if err = json.Unmarshal(dresp.PrevKvs[0].Value, oldJob); err != nil {
			err = nil
			return
		}
	}

	return
}

// ListJob 列出etcd中的所有任务，返回多个任务组成的切片和错误信息
func (m *JobManager) ListJob() (jobArr []common.Job, err error) {
	job := new(common.Job)

	// 获取任务目录
	dir := common.JOB_SAVE_DIR
	gresp, err := m.kv.Get(context.TODO(), dir, clientv3.WithPrefix())
	if err != nil {
		return
	}

	// 获取键值对并反序列化
	for _, kvPair := range gresp.Kvs {
		if err = json.Unmarshal(kvPair.Value, job); err != nil {
			err = nil
			continue
		}
		jobArr = append(jobArr, *job)
	}
	return
}

// KillJob 强制结束一个任务，输入任务的名称，把任务写入到etcd的“/cron/killer/"目录下，返回错误信息
func (m *JobManager) KillJob(jobName string) (err error) {
	// 获取key
	jobKey := common.JOB_KILLER_DIR + jobName
	// 创建租约
	resp, err := m.lease.Grant(context.TODO(), 1)
	if err != nil {
		return
	}

	// 设置标记
	_, err = m.kv.Put(context.TODO(), jobKey, "", clientv3.WithLease(resp.ID))
	if err != nil {
		return
	}
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

	//创建租约
	lease := clientv3.NewLease(cli)

	JobMgr = &JobManager{
		client: cli,
		kv:     kv,
		lease:  lease,
	}
	return
}
