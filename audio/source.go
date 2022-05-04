package audio

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

// 获取资源
func getSource(source string) (nsource string, err error) {
	nsource = source
	// source := fmt.Sprintf("%s/%s", conf.DIR_ASSETS, name)
	// fmt.Println("source::", source)
	// TODO::获取采样率
	stdout, err := exec.Command("ffprobe", source).CombinedOutput()
	// fmt.Println("err::", err, string(stdout))
	if err != nil {
		return
	}
	// 获取采样率正常
	regHz := regexp.MustCompile(`22050`)
	regType := regexp.MustCompile(`\.ape|\.wma`)
	// fmt.Println(regHz.MatchString(string(stdout)), regType.MatchString(string(stdout)))
	if !regHz.MatchString(string(stdout)) && !regType.MatchString(string(stdout)) {
		return
	}
	// log.Println("SR too low 22050")
	regExt := regexp.MustCompile(`\..*$`)
	// TODO::转换采样率和比特率
	// nsource := fmt.Sprintf("%s/%s", DIR_ASSETS, strings.Replace(name, ".", "c.", -1))
	nsource = fmt.Sprintf("%s/%s", filepath.Dir(source), regExt.ReplaceAllString(filepath.Base(source), ".mp3"))
	// fmt.Println("new source::", nsource)
	stdout, err = exec.Command("ffmpeg", "-i", source, "-ar", "44100", "-ab", "128k", nsource).CombinedOutput()
	if err != nil {
		return
	}
	// 删除老文件
	err = os.Remove(source)
	if err != nil {
		return
	}
	// fmt.Println(string(stdout))
	return nsource, nil
}
