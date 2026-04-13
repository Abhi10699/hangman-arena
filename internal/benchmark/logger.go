package benchmark

import (
	"encoding/json"
	"fmt"
	"os"
)

type AttemptLog struct {
	CurrentGuess   string         `json:"currentGuess"`
	IncorrectChars []string       `json:"incorrectChars"`
	Guess          string         `json:"guess"`
	Metadata       map[string]any `json:"metadata,omitempty"`
}

type GameLog struct {
	Word          string       `json:"word"`
	GameIndex     int          `json:"gameIndex"`
	Result        string       `json:"result"`
	TotalAttempts int          `json:"totalAttempts"`
	Attempts      []AttemptLog `json:"attempts"`
}

func LogGame(logData GameLog, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(logData); err != nil {
		return fmt.Errorf("failed to encode log data: %w", err)
	}
	return nil
}
