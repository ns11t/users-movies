package shared

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/ns11t/users-movies/shared/model"
)

const (
	ContextSession = "SESSION"
	ContextPerson  = "PERSON"
)

// Get session from context
func CToSession(c *gin.Context) (*Session, error) {
	sessionCtx, ok := c.Get(ContextSession)
	if !ok {
		return nil, errors.New("Cannot get session from context")
	}
	s, ok := sessionCtx.(*Session)
	if !ok {
		return nil, errors.New("Cannot get session from context: incorrect session object type")
	}
	return s, nil
}

// Get person from context
func CToPerson(c *gin.Context) (*model.Person, error) {
	personCtx, ok := c.Get(ContextPerson)
	if !ok {
		return nil, errors.New("Cannot get person from context")
	}
	person, ok := personCtx.(*model.Person)
	if !ok {
		return nil, errors.New("Cannot get person from context: incorrect person object type")
	}
	return person, nil
}
