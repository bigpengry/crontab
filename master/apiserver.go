package master

import (
	"net"
	"net/http"
	"strconv"
	"time"
)

//设置单例
var(
	 G_APIServer *APIServer
)


//任务的HTTP接口
type APIServer struct {
	httpServer *http.Server
}

//保存任务
func hanleJobSave(w http.ResponseWriter,r *http.Request){

}

//初始化服务
func InitAPIServer()(err error){
	//配置路由
	mux:=&http.ServeMux{}
	mux.HandleFunc("/job/save",hanleJobSave)

	//启动TCP监听
	listener,err:=net.Listen("tcp",":"+strconv.Itoa(G_config.APIPort))
	if err!=nil {
		return
	}

	//配置HTTP服务
	httpServer:=&http.Server{
		ReadTimeout:time.Duration(G_config.APIReadTimeOut)*time.Millisecond,
		WriteTimeout:time.Duration(G_config.APIWriteTimeOut)*time.Millisecond,
		Handler:mux,
	}

	//赋值单例
	G_APIServer=&APIServer{
		httpServer:httpServer,
	}

	//启动HTTP服务
	go httpServer.Serve(listener)
	return
}