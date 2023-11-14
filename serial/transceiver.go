package serial

import (
	"encoding/json"
	"errors"
	// "fmt"
	"math"
	"time"
)

type Clienter interface {
	Request(topic string, payload map[string]any) (string, error)
	Disconnect(delay uint) error
	IsConnected() bool
}

// 发送
type sendFunc func(topic string, payload any) error

type Transceive[T int | string] struct {
	session string
	res     map[T]chan string
	err     map[T]chan error
}

func NewTransceive[T int | string](keys []T, session string) *Transceive[T] {
	transceiver := &Transceive[T]{session: session}
	transceiver.init(keys)
	return transceiver
}

func (t *Transceive[T]) init(keys []T) {
	t.res = make(map[T]chan string)
	t.err = make(map[T]chan error)
	for _, i := range keys {
		t.res[i] = make(chan string, 3)
		t.err[i] = make(chan error, 3)
	}
}

// 监听
func (t *Transceive[T]) Listen(topic string, res []byte) {
	cmd, data, err := parse(t.session, res)
	// fmt.Println(data, err)
	if err != nil {
		t.err[T(cmd)] <- err
		return
	}

	t.res[T(cmd)] <- data

}
func (t *Transceive[T]) Request(topic string, payload map[string]any, send sendFunc) (res string, err error) {
	cmd := T(0)
	if payload != nil {
		cmd = payload["cmd"].(T)
	}
	//
	// go send(topic, payload)
	select {
	case res = <-t.res[cmd]:
		// fmt.Println(cmd, "::", res)
		return res, nil
	case err = <-t.err[cmd]:
		return "", err
	case <-time.After(1 * time.Second):
		return "", errors.New("response timeout")
	}

}

// ========== 解析res

func parse(session string, data []byte) (int, string, error) {
	var res []int
	if err := json.Unmarshal(data, &res); err != nil {
		return 0, "", err
	}

	// fmt.Println(res)
	ret := ""
	h := res[0]
	v := res[1]

	if math.Abs(float64(h-130)) > math.Abs(float64(v-130)) {
		if h < 20 {
			ret = "left"
		}
		if h > 230 {
			ret = "right"
		}
	} else {
		if v < 20 {
			ret = "top"
		}
		if v > 230 {
			ret = "bottom"
		}
    ret = "reset"
	}

	return 0, ret, nil
}
