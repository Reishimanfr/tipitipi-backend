package routes

import (
	"bash06/strona-fundacja/src/backend/core"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Db core.Database
}

type Config struct {
	Router *gin.Engine
}

func NewHandler(cfg *Config, db *core.Database) {
	h := &Handler{}
	h.Db = *db

	cfg.Router.Group("/").
		HEAD("/heartbeat", h.Heartbeat)

	// Blog stuff
	cfg.Router.Group("/api/blog").
		DELETE("/delete/:id", h.DeleteBlogPost).
		POST("/create", h.CreateBlogPost).
		PATCH("/edit/:id", h.EditBlogPost).
		GET("/getPostById/:id", h.GetPostById)
	// GET("/getPostByData", h.GetPostByData)
}
