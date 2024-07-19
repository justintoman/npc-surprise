package services

import (
	"fmt"
	"log/slog"

	"github.com/justintoman/npc-surprise/db"
	"github.com/justintoman/npc-surprise/stream"
)

type CharacterService struct {
	db     db.Db
	stream stream.StreamingServer
}

func NewCharacterService(db db.Db, stream stream.StreamingServer) CharacterService {
	return CharacterService{
		db:     db,
		stream: stream,
	}
}

func (s *CharacterService) Create(input db.CreateCharacterPayload) (db.CharacterWithActions, error) {
	character, err := s.db.Character.Create(input)
	if err != nil {
		slog.Error("Error creating character", "error", err)
		return db.CharacterWithActions{}, err
	}
	data := db.CharacterWithActions{
		Character: character,
		Actions:   make([]db.Action, 0),
	}
	return data, nil
}

func (s *CharacterService) Update(input db.Character) (db.CharacterWithActions, error) {
	character, err := s.db.Character.Update(input)
	if err != nil {
		slog.Error("Error updating character", "error", err)
		return db.CharacterWithActions{}, err
	}
	actions, err := s.db.Action.GetAll(character.Id)
	if err != nil {
		slog.Error("Error fetching actions after updating character", "error", err)
		return db.CharacterWithActions{}, err
	}
	data := db.CharacterWithActions{
		Character: character,
		Actions:   actions,
	}
	return data, nil
}

func (s *CharacterService) GetAllWithActions() ([]db.CharacterWithActions, error) {
	characters, err := s.db.Character.GetAll()
	if err != nil {
		slog.Error("Error fetching characters", "error", err)
		return []db.CharacterWithActions{}, err
	}
	charsWithActions := make([]db.CharacterWithActions, len(characters))
	for i, character := range characters {
		actions, err := s.db.Action.GetAll(character.Id)
		if err != nil {
			slog.Error("error getting actions for character", "error", err, "characterId", character.Id)
			return []db.CharacterWithActions{}, err
		}
		charsWithActions[i] = db.CharacterWithActions{
			Character: character,
			Actions:   actions,
		}
	}
	return charsWithActions, nil
}

func (s *CharacterService) GetAllAssignedWithActionsRedacted(playerId int) ([]db.CharacterWithActions, error) {
	characters, err := s.db.Character.GetAllByPlayerId(playerId)
	if err != nil {
		slog.Error("Error fetching characters", "error", err)
		return []db.CharacterWithActions{}, err
	}
	charsWithActions := make([]db.CharacterWithActions, len(characters))
	for i, character := range characters {
		actions, err := s.db.Action.GetAll(character.Id)
		if err != nil {
			slog.Error("error getting actions for character", "error", err, "characterId", character.Id)
			return []db.CharacterWithActions{}, err
		}
		fields, err := s.db.Character.GetRevealedFields(character.Id)
		if err != nil {
			slog.Error("error getting revealed fields for character", "error", err, "characterId", character.Id)
			return []db.CharacterWithActions{}, err
		}
		redactCharacter(&character, fields)
		charsWithActions[i] = db.CharacterWithActions{
			Character: character,
			Actions:   actions,
		}
	}
	return charsWithActions, nil
}

func (s *CharacterService) Redact(character db.CharacterWithActions) (db.CharacterWithActions, error) {
	fields, err := s.db.Character.GetRevealedFields(character.Id)
	if err != nil {
		slog.Error("error getting revealed fields for character", "error", err, "characterId", character.Id)
		return db.CharacterWithActions{}, err
	}
	actions := make([]db.Action, 0)
	for _, action := range character.Actions {
		if action.CharacterId != character.Id {
			continue
		}
		actions = append(actions, action)
	}
	redactCharacter(&character.Character, fields)
	redacted := db.CharacterWithActions{
		Character: character.Character,
		Actions:   actions,
	}
	return redacted, nil
}

func redactCharacter(character *db.Character, revealedFields db.CharacterReveleadFields) {
	if !revealedFields.Age {
		character.Age = ""
	}
	if !revealedFields.Gender {
		character.Gender = ""
	}
	if !revealedFields.Race {
		character.Race = ""
	}
	if !revealedFields.Name {
		character.Name = ""
	}
	if !revealedFields.Description {
		character.Description = ""
	}
	if !revealedFields.Appearance {
		character.Appearance = ""
	}
}

