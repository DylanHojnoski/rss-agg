package main

import (
	"errors"
	"fmt"
	"net/http"
	"rssagg/internal/auth"
	"rssagg/internal/database"

	"github.com/google/uuid"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User) 

func (apiCfG *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        userID, err := auth.VerifyJWT(r)
        if err != nil {
            respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
            return
        }

        parsedID, err := uuid.Parse(userID)
        if err != nil {
            respondWithError(w, 400, fmt.Sprintf("Couldn't parse ID", err))
            return
        }

        user, err := apiCfG.DB.GetUserByID(r.Context(), parsedID)
        if err != nil {
            respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
            return
        }

        handler(w, r, user)
    }

}

func getUserID(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {

    userID, err := auth.VerifyJWT(r)
    if err != nil {
        return uuid.UUID{}, errors.New("No user id")
    }

    id, err := uuid.Parse(userID)
    if err != nil {
        return uuid.UUID{}, errors.New("No user id")
    }

    return id, nil
}
