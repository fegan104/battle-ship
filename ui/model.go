package ui

import (
	"fmt"
	"time"

	"battle-ship/game"
	bnet "battle-ship/net"

	tea "github.com/charmbracelet/bubbletea"
)

// GameState represents the current state of the game
type GameState int

const (
	StateMenu GameState = iota
	StatePlacement
	StateBattle
	StateGameOver
	// Multiplayer states
	StateMPMenu        // New: Choose Host or Join
	StateMPHostWaiting // Waiting for room creation/opponent
	StateMPJoinInput   // Entering room code
	StateMPConnecting  // Connecting to server
	StateMPPlacement
	StateMPWaitingForOpponent
	StateMPBattle
)

// GameMode represents the type of game being played
type GameMode int

const (
	ModeVsAI GameMode = iota
	ModeMultiplayer
)

// Model represents the application state
type Model struct {
	State             GameState
	GameMode          GameMode
	PlayerBoard       *game.Board
	AIBoard           *game.Board
	OpponentBoard     *game.Board // Used in multiplayer
	AI                *game.AI
	CursorRow         int
	CursorCol         int
	PlayerTurn        bool
	Message           string
	PlayerWon         bool
	ShipsToPlace      []*game.Ship
	CurrentShipIndex  int
	PlacingHorizontal bool

	// Menu selection
	MenuSelection int

	// Multiplayer
	Connection    *bnet.Connection
	ServerAddress string
	RoomCode      string
	IsHost        bool
	ShipsPlaced   bool
	OpponentReady bool
	LastAttackRow int
	LastAttackCol int
}

// NewModel creates a new game model
func NewModel() Model {
	return Model{
		State:             StateMenu,
		GameMode:          ModeVsAI,
		PlayerBoard:       game.NewBoard(),
		AIBoard:           game.NewBoard(),
		OpponentBoard:     game.NewBoard(),
		AI:                game.NewAI(),
		CursorRow:         0,
		CursorCol:         0,
		PlayerTurn:        true,
		ShipsToPlace:      game.ShipDefinitions(),
		CurrentShipIndex:  0,
		PlacingHorizontal: true,
		MenuSelection:     0,
		ServerAddress:     "battleship-server-350181966586.us-central1.run.app", // Default central server or localhost:8080 for local development
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
			m.cleanup()
			return m, tea.Quit
		}

		switch m.State {
		case StateMenu:
			return m.updateMenu(msg)
		case StateMPMenu:
			return m.updateMPMenu(msg)
		case StatePlacement:
			return m.updatePlacement(msg)
		case StateBattle:
			return m.updateBattle(msg)
		case StateGameOver:
			return m.updateGameOver(msg)
		case StateMPHostWaiting:
			return m.updateMPHostWaiting(msg)
		case StateMPJoinInput:
			return m.updateMPJoinInput(msg)
		case StateMPPlacement:
			return m.updateMPPlacement(msg)
		case StateMPWaitingForOpponent:
			return m.updateMPWaiting(msg)
		case StateMPBattle:
			return m.updateMPBattle(msg)
		}

	case aiTurnMsg:
		return m.handleAITurn()

	case connectionEstablishedMsg:
		return m.handleConnectionEstablished(msg)

	case connectionErrorMsg:
		m.Message = "Connection error: " + msg.err.Error()
		m.State = StateMenu
		return m, nil

	case roomCreatedMsg:
		m.RoomCode = msg.code
		m.State = StateMPHostWaiting
		m.Message = fmt.Sprintf("Room Created! Code: %s. Waiting for opponent...", m.RoomCode)
		return m, m.messageLoop()

	case playerJoinedMsg:
		m.Message = "Player joined! Game starting..."
		newModel, cmd := m.startGame()
		m = newModel.(Model)
		return m, tea.Batch(cmd, m.messageLoop())

	case joinErrorMsg:
		m.Message = "Error: " + msg.err
		m.State = StateMPMenu
		return m, nil

	case opponentReadyMsg:
		newModel, cmd := m.handleOpponentReady()
		m = newModel.(Model)
		return m, tea.Batch(cmd, m.messageLoop())

	case opponentAttackMsg:
		newModel, cmd := m.handleOpponentAttack(msg)
		m = newModel.(Model)
		return m, tea.Batch(cmd, m.messageLoop())

	case attackResultMsg:
		newModel, cmd := m.handleAttackResult(msg)
		m = newModel.(Model)
		return m, tea.Batch(cmd, m.messageLoop())

	case opponentGameOverMsg:
		m.State = StateGameOver
		m.PlayerWon = msg.youWon
		return m, nil

	case opponentLeftMsg:
		m.Message = "Opponent disconnected."
		m.State = StateMenu // Or game over
		m.cleanup()
		return m, nil
	}

	return m, nil
}

