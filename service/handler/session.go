package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ns11t/users-movies/shared"
	"github.com/ns11t/users-movies/shared/datastore"
	"github.com/ns11t/users-movies/shared/model"
	"golang.org/x/crypto/bcrypt"
)

// Generate new token for user
func AuthorizeUser(c *gin.Context) {
	var person model.Person
	err := c.BindJSON(&person)
	if err != nil {
		respondWithErr(c, err, http.StatusBadRequest, model.ErrorCodeBadUserModel)
		return
	}
	err = person.ValidateLoginPassword()
	if err != nil {
		respondWithErr(c, err, http.StatusBadRequest, model.ErrorCodeBadUserModel)
		return
	}

	personStored, err := datastore.GetPersonPasswordByLogin(person.Login)
	if err != nil {
		respondWithErr(c, err, http.StatusNotFound, model.ErrorCodeUserNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(personStored.Password), []byte(person.Password))
	if err != nil {
		respondWithErr(c, err, http.StatusUnauthorized, model.ErrorCodeUserIncorrectPassword)
		return
	}

	sess, err := shared.CToSession(c)
	if err != nil {
		respondWithErr(c, err, http.StatusInternalServerError, model.ErrorCodeInternal)
		return
	}

	token, err := sess.GenerateUserToken(person.Login)
	if err != nil {
		respondWithErr(c, err, http.StatusInternalServerError, model.ErrorCodeInternal)
		return
	}
	c.JSON(http.StatusCreated, map[string]string{"token": token})
}
