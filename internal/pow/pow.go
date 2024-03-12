package pow

import (
	"errors"
)

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrNilConfig      = errors.New("nil config")
	ErrNilConnection  = errors.New("nil connection")
	ErrTimeout        = errors.New("io timeout")
	ErrGetRequest     = errors.New("wrong request")
	ErrGeneratePuzzle = errors.New("error generating puzzle")
	ErrGetSolution    = errors.New("error get solution")
	ErrWrongAnswer    = errors.New("client got wrong  answer")
)

type powConf interface {
	GetNumberOfTimesRAppliesF() int64
	GetSizeOfEachValue() int64
}

type metricker interface {
	GetDifficulties() uint
}

type pow struct {
	sizeOfEachValue        int64
	numberOfTimesRAppliesF int64
}

func (m *pow) GeneratePuzzle() (uint64, *Puzzle, error) {
	x0, puzzle, err := generatePuzzle(Algosin, m.sizeOfEachValue, m.numberOfTimesRAppliesF)
	if err != nil {
		return 0, nil, err
	}

	return x0, puzzle, nil
}

func NewPoW(
	cfg powConf,
) (*pow, error) {

	if cfg == nil {
		return nil, ErrNilConfig
	}

	p := &pow{
		numberOfTimesRAppliesF: cfg.GetNumberOfTimesRAppliesF(),
		sizeOfEachValue:        cfg.GetSizeOfEachValue(),
	}
	return p, nil
}
