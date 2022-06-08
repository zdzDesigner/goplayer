package conf

import (
	"strings"
)

var prefix = "-"

// - xxx -> fff/ggg/xxx.mp3
func FilePath(name string) string {
	return AUDIO_NAMES[Index(name)]
}

// - xxxx -> xxxx
func ClearPrefix(name string) string {
	return strings.TrimLeft(name, prefix)
}

// 带前缀的文件名
func PrifixFileName(name string) string {
	return prefix + FileName(name)
}

// 文件名
func FileName(name string) string {
	strs := strings.Split(name, "/")
	return strs[len(strs)-1]
}

// 获取当前歌曲索引地址
func Index(name string) int {
	for i, n := range AUDIO_NAMES {
		if ClearPrefix(FileName(name)) == FileName(n) {
			return i
		}
	}
	return -1
}

// 下一个
func NextIndex(name string) int {
	index := Index(name) + 1
	if index == len(AUDIO_NAMES) {
		index = 0
	}
	return index
}
