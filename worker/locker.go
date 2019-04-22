package worker

import (
	"context"
	"errors"

	"github.com/bigpengry/crontab/common"

	"github.com/coreos/etcd/clientv3"
)

// JobLock 任务锁
type JobLock struct {
	kv       clientv3.KV
	lease    clientv3.Lease
	jobName  string
	cancel   context.CancelFunc
	leaseID  clientv3.LeaseID
	isLocked bool
}

// NewJobLock 构造函数
func NewJobLock(k clientv3.KV, l clientv3.Lease, jobName string) *JobLock {
	return &JobLock{
		kv:      k,
		lease:   l,
		jobName: jobName,
	}
}

// Lock 上锁
func (l *JobLock) Lock() (err error) {
	// 创建租约
	gresp, err := l.lease.Grant(context.TODO(), 5)
	if err != nil {
		return
	}
	leaseID := gresp.ID
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer l.lease.Revoke(context.TODO(), leaseID)

	// 续租
	ch, err := l.lease.KeepAlive(ctx, leaseID)
	if err != nil {
		cancel()
		l.lease.Revoke(context.TODO(), leaseID)
		return
	}
	go func() {
		resp := new(clientv3.LeaseKeepAliveResponse)
		select {
		case resp = <-ch:
			if resp == nil {
				return
			}

		}
	}()

	txn := l.kv.Txn(context.TODO())
	lockKey := common.JOB_LOCK_DIR + l.jobName
	// 抢锁
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseID))).
		Else(clientv3.OpGet(lockKey))

	tresp, err := txn.Commit()
	if err != nil {
		cancel()
		l.lease.Revoke(context.TODO(), leaseID)
		return
	}
	if !tresp.Succeeded {
		cancel()
		l.lease.Revoke(context.TODO(), leaseID)
		err = errors.New("锁已被占用")
		return
	}

	l.leaseID = leaseID
	l.cancel = cancel
	l.isLocked = true
	return
}

// UnLock 解锁
func (l *JobLock) UnLock() {
	if l.isLocked {
		l.cancel()
		l.lease.Revoke(context.TODO(), l.leaseID)
	}
}
