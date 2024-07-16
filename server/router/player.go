package router

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/justintoman/npc-surprise/db"
	"github.com/justintoman/npc-surprise/stream"
)

func (r *Router) PlayerMiddleware(c *gin.Context) {
	requestKey, err := c.Cookie("admin_key")
	if err != nil {
		requestKey = ""
	}
	isAdmin := requestKey == r.AdminKey && requestKey != ""

	c.Set("isAdmin", isAdmin)
	if isAdmin {
		c.Set("player_id", "admin")
		c.Set("player_name", "Admin")
		c.SetCookie("player_id", "admin", 0, "/", "", false, true)
		c.SetCookie("player_name", "Admin", 0, "/", "", false, true)
		return
	}

	playerId, err := c.Cookie("player_id")
	if err != nil {
		c.AbortWithStatusJSON(401, ErrorResponse{Message: "Unauthorized. Login with a name.", Status: 401})
		return
	}

	player, err := r.db.Player.Get(playerId)
	if err != nil {
		// invalid user ID I guess
		// clear the cookie
		c.SetCookie("player_id", "", -1, "/", "", false, true)
		c.AbortWithStatusJSON(404, ErrorResponse{Message: "Invalid Player. Try logging in again.", Status: 404})
		return
	}

	c.Set("player", player)
	c.Next()
}

type AssignInput struct {
	Type     string `json:"type"`
	Id       int    `json:"id"`
	PlayerId int    `json:"player_id"`
}

func (r Router) AssignToPlayer(c *gin.Context, input *AssignInput) error {
	switch input.Type {
	case "character":
		character, err := r.db.Character.Get(input.Id)
		if err != nil {
			slog.Error("error getting character to assign", "error", err, "character_id", input.Id)
			return err
		}
		if character.PlayerId == input.PlayerId {
			slog.Info("already assigned", "character_id", input.Id, "player_id", input.PlayerId)
			return nil
		}
		if character.PlayerId != 0 {
			r.UnassignFromPlayer(c, input)
		}
		character.PlayerId = input.PlayerId
		character, err = r.db.Character.Update(db.UpdateCharacterPayload{
			Character: character,
		})
		if err != nil {
			return err
		}

		// if there any any actions from this character that are assigned, unassign them
		actions, err := r.db.Action.GetAll(input.Id)
		if err != nil {
			return err
		}
		for _, action := range actions {
			if action.PlayerId != input.PlayerId {
				r.UnassignFromPlayer(c, &AssignInput{
					Type:     "action",
					Id:       action.Id,
					PlayerId: input.PlayerId,
				})
			}
		}

		revealedFields, err := r.db.Character.GetRevealedFields(strconv.Itoa(input.Id))
		if err != nil {
			slog.Error("error getting revealed fields", "error", err)
			return err
		}

		redactCharacter(&character, revealedFields)

		payload := stream.MessagePayload{
			Type: "assign-character",
			Data: CharacterWithActions{Character: character, Actions: make([]db.Action, 0)},
		}
		r.stream.SendMessage(strconv.Itoa(input.PlayerId), payload)
		return nil
	case "action":
		action, err := r.db.Action.Get(strconv.Itoa(input.Id))
		if err != nil {
			return err
		}
		if action.PlayerId == input.PlayerId {
			slog.Info("already assigned", "action_id", input.Id, "player_id", input.PlayerId)
			return nil
		}
		action.PlayerId = input.PlayerId
		action, err = r.db.Action.Update(db.UpdateActionPayload{
			Action: action,
		})
		if err != nil {
			return err
		}
		payload := stream.MessagePayload{
			Type: "assign-action",
			Data: action,
		}
		r.stream.SendMessage(strconv.Itoa(input.PlayerId), payload)
		return nil
	default:
		return fmt.Errorf("invalid type: %s", input.Type)
	}
}

func (r Router) UnassignFromPlayer(c *gin.Context, input *AssignInput) error {
	switch input.Type {
	case "character":
		_, err := r.db.Character.Update(db.UpdateCharacterPayload{
			Character: db.Character{
				Id:       input.Id,
				PlayerId: 0,
			},
		})
		if err != nil {
			return err
		}

		// unassign all actions from this character
		actions, err := r.db.Action.GetAll(input.Id)
		if err != nil {
			return err
		}
		for _, action := range actions {
			if action.PlayerId != 0 {
				r.db.Action.Update(db.UpdateActionPayload{
					Action: db.Action{
						Id:       action.Id,
						PlayerId: 0,
					},
				})
			}
		}

		payload := stream.MessagePayload{
			Type: "unassign-character",
			Data: input.Id,
		}
		r.stream.SendMessage(strconv.Itoa(input.PlayerId), payload)
		return nil
	case "action":
		_, err := r.db.Action.Update(db.UpdateActionPayload{
			Action: db.Action{
				Id:       input.Id,
				PlayerId: 0,
			},
		})
		if err != nil {
			return err
		}

		payload := stream.MessagePayload{
			Type: "unassign-action",
			Data: input.Id,
		}
		r.stream.SendMessage(strconv.Itoa(input.PlayerId), payload)
		return nil
	default:
		return fmt.Errorf("invalid type: %s", input.Type)
	}
}

func (r *Router) DeletePlayer(c *gin.Context) error {
	id := c.Params.ByName("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		slog.Error("invalid player id", "error", err, "player_id", id)
	}
	characters, err := r.db.Character.GetAllByPlayerId(id)
	if err != nil {
		slog.Error("Error trying to unassign characters from player", "error", err)
		return err
	}
	for _, character := range characters {
		err = r.UnassignFromPlayer(c, &AssignInput{
			Type:     "character",
			Id:       character.Id,
			PlayerId: idInt,
		})
		if err != nil {
			slog.Error("Error trying to unassign characters from player", "error", err, "character", character)
			return err
		}
	}
	err = r.db.Player.Delete(id)
	if err != nil {
		return err
	}
	r.stream.SendMessage("admin", stream.MessagePayload{Type: "player-deleted", Data: idInt})
	return nil
}
