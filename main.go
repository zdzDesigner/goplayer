package main

import (
	"fmt"
	"player/app"
	"player/conf"
	"player/event"
	"player/ui"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			panic(err)
		}
	}()

	names := conf.List()
	// fmt.Println(names)
	event.Evt.On("choose", func(name string) {
		app.Force <- struct{}{} // 强制结束
		go app.Music(name)

		ui.Log(fmt.Sprintln(name, "----"))
	})

	event.Evt.On("next", func(name string) {
		index := conf.Index(name) + 1
		name = names[index]
		go app.Music(name)

		ui.Nui.Layout.ListVIndex(index)
		ui.Log(fmt.Sprintln(name, "----"))
	})

	go app.Music(names[0])

	// 视图
	ui.View(names)
}
