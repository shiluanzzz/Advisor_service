package cronjob

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"service/model"
	"service/service"
	"service/utils/errmsg"
	"service/utils/logger"
	"time"
)

// CronJob 定时job主体
type CronJob struct {
	OrderId  int64
	UserId   int64
	RushTime int64
	CronId   int
}

// Run 定时任务的运行逻辑
func (order *CronJob) Run() {
	now := time.Now().Unix()
	if now-order.RushTime > 60 {
		code := service.ChangeOrderStatus(order.OrderId, order.UserId, model.Rush, model.Pending)
		// 这个订单被顾问回答了 或者执行成功了
		if code == errmsg.ErrorOrderHasCompleted || code == errmsg.SUCCESS {
			CloseJob(order)
		} else {
			logger.Log.Error("定时任务执行失败,", zap.String("errmsg", errmsg.GetErrMsg(code)))
		}
	}
}

var closeJobChan chan CronJob
var jobMap map[int]*CronJob
var C *cron.Cron

func init() {
	closeJobChan = make(chan CronJob, 100)
	jobMap = make(map[int]*CronJob)
	C = cron.New()
	C.Start()
	// 单开一个协程专门用来关闭不需要的定时任务
	go closeTask(C, closeJobChan, &jobMap)

}
func AddJob(job *CronJob) int {
	jobId, err := C.AddJob("@every 1m", job)
	if err != nil {
		logger.Log.Error("创建定时任务失败", zap.Error(err))
		return errmsg.ErrorCronAddJob
	} else {
		job.CronId = int(jobId)
		jobMap[int(jobId)] = job
		logger.Log.Info("创建定时任务成功", zap.Int("cronId", job.CronId))
		return errmsg.SUCCESS
	}
}
func CloseJob(job *CronJob) {
	// 只负责结束job，状态的查询放在service的事务里去做
	closeJobChan <- *job
}
func closeTask(c *cron.Cron, closeChan chan CronJob, jobMap *map[int]*CronJob) {
	for job := range closeChan {
		logger.Log.Info("关闭定时监控任务",
			zap.Int("cronId", job.CronId),
		)
		c.Remove(cron.EntryID(job.CronId))
		delete(*jobMap, job.CronId)
	}
}

// 系统异常退出 保存job TODO
