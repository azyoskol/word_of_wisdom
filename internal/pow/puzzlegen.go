package pow

import (
	"crypto/rand"
	"math/big"
)

type Puzzle struct {
	Xk uint64

	K int64

	N int64

	Checksum string
}

func randUint64(ceil uint64) (uint64, error) {
	val, err := rand.Int(rand.Reader, big.NewInt(int64(ceil)))
	if err != nil {
		return 0, err
	}
	return val.Uint64(), nil
}

func generatePuzzle(
	algo func(uint64, float64) uint64,
	sizeOfEachValue, numberOfTimesRAppliesF int64,
) (uint64, *Puzzle, error) {
	max := (uint64(1) << sizeOfEachValue) - 1
	maxF := float64(max)

	x0, err := randUint64(max)
	if err != nil {
		return 0, nil, err
	}

	seq := make([]uint64, numberOfTimesRAppliesF+1)
	seq[numberOfTimesRAppliesF] = x0

	xk := x0
	for index := uint64(1); index <= uint64(numberOfTimesRAppliesF); index++ {
		xk = algo(xk, maxF) ^ index
		seq[uint64(numberOfTimesRAppliesF)-index] = xk
	}
	checkSum := checksum(seq)

	return x0, &Puzzle{
		Xk:       xk,
		K:        numberOfTimesRAppliesF,
		N:        sizeOfEachValue,
		Checksum: checkSum,
	}, nil
}
