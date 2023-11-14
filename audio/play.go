package audio

import (
	"os"

	"player/conf"
	"player/ctrl/event"
	"player/ui"
	"sync"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

var (
	ctrl     *beep.Ctrl
	m        sync.Mutex
	Force    = make(chan struct{}, 0) // 强制结束,Force 外部信号停止内部执行
	PlayName = ""
	isinit   = false
	isplay   = false
)

func Music(name string) {
	if isplay {
		Stop()
	}
	if Play(name) {
		event.Evt.Emit("NEXT", event.NewNext(conf.PrifixFileName(PlayName), -1))
	}
}

// 播放
func Play(name string) (ok bool) {
	var err error
	finish := make(chan struct{}, 0) // 单曲完成播放
	defer func() {
		isplay = false
		speaker.Clear()
		if err != nil {
			ok = true
		}
	}()
	Lock()
	PlayName = name
	Unlock()

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

	streamer, format, err := mp3.Decode(file)
	if err != nil {
		return false
	}
	defer streamer.Close()

	// go func() {
	// time.Sleep(time.Second * 2)
	// err = stm.Seek(stm.Len() - 20000)
	// if err != nil {
	// 	ui.Log("len", stm.Len(), "position:", stm.Position(), err.Error())
	// }
	// }()

	// 采样率
	if isinit == false {
		if err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/30)); err != nil {
			return false
		}
		isinit = true
	}

	ctrl = &beep.Ctrl{Streamer: beep.Loop(1, streamer)}
	// volume := &effects.Volume{Streamer: streamer, Base: 2}
	// speaker.Play(volume)

	// go func() {
	// 	time.Sleep(time.Second * 5)
	//
	// 	speaker.Lock()
	// 	streamer.Seek(streamer.Len() - 100)
	// 	speaker.Unlock()
	// }()

	speaker.Play(streamer, beep.StreamerFunc(func(_ [][2]float64) (n int, ok bool) {
		if streamer.Position()%100 == 0 {
			seek := streamer.Position()
			total := streamer.Len()
			ui.Nui.Update()
			ui.Log.Progress(seek * 100 / total)

			// ui.Log.Update(fmt.Sprintf("%d%s", seek*100/total, "%"))
		}
		// err = stm.Seek(stm.Len() - 80000)
		// if err != nil {
		// 	ui.Log("len", stm.Len(), "position:", stm.Position(), err.Error())
		// }
		if streamer.Position()+2000 >= streamer.Len() {
			finish <- struct{}{}
			return 0, false
		}
		return 0, true
	}))

	// 播放完成
	// speaker.Play(beep.Seq(stm, beep.Callback(func() {
	// 	finish <- struct{}{}
	// })))
	isplay = true

	select {
	case <-finish:
		// fmt.Println("============finish=============")
		return true
	case <-Force:
		// fmt.Println("============Force=============")
		return false
	}
}

// 停止
func Stop() {
	// fmt.Println(cap(Force), len(Force))
	// if cap(Force) > len(Force) {
	Force <- struct{}{} // 强制结束
	// }
}

// 暂停
func Paused() {
	speaker.Lock()
	ctrl.Paused = !ctrl.Paused
	speaker.Unlock()
}

func Lock() {
	m.Lock()
}

func Unlock() {
	m.Unlock()
}
