package main

import (
	"fmt"
	"github.com/codingXiang/configer"
	go_cron "github.com/codingXiang/go-cron"
	"github.com/codingXiang/go-logger"
	"github.com/codingXiang/go-orm"
	"github.com/robfig/cron/v3"
)

func init() {
	logger.Log = logger.NewLogger(logger.Logger{
		Level:  "info",
		Format: "text",
	})
	//初始化資料庫
	db := configer.NewConfigerCore("yaml", "database", "./config", ".")
	db.SetAutomaticEnv("")
	//設定資料庫
	if o, err := orm.NewOrm("database", db); err == nil {
		orm.DatabaseORM = o
		// 建立 Table Schema (Module)
		logger.Log.Debug("create table")
		{
			//設定排程
			_ = orm.DatabaseORM.CheckTable(true, go_cron.Scheduler{})
			_ = orm.DatabaseORM.CheckTable(true, go_cron.SchedulerLog{})
		}
	} else {
		logger.Log.Error(err.Error())
		panic(err.Error())
	}
	//初始化 redis 參數
	redis := configer.NewConfigerCore("yaml", "redis", "./config", "./example")
	redis.SetAutomaticEnv("")
	//設定 redis
	if data, err := orm.NewRedisClient("redis", redis); err == nil {
		orm.RedisORM = data
	} else {
		logger.Log.Error(err.Error())
		panic(err.Error())
	}
}

func main() {
	test := &Test{}
	test.SetName("test")

	test1 := &Test{}
	test1.SetName("test1")

	missions := go_cron.NewMission()
	missions.AddMission(test)
	missions.AddMission(test1)
	core := go_cron.NewGoCron(orm.RedisORM, missions, cron.WithSeconds())

	schedulerRepo := go_cron.NewSchedulerRepository(orm.DatabaseORM.GetInstance())
	schedulerSvc := go_cron.NewSchedulerService(core, schedulerRepo)

	//加入排程核心
	go_cron.AddCore(core, test, test1)
	//加入 Service
	go_cron.AddSvc(schedulerSvc, test, test1)
	//開始監聽
	go_cron.StartSchedulerListener("* * * * * *", schedulerSvc)
	schedulerSvc.CreateScheduler(&go_cron.Scheduler{TaskName: "test", Spec: "* * * * * *"})

	select {}
}

type Test struct {
	go_cron.BasicJob
}

func (t *Test) Run() {
	fmt.Println(t.GetHostName() + " " + t.GetName())
	t.UpdateRedisData(nil)
}
