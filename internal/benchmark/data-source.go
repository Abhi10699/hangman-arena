package benchmark

import (
	"bufio"
	"fmt"
	"iter"
	"log"
	"os"
)

type Datasource interface {
	GetNextWord() (string, error)
}

type SingleWordDataSource struct {
	Word string
}

func NewSingleWordDataSource(word string) *SingleWordDataSource {
	return &SingleWordDataSource{
		Word: word,
	}
}

func (ds *SingleWordDataSource) GetNextWord() (string, error) {
	return ds.Word, nil
}

// word list data source

type WordIterator struct {
	Next func() (string, bool)
	Stop func()
}
type WordListDataSource struct {
	Path  string
	Limit int
	Words WordIterator
}

func NewWordListDataSource(path string, limit int) *WordListDataSource {
	filePtr, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	// file content reader

	var wordIter iter.Seq[string] = func(yield func(string) bool) {
		defer filePtr.Close()

		scanner := bufio.NewScanner(filePtr)
		for scanner.Scan() {
			if !yield(scanner.Text()) {
				break
			}
		}
	}

	next, stop := iter.Pull(wordIter)
	return &WordListDataSource{
		Path:  path,
		Limit: limit,
		Words: WordIterator{
			Next: next,
			Stop: stop,
		},
	}
}

func (ds *WordListDataSource) GetNextWord() (string, error) {

	nextWord, hasNext := ds.Words.Next()
	if !hasNext {
		return "", fmt.Errorf("over")
	}

	return nextWord, nil
}
