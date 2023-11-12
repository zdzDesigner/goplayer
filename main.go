package main

import (
	"fmt"
	// "os"
	// "os/signal"
	"player/audio"
	"player/conf"
	"player/ctrl"
	"player/ctrl/event"
	"player/serial"
	// "syscall"

	"player/ui"
	"player/util"
	"time"
)

var name = make(chan string, 1)

var names []string

func main() {
	// go pprof()
	// 歌曲列表
	names = conf.List()
	fmt.Println(names)

	// 选择一首
	event.Evt.On("CHOOSE", func(t any) {
		name := event.StringVal(t)
		// audio.Stop()
		go audio.Music(conf.FilePath(name))

		ui.Log(name)
		// ui.Log("0", "...")
	})

	// 上一首
	event.Evt.On("PREV", func(t any) {
		name, index := event.NextVal(t)
		if index == -1 {
			index = conf.PrevIndex(name)
			name = names[index]
		}
		// audio.Stop()
		// time.Sleep(time.Second * 3)
		// fmt.Println("index:",index)
		go audio.Music(name)

		// return
		ui.Nui.Layout.CursorIndex(index)
		// ui.Log(conf.PrifixFileName(name))
		// ui.Log(name, "...")
		// ui.Log("defer NEXT", name, index)
	})
	// 下一首
	event.Evt.On("NEXT", func(t any) {
		name, index := event.NextVal(t)
		if index == -1 {
			index = conf.NextIndex(name)
			name = names[index]
		}
		// time.Sleep(time.Second * 3)
		// fmt.Println("index:",index)
		go audio.Music(name)

		// return
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
	// serialCtrl()

	// 视图
	ui.View(names)

	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//
	// <-sigs

}

func serialCtrl() {
	throttler := util.NewThrottler(time.Millisecond * 500)

	serialer, err := serial.NewSerialer("", &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200})
	if err != nil {
		return
	}
	for {
		data, err := serialer.Request("", nil)
		if err != nil {
			continue
		}
		if data == "right" {
			throttler.Do(func() { event.Evt.Emit("NEXT", event.NewNext(conf.PrifixFileName(audio.PlayName), -1)) })
		} else if data == "left" {
			throttler.Do(func() { event.Evt.Emit("PREV", event.NewNext(conf.PrifixFileName(audio.PlayName), -1)) })
		} else if data == "top" {
			throttler.Do(func() { event.Evt.Emit("TOP", event.NewNext(conf.PrifixFileName(audio.PlayName), -1)) })
		} else if data == "bottom" {
			throttler.Do(func() { event.Evt.Emit("BOTTOM", event.NewNext(conf.PrifixFileName(audio.PlayName), -1)) })
		}

	}

}
