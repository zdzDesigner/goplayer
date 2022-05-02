package main

import (
	"player/app"
	"player/conf"
	"player/ui"
)

func main() {
	var err error // currIndex = make(chan int)
	// ch  = make(chan struct{}, 1) // 带缓冲区的通道, 允许写入和读出
	// end       = make(chan struct{})
	//
	defer func() {
		if err != nil {
			panic(err)
		}
	}()
	names, err := conf.AudioList()
	if err != nil {
		return
	}
	// fmt.Println(names)
	go app.Music(names[0])
	// //
	// // // fmt.Println("pid::", os.Getpid())
	// currIndex <- app.Control(ch, 0, names)
	// fmt.Println("currIndex first::", currIndex)
	//
	// // play(names[0], ch)
	// <-end
	ui.View(names)
}
