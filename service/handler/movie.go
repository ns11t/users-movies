package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ns11t/users-movies/shared/datastore"
	"github.com/ns11t/users-movies/shared/model"
)

// Get movies by filters
func GetMovies(c *gin.Context) {
	params := c.Request.URL.Query()
	genresStr := params.Get("genres")
	var genres []string
	if genresStr != "" {
		genres = strings.Split(genresStr, ",")
	}

	yearStr := params.Get("year")
	var year int
	var err error
	if yearStr != "" {
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			respondWithErr(c, err, http.StatusBadRequest, model.ErrorCodeBadMovieYearFormat)
			return
		}
	}

	limit, offset, err := getLimitOffset(params)
	if err != nil {
		respondWithErr(c, err, http.StatusBadRequest, model.ErrorCodeBadLimitOffsetFormat)
		return
	}

	movies, total, err := datastore.GetMoviesByGenresYear(genres, year, limit, offset)
	if err != nil {
		respondWithErr(c, err, http.StatusInternalServerError, model.ErrorCodeInternal)
		return
	}

	respondSearchItems(c, total, offset, limit, movies)
}
