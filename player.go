package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
)

type direction uint8
type grid [][][]direction

type state struct {
	grid          grid
	playerScore   uint8
	opponentScore uint8
	myTurn        bool
}

const (
	up direction = iota
	down
	left
	right
)

func mkGrid(boardSize int) grid {
	a := make([][][]direction, boardSize)
	for i := range a {
		a[i] = make([][]direction, boardSize)
		for j := range a[i] {
			a[i][j] = []direction{}
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

func showGrid(g grid) string {
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

func dirInSlice(a direction, list []direction) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getOffCoordAndDir(x int, y int, dir direction) (int, int, direction) {
	switch dir {
	case up:
		return x, y + 1, down
	case down:
		return x, y - 1, up
	case left:
		return x - 1, y, right
	case right:
		return x + 1, y, left
	}
	panic("Unhandled dir")
}

func playSide(x int, y int, s *state, dir direction) {
	offx, offy, offdir := getOffCoordAndDir(x, y, dir)
	putSide(x, y, s, dir)

	if inBoard(len(s.grid), offx, offy) {
		putSide(offx, offy, s, offdir)
	}
}

func indexOf(list []direction, d direction) int {
	for k, v := range list {
		if d == v {
			return k
		}
	}
	return -1 //not found.
}

func removeDirSlice(list []direction, d direction) []direction {
	index := indexOf(list, d)
	return append(list[:index], list[index+1:]...)
}

func putSide(x int, y int, s *state, dir direction) {
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

func scorePut(g grid, x int, y int, dir direction) int {
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

func copyGrid(g grid) grid {
	duplicate := make([][][]direction, len(g))
	for i := range g {
		duplicate[i] = make([][]direction, len(g[i]))
		for j := range g[i] {
			duplicate[i][j] = make([]direction, len(g[i][j]))
			copy(duplicate[i][j], g[i][j])
		}
	}
	return duplicate
}

func findActionInCorridor(s state) (int, int, direction, int) {

	bestX := -1
	bestY := -1
	bestDir := up
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

func findAction(s state) (int, int, direction, int) {
	type Key struct {
		x, y int
		dir  direction
	}

	if hasReachedAMidState(s.grid) {
		return findActionInCorridor(s)
	}

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

func hasReachedAMidState(g grid) bool {
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

func exploreCorridor(s state, coloredGrid [][]int, color int) {
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

func computeCorridors(s state) ([][]int, int) {
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

func showDir(d direction) int32 {
	switch d {
	case up:
		return 'T'
	case down:
		return 'B'
	case left:
		return 'L'
	case right:
		return 'R'
	}
	panic("Unhandled dir")
}

func heuristic(state state) int {
	return int(state.playerScore) - int(state.opponentScore)
}

func copyState(s state) state {
	return state{
		grid:          copyGrid(s.grid),
		playerScore:   s.playerScore,
		opponentScore: s.opponentScore,
	}
}

func getChildStates(s state) []state {
	result := make([]state, 0, 7*7*4)

	for i := range s.grid {
		for j := range s.grid[i] {
			for _, d := range s.grid[i][j] {
				newState := copyState(s)
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

func isTerminal(state state) bool {
	for i := range state.grid {
		for j := range state.grid[i] {
			if len(state.grid[i][j]) > 0 {
				return false
			}
		}
	}
	return true
}

func minimax(state state, depth uint8, maximizingPlayer bool) int {
	if depth == 0 || isTerminal(state) {
		return heuristic(state)
	}

	if maximizingPlayer {
		value := math.MinInt32
		for _, child := range getChildStates(state) {
			value = Max(value, minimax(child, depth-1, false))
		}
		return value
	}

	value := math.MaxInt32
	for _, child := range getChildStates(state) {
		value = Min(value, minimax(child, depth-1, true))
	}
	return value
}

func main() {
	// boardSize: The size of the board.
	var boardSize int
	fmt.Scan(&boardSize)

	// playerID: The ID of the player. 'A'=first player, 'B'=second player.
	var playerID string
	fmt.Scan(&playerID)

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

			parsedSides := make([]direction, 0)
			for _, char := range sides {
				switch char {
				case 'L':
					parsedSides = append(parsedSides, left)
				case 'R':
					parsedSides = append(parsedSides, right)
				case 'T':
					parsedSides = append(parsedSides, up)
				case 'B':
					parsedSides = append(parsedSides, down)
				default:
					panic(fmt.Sprintf("Unknown %c", char))
				}
			}

			x := box[0] - 'A'
			y := box[1] - '1'
			g[y][x] = parsedSides
		}

		fmt.Fprintln(os.Stderr, showGrid(g))

		s := state{
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
		fmt.Printf("%c%c %c\n", x+'A', y+'1', showDir(dir))
	}
}
