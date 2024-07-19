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

func (s *ActionService) Assign(actionId int) (db.Action, error) {
	action, err := s.db.Action.Get(actionId)
	if err != nil {
		return db.Action{}, err
	}
	if action.PlayerId != 0 {
		slog.Info("already assigned", "actionId", actionId, "current playerId", action.PlayerId)
		return db.Action{}, nil
	}
	character, err := s.db.Character.Get(action.CharacterId)
	if err != nil {
		slog.Error("Error getting character", "error", err)
		return db.Action{}, err
	}
	if character.PlayerId == 0 {
		slog.Error("character not assigned to a player", "characterId", action.CharacterId)
		return db.Action{}, nil
	}
	action.PlayerId = character.PlayerId
	action, err = s.db.Action.Update(action)
	if err != nil {
		slog.Error("Error assigning action", "error", err)
		return db.Action{}, err
	}
	return action, nil
}

func (s *ActionService) Unassign(actionId int) (int, error) {
	action, err := s.db.Action.Get(actionId)
	if err != nil {
		slog.Error("Error getting action", "error", err)
		return 0, err
	}
	if action.PlayerId == 0 {
		slog.Info("already unassigned", "actionId", actionId)
		return 0, nil
	}
	oldPlayerId := action.PlayerId
	action.PlayerId = 0
	action, err = s.db.Action.Update(action)
	if err != nil {
		slog.Error("Error unassigning action", "error", err)
		return 0, err
	}
	return oldPlayerId, nil
}

func (s *ActionService) Delete(id int) error {
	err := s.db.Action.Delete(id)
	if err != nil {
		slog.Error("Error deleting action", "error", err)
		return err
	}
	return nil
}
