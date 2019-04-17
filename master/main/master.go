package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/bigpengry/crontab/master"
)

var (
	err        error
	configPath string
)

// initEnv 初始化线程数
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// initArgs 加载配置文件路径
func initArgs() {
	flag.StringVar(&configPath, "config", "./master.json", "指定master.json")
	flag.Parse()
}

func main() {

	// 初始化配置文件
	initArgs()
	// 初始化线程
	initEnv()
	// 加载配置
	if err = master.InitConfig(configPath); err != nil {
		fmt.Println(err)
		return
	}

	// 初始化任务管理器
	if err = master.InitJobManager(); err != nil {
		fmt.Println(err)
		return
	}

	// 启动HTTP服务
	if err = master.InitAPIServer(); err != nil {
		fmt.Println(err)
		return
	}

	for {
		time.Sleep(1 * time.Second)
	}

}
