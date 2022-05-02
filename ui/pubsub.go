package ui

var PS PubSuber

func init() {
	PS = NewPubSub()
}

type Handler func(string)

type PubSuber interface {
	On(string, Handler)
	Emit(string, ...string)
}

type PubSub struct {
	pool map[string][]Handler
}

func (ps *PubSub) On(key string, handler Handler) {
	ps.pool[key] = append(ps.pool[key], handler)
}

func (ps *PubSub) Emit(key string, strs ...string) {
	handlers := ps.pool[key]
	if handlers == nil {
		return
	}
	for _, handler := range handlers {
		handler(strs[0])
	}
}

func NewPubSub() PubSuber {
	return &PubSub{
		pool: make(map[string][]Handler),
	}
}

func Do() {
	pubsub := NewPubSub()

	pubsub.On("key", func(str string) {
		Log(str)
	})

	pubsub.Emit("key", "handler ccc")
}
