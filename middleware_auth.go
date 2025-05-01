package main

import (
	"fmt"
	"net/http"
	"rssagg/internal/auth"
	"rssagg/internal/database"

	"github.com/google/uuid"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User) 

// func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
//     return func(w http.ResponseWriter, r *http.Request) {
//         apiKey, err := auth.GetAPIKey((r.Header))
//         if err != nil {
//             respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
//             return
//         }
//
//         // user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
//         // if err != nil {
//         //     respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
//         //     return
//         // }
//
//         handler(w, r, user)
//     }
// }

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
