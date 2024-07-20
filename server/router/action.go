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
	r.stream.SendAdminActionMessage(action)
	return nil
}

func (r *Router) UpdateAction(c *gin.Context, input *db.Action) error {
	action, err := r.ActionService.Update(*input)
	if err != nil {
		return err
	}
	r.stream.SendAdminActionMessage(action)
	return nil
}

type AssignActionInput struct {
	CharacterId int `uri:"characterId" binding:"required,gt=0"`
	ActionId    int `uri:"actionId" binding:"required,gt=0"`
}

func (r Router) RevealAction(c *gin.Context) error {
	var input AssignActionInput
	err := c.BindUri(&input)
	if err != nil {
		return err
	}

	playerId, action, err := r.ActionService.Reveal(input.ActionId)
	if err != nil {
		return err
	}
	r.stream.SendPlayerActionMessage(playerId, action)
	return nil
}

type UnassignActionInput struct {
	ActionId int `uri:"actionId" binding:"required,gt=0"`
}

func (r Router) HideAction(c *gin.Context) error {
	var input UnassignActionInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		return err
	}
	playerId, action, err := r.ActionService.Hide(input.ActionId)
	if err != nil {
		return err
	}
	r.stream.SendHideActionMessage(playerId, action)
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
