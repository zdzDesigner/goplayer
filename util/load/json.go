package load

import (
	"io/ioutil"
	"os"
)

type Json struct {
	path string
}

func (j *Json) Write(content string) (int, error) {
	file, err := os.OpenFile(j.path, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	return file.WriteString(content)
}
func (j *Json) Read() (string, error) {
	// var content = make([]byte, 0)
	file, err := os.OpenFile(j.path, os.O_RDONLY, 0666)
	// file, err := os.Open(j.path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	// _, err = file.Read(content)
	return string(content), err
}

func NewJson(path string) *Json {
	return &Json{path: path}
}
