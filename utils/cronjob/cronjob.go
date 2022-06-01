package cronjob

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"service/model"
	"service/service"
	"service/utils"
	"service/utils/errmsg"
	"service/utils/logger"
	"time"
)

// CronJob 定时job主体
type CronJob struct {
	OrderId    int64
	UserId     int64
	RushTime   int64
	CreateTime int64
	CronId     int
	// CronType=0 rush->pending =1:pending->expired
	CronType int
}

const (
	RushOrderType = iota
	PendingOrderType
)

// Run 定时任务的运行逻辑
func (order *CronJob) Run() {
	now := time.Now().Unix()
	// 隔1分钟扫描一次
	switch order.CronType {
	case RushOrderType:
		if now-order.RushTime > utils.RushOrder2PendingTime*60 {
			code := service.ChangeOrderStatus(order.OrderId, order.UserId, model.Rush, model.Pending)
			// 这个订单被顾问回答了 或者执行成功了
			if code == errmsg.ErrorOrderHasCompleted || code == errmsg.SUCCESS {
				CloseJob(order)
			} else {
				logger.Log.Error("定时任务0执行失败,", zap.String("errmsg", errmsg.GetErrMsg(code)))
			}
		}
	case PendingOrderType:
		if now-order.CreateTime > 60*60*24 {
			code := service.ChangeOrderStatus(order.OrderId, order.UserId, model.Pending, model.Expired)
			if code == errmsg.SUCCESS {
				CloseJob(order)
			} else {
				logger.Log.Error("定时任务1执行失败,", zap.String("errmsg", errmsg.GetErrMsg(code)))
			}
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
	logMsg := fmt.Sprintf("创建定时任务%d成功", job.CronType)
	if err != nil {
		logger.Log.Error(logMsg, zap.Error(err))
		return errmsg.ErrorCronAddJob
	} else {
		job.CronId = int(jobId)
		jobMap[int(jobId)] = job
		logger.Log.Info(logMsg, zap.Int("cronId", job.CronId))
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
