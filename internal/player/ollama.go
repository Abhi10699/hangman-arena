package player

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	ollama "github.com/ollama/ollama/api"
)

type OllamaHost string

const (
	Cloud OllamaHost = "cloud"
	Local OllamaHost = "local"
)

type OllamaGuesser struct {
	ModelId      string
	SystemPrompt string
	Messages     []ollama.Message
	Client       *ollama.Client
}

func NewOllamaGuesser(modelId string, host OllamaHost) *OllamaGuesser {
	resolvedHost, err := url.Parse("http://localhost:11434")
	if err != nil {
		log.Fatal(err)
	}
	// if host == Cloud {
	// 	resolvedHost, err = url.Parse("https://ollama.com/api")
	// }

	systemPrompt := `
	You are a Hangman solver bot. Your goal is to win the game using optimal linguistic strategy.
	RULES:
	1. You are provided with the current state of the word (e.g., "h _ _ p _") and a list of incorrect guesses.
	2. You must analyze the word length and patterns to choose the most likely letter.
	3. You MUST respond with a valid JSON object only. No conversational text.

	RESPONSE FORMAT:
	{
		"guessedCharacter": string,
		"reason": string
	}
	`
	client := ollama.NewClient(resolvedHost, &http.Client{})

	// check connection

	heartbeatErr := client.Heartbeat(context.Background())
	if heartbeatErr != nil {
		log.Fatal(heartbeatErr)
	}

	return &OllamaGuesser{
		ModelId:      modelId,
		Client:       client,
		SystemPrompt: systemPrompt,
		Messages: []ollama.Message{
			{
				Role:    "system",
				Content: systemPrompt,
			},
		},
	}
}

func (o *OllamaGuesser) GuessCharacter(incorrectGuesses []string, currentGuess string) GuessResult {
	chatReq := o.buildChatRequest(incorrectGuesses, currentGuess)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	modelResponse := o.callModel(ctx, chatReq)

	o.Messages = append(o.Messages, ollama.Message{
		Role:    "assistant",
		Content: modelResponse.RawMessage,
	})
	return GuessResult{
		GuessedCharacter: modelResponse.GuessedCharacter,
		Metadata: map[string]any{
			"reason": modelResponse.Reason,
		},
	}
}

func (o *OllamaGuesser) ResetState() {
	o.Messages = []ollama.Message{
		{
			Role:    "system",
			Content: o.SystemPrompt,
		},
	}
}

func (o *OllamaGuesser) buildChatRequest(incorrectGuesses []string, currentGuess string) *ollama.ChatRequest {
	stream := false
	responseSpecifications := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"guessedCharacter": "string",
			"reason":           "string",
		},
		"required": []string{"guessedCharacter", "reason"},
	}

	jsonBody, jsonErr := json.Marshal(responseSpecifications)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	userPrompt := fmt.Sprintf(
		"Current Word: %s\nTotal Lives: %v\nIncorrect Guesses:%s\nReturn only JSON.",
		strings.Join(strings.Split(currentGuess, ""), "  "),
		9-len(incorrectGuesses),
		strings.Join(incorrectGuesses, ","),
	)

	fmt.Println(userPrompt)

	messages := append(o.Messages, ollama.Message{
		Role:    "user",
		Content: userPrompt,
	})

	return &ollama.ChatRequest{
		Model:    o.ModelId,
		Messages: messages,
		Stream:   &stream,
		Think:    &ollama.ThinkValue{Value: false},
		Format:   jsonBody,
	}
}

type ModelResponse struct {
	RawMessage       string
	GuessedCharacter string
	Reason           string
	Thoughts         string
}

func (o *OllamaGuesser) callModel(ctx context.Context, chatReq *ollama.ChatRequest) ModelResponse {
	var jsonResponse struct {
		GuessedCharacter string `json:"guessedCharacter"`
		Reason           string `json:"reason"`
	}

	var rawMessage string

	respFunc := func(resp ollama.ChatResponse) error {
		if resp.Done {
			// parse json
			rawMessage = resp.Message.Content
			fmt.Println(rawMessage)
			parseErr := extractAndParse(rawMessage, &jsonResponse)
			return parseErr
		}
		return nil
	}

	chatErr := o.Client.Chat(ctx, chatReq, respFunc)
	if chatErr != nil {
		log.Fatal(chatErr)
	}

	return ModelResponse{
		RawMessage:       rawMessage,
		GuessedCharacter: jsonResponse.GuessedCharacter,
		Reason:           jsonResponse.Reason,
		Thoughts:         "", // future use
	}
}

func extractAndParse(raw string, target interface{}) error {
	raw = strings.ReplaceAll(raw, `\_`, `_`)
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start == -1 || end == -1 {
		return fmt.Errorf("no JSON found")
	}
	return json.Unmarshal([]byte(raw[start:end+1]), target)
}
