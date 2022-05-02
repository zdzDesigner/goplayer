package app

import (
	"fmt"
	"os"
	"player/conf"
	"player/ui"
	"player/util"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var forck chan struct{}

func Command(end, ch chan struct{}, currIndex chan int, names []string) {
	var (
		order string
		cmds  = []string{"next", "pre", "del!", "jump", ".exit"}
	)
	for {
		fmt.Scanf("%s", &order)
		fmt.Println("order::", order, "expect::", strings.Join(cmds, " | "))
		if util.Contains(cmds, order) {
			ch <- struct{}{}
			// time.Sleep(time.Second)
			i := <-currIndex
			if order == "del!" {
				conf.DelSong(names[i])
			}
			switch order {
			case "next":
				i++
			case "del!":
				i++
			case "pre":
				i--
			case "jump":
				i += 10
			case ".exit":
				end <- struct{}{}
			default:

			}

			order = ""
			ch = make(chan struct{}, 1)

			// go play(names[1], ch)
			go func() {
				fmt.Println("i::", i)
				currIndex <- Control(ch, i, names)
				// fmt.Println("currIndex::", i)
			}()
		}

	}
}

func reIndex(index, total int) int {
	fmt.Println("index::", index)
	if index > total-1 {
		return 0
	}
	if index < 0 {
		return total - 1
	}
	return index
}

func Control(forck chan struct{}, index int, names []string) int {
	// TODO::bug p1
	defer func() {
		if err := recover(); err != nil {
			index = index + 1
			Control(forck, index, names)
		}
		// return index
	}()
	// fmt.Println("currIndex ----current::", index)

	index = reIndex(index, len(names))
	// if play(names[util.Random(len(names))], forck) {
	if Play(names[index]) {
		return Control(forck, index+1, names)
	}
	// fmt.Println("currIndex ----current::", index)
	return index
}

func Listen() {
	ui.PS.On("play", func(name string) {
		// forck <- struct{}{}
		ui.Log(fmt.Sprintln(name, "----"))
		// time.Sleep(time.Second * 3)
		// forck = make(chan struct{}, 1)
		// Play(name)
	})
}

func Music(name string) {
	forck = make(chan struct{}, 1)
	Listen()
	Play(name)
}

// forck 外部信号停止内部执行, 可以使用context
func Play(name string) bool {
	var err error
	done := make(chan struct{}, 1)
	defer func() {
		// fmt.Println("play end==============")
		if err != nil {
			panic(err)
		}
	}()

	// source := name
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

	// fmt.Println("bfmt.SampleRate::", bfmt.SampleRate)
	if err = speaker.Init(bfmt.SampleRate, bfmt.SampleRate.N(time.Second/10)); err != nil {
		return false
	}
	// fmt.Println("speaker.Init aflter")

	speaker.Play(beep.Seq(stm, beep.Callback(func() {
		// fmt.Println("------------end---------------")
		done <- struct{}{}
	})))

	select {
	case <-done:
		return true
	case <-forck:
		return false
	}
}
