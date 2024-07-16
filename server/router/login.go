package router

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/justintoman/npc-surprise/db"
)

type NameInput struct {
	Name string `json:"name"`
}

func (r Router) Login(c *gin.Context, input *NameInput) error {

	if input.Name == r.AdminKey {
		c.SetCookie("player_id", "admin", 0, "/", "", false, true)
		c.SetCookie("player_name", "Admin", 0, "/", "", false, true)
		c.SetCookie("admin_key", input.Name, 0, "/", "", false, true)
		return nil
	}

	playerIdStr, err := c.Cookie("player_id")
	if err != nil {
		// need to create a new player
		player, err := r.db.Player.Create(input.Name)
		if err != nil {
			// whoops 🤷🤷
			c.AbortWithStatusJSON(500, ErrorResponse{Message: "Internal server error", Status: 500})
			return err
		}

		// set player_id in da cookie
		c.SetCookie("player_id", strconv.Itoa(player.Id), 0, "/", "", false, true)
		c.SetCookie("player_name", player.Name, 0, "/", "", false, true)
		c.JSON(200, player)
		return nil
	}

	playerId, err := strconv.Atoi(playerIdStr)
	if err != nil {
		c.SetCookie("player_id", "", -1, "/", "", false, true)
		c.AbortWithStatusJSON(400, ErrorResponse{Message: "Invalid player ID", Status: 400})
		return err
	}

	// existing player, let's update their name
	player, err := r.db.Player.Update(db.Player{
		Id:   playerId,
		Name: input.Name,
	})

	if err != nil {
		// whoops 🤷🤷
		c.AbortWithStatusJSON(500, ErrorResponse{Message: "Internal server error", Status: 500})
		return err
	}

	// set player_id in da cookie
	c.SetCookie("player_id", strconv.Itoa(player.Id), 0, "/", "", false, true)
	c.SetCookie("player_name", player.Name, 0, "/", "", false, true)
	c.JSON(200, player)
	return nil
}
