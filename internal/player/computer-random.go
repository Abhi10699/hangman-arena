package player

import (
	"math/rand"
	"slices"
	"strings"
)

var alphabets = strings.Split("abcdefghijklmnopqrstuvwxyz", "")

type RandomCharacterGuesser struct{}

func (r *RandomCharacterGuesser) GuessCharacter(incorrectGuesses []string, currentGuess string) GuessResult {

	randomCharIndex := rand.Intn(len(alphabets))
	guessedChar := alphabets[randomCharIndex]

	for {
		isInIncorrectGuesses := slices.Contains(incorrectGuesses, guessedChar)
		isInCorrectGuesses := strings.Contains(currentGuess, guessedChar)

		if !(isInCorrectGuesses || isInIncorrectGuesses) {
			// we found a unique character
			break
		}

		randomCharIndex = rand.Intn(len(alphabets))
		guessedChar = alphabets[randomCharIndex]
	}

	return GuessResult{
		GuessedCharacter: guessedChar,
		Metadata:         nil,
	}
}

func (r *RandomCharacterGuesser) ResetState() {}
