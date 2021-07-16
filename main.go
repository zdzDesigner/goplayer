package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"player/util"
	"regexp"
	"strings"
	"time"

	"github.com/faiface/beep"

	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var (
	DIR_ASSETS = fmt.Sprintf("%s/source/assets/enya", strings.TrimRight(util.Dir(), "/palyer"))
)

func main() {
	var err error
	defer func() {
		if err != nil {
			panic(err)
		}
	}()

	fs, err := ioutil.ReadDir(DIR_ASSETS)
	if err != nil {
		return
	}

	for _, f := range fs {
		fmt.Println(f.Name())
		// to::plan1
		play(f.Name())

		//

	}

	// f, err := os.(DIR_ASSETS)
	// fmt.Println(f, err)
	// fi, err := f.Readdir()
	// play()

}

func play(name string) bool {
	var err error
	defer func() {
		fmt.Println("play end")
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
	speaker.Init(bfmt.SampleRate, bfmt.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(stm, beep.Callback(func() {
		done <- true
	})))
	return <-done
}

// 获取资源
func getSource(name string) (source string, err error) {

	source = fmt.Sprintf("%s/%s", DIR_ASSETS, name)
	fmt.Println("source::", source)
	// TODO::获取采样率
	stdout, err := exec.Command("ffprobe", source).CombinedOutput()
	if err != nil {
		return
	}
	// 获取采样率正常
	reg := regexp.MustCompile(`22050`)
	if !reg.MatchString(string(stdout)) {
		return

	}
	log.Println("SR too low 22050")
	// TODO::转换采样率和比特率
	newsource := fmt.Sprintf("%s/%s", DIR_ASSETS, strings.Replace(name, ".", "c.", -1))
	stdout, err = exec.Command("ffmpeg", "-i", source, "-ar", "44100", "-ab", "128k", newsource).CombinedOutput()
	if err != nil {
		return
	}
	// 删除老文件
	err = os.Remove(source)
	if err != nil {
		return
	}
	fmt.Println(string(stdout))
	return newsource, nil

}
