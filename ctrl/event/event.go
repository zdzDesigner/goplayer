package event

import "sync"

var Evt PubSuber

func init() {
	Evt = NewPubSub()
}

func StringVal(it interface{}) string {
	return it.(string)
}

type next struct {
	name  string
	index int
}

func NewNext(name string, index int) *next {
	return &next{name, index}
}

func NextVal(it interface{}) (string, int) {
	n := it.(*next)
	return n.name, n.index
}

type Handler func(interface{})

type PubSuber interface {
	On(string, Handler)
	Emit(string, interface{})
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

func (ps *PubSub) Emit(key string, val interface{}) {
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

func NewPubSub() PubSuber {
	return &PubSub{
		pool: make(map[string][]Handler),
	}
}
