# Go Battleship

A terminal-based Battleship game written in Go, featuring a robust TUI (Terminal User Interface) and TCP-based multiplayer support.

## Project Structure

The project is organized into three main packages that handle distinct responsibilities:

### 1. `game/` (Core Logic)
Contains the platform-agnostic game rules and state.
- **`board.go`**: Manages the grid state (Hit, Miss, Empty, Ship), ship placement validation, and attack logic.
- **`ship.go`**: Defines ship types, lengths, and tracks their health/sunk status.
- **`ai.go`**: Implements a computer opponent with "hunt and sink" logic.

### 2. `ui/` (User Interface)
Handles the TUI using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework, following The Elm Architecture (Model-View-Update).
- **`model.go`**: The central state store. It holds the game boards, current state (Menu, Placement, Battle), and handles input events.
- **`view.go`**: Renders the UI strings. It draws the boards, ships, and menus using [Lipgloss](https://github.com/charmbracelet/lipgloss) for styling.
- **`styles.go`**: Defines the color palette and layout styles.

### 3. `net/` (Networking)
Manages TCP communication for multiplayer.
- **`network.go`**: Implements a custom JSON-based protocol over TCP. Handles connection establishment (Server/Client) and message exchange (Attacks, Results, Game Over).

### `main.go`
The entry point that initializes the Bubble Tea program and starts the application.

## How it Works

The application runs in a loop driven by the Bubble Tea framework:
1.  **Update**: Listens for keypresses or network messages. It updates the `Model` (e.g., moves cursor, recording a hit).
2.  **View**: Generates a string representation of the current `Model` to display in the terminal.

In **Multiplayer Mode**, the game uses an async message loop. When a player fires, an `Attack` message is sent over TCP. The opponent receives it, updates their board, and sends back an `AttackResult` (Hit/Miss/Sunk).

## How to Run

Ensure you have Go installed (1.21+ recommended).

### Run directly
```bash
go run .
```

### Build and Run
```bash
go build -o battleship .
./battleship
```

## How to Play

### Game Modes
1.  **Play vs AI**: Classic single-player mode against the computer.
2.  **Host Game**: Starts a TCP server on port `8080`. Wait for a friend to connect.
3.  **Join Game**: Connect to a hosting player by entering their IP address (or `localhost` for local testing).

### Controls

| Action | Keys |
|--------|------|
| **Navigation** | Arrow Keys or Vim Keys (`H`, `J`, `K`, `L`) |
| **Select / Fire** | `Enter` |
| **Rotate Ship** | `R` (Deployment phase only) |
| **Quit** | `Q` or `Ctrl+C` |

### Rules
1.  **Placement Phase**: Position your fleet of 5 ships. Ships cannot overlap.
2.  **Battle Phase**: Take turns firing at coordinates on the enemy map.
    - **Hit**: You get another turn (vs AI) or pass turn (Multiplayer, depending on rule implementation). *Note: standard turns alternate.*
    - **Sunk**: Destroy all enemy ships to win.
