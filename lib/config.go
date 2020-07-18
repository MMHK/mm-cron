// config
package lib

import (
	"encoding/json"
	"os"
	"sync"
)

// swagger:model Task
type Task struct {
	// in:body
	
	// job time cron format
	// required: true
	Time string `json:"time"`
	// Command Line Job
	// required: true
	CMD  string `json:"cmd"`
}

type Config struct {
	lock      *sync.RWMutex
	sava_file string
	Tasks     []*Task `json:"task"`
	WebRoot   string  `json:"web_root"`
	Listen    string  `json:"listen"`
}

func NewConfig(filename string) (err error, c *Config) {
	c = &Config{
		lock: new(sync.RWMutex),
	}
	c.sava_file = filename
	err = c.load(filename)
	return
}

func (c *Config) load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		Log.Error(err)
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		Log.Error(err)
		return err
	}
	return nil
}

func (c *Config) Save() error {
	file, err := os.Create(c.sava_file)
	if err != nil {
		Log.Error(err)
		return err
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	defer file.Close()
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		Log.Error(err)
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		Log.Error(err)
		return err
	}
	return nil
}
