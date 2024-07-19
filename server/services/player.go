package services

import (
	"log/slog"

	"github.com/justintoman/npc-surprise/db"
	"github.com/justintoman/npc-surprise/stream"
)

type PlayerService struct {
	db     db.Db
	stream stream.StreamingServer
}

func NewPlayerService(db db.Db, stream stream.StreamingServer) PlayerService {
	return PlayerService{
		db:     db,
		stream: stream,
	}
}

func (s *PlayerService) Create(name string) (db.Player, error) {
	player, err := s.db.Player.Create(db.CreatePlayerPayload{Name: name})
	if err != nil {
		slog.Error("Error creating player", "error", err)
		return db.Player{}, err
	}
	return player, nil
}

func (s *PlayerService) Update(input db.Player) (db.Player, error) {
	player, err := s.db.Player.Update(input)
	if err != nil {
		slog.Error("Error updating player", "error", err)
		return db.Player{}, err
	}
	return player, nil
}

func (s *PlayerService) GetAll() ([]db.Player, error) {
	players, err := s.db.Player.GetAll()
	if err != nil {
		slog.Error("Error fetching players", "error", err)
		return []db.Player{}, err
	}
	return players, nil
}

func (s *PlayerService) Delete(id int) error {
	err := s.db.Player.Delete(id)
	if err != nil {
		slog.Error("Error fetching player", "error", err)
		return err
	}
	return nil
}
