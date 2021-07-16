package util

import (
	"os"
)

var dir = ""

func Dir() string {
	var err error
	if dir != "" {
		return dir
	}

	defer func() {
		if err != nil {
			panic(err)
		}
	}()
	dir, err := os.Getwd()
	return dir

}
