package services

import (
	"fmt"
	"log/slog"

	"github.com/justintoman/npc-surprise/pkg/db"
	"github.com/justintoman/npc-surprise/pkg/stream"
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

func (s *CharacterService) Create(input db.CreateCharacterPayload) (db.CharacterWithActions, db.CharacterReveleadFields, error) {
	character, fields, err := s.db.Character.Create(input)
	if err != nil {
		slog.Error("Error creating character", "error", err)
		return db.CharacterWithActions{}, db.CharacterReveleadFields{}, err
	}
	data := db.CharacterWithActions{
		Character: character,
		Actions:   make([]db.Action, 0),
	}
	return data, fields, nil
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

func (s *CharacterService) GetAllWithActionsAndFields() ([]db.CharacterWithActions, []db.CharacterReveleadFields, error) {
	characters, err := s.db.Character.GetAll()
	if err != nil {
		slog.Error("Error fetching characters", "error", err)
		return []db.CharacterWithActions{}, []db.CharacterReveleadFields{}, err
	}
	fields := make([]db.CharacterReveleadFields, len(characters))
	charsWithActions := make([]db.CharacterWithActions, len(characters))
	for i, character := range characters {
		actions, err := s.db.Action.GetAll(character.Id)
		if err != nil {
			slog.Error("error getting actions for character", "error", err, "characterId", character.Id)
			return []db.CharacterWithActions{}, []db.CharacterReveleadFields{}, err
		}
		charsWithActions[i] = db.CharacterWithActions{
			Character: character,
			Actions:   actions,
		}
		fields[i], err = s.db.Character.GetRevealedFields(character.Id)
		if err != nil {
			slog.Error("error getting revealed fields for character", "error", err, "characterId", character.Id)
			return []db.CharacterWithActions{}, []db.CharacterReveleadFields{}, err
		}
	}
	return charsWithActions, fields, nil
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
		if !action.Revealed {
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

func (s *CharacterService) Assign(characterId int, playerId int) (*int, db.CharacterWithActions, error) {
	character, err := s.db.Character.Get(characterId)
	if err != nil {
		slog.Error("error getting character to assign", "error", err, "characterId", characterId)
		return nil, db.CharacterWithActions{}, err
	}

	actions, err := s.db.Action.GetAll(characterId)
	if err != nil {
		return nil, db.CharacterWithActions{}, err
	}

	prevPlayerId := character.PlayerId
	character.PlayerId = &playerId
	character, err = s.db.Character.Update(character)
	if err != nil {
		slog.Error("error updating character", "error", err)
		return nil, db.CharacterWithActions{}, err
	}

	withActions := db.CharacterWithActions{
		Character: character,
		Actions:   actions,
	}
	return prevPlayerId, withActions, nil
}

func (s *CharacterService) Unassign(characterId int) (int, db.CharacterWithActions, error) {
	character, err := s.db.Character.Get(characterId)
	if err != nil {
		slog.Error("error getting character to unassign", "error", err, "characterId", characterId)
		return 0, db.CharacterWithActions{}, err
	}
	if character.PlayerId == nil {
		slog.Error("character not assigned to a player", "characterId", characterId)
		return 0, db.CharacterWithActions{}, fmt.Errorf("character not assigned to a player, can't unassign from nobody")
	}
	prevPlayerId := *character.PlayerId
	character.PlayerId = nil
	character, err = s.db.Character.Update(character)
	if err != nil {
		slog.Error("error updating character", "error", err)
		return 0, db.CharacterWithActions{}, err
	}

	// hide any currently revealed actions
	actions, err := s.db.Action.GetAll(characterId)
	if err != nil {
		return 0, db.CharacterWithActions{}, err
	}
	for i, action := range actions {
		if action.Revealed {
			action.Revealed = false
			actions[i], err = s.db.Action.Update(action)
			if err != nil {
				slog.Error("error unassigning action from previous player", "error", err)
				return 0, db.CharacterWithActions{}, err
			}
		}
	}

	withActions := db.CharacterWithActions{
		Character: character,
		Actions:   actions,
	}

	return prevPlayerId, withActions, nil
}

func (s *CharacterService) UpdateRevealedFields(input db.CharacterReveleadFields) (db.CharacterWithActions, db.CharacterReveleadFields, error) {
	fields, err := s.db.Character.UpdateRevealedFields(input)
	if err != nil {
		slog.Error("error getting character to update", "error", err, "characterId", input.CharacterId)
		return db.CharacterWithActions{}, db.CharacterReveleadFields{}, err
	}

	character, err := s.db.Character.Get(fields.CharacterId)
	if err != nil {
		slog.Error("error getting character to update", "error", err, "characterId", input.CharacterId)
		return db.CharacterWithActions{}, db.CharacterReveleadFields{}, err
	}

	actions, err := s.db.Action.GetAll(character.Id)
	if err != nil {
		slog.Error("error getting actions for character", "error", err, "characterId", character.Id)
		return db.CharacterWithActions{}, db.CharacterReveleadFields{}, err
	}
	withActions := db.CharacterWithActions{
		Character: character,
		Actions:   actions,
	}
	return withActions, fields, nil
}

func (s *CharacterService) Delete(id int) error {
	err := s.db.Character.Delete(id)
	if err != nil {
		slog.Error("error deleting character", "error", err)
		return err
	}
	return nil
}
