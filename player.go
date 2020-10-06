package main

import (
	"fmt"
	"math/rand"
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
			a[i][j] = []Direction{}
		}
	}
	return a
}

func mkIntGrid(boardSize int, fill int) [][]int {
	a := make([][]int, boardSize)
	for i := range a {
		a[i] = make([]int, boardSize)
		for j := range a[i] {
			a[i][j] = fill
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

func showIntGrid(g [][]int) string {
	res := ""
	for i := len(g) - 1; i >= 0; i-- {
		for j := range g[i] {
			if g[i][j] == -1 {
				res += "."
			} else {
				res += strconv.Itoa(g[i][j])
			}
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

func inBoard(boardSize int, x int, y int) bool {
	return x >= 0 && x < boardSize && y >= 0 && y < boardSize
}

func scorePut(g Grid, x int, y int, dir Direction) int {
	score := 0
	switch len(g[y][x]) {
	case 4:
		score += 0
	case 3:
		score += 0
	case 2:
		score -= 90
	case 1:
		score += 100
	}

	offX, offY, _ := getOffCoordAndDir(x, y, dir)
	if inBoard(len(g), offX, offY) {
		switch len(g[offY][offX]) {
		case 4:
			score += 0
		case 3:
			score += 0
		case 2:
			if len(g[y][x]) == 1 {
				score += 150
			} else {
				score -= 90
			}

		case 1:
			score += 100
		}
	}

	return score
}

func findActionInCorridor(g Grid, coloredGrid [][]int, minColor int) (int, int, Direction) {
	for i := range g {
		for j := range g[i] {
			for _, d := range g[i][j] {
				if coloredGrid[i][j] == minColor {
					return j, i, d
				}
			}
		}
	}
	panic("No action found")
}

func findAction(g Grid) (int, int, Direction, int) {
	type Key struct {
		x, y int
		dir  Direction
	}

	allActionsScored := map[Key]int{}

	for i := range g {
		for j := range g[i] {
			for _, d := range g[i][j] {
				if len(g[i][j]) > 0 {
					allActionsScored[Key{j, i, d}] = scorePut(g, j, i, d)
				}
			}
		}
	}

	bestScore := -1000

	for _, s := range allActionsScored {
		if s > bestScore {
			bestScore = s
		}
	}

	if bestScore == -1000 {
		panic("No action found")
	}

	allBests := make([]Key, 0)

	for key, s := range allActionsScored {
		if s == bestScore {
			allBests = append(allBests, key)
		}
	}

	best := allBests[rand.Intn(len(allBests))]

	return best.x, best.y, best.dir, bestScore
}

func hasReachedAMidState(g Grid) bool {
	for i := range g {
		for j := range g[i] {
			if len(g[i][j]) == 0 {
				continue
			} else if len(g[i][j]) == 1 {
				return false
			} else if len(g[i][j]) == 2 {
				continue
			} else {
				for _, d := range g[i][j] {
					offX, offY, _ := getOffCoordAndDir(j, i, d)
					if !inBoard(len(g), offX, offY) || len(g[offY][offX]) > 2 {
						return false
					}
				}
			}
		}
	}
	return true
}

func exploreCorridor(g Grid, coloredGrid [][]int, color int, x int, y int) {
	coloredGrid[y][x] = color
	sidesPlayable := len(g[y][x])
	switch sidesPlayable {
	case 0, 3, 4:
		return
	case 1, 2:
		for _, d := range g[y][x] {
			offX, offY, _ := getOffCoordAndDir(x, y, d)
			if inBoard(len(g), offX, offY) && coloredGrid[offY][offX] == -1 && len(g[offY][offX]) <= 2 {
				exploreCorridor(g, coloredGrid, color, offX, offY)
			}
		}
	}
}

func computeCorridors(g Grid) ([][]int, int) {
	res := mkIntGrid(len(g), -1)
	currentColor := -1
	for i := range g {
		for j := range g[i] {
			if res[i][j] == -1 && len(g[i][j]) > 0 {
				currentColor++
				exploreCorridor(g, res, currentColor, j, i)
			}
		}
	}
	return res, currentColor
}

func showDir(d Direction) int32 {
	switch d {
	case Up:
		return 'T'
	case Down:
		return 'B'
	case Left:
		return 'L'
	case Right:
		return 'R'
	}
	panic("Unhandled dir")
}

func bestColor(coloredGrid [][]int) (int, int) {
	counts := map[int]int{}

	for i := range coloredGrid {
		for j := range coloredGrid[i] {
			color := coloredGrid[j][i]
			if color != -1 {
				counts[color] += 1
			}
		}
	}

	minScore := 0
	minColor := -1
	for k, v := range counts {
		if minColor == -1 || v < minScore {
			minScore = v
			minColor = k
		}
	}
	return minColor, minScore
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
				default:
					panic(fmt.Sprintf("Unknown %c", char))
				}
			}

			x := box[0] - 'A'
			y := box[1] - '1'
			g[y][x] = parsedSides
		}

		fmt.Fprintln(os.Stderr, showGrid(g))

		if hasReachedAMidState(g) {
			coloredGrid, nbColors := computeCorridors(g)
			fmt.Fprintf(os.Stderr, "colored from %d to %d is\n%s", 0, nbColors, showIntGrid(coloredGrid))

			minColor, minScore := bestColor(coloredGrid)

			fmt.Fprintf(os.Stderr, "best color is %d with score %d", minColor, minScore)

			x, y, dir := findActionInCorridor(g, coloredGrid, minColor)
			fmt.Println(fmt.Sprintf("%c%c %c", x+'A', y+'1', showDir(dir)))
		} else {
			x, y, dir, score := findAction(g)
			fmt.Fprintf(os.Stderr, "best score is %d", score)

			// fmt.Fprintln(os.Stderr, "Debug messages...")
			fmt.Println(fmt.Sprintf("%c%c %c", x+'A', y+'1', showDir(dir)))
		}
	}
}
