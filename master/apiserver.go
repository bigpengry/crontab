package master

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/bigpengry/crontab/common"
)

// APIServer 是一个apiServer对象的全局化实例
var (
	APIServer *http.Server
)

//hanleJobSave 保存任务
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
		resp.ErrorType = -2
		resp.Message = err.Error()
		byte, err := resp.ResponseMarshal()
		if err == nil {
			w.Write(byte)
		}
		return
	}
	//4.保存到ETCD
	oldJob, err := GJobManager.SaveJob(job)
	if err != nil {
		resp.ErrorType = -2
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

// handleJobDelete 删除任务
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
	// 获取任务名
	jobName := r.PostForm.Get("name")
	oldJob, err := GJobManager.DeleteJob(jobName)
	if err != nil {
		resp.ErrorType = -2
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

// handleJobList 返回所有任务信息
func handleJobList(w http.ResponseWriter, r *http.Request) {
	resp := new(common.Response)
	//获取任务列表
	jobList, err := GJobManager.ListJob()
	if err != nil {
		resp.ErrorType = -2
		resp.Message = err.Error()
		byte, err := resp.ResponseMarshal()
		if err == nil {
			w.Write(byte)
		}
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

// handleJobKill 强制杀死任务
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

	// 获取任务名称
	jobName := r.PostForm.Get("name")

	if err := GJobManager.KillJob(jobName); err != nil {
		resp.ErrorType = -2
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

// InitAPIServer 初始化服务
func InitAPIServer() (err error) {
	// 配置路由
	mux := &http.ServeMux{}
	mux.HandleFunc("/job/save", hanleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)
	// 静态文件目录
	staticDir := http.Dir(GConfig.WebRoot)
	staticHandler := http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler))

	// 启动TCP监听
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(GConfig.APIPort))
	if err != nil {
		return
	}

	// 配置HTTP服务
	httpServer := &http.Server{
		ReadTimeout:  time.Duration(GConfig.APIReadTimeOut) * time.Millisecond,
		WriteTimeout: time.Duration(GConfig.APIWriteTimeOut) * time.Millisecond,
		Handler:      mux,
	}

	// 赋值单例
	APIServer = httpServer

	// 启动HTTP服务
	go httpServer.Serve(listener)
	return
}
