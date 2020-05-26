package go_cron

import (
	"errors"
	cronV3 "github.com/robfig/cron/v3"
)

type mission struct {
	content map[string]cronV3.Job
}

var Missions *mission

func init() {
	Missions = new(mission)
}

//AddMission 擴充任務
func (m *mission) AddMission(name string, job cronV3.Job) *mission {
	m.content[name] = job
	return m
}

//GetJob 取得 job
func (m *mission) GetJob(name string) (cronV3.Job, error) {
	if out := m.content[name]; out != nil {
		return out, nil
	} else {
		return nil, errors.New("task name " + name + " is not binding job")
	}
}
