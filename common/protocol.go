package common

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gorhill/cronexpr"
)

// Job 任务对象
type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

// JobEvent 任务事件
type JobEvent struct {
	EventType int
	*Job
}

// Response 用于存储将要返回的任务信息
type Response struct {
	ErrorType int         `json:"errorType"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}

// Task 任务的调度单位
type Task struct {
	*Job
	Expr     *cronexpr.Expression
	NextTime time.Time
}

// TaskStatus 任务状态,
type TaskStatus struct {
	*Job
	ScheduleTime time.Time
	ExecuteTime  time.Time
	Context      context.Context
	Cancle       context.CancelFunc
}

// TaskExecResult 任务执行结果
type TaskExecResult struct {
	*TaskStatus
	Output    []byte
	Error     error
	StartTime time.Time
	EndTime   time.Time
}

// TaskLog 任务执行日志
type TaskLog struct {
	JobName      string `bson:"jobName"`
	Command      string `bson:"command"`
	Output       string `bson:"output"`
	Error        string `bson:"error"`
	ScheduleTime int64  `bson:"scheduleTime	"`
	ExecuteTime  int64  `bson:"executeTime"`
	StartTime    int64  `bson:"startTime"`
	EndTime      int64  `bson:"endTime"`
}

type TaskFilter struct {
	JobName string`bson:"jobName"`
}

type SortLogByStartTime struct {
	SortOrder int `bson:"startTime"`
}

// NewResponse 构造函数
func NewResponse(errorType int, msg string, data interface{}) *Response {
	return &Response{
		ErrorType: errorType,
		Message:   msg,
		Data:      data,
	}
}

// NewJobEvent 构造函数
func NewJobEvent(eventType int, j *Job) *JobEvent {
	return &JobEvent{
		EventType: eventType,
		Job:       j,
	}
}

// NewTask 构造函数
func NewTask(j *Job, expr *cronexpr.Expression) *Task {
	return &Task{
		Job:      j,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}
}

// NewTaskLog 构造函数
func NewTaskLog(jobName, cmd, output, err string,
	scheduleTime, executeTime, startTime, endTime int64) *TaskLog {
	return &TaskLog{
		JobName:      jobName,
		Command:      cmd,
		Output:       output,
		ScheduleTime: scheduleTime,
		ExecuteTime:  executeTime,
		StartTime:    startTime,
		EndTime:      endTime,
		Error:        err,
	}
}

// NewTaskStatus 构造函数
func NewTaskStatus(t *Task) *TaskStatus {
	ctx, cancel := context.WithCancel(context.Background())
	return &TaskStatus{
		Job:          t.Job,
		ScheduleTime: t.NextTime,
		ExecuteTime:  time.Now(),
		Context:      ctx,
		Cancle:       cancel,
	}
}

// MarshalResponse 序列化将要返回的信息
func (r *Response) MarshalResponse() (respArr []byte, err error) {
	respArr, err = json.Marshal(r)
	return
}
