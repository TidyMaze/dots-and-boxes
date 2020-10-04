package main

import (
	"fmt"
	"os"
	"strconv"
)

type Direction int
type Grid [][][]Direction

const (
	Up Direction = iota
	Down
	Left
	Right
)

func mkGrid(boardSize int) Grid {
	a := make([][][]Direction, boardSize)
	for i := range a {
		a[i] = make([][]Direction, boardSize)
		for j := range a[i] {
			a[i][j] = []Direction{Left, Up, Right, Down}
		}
	}
	return a
}

func showGrid(g Grid) string {
	res := ""
	for i := len(g) - 1; i >= 0; i-- {
		for j := range g[i] {
			res += strconv.Itoa(len(g[i][j]))
		}
		res += "\n"
	}
	return res
}

func dirInSlice(a Direction, list []Direction) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getOffCoordAndDir(x int, y int, dir Direction) (int, int, Direction) {
	switch dir {
	case Up:
		return x, y + 1, Down
	case Down:
		return x, y - 1, Up
	case Left:
		return x - 1, y, Right
	case Right:
		return x + 1, y, Left
	}
	panic("Unhandled dir")
}

func playSide(x int, y int, g Grid, dir Direction) {
	offx, offy, offdir := getOffCoordAndDir(x, y, dir)
	putSide(x, y, g, dir)
	putSide(offx, offy, g, offdir)
}

func indexOf(list []Direction, d Direction) int {
	for k, v := range list {
		if d == v {
			return k
		}
	}
	return -1 //not found.
}

func removeDirSlice(list []Direction, d Direction) []Direction {
	index := indexOf(list, d)
	return append(list[:index], list[index+1:]...)
}

func putSide(x int, y int, g Grid, dir Direction) {
	if !dirInSlice(dir, g[y][x]) {
		panic(fmt.Sprintf("cannot put side %d at %d %d already contained", dir, x, y))
	}
	g[y][x] = removeDirSlice(g[y][x], dir)
}

func findAction(g Grid) (int, int, Direction) {
	for i := range g {
		for j := range g[i] {
			for _, d := range []Direction{Up, Down, Left, Right} {
				if !dirInSlice(d, g[i][j]) {
					return j, i, d
				}
			}
		}
	}
	panic("No action found")
}

func main() {
	// boardSize: The size of the board.
	var boardSize int
	fmt.Scan(&boardSize)

	// playerId: The ID of the player. 'A'=first player, 'B'=second player.
	var playerId string
	fmt.Scan(&playerId)

	for {
		// playerScore: The player's score.
		// opponentScore: The opponent's score.
		var playerScore, opponentScore int
		fmt.Scan(&playerScore, &opponentScore)

		// numBoxes: The number of playable boxes.
		var numBoxes int
		fmt.Scan(&numBoxes)

		g := mkGrid(boardSize)

		for i := 0; i < numBoxes; i++ {
			// box: The ID of the playable box.
			// sides: Playable sides of the box.
			var box, sides string
			fmt.Scan(&box, &sides)

			parsedSides := make([]Direction, 0)
			for _, char := range sides {
				switch char {
				case 'L':
					parsedSides = append(parsedSides, Left)
				case 'R':
					parsedSides = append(parsedSides, Right)
				case 'T':
					parsedSides = append(parsedSides, Up)
				case 'B':
					parsedSides = append(parsedSides, Down)
				}
			}

			x := box[0] - 'A'
			y := box[1] - '1'
			g[y][x] = parsedSides
		}

		fmt.Fprintln(os.Stderr, showGrid(g))

		// fmt.Fprintln(os.Stderr, "Debug messages...")
		fmt.Println("A1 T MSG bla bla bla...") // <box> <side> [MSG Optional message]
	}
}
