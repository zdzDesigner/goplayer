package main

import (
	"fmt"
	"player/audio"
	"player/conf"
	"player/ctrl"
	"player/ctrl/event"
	"player/serial"

	"player/ui"
	"player/util"
	"time"
)

func main() {
	// go pprof()
	// 歌曲列表
	names := conf.List()
	fmt.Println(names)

	// 选择一首
	event.Evt.On("CHOOSE", func(t any) {
		name := event.StringVal(t)
		go audio.Music(conf.FilePath(name))
		ui.Log(name)
		// ui.Log("0", "...")
	})

	// 第一首
	event.Evt.On("TOP", func(t any) {
		go audio.Music(names[0])
		ui.Nui.Layout.CursorIndex(0)
	})
	// 最后一首
	event.Evt.On("BOTTOM", func(t any) {
		lastindex := len(names) - 1
		go audio.Music(names[lastindex])
		ui.Nui.Layout.CursorIndex(lastindex)
	})
	// 上一首
	event.Evt.On("PREV", func(t any) {
		name, index := event.NextVal(t)
		if index == -1 {
			index = conf.PrevIndex(name)
			name = names[index]
		}
		go audio.Music(name)
		ui.Nui.Layout.CursorIndex(index)
	})
	// 下一首
	event.Evt.On("NEXT", func(t any) {
		name, index := event.NextVal(t)
		if index == -1 {
			index = conf.NextIndex(name)
			name = names[index]
		}
		go audio.Music(name)
		ui.Nui.Layout.CursorIndex(index)
		// ui.Log(conf.PrifixFileName(name))
		// ui.Log(name, "...")
		// ui.Log("defer NEXT", name, index)
	})

	event.Evt.On("DELETE", func(t any) {
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

	event.Evt.On("AUDIO_CTRL", func(state any) {
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
	go serialCtrl()

	// 视图
	ui.View(names)

}

func main2() {
	serialCtrl()

	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//
	// <-sigs
}

func serialCtrl() {
	throttler := util.NewThrottler(time.Millisecond * 500)

	// serialer, err := serial.NewSerialer("", &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200})
	// serialer, err := serial.NewSerialer("", &serial.Config{Name: "/dev/ttyUSB1", Baud: 115200})
	// serialer, err := serial.NewSerialer("", &serial.Config{Name: "/dev/ttyCH341USB1", Baud: 115200})
	// serialer, err := serial.NewSerialer("", &serial.Config{Name: "/dev/ttyCH341USB2", Baud: 115200})
	serialer, err := serial.NewSerialer("", &serial.Config{Name: "/dev/ttyCH341USB0", Baud: 115200})
	if err != nil {
		fmt.Println(err)
		return
	}
	count := 0
	for {
		data, err := serialer.Request("", nil)
		if err != nil {
			continue
		}
		// fmt.Println("-----------------")

		count++
		throttler.Do(func() {
			// fmt.Println(data, count)
			if data == "right" {
				event.Evt.Emit("NEXT", event.NewNext(conf.PrifixFileName(audio.PlayName), -1))
			}
			if data == "left" {
				event.Evt.Emit("PREV", event.NewNext(conf.PrifixFileName(audio.PlayName), -1))
			}
			if data == "top" {
				event.Evt.Emit("TOP", event.NewNext(conf.PrifixFileName(audio.PlayName), -1))
			}
			if data == "bottom" {
				event.Evt.Emit("BOTTOM", event.NewNext(conf.PrifixFileName(audio.PlayName), -1))
			}
			count = 0
		}, true)
		// if data == "right" {
		// 	count++
		// 	throttler.Do(func() {
		// 		fmt.Println("right", count)
		// 		count = 0
		// 		// event.Evt.Emit("NEXT", event.NewNext(conf.PrifixFileName(audio.PlayName), -1))
		// 	}, true)
		// } else if data == "left" {
		// 	count++
		// 	throttler.Do(func() {
		// 		fmt.Println("left", count)
		// 		count = 0
		// 		event.Evt.Emit("PREV", event.NewNext(conf.PrifixFileName(audio.PlayName), -1))
		// 	}, false)
		// } else if data == "top" {
		// 	count++
		// 	throttler.Do(func() {
		// 		fmt.Println("top", count)
		// 		count = 0
		// 		event.Evt.Emit("TOP", event.NewNext(conf.PrifixFileName(audio.PlayName), -1))
		// 	}, false)
		// } else if data == "bottom" {
		// 	count++
		// 	throttler.Do(func() {
		// 		fmt.Println("bottom", count)
		// 		count = 0
		// 		event.Evt.Emit("BOTTOM", event.NewNext(conf.PrifixFileName(audio.PlayName), -1))
		// 	}, false)
		// } else {
		// 	fmt.Println("reset")
		// }

	}

}
