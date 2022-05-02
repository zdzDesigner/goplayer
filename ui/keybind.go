package ui

import (
	"errors"
	"fmt"
	"log"
	"player/conf"

	"github.com/jroimartin/gocui"
)

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// 高亮设置
func setCursor(v *gocui.View, val int, validator func(int, int) bool) (oy, cy int, err error) { // 滚动
	ox, oy := v.Origin() // 外层数据点
	cx, cy := v.Cursor() // 高亮光标点

	cy = cy + val
	if !validator(oy, cy) {
		return 0, 0, errors.New("cursor limit")
	}
	if err = v.SetCursor(cx, cy); err != nil { // 移动,  内层数据越界了(err)
		v.SetOrigin(ox, oy+val)
		err = nil
	}

	// Log(fmt.Sprintln("oy:", oy, "cy:", cy))
	return
}

func keyUp(_ *gocui.Gui, v *gocui.View) (err error) {
	if v == nil {
		return errors.New("keyup view nil")
	}

	setCursor(v, -1, func(oy, cy int) bool {
		return oy+cy > -1
	})
	return
}

func keyDown(_ *gocui.Gui, v *gocui.View) (err error) {
	names, err := conf.AudioList()
	if err != nil {
		return
	}

	listLength := len(names) - 1
	if v == nil {
		return errors.New("keydown view nil")
	}

	setCursor(v, 1, func(oy, cy int) bool {
		return oy+cy <= listLength
	})
	return
}

func end(_ *gocui.Gui, v *gocui.View) (err error) {
	if v == nil {
		return errors.New("keydown view nil")
	}
	v.SetCursor(0, 41)
	v.SetOrigin(0, 7)
	Log(fmt.Sprintln("end"))
	return
}

func enter(_ *gocui.Gui, v *gocui.View) (err error) {
	if v == nil {
		return errors.New("keydown view nil")
	}
	_, cy := v.Cursor()
	cyline, _ := v.Line(cy)
	// Log(cyline)
	PS.Emit("play", cyline)
	return
}

// 键盘
func Keybind(g *gocui.Gui) {
	var err error
	defer func() {
		if err != nil {
			log.Panic(err)
		}
	}()

	if err = g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return
	}

	if err = g.SetKeybinding(ListView, 'k', gocui.ModNone, keyUp); err != nil {
		return
	}
	if err = g.SetKeybinding(ListView, 'j', gocui.ModNone, keyDown); err != nil {
		return
	}
	if err = g.SetKeybinding(ListView, gocui.KeyCtrlG, gocui.ModNone, end); err != nil {
		return
	}

	if err = g.SetKeybinding(ListView, gocui.KeyEnter, gocui.ModNone, enter); err != nil {
		return
	}
	if err = g.SetKeybinding(ListView, gocui.KeySpace, gocui.ModNone, enter); err != nil {
		return
	}
}
