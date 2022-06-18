package ui

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/willf/pad"
)

var space = pad.Left("", 10, "-")

type update_data struct {
	Data []interface{}
}

type Logger func(...interface{})

func (f Logger) Update(val ...interface{}) {
	data := update_data{Data: append(store, append([]interface{}{space}, val[:]...)...)}
	f(data)
}

var (
	Log   Logger
	store []interface{}
)

// Curry Logger
func RegistLogger(g *gocui.Gui) Logger {
	Log = func(val ...interface{}) {
		if updateDate, ok := val[0].(update_data); !ok { // 是否为更新数据
			store = val
		} else {
			val = updateDate.Data
		}
		maxX, _ := g.Size()
		// str := fmt.Sprintln(append(val, []interface{}{maxX, maxY})...)
		str := fmt.Sprintln(val...)
		length := maxX - len(str)
		str = strings.ReplaceAll(str, space, pad.Right("", len(space)+length+2, " "))
		stdout(g, str)
	}
	return Log
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
