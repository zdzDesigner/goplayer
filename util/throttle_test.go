package util

import (
	"fmt"
	"testing"
	"time"
)

func Test_NewThrottler(t *testing.T) {
	throttler := NewThrottler(time.Millisecond * 5500)
	throttler.Do(func() {
		fmt.Println("vvvvv")
		t.Error("xxxx")
		t.Errorf("Error:%s\n", "vvvvv")
	}, true)

	ch := make(chan struct{})
	go func() {}()
	<-ch

}

// go test -v /home/zdz/Documents/Try/Go/music/player/util/ -count=1 -test.run Test_NewThrottler
