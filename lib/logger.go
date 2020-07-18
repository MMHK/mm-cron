package lib

import (
	"github.com/op/go-logging"
)

var Log = logging.MustGetLogger("mm-cron")

func init() {
	format := logging.MustStringFormatter(
		`mm-cron %{color} %{shortfunc} %{level:.4s} %{shortfile}
%{id:03x}%{color:reset} %{message}`,
	)
	logging.SetFormatter(format)
}
