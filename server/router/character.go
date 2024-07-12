package router

import (
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/justintoman/npc-surprise/db"
	"github.com/justintoman/npc-surprise/stream"
)

type CreateCharacterInput struct {
	db.CharacterBase
}

func (r Router) CreateCharacter(c *gin.Context, input *CreateCharacterInput) (*db.Character, error) {
	return r.db.Character.Create(db.CreateCharacterPayload{CharacterBase: input.CharacterBase})
}

type GetCharactersOuptut struct {
	Characters []CharacterWithActions `json:"characters"`
}

type CharacterWithActions struct {
	db.Character `json:",inline"`
	Actions      []db.Action `json:"actions"`
}

func (r Router) GetCharacters(c *gin.Context) (*GetCharactersOuptut, error) {
	characters, err := r.db.Character.GetAll()
	if err != nil {
		return nil, err
	}
	results := make([]CharacterWithActions, len(*characters))
	for i, character := range *characters {
		actions, err := r.db.Action.GetAll(character.Id)
		if err != nil {
			return nil, err
		}
		results[i] = CharacterWithActions{Character: character, Actions: *actions}
	}
	return &GetCharactersOuptut{Characters: results}, nil
}

type UpdateCharacterInput struct {
	db.Character
}

func (r Router) UpdateCharacter(c *gin.Context, input *UpdateCharacterInput) (*db.Character, error) {
	return r.db.Character.Update(db.UpdateCharacterPayload{Character: input.Character})
}

type UpdateRevealedFieldsInput struct {
	db.CharacterReveleadFields
}

func (r Router) UpdateRevealedFields(c *gin.Context, input *UpdateRevealedFieldsInput) (*db.CharacterReveleadFields, error) {
	result, err := r.db.Character.UpdateRevealedFields(db.UpdateRevealedFieldsPayload{CharacterReveleadFields: input.CharacterReveleadFields})
	if err != nil {
		slog.Error("Error updating revealed fields", "error", err)
		return nil, err
	}

	character, err := r.db.Character.Get(result.CharacterId)
	if err != nil {
		slog.Error("Error fetching character after updating revealed fields", "error", err)
		return nil, err
	}

	if character.PlayerId == 0 {
		return result, nil
	}

	if !result.Age {
		character.Age = ""
	}
	if !result.Gender {
		character.Gender = ""
	}
	if !result.Race {
		character.Race = ""
	}
	if !result.Name {
		character.Name = ""
	}
	if !result.Description {
		character.Description = ""
	}
	if !result.Appearance {
		character.Appearance = ""
	}

	payload := stream.MessagePayload{
		Type: "character",
		Data: character,
	}

	r.stream.SendMessage(strconv.Itoa(character.PlayerId), payload)

	return result, nil
}
