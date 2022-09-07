package main

import (
	"player/audio"
	"player/conf"
	"player/ctrl"
	"player/ctrl/event"
	"player/ui"
)

func main() {
	go pprof()
	// 歌曲列表
	names := conf.List()
	// fmt.Println(names)

	// 选择一首
	event.Evt.On("CHOOSE", func(t interface{}) {
		name := event.StringVal(t)
		audio.Stop()
		go audio.Music(conf.FilePath(name))

		// ui.Log(name)
		// ui.Log("0", "...")
	})

	// 下一首
	event.Evt.On("NEXT", func(t interface{}) {
		name, index := event.NextVal(t)
		if index == -1 {
			index = conf.NextIndex(name)
			name = names[index]
		}
		go audio.Music(name)

		// return
		ui.Nui.Layout.CursorIndex(index)
		ui.Log(conf.PrifixFileName(name))
		// ui.Log(name, "...")
		// ui.Log("defer NEXT", name, index)
	})

	event.Evt.On("DELETE", func(t interface{}) {
		name := conf.ClearPrefix(event.StringVal(t))
		// ui.Log(fmt.Sprintln("DELETE", name, "PlayName", conf.FileName(audio.PlayName)))
		if name == conf.FileName(audio.PlayName) { // 正在播放
			audio.Stop()
			nextname := names[conf.NextIndex(name)] // 先查到下一个播放的歌曲名
			defer func() {                          // 等待更新后再根据目标歌曲名查到index
				event.Evt.Emit("NEXT", event.NewNext(nextname, conf.Index(nextname)))
			}()
		}
		conf.DelSong(name)
		names = conf.UpdateList()
		ui.Nui.Layout.UpdateList(names)
	})

	event.Evt.On("AUDIO_CTRL", func(state interface{}) {
		val, ok := state.(string)
		if !ok {
			return
		}
		// ui.Log("CTRL::", val)
		switch val {
		case "PAUSE":
			audio.Paused()
		case "NEXT":
			// ui.Log("CTRL:: NEXT", audio.PlayName)
			event.Evt.Emit("NEXT", event.NewNext(conf.PrifixFileName(audio.PlayName), -1))
		}
	})
	go audio.Music(names[0])
	go ctrl.ListenGlobal()

	// 视图
	ui.View(names)
}
