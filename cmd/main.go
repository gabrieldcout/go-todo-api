package main

import (
	"go-todo-api/internal/db"
	"go-todo-api/internal/routes"
	"log"
)

func main() {
	db.ConnectDatabase()

	r := routes.SetupRoutes()

	log.Println("Servidor rodando em http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
