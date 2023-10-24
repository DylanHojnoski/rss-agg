package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rssagg/internal/database"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeed (w http.ResponseWriter, r *http.Request, user database.User) {
    type parameters struct {
        URL string `json:"url"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Error parsing JSON: %v", err))
        return
    }

    rssFeed, err := urlToFeed(params.URL)
    if err != nil {
        log.Println("Error fetching feed:", err)
        return
    }

    image:= sql.NullString{}
    if rssFeed.Channel.Image.Url!= "" {
        image.String = rssFeed.Channel.Image.Url
        image.Valid = true
    }

    feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
        ID: uuid.New(),
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
        Name: rssFeed.Channel.Title,
        Url: params.URL,
        Image: image,
        UserID: user.ID,
    })
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't create feed: %v", err))
        return
    }

    respondWithJSON(w, 201, databaseFeedToFeed(feed))
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
    feeds, err := apiCfg.DB.GetFeeds(r.Context())
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't get feeds: %v", err))
        return
    }

    respondWithJSON(w, 201, databaseFeedsToFeeds(feeds))
}

func (apiCfg *apiConfig) handlerGetFeedPosts(w http.ResponseWriter, r *http.Request) {
    feedID, err := uuid.Parse(chi.URLParam(r, "feedID"))
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't get posts: %v", err))
    }

    posts, err := apiCfg.DB.GetPostsForFeed(r.Context(), database.GetPostsForFeedParams{
        FeedID: feedID,
        Limit: 10,
    })
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't get posts: %v", err))
        return
    }

    respondWithJSON(w, 201, databasePostsToPosts(posts))
}

func (apiCfg *apiConfig) handlerGetFeedPostsBeforeDate(w http.ResponseWriter, r *http.Request) {
    feedID, err := uuid.Parse(chi.URLParam(r, "feedID"))
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't get posts: %v", err))
    }

    date, err := time.Parse(time.RFC3339,chi.URLParam(r, "date"))
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't get posts: %v", err))
    }

    posts, err := apiCfg.DB.GetPostsForFeedBeforeDate(r.Context(), database.GetPostsForFeedBeforeDateParams{
        FeedID: feedID,
        PublishedAt: date,
        Limit: 10,
    })
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't get posts: %v", err))
        return
    }

    respondWithJSON(w, 201, databasePostsToPosts(posts))
}