func (s *CharacterService) Assign(characterId int, playerId int) (db.CharacterWithActions, error) {
	character, err := s.db.Character.Get(characterId)
	if err != nil {
		slog.Error("error getting character to assign", "error", err, "characterId", characterId)
		return db.CharacterWithActions{}, err
	}
	if character.PlayerId == playerId {
		slog.Error("already assigned to target player", "characterId", characterId, "playerId", playerId)
		return db.CharacterWithActions{}, fmt.Errorf("character already assigned to target player")
	}

	// if there any any actions from this character that are assigned, unassign them
	actions, err := s.db.Action.GetAll(characterId)
	if err != nil {
		return db.CharacterWithActions{}, err
	}
	for i, action := range actions {
		if action.PlayerId != 0 {
			prevPlayerId := action.PlayerId
			action.PlayerId = 0
			actions[i], err = s.db.Action.Update(action)
			if err != nil {
				slog.Error("error unassigning action from previous player", "error", err)
				return db.CharacterWithActions{}, err
			}

			s.stream.SendUnassignActionMessage(prevPlayerId, action.Id)
		}
	}

	prevPlayerId := character.PlayerId
	character.PlayerId = playerId
	character, err = s.db.Character.Update(character)
	if err != nil {
		slog.Error("error updating character", "error", err)
		return db.CharacterWithActions{}, err
	}
	if prevPlayerId != character.PlayerId {
		// if this was assigned from a previous player, send an unassign message to them
		s.stream.SendUnassignCharacterMessage(character.PlayerId, character.Id)
	}

	withActions := db.CharacterWithActions{
		Character: character,
		Actions:   actions,
	}
	return withActions, nil
}

func (s *CharacterService) Unassign(characterId int) (int, error) {
	character, err := s.db.Character.Get(characterId)
	if err != nil {
		slog.Error("error getting character to unassign", "error", err, "characterId", characterId)
		return 0, err
	}
	if character.PlayerId == 0 {
		slog.Error("character not assigned to a player", "characterId", characterId)
		return 0, fmt.Errorf("character not assigned to a player, can't unassign from nobody")
	}
	prevPlayerId := character.PlayerId
	character.PlayerId = 0
	character, err = s.db.Character.Update(character)
	if err != nil {
		slog.Error("error updating character", "error", err)
		return 0, err
	}
	actions, err := s.db.Action.GetAll(characterId)
	if err != nil {
		slog.Error("error getting actions for character", "error", err, "characterId", characterId)
		return 0, err
	}
	for _, action := range actions {
		if action.PlayerId == 0 {
			continue
		}
		prevPlayerId := action.PlayerId
		action.PlayerId = 0
		action, err = s.db.Action.Update(action)
		if err != nil {
			slog.Error("error unassigning action from previous player", "error", err)
			return 0, err
		}
		s.stream.SendUnassignActionMessage(prevPlayerId, action.Id)
	}
	return prevPlayerId, nil
}

func (s *CharacterService) UpdateRevealedFields(input db.CharacterReveleadFields) (db.CharacterWithActions, error) {
	fields, err := s.db.Character.UpdateRevealedFields(input)
	if err != nil {
		slog.Error("error getting character to update", "error", err, "characterId", input.CharacterId)
		return db.CharacterWithActions{}, err
	}

	character, err := s.db.Character.Get(fields.CharacterId)
	if err != nil {
		slog.Error("error getting character to update", "error", err, "characterId", input.CharacterId)
		return db.CharacterWithActions{}, err
	}

	actions, err := s.db.Action.GetAll(character.Id)
	if err != nil {
		slog.Error("error getting actions for character", "error", err, "characterId", character.Id)
		return db.CharacterWithActions{}, err
	}
	withActions := db.CharacterWithActions{
		Character: character,
		Actions:   actions,
	}
	return withActions, nil
}

func (s *CharacterService) Delete(id int) error {
	err := s.db.Character.Delete(id)
	if err != nil {
		slog.Error("error deleting character", "error", err)
		return err
	}
	return nil
}
