package util

import (
	"fmt"
	"reflect"
	"runtime"
)

func Trace() (funcTrace string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	funcTrace = fmt.Sprintf("%s,:%d %s\n", frame.File, frame.Line, frame.Function)
	fmt.Println(funcTrace)
	return
}

func FuncName(fn interface{}) (name string) {
	name = runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	return
}
