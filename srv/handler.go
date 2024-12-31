package srv

import "bash06/tipitipi-backend/middleware"

func (s *Server) InitHandler() {
	// Routes that are available without having to provide a token
	public := s.Router.Group("/")
	{
		// Logs admin users in returning an opaque token
		public.POST("/admin/login", s.Authorize)

		// Returns a blog post by it's ID
		public.GET("/blog/post/:id", s.BlogGetOne)

		// Returns multiple blog posts based on parameters
		public.GET("/blog/posts", s.BlogGetBulk)

		// Serves an image based on it's filename
		public.GET("/proxy", s.Proxy)

		// Get EVERYTHING that's available out there (who gives a fuck?)
		public.GET("/gallery/everything", s.GetEverything)

		// Get info on all available gallery groups (like how many images they have)
		public.GET("/gallery/groups/all/info", s.GalleryGetGroupsBulk)

		// public.GET("/gallery/groups/all/images", s.GalleryGetImagesBulk)

		// Get info on a specified gallery group
		// public.GET("/gallery/groups/:groupId/info", s.GalleryGetGroupOne)

		// Get image from a specified gallery group
		// TESTING NEEDED
		// public.GET("/gallery/groups/:groupId/images", s.GalleryGetImagesOne)
	}

	// Routes that need the Authorization header with an opaque token
	protected := s.Router.Group("/")
	protected.Use(middleware.AuthMiddleware(s.Db, s.Log))
	{
		blog := protected.Group("/blog/post")
		{
			blog.DELETE("/:id", s.BlogDeleteOne)
			blog.POST("/", s.BlogCreateOne)
			blog.PATCH("/:id", s.BlogEditOne)
		}

		gallery := protected.Group("/gallery")
		{
			// Initializes a new gallery group
			// TESTING NEEDED
			gallery.POST("/groups/new/:name", s.GalleryCreateOne)

			// // Post an image to a specified group
			gallery.POST("/groups/:groupId/images", s.GalleryPostBulk)

			// // Delete an image from a specified group
			gallery.DELETE("/groups/:groupId/images/:imageId", s.GalleryDeleteOne)

			// Delete all images from a group (without deleting the group)
			gallery.DELETE("/groups/:groupId/images", s.GalleryDeleteAll)

			// Delete an entire gallery group
			gallery.DELETE("/groups/:groupId/", s.GalleryDelete)
		}

		admin := protected.Group("/admin")
		{
			// Validates a provided access token
			admin.POST("/validate", s.ValidateToken)

			// Updates admin credentials
			admin.PATCH("/update", s.UpdateCredentials)

			// Deauthorizes all active opaque tokens (including the one used in recent requests)
			admin.DELETE("/deauth", s.Deauth)
		}
	}
}
