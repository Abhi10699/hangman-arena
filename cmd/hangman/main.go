package main

import (
	"flag"
	"fmt"
	"hangman/internal/benchmark"
	"hangman/internal/player"
)

func main() {
	fmt.Println("Hangmang CLI")

	wordToGuess := flag.String("word", "cow", "Word to guess")
	totalGames := flag.Int("total-games", 10, "Number of games")

	flag.Parse()

	ollamaPlayer := player.NewOllamaGuesser("gemma4:latest", player.Local)

	benchmarkRunner := benchmark.NewBenchmarkRunner(
		*wordToGuess,
		*totalGames,
		ollamaPlayer,
	)
	benchmarkRunner.RunBenchmark()
	fmt.Println(benchmarkRunner.Stats)

}
