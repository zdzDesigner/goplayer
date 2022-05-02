package ui

import (
	"log"
	"player/ui/cuidecor"

	"github.com/jroimartin/gocui"
)

func View(names []string) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	ui := NewUI(g)
	ui.Log()
	ui.ForceUpdate()
	ui.Layout(names)
	ui.Keybind()

	// go update(g)
	// gocui.KeyCtrlC  mod + ctrl + c
	// if err := g.SetKeybinding("hello", gocui.KeyEsc, gocui.ModNone, quit); err != nil {
	// 绑定退出
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func NewUI(g *gocui.Gui) *UI {
	return &UI{g}
}

type UI struct {
	g *gocui.Gui
}

func (ui *UI) Layout(names []string) {
	ui.g.SetManagerFunc(NewLayout(ui.g, names).Manage)
}

func (ui *UI) ForceUpdate() {
	cuidecor.Update = cuidecor.ForceUpdate(ui.g)
}

func (ui *UI) Keybind() {
	Keybind(ui.g)
}

func (ui *UI) Log() {
	Log = RegistLogger(ui.g)
}
