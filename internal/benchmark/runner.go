package benchmark

import (
	"fmt"
	"hangman/internal/game"
	"hangman/internal/player"
)

type BenchmarkStats struct {
	Wins   int
	Losses int
}

type BenchmarkRunner struct {
	Player     player.Player
	WordSource Datasource
	TotalGames int

	Stats BenchmarkStats
}

func NewBenchmarkRunner(word Datasource, totalGames int, p player.Player) *BenchmarkRunner {
	return &BenchmarkRunner{
		Player:     p,
		WordSource: word,
		TotalGames: totalGames,
		Stats: BenchmarkStats{
			Wins:   0,
			Losses: 0,
		},
	}
}

func (runner *BenchmarkRunner) RunBenchmark() {
	for gameIndex := range runner.TotalGames {

		fmt.Println("Game: ", gameIndex)

		var attempts []AttemptLog

		wordToGuess, err := runner.WordSource.GetNextWord()
		if err != nil || len(wordToGuess) <= 1 {
			fmt.Errorf("No valid words to test")
			break
		}

		gameSession := game.NewSession(wordToGuess)
		for {

			if gameSession.IsGameOver() {
				// check if the game ended in a win or loss
				fmt.Println(gameSession.CurrentGuess)

				result := "loss"
				if gameSession.CurrentGuess == gameSession.TargetWord {
					runner.Stats.Wins++
					result = "win"
				} else {
					runner.Stats.Losses++
				}

				logData := GameLog{
					Word:          wordToGuess,
					GameIndex:     gameIndex,
					Result:        result,
					TotalAttempts: len(attempts),
					Attempts:      attempts,
				}

				if err := LogGame(logData, "results/benchmark_log.jsonl"); err != nil {
					fmt.Println("Error logging game:", err)
				}

				break
			}

			// Capture state before guess
			currentGuess := gameSession.CurrentGuess
			incorrectChars := make([]string, len(gameSession.IncorrectChars))
			copy(incorrectChars, gameSession.IncorrectChars)

			guessResult := runner.Player.GuessCharacter(gameSession.IncorrectChars, gameSession.CurrentGuess)

			attempts = append(attempts, AttemptLog{
				CurrentGuess:   currentGuess,
				IncorrectChars: incorrectChars,
				Guess:          guessResult.GuessedCharacter,
				Metadata:       guessResult.Metadata,
			})

			incorrectGuessesErr := gameSession.MakeGuess(guessResult.GuessedCharacter)
			if incorrectGuessesErr != nil {
				// do nothing
			}
		}

		runner.Player.ResetState()
	}
}
