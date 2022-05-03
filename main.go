package main

import (
	"fmt"
	"player/app"
	"player/conf"
	"player/event"
	"player/ui"
)

func main() {
	go pprof()
	// 歌曲列表
	names := conf.List()
	// fmt.Println(names)

	// 选择一首
	event.Evt.On("choose", func(name string) {
		app.Force <- struct{}{} // 强制结束
		go app.Music(conf.FilePath(name))

		ui.Log(fmt.Sprintln(conf.FileName(name), "..."))
	})

	// 下一首
	event.Evt.On("next", func(name string) {
		index := conf.NextIndex(name)
		name = names[index]
		go app.Music(name)

		ui.Nui.Layout.CursorIndex(index)
		ui.Log(fmt.Sprintln(conf.FileName(name), "..."))
	})

	go app.Music(names[0])

	// 视图
	ui.View(names)
}
