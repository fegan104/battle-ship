package ui

import (
	"fmt"
	"strings"

	"battle-ship/game"

	"github.com/charmbracelet/lipgloss"
)

// Menu options
var menuOptions = []string{
	"Play vs AI",
	"Multiplayer", // Changed from "Host Game" to just "Multiplayer"
}

// renderMenuWithSelection renders the main menu with selection
func (m Model) renderMenuWithSelection() string {
	title := bigTitleStyle.Render(`
 ██████╗  █████╗ ████████╗████████╗██╗     ███████╗███████╗██╗  ██╗██╗██████╗ 
 ██╔══██╗██╔══██╗╚══██╔══╝╚══██╔══╝██║     ██╔════╝██╔════╝██║  ██║██║██╔══██╗
 ██████╔╝███████║   ██║      ██║   ██║     █████╗  ███████╗███████║██║██████╔╝
 ██╔══██╗██╔══██║   ██║      ██║   ██║     ██╔══╝  ╚════██║██╔══██║██║██╔═══╝ 
 ██████╔╝██║  ██║   ██║      ██║   ███████╗███████╗███████║██║  ██║██║██║     
 ╚═════╝ ╚═╝  ╚═╝   ╚═╝      ╚═╝   ╚══════╝╚══════╝╚══════╝╚═╝  ╚═╝╚═╝╚═╝     
`)

	subtitle := subtitleStyle.Render("A classic naval combat game")

	var menuItems strings.Builder
	menuItems.WriteString("\n\n")
	for i, option := range menuOptions {
		if i == m.MenuSelection {
			menuItems.WriteString(selectedMenuStyle.Render("▸ " + option))
		} else {
			menuItems.WriteString(menuItemStyle.Render("  " + option))
		}
		menuItems.WriteString("\n")
	}

	help := helpStyle.Render("\n↑↓: Select  |  Enter: Confirm  |  Q: Quit")

	errorMsg := ""
	if m.Message != "" {
		errorMsg = "\n\n" + messageStyle.Render(m.Message)
	}

	return containerStyle.Render(title + "\n" + subtitle + menuItems.String() + help + errorMsg)
}

// renderMPMenu renders the multiplayer submenu
func (m Model) renderMPMenu() string {
	title := titleStyle.Render("MULTIPLAYER")

	options := []string{
		"Host Game (Create Room)",
		"Join Game (Enter Code)",
	}

	var menuItems strings.Builder
	menuItems.WriteString("\n\n")
	for i, option := range options {
		if i == m.MenuSelection {
			menuItems.WriteString(selectedMenuStyle.Render("▸ " + option))
		} else {
			menuItems.WriteString(menuItemStyle.Render("  " + option))
		}
		menuItems.WriteString("\n")
	}

	help := helpStyle.Render("\n↑↓: Select  |  Enter: Confirm  |  Esc: Back")

	errorMsg := ""
	if m.Message != "" {
		errorMsg = "\n\n" + messageStyle.Render(m.Message)
	}

	return containerStyle.Render(title + menuItems.String() + help + errorMsg)
}

// renderMenu renders the main menu screen (legacy, now uses renderMenuWithSelection)
func renderMenu() string {
	title := bigTitleStyle.Render(`
 ██████╗  █████╗ ████████╗████████╗██╗     ███████╗███████╗██╗  ██╗██╗██████╗ 
 ██╔══██╗██╔══██╗╚══██╔══╝╚══██╔══╝██║     ██╔════╝██╔════╝██║  ██║██║██╔══██╗
 ██████╔╝███████║   ██║      ██║   ██║     █████╗  ███████╗███████║██║██████╔╝
 ██╔══██╗██╔══██║   ██║      ██║   ██║     ██╔══╝  ╚════██║██╔══██║██║██╔═══╝ 
 ██████╔╝██║  ██║   ██║      ██║   ███████╗███████╗███████║██║  ██║██║██║     
 ╚═════╝ ╚═╝  ╚═╝   ╚═╝      ╚═╝   ╚══════╝╚══════╝╚══════╝╚═╝  ╚═╝╚═╝╚═╝     
`)

	subtitle := subtitleStyle.Render("A classic naval combat game")

	instructions := helpStyle.Render("\n\n  Press ENTER to start\n  Press Q to quit")

	return containerStyle.Render(title + "\n" + subtitle + instructions)
}

