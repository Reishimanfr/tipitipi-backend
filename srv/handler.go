package srv

import "bash06/tipitipi-backend/middleware"

func (s *Server) InitHandler() {
	// Routes that are available without having to provide a token
	public := s.Router.Group("/")
	{
		// Admin login route
		public.POST("/admin/login", s.Authorize)

		// Returns a blog post by it's ID
		public.GET("/blog/post/:id", s.BlogGetOne)

		// Returns multiple blog posts based on parameters
		public.GET("/blog/posts", s.BlogGetBulk)

		// Serves an image based on it's filename
		public.GET("/proxy", s.Proxy)

		// Returns a gallery group by it's ID
		public.GET("/gallery/:id", s.GalleryGetOne)

		// Returns multiple gallery groups based on parameters
		public.GET("/gallery", s.GalleryGetBulk)

		// Returns the content of a page by it's name
		public.GET("/page/:name", s.PageGetOne)
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
			// Creates a new gallery group
			gallery.POST("/", s.GalleryCreateOne)

			// Adds images to a gallery group
			gallery.POST("/:id/images", s.GalleryPostBulk)

			// Deletes images from a gallery group
			gallery.DELETE("/:id/images", s.GalleryDeleteMany)

			// Deletes a gallery group
			gallery.DELETE("/:id", s.GalleryDeleteGroup)
		}

		admin := protected.Group("/admin")
		{
			// Validates a provided access token
			admin.POST("/validate", s.ValidateToken)

			// Updates admin credentials
			admin.PATCH("/update", s.UpdateCredentials)

			// Deauthorizes all active tokens
			admin.DELETE("/deauth", s.Deauth)
		}

		page := protected.Group("/page")
		{
			// Updates the content of a page
			page.PUT("/:name", s.PageUpdateOne)
		}
	}
}
