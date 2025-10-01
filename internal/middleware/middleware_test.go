package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-todo-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(JWTAuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no userID"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"userID": userID})
	})

	// Generate valid token
	token, err := utils.GenerateAccessToken(1)
	assert.NoError(t, err)

	// Request with valid token
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Request without token
	reqNoToken := httptest.NewRequest("GET", "/protected", nil)
	wNoToken := httptest.NewRecorder()
	r.ServeHTTP(wNoToken, reqNoToken)

	assert.Equal(t, http.StatusUnauthorized, wNoToken.Code)

	// Request with invalid token
	reqInvalid := httptest.NewRequest("GET", "/protected", nil)
	reqInvalid.Header.Set("Authorization", "Bearer invalidtoken")
	wInvalid := httptest.NewRecorder()
	r.ServeHTTP(wInvalid, reqInvalid)

	assert.Equal(t, http.StatusUnauthorized, wInvalid.Code)

	// Invalid Bearer
	reqNoBearer := httptest.NewRequest("GET", "/protected", nil)
	reqNoBearer.Header.Set("Authorization", "invalid")
	wNoBearer := httptest.NewRecorder()
	r.ServeHTTP(wNoBearer, reqNoBearer)
	assert.Equal(t, http.StatusUnauthorized, wNoBearer.Code)
}
