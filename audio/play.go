package audio

import (
	"os"
	"player/conf"
	"player/ctrl/event"
	"player/ui"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var (
	ctrl     *beep.Ctrl
	m        sync.Mutex
	Force    = make(chan struct{}, 1) // 强制结束,Force 外部信号停止内部执行
	PlayName = ""
)

func Music(name string) {
	if Play(name) {
		event.Evt.Emit("NEXT", event.NewNext(conf.PrifixFileName(PlayName), -1))
	}
}

// 播放
func Play(name string) (ok bool) {
	var err error
	finish := make(chan struct{}, 1) // 单曲完成播放
	defer func() {
		if err != nil {
			speaker.Close()
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

	stm, bfmt, err := mp3.Decode(file)
	if err != nil {
		return false
	}
	defer stm.Close()

	// go func() {
	// time.Sleep(time.Second * 2)
	// err = stm.Seek(stm.Len() - 20000)
	// if err != nil {
	// 	ui.Log("len", stm.Len(), "position:", stm.Position(), err.Error())
	// }
	// }()

	// 采样率
	if err = speaker.Init(bfmt.SampleRate, bfmt.SampleRate.N(time.Second/10)); err != nil {
		return false
	}

	ctrl = &beep.Ctrl{Streamer: stm}

	speaker.Play(ctrl, beep.StreamerFunc(func(_ [][2]float64) (n int, ok bool) {
		if stm.Position()%100 == 0 {
			seek := stm.Position()
			total := stm.Len()
			ui.Nui.Update()
			// ui.Log.Update(fmt.Sprintf("%d%s", seek*100/total, "%"))
			ui.Log.Progress(seek * 100 / total)
		}
		// err = stm.Seek(stm.Len() - 80000)
		// if err != nil {
		// 	ui.Log("len", stm.Len(), "position:", stm.Position(), err.Error())
		// }
		if stm.Position()+2000 >= stm.Len() {
			finish <- struct{}{}
			return 0, false
		}
		return 0, true
	}))

	// beep.Callback(func() {
	// finish <- struct{}{}
	// })

	// 播放完成
	// speaker.Play(beep.Seq(stm, beep.Callback(func() {
	// 	finish <- struct{}{}
	// })))

	select {
	case <-finish:
		return true
	case <-Force:
		return false
	}
}

// 停止
func Stop() {
	if cap(Force) > len(Force) {
		Force <- struct{}{} // 强制结束
	}
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
