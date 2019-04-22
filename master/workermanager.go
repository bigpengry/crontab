package master

import (
	"context"
	"strings"
	"time"

	"github.com/bigpengry/crontab/common"
	"github.com/coreos/etcd/clientv3"
)

// WorkerManager 节点
type WorkerManager struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

// WorkerMgr 全局变量
var (
	WorkerMgr *WorkerManager
)

// ListWorkers 获取在线worker列表
func (m *WorkerManager) ListWorkers() (workerArr []string, err error) {

	// 初始化数组
	workerArr = make([]string, 0)

	// 获取目录下所有Kv
	getResp, err := m.kv.Get(context.TODO(), common.JOB_WORKER_DIR, clientv3.WithPrefix())
	if err != nil {
		return
	}

	// 解析每个节点的IP
	for _, kv := range getResp.Kvs {
		// kv.Key : /cron/workers/192.168.2.1
		workerIP := strings.TrimPrefix(string(kv.Key), common.JOB_WORKER_DIR)
		workerArr = append(workerArr, workerIP)
	}
	return
}

// InitWorkerManager 初始化节点管理器
func InitWorkerManager() (err error) {
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

	WorkerMgr = &WorkerManager{
		client: cli,
		kv:     kv,
		lease:  lease,
	}
	return
}
