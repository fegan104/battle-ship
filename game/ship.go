package game

// Ship represents a ship in the game
type Ship struct {
	Name      string
	Length    int
	Positions [][2]int // [row, col] pairs
	Hits      []bool   // which positions have been hit
}

// NewShip creates a new ship
func NewShip(name string, length int) *Ship {
	return &Ship{
		Name:   name,
		Length: length,
	}
}

// IsSunk returns true if all positions have been hit
func (s *Ship) IsSunk() bool {
	if len(s.Hits) == 0 {
		return false
	}
	for _, hit := range s.Hits {
		if !hit {
			return false
		}
	}
	return true
}

// ShipDefinitions returns the standard set of ships for Battleship
func ShipDefinitions() []*Ship {
	return []*Ship{
		NewShip("Carrier", 5),
		NewShip("Battleship", 4),
		NewShip("Cruiser", 3),
		NewShip("Submarine", 3),
		NewShip("Destroyer", 2),
	}
}
