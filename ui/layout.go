package ui

import (
	"fmt"
	"time"

	"github.com/jroimartin/gocui"
)

var (
	count    = 0
	LogView  = "Log"
	ListView = "List"
)

func NewLayout(g *gocui.Gui, names []string) *Layout {
	return &Layout{g, names}
}

type Layout struct {
	g     *gocui.Gui
	names []string
}

func (l *Layout) Manage(*gocui.Gui) (err error) {
	l.logView()
	l.listView()
	return
}

func (l *Layout) logView() (err error) {
	maxX, maxY := l.g.Size()
	if v, err := l.g.SetView(LogView, 0, maxY-5, maxX-1, maxY-1); err != nil {
		v.Title = LogView
	}
	return
}

func (l *Layout) listView() (err error) {
	maxX, maxY := l.g.Size()
	v, err := l.g.SetView(ListView, 0, 0, maxX-1, maxY-6)
	if err != nil {
		return
	}
	// v.Title = "view title"
	// v.Wrap = true
	// 滚动到底部
	v.Highlight = true
	v.FgColor = gocui.ColorBlue
	// v.SelFgColor = gocui.ColorYellow
	v.SelFgColor = gocui.AttrBold + gocui.ColorYellow
	// time.Sleep(time.Second * 3)
	v.Clear()
	// return

	// fmt.Println(names)
	// return
	// dir, _ := os.Getwd()
	// tardir := filepath.Join(dir, "../../")
	// fs, _ := ioutil.ReadDir(tardir)

	// count = count + 1
	// stdout(g, fmt.Sprint(count))
	// stdout(g, fmt.Sprint(len(fs)))
	for _, name := range l.names {
		// fmt.Fprintln(v, name, i) //  写入到stdout
		fmt.Fprintln(v, name) //  写入到stdout
	}
	// cuidecor.Update()
	// update(g)
	// fmt.Fprintln(v, "Hello world!")
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
