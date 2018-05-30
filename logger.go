package raven

import (
	"fmt"
	"strings"
	"time"
)

const FATAL_LEVEL = 0
const ERROR_LEVEL = 0
const WRN_LEVEL = 1
const INFO_LEVEL = 2
const DBG_LEVEL = 3

var _ Logger = (*DummyLogger)(nil)
var _ Logger = (*FmtLogger)(nil)

type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warning(...interface{})
	Error(...interface{})
	Fatal(...interface{})
}

//
// This is used when no logger is specified.
//
type DummyLogger struct {
}

func (this *DummyLogger) Debug(...interface{}) {
	return
}

func (this *DummyLogger) Info(...interface{}) {
	return
}

func (this *DummyLogger) Warning(...interface{}) {
	return
}

func (this *DummyLogger) Error(...interface{}) {
	return
}

func (this DummyLogger) Fatal(...interface{}) {
	return
}

//
// Helper logger.
//
type FmtLogger struct {
	Level int
}

func (this FmtLogger) Debug(v ...interface{}) {

	if this.Level < DBG_LEVEL {
		return
	}

	strArr := make([]string, 0, len(v))
	for _, m := range v {
		strArr = append(strArr, fmt.Sprintf("[%v]", m))
	}

	fmt.Printf("%s [DBG] %s\n",
		time.Now().Format(time.StampMilli), strings.Join(strArr, " "))
	return
}

func (this FmtLogger) Info(v ...interface{}) {

	if this.Level < INFO_LEVEL {
		return
	}

	strArr := make([]string, 0, len(v))
	for _, m := range v {
		strArr = append(strArr, fmt.Sprintf("[%v]", m))
	}

	fmt.Printf("%s [INF] %s\n",
		time.Now().Format(time.StampMilli), strings.Join(strArr, " "))
	return
}

func (this FmtLogger) Warning(v ...interface{}) {

	if this.Level < WRN_LEVEL {
		return
	}

	strArr := make([]string, 0, len(v))
	for _, m := range v {
		strArr = append(strArr, fmt.Sprintf("[%v]", m))
	}

	fmt.Printf("%s [WRN] %s\n",
		time.Now().Format(time.StampMilli), strings.Join(strArr, " "))
	return
}

func (this FmtLogger) Error(v ...interface{}) {

	if this.Level < ERROR_LEVEL {
		return
	}

	strArr := make([]string, 0, len(v))
	for _, m := range v {
		strArr = append(strArr, fmt.Sprintf("[%v]", m))
	}

	fmt.Printf("%s [ERR] %s\n",
		time.Now().Format(time.StampMilli), strings.Join(strArr, " "))
	return
}

func (this FmtLogger) Fatal(v ...interface{}) {
	if this.Level < ERROR_LEVEL {
		return
	}
	strArr := make([]string, 0, len(v))
	for _, m := range v {
		strArr = append(strArr, fmt.Sprintf("[%v]", m))
	}

	fmt.Printf("%s [FAL] %s\n",
		time.Now().Format(time.StampMilli), strings.Join(strArr, " "))

	return
}
