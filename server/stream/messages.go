package stream

import "github.com/justintoman/npc-surprise/db"

const AdminPlayerId = 0

type AssignActionMessage struct {
	Type string    `json:"type" validate:"required,eq=assign-action"`
	Data db.Action `json:"action" validate:"required"`
}

func (stream *EventStream) SendAssignActionMessage(action db.Action) {
	payload := AssignActionMessage{
		Type: "assign-action",
		Data: action,
	}
	stream.sendMessage(action.PlayerId, payload)
	stream.sendAdminMessage(payload)
}

type AssignCharacterMessage struct {
	Type string                  `json:"type" validate:"required,eq=assign-character"`
	Data db.CharacterWithActions `json:"character" validate:"required"`
}

func (stream *EventStream) SendAssignCharacterMessage(character db.CharacterWithActions) {
	payload := AssignCharacterMessage{
		Type: "assign-character",
		Data: character,
	}
	stream.sendMessage(character.PlayerId, payload)
	stream.sendAdminMessage(payload)
}

type UnassignMessage struct {
	Type string `json:"type" validate:"required,oneof=unassign-action unassign-character"`
	Data int    `json:"data" validate:"required"` // action id or character id
}

func (stream *EventStream) SendUnassignActionMessage(playerId int, actionId int) {
	payload := UnassignMessage{
		Type: "unassign-action",
		Data: actionId,
	}
	stream.sendAdminMessage(payload)
	stream.sendMessage(playerId, payload)
}

func (stream *EventStream) SendUnassignCharacterMessage(playerId int, characterId int) {
	payload := UnassignMessage{
		Type: "unassign-character",
		Data: characterId,
	}
	stream.sendAdminMessage(payload)
	stream.sendMessage(playerId, payload)
}

type CharacterMessage struct {
	Type string                  `json:"type" validate:"required,eq=character"`
	Data db.CharacterWithActions `json:"data" validate:"required"`
}

func (stream *EventStream) SendAdminCharacterMessage(character db.CharacterWithActions) {
	payload := CharacterMessage{
		Type: "character",
		Data: character,
	}
	stream.sendAdminMessage(payload)
}
func (stream *EventStream) SendPlayerCharacterMessage(character db.CharacterWithActions) {
	if character.PlayerId == 0 {
		// not assigned, no one to send it to
		return
	}
	payload := CharacterMessage{
		Type: "character",
		Data: character,
	}
	stream.sendMessage(character.PlayerId, payload)
}

type ActionMessage struct {
	Type string    `json:"type" validate:"required,eq=action"`
	Data db.Action `json:"data" validate:"required"`
}

func (stream *EventStream) SendActionMessage(action db.Action) {
	payload := ActionMessage{
		Type: "action",
		Data: action,
	}
	stream.sendAdminMessage(payload)
	if action.PlayerId != 0 {
		stream.sendMessage(action.PlayerId, payload)
	}
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

func (stream *EventStream) SendInitMessage(playerId int, characters []db.CharacterWithActions) {
	payload := InitMessage{
		Type: "init",
		Data: InitMessageData{
			Players:    make([]PlayerWithStatus, 0),
			Characters: characters,
		},
	}
	stream.sendMessage(playerId, payload)
}

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
	payload := InitMessage{
		Type: "init",
		Data: InitMessageData{
			Players:    playersWithStatus,
			Characters: characters,
		},
	}
	stream.sendAdminMessage(payload)
}

type PlayerConnectedMessage struct {
	Type string    `json:"type" validate:"required,eq=player-connected"`
	Data db.Player `json:"data" validate:"required"`
}

func (stream *EventStream) SendPlayerConnectedMessage(player db.Player) {
	payload := PlayerConnectedMessage{
		Type: "player-connected",
		Data: player,
	}
	stream.sendAdminMessage(payload)
}

type PlayerDisconnectedMessage struct {
	Type string `json:"type" validate:"required,eq=player-disconnected"`
	Data int    `json:"data" validate:"required"` // player id
}

func (stream *EventStream) SendPlayerDisconnectedMessage(playerId int) {
	payload := PlayerDisconnectedMessage{
		Type: "player-disconnected",
		Data: playerId,
	}
	stream.sendAdminMessage(payload)
}

type DeleteMessage struct {
	Type string `json:"type" validate:"required,oneof=delete-action delete-character delete-player"`
	Data int    `json:"data" validate:"required"` // player, character, action id
}

func (stream *EventStream) SendDeleteActionMessage(id int) {
	payload := DeleteMessage{
		Type: "delete-action",
		Data: id,
	}
	stream.sendAdminMessage(payload)
}

func (stream *EventStream) SendDeleteCharacterMessage(id int) {
	payload := DeleteMessage{
		Type: "delete-character",
		Data: id,
	}
	stream.sendAdminMessage(payload)
}

func (stream *EventStream) SendDeletePlayerMessage(id int) {
	payload := DeleteMessage{
		Type: "delete-player",
		Data: id,
	}
	stream.sendAdminMessage(payload)
}
