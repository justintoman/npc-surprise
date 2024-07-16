package stream

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/justintoman/npc-surprise/db"
)

// This file is heavily inspired by
// https://github.com/gin-gonic/examples/blob/master/server-sent-event/main.go

type StreamingServer interface {
	NewUserStream() (gin.HandlerFunc, StreamHandlerFunc)
	Close(ClientChan)
	SendMessage(clientId string, message MessagePayload)
	Listen(context.Context)
	GetClients() []UserClient
}

type UserClient struct {
	IsAdmin bool   `json:"is_admin"`
	Id      string `json:"id"`
	Name    string `json:"name"`
}

func New(db db.Db) StreamingServer {
	eventStream := &EventStream{
		Db:            db,
		Message:       make(chan Message),
		NewClients:    make(chan ClientChan),
		ClosedClients: make(chan ClientChan),
		TotalClients:  make(map[string]ClientChan),
	}
	return eventStream
}

func (stream *EventStream) Close(clientChan ClientChan) {
	stream.ClosedClients <- clientChan
}

func (stream *EventStream) SendMessage(clientId string, message MessagePayload) {
	slog.Info("Sending message", "clientId", clientId, "message", message)
	stream.Message <- Message{PlayerId: clientId, Payload: message}
}

func (stream *EventStream) GetClients() []UserClient {
	clients := make([]UserClient, 0, len(stream.TotalClients))
	for _, client := range stream.TotalClients {
		clients = append(clients, client.UserClient)
	}
	return clients
}

type ClientChan struct {
	UserClient
	Channel chan MessagePayload
}

type Message struct {
	PlayerId string
	Payload  MessagePayload
}

type MessagePayload struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type EventStream struct {
	Db            db.Db
	Message       chan Message
	NewClients    chan ClientChan
	ClosedClients chan ClientChan
	// map of clients by string topic
	TotalClients map[string]ClientChan
}

func (stream *EventStream) Listen(c context.Context) {
	for {
		select {
		// Add new available client
		case client := <-stream.NewClients:
			if client.IsAdmin {
				stream.TotalClients["admin"] = client
				stream.sendAdminConnectedData(client.Channel)
			} else {
				stream.TotalClients[client.Id] = client
				stream.sendPlayerConnectedData(client.Id)
				admin, ok := stream.TotalClients["admin"]
				if !ok {
					slog.Error("Could not find admin client")
					continue
				}
				type PlayerConnectedPayload struct {
					Id       int    `json:"id"`
					Name     string `json:"name"`
					IsOnline bool   `json:"is_online"`
				}
				id, err := strconv.Atoi(client.Id)
				if err != nil {
					slog.Error("Could not convert player id to int", "error", err)
					continue
				}
				admin.Channel <- MessagePayload{Type: "player-connected", Data: PlayerConnectedPayload{Id: id, Name: client.Name, IsOnline: true}}
			}
			names := make([]string, 0)
			for _, client := range stream.TotalClients {
				names = append(names, client.Name)
			}
			slog.Info(fmt.Sprintf("Client added. %d registered clients. Clients: %+v", len(stream.TotalClients), names))

		// Remove closed client
		case client := <-stream.ClosedClients:
			slog.Info("Client closed", "id", client.Id)
			if client.IsAdmin {
				delete(stream.TotalClients, "admin")
				continue
			}
			delete(stream.TotalClients, client.Id)
			close(client.Channel)

			// tell the admin a user left
			admin, ok := stream.TotalClients["admin"]
			if !ok {
				// no admin client connected
				continue
			}
			playerId, err := strconv.Atoi(client.UserClient.Id)
			if err != nil {
				slog.Error("Could not convert player id to int", "error", err)
				continue
			}
			admin.Channel <- MessagePayload{Type: "player-disconnected", Data: playerId}
			slog.Info(fmt.Sprintf("Removed client. %d registered clients", len(stream.TotalClients)))

		// Broadcast message to client
		case eventMsg := <-stream.Message:
			client, ok := stream.TotalClients[eventMsg.PlayerId]
			if !ok {
				slog.Error(fmt.Sprintf("Attempted to send message to a client that doesn't exist. Id: %s", eventMsg.PlayerId))
				continue
			}

			client.Channel <- eventMsg.Payload
		case <-c.Done():
			return
		}
	}
}

