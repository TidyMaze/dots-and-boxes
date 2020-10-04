package main

import (
	"fmt"
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
			a[i][j] = []Direction{}
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

func putSide(x int, y int, g Grid, dir Direction) {
	if dirInSlice(dir, g[y][x]) {
		panic(fmt.Sprintf("cannot put side %d at %d %d already contained", dir, x, y))
	}
	g[y][x] = append(g[y][x], dir)
}

func main() {
	g := mkGrid(2)
	println(showGrid(g))

	playSide(1, 1, g, Down)
	playSide(0, 1, g, Down)

	println(showGrid(g))

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

		for i := 0; i < numBoxes; i++ {
			// box: The ID of the playable box.
			// sides: Playable sides of the box.
			var box, sides string
			fmt.Scan(&box, &sides)
		}

		// fmt.Fprintln(os.Stderr, "Debug messages...")
		fmt.Println("A1 B MSG bla bla bla...") // <box> <side> [MSG Optional message]
	}
}
