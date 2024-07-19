package router

import (
	"github.com/gin-gonic/gin"
	"github.com/justintoman/npc-surprise/db"
)

func (r *Router) CreateAction(c *gin.Context, input *db.CreateActionPayload) error {
	action, err := r.ActionService.Create(*input)
	if err != nil {
		return err
	}
	r.stream.SendActionMessage(action)
	return nil
}

func (r *Router) UpdateAction(c *gin.Context, input *db.Action) error {
	action, err := r.ActionService.Update(*input)
	if err != nil {
		return err
	}
	r.stream.SendActionMessage(action)
	return nil
}

type AssignActionInput struct {
	ActionId int `uri:"id" binding:"required,gt=0"`
}

func (r Router) AssignAction(c *gin.Context) error {
	var input AssignActionInput
	err := c.BindUri(&input)
	if err != nil {
		return err
	}

	action, err := r.ActionService.Assign(input.ActionId)
	if err != nil {
		return err
	}
	r.stream.SendActionMessage(action)
	return nil
}

type UnassignActionInput struct {
	ActionId int `uri:"id" binding:"required,gt=0"`
}

func (r Router) UnassignAction(c *gin.Context) error {
	var input UnassignActionInput
	err := c.BindUri(&input)
	if err != nil {
		return err
	}
	unassignedPlayerId, err := r.ActionService.Unassign(input.ActionId)
	if err != nil {
		return err
	}
	r.stream.SendUnassignActionMessage(unassignedPlayerId, input.ActionId)
	return nil

}

type DeleteInput struct {
	Id int `uri:"id" binding:"required,gt=0"`
}

func (r *Router) DeleteAction(c *gin.Context) error {
	var input DeleteInput
	err := c.BindUri(&input)
	if err != nil {
		return err
	}

	err = r.ActionService.Delete(input.Id)
	if err != nil {
		return err
	}
	r.stream.SendDeleteActionMessage(input.Id)
	return nil
}
