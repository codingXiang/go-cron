package go_cron

import (
	"github.com/astaxie/beego/validation"
	"github.com/codingXiang/cxgateway/delivery"
	"github.com/codingXiang/cxgateway/pkg/e"
	"github.com/codingXiang/cxgateway/pkg/i18n"
	"github.com/codingXiang/cxgateway/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
	"strconv"
)

type Repository interface {
	GetSchedulerList(data map[string]interface{}) ([]*Scheduler, error)
	GetScheduler(data SchedulerInterface) (*Scheduler, error)
	CreateScheduler(data SchedulerInterface) (*Scheduler, error)
	UpdateScheduler(data SchedulerInterface) (*Scheduler, error)
	ModifyScheduler(m SchedulerInterface, data map[string]interface{}) (*Scheduler, error)
	DeleteScheduler(data SchedulerInterface) error
}

type SchedulerRepository struct {
	orm *gorm.DB
}

func NewSchedulerRepository(orm *gorm.DB) Repository {
	return &SchedulerRepository{orm: orm}
}

func (s *SchedulerRepository) GetSchedulerList(data map[string]interface{}) ([]*Scheduler, error) {
	var (
		err error
		in  = make([]*Scheduler, 0)
	)

	err = s.orm.Find(&in, data).Error
	return in, err
}

func (s *SchedulerRepository) GetScheduler(data SchedulerInterface) (*Scheduler, error) {
	var (
		err error
		in  = data.(*Scheduler)
	)
	err = s.orm.First(&in).Error
	return in, err
}

func (s *SchedulerRepository) CreateScheduler(data SchedulerInterface) (*Scheduler, error) {
	var (
		err error
		in  = data.(*Scheduler)
	)
	err = s.orm.Create(&in).Error
	return in, err
}

func (s *SchedulerRepository) UpdateScheduler(data SchedulerInterface) (*Scheduler, error) {
	var (
		err error
		in  = data.(*Scheduler)
	)
	err = s.orm.Save(&in).Error
	return in, err
}

func (s *SchedulerRepository) ModifyScheduler(m SchedulerInterface, data map[string]interface{}) (*Scheduler, error) {
	var (
		err error
		in  = m.(*Scheduler)
	)
	err = s.orm.Model(&in).Updates(data).Error
	return in, err
}

func (s *SchedulerRepository) DeleteScheduler(data SchedulerInterface) error {
	var (
		err error
		in  = data.(*Scheduler)
	)
	err = s.orm.Delete(&in).Error
	return err
}

type Service interface {
	GetSchedulerList(data map[string]interface{}) ([]*Scheduler, error)
	GetScheduler(data SchedulerInterface) (*Scheduler, error)
	CreateScheduler(data SchedulerInterface) (*Scheduler, error)
	UpdateScheduler(data SchedulerInterface) (*Scheduler, error)
	ModifyScheduler(m SchedulerInterface, data map[string]interface{}) (*Scheduler, error)
	DeleteScheduler(data SchedulerInterface) error
}

type SchedulerService struct {
	repo Repository
	core GoCronInterface
}

func NewSchedulerService(core GoCronInterface, repo Repository) Service {
	return &SchedulerService{core: core, repo: repo}
}

func (s *SchedulerService) GetSchedulerList(data map[string]interface{}) ([]*Scheduler, error) {
	return s.repo.GetSchedulerList(data)
}

func (s *SchedulerService) GetScheduler(data SchedulerInterface) (*Scheduler, error) {
	return s.repo.GetScheduler(data)
}

func (s *SchedulerService) CreateScheduler(data SchedulerInterface) (*Scheduler, error) {
	if scheduler, err := s.repo.CreateScheduler(data); err == nil {
		err = s.core.AddScheduler(scheduler)
		s.core.Run()
		return scheduler, err
	} else {
		return nil, err
	}
}

func (s *SchedulerService) UpdateScheduler(data SchedulerInterface) (*Scheduler, error) {
	return s.repo.UpdateScheduler(data)
}

func (s *SchedulerService) ModifyScheduler(m SchedulerInterface, data map[string]interface{}) (*Scheduler, error) {
	return s.repo.ModifyScheduler(m, data)
}

func (s *SchedulerService) DeleteScheduler(data SchedulerInterface) error {
	return s.repo.DeleteScheduler(data)
}

type HttpHandler interface {
	GetSchedulerList(c *gin.Context) error
	GetScheduler(c *gin.Context) error
	CreateScheduler(c *gin.Context) error
	UpdateScheduler(c *gin.Context) error
	ModifyScheduler(c *gin.Context) error
	DeleteScheduler(c *gin.Context) error
}

const (
	MODULE = "cron"
)

type SchedulerHttpHandler struct {
	i18nMsg i18n.I18nMessageHandlerInterface
	gateway delivery.HttpHandler
	svc     Service
}

func NewSchedulerHttpHandler(gateway delivery.HttpHandler, svc Service) HttpHandler {
	var handler = &SchedulerHttpHandler{
		i18nMsg: i18n.NewI18nMessageHandler(MODULE),
		gateway: gateway,
		svc:     svc,
	}
	v1 := gateway.GetApiRoute().Group("/v1/cron")
	v1.GET("", e.Wrapper(handler.GetSchedulerList))
	v1.GET("/:id", e.Wrapper(handler.GetScheduler))
	v1.POST("", e.Wrapper(handler.CreateScheduler))
	v1.PUT("/:id", e.Wrapper(handler.UpdateScheduler))
	v1.PATCH("/:id", e.Wrapper(handler.ModifyScheduler))
	v1.DELETE("/:id", e.Wrapper(handler.DeleteScheduler))
	return handler
}

