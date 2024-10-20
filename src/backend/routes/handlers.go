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
			// Get info on all available gallery groups (like how many images they have)
			gallery.GET("/groups/all/info", h.getInfoAllGroups)

			// Get info on a specified gallery group
			gallery.GET("/groups/:groupId/info", h.getInfoOnGroup)

			// Get images from all gallery groups
			// TESTING NEEDED
			gallery.GET("/groups/all/images", h.getImagesAllGroups)

			// Get image from a specified gallery group
			// TESTING NEEDED
			gallery.GET("/groups/:groupId/images", h.getImagesFromGroup)

			// Initializes a new gallery group
			// TESTING NEEDED
			gallery.POST("/groups/:groupId/:name", h.createGalleryGroup)

			// Post an image to a specified group
			// TESTING NEEDED
			gallery.POST("/groups/:groupId/images", h.postImageToGroup)

			// Delete an image from a specified group
			// TESTING NEEDED
			gallery.DELETE("/groups/:groupId/images/:imageId", h.deleteImageFromGroup)

			// Delete all images from a group (without deleting the group)
			// TESTING NEEDED
			gallery.DELETE("/groups/:groupId/images", h.deleteAllFromGroup)

			// Delete an entire gallery group
			// TESTING NEEDED
			gallery.DELETE("/groups/:groupId/", h.deleteGroup)
		}

		admin := protected.Group("/admin")
		{
			admin.POST("/validate", h.validateJWT)
			admin.PATCH("/update", h.updateCreds)
		}
	}
}
