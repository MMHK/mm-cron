// config
package lib

import (
	"encoding/json"
	"os"
	"path/filepath"
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

func (this *Config) SaveTasks() (error) {
	tasksConfigPath := filepath.Join(filepath.Dir(this.sava_file), "tasks.json")
	var file *os.File
	if _, err := os.Stat(tasksConfigPath); err != nil && os.IsNotExist(err) {
		file, err = os.Create(tasksConfigPath)
		if err != nil {
			Log.Error(err)
			return err
		}
		defer file.Close()
	} else {
		file, err = os.Open(tasksConfigPath)
		if err != nil {
			Log.Error(err)
			return err
		}
		defer file.Close()
	}
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "   ")
	err := encoder.Encode(this.Tasks)
	if err != nil {
		Log.Error(err)
		return err
	}
	
	return nil
}

func (this *Config) LoadTasks() (error) {
	tasksConfigPath := filepath.Join(filepath.Dir(this.sava_file), "tasks.json")
	if _, err := os.Stat(tasksConfigPath); err != nil && os.IsNotExist(err) {
		return err
	}
	file, err := os.Open(tasksConfigPath)
	if err != nil {
		Log.Error(err)
		return err
	}
	defer file.Close()
	
	decoder := json.NewDecoder(file)
	return decoder.Decode(&this.Tasks)
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
