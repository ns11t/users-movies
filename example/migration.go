package main

// Script for generating test data

import (
	"log"

	"github.com/ns11t/users-movies/config"
	"github.com/ns11t/users-movies/example/data"
	"github.com/ns11t/users-movies/shared/datastore"
)

func main() {
	dbConnStr, _, _ := config.GetConfigValues()
	db, err := datastore.Connect(dbConnStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = datastore.DropDB(db)
	if err != nil {
		panic(err)
	}
	err = datastore.InitDB(db)
	if err != nil {
		panic(err)
	}
	data.GenerateGenres(db)
	data.GenerateMovies(db)

	log.Println("test data about movies and genres were successfully generated")
}
