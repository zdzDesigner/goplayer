package conf

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"player/util"
	"regexp"
	"strings"
)

var AUDIO_NAMES []string // 歌曲播放地址

// 获取歌曲列表
func List() []string {
	var err error
	if len(AUDIO_NAMES) > 0 {
		return AUDIO_NAMES
	}
	AUDIO_NAMES, err = AudioList()
	util.MustNoErr(err)
	return AUDIO_NAMES
}

func UpdateList() []string {
	var err error
	AUDIO_NAMES, err = AudioList()
	util.MustNoErr(err)

	return AUDIO_NAMES
}

// 音频列表
func AudioList() (names []string, err error) {
	_, list, err := GetIgnoreDetail()
	if err != nil {
		return
	}
	if len(AUDIO_NAMES) > 0 {
		names = cacheAudioList(list)
	} else {
		names = runtimeAudioList(list)
	}
	// fmt.Println("names:",names)

	if len(names) == 0 {
		err = errors.New("no music file in current dir")
	}
	return
}

func cacheAudioList(list []string) (names []string) {
	for _, filepath := range AUDIO_NAMES {
		name := regexp.MustCompile(`^.*/`).ReplaceAllString(filepath, "")
		if !util.Contains(list, name) {
			names = append(names, filepath)
		}
	}
	return
}

func runtimeAudioList(list []string) (names []string) {
	var queueNames []string
	err := deepDir(DIR_ASSETS, func(name, dir string) {
		if !validExt(name) {
			return
		}
		if !util.Contains(list, name) {
			// name = regexp.MustCompile(`（`).ReplaceAllString(name, "(")
			// name = regexp.MustCompile(`）`).ReplaceAllString(name, ")")
			queueNames = append(queueNames, fmt.Sprintf("%s/%s", dir, name))
		}
	})
	if err != nil {
		panic(err)
	}
	for _, index := range util.RandomMutil(len(queueNames), 0, len(queueNames)-1) {
		names = append(names, queueNames[index])
	}
	// fmt.Println("names:", names)
	return
}

func validExt(name string) bool {
	exts := []string{"mp3", "wav", "wma", "ape"}
	return util.Contains(exts, strings.TrimLeft(path.Ext(name), "."))
}

// 深度递归文件夹
func deepDir(dirpath string, fn func(string, string)) (err error) {
	// fs, err := ioutil.ReadDir(dir)
	dirs, err := os.ReadDir(dirpath)
	if err != nil {
		return
	}

	for _, dir := range dirs {
		f, err := dir.Info()
		if err != nil {
			continue
		}
		if f.IsDir() {
			if deepDir(fmt.Sprintf("%s/%s", dirpath, f.Name()), fn) != nil {
				continue
			}
		}
		fn(f.Name(), dirpath)
	}
	return
}

// 获取忽略歌名
func GetIgnoreDetail() (addr string, list []string, err error) {
	defer func() {
		if err != nil {
			fmt.Println(err)
		}
	}()

	encoded := base64.URLEncoding.EncodeToString([]byte(DIR_ASSETS))
	// fmt.Println(encoded, util.Dir())
	addr = fmt.Sprintf("%s/assets/%s.json", util.Dir(), encoded)
	f, err := os.OpenFile(addr, os.O_RDWR, 0)
	if err != nil {
		f, err = os.Create(addr)
		util.MustNoErr(err)
		f.WriteString("[]")
	}
	defer f.Close()

	// txt, err := ioutil.ReadFile(addr)
	txt, err := os.ReadFile(addr)
	util.MustNoErr(json.Unmarshal(txt, &list))

	// fmt.Println(list)
	return addr, list, nil
}

// 删除资源
func DelSong(name string) {
	var err error
	defer func() {
		if err != nil {
			fmt.Println("del song err::", err)
		}
	}()

	addr, list, err := GetIgnoreDetail()
	list = append(list, name)
	bts, err := json.Marshal(list)
	if err != nil {
		return
	}
	// fmt.Println("bts::", string(bts))

	f, err := os.OpenFile(addr, os.O_RDWR, 0)
	if err != nil {
		return
	}
	_, err = f.Write(bts)
}
