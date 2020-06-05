package go_cron

import (
	"encoding/json"
	"errors"
	"github.com/codingXiang/go-logger"
	"github.com/codingXiang/go-orm"
	"os"
	"time"
)

type BasicJobInterface interface {
	SetCore(core GoCronInterface)
	GetCore() GoCronInterface
	SetName(name string)
	GetName() string
	SetHostName(hostName string)
	GetHostName() string
	SetSvc(svc Service)
	GetSvc() Service
	UpdateRedisData(errMsg error)
	Run()
}

type BasicJob struct {
	Name     string
	HostName string
	Core     GoCronInterface
	Svc      Service
}

func (b *BasicJob) SetSvc(svc Service) {
	b.Svc = svc
}

func (b *BasicJob) GetSvc() Service {
	return b.Svc
}

func (b *BasicJob) SetCore(core GoCronInterface) {
	b.Core = core
}

func (b *BasicJob) GetCore() GoCronInterface {
	return b.Core
}

func (b *BasicJob) SetName(name string) {
	b.Name = name
}

func (b *BasicJob) SetHostName(hostName string) {
	b.HostName = hostName
}

func (b *BasicJob) GetName() string {
	return b.Name
}

func (b *BasicJob) GetHostName() string {
	return b.HostName
}

func (b *BasicJob) Run() {}

func (b *BasicJob) UpdateRedisData(errMsg error) {
	key := "cron_" + b.GetName()
	hostname, err := os.Hostname()
	isSuccess := true
	if err != nil {
		isSuccess = false
		logger.Log.Error(err.Error())
	}
	b.SetHostName(hostname)
	if data, err := ParseCronRedisData(key); err == nil {
		SetCronRedisData(key, b.GetCore(), data)
		//上傳 log
		if errMsg == nil {
			errMsg = errors.New("")
		}
		b.GetSvc().CreateSchedulerLog(&SchedulerLog{
			TaskName: data.Scheduler.TaskName,
			HostName: hostname,
			Success:  isSuccess,
			Message:  errMsg.Error(),
		})
	}
}

func ParseCronRedisData(key string) (data *RedisData, err error) {
	var (
		val string
	)
	val, err = orm.RedisORM.GetValue(key)
	//從 redis 中取得資料發生錯誤
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(val), &data)
	//轉換格式
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	return
}

func SetCronRedisData(key string, core GoCronInterface, data *RedisData) (err error) {
	var (
		val []byte
	)
	//轉換 redis data 為 byte
	val, err = json.Marshal(&data)
	if err != nil {
		return
	}
	entry := core.GetCore().Entry(data.ID)

	duration := entry.Next.Sub(time.Now())
	//存入資料
	err = orm.RedisORM.SetKeyValue(key, string(val), duration*3)
	if err != nil {
		logger.Log.Error("err", err.Error())
		return
	}
	logger.Log.Debug("set", key, " more", duration.String())
	return
}

func AddCore(core GoCronInterface, jobs ...BasicJobInterface) {
	for _, job := range jobs {
		job.SetCore(core)
	}
}

func AddSvc(svc Service, jobs ...BasicJobInterface) {
	for _, job := range jobs {
		job.SetSvc(svc)
	}
}
