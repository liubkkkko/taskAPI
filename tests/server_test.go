package tests

import (
	"github.com/joho/godotenv"
	"github.com/liubkkkko/firstAPI/api/controllers"
	"github.com/liubkkkko/firstAPI/api/seed"
	"github.com/liubkkkko/firstAPI/api/tokenstorage"
	"log"
	"os"
	"testing"
)

var testServer = controllers.Server{}

func TestRun(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		log.Println("We are getting the env values")
	}

	testServer.Initialize(
		os.Getenv("TEST_DB_DRIVER"),
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_NAME"))

	tokenstorage.RedisStart(
		os.Getenv("REDIS_ADDR"),
		os.Getenv("REDIS_PASSWORD"),
		os.Getenv("REDIS_DB"))

	seed.Load(testServer.DB)

	testServer.Run(":8080")

}

//
//func Database() {
//	//connect to postgres
//	var err error
//	dbUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
//	server.DB, err = gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
//	if err != nil {
//		fmt.Printf("Cannot connect to %s database", os.Getenv("TestDbDriver"))
//		log.Fatal("This is the error:", err)
//	} else {
//		fmt.Printf("We are connected to the %s database", os.Getenv("TestDbDriver"))
//	}
//
//}
