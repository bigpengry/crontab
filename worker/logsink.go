package worker

import (
	"context"
	"time"

	"github.com/bigpengry/crontab/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// LogSink mongodb存储日志
type LogSink struct {
	client        *mongo.Client
	logCollection *mongo.Collection
	logChan       chan *common.TaskLog
}

// Log 日志处理
var (
	Log *LogSink
)

// write 写入日志
func (s *LogSink) write() {
	log := new(common.TaskLog)
	buffer := make([]interface{}, 0)
	timer := time.NewTimer(time.Duration(Conf.LogCommitTimeOut) * time.Millisecond)
	for {
		select {
		case log = <-s.logChan:

			buffer = append(buffer, log)
			if len(buffer) >= Conf.LogBufferSize {
				s.logCollection.InsertMany(context.TODO(), buffer)
				buffer = make([]interface{}, 0)
				timer.Reset(time.Duration(Conf.LogCommitTimeOut) * time.Millisecond)
			}
		case <-timer.C:
			s.logCollection.InsertMany(context.TODO(), buffer)
			buffer = make([]interface{}, 0)
			timer.Reset(time.Duration(Conf.LogCommitTimeOut) * time.Millisecond)
		}
	}
}

// Append 发送数据，如果chan溢出，则直接丢弃
func (s *LogSink) Append(log *common.TaskLog) {
	select {
	case s.logChan <- log:
	default:
	}
}

// InitLogSink 初始化日志
func InitLogSink() (err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(Conf.MongoDBConnectTimeOut)*time.Millisecond)
	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(Conf.MongoDBURI))
	if err!=nil{
		return
	}
	Log = &LogSink{
		client:        cli,
		logCollection: cli.Database("cron").Collection("log"),
		logChan:       make(chan *common.TaskLog, 1000),
	}

	go Log.write()
	return
}
