package app

import (
	"fmt"
	"os"
	"player/util"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func Ctrl(end chan struct{}, ch chan struct{}, currIndex chan int, names []string) {
	var (
		order string
		cmds  = []string{"next", "pre", "del!", "jump", "exit"}
	)
	for {
		fmt.Scanf("%s", &order)
		fmt.Println("order::", order, "expect::", strings.Join(cmds, " | "))
		if util.Contains(cmds, order) {
			ch <- struct{}{}
			// time.Sleep(time.Second)
			i := <-currIndex
			if order == "del!" {
				DelSong(names[i])
			}
			switch order {
			case "next":
			case "del!":
				i++
			case "pre":
				i--
			case "jump":
				i += 10
			case "exit":
				end <- struct{}{}
			default:

			}

			order = ""
			ch = make(chan struct{}, 1)

			// go play(names[1], ch)
			go func() {
				currIndex <- Control(ch, i, names)
				fmt.Println("currIndex::", i)
			}()
		}

	}
}

func Control(trigger chan struct{}, index int, names []string) int {
	// TODO::bug p1
	defer func() {
		if err := recover(); err != nil {
			index = index + 1
			Control(trigger, index, names)
		}
		// return index
	}()
	// fmt.Println("currIndex ----current::", index)
	if index > len(names)-1 {
		index = 0
	}
	if index < 0 {
		index = len(names) - 1
	}

	if play(names[util.Random(len(names))], trigger) {
		// if play(names[index], trigger) {
		return Control(trigger, index+1, names)
	}
	fmt.Println("currIndex ----current::", index)
	// if index == 0 {
	// 	index = 1
	// }
	return index
}

func play(name string, trigger chan struct{}) bool {
	var err error
	done := make(chan struct{}, 1)
	defer func() {
		fmt.Println("play end==============")
		if err != nil {
			panic(err)
		}
	}()

	source, err := getSource(name)
	if err != nil {
		return false
	}

	// source, err := getSource("output.mp3")
	fmt.Println(source)
	file, err := os.Open(source)
	if err != nil {
		return false
	}

	stm, bfmt, err := mp3.Decode(file)
	defer stm.Close()

	fmt.Println("bfmt.SampleRate::", bfmt.SampleRate)
	if err = speaker.Init(bfmt.SampleRate, bfmt.SampleRate.N(time.Second/10)); err != nil {
		return false
	}
	fmt.Println("speaker.Init aflter")

	speaker.Play(beep.Seq(stm, beep.Callback(func() {
		fmt.Println("------------end---------------")
		done <- struct{}{}
	})))

	select {
	case <-done:
		return true
	case <-trigger:
		return false
	}
}
