package routes

import (
	"bash06/strona-fundacja/src/backend/core"

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

	public := cfg.Router.Group("/")
	{
		public.HEAD("/heartbeat", h.Heartbeat)
		public.POST("/admin/login", h.AdminLogin)
		public.GET("/blog/post/:id", h.getOne)
		public.GET("/blog/posts", h.posts)
	}

	protected := cfg.Router.Group("/")
	// protected.Use(middleware.AuthMiddleware())
	{
		blog := protected.Group("/blog/post")
		{
			blog.DELETE("/:id", h.delete)
			blog.POST("/", h.create)
			blog.PATCH("/:id", h.edit)
		}

		gallery := protected.Group("/gallery")
		{
			gallery.POST("/", h.createGallery)
			gallery.DELETE("/:id", h.deleteGallery)
		}

		admin := protected.Group("/admin")
		{
			admin.PATCH("/account", h.changePassword)
			admin.POST("/validate", h.validateJWT)
		}
	}
}
