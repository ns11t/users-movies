package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ns11t/users-movies/shared"
	"github.com/ns11t/users-movies/shared/datastore"
	"github.com/ns11t/users-movies/shared/model"
	"golang.org/x/crypto/bcrypt"
)

// Register new user
func RegisterUser(c *gin.Context) {
	var person model.Person
	err := c.BindJSON(&person)
	if err != nil {
		respondWithErr(c, err, http.StatusBadRequest, model.ErrorCodeBadUserModel)
		return
	}

	err = person.Validate()
	if err != nil {
		respondWithErr(c, err, http.StatusBadRequest, model.ErrorCodeBadUserModel)
		return
	}

	exists, err := datastore.CheckPersonExists(person.Login)
	if err != nil {
		respondWithErr(c, err, http.StatusInternalServerError, model.ErrorCodeInternal)
		return
	}
	if exists {
		respondWithErr(c, err, http.StatusConflict, model.ErrorCodeUserAlreadyExists)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(person.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithErr(c, err, http.StatusInternalServerError, model.ErrorCodeInternal)
		return
	}

	person.Password = string(hashedPassword)
	_, err = datastore.InsertPerson(person)
	if err != nil {
		respondWithErr(c, err, http.StatusInternalServerError, model.ErrorCodeInternal)
		return
	}

	respondWithSuccess(c, http.StatusCreated)
}

// Get all movies that belongs to current user
func UserGetMovies(c *gin.Context) {
	person, err := shared.CToPerson(c)
	if err != nil {
		respondWithErr(c, err, http.StatusInternalServerError, model.ErrorCodeInternal)
		return
	}

	params := c.Request.URL.Query()
	limit, offset, err := getLimitOffset(params)
	if err != nil {
		respondWithErr(c, err, http.StatusBadRequest, model.ErrorCodeBadLimitOffsetFormat)
		return
	}

	movies, total, err := datastore.GetMoviesByPerson(person.Id, limit, offset)
	if err != nil {
		respondWithErr(c, err, http.StatusInternalServerError, model.ErrorCodeInternal)
		return
	}

	respondSearchItems(c, total, offset, limit, movies)
}

// Add movie to user's list
func UserAddMovie(c *gin.Context) {
	person, err := shared.CToPerson(c)
	if err != nil {
		respondWithErr(c, err, http.StatusInternalServerError, model.ErrorCodeInternal)
		return
	}

	movieId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondWithErr(c, err, http.StatusBadRequest, model.ErrorCodeBadMovieIdFormat)
		return
	}

	exists, err := datastore.CheckMovieExists(movieId)
	if err != nil {
		respondWithErr(c, err, http.StatusInternalServerError, model.ErrorCodeInternal)
		return
	}
	if !exists {
		respondWithErr(c, err, http.StatusNotFound, model.ErrorCodeMovieNotFound)
		return
	}

	exists, err = datastore.CheckPersonMovieExists(person.Id, movieId)
	if err != nil {
		respondWithErr(c, err, http.StatusInternalServerError, model.ErrorCodeInternal)
		return
	}
	if exists {
		respondWithErr(c, err, http.StatusConflict, model.ErrorCodePersonMovieAlreadyExists)
		return
	}

	_, err = datastore.InsertPersonMovie(person.Id, movieId)
	if err != nil {
		respondWithErr(c, err, http.StatusInternalServerError, model.ErrorCodeInternal)
		return
	}

	respondWithSuccess(c, http.StatusCreated)
}

// Delete movie from user's list
func UserDeleteMovie(c *gin.Context) {
	person, err := shared.CToPerson(c)
	if err != nil {
		respondWithErr(c, err, http.StatusInternalServerError, model.ErrorCodeInternal)
		return
	}

	movieIdStr := c.Param("id")
	movieId, err := strconv.Atoi(movieIdStr)
	if err != nil {
		respondWithErr(c, err, http.StatusBadRequest, model.ErrorCodeBadMovieIdFormat)
		return
	}

	exists, err := datastore.CheckPersonMovieExists(person.Id, movieId)
	if err != nil {
		respondWithErr(c, err, http.StatusInternalServerError, model.ErrorCodeInternal)
		return
	}
	if !exists {
		respondWithErr(c, err, http.StatusNotFound, model.ErrorCodeUserMovieNotFound)
		return
	}

	_, err = datastore.DeletePersonMovie(person.Id, movieId)
	if err != nil {
		respondWithErr(c, err, http.StatusInternalServerError, model.ErrorCodeInternal)
		return
	}
	respondWithSuccess(c, http.StatusOK)
}
