package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"rssagg/internal/database"
	"strconv"
	"time"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func createFeed (url string, apiCfg *apiConfig, w http.ResponseWriter, r *http.Request) (Feed, error) {

    rssFeed, err := urlToFeed(url)
    var finalFeed Feed
    if err != nil {
        log.Println("Error fetching feed:", err)
        return finalFeed, err
    }

    if (rssFeed.Channel.Title == "") {
        return finalFeed, errors.New("Invalid RSS URL: " + url)
    }

    description := sql.NullString{}
    if rssFeed.Channel.Description != "" {
        description.String = rssFeed.Channel.Description
        description.Valid = true
    }

    image := sql.NullString{}
    if rssFeed.Channel.Image.Url != "" {
        image.String = rssFeed.Channel.Image.Url
        image.Valid = true
    }

    feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
        ID: uuid.New(),
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
        Name: rssFeed.Channel.Title,
        Description: description,
        Url: url,
        Image: image,
    })
    if err != nil {
        return finalFeed, err
    }

    return databaseFeedToFeed(feed), err
}

func (apiCfg *apiConfig) handlerCreateFeed (w http.ResponseWriter, r *http.Request) {
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

    feed, err := createFeed(params.URL, apiCfg, w, r)
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't create feed: %v", err))
    }

    respondWithJSON(w, 201, feed)
}

func (apiCfg *apiConfig) handlerGetFeedForID(w http.ResponseWriter, r *http.Request) {
    feedID, err := uuid.Parse(chi.URLParam(r, "feedID"))

    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't get feed: %v", err))
    }

    feed, err := apiCfg.DB.GetFeedForID(r.Context(), feedID)
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't get feeds: %v", err))
        return
    }

    respondWithJSON(w, 201, databaseFeedForIDToFeed(feed))
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
    feeds, err := apiCfg.DB.GetFeeds(r.Context())
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't get feeds: %v", err))
        return
    }

    respondWithJSON(w, 201, databaseFeedsRowToFeeds(feeds))
}

func (apiCfg *apiConfig) handlerGetFeedsForCategory(w http.ResponseWriter, r *http.Request) {
    categoryID, err := uuid.Parse(chi.URLParam(r, "categoryID"))
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't get category: %v", err))
    }

    feeds, err := apiCfg.DB.GetFeedsForCategory(r.Context(), categoryID)
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't get feeds: %v", err))
        return
    }

    respondWithJSON(w, 201, databaseFeedsForCategoryRowToFeeds(feeds))
}

func (apiCfg *apiConfig) handlerGetFeedCategories(w http.ResponseWriter, r *http.Request) {
    categories, err := apiCfg.DB.GetCategories(r.Context())
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't get feeds: %v", err))
        return
    }

    respondWithJSON(w, 201, databaseCategoriesToCategories(categories))
}

func (apiCfg *apiConfig) handlerGetFollowedFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
    feeds, err := apiCfg.DB.GetFollowedFeeds(r.Context(), user.ID)
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't get feeds: %v", err))
        return
    }

    respondWithJSON(w, 201, databaseFollwedFeedsRowToFeeds(feeds))
}

func (apiCfg *apiConfig) handlerGetFeedPosts(w http.ResponseWriter, r *http.Request) {
    feedID, err := uuid.Parse(chi.URLParam(r, "feedID"))
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't get posts: %v", err))
    }

    userID, userErr := getUserID(w, r)

    order := r.URL.Query().Get("order")
    limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
    if err != nil {
        limit = 10
    }

    if order == "asc" {
        var posts []database.GetPostsForFeedAscRow
        if userErr != nil {
            posts, err = apiCfg.DB.GetPostsForFeedAsc(r.Context(), database.GetPostsForFeedAscParams{
                FeedID: feedID,
                Limit: int32(limit),
                Userid: uuid.Nil,
            })

        } else {
            posts, err = apiCfg.DB.GetPostsForFeedAsc(r.Context(), database.GetPostsForFeedAscParams{
                FeedID: feedID,
                Limit: int32(limit),
                Userid: userID,
            })
        }

        if err != nil {
            respondWithError(w,400, fmt.Sprintf("Couldn't get posts: %v", err))
            return
        }

        respondWithJSON(w, 200, databasePostsForFeedAscToPosts(posts))
    } else {
        var posts []database.GetPostsForFeedRow
        if userErr != nil {
            posts, err = apiCfg.DB.GetPostsForFeed(r.Context(), database.GetPostsForFeedParams{
                FeedID: feedID,
                Limit: int32(limit),
                Userid: uuid.Nil,
            })

        } else {
            posts, err = apiCfg.DB.GetPostsForFeed(r.Context(), database.GetPostsForFeedParams{
                FeedID: feedID,
                Limit: int32(limit),
                Userid: userID,
            })
        }

        if err != nil {
            respondWithError(w,400, fmt.Sprintf("Couldn't get posts: %v", err))
            return
        }

        respondWithJSON(w, 200, databasePostsForFeedToPosts(posts))
    }
}

func (apiCfg *apiConfig) handlerGetFeedPostsDate(w http.ResponseWriter, r *http.Request) {
    feedID, err := uuid.Parse(chi.URLParam(r, "feedID"))
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't get posts: %v", err))
    }

    date, err := time.Parse(time.RFC3339,chi.URLParam(r, "date"))
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't get posts: %v", err))
    }

    userID, userErr := getUserID(w, r)

    order := r.URL.Query().Get("order")
    limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
    if err != nil {
        limit = 10
    }

    if order == "asc" {
        var posts []database.GetPostsForFeedAfterDateRow
        if userErr != nil {
            posts, err = apiCfg.DB.GetPostsForFeedAfterDate(r.Context(), database.GetPostsForFeedAfterDateParams{
                FeedID: feedID,
                PublishedAt: date,
                Limit: int32(limit),
                Userid: uuid.Nil,
            })
        } else {
            posts, err = apiCfg.DB.GetPostsForFeedAfterDate(r.Context(), database.GetPostsForFeedAfterDateParams{
                FeedID: feedID,
                PublishedAt: date,
                Limit: int32(limit),
                Userid: userID,
            })
        }

        if err != nil {
            respondWithError(w,400, fmt.Sprintf("Couldn't get posts: %v", err))
            return
        }

        respondWithJSON(w, 201, databasePostsForFeedAfterDateToPosts(posts))
    } else {
        var posts []database.GetPostsForFeedBeforeDateRow
        if userErr != nil {
            posts, err = apiCfg.DB.GetPostsForFeedBeforeDate(r.Context(), database.GetPostsForFeedBeforeDateParams{
                FeedID: feedID,
                PublishedAt: date,
                Limit: int32(limit),
                Userid: uuid.Nil,
            })
        } else {
            posts, err = apiCfg.DB.GetPostsForFeedBeforeDate(r.Context(), database.GetPostsForFeedBeforeDateParams{
                FeedID: feedID,
                PublishedAt: date,
                Limit: int32(limit),
                Userid: userID,
            })
        }

        if err != nil {
            respondWithError(w,400, fmt.Sprintf("Couldn't get posts: %v", err))
            return
        }

        respondWithJSON(w, 201, databasePostsForFeedBeforeDateToPosts(posts))
    }
}

