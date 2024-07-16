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
	character, err := r.db.Character.Create(db.CreateCharacterPayload{CharacterBase: input.CharacterBase})
	if err != nil {
		slog.Error("Error creating character", "error", err)
		return nil, err
	}
	slog.Info("created character", "character", character)
	data := CharacterWithActions{
		Character: *character,
		Actions:   make([]db.Action, 0),
	}
	r.stream.SendMessage("admin", stream.MessagePayload{Type: "character", Data: data})
	return character, nil
}

type CharacterWithActions struct {
	db.Character `json:",inline"`
	Actions      []db.Action `json:"actions"`
}

func (r Router) GetCharacters(c *gin.Context) (*[]CharacterWithActions, error) {
	characters, err := r.db.Character.GetAll()
	if err != nil {
		slog.Error("Error fetching characters", "error", err)
		return nil, err
	}
	results := make([]CharacterWithActions, len(characters))
	for i, character := range characters {
		actions, err := r.db.Action.GetAll(character.Id)
		if err != nil {
			slog.Error("Error fetching actions for each character", "error", err, "character", character)
			return nil, err
		}
		results[i] = CharacterWithActions{Character: character, Actions: actions}
	}
	return &results, nil
}

type UpdateCharacterInput struct {
	db.Character
}

func (r Router) UpdateCharacter(c *gin.Context, input *UpdateCharacterInput) (*db.Character, error) {
	character, err := r.db.Character.Update(db.UpdateCharacterPayload{Character: input.Character})
	if err != nil {
		slog.Error("Error updating character", "error", err)
		return nil, err
	}
	slog.Info("updated character", "character", character)
	actions, err := r.db.Action.GetAll(character.Id)
	if err != nil {
		slog.Error("Error fetching actions after updating character", "error", err)
		return nil, err
	}
	data := CharacterWithActions{
		Character: character,
		Actions:   actions,
	}
	r.stream.SendMessage("admin", stream.MessagePayload{Type: "character", Data: data})
	return &character, nil
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

	character, err := r.getRedactedCharacter(input.CharacterId)
	if err != nil {
		return nil, err
	}

	payload := stream.MessagePayload{
		Type: "character",
		Data: character,
	}

	r.stream.SendMessage(strconv.Itoa(character.PlayerId), payload)

	return result, nil
}

func (r Router) getRedactedCharacter(id int) (*CharacterWithActions, error) {
	character, err := r.db.Character.Get(id)
	if err != nil {
		slog.Error("Error fetching character", "error", err)
		return nil, err
	}

	revealedFields, err := r.db.Character.GetRevealedFields(strconv.Itoa(character.Id))
	if err != nil {
		slog.Error("Error fetching character revealed fields", "error", err)
		return nil, err
	}

	redactCharacter(&character, revealedFields)

	actions, err := r.db.Action.GetAllAssigned(strconv.Itoa(character.PlayerId))
	if err != nil {
		slog.Error("Error fetching assigned actions for player", "error", err, "player_id", character.PlayerId)
		return nil, err
	}

	characterActions := make([]db.Action, 0)
	for _, action := range actions {
		if action.CharacterId != character.Id {
			continue
		}
		characterActions = append(characterActions, action)
	}

	return &CharacterWithActions{
		Character: character,
		Actions:   characterActions,
	}, nil
}

func redactCharacter(character *db.Character, revealedFields db.CharacterReveleadFields) {
	if !revealedFields.Age {
		character.Age = ""
	}
	if !revealedFields.Gender {
		character.Gender = ""
	}
	if !revealedFields.Race {
		character.Race = ""
	}
	if !revealedFields.Name {
		character.Name = ""
	}
	if !revealedFields.Description {
		character.Description = ""
	}
	if !revealedFields.Appearance {
		character.Appearance = ""
	}
	return
}

func (r Router) DeleteCharacter(c *gin.Context) error {
	id := c.Params.ByName("id")
	return r.db.Character.Delete(id)
}
