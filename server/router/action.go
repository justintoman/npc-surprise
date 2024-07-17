package router

import (
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/justintoman/npc-surprise/db"
	"github.com/justintoman/npc-surprise/stream"
)

type CreateActionInput struct {
	db.ActionBase
}

func (r *Router) CreateAction(c *gin.Context, input *CreateActionInput) (db.Action, error) {
	action, err := r.db.Action.Create(db.CreateActionPayload{ActionBase: input.ActionBase})
	if err != nil {
		slog.Error("Error updating action", "error", err)
		return db.Action{}, err
	}
	r.stream.SendMessage("admin", stream.MessagePayload{Type: "action", Data: action})
	return action, nil
}

type UpdateActionInput struct {
	db.Action
}

func (r *Router) UpdateAction(c *gin.Context, input *UpdateActionInput) (db.Action, error) {
	_, err := r.db.Action.Update(db.UpdateActionPayload{Action: input.Action})
	if err != nil {
		slog.Error("Error updating action", "error", err)
		return db.Action{}, err
	}
	action, err := r.db.Action.Get(strconv.Itoa(input.Id))
	if err != nil {
		slog.Error("Error fetching updated action", "error", err)
		return db.Action{}, err
	}
	r.stream.SendMessage("admin", stream.MessagePayload{Type: "action", Data: action})
	if action.PlayerId != 0 {
		r.stream.SendMessage(strconv.Itoa(action.PlayerId), stream.MessagePayload{Type: "action", Data: action})
	}
	return action, nil
}

func (r *Router) DeleteAction(c *gin.Context) error {
	id := c.Params.ByName("id")
	err := r.db.Action.Delete(id)
	if err != nil {
		slog.Error("Error deleting action", "error", err)
		return err
	}
	r.stream.SendMessage("admin", stream.MessagePayload{Type: "delete-action", Data: id})
	return nil
}