type StreamHandlerFunc func(*gin.Context) (bool, error)

func (stream *EventStream) NewUserStream() (gin.HandlerFunc, StreamHandlerFunc) {
	middleware := func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")

		id, err := c.Cookie("player_id")

		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"message": "Missing player_id cookie", "status": 400})
			return
		}

		name, err := c.Cookie("player_name")
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"message": "Missing player_id cookie", "status": 400})
			return
		}

		isAdmin := c.MustGet("isAdmin").(bool)
		clientChan := ClientChan{
			Channel: make(chan MessagePayload),
			UserClient: UserClient{
				IsAdmin: isAdmin,
				Id:      id,
				Name:    name,
			},
		}
		stream.NewClients <- clientChan

		go func() {
			<-c.Writer.CloseNotify()
			stream.Close(clientChan)
		}()

		c.Set("clientChan", clientChan)
		c.Next()
	}
	handler := func(c *gin.Context) (bool, error) {
		v, ok := c.Get("clientChan")
		if !ok {
			return false, fmt.Errorf("clientChan not found. Is Middleware applied?")
		}
		clientChan, ok := v.(ClientChan)
		if !ok {
			return false, fmt.Errorf("clientChan not a ClientChan. Is Middleware applied?")
		}
		c.Stream(func(w io.Writer) bool {
			// Stream message to client from message channel
			if msg, ok := <-clientChan.Channel; ok {
				c.SSEvent("message", msg)
				return true
			}
			return false
		})
		return true, nil
	}
	return middleware, handler
}

func (stream *EventStream) sendPlayerConnectedData(player_id string) error {
	client, ok := stream.TotalClients[player_id]
	if !ok {
		return fmt.Errorf("client not found for player_id %s", player_id)
	}

	characters, err := stream.Db.Character.GetAllByPlayerId(player_id)
	if err != nil {
		return err
	}

	charsWithActions := make([]map[string]any, len(characters))
	for i, character := range characters {
		actions, err := stream.Db.Action.GetAll(character.Id)
		if err != nil {
			return err
		}
		charsWithActions[i] = map[string]any{
			"id":          character.Id,
			"name":        character.Name,
			"race":        character.Race,
			"gender":      character.Gender,
			"age":         character.Age,
			"description": character.Description,
			"appearance":  character.Appearance,
			"actions":     actions,
		}
	}
	client.Channel <- MessagePayload{Type: "connected", Data: map[string]any{
		"players":    make([]map[string]any, 0),
		"characters": charsWithActions,
	}}
	return nil
}

func (stream *EventStream) sendAdminConnectedData(channel chan MessagePayload) error {

	characters, err := stream.Db.Character.GetAll()
	if err != nil {
		return err
	}

	charsWithActions := make([]map[string]any, len(characters))
	for i, character := range characters {
		actions, err := stream.Db.Action.GetAll(character.Id)
		if err != nil {
			return err
		}
		charsWithActions[i] = map[string]any{
			"id":          character.Id,
			"name":        character.Name,
			"race":        character.Race,
			"gender":      character.Gender,
			"age":         character.Age,
			"description": character.Description,
			"appearance":  character.Appearance,
			"actions":     actions,
		}
	}
	dbPlayers, err := stream.Db.Player.GetAll()
	if err != nil {
		return err
	}
	onlinePlayers := stream.GetClients()

	type PlayerResponse struct {
		db.Player
		IsOnline bool `json:"is_online"`
	}

	result := make([]PlayerResponse, len(*dbPlayers))
	for i, player := range *dbPlayers {
		result[i] = PlayerResponse{Player: player}
		for _, onlinePlayer := range onlinePlayers {
			if result[i].IsOnline {
				continue
			}
			if strconv.Itoa(player.Id) == onlinePlayer.Id {
				result[i].IsOnline = true
			}
		}
	}
	channel <- MessagePayload{Type: "connected", Data: map[string]any{
		"players":    result,
		"characters": charsWithActions,
	}}
	return nil
}
