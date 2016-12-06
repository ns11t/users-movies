package datastore

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Connect(connectionStr string) (*sql.DB, error) {
	var err error
	db, err = sql.Open("postgres", connectionStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Initialize a database structure
func InitDB(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " +
		`person("id" SERIAL PRIMARY KEY,` +
		`"login" varchar(50) UNIQUE,` +
		`"password" varchar(100) NOT NULL,` +
		`"name" varchar(50) NOT NULL,` +
		`"age" smallint,` +
		`"phone_number" varchar(100))`)
	if err != nil {
		return err
	}

	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS login_idx ON person (login);")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " +
		`genre("id" SERIAL PRIMARY KEY,` +
		`"name" varchar(50) NOT NULL)`)
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " +
		`movie("id" SERIAL PRIMARY KEY,` +
		`"name" varchar(50) NOT NULL,` +
		`"year" smallint,` +
		`"genre_id" integer NOT NULL REFERENCES genre (id))`)
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " +
		`person_movie("id" SERIAL PRIMARY KEY,` +
		`"person_id" integer NOT NULL REFERENCES person (id),` +
		`"movie_id" integer NOT NULL REFERENCES movie (id))`)
	if err != nil {
		return err
	}

	return nil
}

// Drop database entities. Used by tests and migration.
func DropDB(db *sql.DB) error {
	if _, err := db.Exec("DROP TABLE IF EXISTS person_movie, movie, genre, person"); err != nil {
		return err
	}
	return nil
}
