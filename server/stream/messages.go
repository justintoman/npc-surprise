package stream

import (
	"log/slog"

	"github.com/justintoman/npc-surprise/db"
)

const AdminPlayerId = 0

type CharacterMessage struct {
	Type string                  `json:"type" validate:"required,eq=character"`
	Data db.CharacterWithActions `json:"data" validate:"required"`
}

type ActionMessage struct {
	Type string    `json:"type" validate:"required,eq=action"`
	Data db.Action `json:"data" validate:"required"`
}

type DeleteMessage struct {
	Type string `json:"type" validate:"required,oneof=delete-action delete-character delete-player"`
	Data int    `json:"data" validate:"required"` // player, character, action id
}

type InitMessage struct {
	Type string          `json:"type" validate:"required,eq=init"`
	Data InitMessageData `json:"data" validate:"required"`
}

type InitMessageData struct {
	Players    []PlayerWithStatus        `json:"players" validate:"required"`
	Characters []db.CharacterWithActions `json:"characters" validate:"required"`
}

type PlayerWithStatus struct {
	Id       int    `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required"`
	IsOnline bool   `json:"isOnline" validate:"required"`
}

type PlayerConnectedMessage struct {
	Type string    `json:"type" validate:"required,eq=player-connected"`
	Data db.Player `json:"data" validate:"required"`
}

type PlayerDisconnectedMessage struct {
	Type string `json:"type" validate:"required,eq=player-disconnected"`
	Data int    `json:"data" validate:"required"` // player id
}

/****************************************
*********** Player Messages *************
*****************************************/

func (stream *EventStream) SendInitPlayerMessage(playerId int, characters []db.CharacterWithActions) {
	payload := InitMessage{
		Type: "init",
		Data: InitMessageData{
			Players:    make([]PlayerWithStatus, 0),
			Characters: characters,
		},
	}
	stream.sendMessage(playerId, payload)
}

func (stream *EventStream) SendPlayerCharacterMessage(character db.CharacterWithActions) {
	if character.PlayerId == nil {
		slog.Error("unabled to reveal character because character not assigned to a player", "characterId", character.Id)
		return
	}
	stream.sendMessage(*character.PlayerId, CharacterMessage{
		Type: "character",
		Data: character,
	})
}

// since actions are either completely revealed or hidden, send to both admin and player
func (stream *EventStream) SendPlayerActionMessage(playerId int, action db.Action) {
	payload := ActionMessage{
		Type: "action",
		Data: action,
	}
	stream.sendAdminMessage(payload)
	stream.sendMessage(playerId, ActionMessage{
		Type: "action",
		Data: action,
	})
}

func (stream *EventStream) SendHideActionMessage(playerId int, action db.Action) {
	stream.sendAdminMessage(ActionMessage{
		Type: "action",
		Data: action,
	})

	// from the player's perspective, the action was deleted
	stream.sendMessage(playerId, DeleteMessage{
		Type: "delete-action",
		Data: action.Id,
	})
}

func (stream *EventStream) SendHideCharacterMessage(playerId int, character db.CharacterWithActions) {
	stream.sendAdminMessage(CharacterMessage{
		Type: "character",
		Data: character,
	})

	// from the player's perspective, the character was deleted
	stream.sendMessage(playerId, DeleteMessage{
		Type: "delete-character",
		Data: character.Id,
	})
}

/****************************************
*********** Admin Messages *************
*****************************************/

func (stream *EventStream) SendInitAdminMessage(players []db.Player, characters []db.CharacterWithActions) {
	connectedPlayers := stream.GetClients()
	playersWithStatus := make([]PlayerWithStatus, 0)

	for _, player := range players {
		if player.Id == 0 {
			continue
		}
		p := PlayerWithStatus{
			Id:       player.Id,
			Name:     player.Name,
			IsOnline: false,
		}
		for _, connectedPlayer := range connectedPlayers {
			if connectedPlayer.Id == p.Id {
				p.IsOnline = true
				break
			}
		}
		playersWithStatus = append(playersWithStatus, p)
	}

	stream.sendAdminMessage(InitMessage{
		Type: "init",
		Data: InitMessageData{
			Players:    playersWithStatus,
			Characters: characters,
		},
	})
}

func (stream *EventStream) SendAdminCharacterMessage(character db.CharacterWithActions) {
	stream.sendAdminMessage(CharacterMessage{
		Type: "character",
		Data: character,
	})
}

func (stream *EventStream) SendAdminActionMessage(action db.Action) {
	stream.sendAdminMessage(ActionMessage{
		Type: "action",
		Data: action,
	})
}

func (stream *EventStream) SendPlayerConnectedMessage(player db.Player) {
	stream.sendAdminMessage(PlayerConnectedMessage{
		Type: "player-connected",
		Data: player,
	})
}

func (stream *EventStream) SendPlayerDisconnectedMessage(playerId int) {
	stream.sendAdminMessage(PlayerDisconnectedMessage{
		Type: "player-disconnected",
		Data: playerId,
	})
}

func (stream *EventStream) SendDeleteCharacterMessage(id int) {
	stream.sendAdminMessage(DeleteMessage{
		Type: "delete-character",
		Data: id,
	})
}

func (stream *EventStream) SendDeleteActionMessage(id int) {
	stream.sendAdminMessage(DeleteMessage{
		Type: "delete-action",
		Data: id,
	})
}

func (stream *EventStream) SendDeletePlayerMessage(id int) {
	stream.sendAdminMessage(DeleteMessage{
		Type: "delete-player",
		Data: id,
	})
}
