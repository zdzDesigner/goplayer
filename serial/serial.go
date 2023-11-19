package serial

import (
	"bufio"
	"encoding/json"

	"fmt"

	"github.com/tarm/serial"
)

type LineHookFunc func([]byte)
type Serial struct {
	client      *serial.Port
	transceiver *Transceive[int]
	isconnected bool
}

type Config struct {
	Name     string
	Baud     int
	Parity   serial.Parity
	StopBits serial.StopBits
}

func NewSerialer(session string, config *Config) (*Serial, error) {
	// fmt.Println("-------- NewSerialer -------")
	client, err := serial.OpenPort(&serial.Config{Name: config.Name, Baud: config.Baud, Parity: config.Parity, StopBits: config.StopBits})
	if err != nil {
		return nil, err
	}
	// fmt.Println("serial client", client)

	transceiver := NewTransceive([]int{-1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 20}, session)

	serialer := &Serial{client: client, transceiver: transceiver}
	go serialer.Listen()

	return serialer, nil
}

func (s *Serial) IsConnected() bool { return s.isconnected }
func (s *Serial) Disconnect(delay uint) error {
	err := s.client.Close()
	if err != nil {
		return err
	}
	s.isconnected = false
	return nil
}
func (s *Serial) Listen() {
	// listen(m.client, topic, m.transceiver.Listen)
	listen(s.client, func(res []byte) {
		// fmt.Println(string(res))
		// return
		s.transceiver.Listen("", res)
	})
}

// topic 暂时无用
func (s *Serial) Request(topic string, payload map[string]any) (res string, err error) {
	return s.transceiver.Request(topic, payload, s.Send)
}
func (s *Serial) Send(topic string, data any) error {
	bts, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return send(s.client, bts)
}

func NewSerial() (*serial.Port, error) {
	config := &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200}
	return serial.OpenPort(config)
}

func send(port *serial.Port, payload []byte) error {
	_, err := port.Write(append(payload, '\n'))
	return err
}

func listen(port *serial.Port, lineHook LineHookFunc) error {
	// fmt.Println("read serial")
	buf := bufio.NewReader(port)
	for {
		// fmt.Println(buf.ReadByte())
		// fmt.Println(buf.ReadByte())
		// fmt.Println(buf.ReadByte())
		// fmt.Println(buf.ReadByte())
		// line, err := buf.ReadString('\n')
		// line = strings.TrimSpace(line)
		line, _, err := buf.ReadLine()
		if err != nil {
			fmt.Println(err)
			return err
		}
		// fmt.Println("-----------", line)
		lineHook(line)
		// fmt.Println(string(line))
	}
}
