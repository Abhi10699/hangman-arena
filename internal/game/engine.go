package game

import (
	"errors"
	"strings"
)

var (
	IncorrectGuess  = errors.New("Incorrect Guess")
	ValidationError = errors.New("ValidationError")
)

type HangmanGameSession struct {
	TargetWord     string
	CurrentGuess   string
	IncorrectChars []string
	Lives          int
}

func NewSession(targetWord string) *HangmanGameSession {
	return &HangmanGameSession{
		TargetWord:     strings.ToLower(targetWord),
		CurrentGuess:   strings.Repeat("_", len(targetWord)),
		IncorrectChars: []string{},
		Lives:          9,
	}
}

func (s *HangmanGameSession) MakeGuess(character string) error {
	gussedCharacter := character[0]
	found := false
	newGuess := []byte(s.CurrentGuess)

	for idx := 0; idx < len(s.TargetWord); idx++ {
		if s.TargetWord[idx] == gussedCharacter {
			newGuess[idx] = gussedCharacter
			found = true
		}
	}

	if !found {
		s.Lives--
		s.IncorrectChars = append(s.IncorrectChars, character)
		return IncorrectGuess
	}

	s.CurrentGuess = string(newGuess)
	return nil
}

func (s *HangmanGameSession) IsGameOver() bool {
	if s.CurrentGuess == s.TargetWord || s.Lives == 0 {
		return true
	}

	return false
}
