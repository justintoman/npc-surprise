package db

import (
	"encoding/json"
	"fmt"

	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

type CreatePlayerPayload struct {
	Name string `json:"name"`
}

type Player struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type PlayerTable struct {
	client *supabase.Client
}

func (db PlayerTable) GetAll() ([]Player, error) {
	query := selectAll(db.from())
	data, _, err := query.Execute()
	var players []Player
	json.Unmarshal(data, &players)
	return players, err
}

func (db PlayerTable) Get(id int) (Player, error) {
	query := selectAll(db.from())
	query = filterById(query, id)
	data, _, err := query.Execute()
	var player Player
	json.Unmarshal(data, &player)
	return player, err
}

func (db PlayerTable) Create(payload CreatePlayerPayload) (Player, error) {
	query := insertSingle(db.from(), payload).Single()
	data, _, err := query.Execute()
	var result Player
	json.Unmarshal(data, &result)
	fmt.Println("create player", result)
	return result, err
}

func (db PlayerTable) Update(player Player) (Player, error) {
	query := insertSingle(db.from(), player).Single()
	data, _, err := query.Execute()
	var result Player
	json.Unmarshal(data, &result)
	fmt.Println("update player", result)
	return result, err
}

func (db PlayerTable) Delete(id int) error {
	query := deleteSingle(db.from())
	query = filterById(query, id)
	_, _, err := query.Execute()
	return err
}

func (table PlayerTable) from() *postgrest.QueryBuilder {
	return table.client.From("players")
}
