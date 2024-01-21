package core

type CallableFunction func(args ...any)

type Signal struct {
	subscribers_ []CallableFunction
}

func NewSignal() *Signal {
	return &Signal{}
}
func (c *Signal) Subscribe(fn CallableFunction) CallableFunction {
	c.subscribers_ = append(c.subscribers_, fn)
	return fn
}

func (c *Signal) Notify(args ...any) {
	for _, v := range c.subscribers_ {
		v(args...)
	}
}
