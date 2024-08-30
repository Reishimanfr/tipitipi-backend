package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"bash06/strona-fundacja/src/backend/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	Db  core.Database
	Log *zap.Logger
	A   core.Argon2idHash
}

type Config struct {
	Router *gin.Engine
}

func NewHandler(cfg *Config, db *core.Database) {
	h := &Handler{}
	h.Db = *db
	h.Log = core.GetLogger()
	h.A = *core.NewArgon2idHash(1, 32, 64*1024, 32, 256)

	cfg.Router.
		HEAD("/heartbeat", h.Heartbeat).
		POST("/admin/login", h.AdminLogin)

	// Define protected routes with a different group
	protected := cfg.Router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.DELETE("/api/blog/delete/:id", h.delete)
		protected.POST("/api/blog/create", h.create)
		protected.PATCH("/api/blog/edit/:id", h.edit)
		protected.GET("/api/blog/post/:id", h.post)
		protected.GET("/api/blog/posts/", h.posts)
		protected.POST("/admin/changePassword", h.changePassword)
	}
}
