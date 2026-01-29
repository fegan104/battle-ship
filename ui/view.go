package ui

import (
	"fmt"
	"strings"

	"battle-ship/game"

	"github.com/charmbracelet/lipgloss"
)

// renderMenu renders the main menu screen
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
	ship := m.ShipsToPlace[m.CurrentShipIndex]
	orientation := "Horizontal"
	if !m.PlacingHorizontal {
		orientation = "Vertical"
	}
	
	shipInfo := fmt.Sprintf("Placing: %s (length: %d) - %s", ship.Name, ship.Length, orientation)
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
	currentShip := m.ShipsToPlace[m.CurrentShipIndex]
	previewPositions := getPreviewPositions(currentShip.Length, m.CursorRow, m.CursorCol, m.PlacingHorizontal)
	canPlace := m.PlayerBoard.CanPlaceShip(currentShip, m.CursorRow, m.CursorCol, m.PlacingHorizontal)

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
	enemyBoard := m.renderEnemyBoard()

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
func (m Model) renderEnemyBoard() string {
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
			cell := m.AIBoard.Cells[r][c]
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
