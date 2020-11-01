package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
)

type Direction uint8
type Grid [][][]Direction

type State struct {
	grid          Grid
	playerScore   uint8
	opponentScore uint8
	myTurn        bool
}

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

func playSide(x int, y int, s *State, dir Direction) {
	offx, offy, offdir := getOffCoordAndDir(x, y, dir)
	putSide(x, y, s, dir)

	if inBoard(len(s.grid), offx, offy) {
		putSide(offx, offy, s, offdir)
	}
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

func putSide(x int, y int, s *State, dir Direction) {
	if !dirInSlice(dir, s.grid[y][x]) {
		panic(fmt.Sprintf("cannot put side %d at %d %d already contained", dir, x, y))
	}
	s.grid[y][x] = removeDirSlice(s.grid[y][x], dir)

	if len(s.grid[y][x]) == 0 {
		if s.myTurn {
			s.playerScore++
		} else {
			s.opponentScore++
		}
	}
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

func copyGrid(g Grid) Grid {
	duplicate := make([][][]Direction, len(g))
	for i := range g {
		duplicate[i] = make([][]Direction, len(g[i]))
		for j := range g[i] {
			duplicate[i][j] = make([]Direction, len(g[i][j]))
			copy(duplicate[i][j], g[i][j])
		}
	}
	return duplicate
}

func findActionInCorridor(s State) (int, int, Direction, int) {

	bestX := -1
	bestY := -1
	bestDir := Up
	bestScore := -1

	g := s.grid

	for i := range g {
		for j := range g[i] {
			for _, d := range g[i][j] {
				modifiedState := copyState(s)
				modifiedState.myTurn = false
				playSide(j, i, &modifiedState, d)

				coloredGrid, nbReachable := computeCorridors(modifiedState)

				if bestScore == -1 || nbReachable < bestScore {
					fmt.Fprintf(os.Stderr, "best is %d with\n%s", nbReachable, showIntGrid(coloredGrid))

					bestScore = nbReachable
					bestX = j
					bestY = i
					bestDir = d
				}
			}
		}
	}

	if bestScore == -1 {
		panic("No action found :/")
	}
	return bestX, bestY, bestDir, bestScore
}

func findAction(s State) (int, int, Direction, int) {
	type Key struct {
		x, y int
		dir  Direction
	}

	if hasReachedAMidState(s.grid) {
		return findActionInCorridor(s)
	} else {
		allActionsScored := map[Key]int{}

		for i := range s.grid {
			for j := range s.grid[i] {
				for _, d := range s.grid[i][j] {
					if len(s.grid[i][j]) > 0 {
						allActionsScored[Key{j, i, d}] = scorePut(s.grid, j, i, d)
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

func exploreCorridor(s State, coloredGrid [][]int, color int) {
	foundOne := true
	for foundOne {
		foundOne = false
		for i := range s.grid {
			for j := range s.grid[i] {
				if len(s.grid[i][j]) == 1 {
					foundOne = true
					coloredGrid[i][j] = color
					playSide(j, i, &s, s.grid[i][j][0])
				}
			}
		}
	}
}

func computeCorridors(s State) ([][]int, int) {
	res := mkIntGrid(len(s.grid), -1)

	exploreCorridor(s, res, 0)

	count := 0
	for i := range s.grid {
		for j := range s.grid[i] {
			if res[i][j] == 0 {
				count++
			}
		}
	}
	return res, count
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

func heuristic(state State) int {
	return int(state.playerScore) - int(state.opponentScore)
}

func copyState(state State) State {
	return State{
		grid:          copyGrid(state.grid),
		playerScore:   state.playerScore,
		opponentScore: state.opponentScore,
	}
}

func getChildStates(state State) []State {
	result := make([]State, 0, 7*7*4)

	for i := range state.grid {
		for j := range state.grid[i] {
			for _, d := range state.grid[i][j] {
				newState := copyState(state)
				playSide(j, i, &newState, d)
				result = append(result, newState)
			}
		}
	}

	return result
}

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func isTerminal(state State) bool {
	for i := range state.grid {
		for j := range state.grid[i] {
			if len(state.grid[i][j]) > 0 {
				return false
			}
		}
	}
	return true
}

func minimax(state State, depth uint8, maximizingPlayer bool) int {
	if depth == 0 || isTerminal(state) {
		return heuristic(state)
	}

	if maximizingPlayer {
		value := math.MinInt32
		for _, child := range getChildStates(state) {
			value = Max(value, minimax(child, depth-1, false))
		}
		return value
	} else {
		value := math.MaxInt32
		for _, child := range getChildStates(state) {
			value = Min(value, minimax(child, depth-1, true))
		}
		return value
	}
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

		s := State{
			grid:          g,
			playerScore:   uint8(playerScore),
			opponentScore: uint8(opponentScore),
			myTurn:        true,
		}

		x, y, dir, score := findAction(s)

		depth := uint8(2)

		resultScore := minimax(s, depth, true)

		fmt.Fprintf(os.Stderr, "best score is %d, minimax score is %d with depth %d", score, resultScore, depth)

		// fmt.Fprintln(os.Stderr, "Debug messages...")
		fmt.Println(fmt.Sprintf("%c%c %c", x+'A', y+'1', showDir(dir)))
	}
}
