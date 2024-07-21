package stream

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/justintoman/npc-surprise/pkg/db"
)

// This file is heavily inspired by
// https://github.com/gin-gonic/examples/blob/master/server-sent-event/main.go

type StreamingServer interface {
	NewUserStream(onAdd OnClientAddedFunc, onRemove OnClientRemovedFunc) (gin.HandlerFunc, StreamHandlerFunc)
	Close(ClientChan)
	Listen(context.Context)
	GetClients() []db.Player

	// player messages

	SendInitPlayerMessage(playerId int, characters []db.CharacterWithActions)

	// Send a redacted character to only the assigned player.
	// Note that admins need a full non-redacted character, so this only sends to the player.
	SendPlayerCharacterMessage(charcter db.CharacterWithActions)
	SendPlayerActionMessage(playerId int, action db.Action)
	SendHideActionMessage(playerId int, action db.Action)
	SendHideCharacterMessage(playerId int, character db.CharacterWithActions)

	// admin messages
	SendInitAdminMessage(players []db.Player, characters []db.CharacterWithActions, fields []db.CharacterReveleadFields)
	SendAdminCharacterMessage(character db.CharacterWithActions)
	SendAdminCharacterMessageWithFields(character db.CharacterWithActions, fields db.CharacterReveleadFields)
	SendAdminActionMessage(action db.Action)
	SendPlayerConnectedMessage(player db.Player)
	SendPlayerDisconnectedMessage(playerId int)
	SendDeleteCharacterMessage(characterId int)
	SendDeleteActionMessage(actionId int)
	SendDeletePlayerMessage(playerId int)
}

func New(db db.Db) StreamingServer {
	eventStream := &EventStream{
		Message:       make(chan Message),
		NewClients:    make(chan ClientChan),
		ClosedClients: make(chan ClientChan),
		TotalClients:  make(map[ClientChan]bool),
	}
	return eventStream
}

func (stream *EventStream) Close(clientChan ClientChan) {
	stream.ClosedClients <- clientChan
}

func (stream *EventStream) sendMessage(playerId int, message any) {
	slog.Info("Sending message", "clientId", playerId)
	stream.Message <- Message{PlayerId: playerId, Payload: message}
}

func (stream *EventStream) sendAdminMessage(message any) {
	for client := range stream.TotalClients {
		if client.Id == AdminPlayerId {
			stream.sendMessage(client.Id, message)
		}
	}

	stream.sendMessage(AdminPlayerId, message)
}

func (stream *EventStream) GetClients() []db.Player {
	clients := make([]db.Player, 0, len(stream.TotalClients))
	for client := range stream.TotalClients {
		clients = append(clients, client.Player)
	}
	return clients
}

type ClientChan struct {
	db.Player
	Channel chan any
}

type Message struct {
	PlayerId int
	Payload  any
}

type EventStream struct {
	Db            db.Db
	Message       chan Message
	NewClients    chan ClientChan
	ClosedClients chan ClientChan
	// map of clients by string topic
	TotalClients map[ClientChan]bool
}

func (stream *EventStream) Listen(c context.Context) {
	for {
		select {
		// Add new available client
		case client := <-stream.NewClients:
			stream.TotalClients[client] = true
			names := make([]string, 0)
			for client := range stream.TotalClients {
				names = append(names, client.Name)
			}
			slog.Info(fmt.Sprintf("Client added. %d registered clients. Clients: %+v", len(stream.TotalClients), names))

		// Remove closed client
		case client := <-stream.ClosedClients:
			slog.Info("Client closed", "id", client.Id, "name", client.Name)
			delete(stream.TotalClients, client)
			close(client.Channel)
			slog.Info(fmt.Sprintf("Removed client. %d registered clients", len(stream.TotalClients)))

		// Broadcast message to client
		case eventMsg := <-stream.Message:
			sentMessage := false
			for client := range stream.TotalClients {
				if client.Id == eventMsg.PlayerId {
					sentMessage = true
					client.Channel <- eventMsg.Payload
				}
			}
			if !sentMessage {
				slog.Error(fmt.Sprintf("Attempted to send message to a client that doesn't exist. Id: %d", eventMsg.PlayerId))
				continue
			}

		case <-c.Done():
			return
		}
	}
}

type StreamHandlerFunc func(*gin.Context) (bool, error)

type OnClientAddedFunc func(db.Player)
type OnClientRemovedFunc func(db.Player)

func (stream *EventStream) NewUserStream(onAdded OnClientAddedFunc, onRemoved OnClientRemovedFunc) (gin.HandlerFunc, StreamHandlerFunc) {
	middleware := func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")

		ctxPlayer, ok := c.Get("player")
		if !ok {
			slog.Info("Player not found. Is Middleware applied?")
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid Player. Try logging in again.", "status": 401})
			return
		}

		player := ctxPlayer.(db.Player)

		slog.Info("Player connected, creating client", "id", player.Id, "name", player.Name)

		clientChan := ClientChan{
			Channel: make(chan any),
			Player: db.Player{
				Id:   player.Id,
				Name: player.Name,
			},
		}
		stream.NewClients <- clientChan

		go onAdded(player)

		go func() {
			<-c.Writer.CloseNotify()
			stream.Close(clientChan)
			go onRemoved(player)
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
