package worker

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/bigpengry/crontab/common"
	"github.com/coreos/etcd/clientv3"
)

// Register 节点注册
type Register struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
	wokerIP string
}

// REG 全局变量
var (
	REG *Register
)

func getWorkerIP() (ipv4 string, err error) {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	for _, addr := range addrs {
		ipNet, isipNet := addr.(*net.IPNet)
		if isipNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String()
				return
			}
		}

	}
	err = errors.New("没有找到网卡IP")
	return
}

func (r *Register) keepOnline() {
	for {
		regKey := common.JOB_WORKER_DIR + REG.wokerIP
		ctx, cancel := context.WithCancel(context.Background())
		resp, err := r.lease.Grant(context.TODO(), 5)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		ch, err := r.lease.KeepAlive(context.TODO(), resp.ID)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		_, err = r.kv.Put(ctx, regKey, "", clientv3.WithLease(resp.ID))
		if err != nil {
			cancel()
		}
		for {
			select {
			case keepAliveResp := <-ch:
				if keepAliveResp == nil {
					cancel()
				}

			}
		}

	}

}

// InitResister 初始化
func InitResister() (err error) {
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

	ip, err := getWorkerIP()
	if err != nil {
		return
	}

	// 创建kv
	kv := clientv3.NewKV(cli)

	// 创建租约
	lease := clientv3.NewLease(cli)

	// 创建watcher
	watcher := clientv3.NewWatcher(cli)

	REG = &Register{
		client:  cli,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
		wokerIP: ip,
	}

	go REG.keepOnline()
	return
}
