package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

// Message types matching the client's expectation
const (
	MsgCreateRoom   = "create_room"
	MsgRoomCreated  = "room_created"
	MsgJoinRoom     = "join_room"
	MsgPlayerJoined = "player_joined" // Notify host
	MsgJoinError    = "join_error"
	MsgGameStart    = "game_start"
	MsgOpponentLeft = "opponent_left"
)

// Message is the wrapper for all network messages types.
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// CreateRoomResponse is the response payload for creating a room.
type CreateRoomResponse struct {
	Code string `json:"code"`
}

// JoinRoomPayload is the request payload for joining a room.
type JoinRoomPayload struct {
	Code string `json:"code"`
}

// ErrorPayload is the payload for error messages.
type ErrorPayload struct {
	Message string `json:"message"`
}

// Client represents a connected player's WebSocket connection and state.
type Client struct {
	conn   *websocket.Conn
	room   *Room
	isHost bool
}

// Room represents a game session between two players.
type Room struct {
	Code  string
	Host  *Client
	Guest *Client
	mu    sync.Mutex
}

// Server manages active rooms and concurrency.
type Server struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

var server = &Server{
	rooms: make(map[string]*Room),
}

func main() {
	http.HandleFunc("/ws", handleConnections)

	port := 8080
	fmt.Printf("Server started on :%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	client := &Client{conn: ws}

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			handleDisconnect(client)
			break
		}

		handleMessage(client, msg)
	}
}

func handleMessage(client *Client, msg Message) {
	switch msg.Type {
	case MsgCreateRoom:
		handleCreateRoom(client)
	case MsgJoinRoom:
		var payload JoinRoomPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			sendError(client, "Invalid payload")
			return
		}
		handleJoinRoom(client, payload.Code)
	default:
		// Relay game messages if in a room
		if client.room != nil {
			relayMessage(client, msg)
		}
	}
}

func handleCreateRoom(client *Client) {
	code := generateRoomCode()
	room := &Room{
		Code: code,
		Host: client,
	}

	server.mu.Lock()
	server.rooms[code] = room
	server.mu.Unlock()

	client.room = room
	client.isHost = true

	// Send room code back to host
	payload, _ := json.Marshal(CreateRoomResponse{Code: code})
	response := Message{
		Type:    MsgRoomCreated,
		Payload: payload,
	}
	client.conn.WriteJSON(response)

	log.Printf("Room created: %s", code)
}

func handleJoinRoom(client *Client, code string) {
	code = strings.ToUpper(code)
	server.mu.Lock()
	room, exists := server.rooms[code]
	server.mu.Unlock()

	if !exists {
		sendError(client, "Room not found")
		return
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	if room.Guest != nil {
		sendError(client, "Room is full")
		return
	}

	room.Guest = client
	client.room = room
	client.isHost = false

	// Notify Guest they joined
	client.conn.WriteJSON(Message{Type: MsgGameStart, Payload: []byte("{}")})

	// Notify Host that Guest joined
	if room.Host != nil {
		room.Host.conn.WriteJSON(Message{Type: MsgPlayerJoined, Payload: []byte("{}")})
	}

	log.Printf("Player joined room: %s", code)
}

func relayMessage(sender *Client, msg Message) {
	room := sender.room
	if room == nil {
		return
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	var target *Client
	if sender == room.Host {
		target = room.Guest
	} else {
		target = room.Host
	}

	if target != nil {
		target.conn.WriteJSON(msg)
	}
}

func handleDisconnect(client *Client) {
	if client.room == nil {
		return
	}

	room := client.room
	room.mu.Lock()
	defer room.mu.Unlock()

	// Notify other player
	var target *Client
	if client == room.Host {
		target = room.Guest
	} else {
		target = room.Host
	}

	if target != nil {
		target.conn.WriteJSON(Message{Type: MsgOpponentLeft, Payload: []byte("{}")})
		target.room = nil // Unlink them so they can join another game? Or just end session.
	}

	// Remove room
	server.mu.Lock()
	delete(server.rooms, room.Code)
	server.mu.Unlock()
}

func sendError(client *Client, message string) {
	payload, _ := json.Marshal(ErrorPayload{Message: message})
	client.conn.WriteJSON(Message{
		Type:    MsgJoinError,
		Payload: payload,
	})
}

func generateRoomCode() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 4)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
