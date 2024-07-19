package db

import (
	"encoding/json"

	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

type CreateActionPayload struct {
	Content     string `json:"content" binding:"required"`
	Order       int    `json:"order" binding:"required"`
	CharacterId int    `json:"characterId" binding:"required"`
}

type Action struct {
	Id          int    `json:"id" binding:"required"`
	Content     string `json:"content" binding:"required"`
	Order       int    `json:"order" binding:"required"`
	CharacterId int    `json:"characterId" binding:"required"`
	PlayerId    int    `json:"playerId"`
}

type ActionTable struct {
	client *supabase.Client
}

func (db ActionTable) GetAll(characterId int) ([]Action, error) {
	query := selectAll(db.from())
	query = filterByCharacterId(query, characterId)
	query = orderActions(query)
	data, _, err := query.Execute()
	var actions []Action
	json.Unmarshal(data, &actions)
	return actions, err
}

func (db ActionTable) GetAllByPlayerId(characterId, playerId int) ([]Action, error) {
	query := selectAll(db.from())
	query = filterByCharacterId(query, characterId)
	query = filterByPlayerId(query, playerId)
	query = orderActions(query)
	data, _, err := query.Execute()
	var actions []Action
	json.Unmarshal(data, &actions)
	return actions, err
}

func (db ActionTable) Get(id int) (Action, error) {
	query := selectAll(db.from())
	query = filterById(query, id)
	data, _, err := query.Execute()
	var action Action
	json.Unmarshal(data, &action)
	return action, err
}

func (db ActionTable) Create(action CreateActionPayload) (Action, error) {
	query := insertSingle(db.from(), action)
	data, _, err := query.Execute()
	var result Action
	json.Unmarshal(data, &result)
	return result, err
}

func (db ActionTable) Update(action Action) (Action, error) {
	query := insertSingle(db.from(), action)
	data, _, err := query.Execute()
	var result Action
	json.Unmarshal(data, &result)
	return result, err
}

func (db ActionTable) Delete(id int) error {
	query := deleteSingle(db.from())
	query = filterById(query, id)
	_, _, err := query.Execute()
	return err
}

func (db ActionTable) from() *postgrest.QueryBuilder {
	return db.client.From("actions")
}

func orderActions(query *postgrest.FilterBuilder) *postgrest.FilterBuilder {
	return query.Order("order", &postgrest.OrderOpts{
		Ascending: true,
	})
}
