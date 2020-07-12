package lib

import (
	"github.com/op/go-logging"
)

var Log = logging.MustGetLogger("mmcron")

func init() {
	format := logging.MustStringFormatter(
		`ipa2s3 %{color} %{shortfunc} %{level:.4s} %{shortfile}
%{id:03x}%{color:reset} %{message}`,
	)
	logging.SetFormatter(format)
}
