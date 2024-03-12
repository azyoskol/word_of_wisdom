package server

import (
	"context"
	"net"

	"golang.org/x/sync/errgroup"
)

type Worker interface {
	Work(*Server, net.Conn, context.Context, *errgroup.Group)
}

type defaultWorker struct {
	conn    chan net.Conn
	pos     uint
	server  *Server
	handler func(net.Conn, challange)
	ctx     context.Context
	g       *errgroup.Group
}

func (w *defaultWorker) Work(srv *Server, conn net.Conn, ctx context.Context, g *errgroup.Group) {
	w.server = srv
	w.ctx = ctx
	w.g = g
	w.handler = srv.Handler
	w.conn <- conn
}

func (w *defaultWorker) run() {
	go func() {
		for {
			select {
			case c := <-w.conn:
				w.handler(c, w.server.Challanger)
				w.server.WorkerPool.Put(w)
			}
		}
	}()
}

func newWorker(i uint) *defaultWorker {
	return &defaultWorker{
		conn: make(chan net.Conn),
		pos:  i,
	}
}
