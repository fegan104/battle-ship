package game

import (
	"math/rand"
)

// AI handles computer opponent logic
type AI struct {
	lastHit      *[2]int
	huntMode     bool
	huntTargets  [][2]int
	attackedCells map[[2]int]bool
}

// NewAI creates a new AI opponent
func NewAI() *AI {
	return &AI{
		attackedCells: make(map[[2]int]bool),
	}
}

// PlaceShipsRandomly places all ships randomly on the board
func (ai *AI) PlaceShipsRandomly(board *Board) {
	ships := ShipDefinitions()

	for _, ship := range ships {
		placed := false
		for !placed {
			row := rand.Intn(BoardSize)
			col := rand.Intn(BoardSize)
			horizontal := rand.Intn(2) == 0

			if board.PlaceShip(ship, row, col, horizontal) {
				placed = true
			}
		}
	}
}

// ChooseAttack selects a cell to attack
func (ai *AI) ChooseAttack() (int, int) {
	// If we have hunt targets from a previous hit, try those first
	if ai.huntMode && len(ai.huntTargets) > 0 {
		for len(ai.huntTargets) > 0 {
			target := ai.huntTargets[0]
			ai.huntTargets = ai.huntTargets[1:]

			if !ai.attackedCells[target] {
				ai.attackedCells[target] = true
				return target[0], target[1]
			}
		}
		ai.huntMode = false
	}

	// Random attack
	for {
		row := rand.Intn(BoardSize)
		col := rand.Intn(BoardSize)
		pos := [2]int{row, col}

		if !ai.attackedCells[pos] {
			ai.attackedCells[pos] = true
			return row, col
		}
	}
}

// RecordHit tells the AI about a successful hit
func (ai *AI) RecordHit(row, col int) {
	ai.lastHit = &[2]int{row, col}
	ai.huntMode = true

	// Add adjacent cells to hunt targets
	adjacent := [][2]int{
		{row - 1, col},
		{row + 1, col},
		{row, col - 1},
		{row, col + 1},
	}

	for _, pos := range adjacent {
		if pos[0] >= 0 && pos[0] < BoardSize && pos[1] >= 0 && pos[1] < BoardSize {
			if !ai.attackedCells[pos] {
				ai.huntTargets = append(ai.huntTargets, pos)
			}
		}
	}
}

// RecordMiss tells the AI about a miss
func (ai *AI) RecordMiss(row, col int) {
	// Nothing special to do on a miss
}
