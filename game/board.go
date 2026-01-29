package game

// CellState represents the state of a cell on the board
type CellState int

const (
	Empty CellState = iota
	ShipCell
	Hit
	Miss
)

// BoardSize is the dimension of the game board
const BoardSize = 10

// Board represents a game board
type Board struct {
	Cells [BoardSize][BoardSize]CellState
	Ships []*Ship
}

// NewBoard creates a new empty board
func NewBoard() *Board {
	return &Board{
		Ships: make([]*Ship, 0),
	}
}

// CanPlaceShip checks if a ship can be placed at the given position
func (b *Board) CanPlaceShip(ship *Ship, row, col int, horizontal bool) bool {
	positions := getShipPositions(ship.Length, row, col, horizontal)
	if positions == nil {
		return false
	}

	for _, pos := range positions {
		if b.Cells[pos[0]][pos[1]] != Empty {
			return false
		}
	}
	return true
}

// PlaceShip places a ship on the board
func (b *Board) PlaceShip(ship *Ship, row, col int, horizontal bool) bool {
	if !b.CanPlaceShip(ship, row, col, horizontal) {
		return false
	}

	positions := getShipPositions(ship.Length, row, col, horizontal)
	ship.Positions = positions
	ship.Hits = make([]bool, len(positions))

	for _, pos := range positions {
		b.Cells[pos[0]][pos[1]] = ShipCell
	}

	b.Ships = append(b.Ships, ship)
	return true
}

// Attack attacks a cell and returns true if it was a hit, and the ship name if sunk
func (b *Board) Attack(row, col int) (hit bool, alreadyAttacked bool, sunkShipName string) {
	if row < 0 || row >= BoardSize || col < 0 || col >= BoardSize {
		return false, true, ""
	}

	cell := b.Cells[row][col]
	if cell == Hit || cell == Miss {
		return false, true, ""
	}

	if cell == ShipCell {
		b.Cells[row][col] = Hit
		// Mark the hit on the ship and check if sunk
		for _, ship := range b.Ships {
			for i, pos := range ship.Positions {
				if pos[0] == row && pos[1] == col {
					ship.Hits[i] = true
					if ship.IsSunk() {
						return true, false, ship.Name
					}
					break
				}
			}
		}
		return true, false, ""
	}

	b.Cells[row][col] = Miss
	return false, false, ""
}

// AllShipsSunk returns true if all ships have been sunk
func (b *Board) AllShipsSunk() bool {
	for _, ship := range b.Ships {
		if !ship.IsSunk() {
			return false
		}
	}
	return len(b.Ships) > 0
}

// getShipPositions calculates positions for a ship placement
func getShipPositions(length, row, col int, horizontal bool) [][2]int {
	positions := make([][2]int, length)

	for i := 0; i < length; i++ {
		var r, c int
		if horizontal {
			r, c = row, col+i
		} else {
			r, c = row+i, col
		}

		if r < 0 || r >= BoardSize || c < 0 || c >= BoardSize {
			return nil
		}
		positions[i] = [2]int{r, c}
	}

	return positions
}

// HasShipAt returns true if there's a ship at the given position
func (b *Board) HasShipAt(row, col int) bool {
	return b.Cells[row][col] == ShipCell || b.Cells[row][col] == Hit
}
