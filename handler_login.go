package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rssagg/internal/auth"
	"rssagg/internal/database"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()


	// parse the GoogleJWT that was POSTed from the front-end
	type parameters struct {
        Params struct {
            GoogleJWT string `json:"googleJWT"`
        } `json:"params"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Couldn't decode parameters")
		return
	}

	// Validate the JWT is valid
	claims, err := auth.ValidateGoogleJWT(params.Params.GoogleJWT)
	if err != nil {
		respondWithError(w, 403, "Invalid google auth")
		return
	}

  user, err :=  cfg.DB.GetUserByEmail(r.Context(), claims.Email);
    if err != nil {
        if  claims.Email != "" {
            user, err = cfg.DB.CreateUser(r.Context(), database.CreateUserParams {
                ID: uuid.New(),
                CreatedAt: time.Now().UTC(),
                UpdatedAt: time.Now().UTC(),
                Name: claims.FirstName,
                Email: claims.Email,
            })
            if err != nil {
                respondWithError(w,400, fmt.Sprintf("Couldn't create user: %v", err))
                return
            }
        }
        respondWithError(w,400, fmt.Sprintf("Couldn't get user email: %v %v", claims.Email, err))
        return
    }

    fmt.Printf("user %v", user);

	// create a JWT for OUR app and give it back to the client for future requests
	tokenString, err := auth.MakeJWT(claims, cfg.JWTSecret)
	if err != nil {
        fmt.Printf("err %v", err);
		respondWithError(w, 500, "Couldn't make authentication token")
		return
	}

	respondWithJSON(w, 200, struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	})
}
