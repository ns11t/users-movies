package service

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ns11t/users-movies/shared"
	"github.com/ns11t/users-movies/shared/datastore"
	"github.com/ns11t/users-movies/shared/model"
)

// Add database session to context
func DbSessionToC(session *shared.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(shared.ContextSession, session)
		c.Next()
	}
}

// Validate user's token
func UserAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := shared.GetTokenFromRequest(c.Request)
		if err != nil {
			respondWithErr(c, err, http.StatusUnauthorized, model.ErrorCodeBadToken)
			return
		}

		sess, err := shared.CToSession(c)
		if err != nil {
			respondWithErr(c, err, http.StatusInternalServerError, model.ErrorCodeInternal)
			return
		}

		login, err := sess.ParseUserToken(token)
		if err != nil {
			respondWithErr(c, err, http.StatusUnauthorized, model.ErrorCodeBadToken)
			return
		}

		person, err := datastore.GetPersonByLogin(login)
		if err != nil {
			respondWithErr(c, err, http.StatusNotFound, model.ErrorCodeUserNotFound)
			return
		}

		c.Set(shared.ContextPerson, person)
		c.Next()
	}
}

func respondWithErr(c *gin.Context, err error, code int, errCode string) {
	log.Println(err)
	c.JSON(code, map[string]interface{}{"success": false, "msg": errCode})
	c.Abort()
}
