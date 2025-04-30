package main

import (
	"log"
	"os"

	"github.com/MCantyDev/city-explorer-server/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	// Read the ENV file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	database.Connect(os.Getenv("DB_NAME"))

	// // Test Cases for the Database Querying
	// query := database.NewQueryBuilder("INSERT").Table("users").Columns("username", "password").Values(2).Build()
	// database.Execute(nil, query, "Jimbob", "saddadasdad")

	// query = database.NewQueryBuilder("INSERT").Table("users").Columns("username", "password").Values(2).Build()
	// _, err = database.Execute(nil, query, "Peter Pan", "qeqeqeafsdfwefrwef")
	// if err != nil {
	// 	fmt.Printf("Error Occured: %s", err)
	// }

	// query = database.NewQueryBuilder("SELECT").Table("users").Build()
	// var users []models.User
	// _, err = database.Execute(&users, query)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, user := range users {
	// 	fmt.Printf("ID: %d, Username: %s, Password: %s\n", user.ID, user.Username, user.Password)
	// }

	// query := database.NewQueryBuilder("UPDATE").Table("users").Columns("password").Where("id = 2")
	// database.Execute(query.Build(), "scoobydoo")

	// query := database.NewQueryBuilder("DELETE").Table("users").Where("id = 1")
	// database.Execute(query.Build())

	// router := gin.Default()
	// routes.SetupRoutes(router)
	// router.Run(":5050")
}
