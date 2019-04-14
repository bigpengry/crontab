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

func (jobManager *JobManager) SvaeJob(job *common.Job) (old *common.Job, err error) {
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
		oldJob := new(common.Job)
		if err = json.Unmarshal(putResponse.PrevKv.Value, oldJob); err != nil {
			return
		}
		return
	}
	return
}

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
