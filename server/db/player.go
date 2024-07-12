package db

import (
	"encoding/json"
	"fmt"

	"github.com/supabase-community/supabase-go"
)

type Player struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type PlayerTable struct {
	client *supabase.Client
}

func (db PlayerTable) GetAll() (*[]Player, error) {
	data, _, err := db.client.From("players").Select("*", "exact", false).Execute()
	var players []Player
	json.Unmarshal(data, &players)
	return &players, err
}

func (db PlayerTable) Get(id string) (Player, error) {
	data, _, err := db.client.From("players").Select("*", "exact", false).Filter("id", "eq", id).Single().Execute()
	var player Player
	json.Unmarshal(data, &player)
	return player, err
}

func (db PlayerTable) Create(name string) (Player, error) {
	data, _, err := db.client.From("players").Insert(map[string]any{"name": name}, true, "", "", "exact").Single().Execute()
	var result Player
	json.Unmarshal(data, &result)
	fmt.Println("create player", result)
	return result, err
}

func (db PlayerTable) Update(player Player) (Player, error) {
	data, _, err := db.client.From("players").Insert(player, true, "", "", "exact").Single().Execute()
	var result Player
	var mapthing map[string]any
	json.Unmarshal(data, &result)
	json.Unmarshal(data, &mapthing)
	fmt.Println("update player", result)
	fmt.Println("update mapthing", mapthing)
	return result, err
}

func (db PlayerTable) Delete(id string) error {
	_, _, err := db.client.From("players").Delete("", "").Filter("id", "eq", id).Execute()
	return err
}

func (db PlayerTable) GetActions(id string) ([]Action, error) {
	data, _, err := db.client.From("actions").Select("*", "exact", false).Filter("player_id", "eq", id).Execute()
	var actions []Action
	json.Unmarshal(data, &actions)
	return actions, err
}

func (db PlayerTable) GetCharacters(id string) ([]Character, error) {
	data, _, err := db.client.From("characters").Select("*", "exact", false).Filter("player_id", "eq", id).Execute()
	var characters []Character
	json.Unmarshal(data, &characters)
	return characters, err
}
