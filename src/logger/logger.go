package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	debug = log.New(os.Stdout, "[DEBUG]", log.Ldate|log.Ltime|log.Lshortfile).Output
	info  = log.New(os.Stdout, "[INFO]", log.Ldate|log.Ltime|log.Lshortfile).Output
	warn  = log.New(os.Stdout, "[WARN]", log.Ldate|log.Ltime|log.Lshortfile).Output
	err   = log.New(os.Stderr, "[ERROR]", log.Ldate|log.Ltime|log.Lshortfile).Output
)

func Debugf(format string, v ...interface{}) {
	debug(2, fmt.Sprintf(format, v...))
}

func Infof(format string, v ...interface{}) {
	info(2, fmt.Sprintf(format, v...))
}

func Warnf(format string, v ...interface{}) {
	warn(2, fmt.Sprintf(format, v...))
}

func Errorf(format string, v ...interface{}) {
	err(2, fmt.Sprintf(format, v...))
}
