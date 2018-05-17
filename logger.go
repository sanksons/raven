package raven

import (
	"fmt"
)

type LOGLEVEL int

const LOG_LEVEL_OFF LOGLEVEL = 0
const LOG_LEVEL_ERR LOGLEVEL = 1
const LOG_LEVEL_WRN LOGLEVEL = 2
const LOG_LEVEL_DBG LOGLEVEL = 3

type Logger struct {
	level LOGLEVEL
}

func (this *Logger) Debug(msg interface{}) {
	this.log(LOG_LEVEL_DBG, msg)
}

func (this *Logger) Warning(msg interface{}) {
	this.log(LOG_LEVEL_WRN, msg)
}

func (this *Logger) Error(msg interface{}) {
	this.log(LOG_LEVEL_ERR, msg)
}

func (this *Logger) log(level LOGLEVEL, msg interface{}) {
	if level <= this.level {
		fmt.Printf("%v\n", msg)
	}
}
