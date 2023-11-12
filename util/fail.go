package util

import (
	"fmt"
	"os"
	"runtime/debug"
)

// 错误+退出
func FailExit(val any) {
	fmt.Printf("\033[31mfail exit\033[0m:%+v\n", val)
	debug.PrintStack()
	os.Exit(1)
}

// 无错
func MustNoErr(err error) {
	if err == nil {
		return
	}
	FailExit(err)
}

// 断言
func Assert(condition bool, errmsg any) {
	if condition {
		return
	}
	FailExit(errmsg)
}
