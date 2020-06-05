package go_cron

import (
	"encoding/json"
	"github.com/codingXiang/go-logger"
	"github.com/codingXiang/go-orm"
	cronV3 "github.com/robfig/cron/v3"
	"time"
)

const (
	Cron = "cron_"
)

type GoCronInterface interface {
	Start()
	Stop()
	GetCore() *cronV3.Cron
	AddScheduler(s SchedulerInterface) (*RedisData, error)
	RemoveScheduler(s SchedulerInterface) error
}

type (
	GoCron struct {
		redis    orm.RedisClientInterface
		missions *mission
		core     *cronV3.Cron
	}
	RedisData struct {
		ID          cronV3.EntryID `json:"id"`
		Scheduler   *Scheduler     `json:"scheduler"`
		PreExecTime string         `json:"PreExecTime"`
	}
)

func StartSchedulerListener(spec string, s Service) {
	logger.Log.Info("Start scheduler listener, check frequency =", spec)
	s.Listen(spec)
}

func NewGoCron(redis orm.RedisClientInterface, missions *mission, opts ...cronV3.Option) GoCronInterface {
	c := &GoCron{
		redis:    redis,
		missions: missions,
		core:     cronV3.New(opts...),
	}
	return c
}

//GetCore 取得核心
func (g *GoCron) GetCore() *cronV3.Cron {
	return g.core
}

//Start 背景啟動排程
func (g *GoCron) Start() {
	g.core.Start()
}

//Stop 結束排程
func (g *GoCron) Stop() {
	g.core.Stop()
}

//CheckCronRecordIsExist 檢查排程是否存在（用於多 instance)
func (g *GoCron) CheckCronRecordIsExist(s SchedulerInterface) (*RedisData, error) {
	key := Cron + s.GetTaskName()
	//檢查 redis 中是否有相關的 key
	val, err := g.redis.GetValue(key)
	if err == nil {
		var (
			data = new(RedisData)
		)
		e := json.Unmarshal([]byte(val), data)
		//轉換格式
		if e != nil {
			logger.Log.Error(e.Error())
			return nil, e
		}
		logger.Log.Debug("data exist in redis")
		return data, nil
	} else {
		logger.Log.Debug("data not in redis")
		return nil, err
	}
}

//RemoveCronRecord 移除排程紀錄（用於多 instance)
func (g *GoCron) RemoveCronRecord(s SchedulerInterface) error {
	key := Cron + s.GetTaskName()
	if err := g.redis.RemoveKey(key); err == nil {
		return nil
	} else {
		return err
	}
}

//AddCron 加入排程
func (g *GoCron) AddCronRecord(entry cronV3.Entry, s SchedulerInterface) error {
	key := Cron + s.GetTaskName()
	data := RedisData{
		ID:        entry.ID,
		Scheduler: s.(*Scheduler),
	}
	tmp, _ := json.Marshal(data)
	logger.Log.Debug("add cron record", string(tmp))
	duration := entry.Next.Sub(time.Now())

	if err := g.redis.SetKeyValue(key, string(tmp), duration*2); err == nil {
		return nil
	} else {
		return err
	}
}

//AddScheduler 新增排程
func (g *GoCron) AddScheduler(s SchedulerInterface) (*RedisData, error) {
	//判斷 scheduler 是否存在
	if data, err := g.CheckCronRecordIsExist(s); data == nil || err != nil { //沒有取得紀錄
		//如果排程狀態為刪除，則跳過
		if s.GetDeleteAt() != nil {
			return nil, nil
		}

		//透過 Scheduler 的 task name 取得 job
		if job, err := g.missions.GetJob(s.GetTaskName()); err == nil {
			if id, err := g.core.AddJob(s.GetSpec(), job); err == nil {
				entry := g.GetCore().Entry(id)
				if err := g.AddCronRecord(entry, s); err != nil {
					return nil, err
				} else {
					return nil, nil
				}
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		//時程更新
		if s.GetSpec() != data.Scheduler.GetSpec() || s.GetDeleteAt() != nil{
			key := Cron + data.Scheduler.TaskName
			if data, err := ParseCronRedisData(key); err == nil {
				g.core.Remove(data.ID)
				orm.RedisORM.SetKeyValue(key, nil, -5)
			}
		}
		if s.GetDeleteAt() != nil {
			orm.DatabaseORM.GetInstance().Unscoped().Delete(s.(*Scheduler))
		}
		return data, nil
	}
}

//RemoveScheduler 新增排程
func (g *GoCron) RemoveScheduler(s SchedulerInterface) error {
	//判斷 scheduler 是否存在
	if data, err := g.CheckCronRecordIsExist(s); data != nil || err == nil { //沒有取得紀錄
		g.core.Remove(data.ID)
		if err := g.RemoveCronRecord(s); err == nil {
			g.Start()
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}
