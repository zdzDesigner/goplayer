package ctrl

import (
	"fmt"
	"io/ioutil"
	"player/ctrl/event"
	"player/ui"
	"player/util"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type IsWatch struct {
	m  sync.Mutex
	is bool
}

func ListenGlobal() {
	watchFile := fmt.Sprintf("%s/.fsnotify", util.Dir())
	// 创建 watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()

	wait := make(chan bool)
	// count := 0
	go func() {
		iswatch := IsWatch{is: false}
		for {
			select {
			case evt, ok := <-watcher.Events:
				if !ok && evt.Op != fsnotify.Write {
					break
				}
				val, err := ioutil.ReadFile(watchFile)
				if err != nil {
					break
				}
				iswatch.m.Lock()
				if iswatch.is {
					iswatch.m.Unlock()
					break
				}
				iswatch.m.Unlock()
				if strings.Contains(string(val), "next") {
					iswatch.m.Lock()
					event.Evt.Emit("AUDIO_CTRL", "NEXT")
					// count = count + 1
					// fmt.Println(string(val), count)
					iswatch.is = true
					iswatch.m.Unlock()
					go func() {
						time.Sleep(time.Microsecond * 100)
						iswatch.m.Lock()
						iswatch.is = false
						iswatch.m.Unlock()
					}()
				}

			case _, err := <-watcher.Errors:
				ui.Log(err)
			}
		}
	}()

	watcher.Add(watchFile)

	<-wait
}
