package master

import (
	"encoding/json"
	"io/ioutil"
)

var (
	G_config *Config
)

type Config struct {
	APIPort         int      `json:"APIPort"`
	APIReadTimeOut  int      `json:"APIReadTimeOut"`
	APIWriteTimeOut int      `json:"APIWriteTimeOut"`
	ETCDEndPoints   []string `json:"etcdEndPoints"`
	ETCDDialTimeOut int      `json:"etcdDialTimeOut"`
}

//加载配置文件
func InitConfig(filename string) (err error) {
	//1.读取文件
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	//2.JSON反序列化
	config := new(Config)
	if err = json.Unmarshal(content, config); err != nil {
		return
	}

	G_config = config
	return
}
