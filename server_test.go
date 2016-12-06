package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ns11t/users-movies/config"
	"github.com/ns11t/users-movies/example/data"
	"github.com/ns11t/users-movies/shared"
	"github.com/ns11t/users-movies/shared/datastore"
	"github.com/ns11t/users-movies/shared/model"
	"github.com/stretchr/testify/assert"
)

var (
	db        *sql.DB
	engine    *gin.Engine
	authToken string
	movies    []model.Movie
)

type searchResponse struct {
	Success bool
	Total   int
}

func init() {
	var err error
	dbConnStr, userTokenSignKey, userTokenVerifyKey := config.GetConfigValues()
	db, err = datastore.Connect(dbConnStr)
	if err != nil {
		panic(err)
	}
	err = datastore.DropDB(db)
	if err != nil {
		panic(err)
	}
	err = datastore.InitDB(db)
	if err != nil {
		panic(err)
	}
	data.GenerateGenres(db)
	movies = data.GenerateMovies(db)

	sess := shared.NewSession(userTokenSignKey, userTokenVerifyKey)
	engine = GetEngine(sess)
}

//get movies by filters
func TestGetMovies(t *testing.T) {
	moviesResp := searchResponse{}

	//1. get all movies
	req, err := http.NewRequest("GET", "/api/v1/movies", nil)
	assert.Nil(t, err)
	resp := httptest.NewRecorder()
	engine.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)

	err = json.Unmarshal(resp.Body.Bytes(), &moviesResp)
	assert.Nil(t, err)
	assert.True(t, moviesResp.Success)
	assert.Equal(t, len(data.TestMovieNames), moviesResp.Total)

	//2. get movies by genres
	req, err = http.NewRequest("GET", "/api/v1/movies?genres="+data.TestGenreNames[0], nil)
	assert.Nil(t, err)
	resp = httptest.NewRecorder()
	engine.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)

	err = json.Unmarshal(resp.Body.Bytes(), &moviesResp)
	assert.Nil(t, err)
	assert.True(t, moviesResp.Success)

	//response should contain total amount of movies that belongs to this genre
	moviesByGenreCount := 0
	for _, movie := range movies {
		if movie.Genre == data.TestGenreNames[0] {
			moviesByGenreCount++
		}
	}
	assert.Equal(t, moviesByGenreCount, moviesResp.Total)

	//3. get movies by year
	req, err = http.NewRequest("GET", "/api/v1/movies?year="+strconv.Itoa(data.TestYears[0]), nil)
	assert.Nil(t, err)
	resp = httptest.NewRecorder()
	engine.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)

	err = json.Unmarshal(resp.Body.Bytes(), &moviesResp)
	assert.Nil(t, err)
	assert.True(t, moviesResp.Success)

	//response should contain total amount of movies that belongs to this year
	moviesByYearCount := 0
	for _, movie := range movies {
		if movie.Year == data.TestYears[0] {
			moviesByYearCount++
		}
	}
	assert.Equal(t, moviesByYearCount, moviesResp.Total)
}

// user registration/authorization
func TestRegAuth(t *testing.T) {
	//register user
	req, err := http.NewRequest("POST", "/api/v1/users", strings.NewReader(`{
		"login": "login_for_new_user_123",
		"password": "psWd45125-Akb",
		"name": "Ivan Ivanov",
		"age": 32,
		"phoneNumber": "89161233030"
	}`))
	assert.Nil(t, err)
	resp := httptest.NewRecorder()
	engine.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code)

	//authorize user
	req, err = http.NewRequest("POST", "/api/v1/session", strings.NewReader(`{
		"login": "login_for_new_user_123",
		"password": "psWd45125-Akb"
	}`))
	assert.Nil(t, err)
	resp = httptest.NewRecorder()
	engine.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code)

	//get auth token from response
	tokenResp := map[string]string{}
	err = json.Unmarshal(resp.Body.Bytes(), &tokenResp)
	assert.Nil(t, err)
	authToken = tokenResp["token"]
	assert.NotEqual(t, 0, len(authToken))
}

// adding/removing/getting user's movies
func TestUserMovies(t *testing.T) {
	//add movie to user
	req, err := http.NewRequest("POST", "/api/v1/users/movies/1", nil)
	assert.Nil(t, err)
	resp := httptest.NewRecorder()
	engine.ServeHTTP(resp, req)
	assert.Equal(t, 401, resp.Code) //expected unauthorized response code

	req.Header.Add("Authorization", "Bearer "+authToken)
	resp = httptest.NewRecorder()
	engine.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code) //expected created response code

	//get user movies
	req, err = http.NewRequest("GET", "/api/v1/users/movies", nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+authToken)
	resp = httptest.NewRecorder()
	engine.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)

	moviesResp := searchResponse{}
	err = json.Unmarshal(resp.Body.Bytes(), &moviesResp)
	assert.Nil(t, err)
	assert.True(t, moviesResp.Success)
	assert.Equal(t, 1, moviesResp.Total)

	//remove user movie
	req, err = http.NewRequest("DELETE", "/api/v1/users/movies/1", nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+authToken)
	resp = httptest.NewRecorder()
	engine.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)

	//get user movies after removing
	req, err = http.NewRequest("GET", "/api/v1/users/movies", nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+authToken)
	resp = httptest.NewRecorder()
	engine.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)

	moviesResp = searchResponse{}
	err = json.Unmarshal(resp.Body.Bytes(), &moviesResp)
	assert.Nil(t, err)
	assert.True(t, moviesResp.Success)
	assert.Equal(t, 0, moviesResp.Total) //expected that user now has 0 movies
}
