package routes

import (
	"github.com/gin-gonic/gin"
)

type Handler struct{}

type Config struct {
	Router *gin.Engine
}

func NewHandler(cfg *Config) {
	h := &Handler{}

	cfg.Router.Group("/").
		HEAD("/heartbeat", h.Heartbeat)

	// Blog stuff
	cfg.Router.Group("/api/blog").
		DELETE("/delete/:id", h.DeleteBlogPost).
		POST("/create", h.CreateBlogPost).
		PATCH("/edit/:id", h.EditBlogPost)
}
