package main

import (
	"log"
	"os"

	"github.com/Farzam-Haji/mock-bank/internal/database"
	"github.com/Farzam-Haji/mock-bank/internal/handlers"
	"github.com/Farzam-Haji/mock-bank/internal/repositories"
	"github.com/Farzam-Haji/mock-bank/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	config := database.Config{
		User: os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Host: os.Getenv("DATABASE_HOST"),
		Port: os.Getenv("DATABASE_PORT"),
		DBname: os.Getenv("DATABASE_DBNAME"),
	}

	db, err := database.Open(config)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	//constructors
	repoInterface := repositories.NewPostgresBankRepo(db)
	serviceInterface := services.NewBankService(repoInterface)
	bankHandler := handlers.NewBankHanler(serviceInterface)


	// Gin
	engine := gin.Default()

	api := engine.Group("/api")
	v1 := api.Group("/v1")
	v1.GET("/ping", handlers.Ping)

	v1.POST("/payments", bankHandler.CreatePayment)

	engine.Run("localhost:8081")
}