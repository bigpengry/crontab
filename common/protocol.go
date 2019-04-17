package common

import "encoding/json"

// Job 任务对象
type Job struct {
	Name     string `json:"name"`     //任务名称
	Command  string `json:"command"`  //shell命令
	CronExpr string `json:"cronExpr"` //cron表达式
}

// Response 用于存储将要返回的任务信息
type Response struct {
	ErrorType int         `json:"errorType"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}

// ResponseMarshal 序列化将要返回的信息
func (resp *Response) ResponseMarshal() (response []byte, err error) {
	response, err = json.Marshal(resp)
	return
}
