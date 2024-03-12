package quote

import (
	"errors"
	"math/rand"
)

var (
	errRandom = errors.New("ops... something go wrong")
)

func GetQuote() (string, error) {
	var quotes []string

	quotes = append(quotes, "If you were waiting for the opportune moment, that was it.")
	quotes = append(quotes, "This guy is too good he shoots fire from his brain")
	quotes = append(quotes, "Everyone is a monster to someone. Since you are so convinced that I am yours. I will be it.")
	quotes = append(quotes, "It took me five tries.")
	quotes = append(quotes, "No, you clearly don't know who you're talking to, so let me clue you in. I am not in danger, Skyler. I am the danger! A guy opens his door and gets shot and you think that of me? No. I am the one who knocks!")
	quotes = append(quotes, "I am running away from my responsibilities. And it feels good.")

	if rand.Intn(1000)%5 == 0 {
		return "", errRandom
	}

	return quotes[rand.Intn(len(quotes))], nil
}