// View implements tea.Model
func (m Model) View() string {
	switch m.State {
	case StateMenu:
		return m.renderMenuWithSelection()
	case StateMPMenu:
		return m.renderMPMenu()
	case StatePlacement:
		return m.renderPlacement()
	case StateBattle:
		return m.renderBattle()
	case StateGameOver:
		return m.renderGameOver()
	case StateMPHostWaiting:
		return m.renderMPHostWaiting()
	case StateMPJoinInput:
		return m.renderMPJoinInput()
	case StateMPConnecting:
		return fmt.Sprintf("Connecting to server at %s...\n\n%s", m.ServerAddress, m.Message)
	case StateMPPlacement:
		return m.renderMPPlacement()
	case StateMPWaitingForOpponent:
		return m.renderMPWaiting()
	case StateMPBattle:
		return m.renderMPBattle()
	default:
		return "Unknown state"
	}
}

// cleanup closes network connections
func (m *Model) cleanup() {
	if m.Connection != nil {
		m.Connection.Close()
	}
}

// updateMenu handles input during the menu state
func (m Model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.MenuSelection > 0 {
			m.MenuSelection--
		}
	case "down", "j":
		if m.MenuSelection < 1 { // Reduced since we moved host/join to submenu
			m.MenuSelection++
		}
	case "enter":
		switch m.MenuSelection {
		case 0: // vs AI
			m.GameMode = ModeVsAI
			m.State = StatePlacement
		case 1: // Multiplayer
			m.GameMode = ModeMultiplayer
			m.State = StateMPMenu
			m.MenuSelection = 0 // Reset for submenu
		}
	}
	return m, nil
}

func (m Model) updateMPMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.State = StateMenu
		m.MenuSelection = 0
		return m, nil
	case "up", "k":
		if m.MenuSelection > 0 {
			m.MenuSelection--
		}
	case "down", "j":
		if m.MenuSelection < 1 {
			m.MenuSelection++
		}
	case "enter":
		switch m.MenuSelection {
		case 0: // Host
			m.IsHost = true
			m.State = StateMPConnecting
			return m, m.connectAndCreateRoom()
		case 1: // Join
			m.IsHost = false
			m.State = StateMPJoinInput
			m.RoomCode = ""
		}
	}
	return m, nil
}

// updatePlacement handles ship placement in single player
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
		m.cleanup()
		return NewModel(), nil
	}
	return m, nil
}

// ========== Multiplayer Methods ==========

// Message types for async operations
type connectionEstablishedMsg struct {
	conn *bnet.Connection
}

type connectionErrorMsg struct {
	err error
}

type roomCreatedMsg struct {
	code string
}

type playerJoinedMsg struct{}

type joinErrorMsg struct {
	err string
}

type opponentReadyMsg struct{}

type opponentAttackMsg struct {
	row int
	col int
}

type attackResultMsg struct {
	hit          bool
	sunkShipName string
}

type opponentGameOverMsg struct {
	youWon bool
}

type opponentLeftMsg struct{}

