package router

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/justintoman/npc-surprise/db"
	"github.com/justintoman/npc-surprise/stream"
	"github.com/loopfz/gadgeto/tonic"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type Router struct {
	stream   stream.StreamingServer
	db       db.Db
	AdminKey string
}

func New(db db.Db, adminKey string) *gin.Engine {
	streamService := stream.New(db)

	router := Router{stream: streamService, db: db, AdminKey: adminKey}

	g := gin.Default()

	g.POST("/login", tonic.Handler(router.Login, 200))
	g.GET("/status", tonic.Handler(router.Status, 200))

	adminRoutes := g.Group("/")
	adminRoutes.Use(router.AdminMiddleware)
	adminRoutes.GET("players", tonic.Handler(router.GetPlayers, 200))
	adminRoutes.DELETE("players/:id", tonic.Handler(router.DeletePlayer, 200))

	assignRoutes := adminRoutes.Group("/assign")
	assignRoutes.POST("", tonic.Handler(router.AssignToPlayer, 200))
	assignRoutes.DELETE("", tonic.Handler(router.UnassignFromPlayer, 200))

	characterRoutes := adminRoutes.Group("/characters")
	characterRoutes.GET("", tonic.Handler(router.GetCharacters, 200))
	characterRoutes.POST("", tonic.Handler(router.CreateCharacter, 200))
	characterRoutes.PUT("/:id", tonic.Handler(router.UpdateCharacter, 200))
	characterRoutes.PUT("/:id/reveal", tonic.Handler(router.UpdateRevealedFields, 200))
	characterRoutes.DELETE("/:id", tonic.Handler(router.DeleteCharacter, 200))

	actionRoutes := adminRoutes.Group("/actions")
	actionRoutes.POST("", tonic.Handler(router.CreateAction, 200))
	actionRoutes.PUT(":id", tonic.Handler(router.UpdateAction, 200))
	actionRoutes.DELETE(":id", tonic.Handler(router.DeleteAction, 200))

	authRoutes := g.Group("/")

	middleware, handler := streamService.NewUserStream()
	authRoutes.GET("/stream", router.PlayerMiddleware, middleware, tonic.Handler(handler, 200))

	ctx := context.Background()
	go streamService.Listen(ctx)

	return g
}

func (r Router) AdminMiddleware(c *gin.Context) {
	requestKey, err := c.Cookie("admin_key")
	if err != nil {
		requestKey = ""
	}

	if requestKey != r.AdminKey {
		c.AbortWithStatusJSON(401, ErrorResponse{Message: "Unauthorized. You are not Justin.", Status: 401})
		return
	}

	c.Next()
}

func (r Router) Status(c *gin.Context) error {
	playerName, err := c.Cookie("player_name")
	if err != nil {
		playerName = ""
	}
	playerId, err := c.Cookie("player_id")
	if err != nil {
		playerId = ""
	}

	requestKey, err := c.Cookie("admin_key")
	if err != nil {
		requestKey = ""
	}
	isAdmin := requestKey == r.AdminKey && requestKey != ""
	response := StatusResponse{
		PlayerId:   playerId,
		PlayerName: playerName,
		IsAdmin:    isAdmin,
	}
	c.JSON(200, response)
	return nil
}

type StatusResponse struct {
	PlayerId   string `json:"player_id,omitempty"`
	PlayerName string `json:"player_name,omitempty"`
	IsAdmin    bool   `json:"is_admin"`
}

func (r Router) GetPlayers(c *gin.Context) ([]PlayerResponse, error) {
	dbPlayers, err := r.db.Player.GetAll()
	if err != nil {
		return nil, err
	}
	onlinePlayers := r.stream.GetClients()

	result := make([]PlayerResponse, len(*dbPlayers))
	for i, player := range *dbPlayers {
		result[i] = PlayerResponse{Player: player}
		for _, onlinePlayer := range onlinePlayers {
			if result[i].IsOnline {
				continue
			}
			if strconv.Itoa(player.Id) == onlinePlayer.Id {
				result[i].IsOnline = true
			}
		}
	}
	return result, nil
}

type PlayerResponse struct {
	db.Player
	IsOnline bool `json:"is_online"`
}