// renderPlacement renders the ship placement phase
func (m Model) renderPlacement() string {
	var sb strings.Builder

	title := titleStyle.Render("PLACE YOUR SHIPS")
	sb.WriteString(title + "\n\n")

	// Current ship info
	var shipInfo string
	if m.CurrentShipIndex < len(m.ShipsToPlace) {
		ship := m.ShipsToPlace[m.CurrentShipIndex]
		orientation := "Horizontal"
		if !m.PlacingHorizontal {
			orientation = "Vertical"
		}

		shipInfo = fmt.Sprintf("Placing: %s (length: %d) - %s", ship.Name, ship.Length, orientation)
	} else {
		shipInfo = "All ships placed!"
	}
	sb.WriteString(messageStyle.Render(shipInfo) + "\n\n")

	// Render the board with placement preview
	sb.WriteString(m.renderPlacementBoard())

	// Instructions
	help := helpStyle.Render("\n↑↓←→: Move  |  R: Rotate  |  Enter: Place Ship  |  Q: Quit")
	sb.WriteString(help)

	return containerStyle.Render(sb.String())
}

// renderPlacementBoard renders the player board during ship placement
func (m Model) renderPlacementBoard() string {
	var sb strings.Builder

	// Column headers
	sb.WriteString("    ")
	for c := 0; c < game.BoardSize; c++ {
		sb.WriteString(headerStyle.Render(fmt.Sprintf(" %c ", 'A'+c)))
	}
	sb.WriteString("\n")

	// Get preview positions
	// Check if we are done placing
	var previewPositions [][2]int
	var canPlace bool

	if m.CurrentShipIndex < len(m.ShipsToPlace) {
		currentShip := m.ShipsToPlace[m.CurrentShipIndex]
		previewPositions = getPreviewPositions(currentShip.Length, m.CursorRow, m.CursorCol, m.PlacingHorizontal)
		canPlace = m.PlayerBoard.CanPlaceShip(currentShip, m.CursorRow, m.CursorCol, m.PlacingHorizontal)
	}

	// Rows
	for r := 0; r < game.BoardSize; r++ {
		sb.WriteString(headerStyle.Render(fmt.Sprintf(" %2d ", r+1)))
		for c := 0; c < game.BoardSize; c++ {
			cell := m.PlayerBoard.Cells[r][c]

			// Check if this is a preview position
			isPreview := false
			for _, pos := range previewPositions {
				if pos[0] == r && pos[1] == c {
					isPreview = true
					break
				}
			}

			if isPreview {
				if canPlace {
					sb.WriteString(validPreviewCell.Render("█"))
				} else {
					sb.WriteString(invalidPreviewCell.Render("█"))
				}
			} else if cell == game.ShipCell {
				sb.WriteString(shipCell.Render("█"))
			} else {
				sb.WriteString(waterCell.Render("~"))
			}
		}
		sb.WriteString("\n")
	}

	return boardStyle.Render(sb.String())
}

// renderBattle renders the battle phase with both boards
func (m Model) renderBattle() string {
	var sb strings.Builder

	title := titleStyle.Render("BATTLE!")
	sb.WriteString(title + "\n")

	// Turn indicator
	turnText := "YOUR TURN - Select a target"
	if !m.PlayerTurn {
		turnText = "ENEMY'S TURN..."
	}
	sb.WriteString(messageStyle.Render(turnText) + "\n\n")

	// Render both boards side by side
	playerBoard := m.renderPlayerBoardBattle()
	enemyBoard := m.renderEnemyBoard(m.AIBoard)

	boards := lipgloss.JoinHorizontal(lipgloss.Top, playerBoard, "    ", enemyBoard)
	sb.WriteString(boards)

	// Message
	if m.Message != "" {
		sb.WriteString("\n" + messageStyle.Render(m.Message))
	}

	// Instructions
	help := helpStyle.Render("\n↑↓←→: Move cursor  |  Enter: Fire  |  Q: Quit")
	sb.WriteString(help)

	return containerStyle.Render(sb.String())
}

// renderPlayerBoardBattle renders the player's board showing their ships
func (m Model) renderPlayerBoardBattle() string {
	var sb strings.Builder

	sb.WriteString(boardTitleStyle.Render("YOUR FLEET") + "\n")

	// Column headers
	sb.WriteString("    ")
	for c := 0; c < game.BoardSize; c++ {
		sb.WriteString(headerStyle.Render(fmt.Sprintf(" %c ", 'A'+c)))
	}
	sb.WriteString("\n")

	// Rows
	for r := 0; r < game.BoardSize; r++ {
		sb.WriteString(headerStyle.Render(fmt.Sprintf(" %2d ", r+1)))
		for c := 0; c < game.BoardSize; c++ {
			cell := m.PlayerBoard.Cells[r][c]
			switch cell {
			case game.Hit:
				sb.WriteString(hitCell.Render("X"))
			case game.Miss:
				sb.WriteString(missCell.Render("•"))
			case game.ShipCell:
				sb.WriteString(shipCell.Render("█"))
			default:
				sb.WriteString(waterCell.Render("~"))
			}
		}
		sb.WriteString("\n")
	}

	return boardStyle.Render(sb.String())
}

