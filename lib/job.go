package lib

import (
	"bytes"
	"os/exec"
	"strings"
)

type CmdTaskWrapper struct {
	Time string `json:"time"`
	Cmd  string `json:"cmd"`
}

func (w CmdTaskWrapper) Run() {
	Log.Info("starting run task:", w.Cmd)

	cmds := strings.Split(w.Cmd, " ")
	cmd := exec.Command(cmds[0], cmds[1:]...)
	var outbuffer bytes.Buffer
	var errbuffer bytes.Buffer
	cmd.Stdout = &outbuffer
	cmd.Stderr = &errbuffer
	err := cmd.Run()
	if err != nil {
		Log.Error(err)
		Log.Error(errbuffer.String())
		return
	}
	Log.Info(outbuffer.String())
}
