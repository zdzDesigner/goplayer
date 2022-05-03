package ui

import (
	"fmt"
	"player/conf"
	"time"

	"github.com/jroimartin/gocui"
)

var (
	count    = 0
	CurView  = "Current"
	ListView = "Song List"
)

func NewLayout(g *gocui.Gui, names []string) *Layout {
	return &Layout{g: g, names: names}
}

type Layout struct {
	g     *gocui.Gui
	listV *gocui.View
	names []string
}

func (l *Layout) Manage(*gocui.Gui) (err error) {
	l.logView()
	l.listView()
	return
}

func (l *Layout) CursorIndex(y int) (err error) {
	Update()
	return KeyIndex(l.g, l.listV, y)
}

func (l *Layout) logView() (err error) {
	maxX, maxY := l.g.Size()
	v, err := l.g.SetView(CurView, 0, maxY-5, maxX-1, maxY-1)
	if err != nil {
		v.Title = CurView
	}
	// Log(conf.FileName(l.names[0]))
	return
}

func (l *Layout) listView() (err error) {
	maxX, maxY := l.g.Size()
	v, err := l.g.SetView(ListView, 0, 0, maxX-1, maxY-6)
	if err != nil {
		return
	}
	v.Title = ListView
	l.listV = v // 添加listV视图

	// 滚动到底部
	v.Highlight = true
	v.FgColor = gocui.ColorBlue
	v.SelFgColor = gocui.AttrBold + gocui.ColorYellow

	v.Clear()
	for _, name := range l.names {
		// fmt.Fprintln(v, name, i) //  写入到stdout
		fmt.Fprintln(v, conf.FileName(name)) //  写入到stdout
	}
	_, err = l.g.SetCurrentView(ListView)
	return
}

func update(g *gocui.Gui) {
	time.Sleep(time.Second * 2)

	g.Update(func(g *gocui.Gui) (err error) {
		// return nil
		v, err := g.View(ListView)
		if err != nil {
			return
		}
		v.Clear()
		_, err = g.SetCurrentView(ListView)
		return
	})
}
