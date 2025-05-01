package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"github.com/golang-jwt/jwt/v5"
)

// GetAPIKey extracts an API Key from
// the headrs of an HTTP request
// Example:
// Authorization: ApiKey {insert apikey here}
func GetAPIKey(headers http.Header) (string, error) {
    val := headers.Get("Authorization")
    if val == "" {
        return "", errors.New("no authentication info found")
    }

    vals := strings.Split(val, " ")
    if len(vals) != 2 {
        return "", errors.New("malformed auth header")
    }

    if vals[0] != "ApiKey" {
        return "", errors.New("malformed first part auth header")
    }
    return vals[1], nil
}

func VerifyJWT(r *http.Request) (string , error) {
        cookies, err := r.Cookie("jwt")
        if err != nil {
            return "", err
        }

        token, err := jwt.ParseWithClaims(cookies.Value, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("KEY")), nil})

        if err != nil {
            return "", err
        } else if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
            return claims.Issuer, nil
        } else {
            return "", errors.New("Invalid JWT")
        }

}
