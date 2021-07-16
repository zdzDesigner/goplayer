package util

import (
	"fmt"
	"io/ioutil"
	"os"
)

// LocalFileByte ..
func LocalFileByte(filepath string) ([]byte, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(fmt.Sprintf("%s%s", dir, filepath))
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return bs, nil
}