// connectAndCreateRoom connects to server and requests a room
func (m Model) connectAndCreateRoom() tea.Cmd {
	return func() tea.Msg {
		conn, err := bnet.Connect(m.ServerAddress)
		if err != nil {
			return connectionErrorMsg{err: err}
		}

		// Send create room request
		if err := conn.Send(bnet.MsgCreateRoom, struct{}{}); err != nil {
			return connectionErrorMsg{err: err}
		}

		return connectionEstablishedMsg{conn: conn}
	}
}

// updateMPHostWaiting handles input while waiting for connection
func (m Model) updateMPHostWaiting(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.cleanup()
		m.State = StateMPMenu
		return m, nil
	}
	return m, nil
}

// updateMPJoinInput handles input while entering room code
func (m Model) updateMPJoinInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.State = StateMPMenu
		return m, nil
	case "enter":
		m.State = StateMPConnecting
		m.Message = "Connecting..."
		return m, func() tea.Msg {
			conn, err := bnet.Connect(m.ServerAddress)
			if err != nil {
				return connectionErrorMsg{err: err}
			}

			// Join room
			if err := conn.Send(bnet.MsgJoinRoom, bnet.JoinRoomPayload{Code: m.RoomCode}); err != nil {
				return connectionErrorMsg{err: err}
			}

			return connectionEstablishedMsg{conn: conn}
		}
	case "backspace":
		if len(m.RoomCode) > 0 {
			m.RoomCode = m.RoomCode[:len(m.RoomCode)-1]
		}
	default:
		// Uppercase only
		if len(msg.String()) == 1 && len(m.RoomCode) < 4 {
			char := msg.String()
			// Very basic implementation, ideally force upper
			m.RoomCode += char
		}
	}
	return m, nil
}

// handleConnectionEstablished handles successful connection
func (m Model) handleConnectionEstablished(msg connectionEstablishedMsg) (tea.Model, tea.Cmd) {
	m.Connection = msg.conn

	// Start the message loop
	return m, m.messageLoop()
}

// messageLoop continuously reads messages from the connection
func (m Model) messageLoop() tea.Cmd {
	return func() tea.Msg {
		msg, err := m.Connection.Receive()
		if err != nil {
			// Assuming naive disconnect
			return connectionErrorMsg{err: err}
		}

		switch msg.Type {
		case bnet.MsgRoomCreated:
			payload, _ := bnet.ParseCreateRoomResponse(msg.Payload)
			return roomCreatedMsg{code: payload.Code}

		case bnet.MsgJoinError:
			payload, _ := bnet.ParseErrorPayload(msg.Payload)
			return joinErrorMsg{err: payload.Message}

		case bnet.MsgGameStart: // Guest joined or Host notified
			return playerJoinedMsg{}

		case bnet.MsgPlayerJoined: // Host notified
			return playerJoinedMsg{}

		case bnet.MsgShipsPlaced:
			return opponentReadyMsg{}

		case bnet.MsgAttack:
			payload, _ := bnet.ParseAttackPayload(msg.Payload)
			return opponentAttackMsg{row: payload.Row, col: payload.Col}

		case bnet.MsgAttackResult:
			payload, _ := bnet.ParseAttackResultPayload(msg.Payload)
			return attackResultMsg{hit: payload.Hit, sunkShipName: payload.SunkShipName}

		case bnet.MsgGameOver:
			payload, _ := bnet.ParseGameOverPayload(msg.Payload)
			return opponentGameOverMsg{youWon: payload.YouWon}

		case bnet.MsgOpponentLeft:
			return opponentLeftMsg{}
		}

		// Continue loop if not handled or non-terminal
		// We rely on the Update function to restart the loop for valid messages
		return nil
	}
}

func (m Model) waitForPlayerJoined() tea.Cmd {
	// Already started message loop in handleConnectionEstablished
	return nil
}

func (m Model) startGame() (tea.Model, tea.Cmd) {
	m.State = StateMPPlacement
	m.Message = "Connected! Place your ships."
	m.CursorRow = 0
	m.CursorCol = 0
	return m, nil
}

