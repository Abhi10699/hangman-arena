package player

type GuessResult struct {
	GuessedCharacter string
	Metadata         map[string]any
}

type Player interface {
	GuessCharacter(incorrectGuesses []string, currentGuess string) GuessResult
	ResetState()
}
