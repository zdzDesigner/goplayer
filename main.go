package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"player/util"
	"regexp"
	"strings"
	"time"

	"github.com/faiface/beep"

	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

// DIR_ASSETS = fmt.Sprintf("%s/source/assets/enya", strings.TrimRight(util.Dir(), "/palyer"))

// DIR_ASSETS = "/home/zdz/temp/xuwei-resource/默认"
// DIR_ASSETS = "/home/zdz/temp/yingwenjingdian-resource"
// DIR_ASSETS = "/home/zdz/temp/like-resource"
var DIR_ASSETS = "/home/zdz/temp/ape-resource" // DIR_ASSETS = "/home/zdz/temp/qyy-resource"
// var DIR_ASSETS = "/home/zdz/temp/ape-resource2" // DIR_ASSETS = "/home/zdz/temp/qyy-resource"
// DIR_ASSETS = "/home/zdz/temp/zhoujielun-resource"

func main() {
	var (
		err       error
		currIndex = make(chan int)
		ch        = make(chan struct{}, 1) // 带缓冲区的通道, 允许写入和读出
		end       = make(chan struct{})
		exts      = []string{"mp3", "wav", "wma", "ape"}
	)
	// strings.Join([]string{"aaa"}, "|")
	defer func() {
		if err != nil {
			panic(err)
		}
	}()

	fs, err := ioutil.ReadDir(DIR_ASSETS)
	if err != nil {
		return
	}
	names := make([]string, 0, len(fs))
	_, list, err := getIgnoreDetail()
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

	go ctrl(end, ch, currIndex, names)

	fmt.Println("pid::", os.Getpid())
	currIndex <- control(ch, 0, names)
	fmt.Println("currIndex first::", currIndex)

	// play(names[0], ch)
	<-end
}

func control(trigger chan struct{}, index int, names []string) int {
	// TODO::bug p1
	defer func() {
		if err := recover(); err != nil {
			index = index + 1
			control(trigger, index, names)
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

	// if play(names[util.Random(len(names))], trigger) {
	if play(names[index], trigger) {
		return control(trigger, index+1, names)
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

// 获取资源
func getSource(name string) (source string, err error) {
	source = fmt.Sprintf("%s/%s", DIR_ASSETS, name)
	fmt.Println("source::", source)
	// TODO::获取采样率
	stdout, err := exec.Command("ffprobe", source).CombinedOutput()
	fmt.Println("err::", err, string(stdout))
	if err != nil {
		return
	}
	// 获取采样率正常
	regHz := regexp.MustCompile(`22050`)
	regType := regexp.MustCompile(`\.ape|\.wma`)
	fmt.Println(regHz.MatchString(string(stdout)), regType.MatchString(string(stdout)))
	if !regHz.MatchString(string(stdout)) && !regType.MatchString(string(stdout)) {
		return
	}
	log.Println("SR too low 22050")
	regExt := regexp.MustCompile(`\..*$`)
	// TODO::转换采样率和比特率
	// newsource := fmt.Sprintf("%s/%s", DIR_ASSETS, strings.Replace(name, ".", "c.", -1))
	newsource := fmt.Sprintf("%s/%s", DIR_ASSETS, regExt.ReplaceAllString(name, ".mp3"))
	fmt.Println(newsource)
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

// 获取忽略歌名
func getIgnoreDetail() (*os.File, []string, error) {
	var (
		err  error
		list []string
	)
	defer func() {
		if err != nil {
			fmt.Println(err)
		}
	}()

	encoded := base64.URLEncoding.EncodeToString([]byte(DIR_ASSETS))
	// fmt.Println(encoded, util.Dir())
	faddr := fmt.Sprintf("%s/assets/%s.json", util.Dir(), encoded)
	f, err := os.OpenFile(faddr, os.O_RDWR, 0)
	if err != nil {
		f, err = os.Create(faddr)
		f.WriteString("[]")
	}
	if err != nil {
		return nil, nil, err
	}

	txt, err := ioutil.ReadFile(faddr)
	fmt.Println(txt)
	if err = json.Unmarshal(txt, &list); err != nil {
		return nil, nil, err
	}
	fmt.Println(list)
	return f, list, nil
}

func ctrl(end chan struct{}, ch chan struct{}, currIndex chan int, names []string) {
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
				delSong(names[i])
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
				currIndex <- control(ch, i, names)
				fmt.Println("currIndex::", i)
			}()
		}

	}
}

// 删除资源
func delSong(name string) {
	var err error
	defer func() {
		if err != nil {
			fmt.Println("del song err::", err)
		}
	}()

	f, list, err := getIgnoreDetail()
	list = append(list, name)
	bts, err := json.Marshal(list)
	if err != nil {
		return
	}
	fmt.Println("bts::", string(bts))
	_, err = f.Write(bts)
	if err != nil {
		return
	}
}
