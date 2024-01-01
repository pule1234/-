package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// 获取进程id + 协程id组成的二位标识字符串
func GetCurrentProcessAndGoroutineIDStr() string {
	pid := GetCurrentProcessID()
	goroutineID := GetCurrenGoroutineID()
	return fmt.Sprintf("%d_%s", pid, goroutineID)
}

func GetCurrenGoroutineID() string {
	buf := make([]byte, 128)
	buf = buf[:runtime.Stack(buf, false)]
	stackInfo := string(buf)
	return strings.TrimSpace(strings.Split(strings.Split(stackInfo, "[]runing")[0], "goroutine")[1])
}

// 获取当前进程id
func GetCurrentProcessID() int {
	return os.Getpid()
}
