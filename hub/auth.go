package hub

import (
	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

type authError struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
}

func (a *authError) Error() string {
	return a.Reason
}

// parseAuthHeader parses AWS cognito id tokens and after validating returns username
func parseAuthHeader(r *http.Request, jwks *keyfunc.Keyfunc) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", &authError{
			Code:   http.StatusUnauthorized,
			Reason: "request lacks authorization header",
		}
	}

	authArray := strings.Split(authHeader, " ")
	if len(authArray) != 2 {
		return "", &authError{
			Code:   http.StatusUnauthorized,
			Reason: "invalid auth header provided",
		}
	}

	if bearer := authArray[0]; strings.ToLower(bearer) != "bearer" {
		return "", &authError{
			Code:   http.StatusUnauthorized,
			Reason: "unsupported authorization scheme",
		}
	}

	jwtString := authArray[1]
	token, _ := jwt.Parse(jwtString, (*jwks).Keyfunc)
	if !token.Valid {
		return "", &authError{
			Code:   http.StatusUnauthorized,
			Reason: "invalid token",
		}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", &authError{
			Code:   http.StatusUnauthorized,
			Reason: "invalid claims",
		}
	}

	username, ok := claims["preferred_username"].(string)
	if !ok || username == "" {
		return "", &authError{
			Code:   http.StatusUnauthorized,
			Reason: "username not found",
		}
	}

	return username, nil
}