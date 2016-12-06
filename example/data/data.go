package data

// Some methods for generating test data for movies and Genres

import (
	"database/sql"
	"math/rand"

	"github.com/ns11t/users-movies/shared/model"
)

var TestGenreNames = []string{
	"Action",
	"Adventure",
	"Comedy",
	"Crime",
	"Drama",
	"Fantasy",
	"Historical",
	"Historical fiction",
	"Horror",
	"Magical realism",
	"Mystery",
	"Paranoid",
	"Philosophical",
	"Political",
	"Romance",
	"Saga",
	"Satire",
	"Science fiction",
	"Slice of Life",
	"Speculative",
	"Thriller",
	"Urban",
	"Western",
}

var TestMovieNames = []string{
	"Allied",
	"Moana",
	"Fantastic Beasts and Where to Find Them",
	"Underworld: Blood Wars",
	"Bad Santa 2",
	"Incarnate",
}

var TestYears = []int{
	1973,
	1999,
	2004,
	2012,
	2016,
}

// Insert genres list into DB
func GenerateGenres(db *sql.DB) {
	for _, genreName := range TestGenreNames {
		_, err := db.Exec("INSERT INTO genre VALUES (default, $1)",
			genreName)
		if err != nil {
			panic(err)
		}
	}
}

// Generate and insert some records for movie table
func GenerateMovies(db *sql.DB) []model.Movie {
	movies := []model.Movie{}
	for _, movieName := range TestMovieNames {
		yearInd := rand.Intn(len(TestYears))
		genreInd := rand.Intn(len(TestGenreNames) - 1)
		var id int
		err := db.QueryRow("INSERT INTO movie VALUES (default, $1, $2, $3) RETURNING id",
			movieName, TestYears[yearInd], genreInd+1).Scan(&id)
		if err != nil {
			panic(err)
		}
		movies = append(movies, model.Movie{
			Id:      id,
			Name:    movieName,
			Year:    TestYears[yearInd],
			GenreId: genreInd + 1,
			Genre:   TestGenreNames[genreInd],
		})
	}
	return movies
}