// renderEnemyBoard renders the enemy board (hiding ship positions)
func (m Model) renderEnemyBoard(board *game.Board) string {
	var sb strings.Builder

	sb.WriteString(boardTitleStyle.Render("ENEMY WATERS") + "\n")

	// Column headers
	sb.WriteString("    ")
	for c := 0; c < game.BoardSize; c++ {
		sb.WriteString(headerStyle.Render(fmt.Sprintf(" %c ", 'A'+c)))
	}
	sb.WriteString("\n")

	// Rows
	for r := 0; r < game.BoardSize; r++ {
		sb.WriteString(headerStyle.Render(fmt.Sprintf(" %2d ", r+1)))
		for c := 0; c < game.BoardSize; c++ {
			cell := board.Cells[r][c]
			isCursor := m.PlayerTurn && r == m.CursorRow && c == m.CursorCol

			if isCursor {
				sb.WriteString(cursorCell.Render("◎"))
			} else {
				switch cell {
				case game.Hit:
					sb.WriteString(hitCell.Render("X"))
				case game.Miss:
					sb.WriteString(missCell.Render("•"))
				default:
					// Don't reveal enemy ships
					sb.WriteString(waterCell.Render("~"))
				}
			}
		}
		sb.WriteString("\n")
	}

	return boardStyle.Render(sb.String())
}

// renderGameOver renders the game over screen
func (m Model) renderGameOver() string {
	var sb strings.Builder

	var resultText string
	if m.PlayerWon {
		resultText = successStyle.Render(`
 ██╗   ██╗██╗ ██████╗████████╗ ██████╗ ██████╗ ██╗   ██╗██╗
 ██║   ██║██║██╔════╝╚══██╔══╝██╔═══██╗██╔══██╗╚██╗ ██╔╝██║
 ██║   ██║██║██║        ██║   ██║   ██║██████╔╝ ╚████╔╝ ██║
 ╚██╗ ██╔╝██║██║        ██║   ██║   ██║██╔══██╗  ╚██╔╝  ╚═╝
  ╚████╔╝ ██║╚██████╗   ██║   ╚██████╔╝██║  ██║   ██║   ██╗
   ╚═══╝  ╚═╝ ╚═════╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚═╝
`)
		sb.WriteString(resultText + "\n")
		sb.WriteString(successStyle.Render("You sunk all enemy ships!") + "\n")
	} else {
		resultText = errorStyle.Render(`
 ██████╗ ███████╗███████╗███████╗ █████╗ ████████╗
 ██╔══██╗██╔════╝██╔════╝██╔════╝██╔══██╗╚══██╔══╝
 ██║  ██║█████╗  █████╗  █████╗  ███████║   ██║   
 ██║  ██║██╔══╝  ██╔══╝  ██╔══╝  ██╔══██║   ██║   
 ██████╔╝███████╗██║     ███████╗██║  ██║   ██║   
 ╚═════╝ ╚══════╝╚═╝     ╚══════╝╚═╝  ╚═╝   ╚═╝   
`)
		sb.WriteString(resultText + "\n")
		sb.WriteString(errorStyle.Render("The enemy sunk all your ships!") + "\n")
	}

	help := helpStyle.Render("\nPress ENTER to play again  |  Press Q to quit")
	sb.WriteString(help)

	return containerStyle.Render(sb.String())
}

// ========== Multiplayer Views ==========

// renderMPHostWaiting renders the waiting for connection screen
func (m Model) renderMPHostWaiting() string {
	title := titleStyle.Render("HOSTING GAME")

	waiting := messageStyle.Render(fmt.Sprintf("\n\nWaiting for opponent to join Room: %s", m.RoomCode))
	hint := helpStyle.Render("\nOther player should select 'Join Game' and enter this code.")

	help := helpStyle.Render("\n\nPress ESC to cancel")

	return containerStyle.Render(title + waiting + hint + help)
}

// renderMPJoinInput renders the join game input screen
func (m Model) renderMPJoinInput() string {
	title := titleStyle.Render("JOIN GAME")

	prompt := messageStyle.Render("\n\nEnter Room Code:")

	address := inputStyle.Render("\n" + m.RoomCode + "█")

	help := helpStyle.Render("\n\nPress ENTER to connect  |  ESC to cancel")

	errorMsg := ""
	if m.Message != "" {
		errorMsg = "\n\n" + messageStyle.Render(m.Message)
	}

	return containerStyle.Render(title + prompt + address + help + errorMsg)
}

