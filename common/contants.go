package common

const (
	// JOB_KILLER_DIR 文件强杀路径
	JOB_KILLER_DIR = "/cron/killer/"
	// JOB_SAVE_DIR 文件存放路径
	JOB_SAVE_DIR = "/cron/jobs/"
	// JOB_LOCK_DIR 文件存放路径
	JOB_LOCK_DIR = "/cron/lock/"
	// JOB_WORKER_DIR 文件存放路径
	JOB_WORKER_DIR = "/cron/worker"
	// JOB_EVENT_SAVE 任务事件类型：保存
	JOB_EVENT_SAVE = 1
	// JOB_EVENT_DELETE 任务事件类型：删除
	JOB_EVENT_DELETE = 2
	// JOB_EVENT_KILL 任务事件类型：强杀
	JOB_EVENT_KILL = 3
)
