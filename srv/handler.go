package srv

import (
	"bash06/tipitipi-backend/flags"
	"bash06/tipitipi-backend/middleware"
	"time"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
)

var (
	cacheTime = time.Minute * time.Duration(*flags.CacheLifetime)
)

func (s *Server) InitHandler() {
	s.Cache = persist.NewMemoryStore(cacheTime)

	// Routes that are available without having to provide a token
	public := s.Router.Group("/")
	{
		// Admin login route
		public.POST("/admin/login", s.Authorize)

		// Returns a blog post by it's ID
		public.GET("/blog/post/:id", s.BlogGetOne, cache.CacheByRequestURI(s.Cache, cacheTime))

		// Returns multiple blog posts based on parameters
		public.GET("/blog/posts", s.BlogGetBulk, cache.CacheByRequestURI(s.Cache, cacheTime))

		// Serves an image based on it's filename
		public.GET("/proxy", s.Proxy, cache.CacheByRequestURI(s.Cache, cacheTime))

		// Returns a gallery group by it's ID
		public.GET("/gallery/:id", s.GalleryGetOne, cache.CacheByRequestURI(s.Cache, cacheTime))

		// Returns multiple gallery groups based on parameters
		public.GET("/gallery", s.GalleryGetBulk, cache.CacheByRequestURI(s.Cache, cacheTime))

		// Returns the content of a page by it's name
		public.GET("/page/:name", s.PageGetOne, cache.CacheByRequestURI(s.Cache, cacheTime))
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
