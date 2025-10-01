package handlers_test

import (
	"go-todo-api/internal/db"
	"go-todo-api/internal/handlers"
	"go-todo-api/internal/models"
	"go-todo-api/internal/testutils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetTasks(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)

	db.DB.Create(&models.Task{Title: "test", UserID: 1})

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/tasks", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handlers.GetTasks(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "test")
}

func TestGetTasksErrorNoUserID(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/tasks", func(c *gin.Context) {
		handlers.GetTasks(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateTaskSuccess(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/tasks", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handlers.CreateTask(c)
	})

	body := `{"title":"Nova Tarefa","description":"Descrição da tarefa"}`
	req, _ := http.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Nova Tarefa")
}

func TestCreateTaskUnauthorized(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)

	r := gin.Default()
	r.POST("/tasks", func(c *gin.Context) {
		handlers.CreateTask(c)
	})

	body := `{"title": "Nova tarefa", "description": "desc"}`
	req, _ := http.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "User not authenticated")
}

func TestCreateTaskInvalidJSON(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/tasks", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handlers.CreateTask(c)
	})

	req, _ := http.NewRequest(http.MethodPost, "/tasks", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid character")
}

func TestCreateTaskInternalServerError(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)

	db.DB.Migrator().DropTable(&models.Task{})

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/tasks", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handlers.CreateTask(c)
	})

	body := `{"title":"Falha","description":"Deve falhar"}`
	req, _ := http.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error creating task")
}

func TestUpdateTaskSuccess(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)
	db.DB.Create(&models.Task{Title: "Antiga", UserID: 1})

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.PUT("/tasks/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handlers.UpdateTask(c)
	})

	body := `{"title":"Updated","description":"New description","done":true}`
	req, _ := http.NewRequest(http.MethodPut, "/tasks/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Updated")
}

func TestUpdateTaskUnauthorized(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)

	r := gin.Default()
	r.PUT("/tasks/:id", func(c *gin.Context) {
		handlers.UpdateTask(c)
	})

	body := `{"title": "ed", "description": "desc", "done": true}`
	req, _ := http.NewRequest(http.MethodPut, "/tasks/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "User not authenticated")
}

func TestUpdateTaskNotFound(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.PUT("/tasks/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handlers.UpdateTask(c)
	})

	body := `{"title":"test","description":"nothing"}`
	req, _ := http.NewRequest(http.MethodPut, "/tasks/999", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Task not found")
}

func TestUpdateTaskInvalidJSON(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)
	db.DB.Create(&models.Task{Title: "Teste", UserID: 1})

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.PUT("/tasks/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handlers.UpdateTask(c)
	})

	req, _ := http.NewRequest(http.MethodPut, "/tasks/1", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid character")
}

func TestUpdateTaskInternalServerError(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)
	db.DB.Create(&models.Task{Title: "Error", UserID: 1})

	db.DB.Migrator().DropTable(&models.Task{})

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.PUT("/tasks/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handlers.UpdateTask(c)
	})

	body := `{"title":"new","description":"desc","done":true}`
	req, _ := http.NewRequest(http.MethodPut, "/tasks/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusInternalServerError || w.Code == http.StatusNotFound)
}

func TestDeleteTaskSuccess(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)
	db.DB.Create(&models.Task{Title: "Delete", UserID: 1})

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.DELETE("/tasks/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handlers.DeleteTask(c)
	})

	req, _ := http.NewRequest(http.MethodDelete, "/tasks/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Task deleted")
}

func TestDeleteTaskUnauthorized(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)

	r := gin.Default()
	r.DELETE("/tasks/:id", func(c *gin.Context) {
		handlers.DeleteTask(c)
	})

	req, _ := http.NewRequest(http.MethodDelete, "/tasks/1", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "User not authenticated")
}

func TestDeleteTaskNotFound(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.DELETE("/tasks/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handlers.DeleteTask(c)
	})

	req, _ := http.NewRequest(http.MethodDelete, "/tasks/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Task not found")
}

func TestDeleteTaskInternalServerError(t *testing.T) {
	db.DB = testutils.SetupTestDB(t)
	db.DB.Create(&models.Task{Title: "Falha", UserID: 1})

	db.DB.Migrator().DropTable(&models.Task{})

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.DELETE("/tasks/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handlers.DeleteTask(c)
	})

	req, _ := http.NewRequest(http.MethodDelete, "/tasks/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusInternalServerError || w.Code == http.StatusNotFound)
}
