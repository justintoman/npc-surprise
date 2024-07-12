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
		character, err := r.db.Character.Update(db.UpdateCharacterPayload{
			Character: db.Character{
				Id:       input.Id,
				PlayerId: input.PlayerId,
			},
		})
		if err != nil {
			return err
		}
		payload := stream.MessagePayload{
			Type: "assign-character",
			Data: character,
		}
		r.stream.SendMessage(strconv.Itoa(input.PlayerId), payload)
		return nil
	case "action":
		action, err := r.db.Action.Update(db.UpdateActionPayload{
			Action: db.Action{
				Id:       input.Id,
				PlayerId: input.PlayerId,
			},
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

func (r Router) UnassignFromPlayer(c *gin.Context) error {
	modelType := c.Query("type")
	idStr := c.Query("id")
	playerIdStr := c.Query("player_id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Error("ID not a valid int", "error", err, "id", idStr)
		return err
	}

	playerId, err := strconv.Atoi(playerIdStr)
	if err != nil {
		slog.Error("Player ID not a valid int", "error", err, "player_id", playerIdStr)
		return err
	}

	switch modelType {
	case "character":
		_, err := r.db.Character.Update(db.UpdateCharacterPayload{
			Character: db.Character{
				Id:       id,
				PlayerId: 0,
			},
		})
		if err != nil {
			return err
		}

		payload := stream.MessagePayload{
			Type: "unassign-character",
			Data: id,
		}
		r.stream.SendMessage(strconv.Itoa(playerId), payload)
		return nil
	case "action":
		_, err := r.db.Action.Update(db.UpdateActionPayload{
			Action: db.Action{
				Id:       id,
				PlayerId: 0,
			},
		})
		if err != nil {
			return err
		}

		payload := stream.MessagePayload{
			Type: "unassign-action",
			Data: id,
		}
		r.stream.SendMessage(strconv.Itoa(playerId), payload)
		return nil
	default:
		return fmt.Errorf("invalid type: %s", modelType)
	}
}
