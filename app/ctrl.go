package app

import (
	"os"
	"player/conf"
	"player/event"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var Force chan struct{}

func Music(name string) {
	if Play(name) {
		event.Evt.Emit("next", conf.FileName(name))
	}
}

// Force 外部信号停止内部执行
func Play(name string) bool {
	var err error
	Force = make(chan struct{}, 1)   // 强制结束
	finish := make(chan struct{}, 1) // 单曲完成播放
	defer func() {
		if err != nil {
			panic(err)
		}
	}()

	source, err := getSource(name)
	if err != nil {
		return false
	}

	// source, err := getSource("output.mp3")
	// fmt.Println("source::", source)
	file, err := os.Open(source)
	if err != nil {
		return false
	}
	defer file.Close()

	stm, bfmt, err := mp3.Decode(file)
	if err != nil {
		return false
	}
	defer stm.Close()

	// 采样率
	if err = speaker.Init(bfmt.SampleRate, bfmt.SampleRate.N(time.Second/10)); err != nil {
		return false
	}

	// 播放完成
	speaker.Play(beep.Seq(stm, beep.Callback(func() {
		finish <- struct{}{}
	})))

	select {
	case <-finish:
		return true
	case <-Force:
		return false
	}
}
