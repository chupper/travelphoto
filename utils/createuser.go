package main

import (
	"flag"
	"log"

	"github.com/chupper/travelphoto/models/user"
	"github.com/chupper/travelphoto/shared/database"
)

// creates the user name / password combination to upload pictures to the blog
func main() {

	// setup database
	database.Configure(database.DbConfig{})
	db := database.DbConnection()

	userName := flag.String("username", "", "new user name")
	password := flag.String("password", "", "password")
	flag.Parse()

	log.Println(*userName, *password)
	if *userName == "" || *password == "" {
		log.Println("Error: username and password need to be populated")
		return
	}

	err := user.CreateUser(db, *userName, *password)

	if err != nil {
		log.Fatal("Error creating account: ", err)
	}

	log.Println("User created: ", *userName)
}
