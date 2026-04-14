package main

import (
	"flag"
	"fmt"
	"hangman/internal/benchmark"
	"hangman/internal/player"
	"os"
)

func main() {
	fmt.Println("Hangman CLI Benchmark")

	totalGames := flag.Int("total-games", 10, "Number of games to run in the benchmark")
	datasource := flag.String("source", "single-word", "The word source to use: 'single-word' or 'word-list'")
	datasourceFile := flag.String("word-list", "", "Path to the line-separated word list file (required if -source=word-list)")
	singleWord := flag.String("word", "hello", "The word to use for 'single-word' source")

	flag.Parse()

	var wordSource benchmark.Datasource

	switch *datasource {
	case "single-word":
		if *singleWord == "" {
			fmt.Println("Error: -word must be provided for 'single-word' source")
			flag.Usage()
			os.Exit(1)
		}
		wordSource = benchmark.NewSingleWordDataSource(*singleWord)
	case "word-list":
		if *datasourceFile == "" {
			fmt.Println("Error: -word-list path must be provided for 'word-list' source")
			flag.Usage()
			os.Exit(1)
		}
		if _, err := os.Stat(*datasourceFile); os.IsNotExist(err) {
			fmt.Printf("Error: word list file '%s' does not exist\n", *datasourceFile)
			os.Exit(1)
		}
		wordSource = benchmark.NewWordListDataSource(*datasourceFile, *totalGames)
	default:
		fmt.Printf("Error: unknown source '%s'. Use 'single-word' or 'word-list'\n", *datasource)
		flag.Usage()
		os.Exit(1)
	}

	ollamaPlayer := player.NewOllamaGuesser("gemma4:latest", player.Local)

	benchmarkRunner := benchmark.NewBenchmarkRunner(
		wordSource,
		*totalGames,
		ollamaPlayer,
	)

	benchmarkRunner.RunBenchmark()
	fmt.Println(benchmarkRunner.Stats)
}
