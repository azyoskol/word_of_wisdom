package server

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"net"

	"github.com/azyoskol/word_of_wisdom/internal/pow"
)

var (
	errAddr          = errors.New("addr must be set like 0.0.0.0:8000")
	errHandler       = errors.New("handler must be set")
	errErrHandler    = errors.New("ErrHandler must be set")
	errRejectHandler = errors.New("RejectHandler must be set")
)

type challange interface {
	GeneratePuzzle() (uint64, *pow.Puzzle, error)
}

type Server struct {
	Addr          string
	WorkerPool    WorkerPool
	Challanger    challange
	Handler       func(net.Conn, challange)
	ErrHandler    func(net.Conn, string)
	RejectHandler func(net.Conn, string)
}

func (s *Server) Startup(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	if s.Addr == "" {
		return errAddr
	}
	if s.Handler == nil {
		return errHandler
	}
	if s.ErrHandler == nil {
		return errErrHandler
	}
	if s.RejectHandler == nil {
		return errRejectHandler
	}
	if s.WorkerPool == nil {
		wp, err := DefaultWorkerPool()
		if err != nil {
			return err
		}
		s.WorkerPool = wp
	}

	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	g.Go(func() error {
		<-ctx.Done()
		defer l.Close()
		return nil
	})

	g.Go(func() error {
		for {
			conn, err := l.Accept()
			if err != nil {
				s.ErrHandler(conn, err.Error())
				continue
			}
			w, err := s.WorkerPool.Get()
			if err != nil {
				s.RejectHandler(conn, err.Error())
				continue
			}
			w.Work(s, conn, ctx, g)
		}
	})

	err = g.Wait()
	if err != nil {
		return err
	}

	return nil
}
