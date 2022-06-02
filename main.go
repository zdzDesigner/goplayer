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
	event.Evt.On("CHOOSE", func(it interface{}) {
		name := event.StringVal(it)
		audio.Force <- struct{}{} // 强制结束
		go audio.Music(conf.FilePath(name))

		ui.Log(fmt.Sprintln(conf.PrifixFileName(name), "..."))
	})

	// 下一首
	event.Evt.On("NEXT", func(it interface{}) {
		name, index := event.NextVal(it)
		if index == 0 {
			index = conf.NextIndex(name)
			name = names[index]
		}
		go audio.Music(name)

		ui.Nui.Layout.CursorIndex(index)
		ui.Log(fmt.Sprintln(conf.PrifixFileName(name), "..."))
		// ui.Log("defer NEXT", name, index)
	})

	event.Evt.On("DELETE", func(it interface{}) {
		name := conf.ClearPrefix(event.StringVal(it))
		// ui.Log(fmt.Sprintln("DELETE", name, "PlayName", conf.FileName(audio.PlayName)))
		if name == conf.FileName(audio.PlayName) { // 正在播放
			audio.Force <- struct{}{}               // 强制结束
			nextname := names[conf.NextIndex(name)] // 先查到下一个播放的歌曲名
			defer func() {                          // 等待更新后再根据目标歌曲名查到index
				event.Evt.Emit("NEXT", event.NewNext(nextname, conf.Index(nextname)))
			}()
		}
		conf.DelSong(name)
		names = conf.UpdateList()
		ui.Nui.Layout.UpdateList(names)
	})

	go audio.Music(names[0])

	// 视图
	ui.View(names)
}
