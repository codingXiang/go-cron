package go_cron

import (
	"errors"
	"github.com/codingXiang/go-orm"
	cronV3 "github.com/robfig/cron/v3"
	"strconv"
)

const (
	cron = "cron_"
)

type GoCronInterface interface {
	Run()
	Stop()
	AddScheduler(s *Scheduler) error
	RemoveScheduler(s *Scheduler) error
}

type GoCron struct {
	redis    orm.RedisClientInterface
	missions *mission
	core     *cronV3.Cron
}

func NewGoCron(redis orm.RedisClientInterface, opts ...cronV3.Option) GoCronInterface {
	c := &GoCron{
		redis: redis,
		core:  cronV3.New(opts...),
	}
	return c
}

//Run 啟動排程
func (g *GoCron) Run() {
	g.core.Run()
}

//Stop 結束排程
func (g *GoCron) Stop() {
	g.core.Stop()
}

//CheckCron 檢查排程是否存在（用於多 instance)
func (g *GoCron) CheckCronRecord(s SchedulerInterface) (int, error) {
	key := cron + s.GetTaskName()
	if val, err := g.redis.GetValue(key); err == nil && val != "" {
		return strconv.Atoi(val)
	} else {
		return 0, err
	}
}

//RemoveCronRecord 移除排程紀錄（用於多 instance)
func (g *GoCron) RemoveCronRecord(s SchedulerInterface) error {
	key := cron + s.GetTaskName()
	if err := g.redis.RemoveKey(key); err == nil {
		return nil
	} else {
		return err
	}
}

//AddCron 加入排程
func (g *GoCron) AddCronRecord(id cronV3.EntryID, s SchedulerInterface) error {
	key := cron + s.GetTaskName()
	if err := g.redis.SetKeyValue(key, id, 0); err == nil {
		return nil
	} else {
		return err
	}
}

//AddScheduler 新增排程
func (g *GoCron) AddScheduler(s *Scheduler) error {
	//判斷 scheduler 是否存在
	if id, err := g.CheckCronRecord(s); id == 0 || err != nil { //沒有取得紀錄
		//透過 Scheduler 的 task name 取得 job
		if job, err := g.missions.GetJob(s.GetTaskName()); err == nil {
			if id, err := g.core.AddJob(s.GetSpec(), *job); err == nil {
				if err := g.AddCronRecord(id, s); err == nil {
					return nil
				} else {
					return err
				}
			} else {
				return err
			}
		} else {
			return err
		}
	} else {
		return errors.New("cron record is exist")
	}
}

//RemoveScheduler 新增排程
func (g *GoCron) RemoveScheduler(s *Scheduler) error {
	//判斷 scheduler 是否存在
	if id, err := g.CheckCronRecord(s); id != 0 || err == nil { //沒有取得紀錄
		if err := g.RemoveCronRecord(s); err == nil {
			g.Run()
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}