// updateMPPlacement handles ship placement in multiplayer
func (m Model) updateMPPlacement(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
				// All ships placed, notify opponent
				m.ShipsPlaced = true
				m.Connection.Send(bnet.MsgShipsPlaced, struct{}{})

				if m.OpponentReady {
					// Both ready, start battle
					m.State = StateMPBattle
					m.CursorRow = 0
					m.CursorCol = 0
					m.PlayerTurn = m.IsHost
					if m.PlayerTurn {
						m.Message = "Battle begins! Your turn."
					} else {
						m.Message = "Battle begins! Opponent's turn."
					}
				} else {
					m.State = StateMPWaitingForOpponent
					m.Message = "Ships placed! Waiting for opponent..."
					// Loop is already running
					return m, nil
				}
			}
		}
	}
	return m, nil
}

// updateMPWaiting handles waiting for opponent
func (m Model) updateMPWaiting(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Just waiting, no input needed
	return m, nil
}

// handleOpponentReady handles when opponent finishes placement
func (m Model) handleOpponentReady() (tea.Model, tea.Cmd) {
	m.OpponentReady = true
	if m.ShipsPlaced {
		m.State = StateMPBattle
		m.CursorRow = 0
		m.CursorCol = 0
		m.PlayerTurn = m.IsHost
		if m.PlayerTurn {
			m.Message = "Battle begins! Your turn."
		} else {
			m.Message = "Battle begins! Opponent's turn."
		}
	}
	return m, nil
}

// updateMPBattle handles multiplayer battle input
func (m Model) updateMPBattle(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		// Check if already attacked
		cell := m.OpponentBoard.Cells[m.CursorRow][m.CursorCol]
		if cell == game.Hit || cell == game.Miss {
			m.Message = "Already attacked this location!"
			return m, nil
		}

		// Send attack to opponent
		m.Connection.Send(bnet.MsgAttack, bnet.AttackPayload{
			Row: m.CursorRow,
			Col: m.CursorCol,
		})
		m.LastAttackRow = m.CursorRow
		m.LastAttackCol = m.CursorCol

		// Wait for result
		// Loop running
		return m, nil
	}
	return m, nil
}

// handleOpponentAttack processes an attack from opponent
func (m Model) handleOpponentAttack(msg opponentAttackMsg) (tea.Model, tea.Cmd) {
	hit, _, sunkShipName := m.PlayerBoard.Attack(msg.row, msg.col)

	// Send result back
	m.Connection.Send(bnet.MsgAttackResult, bnet.AttackResultPayload{
		Row:          msg.row,
		Col:          msg.col,
		Hit:          hit,
		SunkShipName: sunkShipName,
	})

	if hit {
		if sunkShipName != "" {
			m.Message = "Opponent sunk your " + sunkShipName + "!"
		} else {
			m.Message = "Opponent hit your ship!"
		}

		if m.PlayerBoard.AllShipsSunk() {
			m.Connection.Send(bnet.MsgGameOver, bnet.GameOverPayload{YouWon: true})
			m.State = StateGameOver
			m.PlayerWon = false
			return m, nil
		}
	} else {
		m.Message = "Opponent missed!"
	}

	// Now it's our turn
	m.PlayerTurn = true
	m.Message += " Your turn."
	return m, nil
}

// handleAttackResult processes the result of our attack
func (m Model) handleAttackResult(msg attackResultMsg) (tea.Model, tea.Cmd) {
	// Update our view of opponent's board
	if msg.hit {
		m.OpponentBoard.Cells[m.LastAttackRow][m.LastAttackCol] = game.Hit
		if msg.sunkShipName != "" {
			m.Message = "HIT! You sunk their " + msg.sunkShipName + "!"
		} else {
			m.Message = "HIT!"
		}
	} else {
		m.OpponentBoard.Cells[m.LastAttackRow][m.LastAttackCol] = game.Miss
		m.Message = "Miss..."
	}

	// Switch turns
	m.PlayerTurn = false
	m.Message += " Opponent's turn."

	return m, nil
}
