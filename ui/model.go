package ui

import (
	"time"

	"battle-ship/game"

	tea "github.com/charmbracelet/bubbletea"
)

// GameState represents the current state of the game
type GameState int

const (
	StateMenu GameState = iota
	StatePlacement
	StateBattle
	StateGameOver
)

// Model represents the application state
type Model struct {
	State             GameState
	PlayerBoard       *game.Board
	AIBoard           *game.Board
	AI                *game.AI
	CursorRow         int
	CursorCol         int
	PlayerTurn        bool
	Message           string
	PlayerWon         bool
	ShipsToPlace      []*game.Ship
	CurrentShipIndex  int
	PlacingHorizontal bool
}

// NewModel creates a new game model
func NewModel() Model {
	return Model{
		State:             StateMenu,
		PlayerBoard:       game.NewBoard(),
		AIBoard:           game.NewBoard(),
		AI:                game.NewAI(),
		CursorRow:         0,
		CursorCol:         0,
		PlayerTurn:        true,
		ShipsToPlace:      game.ShipDefinitions(),
		CurrentShipIndex:  0,
		PlacingHorizontal: true,
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		switch m.State {
		case StateMenu:
			return m.updateMenu(msg)
		case StatePlacement:
			return m.updatePlacement(msg)
		case StateBattle:
			return m.updateBattle(msg)
		case StateGameOver:
			return m.updateGameOver(msg)
		}
	
	case aiTurnMsg:
		return m.handleAITurn()
	}

	return m, nil
}

// View implements tea.Model
func (m Model) View() string {
	switch m.State {
	case StateMenu:
		return renderMenu()
	case StatePlacement:
		return m.renderPlacement()
	case StateBattle:
		return m.renderBattle()
	case StateGameOver:
		return m.renderGameOver()
	default:
		return "Unknown state"
	}
}

// updateMenu handles input during the menu state
func (m Model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.State = StatePlacement
	}
	return m, nil
}

// updatePlacement handles input during ship placement
func (m Model) updatePlacement(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.CursorRow > 0 {
			m.CursorRow--
		}
	case "down", "j":
		if m.CursorRow < game.BoardSize-1 {
			m.CursorRow++
		}
	case "left", "h":
		if m.CursorCol > 0 {
			m.CursorCol--
		}
	case "right", "l":
		if m.CursorCol < game.BoardSize-1 {
			m.CursorCol++
		}
	case "r":
		m.PlacingHorizontal = !m.PlacingHorizontal
	case "enter":
		ship := m.ShipsToPlace[m.CurrentShipIndex]
		if m.PlayerBoard.PlaceShip(ship, m.CursorRow, m.CursorCol, m.PlacingHorizontal) {
			m.CurrentShipIndex++
			if m.CurrentShipIndex >= len(m.ShipsToPlace) {
				// All ships placed, start battle
				m.AI.PlaceShipsRandomly(m.AIBoard)
				m.State = StateBattle
				m.CursorRow = 0
				m.CursorCol = 0
				m.Message = "All ships placed! Fire at will!"
			}
		}
	}
	return m, nil
}

// updateBattle handles input during the battle phase
func (m Model) updateBattle(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if !m.PlayerTurn {
		return m, nil
	}

	switch msg.String() {
	case "up", "k":
		if m.CursorRow > 0 {
			m.CursorRow--
		}
	case "down", "j":
		if m.CursorRow < game.BoardSize-1 {
			m.CursorRow++
		}
	case "left", "h":
		if m.CursorCol > 0 {
			m.CursorCol--
		}
	case "right", "l":
		if m.CursorCol < game.BoardSize-1 {
			m.CursorCol++
		}
	case "enter":
		hit, alreadyAttacked, sunkShipName := m.AIBoard.Attack(m.CursorRow, m.CursorCol)
		if alreadyAttacked {
			m.Message = "Already attacked this location!"
			return m, nil
		}

		if hit {
			if sunkShipName != "" {
				m.Message = "HIT! You sunk their " + sunkShipName + "!"
			} else {
				m.Message = "HIT!"
			}
			if m.AIBoard.AllShipsSunk() {
				m.State = StateGameOver
				m.PlayerWon = true
				return m, nil
			}
		} else {
			m.Message = "Miss..."
		}

		// AI's turn
		m.PlayerTurn = false
		return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
			return aiTurnMsg{}
		})
	}
	return m, nil
}

// aiTurnMsg is sent when it's time for the AI to make a move
type aiTurnMsg struct{}

// handleAITurn processes the AI's attack
func (m Model) handleAITurn() (tea.Model, tea.Cmd) {
	row, col := m.AI.ChooseAttack()
	hit, _, sunkShipName := m.PlayerBoard.Attack(row, col)

	if hit {
		m.AI.RecordHit(row, col)
		if sunkShipName != "" {
			m.Message += " | Enemy sunk your " + sunkShipName + "!"
		} else {
			m.Message += " | Enemy hit your ship!"
		}
		if m.PlayerBoard.AllShipsSunk() {
			m.State = StateGameOver
			m.PlayerWon = false
			return m, nil
		}
	} else {
		m.AI.RecordMiss(row, col)
		m.Message += " | Enemy missed."
	}

	m.PlayerTurn = true
	return m, nil
}

// updateGameOver handles input during game over
func (m Model) updateGameOver(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Reset the game
		return NewModel(), nil
	}
	return m, nil
}
