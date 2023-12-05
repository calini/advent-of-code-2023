package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"unicode"

	"github.com/calini/std/strutils"
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
	for sc.Scan() {
		line := sc.Text()
		left, right := BoundingDigits(line)
		if left == -1 || right == -1 {
			slog.Error("All lines should have at least one digit")
			os.Exit(3)
		}
		sum += 10*left + right
	}

	// result
	fmt.Printf(strconv.Itoa(sum))
}

// BoundingDigits gets the bounding digits in the string
func BoundingDigits(line string) (int, int) {
	left := FirstDigit(line)
	right := FirstDigit(strutils.Reverse(line))

	return left, right
}

// FirstDigit finds the first instance of a digit represented in the string, or -1 if none found
func FirstDigit(line string) int {
	for _, ch := range line {
		if unicode.IsDigit(ch) {
			digit, err := strconv.Atoi(string(ch))
			if err != nil {
				slog.Error("Conversion error")
			}
			slog.Debug("found", "string", line, "digit", digit)
			return digit
		}
	}
	return -1
}
