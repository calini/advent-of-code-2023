package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
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

	// preload the first line with a fake upper line
	var line string
	if sc.Scan() {
		line = sc.Text()
		slog.Debug("read", "line", line)
	}
	upper := strings.Repeat(".", len(line))
	// scan through
	for sc.Scan() {
		lower := sc.Text()
		numbers := PartNumbers(upper, line, lower)

		// shift lines
		upper = line
		line = lower

		// add results to sum
		for _, nr := range numbers {
			sum += nr
		}
	}

	// post-load an extra lower line
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
func PartNumbers(upper, middle, lower string) []int {
	slog.Debug("Checking lines", "upper", upper, "middle", middle, "lower", lower)
	return []int{}
}
