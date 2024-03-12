package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"flag"
	"net"

	"github.com/azyoskol/word_of_wisdom/config"
	"github.com/azyoskol/word_of_wisdom/internal/log"
	"github.com/azyoskol/word_of_wisdom/internal/pow"
	"github.com/azyoskol/word_of_wisdom/internal/state"
	"github.com/azyoskol/word_of_wisdom/utils"
)

var (
	errGetPuzzle       = errors.New("can't reacive puzzle")
	errUnmarshalPuzzle = errors.New("can't unmarshal puzzle")
	errSolveChallange  = errors.New("can't solve challange")
	errSendSolution    = errors.New("can't send solution")
	errGetWordOfWisdom = errors.New("can't get word of wisdom")
)

func main() {
	cfg, err := config.NewServerConfig()
	if err != nil {
		log.Fatalw("Error loading configuration", log.M{"err": err})
	}
	netPath := flag.String("server", cfg.GetAddress(), "connect to server string")
	log.Infow("Connecting to server", log.M{"path": *netPath})

	conn, err := net.Dial("tcp", *netPath)
	if err != nil {
		log.Fatalw("unable to connect to server ",
			log.M{
				"err":  err,
				"path": *netPath,
			},
		)
	}
	ctx := state.Context{
		Conn: conn,
	}

	ctx.SetState(&RequestChallenge{})
	err = ctx.Do()
	if err != nil {
		log.Fatalw("unable to get puzzle",
			log.M{
				"err":  err,
				"path": *netPath,
			},
		)
	}

	ctx.SetState(&SolveChallenge{})
	err = ctx.Do()
	if err != nil {
		log.Fatalw("unable to solve puzzle",
			log.M{
				"err":  err,
				"path": *netPath,
			},
		)
	}

	ctx.SetState(&CheckChallange{})
	err = ctx.Do()
	if err != nil {
		log.Fatalw("error read quote from server",
			log.M{
				"err":  err,
				"path": *netPath,
			},
		)
	}

	log.Infow("got quote from server",
		log.M{"quote": ctx.Quote},
	)
}

type RequestChallenge struct{}

func (s *RequestChallenge) Handle(c *state.Context) error {
	data, err := utils.ReadMessage(c.Conn)
	if err != nil {
		return errGetPuzzle
	}

	err = json.Unmarshal(data, &c.Puzzle)
	if err != nil {
		return errUnmarshalPuzzle
	}

	return nil
}

type SolveChallenge struct{}

func (s *SolveChallenge) Handle(c *state.Context) error {
	solution, err := pow.Solution(&c.Puzzle, pow.Algosin)
	if err != nil {
		return errSolveChallange
	}

	reply := make([]byte, 8)
	binary.BigEndian.PutUint64(reply, solution)
	err = utils.WriteMessage(c.Conn, reply)
	if err != nil {
		return errSendSolution
	}

	return nil
}

type CheckChallange struct{}

func (s *CheckChallange) Handle(c *state.Context) error {

	data, err := utils.ReadMessage(c.Conn)
	if err != nil {
		return errGetWordOfWisdom
	}

	c.Quote = string(data)

	return nil
}
