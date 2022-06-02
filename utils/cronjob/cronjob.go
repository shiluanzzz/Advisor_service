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

const (
	RushOrderType = iota
	PendingOrderType
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

// Run 定时任务的运行逻辑
func (order *CronJob) Run() {
	now := time.Now()
	// 隔1分钟扫描一次
	switch order.CronType {
	case RushOrderType:
		runTime := time.Unix(order.RushTime+utils.RushOrder2PendingTime*60, 0)
		if runTime.After(now) {
			// service层
			code := service.ChangeOrderStatus(order.OrderId, order.UserId, model.Rush, model.Pending)
			// 这个订单被顾问回答了 或者执行成功了
			if code == errmsg.ErrorOrderHasCompleted || code == errmsg.SUCCESS {
				CloseJob(order)
			} else {
				logger.Log.Error("定时任务0执行失败,", zap.String("errorMsg", errmsg.GetErrMsg(code)))
			}
		}
	case PendingOrderType:
		runTime := time.Unix(order.CreateTime+utils.PendingOrder2ExpireTime*60, 0)
		if runTime.After(now) {
			// service层
			code := service.ChangeOrderStatus(order.OrderId, order.UserId, model.Pending, model.Expired)
			if code == errmsg.SUCCESS {
				CloseJob(order)
			} else {
				logger.Log.Error("定时任务1执行失败,", zap.String("errorMsg", errmsg.GetErrMsg(code)))
			}
		}
	}
}

var closeJobChan chan CronJob
var jobMap map[int]*CronJob
var C *cron.Cron

func InitCronJob() {
	closeJobChan = make(chan CronJob, 100)
	jobMap = make(map[int]*CronJob)
	C = cron.New()
	C.Start()
	// 单开一个协程专门用来关闭不需要的定时任务
	go closeTask(C, closeJobChan, &jobMap)
	// 恢复数据库中可能存在的任务
	go recoverJobs()

}

// AddJob 往cron中添加定时任务
func AddJob(job *CronJob) int {
	jobId, err := C.AddJob("@every 1m", job)
	logMsg := fmt.Sprintf("创建定时任务%d", job.CronType)
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

// CloseJob 关闭cron中的定时任务
func CloseJob(job *CronJob) {
	// 只负责结束job，状态的查询放在service的事务里去做
	closeJobChan <- *job
}

// closeTask 单开一个协程去关闭掉cron中的任务
func closeTask(c *cron.Cron, closeChan chan CronJob, jobMap *map[int]*CronJob) {
	for job := range closeChan {
		logger.Log.Info("关闭定时监控任务",
			zap.Int("cronId", job.CronId),
		)
		c.Remove(cron.EntryID(job.CronId))
		delete(*jobMap, job.CronId)
	}
}

// 系统异常退出 从数据库重新读取表 看有没有需要重新导入表中的写法
func recoverJobs() {
	code, res := service.GetManyTableItemsByWhere(service.ORDERTABLE,
		map[string]interface{}{"status": model.Pending},
		[]string{"id", "user_id", "create_time"},
	)
	if code == errmsg.SUCCESS {
		if len(res) != 0 {
			logger.Log.Info(fmt.Sprintf("从数据库中查询到%d条需要定时处理的普通订单", len(res)))
		}
		for _, v := range res {
			job := CronJob{
				OrderId:    v["id"].(int64),
				UserId:     v["user_id"].(int64),
				CreateTime: v["create_time"].(int64),
				CronType:   PendingOrderType,
			}
			_ = AddJob(&job)
		}
	}

	code, res = service.GetManyTableItemsByWhere(service.ORDERTABLE,
		map[string]interface{}{"status": model.Rush},
		[]string{"id", "user_id", "rush_time"},
	)
	if code == errmsg.SUCCESS {
		if len(res) != 0 {
			logger.Log.Info(fmt.Sprintf("从数据库中查询到%d条需要定时处理的加急订单", len(res)))
		}
		for _, v := range res {
			job := CronJob{
				OrderId:    v["id"].(int64),
				UserId:     v["user_id"].(int64),
				CreateTime: v["rush_time"].(int64),
				CronType:   RushOrderType,
			}
			_ = AddJob(&job)
		}
	}
}
