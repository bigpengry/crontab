package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/bigpengry/crontab/worker"
)

var (
	err      error
	confPath string
)

// initEnv 初始化线程数
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// initArgs 解析命令行参数
func initArgs() {
	// worker -config ./worker.json
	// worker -h
	flag.StringVar(&confPath, "config", "./worker.json", "指定worker.json")
	flag.Parse()
}

func main() {

	// 初始化配置文件
	initArgs()
	// 初始化线程
	initEnv()
	// 加载配置
	if err = worker.InitConfig(confPath); err != nil {
		fmt.Println(err)
		return
	}

	if err = worker.InitScheduler(); err != nil {
		fmt.Println(err)
		return
	}

	if err = worker.InitLogSink(); err != nil {
		fmt.Println(err)
		return
	}

	if err = worker.InitJobManager(); err != nil {
		fmt.Println(err)
		return
	}

	if err = worker.InitResister(); err != nil {
		fmt.Println(err)
		return
	}
	for {
		time.Sleep(1 * time.Second)
	}

}
