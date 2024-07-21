package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/justintoman/npc-surprise/pkg/db"
)

type PlayerResponse struct {
	db.Player `json:",inline"`
	IsOnline  bool `json:"isOnline"`
}

type DeletePlayerInput struct {
	Id int `uri:"id" binding:"required" validate:"gt=0,required"`
}

func (r *Router) DeletePlayer(c *gin.Context) error {
	var input DeletePlayerInput
	err := c.BindUri(&input)
	if err != nil {
		return err
	}

	slog.Info("deleting player", "playerId", input.Id)

	err = r.PlayerService.Delete(input.Id)
	if err != nil {
		return err
	}
	r.stream.SendDeletePlayerMessage(input.Id)
	return nil
}
