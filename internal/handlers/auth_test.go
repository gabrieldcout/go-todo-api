package handlers

import (
	"fmt"
	"go-todo-api/internal/db"
	"go-todo-api/internal/models"
	"go-todo-api/internal/testutils"
	"go-todo-api/internal/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSignupSuccess(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)

	r := gin.Default()
	r.POST("/signup", Signup)

	body := `{"email":"teste@example.com","password":"123456"}`
	req, _ := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println("Status:", w.Code)
	fmt.Println("Body:", w.Body.String())

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "User created successfully")
}

func TestSignupInvalidJSON(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)

	r := gin.Default()
	r.POST("/signup", Signup)

	body := `{"email":1, "password":true}`
	req, _ := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid data")
}

func TestSignupDuplicateEmail(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)

	db.DB.Create(&models.User{Email: "teste@example.com", PasswordHash: "senha"})

	r := gin.Default()
	r.POST("/signup", Signup)

	body := `{"email":"teste@example.com","password":"nova"}`
	req, _ := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error creating user")
}

func TestLoginSuccess(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)

	hashed, _ := utils.HashPassword("123456")
	db.DB.Create(&models.User{Email: "teste@example.com", PasswordHash: hashed})

	r := gin.Default()
	r.POST("/login", Login)

	body := `{"email":"teste@example.com","password":"123456"}`
	req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
	cookie := w.Result().Cookies()
	assert.NotEmpty(t, cookie)
}

func TestLoginInvalidJSON(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)

	r := gin.Default()
	r.POST("/login", Login)

	body := `{"email":123,"password":false}`
	req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginUserNotFound(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)

	r := gin.Default()
	r.POST("/login", Login)

	body := `{"email":"naoexiste@example.com","password":"123"}`
	req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid username or password")
}

func TestLoginWrongPassword(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)

	hashed, _ := utils.HashPassword("correta")
	db.DB.Create(&models.User{Email: "teste@example.com", PasswordHash: hashed})

	r := gin.Default()
	r.POST("/login", Login)

	body := `{"email":"teste@example.com","password":"errada"}`
	req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRefreshTokenSuccess(t *testing.T) {
	token, _ := utils.GenerateRefreshToken(1)

	r := gin.Default()
	r.GET("/refresh", RefreshToken)

	req, _ := http.NewRequest(http.MethodGet, "/refresh", nil)
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: token,
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "access_token")
}

func TestRefreshTokenMissingCookie(t *testing.T) {
	r := gin.Default()
	r.GET("/refresh", RefreshToken)

	req, _ := http.NewRequest(http.MethodGet, "/refresh", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Refresh token missing")
}

func TestRefreshTokenInvalidToken(t *testing.T) {
	r := gin.Default()
	r.GET("/refresh", RefreshToken)

	req, _ := http.NewRequest(http.MethodGet, "/refresh", nil)
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: "token_invalido",
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid refresh token")
}

func TestLogoutSuccess(t *testing.T) {
	r := gin.Default()
	r.POST("/logout", Logout)

	req, _ := http.NewRequest(http.MethodPost, "/logout", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Logout successful")
}