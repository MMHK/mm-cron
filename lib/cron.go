package lib

import (
	"net/http"

	"github.com/mm-sam/jobrunner"
)

type CronTab struct {
	config *Config
	http   *http.Server
}

func NewCronTab(conf *Config) *CronTab {
	return &CronTab{
		config: conf,
	}
}

func (c *CronTab) HTTPService() {
	httpService := NewHTTP(c.config)
	c.http = httpService.Start()
}

func (c *CronTab) Start() {
	jobrunner.Start()

	for _, task := range c.config.Tasks {

		t := CmdTaskWrapper{Cmd: task.CMD, Time: task.Time}
		err := jobrunner.Schedule(t.Time, t)
		if err != nil {
			Log.Error(err)
		}
	}
}

func (c *CronTab) Stop() {
	jobrunner.Stop()
	err := c.http.Close()
	if err != nil {
		Log.Error(err)
	}
}
