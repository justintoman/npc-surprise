package router

import (
	"github.com/gin-gonic/gin"
	"github.com/justintoman/npc-surprise/db"
)

type CreateActionInput struct {
	db.ActionBase
}

func (r *Router) CreateAction(c *gin.Context, input *CreateActionInput) (db.Action, error) {
	return r.db.Action.Create(db.CreateActionPayload{ActionBase: input.ActionBase})
}

type UpdateActionInput struct {
	db.Action
}

func (r *Router) UpdateAction(c *gin.Context, input *UpdateActionInput) (db.Action, error) {
	return r.db.Action.Update(db.UpdateActionPayload{Action: input.Action})
}

func (r *Router) DeleteAction(c *gin.Context) error {
	id := c.Params.ByName("id")
	return r.db.Action.Delete(id)
}
