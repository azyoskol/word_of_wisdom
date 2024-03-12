package state

import (
	"net"

	"github.com/azyoskol/word_of_wisdom/internal/pow"
)

type Context struct {
	state  State
	Conn   net.Conn
	Answer uint64
	Puzzle pow.Puzzle
	Quote  string
}

type State interface {
	Handle(*Context) error
}

func (c *Context) SetState(state State) {
	c.state = state
}

func (c *Context) Do() error {
	return c.state.Handle(c)
}
