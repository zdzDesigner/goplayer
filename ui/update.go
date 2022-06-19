package ui

import "player/lib/gocui"

type Updater func()

var Update Updater

// 强制更新
func ForceUpdate(g *gocui.Gui) Updater {
	Update = func() {
		g.Update(func(g *gocui.Gui) error { return nil })
	}
	return Update
}
