package raven

import (
	"fmt"
	"strings"
	"time"
)

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
}

func (this FmtLogger) Debug(v ...interface{}) {
	strArr := make([]string, 0, len(v))
	for _, m := range v {
		strArr = append(strArr, fmt.Sprintf("[%v]", m))
	}

	fmt.Printf("%s [DBG] %s\n",
		time.Now().Format(time.StampMilli), strings.Join(strArr, " "))
	return
}

func (this FmtLogger) Info(v ...interface{}) {

	strArr := make([]string, 0, len(v))
	for _, m := range v {
		strArr = append(strArr, fmt.Sprintf("[%v]", m))
	}

	fmt.Printf("%s [INF] %s\n",
		time.Now().Format(time.StampMilli), strings.Join(strArr, " "))
	return
}

func (this FmtLogger) Warning(v ...interface{}) {

	strArr := make([]string, 0, len(v))
	for _, m := range v {
		strArr = append(strArr, fmt.Sprintf("[%v]", m))
	}

	fmt.Printf("%s [WRN] %s\n",
		time.Now().Format(time.StampMilli), strings.Join(strArr, " "))
	return
}

func (this FmtLogger) Error(v ...interface{}) {

	strArr := make([]string, 0, len(v))
	for _, m := range v {
		strArr = append(strArr, fmt.Sprintf("[%v]", m))
	}

	fmt.Printf("%s [ERR] %s\n",
		time.Now().Format(time.StampMilli), strings.Join(strArr, " "))
	return
}

func (this FmtLogger) Fatal(v ...interface{}) {
	strArr := make([]string, 0, len(v))
	for _, m := range v {
		strArr = append(strArr, fmt.Sprintf("[%v]", m))
	}

	fmt.Printf("%s [FAL] %s\n",
		time.Now().Format(time.StampMilli), strings.Join(strArr, " "))

	return
}
