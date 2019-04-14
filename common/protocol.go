package common

import "encoding/json"

//定时任务
type Job struct {
	Name     string `json:"name"`     //任务名称
	Command  string `json:"command"`  //shell命令
	CronExpr string `json:"cronExpr"` //cron表达式
}

//HTTP接口
type Response struct {
	ErrorType int         `json:"errorType"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}

//对reponse对象进行json序列化
func (resp *Response) ResponseMarshal() (response []byte, err error) {
	response, err = json.Marshal(resp)
	return
}
