package master

import (
	"github.com/coreos/etcd/clientv3"
	"time"
)

var(
	G_jobManager *JobManager
)

//任务管理器
type JobManager struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
}

func InitJobManager()(err error){
	config:=clientv3.Config{
		Endpoints:G_config.ETCDEndPoints,
		DialTimeout:time.Duration(G_config.ETCDDialTimeOut)*time.Millisecond,
	}

	client,err:=clientv3.New(config)
	if err!=nil {
		return
	}

	kv:=clientv3.NewKV(client)
	lease:=clientv3.NewLease(client)
	G_jobManager=&JobManager{
		client:client,
		kv:kv,
		lease:lease,
	}
	return
}
