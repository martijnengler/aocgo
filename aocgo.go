// Package aocgo provides functions to get your Advent of Code puzzle inputs in a non-intrusive way.
package aocgo

import (
	"strings"

	"dalton.dog/aocgo/internal/dirparse"
	"dalton.dog/aocgo/internal/resources"
	"dalton.dog/aocgo/internal/session"
	"github.com/charmbracelet/log"
)

// GetInputAsByteArray will return the user's puzzle input, as determined by the file's working directory, as an array of bytes.
func GetInputAsByteArray() []byte {
	year, day, err := dirparse.GetYearAndDayFromCWD()
	if err != nil {
		log.Fatal(err)
	}

	return getData(year, day)
}

// GetInputAsString will return the user's puzzle input, as determined by the file's working directory, as a single string.
func GetInputAsString() string {
	return string(GetInputAsByteArray())
}

// GetInputAsLineArray will return the user's puzzle input, as determined by the file's working directory, as an array of strings, split on newline.
func GetInputAsLineArray() []string {
	return strings.Split(GetInputAsString(), "\n")
}

func getData(year int, day int) []byte {
	userToken, err := session.GetSessionToken()
	if err != nil {
		log.Fatal(err)
	}
	puzzle := resources.LoadOrCreatePuzzle(year, day, userToken)
	input, err := puzzle.GetUserInput()
	if err != nil {
		log.Fatal(err)
	}
	return input
}
