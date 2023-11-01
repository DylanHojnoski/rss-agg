package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

/*func GetEmail(header http.Header) (string, error) {
    val := headers.Get("Authorization")

    token, err := jwt.ParseWithClaims(
        tokenString,
        &claimStruct,
        func(t *jwt.Token) (interface{}, error) {
            pem, err := getGooglePublicKey(fmt.Sprintf("%s", t.Header["kid"]))
            if err != nil {
                return nil, err
            }
            
            key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
            if err != nil {
                return nil, err
            }
            
            return key, nil
        },
    )

}*/

type GoogleClaims struct {
    Email string `json:"email"`
    EmailVerified bool `json:"email_verified"`
    FirstName string `json:"given_name"`
    LastName string `json:"family_name"`
    jwt.MapClaims
}

func ValidateGoogleJWT(tokenString string) (GoogleClaims, error) {
    claimStruct :=  GoogleClaims{}

    token, err := jwt.ParseWithClaims(
        tokenString,
        &claimStruct,
        func(t *jwt.Token) (interface{}, error) {
            pem, err := getGooglePublicKey(fmt.Sprintf("%s", t.Header["kid"]))
            if err != nil {
                return nil, err
            }
            
            key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
            if err != nil {
                return nil, err
            }

            return key, nil
        },
    )

    if err != nil {
        fmt.Printf("Invalid token %v\n" , err);
        return GoogleClaims{}, err
    }

    claims, ok := token.Claims.(*GoogleClaims)
    if !ok {
        fmt.Println("Invalid Google JWT");
        return GoogleClaims{}, errors.New("Invalid Google JWT")
    }
    return *claims, nil

    /*issuer, err := claimStruct.GetIssuer()
    if err != nil || issuer != "accounts.google.com" && issuer != "https://accounts.google.com" {
        return GoogleClaims{}, errors.New("iss is invalid")
    }

    audience, err := claimStruct.GetAudience()
    if err != nil || audience !=  "client_Id" {
        return GoogleClaims{}, errors.New("aud is invalid")
    }

    expiresAt, err := claimStruct.GetExpirationTime()
    if err != nil || expiresAt.Before(time.Now().UTC()){
        return GoogleClaims{}, errors.New("JWT is expired")
    }
    
    return *claims, nil */
}

func getGooglePublicKey(keyID string) (string, error) {
    resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
    if err != nil {
        return "", err
    }

    dat, err := io.ReadAll((resp.Body))
    if err != nil {
        return "", err
    }

    myResp := map[string]string{}
    err = json.Unmarshal(dat, &myResp)
    if err != nil {
        return "", err
    }

    key, ok := myResp[keyID]
    if !ok {
        return "", errors.New("key not found")
    }
    
    return key, nil
}

func MakeJWT(claims jwt.Claims, secret string) (string, error){
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
fmt.Printf("secret %v\n", secret);
    signedToken, err := token.SignedString([]byte(secret))

    return signedToken, err 
}
