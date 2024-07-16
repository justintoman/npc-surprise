package db

import (
	"encoding/json"
	"strconv"

	"github.com/supabase-community/supabase-go"
)

type CreateActionPayload struct {
	ActionBase
}

type UpdateActionPayload struct {
	Action
}

type Action struct {
	ActionBase
	Id       int `json:"id"`
	PlayerId int `json:"player_id"`
}

type ActionBase struct {
	CharacterId int    `json:"character_id"`
	Type        string `json:"type,omitempty"`
	Content     string `json:"content,omitempty"`
	Direction   string `json:"direction,omitempty"`
}

type ActionTable struct {
	client *supabase.Client
}

func (db ActionTable) GetAll(characterId int) ([]Action, error) {
	data, _, err := db.client.From("actions").Select("*", "exact", false).Filter("character_id", "eq", strconv.Itoa(characterId)).Execute()
	var actions []Action
	json.Unmarshal(data, &actions)
	return actions, err
}

func (db ActionTable) GetAllAssigned(playerId string) ([]Action, error) {
	data, _, err := db.client.From("actions").Select("*", "exact", false).Filter("player_id", "eq", playerId).Execute()
	var actions []Action
	json.Unmarshal(data, &actions)
	return actions, err
}

func (db ActionTable) Get(id string) (Action, error) {
	data, _, err := db.client.From("actions").Select("*", "exact", false).Filter("id", "eq", id).Single().Execute()
	var action Action
	json.Unmarshal(data, &action)
	return action, err
}

func (db ActionTable) Create(action CreateActionPayload) (Action, error) {
	data, _, err := db.client.From("actions").Insert(action, true, "", "", "exact").Single().Execute()
	var result Action
	json.Unmarshal(data, &result)
	return result, err
}

func (db ActionTable) Update(action UpdateActionPayload) (Action, error) {
	data, _, err := db.client.From("actions").Insert(action, true, "", "", "exact").Execute()
	var result Action
	json.Unmarshal(data, &result)
	return result, err
}

func (db ActionTable) Delete(id string) error {
	_, _, err := db.client.From("actions").Delete("", "").Filter("id", "eq", id).Execute()
	return err
}
