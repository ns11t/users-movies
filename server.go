package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ns11t/users-movies/config"
	"github.com/ns11t/users-movies/service"
	"github.com/ns11t/users-movies/service/handler"
	"github.com/ns11t/users-movies/shared"
	"github.com/ns11t/users-movies/shared/datastore"
)

func main() {
	dbConnStr, userTokenSignKey, userTokenVerifyKey := config.GetConfigValues()

	db, err := datastore.Connect(dbConnStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = datastore.InitDB(db)
	if err != nil {
		panic(err)
	}

	sess := shared.NewSession(userTokenSignKey, userTokenVerifyKey)

	GetEngine(sess).Run(":8080")
}

func GetEngine(sess *shared.Session) *gin.Engine {
	r := gin.Default()
	r.Use(service.DbSessionToC(sess))

	r.StaticFile("/help/api", "docs/api.html")

	v1 := r.Group("api/v1")
	{
		v1.POST("/users", handler.RegisterUser)
		v1.POST("/session", handler.AuthorizeUser)
		v1.GET("/movies", handler.GetMovies)

		v1.Use(service.UserAuthentication())
		v1.GET("/users/movies", handler.UserGetMovies)
		v1.POST("/users/movies/:id", handler.UserAddMovie)
		v1.DELETE("/users/movies/:id", handler.UserDeleteMovie)
	}

	return r
}
