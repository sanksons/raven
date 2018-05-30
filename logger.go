package raven

import (
	"fmt"
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
	fmt.Println(v)
	return
}

func (this FmtLogger) Info(v ...interface{}) {
	fmt.Println(v)
	return
}

func (this FmtLogger) Warning(v ...interface{}) {
	fmt.Println(v)
	return
}

func (this FmtLogger) Error(v ...interface{}) {
	fmt.Println(v)
	return
}

func (this FmtLogger) Fatal(v ...interface{}) {
	fmt.Println(v)
	return
}
