package worker

import (
	"math/rand"
	"os/exec"
	"time"

	"github.com/bigpengry/crontab/common"
)

// ExecuteJob 执行任务
func ExecuteJob(s *common.TaskStatus) {
	go func() {
		result := new(common.TaskExecResult)
		result.TaskStatus = s
		// 上锁
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		jobLock := NewJobLock(JobMgr.kv, JobMgr.lease, s.Name)
		result.StartTime = time.Now()
		err := jobLock.Lock()
		defer jobLock.UnLock()
		if err != nil {
			result.Error = err
			result.EndTime = time.Now()
			Queue.taskResultChan <- result
			return
		}
		result.StartTime = time.Now()
		cmd := exec.CommandContext(s.Context, "/bin/bash", "-c", s.Command)
		output, err := cmd.CombinedOutput()

		result.EndTime = time.Now()
		result.Output = output
		result.Error = err

		Queue.taskResultChan <- result

	}()
}
