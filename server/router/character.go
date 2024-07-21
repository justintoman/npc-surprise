package router

import (
	"github.com/gin-gonic/gin"
	"github.com/justintoman/npc-surprise/db"
)

type CreateCharacterInput struct {
	Name        string `json:"name" binding:"required"`
	Race        string `json:"race,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Age         string `json:"age,omitempty"`
	Description string `json:"description,omitempty"`
	Appearance  string `json:"appearance,omitempty"`
}

func (r Router) CreateCharacter(c *gin.Context, input *db.CreateCharacterPayload) error {
	character, fields, err := r.CharacterService.Create(*input)
	if err != nil {
		return err
	}
	r.stream.SendAdminCharacterMessageWithFields(character, fields)
	return nil
}

func (r Router) UpdateCharacter(c *gin.Context, input *db.Character) error {
	character, err := r.CharacterService.Update(*input)
	if err != nil {
		return err
	}
	playerCharacter, err := r.CharacterService.Redact(character)
	if err != nil {
		return err
	}
	r.stream.SendAdminCharacterMessage(character)
	r.stream.SendPlayerCharacterMessage(playerCharacter)
	return nil
}

type AssignCharacterInput struct {
	CharacterId int `uri:"characterId" binding:"required,gt=0"`
	PlayerId    int `uri:"playerId" binding:"required,gt=0"`
}

func (r Router) AssignCharacter(c *gin.Context) error {
	var input AssignCharacterInput
	err := c.BindUri(&input)
	if err != nil {
		return err
	}

	prevPlayerId, character, err := r.CharacterService.Assign(input.CharacterId, input.PlayerId)
	if err != nil {
		return err
	}
	redacted, err := r.CharacterService.Redact(character)
	if err != nil {
		return err
	}
	r.stream.SendPlayerCharacterMessage(redacted)
	if prevPlayerId != nil {
		r.stream.SendHideCharacterMessage(*prevPlayerId, character)
	} else {
		r.stream.SendAdminCharacterMessage(character)
	}
	return nil
}

type UnassignCharacterInput struct {
	CharacterId int `uri:"characterId" binding:"required,gt=0"`
}

func (r Router) UnassignCharacter(c *gin.Context) error {
	var input UnassignCharacterInput
	err := c.BindUri(&input)
	if err != nil {
		return err
	}

	playerId, character, err := r.CharacterService.Unassign(input.CharacterId)
	if err != nil {
		return err
	}
	r.stream.SendHideCharacterMessage(playerId, character)
	return nil
}

type CharacterReveleadFieldsInput struct {
	CharacterId int   `json:"characterId" validate:"required"`
	Name        *bool `json:"name" validate:"required"`
	Race        *bool `json:"race" validate:"required"`
	Gender      *bool `json:"gender" validate:"required"`
	Age         *bool `json:"age" validate:"required"`
	Description *bool `json:"description" validate:"required"`
	Appearance  *bool `json:"appearance" validate:"required"`
}

func (r Router) UpdateRevealedFields(c *gin.Context, input *CharacterReveleadFieldsInput) error {
	character, fields, err := r.CharacterService.UpdateRevealedFields(db.CharacterReveleadFields{
		CharacterId: input.CharacterId,
		Name:        *input.Name,
		Race:        *input.Race,
		Gender:      *input.Gender,
		Age:         *input.Age,
		Description: *input.Description,
		Appearance:  *input.Appearance,
	})
	if err != nil {
		return err
	}
	redacted, err := r.CharacterService.Redact(character)
	if err != nil {
		return err
	}
	r.stream.SendAdminCharacterMessageWithFields(character, fields)
	r.stream.SendPlayerCharacterMessage(redacted)
	return nil
}

type DeleteCharacterInput struct {
	Id int `uri:"id" binding:"required,gt=0"`
}

func (r Router) DeleteCharacter(c *gin.Context) error {
	var input DeleteCharacterInput
	err := c.BindUri(&input)
	if err != nil {
		return err
	}

	err = r.db.Character.Delete(input.Id)
	if err != nil {
		return err
	}

	r.stream.SendDeleteCharacterMessage(input.Id)
	return nil
}
