package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"player/app"
	"player/conf"
	"player/util"
	"strings"
)

func main() {
	var (
		err       error
		currIndex = make(chan int)
		ch        = make(chan struct{}, 1) // 带缓冲区的通道, 允许写入和读出
		end       = make(chan struct{})
		exts      = []string{"mp3", "wav", "wma", "ape"}
	)

	defer func() {
		if err != nil {
			panic(err)
		}
	}()

	fs, err := ioutil.ReadDir(conf.DIR_ASSETS)
	if err != nil {
		return
	}
	names := make([]string, 0, len(fs))
	_, list, err := app.GetIgnoreDetail()
	if err != nil {
		return
	}
	for _, f := range fs {

		name := f.Name()
		fmt.Println(name)
		// fmt.Println(f.Sys())
		if !util.Contains(exts, strings.TrimLeft(path.Ext(name), ".")) {
			continue
		}
		if !util.Contains(list, name) {
			names = append(names, name)
		}
	}

	go app.Ctrl(end, ch, currIndex, names)

	fmt.Println("pid::", os.Getpid())
	currIndex <- app.Control(ch, 0, names)
	fmt.Println("currIndex first::", currIndex)

	// play(names[0], ch)
	<-end
}
