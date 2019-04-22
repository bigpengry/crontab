package worker

import (
	"errors"
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"

	"github.com/bigpengry/crontab/common"
)

// SchedulingQueue 任务调度队列
type SchedulingQueue struct {
	jobEventChan   chan *common.JobEvent
	schedule       map[string]*common.Task
	execTable      map[string]*common.TaskStatus
	taskResultChan chan *common.TaskExecResult
}

// Queue 全局变量
var (
	Queue *SchedulingQueue
)

// ExecuteSchedule 执行任务调度
func (q *SchedulingQueue) ExecuteSchedule() (NextSchedulingTime time.Duration) {
	recentTime := new(time.Time)
	now := time.Now()

	if len(q.schedule) == 0 {
		NextSchedulingTime = 1 * time.Second
		return
	}

	for _, task := range q.schedule {
		if task.NextTime.Before(now) || task.NextTime.Equal(now) {
			q.ExecuteTask(task)
			// 更新执行时间
			task.NextTime = task.Expr.Next(now)
		}
		//最近即将过期任务的时间
		if recentTime.IsZero() || task.NextTime.Before(*recentTime) {
			*recentTime = task.NextTime
		}
	}

	// 计算下次执行任务的间隔
	NextSchedulingTime = (recentTime).Sub(now)
	return
}

// handleJobEvent 处理事件
func (q *SchedulingQueue) handleJobEvent(event *common.JobEvent) {
	task := new(common.Task)

	// 判断任务事件类型
	switch event.EventType {
	case common.JOB_EVENT_SAVE:
		expr, err := cronexpr.Parse(event.CronExpr)
		if err != nil {
			fmt.Println(err)
			return
		}
		task = common.NewTask(event.Job, expr)
		q.schedule[event.Name] = task
	case common.JOB_EVENT_DELETE:
		_, ok := q.schedule[event.Name]
		if ok {
			delete(q.schedule, event.Name)
		}
	case common.JOB_EVENT_KILL:
		taskStatus, ok := q.execTable[event.Name]
		if ok {
			taskStatus.Cancle()
		}
	}
}

func (q *SchedulingQueue) handleTaskResult(r *common.TaskExecResult) {
	delete(q.execTable, r.Name)
	err := ""
	if r.Error != nil {
		err = r.Error.Error()
	}
	if r.Error != errors.New("锁已被占用") {
		taskLog := common.NewTaskLog(r.Name, r.Command, string(r.Output), err,
									 r.ScheduleTime.UnixNano()/1000000,
									 r.ExecuteTime.UnixNano()/1000000,
									 r.StartTime.UnixNano()/1000000,
									 r.EndTime.UnixNano()/1000000)
		Log.Append(taskLog)
	}
	fmt.Println("任务执行完成", r.Name, string(r.Output), r.Error)
}

// ExecuteTask 任务执行模块
func (q *SchedulingQueue) ExecuteTask(t *common.Task) {

	taskStatus, ok := q.execTable[t.Name]
	if ok {
		fmt.Println("尚未退出，跳过执行：", taskStatus.Name)
		return
	}
	taskStatus = common.NewTaskStatus(t)
	q.execTable[t.Name] = taskStatus
	//fmt.Println("执行任务：", taskStatus.Name, taskStatus.ScheduleTime, taskStatus.ExecuteTime)
	ExecuteJob(taskStatus)
}

// Scheduler 任务调度器
func (q *SchedulingQueue) Scheduler() {
	nextTime := q.ExecuteSchedule()
	timer := time.NewTimer(nextTime)
	jobEvent := new(common.JobEvent)
	result := new(common.TaskExecResult)
	for {
		select {
		case jobEvent = <-q.jobEventChan:
			q.handleJobEvent(jobEvent)
		case <-timer.C:
		case result = <-q.taskResultChan:
			q.handleTaskResult(result)
		}

		// 重置调度间隔
		nextTime = q.ExecuteSchedule()
		timer.Reset(nextTime)
	}
}

// InitScheduler 初始化任务调度
func InitScheduler() (err error) {
	Queue = &SchedulingQueue{
		jobEventChan:   make(chan *common.JobEvent, 1000),
		schedule:       make(map[string]*common.Task),
		execTable:      make(map[string]*common.TaskStatus),
		taskResultChan: make(chan *common.TaskExecResult, 1000),
	}

	// 启动调度器
	go Queue.Scheduler()
	return
}
