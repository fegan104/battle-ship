package net

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// MessageType identifies the type of network message
type MessageType string

const (
	// Game Messages
	MsgShipsPlaced  MessageType = "ships_placed"
	MsgAttack       MessageType = "attack"
	MsgAttackResult MessageType = "attack_result"
	MsgGameOver     MessageType = "game_over"

	// Control Messages
	MsgCreateRoom   MessageType = "create_room"
	MsgRoomCreated  MessageType = "room_created"
	MsgJoinRoom     MessageType = "join_room"
	MsgPlayerJoined MessageType = "player_joined"
	MsgJoinError    MessageType = "join_error"
	MsgGameStart    MessageType = "game_start"
	MsgOpponentLeft MessageType = "opponent_left"
)

// Message is the wrapper for all network messages
type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// Payloads

type CreateRoomResponse struct {
	Code string `json:"code"`
}

type JoinRoomPayload struct {
	Code string `json:"code"`
}

type ErrorPayload struct {
	Message string `json:"message"`
}

type AttackPayload struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type AttackResultPayload struct {
	Row          int    `json:"row"`
	Col          int    `json:"col"`
	Hit          bool   `json:"hit"`
	SunkShipName string `json:"sunk_ship_name,omitempty"`
}

type GameOverPayload struct {
	YouWon bool `json:"you_won"`
}

// Connection wraps a WebSocket connection
type Connection struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

// Send sends a message over the connection
func (c *Connection) Send(msgType MessageType, payload interface{}) error { // Changed receiver to pointer
	c.mu.Lock()
	defer c.mu.Unlock()

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	msg := Message{
		Type:    msgType,
		Payload: payloadBytes,
	}

	return c.conn.WriteJSON(msg)
}

// Receive receives a message from the connection
func (c *Connection) Receive() (*Message, error) { // Changed receiver to pointer
	var msg Message
	err := c.conn.ReadJSON(&msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// Close closes the connection
func (c *Connection) Close() error { // Changed receiver to pointer
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.Close()
}

// Connect connects to the central server
func Connect(address string) (*Connection, error) {
	url := fmt.Sprintf("ws://%s/ws", address)
	log.Printf("Connecting to %s", url)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	return &Connection{
		conn: conn,
	}, nil
}

// Helper functions for parsing payloads

func ParseCreateRoomResponse(payload json.RawMessage) (*CreateRoomResponse, error) {
	var p CreateRoomResponse
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func ParseErrorPayload(payload json.RawMessage) (*ErrorPayload, error) {
	var p ErrorPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func ParseAttackPayload(payload json.RawMessage) (*AttackPayload, error) {
	var p AttackPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func ParseAttackResultPayload(payload json.RawMessage) (*AttackResultPayload, error) {
	var p AttackResultPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func ParseGameOverPayload(payload json.RawMessage) (*GameOverPayload, error) {
	var p GameOverPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	return &p, nil
}
