package main

import (
	"fmt"
	"player/audio"
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
	event.Evt.On("CHOOSE", func(name string) {
		audio.Force <- struct{}{} // 强制结束
		go audio.Music(conf.FilePath(name))

		ui.Log(fmt.Sprintln(conf.PrifixFileName(name), "..."))
	})

	// 下一首
	event.Evt.On("NEXT", func(name string) {
		index := conf.NextIndex(name)
		name = names[index]
		go audio.Music(name)

		ui.Nui.Layout.CursorIndex(index)
		ui.Log(fmt.Sprintln(conf.PrifixFileName(name), "..."))
		ui.Log("defer NEXT", name, index)
	})

	event.Evt.On("DELETE", func(name string) {
		ui.Log(fmt.Sprintln("DELETE", name, "PlayName", conf.FileName(audio.PlayName)))
		if conf.ClearPrefix(name) == conf.FileName(audio.PlayName) { // 正在播放
			audio.Force <- struct{}{} // 强制结束
			defer func() {
				event.Evt.Emit("NEXT", name)
			}()
		}
		conf.DelSong(conf.ClearPrefix(name))
		names = conf.UpdateList()
		ui.Nui.Layout.UpdateList(names)
		ui.Log(" NEXT")
		// 2 不在播放
		// go audio.Music(conf.FilePath(name))

		// ui.Log(fmt.Sprintln(conf.PrifixFileName(name), "..."))
	})

	go audio.Music(names[0])

	// 视图
	ui.View(names)
}
