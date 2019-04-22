package master

import (
	"context"
	"time"

	"github.com/bigpengry/crontab/common"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// LogManager 日志管理
type LogManager struct {
	client        *mongo.Client
	logCollection *mongo.Collection
}

// LogMgr
var (
	LogMgr *LogManager
)

// ListLog 列出任务日志
func (m *LogManager) ListLog(name string, skip, limit int) (logArr []common.TaskLog, err error) {
	filter := &common.TaskFilter{JobName: name}
	logSort := &common.SortLogByStartTime{SortOrder: -1}
	taskLog := new(common.TaskLog)
	cursor, err := LogMgr.logCollection.Find(context.TODO(), filter,
		options.Find().SetSort(logSort),
		options.Find().SetSkip(int64(skip)),
		options.Find().SetLimit(int64(limit)))
	if err != nil {
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		if err = cursor.Decode(taskLog); err != nil {
			continue
		}
		logArr = append(logArr, *taskLog)
	}
	return
}

// InitLogManager 日志管理器
func InitLogManager() (err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(Conf.MongoDBConnectTimeOut)*time.Millisecond)
	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(Conf.MongoDBURI))
	if err != nil {
		return
	}
	LogMgr = &LogManager{
		client:        cli,
		logCollection: cli.Database("cron").Collection("log"),
	}
	return
}
