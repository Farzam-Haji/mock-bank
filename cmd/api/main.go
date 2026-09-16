package main

import (
	"github.com/Farzam-Haji/mock-bank/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {


	// Gin
	engine := gin.Default()
	api := engine.Group("/api")
	v1 := api.Group("/v1")
	v1.GET("/ping", handlers.Ping)
	engine.Run("localhost:8081")
}