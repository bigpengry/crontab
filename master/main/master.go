package main

import (
	"flag"
	"fmt"
	"github.com/bigpengry/crontab/master"
	"runtime"
)

var(
	err error
	configPath string
)

func initEnv()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func initArgs(){
	flag.StringVar(&configPath,"config","./master.json","指定master.json")
	flag.Parse()
}

func main() {

	//初始化配置文件
	initArgs()
	//初始化线程
	initEnv()

	//加载配置
	if err=master.InitConfig(configPath);err!=nil {
		goto ERR
	}

	//任务管理器
	if err=master.InitJobManager();err!=nil {
		goto ERR
	}

	//启动HTTP服务
	if err=master.InitAPIServer();err!=nil {
		goto ERR
	}

	return

	ERR:
		fmt.Println(err)
}
