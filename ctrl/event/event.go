package event

import "sync"

var ChooseName = make(chan string)

var Evt PubSuber

func init() {
	Evt = NewPubSub()
}

func NewPubSub() PubSuber {
	return &PubSub{
		pool: make(map[string][]Handler),
	}
}

func StringVal(it any) string {
	return it.(string)
}

type next struct {
	name  string
	index int
}

func NewNext(name string, index int) *next {
	return &next{name, index}
}

func NextVal(it any) (string, int) {
	n := it.(*next)
	return n.name, n.index
}

type H interface {
	string | next
}

type Handler func(any)

type PubSuber interface {
	On(string, Handler)
	Emit(string, any)
}

type PubSub struct {
	sync.Mutex
	pool map[string][]Handler
}

func (ps *PubSub) On(key string, handler Handler) {
	ps.Lock()
	defer ps.Unlock()
	ps.pool[key] = append(ps.pool[key], handler)
}

func (ps *PubSub) Emit(key string, val any) {
	ps.Lock()
	handlers := ps.pool[key]
	ps.Unlock()
	if handlers == nil {
		return
	}
	// copyHands := make([]Handler, len(handlers))
	// copy(copyHands, handlers)
	for _, handler := range handlers {
		handler(val)
	}
}
