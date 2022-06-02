package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type Logger func(...interface{})

var Log Logger

// Curry Logger
func RegistLogger(g *gocui.Gui) Logger {
	return func(val ...interface{}) {
		stdout(g, fmt.Sprintln(val...))
	}
}

// DOTO:: 记录日志列表, 可以上下滚动查询
func stdout(g *gocui.Gui, log string) (err error) {
	v, err := g.View(CurView)
	if err != nil {
		return
	}
	v.Clear()            // 清除输出
	fmt.Fprintln(v, log) // 写入到stdout

	return
}
