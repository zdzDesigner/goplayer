package app

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"player/conf"
	"player/util"
	"regexp"
)

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

	encoded := base64.URLEncoding.EncodeToString([]byte(conf.DIR_ASSETS))
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
	fmt.Println("bts::", string(bts))
	_, err = f.Write(bts)
	if err != nil {
		return
	}
}

// 获取资源
func getSource(name string) (source string, err error) {
	source = fmt.Sprintf("%s/%s", conf.DIR_ASSETS, name)
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
	newsource := fmt.Sprintf("%s/%s", conf.DIR_ASSETS, regExt.ReplaceAllString(name, ".mp3"))
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
