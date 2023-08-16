package api

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/liubkkkko/firstAPI/api/controllers"
	"github.com/liubkkkko/firstAPI/api/seed"
	"github.com/liubkkkko/firstAPI/api/tokenstorage"
)

var server = controllers.Server{}

func Run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		log.Println("We are getting the env values")
	}

	server.Initialize(
		os.Getenv("DB_DRIVER"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"))

	tokenstorage.RedisStart(
		os.Getenv("REDIS_ADDR"),
		os.Getenv("REDIS_PASSWORD"),
		os.Getenv("REDIS_DB"))

	seed.Load(server.DB)

	server.Run(":8080")

}
