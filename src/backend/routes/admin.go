package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"bash06/strona-fundacja/src/backend/middleware"
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

const (
	err_invalid_json             = "Malformed or invalid JSON body"
	err_hashing_failure          = "Failed to hash password"
	err_auth_invalid_credentials = "Invalid credentials"
	err_auth_jwt_failure         = "Failed to generate new token"
	err_auth_no_creds            = "You must provide at least one credential to be changed"
	ok_auth_update_success       = "Admin user credentials updated successfully"
)

// /admin/login
func (h *Handler) auth(c *gin.Context) {
	body := new(AuthBody)
	adminUser := new(core.AdminUser)

	if err := c.ShouldBindJSON(body); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": err_invalid_json,
		})
		return
	}

	if err := h.Db.First(&adminUser, "id = 1").Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			h.Log.Error("Error while checking if admin user exists", zap.Error(err))

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": errSqlQuery,
			})
			return
		}
	}

	if adminUser == nil {
		hs, err := h.A.GenerateHash([]byte(body.Password), nil)

		if err != nil {
			h.Log.Error("Error while hashing the provided default admin password", zap.Error(err))

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   err_hashing_failure,
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

		h.Db.Create(adminUser)
	}

	err := h.A.Compare([]byte(adminUser.Hash), []byte(adminUser.Salt), []byte(body.Password))

	if body.Username != adminUser.Username || err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   err_auth_invalid_credentials,
			"message": nil,
		})
		return
	}

	token, err := middleware.GenerateJWT(body.Username)
	if err != nil {
		h.Log.Error("Error while generating JWT token", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err_auth_jwt_failure,
			"message": nil,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (h *Handler) updateCreds(c *gin.Context) {
	creds := new(AuthBody)

	if err := c.ShouldBind(&creds); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": err_invalid_json,
		})
		return
	}

	newPass := strings.Trim(creds.Password, "")
	newUser := strings.Trim(creds.Username, "")

	if newPass == "" && newUser == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err_auth_no_creds,
			"message": nil,
		})
		return
	}

	newCreds := &core.AdminUser{
		ID: 1,
	}

	if newPass != "" {
		hs, err := h.A.GenerateHash([]byte(creds.Password), nil)
		if err != nil {
			h.Log.Error("Error while hashing password", zap.Error(err))

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err_hashing_failure,
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

	if err := h.Db.Model(&core.AdminUser{}).Where("id = 1").Updates(newCreds).Error; err != nil {
		h.Log.Error("Error while updating admin credentials", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": errSqlQuery,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": ok_auth_update_success,
		"error":   nil,
	})
}

func (h *Handler) validateJWT(c *gin.Context) {
	// The reason we can do this is that the JWT middleware
	// already handles everything for us and writing the same
	// code again would just be a waste of time
	c.JSON(http.StatusOK, gin.H{
		"message": "JWT token is valid",
	})
}
