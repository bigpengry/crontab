package master

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bigpengry/crontab/common"

	"github.com/coreos/etcd/clientv3"
)

var (
	G_jobManager *JobManager
)

//任务管理器
type JobManager struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

//保存任务
func (jobManager *JobManager) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	//任务的key
	jobKey := "/cron/jobs/" + job.Name
	//反序列化
	jobValue, err := json.Marshal(job)
	if err != nil {
		return
	}
	//保存到ETCD
	putResponse, err := jobManager.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV())
	if err != nil {
		return
	}
	if putResponse.PrevKv != nil {
		if err = json.Unmarshal(putResponse.PrevKv.Value, oldJob); err != nil {
			return
		}
		return
	}
	return
}

//删除任务
func (jobManager *JobManager) DeleteJob(jobName string) (oldJob *common.Job, err error) {
	//得到key
	jobKey := "/cron/jobs/" + jobName

	//从etcd中删除key
	deleteRespoonse, err := jobManager.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV())
	if err != nil {
		return
	}

	//返回被删除的任务信息
	if len(deleteRespoonse.PrevKvs) != 0 {

		if err = json.Unmarshal(deleteRespoonse.PrevKvs[0].Value, oldJob); err != nil {
			err = nil
			return
		}
	}

	return
}

//列出所有任务
func (jobManager *JobManager) ListJob() (jobList []*common.Job, err error) {
	job := new(common.Job)
	jobList = make([]*common.Job, 0)
	directory := "/cron/jobs"
	getResponse, err := jobManager.kv.Get(context.TODO(), directory, clientv3.WithPrefix())
	if err != nil {
		return
	}
	for _, kvPair := range getResponse.Kvs {
		if err = json.Unmarshal(kvPair.Value, job); err != nil {
			err = nil
			continue
		}
		jobList = append(jobList, job)
	}
	return
}

//初始化任务管理器
func InitJobManager() (err error) {
	config := clientv3.Config{
		Endpoints:   G_config.ETCDEndPoints,
		DialTimeout: time.Duration(G_config.ETCDDialTimeOut) * time.Millisecond,
	}

	client, err := clientv3.New(config)
	if err != nil {
		return
	}

	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)
	G_jobManager = &JobManager{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}
