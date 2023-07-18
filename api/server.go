package api

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/liubkkkko/firstAPI/api/controllers"
	"github.com/liubkkkko/firstAPI/api/seed"
)

var server = controllers.Server{}

func Run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	RedisDbInt, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	server.Initialize(
		os.Getenv("DB_DRIVER"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		os.Getenv("REDIS_ADDR"),
		os.Getenv("REDIS_PASSWORD"),
		RedisDbInt)

	seed.Load(server.DB)

	server.Run(":8080")

}
