package go_cron

import (
	"time"
)

//Scheduler 排程
type Scheduler struct {
	ID        int       `json:"id" gorm:"primary_key;auto_increment"`
	Spec      string    `json:"spec" gorm:"not_null;comment:'週期'"`
	TaskName  string    `json:"taskName" gorm:"Column:task_name;not_null;comment:'任務名稱'"`
	CreateAt  time.Time `json:"createAt" sql:"default:CURRENT_TIMESTAMP;comment:'建立時間'"`
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

//SchedulerLog 排程紀錄
type SchedulerLog struct {
	ID       int       `json:"id" gorm:"primary_key;auto_increment"`
	TaskName string    `json:"taskName" gorm:"Column:task_name;not_null;comment:'任務名稱'"`
	HostName string    `json:"hostname" gorm:"comment:'host 名稱'"`
	Success  bool      `json:"isSuccess" gorm:"Column:isSuccess"`
	Message  string    `json:"message" gorm:"Column:message"`
	CreateAt time.Time `json:"createAt" sql:"default:CURRENT_TIMESTAMP;comment:'執行時間'"`
}

//此為自動產生的 Interface，建議不要進行更動
type SchedulerInterface interface {
	GetID() int
	SetID(in int) *Scheduler
	GetSpec() string
	SetSpec(in string) *Scheduler
	GetTaskName() string
	SetTaskName(in string) *Scheduler
	GetCreateAt() time.Time
	SetCreateAt(in time.Time) *Scheduler
	GetDeleteAt() *time.Time
}

//此為自動產生的 Interface，建議不要進行更動
type SchedulerLogInterface interface {
	GetID() int
	SetID(in int) *SchedulerLog
	IsSuccess() bool
	SetSuccess(in bool) *SchedulerLog
	GetCreateAt() time.Time
	SetCreateAt(in time.Time) *SchedulerLog
}

//此為自動產生的方法，建議不要更動
func NewSchedulerLog() SchedulerLogInterface {
	return &SchedulerLog{}
}

//此為自動產生的方法，建議不要更動
func NewScheduler() SchedulerInterface {
	return &Scheduler{}
}

//此為自動產生的方法，建議不要更動
func (g *Scheduler) GetID() int {

	return g.ID
}

//此為自動產生的方法，建議不要更動
func (g *Scheduler) SetID(in int) *Scheduler {
	g.ID = in
	return g
}

//此為自動產生的方法，建議不要更動
func (g *Scheduler) GetSpec() string {

	return g.Spec
}

//此為自動產生的方法，建議不要更動
func (g *Scheduler) SetSpec(in string) *Scheduler {
	g.Spec = in
	return g
}

//此為自動產生的方法，建議不要更動
func (g *Scheduler) GetTaskName() string {

	return g.TaskName
}

//此為自動產生的方法，建議不要更動
func (g *Scheduler) SetTaskName(in string) *Scheduler {
	g.TaskName = in
	return g
}

//此為自動產生的方法，建議不要更動
func (g *Scheduler) GetCreateAt() time.Time {

	return g.CreateAt
}

//此為自動產生的方法，建議不要更動
func (g *Scheduler) SetCreateAt(in time.Time) *Scheduler {
	g.CreateAt = in
	return g
}

//此為自動產生的方法，建議不要更動
func (g *SchedulerLog) GetID() int {

	return g.ID
}

//此為自動產生的方法，建議不要更動
func (g *SchedulerLog) SetID(in int) *SchedulerLog {
	g.ID = in
	return g
}

//此為自動產生的方法，建議不要更動
func (g *SchedulerLog) IsSuccess() bool {

	return g.Success
}

//此為自動產生的方法，建議不要更動
func (g *SchedulerLog) SetSuccess(in bool) *SchedulerLog {
	g.Success = in
	return g
}

//此為自動產生的方法，建議不要更動
func (g *SchedulerLog) GetCreateAt() time.Time {

	return g.CreateAt
}

//此為自動產生的方法，建議不要更動
func (g *SchedulerLog) SetCreateAt(in time.Time) *SchedulerLog {
	g.CreateAt = in
	return g
}

//此為自動產生的方法，建議不要更動
func (g *Scheduler) GetDeleteAt() *time.Time {

	return g.DeletedAt
}
