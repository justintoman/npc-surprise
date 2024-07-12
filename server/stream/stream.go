package stream

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/gin-gonic/gin"
)

// This file is heavily inspired by
// https://github.com/gin-gonic/examples/blob/master/server-sent-event/main.go

type StreamingServer interface {
	NewUserStream() (gin.HandlerFunc, StreamHandlerFunc)
	Close(ClientChan)
	SendMessage(topic string, message MessagePayload)
	Listen(context.Context)
	GetClients() []UserClient
}

type UserClient struct {
	Id   string
	Name string
}

func New() StreamingServer {
	eventStream := &EventStream{
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

func (stream *EventStream) SendMessage(topic string, message MessagePayload) {
	stream.Message <- Message{PlayerId: topic, Payload: message}
}

func (stream *EventStream) GetClients() []UserClient {
	clients := make([]UserClient, 0, len(stream.TotalClients))
	for _, client := range stream.TotalClients {
		clients = append(clients, *client.UserClient)
	}
	return clients
}

type ClientChan struct {
	*UserClient
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
			stream.TotalClients[client.Id] = client
			if client.Id == "1" {
				stream.TotalClients["admin"] = client
			} else {
				// tell the admin there's a new user
				admin, ok := stream.TotalClients["admin"]
				if !ok {
					slog.Error("Could not find admin client")
					continue
				}
				admin.Channel <- MessagePayload{Type: "new_client", Data: client.UserClient}
			}

			slog.Info(fmt.Sprintf("Client added. %d registered clients", len(stream.TotalClients)))

		// Remove closed client
		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client.Id)
			close(client.Channel)

			// tell the admin a user left
			admin := stream.TotalClients["admin"]
			admin.Channel <- MessagePayload{Type: "closed_client", Data: client.UserClient}
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

		clientChan := ClientChan{Channel: make(chan MessagePayload), UserClient: &UserClient{Id: id, Name: name}}
		stream.NewClients <- clientChan

		defer stream.Close(clientChan)

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
