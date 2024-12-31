package srv

import (
	"bash06/tipitipi-backend/core"
	"bash06/tipitipi-backend/middleware"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) Authorize(c *gin.Context) {
	body := new(AuthBody)
	var adminUser *core.AdminUser

	if err := c.ShouldBindJSON(body); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "Malformed or invalid JSON body",
		})
		return
	}

	if err := s.Db.First(&adminUser).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			s.Log.Error("Error while checking if admin user exists", zap.Error(err))

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Something went wrong while processing your request",
			})
			return
		}
	}

	if adminUser.Username == "" {
		hs, err := s.Argon.GenerateHash([]byte(body.Password), nil)

		if err != nil {
			s.Log.Error("Error while hashing the provided default admin password", zap.Error(err))

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to hash password",
				"message": nil,
			})
			return
		}

		adminUser = &core.AdminUser{
			ID:       1,
			Username: "admin",
			Hash:     string(hs.Hash),
			Salt:     string(hs.Salt),
		}

		s.Db.Create(adminUser)
	}

	err := s.Argon.Compare([]byte(adminUser.Hash), []byte(adminUser.Salt), []byte(body.Password))

	if body.Username != adminUser.Username || err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid credentials",
			"message": nil,
		})
		return
	}

	token, err := middleware.GenerateJWT(body.Username)
	if err != nil {
		s.Log.Error("Error while generating JWT token", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate new token",
			"message": nil,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (s *Server) UpdateCredentials(c *gin.Context) {
	creds := new(AuthBody)

	if err := c.ShouldBind(&creds); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "Malformed or invalid JSON body",
		})
		return
	}

	newPass := strings.Trim(creds.Password, "")
	newUser := strings.Trim(creds.Username, "")

	if newPass == "" && newUser == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "You must provide at least one credential to be changed",
			"message": nil,
		})
		return
	}

	newCreds := &core.AdminUser{
		ID: 1,
	}

	if newPass != "" {
		hs, err := s.Argon.GenerateHash([]byte(creds.Password), nil)
		if err != nil {
			s.Log.Error("Error while hashing password", zap.Error(err))

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to hash password",
				"error":   err.Error(),
			})
			return
		}

		newCreds.Hash = string(hs.Hash)
		newCreds.Salt = string(hs.Salt)
	}

	if newUser != "" {
		newCreds.Username = newUser
	}

	if err := s.Db.Model(&core.AdminUser{}).Where("id = 1").Updates(newCreds).Error; err != nil {
		s.Log.Error("Error while updating admin credentials", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong while processing your request",
			"error":   err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) ValidateToken(c *gin.Context) {
	// The reason we can do this is that the JWT middleware
	// already handles everything for us and writing the same
	// code again would just be a waste of time
	c.Status(http.StatusOK)
}
