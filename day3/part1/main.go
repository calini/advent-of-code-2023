package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Handler struct {
	Level slog.Level
}

// Add a simple DEBUG ENV var toggle for debug logs
func init() {
	debug, ok := os.LookupEnv("DEBUG")
	if ok && debug != "0" && debug != "false" && debug != "FALSE" {
		// TODO change just the Default logger Level once Go 1.22 drops
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})))
	}
}

func main() {
	if len(os.Args) != 2 {
		slog.Error("Usage: main input_file")
		os.Exit(1)
	}

	// file opening
	file, err := os.Open(os.Args[1])
	defer file.Close()

	if err != nil {
		slog.Error("Error opening file: %s", err)
		os.Exit(2)
	}

	// file reading
	sc := bufio.NewScanner(file)

	// business logic
	sum := 0

	// We need a 3-row sliding window scanning down,
	// checking the middle row for numbers, then checking if the numbers are a 'part'
	// by checking for any non '.' characters

	// pad the table a fake upper line
	var line string
	if sc.Scan() {
		line = "." + sc.Text() + "." // pad vertically
	}
	upper := strings.Repeat(".", len(line))
	// scan through
	for sc.Scan() {
		lower := "." + sc.Text() + "." // pad veritcally
		numbers := PartNumbers(upper, line, lower)

		// shift lines
		upper = line
		line = lower

		// add results to sum
		for _, nr := range numbers {
			sum += nr
		}
	}

	// pad with an extra lower line
	lower := strings.Repeat(".", len(line))
	numbers := PartNumbers(upper, line, lower)
	for _, nr := range numbers {
		sum += nr
	}

	// result
	fmt.Printf("Result: %d\n", sum)
}

// Finds "part numbers" in the middle row.
// "Part number" = A number that has a neighbouring special character (non-number, not ".")
func PartNumbers(upper, middle, lower string) (numbers []int) {
	// find numbers in the midle row
	var middleRunes = []rune(middle)
	for i := 0; i < len(middleRunes); i++ {
		if unicode.IsDigit(middleRunes[i]) {
			idxLeft := i
			for unicode.IsDigit(middleRunes[i]) {
				i++
			}
			idxRight := i

			number, err := strconv.Atoi(string(middleRunes[idxLeft:idxRight]))
			if err != nil {
				slog.Error("Cannot convert number")
			}

			if IsPartNumber([]rune(upper), []rune(middle), []rune(lower), idxLeft, idxRight) {
				numbers = append(numbers, number)
				slog.Debug("Found", "number", number)
			}
		}

	}
	return numbers
}

// IsPartNumber checks around a number for special characters
func IsPartNumber(upper, middle, lower []rune, leftIndex, rightIndex int) bool {

	if IsSpecial(middle[leftIndex-1]) || IsSpecial(middle[rightIndex]) {
		return true
	}

	for i := 0; i < rightIndex-leftIndex+2; i++ {
		if IsSpecial(lower[leftIndex+i-1]) || IsSpecial(upper[leftIndex+i-1]) {
			return true
		}
	}

	return false
}

// IsSpecial returns true if passed a non-'.' and non-digit character
func IsSpecial(ch rune) bool {
	if ch == '.' || unicode.IsDigit(ch) {
		return false
	}
	return true
}
