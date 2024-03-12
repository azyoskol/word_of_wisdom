package server

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
	"os"

	"github.com/azyoskol/word_of_wisdom/internal/log"
	"github.com/azyoskol/word_of_wisdom/internal/pow"
	"github.com/azyoskol/word_of_wisdom/internal/quote"
	"github.com/azyoskol/word_of_wisdom/internal/state"
	"github.com/azyoskol/word_of_wisdom/utils"
)

var (
	errNilConfig     = errors.New("configuration is nil")
	errNilConnection = errors.New("connection is nil")
)

type configurer interface {
	GetAddress() string
}

func NewPowServer(cfg configurer, c challange) (*Server, error) {
	if cfg == nil {
		return nil, errNilConfig
	}
	wp, err := DefaultWorkerPool(1, 2)
	if err != nil {
		log.Fatalw("Server", log.M{"err": err})
	}

	server := &Server{
		Addr:          cfg.GetAddress(),
		Challanger:    c,
		Handler:       powHandler,
		ErrHandler:    powErrHandler,
		RejectHandler: powRejectHandler,
		WorkerPool:    wp,
	}

	return server, nil
}

func powHandler(conn net.Conn, c challange) {
	defer conn.Close()
	protocolCtx := state.Context{
		Conn: conn,
	}

	protocolCtx.SetState(&CreateChallange{
		c,
	})
	err := protocolCtx.Do()
	if err != nil {
		log.Fatalw("Can't create challange", log.M{"err": err})
	}
	protocolCtx.SetState(&ValidateChallange{})
	err = protocolCtx.Do()
	if err != nil {
		log.Fatalw("Can't validate challange", log.M{"err": err})
	}

	protocolCtx.SetState(&SendWordOfWisdom{})
	err = protocolCtx.Do()
	if err != nil {
		log.Errorw("error write to client", log.M{"err": err})
	}
}

func powErrHandler(conn net.Conn, err string) {
	if conn != nil {
		conn.Close()
	}

	log.Fatalw("run server error", log.M{"err": err})
}

func powRejectHandler(conn net.Conn, err string) {
	if conn != nil {
		conn.Close()
	}
	log.Fatalw("reject connect error", log.M{"err": err})
}

type CreateChallange struct {
	challange challange
}

func (s *CreateChallange) Handle(c *state.Context) error {
	x0, puzzle, err := s.challange.GeneratePuzzle()
	if err != nil {
		return pow.ErrGeneratePuzzle
	}
	c.Answer = x0

	data, err := json.Marshal(puzzle)
	if err != nil {
		return pow.ErrGeneratePuzzle
	}

	err = utils.WriteMessage(c.Conn, data)
	if err != nil {
		return pow.ErrTimeout
	}
	return nil
}

type ValidateChallange struct{}

func (s *ValidateChallange) Handle(c *state.Context) error {
	data, err := utils.ReadMessage(c.Conn)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return pow.ErrTimeout
		}
		if os.IsTimeout(err) {
			return pow.ErrTimeout
		}
		return err
	}

	candidate := binary.BigEndian.Uint64(data)

	if candidate != c.Answer {
		return pow.ErrWrongAnswer
	}
	return nil
}

type SendWordOfWisdom struct {
}

func (s *SendWordOfWisdom) Handle(c *state.Context) error {
	if c.Conn == nil {
		return errNilConnection
	}
	quote, err := quote.GetQuote()
	if err != nil {
		return err
	}

	err = utils.WriteMessage(c.Conn, []byte(quote))
	if err != nil {
		return err
	}

	return nil
}
