package router

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/justintoman/npc-surprise/db"
	"github.com/justintoman/npc-surprise/services"
	"github.com/justintoman/npc-surprise/stream"
	"github.com/loopfz/gadgeto/tonic"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type Router struct {
	stream           stream.StreamingServer
	db               db.Db
	AdminKey         string
	ActionService    services.ActionService
	CharacterService services.CharacterService
	PlayerService    services.PlayerService
}

func New(db db.Db, adminKey string) *gin.Engine {
	streamService := stream.New(db)

	router := Router{
		stream:           streamService,
		db:               db,
		AdminKey:         adminKey,
		ActionService:    services.NewActionService(db),
		CharacterService: services.NewCharacterService(db, streamService),
		PlayerService:    services.NewPlayerService(db, streamService),
	}

	g := gin.Default()

	g.POST("/login", tonic.Handler(router.Login, 200))
	g.GET("/status", tonic.Handler(router.Status, 200))

	adminRoutes := g.Group("/")
	adminRoutes.Use(router.AdminMiddleware)
	adminRoutes.DELETE("players/:id", tonic.Handler(router.DeletePlayer, 200))

	characterRoutes := adminRoutes.Group("/characters")
	characterRoutes.POST("", tonic.Handler(router.CreateCharacter, 200))
	characterRoutes.PUT("/:characterId", tonic.Handler(router.UpdateCharacter, 200))
	characterRoutes.PUT("/:characterId/assign/:playerId", tonic.Handler(router.AssignCharacter, 200))
	characterRoutes.PUT("/:characterId/unassign", tonic.Handler(router.UnassignCharacter, 200))
	characterRoutes.PUT("/:characterId/reveal", tonic.Handler(router.UpdateRevealedFields, 200))
	characterRoutes.DELETE("/:characterId", tonic.Handler(router.DeleteCharacter, 200))

	actionRoutes := characterRoutes.Group("/:characterId/actions")
	actionRoutes.POST("", tonic.Handler(router.CreateAction, 200))
	actionRoutes.PUT(":actionId", tonic.Handler(router.UpdateAction, 200))
	actionRoutes.PUT(":actionId/reveal", tonic.Handler(router.RevealAction, 200))
	actionRoutes.PUT(":actionId/hide", tonic.Handler(router.HideAction, 200))
	actionRoutes.DELETE(":actionId", tonic.Handler(router.DeleteAction, 200))

	authRoutes := g.Group("/")

	middleware, handler := streamService.NewUserStream(router.onPlayerConnected, router.onPlayerDisconnected)
	authRoutes.GET("/stream", router.PlayerMiddleware, middleware, tonic.Handler(handler, 200))

	ctx := context.Background()
	go streamService.Listen(ctx)

	return g
}

func (r *Router) onPlayerConnected(player db.Player) {
	if player.Id == stream.AdminPlayerId {
		characters, fields, err := r.CharacterService.GetAllWithActionsAndFields()
		if err != nil {
			slog.Error("error getting characters for player", "error", err)
			return
		}
		players, err := r.PlayerService.GetAll()
		if err != nil {
			slog.Error("error getting players", "error", err)
			return
		}
		r.stream.SendInitAdminMessage(players, characters, fields)
	} else {
		r.stream.SendPlayerConnectedMessage(player)
		characters, err := r.CharacterService.GetAllAssignedWithActionsRedacted(player.Id)
		if err != nil {
			slog.Error("error getting characters for player", "error", err, "playerId", player.Id)
			return
		}
		r.stream.SendInitPlayerMessage(player.Id, characters)
	}
}

func (r *Router) onPlayerDisconnected(player db.Player) {
	r.stream.SendPlayerDisconnectedMessage(player.Id)
}
