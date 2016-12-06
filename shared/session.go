package shared

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Session contains keys for signing and validating tokens
type Session struct {
	userTokenSignKey   []byte
	userTokenVerifyKey []byte
}

func NewSession(userTokenSignKey, userTokenVerifyKey []byte) *Session {
	return &Session{
		userTokenSignKey:   userTokenSignKey,
		userTokenVerifyKey: userTokenVerifyKey,
	}
}

// Create a new token for user
func (s *Session) GenerateUserToken(personId string) (string, error) {
	jwtToken := jwt.New(jwt.GetSigningMethod("RS256"))
	jwtToken.Claims["userid"] = personId
	jwtToken.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	return jwtToken.SignedString(s.userTokenSignKey)
}

// Get userid from token
func (s *Session) ParseUserToken(token string) (string, error) {
	var jwtToken, jwtErr = jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return s.userTokenVerifyKey, nil
	})

	if jwtErr != nil {
		return "", jwtErr
	}

	if !jwtToken.Valid {
		return "", errors.New("invalid token")
	}

	login, ok := jwtToken.Claims["userid"].(string)
	if !ok {
		return "", errors.New("cannot parse token: incorrect userid")
	}

	return login, nil
}

// Get token from http Request Authorization header.
func GetTokenFromRequest(r *http.Request) (string, error) {
	var token string
	_, err := fmt.Sscanf(r.Header.Get("Authorization"), "Bearer %s", &token)
	if err != nil || len(token) == 0 {
		return "", errors.New("token not found")
	}
	return token, nil
}
