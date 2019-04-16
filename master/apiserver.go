package master

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/bigpengry/crontab/common"
)

//设置单例
var (
	G_APIServer *APIServer
)

//任务的HTTP接口
type APIServer struct {
	httpServer *http.Server
}

//保存任务
func hanleJobSave(w http.ResponseWriter, r *http.Request) {
	job := new(common.Job)
	resp := new(common.Response)
	//1.解析表单(错误处理可以改进)
	if err := r.ParseForm(); err != nil {
		resp.ErrorType = -1
		resp.Message = err.Error()
		byte, err := resp.ResponseMarshal()
		if err == nil {
			w.Write(byte)
		}
		return
	}
	//2.获取job字段
	postJob := r.PostForm.Get("job")
	//3.反序列化
	if err := json.Unmarshal([]byte(postJob), job); err != nil {
		resp.ErrorType = -1
		resp.Message = err.Error()
		byte, err := resp.ResponseMarshal()
		if err == nil {
			w.Write(byte)
		}
		return
	}
	//4.保存到ETCD
	oldJob, err := G_jobManager.SaveJob(job)
	if err != nil {
		resp.ErrorType = -1
		resp.Message = err.Error()
		byte, err := resp.ResponseMarshal()
		if err == nil {
			w.Write(byte)
		}
		return
	}

	resp.ErrorType = 0
	resp.Message = "success"
	resp.Data = oldJob
	byte, err := resp.ResponseMarshal()
	if err == nil {
		w.Write(byte)
	}

	return
}

//删除任务接口
func handleJobDelete(w http.ResponseWriter, r *http.Request) {
	resp := new(common.Response)
	if err := r.ParseForm(); err != nil {
		resp.ErrorType = -1
		resp.Message = err.Error()
		byte, err := resp.ResponseMarshal()
		if err == nil {
			w.Write(byte)
		}
		return
	}
	//获取任务名
	jobName := r.PostForm.Get("name")
	oldJob, err := G_jobManager.DeleteJob(jobName)
	if err != nil {
		resp.ErrorType = -1
		resp.Message = err.Error()
		byte, err := resp.ResponseMarshal()
		if err == nil {
			w.Write(byte)
		}
		return
	}
	resp.ErrorType = 0
	resp.Message = "success"
	resp.Data = oldJob
	byte, err := resp.ResponseMarshal()
	if err == nil {
		w.Write(byte)
	}

}

//任务列表接口
func handleJobList(w http.ResponseWriter, r *http.Request) {
	resp := new(common.Response)
	//获取任务列表
	jobList, err := G_jobManager.ListJob()
	if err != nil {
		//错误处理
		return
	}
	resp.ErrorType = 0
	resp.Message = "success"
	resp.Data = jobList
	byte, err := resp.ResponseMarshal()
	if err == nil {
		w.Write(byte)
	}

}

//强制杀死任务
func handleJobKill(w http.ResponseWriter, r *http.Request) {
	resp := new(common.Response)
	if err := r.ParseForm(); err != nil {
		resp.ErrorType = -1
		resp.Message = err.Error()
		byte, err := resp.ResponseMarshal()
		if err == nil {
			w.Write(byte)
		}
		return
	}

	//获取任务名称
	jobName := r.PostForm.Get("name")

	if err := G_jobManager.KillJob(jobName); err != nil {
		resp.ErrorType = -1
		resp.Message = err.Error()
		byte, err := resp.ResponseMarshal()
		if err == nil {
			w.Write(byte)
		}
		return

	}
	resp.ErrorType = 0
	resp.Message = "success"
	byte, err := resp.ResponseMarshal()
	if err == nil {
		w.Write(byte)
	}

}

//初始化服务
func InitAPIServer() (err error) {
	//配置路由
	mux := &http.ServeMux{}
	mux.HandleFunc("/job/save", hanleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)
	//静态文件目录
	staticDir := http.Dir("./webroot")

	//启动TCP监听
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(G_config.APIPort))
	if err != nil {
		return
	}

	//配置HTTP服务
	httpServer := &http.Server{
		ReadTimeout:  time.Duration(G_config.APIReadTimeOut) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.APIWriteTimeOut) * time.Millisecond,
		Handler:      mux,
	}

	//赋值单例
	G_APIServer = &APIServer{
		httpServer: httpServer,
	}

	//启动HTTP服务
	go httpServer.Serve(listener)
	return
}
