package datastore

import (
	"database/sql"
)

// Checks if movie with given id currently belongs to person
func CheckPersonMovieExists(personId, movieId int) (bool, error) {
	row := db.QueryRow("SELECT 1 FROM person_movie WHERE person_id=$1 AND movie_id=$2", personId, movieId)
	var exists int
	err := row.Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists > 0, nil
}

func InsertPersonMovie(personId, movieId int) (sql.Result, error) {
	return db.Exec("INSERT INTO person_movie VALUES (default, $1, $2)", personId, movieId)
}

func DeletePersonMovie(personId, movieId int) (sql.Result, error) {
	return db.Exec("DELETE FROM person_movie WHERE person_id=$1 AND movie_id=$2", personId, movieId)
}
