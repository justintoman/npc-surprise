package db

import (
	"encoding/json"
	"log/slog"
	"strconv"

	"github.com/supabase-community/supabase-go"
)

type CreateCharacterPayload struct {
	CharacterBase
}

type UpdateCharacterPayload struct {
	Character
}

type Character struct {
	CharacterBase `json:",inline"`
	Id            int `json:"id"`
	PlayerId      int `json:"player_id,omitempty"`
}

type CharacterBase struct {
	Name        string `json:"name,omitempty"`
	Race        string `json:"race,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Age         string `json:"age,omitempty"`
	Description string `json:"description,omitempty"`
	Appearance  string `json:"appearance,omitempty"`
}

type CharacterReveleadFields struct {
	CharacterId int  `json:"character_id"`
	Name        bool `json:"name"`
	Race        bool `json:"race"`
	Gender      bool `json:"gender"`
	Age         bool `json:"age"`
	Description bool `json:"description"`
	Appearance  bool `json:"appearance"`
}

type CharacterTable struct {
	client *supabase.Client
}

func (db CharacterTable) GetAll() (*[]Character, error) {
	data, _, err := db.client.From("characters").Select("*", "exact", false).Execute()
	var characters *[]Character
	json.Unmarshal(data, characters)
	return characters, err
}

func (db CharacterTable) Get(id int) (*Character, error) {
	data, _, err := db.client.From("characters").Select("*", "exact", false).Filter("id", "eq", strconv.Itoa(id)).Single().Execute()
	var character *Character
	json.Unmarshal(data, character)
	return character, err
}

func (db CharacterTable) Create(character CreateCharacterPayload) (*Character, error) {
	data, _, err := db.client.From("characters").Insert(character, true, "", "", "exact").Single().Execute()
	if err != nil {
		slog.Error("Error creating character", "error", err)
		return nil, err
	}
	var result Character
	err = json.Unmarshal(data, &result)
	if err != nil {
		slog.Error("Error unmarshalling character", "error", err)
		return nil, err
	}

	_, _, err = db.client.From("character_revealed_fields").Insert(map[string]any{"character_id": result.Id}, true, "", "", "exact").Execute()
	if err != nil {
		slog.Error("Error creating character_revelead_fields", "error", err)
		return nil, err
	}

	return &result, err
}

func (db CharacterTable) Update(character UpdateCharacterPayload) (*Character, error) {
	data, _, err := db.client.From("characters").Insert(character, true, "", "", "exact").Execute()
	var result *Character
	json.Unmarshal(data, result)
	return result, err
}

func (db CharacterTable) Delete(id string) error {
	_, _, err := db.client.From("characters").Delete("", "").Filter("id", "eq", id).Execute()
	return err
}

func (db CharacterTable) GetRevealedFields(id string) (CharacterReveleadFields, error) {
	data, _, err := db.client.From("actions").Select("*", "exact", false).Filter("character_id", "eq", id).Single().Execute()
	var revealedFields CharacterReveleadFields
	json.Unmarshal(data, &revealedFields)
	return revealedFields, err
}

type UpdateRevealedFieldsPayload struct {
	CharacterReveleadFields
}

func (db CharacterTable) UpdateRevealedFields(revealedFields UpdateRevealedFieldsPayload) (*CharacterReveleadFields, error) {
	data, _, err := db.client.From("actions").Insert(revealedFields, true, "", "", "exact").Execute()
	var result CharacterReveleadFields
	json.Unmarshal(data, &result)
	return &result, err
}
