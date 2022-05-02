package conf

import (
	"strings"
)

var prefix = "-"

// - xxxx -> xxxx
func ParsePrefix(name string) string {
	return strings.TrimLeft(name, prefix)
}

// - xxxx -> fff/ggg/xxx.mp3
func FilePath(name string) string {
	return names[Index(name)]
}

// 文件名
func FileName(name string) string {
	strs := strings.Split(name, "/")
	return prefix + strs[len(strs)-1]
}

// 获取当前歌曲索引地址
func Index(name string) int {
	for i, n := range names {
		if name == FileName(n) {
			return i
		}
	}
	return -1
}
