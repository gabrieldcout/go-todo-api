package routes

import (
	"go-todo-api/internal/handlers"
	"go-todo-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.POST("/signup", handlers.Signup)
	r.POST("/login", handlers.Login)
	r.POST("/refresh", handlers.RefreshToken)
	r.POST("/logout", handlers.Logout)

	// Rotas protegidas
	auth := r.Group("/")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		auth.GET("/tasks", handlers.GetTasks)
		auth.POST("/tasks", handlers.CreateTask)
		auth.PUT("/tasks/:id", handlers.UpdateTask)
		auth.DELETE("/tasks/:id", handlers.DeleteTask)
	}

	return r
}
