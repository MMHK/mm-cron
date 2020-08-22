package main

import (
	"mm-cron/lib"
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	conf_path := flag.String("c", "conf.json", "config json file")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	err, conf := lib.NewConfig(*conf_path)

	if err != nil {
		lib.Log.Error(err)
		return
	}
	
	err = conf.LoadTasks()
	if err != nil {
		lib.Log.Error(err)
	}

	cron := lib.NewCronTab(conf)

	cron.Start()
	cron.HTTPService()

	lib.Log.Info("Enter `Ctrl + C` 2 Exit")
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	cron.Stop()
	os.Exit(1)
}
