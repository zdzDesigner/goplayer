package ui

import (
	"fmt"
	"player/lib/gocui"
	"regexp"
	"strings"

	"github.com/willf/pad"
)

var space = pad.Left("", 1, "@")

type Progresser interface {
	Val() int
}

type update_data struct {
	Data []interface{}
}

type Logger func(...interface{})

func (f Logger) Progress(val int) {
	// data := update_data{Data: append(store, val[:]...)}
	// f(data)
}

func (f Logger) Update(val ...interface{}) {
	data := update_data{Data: append(store, val[:]...)}
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
			store = append(val, space)
			// store = val
		} else {
			val = updateDate.Data
		}
		maxX, _ := g.Size()

		str := fmt.Sprintln(val...)
		// str := fmt.Sprintf(pad.Left("", len(val), "%s"), val...)
		curlen := len(hantoen(str))
		length := maxX - 2 - curlen
		// str = strings.ReplaceAll(str, space, pad.Left("", len(space)+length, "-"))
		str = strings.ReplaceAll(str, space, pad.Left("", len(space)+length, " "))
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
	v.Clear() // 清除输出
	// fmt.Fprintln(v, fmt.Sprint("aa 绷带俱乐部拉乌33")) //  写入到stdout
	fmt.Fprintln(v, log) // 写入到stdout

	return
}

func hantoen(str string) string {
	return regexp.MustCompile("[\u4e00-\u9fa5]").ReplaceAllLiteralString(str, "11")
}
