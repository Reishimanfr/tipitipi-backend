package srv

import "bash06/strona-fundacja/src/backend/middleware"

func (s *Server) InitHandler() {
	public := s.Router.Group("/")
	{
		public.POST("/admin/login", s.Authorize)
		public.GET("/blog/post/:id", s.BlogGetOne)
		public.GET("/blog/posts", s.BlogGetBulk)
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

	protected := s.Router.Group("/")
	protected.Use(middleware.AuthMiddleware())
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
			admin.POST("/validate", s.ValidateToken)
			admin.PATCH("/update", s.UpdateCredentials)
		}
	}
}