func (g *SchedulerHttpHandler) GetSchedulerList(c *gin.Context) error {
	var (
		data = map[string]interface{}{}
	)
	g.i18nMsg.SetCore(util.GetI18nData(c))

	//抓取 query string
	if in, isExist := c.GetQuery("spec"); isExist {
		data["spec"] = in
	}

	if result, err := g.svc.GetSchedulerList(data); err != nil {
		return g.i18nMsg.GetError(err)
	} else {
		c.JSON(g.i18nMsg.GetSuccess(result))
		return nil
	}
}

func (g *SchedulerHttpHandler) GetScheduler(c *gin.Context) error {
	var (
		data  = new(Scheduler)
		strId = c.Params.ByName("id")
	)
	//將 middleware 傳入的 i18n 進行轉換
	g.i18nMsg.SetCore(util.GetI18nData(c))

	if id, err := strconv.Atoi(strId); err == nil {
		data.SetID(id)
	} else {
		g.i18nMsg.ParameterIntError(strId, err)
	}

	if result, err := g.svc.GetScheduler(data); err != nil {
		return g.i18nMsg.GetError(err)
	} else {
		c.JSON(g.i18nMsg.GetSuccess(result))
	}
	return nil
}

func (g *SchedulerHttpHandler) CreateScheduler(c *gin.Context) error {
	var (
		valid = new(validation.Validation)
		data  = new(Scheduler)
	)
	//將 middleware 傳入的 i18n 進行轉換
	g.i18nMsg.SetCore(util.GetI18nData(c))
	//綁定參數
	var err = c.ShouldBindWith(&data, binding.JSON)
	if err != nil || data == nil {
		return g.i18nMsg.ParameterFormatError()
	}

	//驗證表單資訊是否填寫充足
	valid.Required(&data.TaskName, "taskName")
	valid.Required(&data.Spec, "spec")

	if err := util.NewRequestHandler().ValidValidation(valid); err != nil {
		return err
	}

	if result, err := g.svc.CreateScheduler(data); err != nil {
		return g.i18nMsg.CreateError(err)
	} else {
		c.JSON(g.i18nMsg.CreateSuccess(result))
		return nil
	}
}

func (g *SchedulerHttpHandler) UpdateScheduler(c *gin.Context) error {
	var (
		data  = new(Scheduler)
		err   error
		strId = c.Params.ByName("id")
	)
	//將 middleware 傳入的 i18n 進行轉換
	g.i18nMsg.SetCore(util.GetI18nData(c))

	if id, err := strconv.Atoi(strId); err == nil {
		data.SetID(id)
	} else {
		g.i18nMsg.ParameterIntError(strId, err)
	}
	//取得 tenant
	if data, err = g.svc.GetScheduler(data); err != nil {
		return g.i18nMsg.GetError(err)
	}

	//綁定參數
	err = c.ShouldBindWith(data, binding.JSON)
	if err != nil || data == nil {
		return g.i18nMsg.ParameterFormatError()
	}

	//更新 tenant
	if result, err := g.svc.UpdateScheduler(data); err != nil {
		return g.i18nMsg.UpdateError(err)
	} else {
		c.JSON(g.i18nMsg.UpdateSuccess(result))
		return nil
	}
}

func (g *SchedulerHttpHandler) ModifyScheduler(c *gin.Context) error {
	var (
		data       = new(Scheduler)
		updateData = new(map[string]interface{})
		strId      = c.Params.ByName("id")
	)
	//將 middleware 傳入的 i18n 進行轉換
	g.i18nMsg.SetCore(util.GetI18nData(c))

	if id, err := strconv.Atoi(strId); err == nil {
		data.SetID(id)
	} else {
		g.i18nMsg.ParameterIntError(strId, err)
	}

	//綁定參數
	err := c.ShouldBindWith(&updateData, binding.JSON)
	if err != nil || data == nil {
		return g.i18nMsg.ParameterFormatError()
	}

	if result, err := g.svc.ModifyScheduler(data, *updateData); err != nil {
		return g.i18nMsg.ModifyError(err)
	} else {
		c.JSON(g.i18nMsg.ModifySuccess(result))
		return nil
	}
}

func (g *SchedulerHttpHandler) DeleteScheduler(c *gin.Context) error {
	var (
		data  = new(Scheduler)
		strId = c.Params.ByName("id")
	)
	//將 middleware 傳入的 i18n 進行轉換
	g.i18nMsg.SetCore(util.GetI18nData(c))

	if id, err := strconv.Atoi(strId); err == nil {
		data.SetID(id)
	} else {
		g.i18nMsg.ParameterIntError(strId, err)
	}

	if err := g.svc.DeleteScheduler(data); err != nil {
		return g.i18nMsg.DeleteError(err)
	} else {
		c.JSON(g.i18nMsg.DeleteSuccess(nil))
		return nil
	}
}
