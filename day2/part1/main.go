package main

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type Game struct {
	Id    int
	Draws []Draw
}
type Draw struct {
	Red   int
	Green int
	Blue  int
}

// Maximum number of individual balls in the bag
var maxDraw = Draw{
	Red:   12,
	Green: 13,
	Blue:  14,
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
		game := ParseGame(line)

		if IsValid(game) {
			sum += game.Id
		}
	}

	// result
	fmt.Printf("Result: %d\n", sum)
}

// Checks if the Game's set of Draws are all valid, compared to the max possible draw set (maxDraw)
func IsValid(game Game) bool {
	for _, draw := range game.Draws {
		if draw.Red > maxDraw.Red {
			return false
		}
		if draw.Green > maxDraw.Green {
			return false
		}
		if draw.Blue > maxDraw.Blue {
			return false
		}
	}
	return true
}

// ParseGame prases a string defining a Game
// Example line: `Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red`
func ParseGame(line string) Game {
	strGameStart, strDrawList, found := strings.Cut(line, ": ")
	if !found {
		log.Fatal("Incorrect format")
	}

	strId := strings.TrimLeft(strGameStart, "Game ")
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Fatal("Game number must be int", err)
	}

	draws := []Draw{}
	strDraws := strings.Split(strDrawList, "; ")

	for _, strDraw := range strDraws {
		strColors := strings.Split(strDraw, ", ")

		draw := Draw{}
		for _, strColor := range strColors {
			strColorAmount, strColorName, found := strings.Cut(strColor, " ")
			if !found {
				log.Fatal("Incorrect format")
			}

			amount, err := strconv.Atoi(strColorAmount)
			if err != nil {
				log.Fatal("Game number must be int", err)
			}

			if strColorName == "red" {
				draw.Red = amount
			}
			if strColorName == "green" {
				draw.Green = amount
			}
			if strColorName == "blue" {
				draw.Blue = amount
			}
		}

		draws = append(draws, draw)
	}

	return Game{Id: id, Draws: draws}
}
