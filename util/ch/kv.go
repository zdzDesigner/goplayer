package ch

type ChKV struct {
	val  map[string]string
	sign chan struct{}
}

func (c *ChKV) Set(val map[string]string) {
	c.val = val
	c.sign <- struct{}{}
}

func (c *ChKV) Get(key string) string {
	return c.val[key]
}

func (c *ChKV) Done() chan struct{} {
	return c.sign
}

func NewChKV() *ChKV {
	chkv := &ChKV{val: make(map[string]string), sign: make(chan struct{})}
	return chkv
}
