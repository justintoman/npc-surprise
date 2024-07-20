package db

import (
	"encoding/json"
	"log/slog"

	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

type Character struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	PlayerId    *int   `json:"playerId"`
	Race        string `json:"race,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Age         string `json:"age,omitempty"`
	Description string `json:"description,omitempty"`
	Appearance  string `json:"appearance,omitempty"`
}

type CharacterWithActions struct {
	Character `json:",inline"`
	Actions   []Action `json:"actions"`
}

type CreateCharacterPayload struct {
	Name        string `json:"name" binding:"required"`
	Race        string `json:"race,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Age         string `json:"age,omitempty"`
	Description string `json:"description,omitempty"`
	Appearance  string `json:"appearance,omitempty"`
}

type CharacterReveleadFields struct {
	CharacterId int  `json:"characterId" binding:"required"`
	Name        bool `json:"name" binding:"required"`
	Race        bool `json:"race" binding:"required"`
	Gender      bool `json:"gender" binding:"required"`
	Age         bool `json:"age" binding:"required"`
	Description bool `json:"description" binding:"required"`
	Appearance  bool `json:"appearance" binding:"required"`
}

type CharacterTable struct {
	client *supabase.Client
}

func (db CharacterTable) GetAll() ([]Character, error) {
	query := selectAll(db.from())
	data, _, err := query.Execute()
	characters := make([]Character, 0)
	json.Unmarshal(data, &characters)
	return characters, err
}

func (db CharacterTable) GetAllByPlayerId(id int) ([]Character, error) {
	query := selectAll(db.from())
	query = filterByPlayerId(query, id)
	data, _, err := query.Execute()
	characters := make([]Character, 0)
	json.Unmarshal(data, &characters)
	return characters, err
}

func (db CharacterTable) Get(id int) (Character, error) {
	query := selectAll(db.from())
	query = filterById(query, id)
	data, _, err := query.Execute()
	var character Character
	json.Unmarshal(data, &character)
	return character, err
}

func (db CharacterTable) Create(character CreateCharacterPayload) (Character, error) {
	query := insertSingle(db.from(), character)
	data, _, err := query.Execute()
	if err != nil {
		slog.Error("Error creating character", "error", err)
		return Character{}, err
	}
	var result Character
	err = json.Unmarshal(data, &result)
	if err != nil {
		slog.Error("Error unmarshalling character", "error", err)
		return Character{}, err
	}

	fields := CharacterReveleadFields{CharacterId: result.Id}
	fieldsQuery := insertSingle(db.fromRevealedFields(), fields)
	_, _, err = fieldsQuery.Execute()
	if err != nil {
		slog.Error("Error creating character_revelead_fields", "error", err)
		return Character{}, err
	}

	return result, err
}

func (db CharacterTable) Update(character Character) (Character, error) {
	query := insertSingle(db.from(), character)
	_, _, err := query.Execute()
	if err != nil {
		slog.Error("Error updating character", "error", err)
		return Character{}, err
	}
	return db.Get(character.Id)
}

func (db CharacterTable) GetRevealedFields(characterId int) (CharacterReveleadFields, error) {
	query := selectAll(db.fromRevealedFields())
	query = filterByCharacterId(query, characterId).Single()
	data, _, err := query.Execute()
	var revealedFields CharacterReveleadFields
	json.Unmarshal(data, &revealedFields)
	return revealedFields, err
}

func (db CharacterTable) UpdateRevealedFields(revealedFields CharacterReveleadFields) (CharacterReveleadFields, error) {
	query := insertSingle(db.fromRevealedFields(), revealedFields)
	data, _, err := query.Execute()
	var result CharacterReveleadFields
	json.Unmarshal(data, &result)
	return result, err
}

func (db CharacterTable) Delete(id int) error {
	query := deleteSingle(db.from())
	query = filterById(query, id)
	_, _, err := query.Execute()
	return err
}

func (table CharacterTable) from() *postgrest.QueryBuilder {
	return table.client.From("characters")
}

func (table CharacterTable) fromRevealedFields() *postgrest.QueryBuilder {
	return table.client.From("character_revealed_fields")
}
