package datastore

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/ns11t/users-movies/shared/model"
)

func GetMoviesByGenresYear(genres []string, year, limit, offset int) ([]model.Movie, int, error) {
	query := "SELECT COUNT(*) OVER(), movie.*, genre.name FROM movie " +
		"JOIN genre ON movie.genre_id=genre.id"

	var where string
	var queryParams []interface{}
	var queryParamsCount int
	if len(genres) > 0 {
		queryParamsCount++
		where += fmt.Sprintf("genre.name = ANY($%d)", queryParamsCount)
		queryParams = append(queryParams, pq.Array(genres))
	}
	if year > 0 {
		if where != "" {
			where += " AND "
		}
		queryParamsCount++
		where += fmt.Sprintf("movie.year = $%d", queryParamsCount)
		queryParams = append(queryParams, year)
	}
	if len(where) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, where)
	}

	query += fmt.Sprintf(" ORDER BY movie.id LIMIT %d OFFSET %d", limit, offset)
	rows, err := db.Query(query, queryParams...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var movies []model.Movie
	var movie model.Movie
	var count int
	for rows.Next() {
		err = rows.Scan(&count, &movie.Id, &movie.Name, &movie.Year, &movie.GenreId, &movie.Genre)
		if err != nil {
			return nil, 0, err
		}
		movies = append(movies, movie)
	}
	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}
	return movies, count, nil
}

func GetMoviesByPerson(personId, limit, offset int) ([]model.Movie, int, error) {
	query := "SELECT COUNT(*) OVER(), movie.*, genre.name FROM movie" +
		" JOIN genre ON movie.genre_id=genre.id" +
		" JOIN person_movie ON movie.id=person_movie.movie_id" +
		" WHERE person_movie.person_id=$1 ORDER BY movie.id LIMIT $2 OFFSET $3"

	rows, err := db.Query(query, personId, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var movies []model.Movie
	var movie model.Movie
	var count int
	for rows.Next() {
		err = rows.Scan(&count, &movie.Id, &movie.Name, &movie.Year, &movie.GenreId, &movie.Genre)
		if err != nil {
			return nil, 0, err
		}
		movies = append(movies, movie)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	return movies, count, nil
}

// Checks if movie with given id exists
func CheckMovieExists(id int) (bool, error) {
	row := db.QueryRow("SELECT 1 FROM movie WHERE id=$1", id)
	var exists int
	err := row.Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists > 0, nil
}