// renderMPPlacement renders multiplayer ship placement
func (m Model) renderMPPlacement() string {
	var sb strings.Builder

	modeLabel := ""
	if m.IsHost {
		modeLabel = " (HOST)"
	} else {
		modeLabel = " (JOINED)"
	}

	title := titleStyle.Render("PLACE YOUR SHIPS" + modeLabel)
	sb.WriteString(title + "\n\n")

	// Current ship info
	var shipInfo string
	if m.CurrentShipIndex < len(m.ShipsToPlace) {
		ship := m.ShipsToPlace[m.CurrentShipIndex]
		orientation := "Horizontal"
		if !m.PlacingHorizontal {
			orientation = "Vertical"
		}
		shipInfo = fmt.Sprintf("Placing: %s (length: %d) - %s", ship.Name, ship.Length, orientation)
	} else {
		shipInfo = "All ships placed! Waiting for opponent..."
	}

	sb.WriteString(messageStyle.Render(shipInfo) + "\n\n")

	// Render the board with placement preview
	sb.WriteString(m.renderPlacementBoard())

	// Instructions
	help := helpStyle.Render("\n↑↓←→: Move  |  R: Rotate  |  Enter: Place Ship  |  Q: Quit")
	sb.WriteString(help)

	if m.Message != "" {
		sb.WriteString("\n" + messageStyle.Render(m.Message))
	}

	return containerStyle.Render(sb.String())
}

// renderMPWaiting renders waiting for opponent to place ships (when we are done but they aren't)
func (m Model) renderMPWaiting() string {
	title := titleStyle.Render("SHIPS PLACED!")

	waiting := messageStyle.Render("\n\nWaiting for opponent to place their ships...")

	// Show player's board
	board := "\n\n" + m.renderPlayerBoardBattle()

	help := helpStyle.Render("\n\nPlease wait...")

	if m.Message != "" {
		help += "\n" + messageStyle.Render(m.Message)
	}

	return containerStyle.Render(title + waiting + board + help)
}

// renderMPBattle renders the multiplayer battle phase
func (m Model) renderMPBattle() string {
	var sb strings.Builder

	modeLabel := ""
	if m.IsHost {
		modeLabel = " (HOST)"
	} else {
		modeLabel = " (JOINED)"
	}

	title := titleStyle.Render("MULTIPLAYER BATTLE!" + modeLabel)
	sb.WriteString(title + "\n")

	// Turn indicator
	var turnText string
	if m.PlayerTurn {
		turnText = "YOUR TURN - Select a target"
	} else {
		turnText = "OPPONENT'S TURN..."
	}
	sb.WriteString(messageStyle.Render(turnText) + "\n\n")

	// Render both boards side by side
	playerBoard := m.renderPlayerBoardBattle()
	opponentBoard := m.renderEnemyBoard(m.OpponentBoard)

	boards := lipgloss.JoinHorizontal(lipgloss.Top, playerBoard, "    ", opponentBoard)
	sb.WriteString(boards)

	// Message
	if m.Message != "" {
		sb.WriteString("\n" + messageStyle.Render(m.Message))
	}

	// Instructions
	help := helpStyle.Render("\n↑↓←→: Move cursor  |  Enter: Fire  |  Q: Quit")
	sb.WriteString(help)

	return containerStyle.Render(sb.String())
}

// renderMPOpponentBoard renders the opponent's board in multiplayer
// Reused renderEnemyBoard by making it accept a Board param instead of defaulting to AIBoard, or just adding this wrapper.
// I updated renderEnemyBoard earlier to take a board param.
// But wait, my previous renderEnemyBoard implementation in view.go hardcoded m.AIBoard. I need to check if I updated it in this write.
// YES, I updated func (m Model) renderEnemyBoard(board *game.Board) string in this file content.

// getPreviewPositions returns the positions where a ship would be placed
func getPreviewPositions(length, row, col int, horizontal bool) [][2]int {
	positions := make([][2]int, 0, length)
	for i := 0; i < length; i++ {
		var r, c int
		if horizontal {
			r, c = row, col+i
		} else {
			r, c = row+i, col
		}
		if r >= 0 && r < game.BoardSize && c >= 0 && c < game.BoardSize {
			positions = append(positions, [2]int{r, c})
		}
	}
	return positions
}
