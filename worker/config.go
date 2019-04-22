package worker

import (
	"encoding/json"
	"io/ioutil"
)

// Conf 是一个Config对象的全局化实例
var (
	Conf *Config
)

// Config 配置文件对象
type Config struct {
	ETCDEndPoints         []string `json:"etcdEndPoints"`
	ETCDDialTimeOut       int      `json:"etcdDialTimeOut"`
	MongoDBURI            string   `json:"mongodbURI“`
	MongoDBConnectTimeOut int      `json:"mongodbConnectTimeOut"`
	LogBufferSize         int      `json:"logBufferSize"`
	LogCommitTimeOut      int      `json:"logCommitTimeOut"`
}

// InitConfig 加载配置文件
func InitConfig(filename string) (err error) {
	conf := new(Config)
	// 读取文件
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	// JSON反序列化
	if err = json.Unmarshal(content, conf); err != nil {
		return
	}

	Conf = conf
	return
}
