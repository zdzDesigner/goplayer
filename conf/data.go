package conf

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"player/util"
	"strings"
)

var names []string // 歌曲播放地址

// 获取当前歌曲索引地址
func Index(name string) int {
	for i, n := range names {
		if name == n {
			return i
		}
	}
	return -1
}

// 获取歌曲列表
func List() []string {
	var err error
	if len(names) > 0 {
		return names
	}
	names, err = AudioList()
	if err != nil {
		panic("get song list fail")
	}
	return names
}

// 音频列表
func AudioList() (names []string, err error) {
	exts := []string{"mp3", "wav", "wma", "ape"}
	_, list, err := GetIgnoreDetail()

	deepDir(DIR_ASSETS, func(name, dir string) {
		if !util.Contains(exts, strings.TrimLeft(path.Ext(name), ".")) {
			return
		}
		if !util.Contains(list, name) {
			names = append(names, fmt.Sprintf("%s/%s", dir, name))
		}
	})

	return
}

// 深度递归文件夹
func deepDir(dir string, fn func(string, string)) (err error) {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, f := range fs {
		if f.IsDir() {
			if deepDir(fmt.Sprintf("%s/%s", dir, f.Name()), fn) != nil {
				continue
			}
		}
		fn(f.Name(), dir)
	}
	return
}

// 获取忽略歌名
func GetIgnoreDetail() (*os.File, []string, error) {
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
	// fmt.Println(txt)
	if err = json.Unmarshal(txt, &list); err != nil {
		return nil, nil, err
	}
	// fmt.Println(list)
	return f, list, nil
}

// 删除资源
func DelSong(name string) {
	var err error
	defer func() {
		if err != nil {
			fmt.Println("del song err::", err)
		}
	}()

	f, list, err := GetIgnoreDetail()
	list = append(list, name)
	bts, err := json.Marshal(list)
	if err != nil {
		return
	}
	// fmt.Println("bts::", string(bts))
	_, err = f.Write(bts)
	if err != nil {
		return
	}
}
