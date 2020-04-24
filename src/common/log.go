package common

import (

"github.com/go-logging"
"os"
)

//日志操作
var (
	Log  logging.Logger
)

func init(){
	var format = logging.MustStringFormatter(
		`%{color}%{time:2016-01-02 15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{shortfile} %{color:reset} %{message}`,
	)

	backend1 := logging.NewLogBackend(os.Stdout, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")

	// Set the backends to be used.
	logging.SetBackend(backend1Leveled, backend2Formatter)

}

