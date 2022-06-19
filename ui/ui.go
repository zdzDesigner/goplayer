package ui

import (
	"log"
	"player/conf"
	"player/lib/gocui"
	"time"
)

var Nui *UI

func View(names []string) {
	g, _ := gocui.NewGui(gocui.OutputTrue, false, gocui.NORMAL, false, map[rune]string{})
	defer g.Close()
	Nui = NewUI(g)
	Nui.RegistLog()
	Nui.ForceUpdate()
	Nui.layout(names)
	Nui.keybind()
	go func() {
		time.Sleep(time.Microsecond * 10)
		Nui.Log(conf.PrifixFileName(names[0]))
	}()

	// <-channel 主循环
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func NewUI(g *gocui.Gui) *UI {
	return &UI{g: g}
}

type UI struct {
	g      *gocui.Gui
	Layout *Layout
	Log    Logger
	Update Updater
}

func (ui *UI) layout(names []string) {
	ui.Layout = NewLayout(ui.g, names)
	ui.g.SetManagerFunc(ui.Layout.Manage)
}

func (ui *UI) ForceUpdate() {
	ui.Update = ForceUpdate(ui.g)
}

func (ui *UI) keybind() {
	Keybind(ui.g)
}

func (ui *UI) RegistLog() {
	ui.Log = RegistLogger(ui.g)
}
