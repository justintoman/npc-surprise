package services

import (
	"log/slog"

	"github.com/justintoman/npc-surprise/db"
)

type ActionService struct {
	db db.Db
}

func NewActionService(db db.Db) ActionService {
	return ActionService{
		db: db,
	}
}

func (s *ActionService) Create(input db.CreateActionPayload) (db.Action, error) {
	action, err := s.db.Action.Create(input)
	if err != nil {
		slog.Error("Error creating action", "error", err)
		return db.Action{}, err
	}
	return action, nil
}

func (s *ActionService) Update(input db.Action) (db.Action, error) {
	action, err := s.db.Action.Update(input)
	if err != nil {
		slog.Error("Error updating action", "error", err)
		return db.Action{}, err
	}
	return action, nil
}

func (s *ActionService) Reveal(actionId int) (int, db.Action, error) {
	action, err := s.db.Action.Get(actionId)
	if err != nil {
		slog.Error("Error getting action to reveal", "error", err)
		return 0, db.Action{}, err
	}
	character, err := s.db.Character.Get(action.CharacterId)
	if err != nil {
		slog.Error("Error getting character for action to reveal", "error", err)
		return 0, db.Action{}, err
	}
	if character.PlayerId == nil {
		slog.Error("character not assigned to a player, cannot reveal action", "characterId", action.CharacterId)
		return 0, db.Action{}, nil
	}
	if action.Revealed {
		slog.Info("already revealed", "actionId", actionId)
		return 0, db.Action{}, nil
	}
	action.Revealed = true
	action, err = s.db.Action.Update(action)
	if err != nil {
		slog.Error("Error revealing action", "error", err)
		return 0, db.Action{}, err
	}
	return *character.PlayerId, action, nil
}

func (s *ActionService) Hide(actionId int) (int, db.Action, error) {
	action, err := s.db.Action.Get(actionId)
	if err != nil {
		slog.Error("Error getting action", "error", err)
		return 0, db.Action{}, err
	}
	if !action.Revealed {
		slog.Info("already hidden", "actionId", actionId)
		return 0, db.Action{}, nil
	}
	character, err := s.db.Character.Get(action.CharacterId)
	if err != nil {
		slog.Error("Error getting character for action to hide", "error", err)
		return 0, db.Action{}, err
	}
	action.Revealed = false
	action, err = s.db.Action.Update(action)
	if err != nil {
		slog.Error("Error unassigning action", "error", err)
		return 0, db.Action{}, err
	}
	return *character.PlayerId, action, nil
}

func (s *ActionService) Delete(id int) error {
	err := s.db.Action.Delete(id)
	if err != nil {
		slog.Error("Error deleting action", "error", err)
		return err
	}
	return nil
}
