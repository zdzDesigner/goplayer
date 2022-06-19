package ui

import (
	"errors"
	"fmt"
	"log"
	"player/conf"
	"player/event"
	"player/lib/gocui"
)

var index = 0

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
		Log(ox, oy)
		v.SetOrigin(ox, oy+val)
		// v.FocusPoint(ox, oy+val)
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

func KeyIndex(_ *gocui.Gui, v *gocui.View, y int) (err error) {
	// v.FocusPoint(0, y)
	return v.SetCursor(0, y)
}

func KeyDown(_ *gocui.Gui, v *gocui.View) (err error) {
	listLength := len(conf.List()) - 1
	if v == nil {
		return errors.New("keydown view nil")
	}

	setCursor(v, 1, func(oy, cy int) bool {
		return oy+cy <= listLength
	})
	return
}

func KeyDel(_ *gocui.Gui, v *gocui.View) (err error) {
	if v == nil {
		return errors.New("keydown view nil")
	}
	_, cy := v.Cursor()
	cyline, _ := v.Line(cy)
	event.Evt.Emit("DELETE", cyline)
	return
}

func KeyAuidoCtrl(_ *gocui.Gui, v *gocui.View) (err error) {
	if v == nil {
		return errors.New("keydown view nil")
	}
	// _, cy := v.Cursor()
	// cyline, _ := v.Line(cy)
	event.Evt.Emit("AUDIO_CTRL", nil)
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
	event.Evt.Emit("CHOOSE", cyline)
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

	if err = g.SetKeybinding("", nil, 'q', gocui.ModNone, quit); err != nil {
		return
	}

	if err = g.SetKeybinding(ListView, nil, 'k', gocui.ModNone, keyUp); err != nil {
		return
	}
	if err = g.SetKeybinding(ListView, nil, 'j', gocui.ModNone, KeyDown); err != nil {
		return
	}
	if err = g.SetKeybinding(ListView, nil, 'd', gocui.ModNone, KeyDel); err != nil {
		return
	}
	if err = g.SetKeybinding(ListView, nil, 's', gocui.ModNone, KeyAuidoCtrl); err != nil {
		return
	}
	if err = g.SetKeybinding(ListView, nil, gocui.KeyCtrlG, gocui.ModNone, end); err != nil {
		return
	}

	if err = g.SetKeybinding(ListView, nil, gocui.KeyEnter, gocui.ModNone, enter); err != nil {
		return
	}
	if err = g.SetKeybinding(ListView, nil, gocui.KeySpace, gocui.ModNone, enter); err != nil {
		return
	}
}
