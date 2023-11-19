package util

import (
	"bytes"
	"io"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func Gbk2Utf8(bts []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(bts), simplifiedchinese.GBK.NewDecoder())
	return io.ReadAll(reader)
}
