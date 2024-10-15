package routes

import (
    ovh "bash06/strona-fundacja/src/backend/aws"
    "bash06/strona-fundacja/src/backend/core"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type Handler struct {
    Db  core.Database
    Log *zap.Logger
    A   core.Argon2idHash
    Ovh *ovh.Worker
}

type Config struct {
    Router *gin.Engine
}

func NewHandler(cfg *Config, db *core.Database, worker *ovh.Worker) {
    h := &Handler{}
    h.Db = *db
    h.Log = core.GetLogger()
    h.A = *core.NewArgon2idHash(1, 32, 64*1024, 32, 256)
    h.Ovh = worker

    public := cfg.Router.Group("/")
    {
        public.POST("/admin/login", h.auth)
        public.GET("/blog/post/:id", h.getOne)
        public.GET("/blog/posts", h.getMany)
        public.GET("/proxy", h.proxy)
    }

    protected := cfg.Router.Group("/")
    // protected.Use(middleware.AuthMiddleware())
    {
        blog := protected.Group("/blog/post")
        {
            blog.DELETE("/:id", h.deleteOne)
            blog.POST("/", h.createOne)
            blog.PATCH("/:id", h.editOne)
        }

        gallery := protected.Group("/gallery")
        {
            gallery.POST("/", h.uploadToGallery)
            gallery.DELETE("/:id", h.deleteFromGallery)
            gallery.GET("/", h.getGallery)
        }

        admin := protected.Group("/admin")
        {
            admin.POST("/validate", h.validateJWT)
            admin.PATCH("/update", h.updateCreds)
        }
    }
}