package router

import (
	"ai-chat-platform/backend/config"
	"ai-chat-platform/backend/handler"
	"ai-chat-platform/backend/middleware"
	"ai-chat-platform/backend/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func New(cfg config.Config, db *gorm.DB) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.CorsAllowOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	healthHandler := handler.HealthHandler{Config: cfg}
	modelHandler := handler.ModelHandler{DB: db}
	channelHandler := handler.ChannelHandler{DB: db}
	abilityHandler := handler.AbilityHandler{DB: db}
	keyHandler := handler.KeyHandler{DB: db}
	optionHandler := handler.OptionHandler{DB: db}
	sessionHandler := handler.SessionHandler{DB: db}
	messageHandler := handler.MessageHandler{DB: db}
	usageLogHandler := handler.UsageLogHandler{DB: db}
	gatewayHandler := handler.GatewayHandler{
		DB:             db,
		DeepSeekClient: service.NewDeepSeekClient(cfg.RequestTimeoutSeconds),
	}

	r.GET("/health", healthHandler.Get)

	api := r.Group("/api")
	{
		api.GET("/models", modelHandler.List)
		api.GET("/models/:id", modelHandler.Get)
		api.POST("/models", modelHandler.Create)
		api.PUT("/models/:id", modelHandler.Update)
		api.DELETE("/models/:id", modelHandler.Delete)

		api.GET("/channels", channelHandler.List)
		api.GET("/channels/:id", channelHandler.Get)
		api.POST("/channels", channelHandler.Create)
		api.PUT("/channels/:id", channelHandler.Update)
		api.DELETE("/channels/:id", channelHandler.Delete)

		api.GET("/abilities", abilityHandler.List)
		api.POST("/abilities", abilityHandler.Create)
		api.PUT("/abilities/:id", abilityHandler.Update)
		api.DELETE("/abilities/:id", abilityHandler.Delete)

		api.GET("/keys", keyHandler.List)
		api.GET("/keys/:id", keyHandler.Get)
		api.POST("/keys", keyHandler.Create)
		api.PUT("/keys/:id", keyHandler.Update)
		api.DELETE("/keys/:id", keyHandler.Delete)

		api.GET("/options", optionHandler.List)
		api.POST("/options", optionHandler.Create)
		api.PUT("/options/:id", optionHandler.Update)
		api.DELETE("/options/:id", optionHandler.Delete)

		api.GET("/sessions", sessionHandler.List)
		api.GET("/sessions/:id", sessionHandler.Get)
		api.POST("/sessions", sessionHandler.Create)
		api.PUT("/sessions/:id", sessionHandler.Update)
		api.DELETE("/sessions/:id", sessionHandler.Delete)

		api.GET("/sessions/:id/messages", messageHandler.List)
		api.POST("/sessions/:id/messages", messageHandler.Create)
		api.DELETE("/messages/:id", messageHandler.Delete)

		api.GET("/usage-logs", usageLogHandler.List)
		api.POST("/usage-logs", usageLogHandler.Create)
		api.DELETE("/usage-logs/:id", usageLogHandler.Delete)
	}

	v1 := r.Group("/v1")
	v1.Use(middleware.APIKeyAuth(db))
	{
		v1.GET("/models", gatewayHandler.ListModels)
		v1.POST("/chat/completions", gatewayHandler.ChatCompletions)
	}

	return r
}
