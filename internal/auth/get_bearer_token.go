package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	const BearerPrefix = "Bearer "
	
	bearer_token := headers.Get("Authorization")

	if bearer_token == "" {
		return "", fmt.Errorf("no authorization header present")
	}

	if !strings.HasPrefix(bearer_token, BearerPrefix) {
		return "", fmt.Errorf("authorization header does not start with correct prefix")
	}

	token := strings.TrimPrefix(bearer_token, BearerPrefix)
	token = strings.TrimSpace(token)
	
	if token == "" {
		return "", fmt.Errorf("no token after bearer prefix")
	}

	return token, nil	
}