package router

import (
	"encoding/json"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/justintoman/npc-surprise/pkg/db"
	"github.com/justintoman/npc-surprise/pkg/stream"
)

type LoginInput struct {
	Name string `json:"name" binding:"required"`
}

func (r Router) Login(c *gin.Context, input *LoginInput) (LoginResponse, error) {
	player, err := r.parsePlayerFromCookie(c)
	if err != nil {
		if input.Name == r.AdminKey {
			// admin logging in
			err = setPlayerCookie(c, db.Player{
				Id:   0,
				Name: "Admin",
			})
			if err != nil {
				return LoginResponse{}, err
			}
			return LoginResponse{
				Id:      0,
				Name:    "Admin",
				IsAdmin: true,
			}, nil
		}

		// just a normie player logging in
		player, err := r.PlayerService.Create(input.Name)
		if err != nil {
			return LoginResponse{}, err
		}

		err = setPlayerCookie(c, player)
		if err != nil {
			return LoginResponse{}, err
		}

		return LoginResponse{
			Id:      player.Id,
			Name:    player.Name,
			IsAdmin: false,
		}, nil
	}

	if player.Id == 0 {
		// they're somehow the admin and logged in but logging in again?
		// well it's not illegal...
		return LoginResponse{
			Id:      0,
			Name:    "Admin",
			IsAdmin: true,
		}, nil
	}

	// they already have a player,
	// so update it with their new name
	player, err = r.PlayerService.Update(db.Player{
		Id:   player.Id,
		Name: input.Name,
	})

	if err != nil {
		clearPlayerCookie(c)
		return LoginResponse{}, err
	}

	err = setPlayerCookie(c, player)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		Id:      player.Id,
		Name:    player.Name,
		IsAdmin: false,
	}, nil
}

type LoginResponse struct {
	Id      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	IsAdmin bool   `json:"isAdmin"`
}

type StatusResponse LoginResponse

func (r Router) Status(c *gin.Context) (StatusResponse, error) {
	player, err := r.parsePlayerFromCookie(c)
	if err != nil {
		return StatusResponse{
			IsAdmin: false,
		}, nil
	}

	response := StatusResponse{
		Id:      player.Id,
		Name:    player.Name,
		IsAdmin: player.Id == stream.AdminPlayerId,
	}
	return response, nil
}

func (r Router) AdminMiddleware(c *gin.Context) {
	player, err := r.parsePlayerFromCookie(c)
	if err != nil {
		c.AbortWithStatusJSON(401, ErrorResponse{Message: "Unauthorized. You are not Justin.", Status: 401})
		return
	}

	if player.Id != stream.AdminPlayerId {
		c.AbortWithStatusJSON(401, ErrorResponse{Message: "Unauthorized. You are not Justin.", Status: 401})
		return
	}

	c.Next()
}

func (r *Router) PlayerMiddleware(c *gin.Context) {
	player, err := r.parsePlayerFromCookie(c)
	if err != nil {
		c.AbortWithStatusJSON(401, ErrorResponse{Message: "Invalid Player. Try logging in again.", Status: 401})
		return
	}

	c.Set("player", player)
	c.Next()
}

func (r *Router) parsePlayerFromCookie(c *gin.Context) (db.Player, error) {
	cookie, err := c.Cookie("player")
	if err != nil {
		// player cookie is boned, unset it
		slog.Info("player cookie is boned, unset it", "error", err)
		clearPlayerCookie(c)
		return db.Player{}, err
	}

	var player db.Player
	err = json.Unmarshal([]byte(cookie), &player)
	if err != nil {
		// player cookie is boned, unset it
		slog.Info("player cookie failed to be unmarshalled, unset it", "error", err)
		clearPlayerCookie(c)
		return db.Player{}, err
	}

	if player.Id == stream.AdminPlayerId {
		slog.Info("player cookie is admin", "playerId", player.Id, "playerName", player.Name)
		return player, nil
	}

	_, err = r.PlayerService.Get(player.Id)
	if err != nil {
		// player cookie is boned, unset it
		slog.Info("cookie seems good but there's no player in the db for it", "error", err, "player", player)
		clearPlayerCookie(c)
		return db.Player{}, err
	}

	slog.Info("player cookie is valid", "playerId", player.Id, "playerName", player.Name)

	return player, nil
}

func setPlayerCookie(c *gin.Context, player db.Player) error {
	cookie, err := json.Marshal(player)
	if err != nil {
		clearPlayerCookie(c)
		return err
	}
	c.SetCookie("player", string(cookie), 0, "/", "", false, true)
	return nil
}

func clearPlayerCookie(c *gin.Context) {
	c.SetCookie("player", "", -1, "/", "", false, true)
}
