package helper

import (
	"fmt"
	"runtime"
	"strings"
)

// TraceCurrentFunc is function to get the current function name when it is called
func TraceCurrentFunc() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	splitFunc := strings.Split(frame.Function, "/")
	funcName := splitFunc[len(splitFunc)-1:]
	return fmt.Sprintf("%s", funcName[0])
}

func TraceCurrentFuncArgs(args ...interface{}) string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	splitFunc := strings.Split(frame.Function, "/")
	funcName := splitFunc[len(splitFunc)-1:]

	return fmt.Sprintf("%s(%+v)", funcName[0], args)
}
