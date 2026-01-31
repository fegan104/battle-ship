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
Manages WebSocket communication for multiplayer.
- **`network.go`**: Implements a custom JSON-based protocol over WebSockets. Handles connection to the central server and message exchange (Room creation, Attacks, Results).

### `main.go`
The entry point that initializes the Bubble Tea program and starts the application.
- **`cmd/server/main.go`**: The central WebSocket server that manages game rooms and relays messages between players.

## How it Works

The application runs in a loop driven by the Bubble Tea framework:
1.  **Update**: Listens for keypresses or network messages. It updates the `Model` (e.g., moves cursor, recording a hit).
2.  **View**: Generates a string representation of the current `Model` to display in the terminal.

In **Multiplayer Mode**, clients connect to a central server via WebSockets.
- **Hosting**: A player creates a room and receives a unique 4-letter code.
- **Joining**: Another player enters that code to join the session.
- **Battle**: The server relays game messages (Attacks, Results) between the two connected players.

## How to Run

Ensure you have Go installed (1.21+ recommended).

### 1. Start the Server
The multiplayer feature requires the central server to be running.

```bash
go run cmd/server/main.go
```
*   The server listens on port `8080`.

### 2. Run the Game Client
Open a new terminal (or multiple for local testing) and run the game.

```bash
go run .
```

## How to Play

### Game Modes
1.  **Play vs AI**: Classic single-player mode against the computer.
2.  **Multiplayer**:
    *   **Host Game**: Create a new room and get a Room Code (e.g., `ABCD`).
    *   **Join Game**: Enter a Room Code to play against a friend.

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
