package master

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/bigpengry/crontab/common"
)

// Svr 是一个apiServer对象的全局化实例
var (
	Svr *http.Server
)

//hanleJobSave 保存任务
func hanleJobSave(w http.ResponseWriter, r *http.Request) {
	job := new(common.Job)
	resp := new(common.Response)
	//1.解析表单(错误处理可以改进)
	if err := r.ParseForm(); err != nil {
		resp = common.NewResponse(-1, err.Error(), nil)
		byte, err := resp.MarshalResponse()
		if err != nil {
			return
		}
		w.Write(byte)
		return
	}
	//2.获取job字段
	postJob := r.PostForm.Get("job")
	//3.反序列化
	if err := json.Unmarshal([]byte(postJob), job); err != nil {
		resp = common.NewResponse(-2, err.Error(), nil)
		byte, err := resp.MarshalResponse()
		if err != nil {
			return
		}
		w.Write(byte)
		return
	}
	//4.保存到ETCD
	oldJob, err := JobMgr.SaveJob(job)
	if err != nil {
		resp = common.NewResponse(-2, err.Error(), nil)
		byte, err := resp.MarshalResponse()
		if err != nil {
			return
		}
		w.Write(byte)
		return
	}

	resp = common.NewResponse(0, "success", oldJob)
	byte, err := resp.MarshalResponse()
	if err != nil {
		return
	}
	w.Write(byte)
}

// handleJobDelete 删除任务
func handleJobDelete(w http.ResponseWriter, r *http.Request) {
	resp := new(common.Response)
	if err := r.ParseForm(); err != nil {
		resp = common.NewResponse(-1, err.Error(), nil)
		byte, err := resp.MarshalResponse()
		if err != nil {
			return
		}
		w.Write(byte)
		return
	}
	// 获取任务名
	jobName := r.PostForm.Get("name")
	oldJob, err := JobMgr.DeleteJob(jobName)
	if err != nil {
		resp = common.NewResponse(-2, err.Error(), nil)
		byte, err := resp.MarshalResponse()
		if err != nil {
			return
		}
		w.Write(byte)
		return
	}
	resp = common.NewResponse(0, "success", oldJob)
	byte, err := resp.MarshalResponse()
	if err != nil {
		return
	}
	w.Write(byte)
}

// handleJobList 返回所有任务信息
func handleJobList(w http.ResponseWriter, r *http.Request) {
	resp := new(common.Response)
	//获取任务列表
	jobList, err := JobMgr.ListJob()
	if err != nil {
		resp = common.NewResponse(-2, err.Error(), nil)
		byte, err := resp.MarshalResponse()
		if err != nil {
			return
		}
		w.Write(byte)
		return
	}
	resp = common.NewResponse(0, "success", jobList)
	byte, err := resp.MarshalResponse()
	if err != nil {
		return
	}
	w.Write(byte)

}

// handleJobKill 强制杀死任务
func handleJobKill(w http.ResponseWriter, r *http.Request) {
	resp := new(common.Response)
	if err := r.ParseForm(); err != nil {
		resp = common.NewResponse(-1, err.Error(), nil)
		byte, err := resp.MarshalResponse()
		if err != nil {
			return
		}
		w.Write(byte)
		return
	}

	// 获取任务名称
	jobName := r.PostForm.Get("name")

	if err := JobMgr.KillJob(jobName); err != nil {
		resp = common.NewResponse(-2, err.Error(), nil)
		byte, err := resp.MarshalResponse()
		if err != nil {
			return
		}
		w.Write(byte)
		return

	}
	resp = common.NewResponse(0, "success", nil)
	byte, err := resp.MarshalResponse()
	if err != nil {
		return
	}
	w.Write(byte)

}

// handleJobLog	查询日志
func handleJobLog(w http.ResponseWriter, r *http.Request) {
	resp := new(common.Response)
	if err := r.ParseForm(); err != nil {
		resp = common.NewResponse(-1, err.Error(), nil)
		byte, err := resp.MarshalResponse()
		if err != nil {
			return
		}
		w.Write(byte)
		return
	}

	name := r.Form.Get("name")
	skipParam := r.Form.Get("skip")
	limitParam := r.Form.Get("limit")

	skip, err := strconv.Atoi(skipParam)
	if err != nil {
		skip = 0
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 20
	}

	logArr, err := LogMgr.ListLog(name, skip, limit)
	if err != nil {
		resp = common.NewResponse(-1, err.Error(), nil)
		byte, err := resp.MarshalResponse()
		if err != nil {
			return
		}
		w.Write(byte)
		return
	}

	resp = common.NewResponse(0, "success", logArr)
	byte, err := resp.MarshalResponse()
	if err != nil {
		return
	}
	w.Write(byte)

}

func handleWorkerList(w http.ResponseWriter, r *http.Request) {
	resp := new(common.Response)
	workerArr, err := WorkerMgr.ListWorkers()
	if err != nil {
		resp = common.NewResponse(-1, err.Error(), nil)
		byte, err := resp.MarshalResponse()
		if err != nil {
			return
		}
		w.Write(byte)
		return
	}

	resp = common.NewResponse(0, "success", workerArr)
	byte, err := resp.MarshalResponse()
	if err != nil {
		return
	}
	w.Write(byte)
}

// InitHTTPServer 初始化服务http服务器
func InitHTTPServer() (err error) {
	// 配置路由
	mux := &http.ServeMux{}
	mux.HandleFunc("/job/save", hanleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)
	mux.HandleFunc("/job/log", handleJobLog)
	mux.HandleFunc("/worker/list", handleWorkerList)
	// 静态文件目录
	staticDir := http.Dir(Conf.WebRoot)
	staticHandler := http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler))

	// 启动TCP监听
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(Conf.APIPort))
	if err != nil {
		return
	}

	// 配置HTTP服务
	httpServer := &http.Server{
		ReadTimeout:  time.Duration(Conf.APIReadTimeOut) * time.Millisecond,
		WriteTimeout: time.Duration(Conf.APIWriteTimeOut) * time.Millisecond,
		Handler:      mux,
	}

	// 赋值单例
	Svr = httpServer

	// 启动HTTP服务
	go httpServer.Serve(listener)
	return
}
