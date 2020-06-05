package go_cron

import (
	"errors"
)

type mission struct {
	content map[string]BasicJobInterface
}

var Missions *mission

func NewMission() *mission {
	return &mission{
		content: make(map[string]BasicJobInterface),
	}
}

//AddMission 擴充任務
func (m *mission) AddMission(job BasicJobInterface) *mission {
	m.content[job.GetName()] = job
	return m
}

//GetJob 取得 job
func (m *mission) GetJob(name string) (BasicJobInterface, error) {
	if out := m.content[name]; out != nil {
		return out, nil
	} else {
		return nil, errors.New("task name " + name + " is not binding job")
	}
}
