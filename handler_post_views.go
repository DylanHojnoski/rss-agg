package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rssagg/internal/database"
	"time"
	"github.com/google/uuid"
)


func (apiCfg *apiConfig) handlerCreatePostView (w http.ResponseWriter, r *http.Request, user database.User) {
    type parameters struct {
        PostID uuid.UUID `json:"id"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Error parsing JSON: %v", err))
        return
    }

    postView, err := apiCfg.DB.CreatePostView(r.Context(), database.CreatePostViewParams{
        ID: uuid.New(),
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
        UserID: user.ID,
        PostID: params.PostID,
    })

    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't create feed follow: %v", err))
        return
    }

    respondWithJSON(w, 201, databasePostViewToPostView(postView))
}
