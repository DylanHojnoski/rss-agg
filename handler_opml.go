package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"rssagg/internal/database"
	"time"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerGetOPML(w http.ResponseWriter, r *http.Request) {
    feeds, err := apiCfg.DB.GetFeeds(r.Context())
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't get feeds: %v", err))
        return
    }

     opml := feedsToOPML(databaseFeedsRowToFeeds(feeds))

     bytes, err := xml.Marshal(opml)
     if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't parse opml: %v", err))
        return
     }

     respondWithFile(w, r, 200, bytes);
}

func (apiCfg *apiConfig) handlerGetOPMLUser(w http.ResponseWriter, r *http.Request, user database.User) {
    feeds, err := apiCfg.DB.GetFollowedFeeds(r.Context(), user.ID)
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't get feeds: %v", err))
        return
    }

     opml := feedsToOPML(databaseFollwedFeedsRowToFeeds(feeds))

     bytes, err := xml.Marshal(opml)
     if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't parse opml: %v", err))
        return
     }

     respondWithFile(w, r, 200, bytes);
}

func (apiCfg *apiConfig) handlerImportOPML(w http.ResponseWriter, r *http.Request) {
    var opml OPML
    decoder := xml.NewDecoder(r.Body)
    err := decoder.Decode(&opml)
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't parse opml: %v", err))
        return
    }

    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    var newFeeds []Feed

    for _, feed := range opml.Body.OutlineContainer.Outlines {
        newFeed, err := createFeed(feed.Url, apiCfg, w, r)
        log.Printf("New feed %s\n", feed.Text);

        if err == nil {
            newFeeds = append(newFeeds, newFeed)
        }  
    }

    respondWithJSON(w, 201, newFeeds)
}

func (apiCfg *apiConfig) handlerImportOPMLUser(w http.ResponseWriter, r *http.Request, user database.User) {
    var opml OPML
    decoder := xml.NewDecoder(r.Body)
    err := decoder.Decode(&opml)
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't parse opml: %v", err))
        return
    }

    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    var newFeeds []Feed
    var followedFeeds []FeedFollow

    for _, feed := range opml.Body.OutlineContainer.Outlines {

        var url = feed.Url;
        newFeed, err := createFeed(feed.Url, apiCfg, w, r)

        if err == nil {
            log.Printf("New feed %s\n", feed.Text);
            newFeeds = append(newFeeds, newFeed)
            url = newFeed.Url
        } 

        feedFollow, err := apiCfg.DB.CreateFeedFollowWithURL(r.Context(), database.CreateFeedFollowWithURLParams{
            ID: uuid.New(),
            CreatedAt: time.Now().UTC(),
            UpdatedAt: time.Now().UTC(),
            UserID: user.ID,
            Url: url,
        })

        if err != nil {
            log.Printf("Error following %s\n", feed.Text)
        } else {
            followedFeeds = append(followedFeeds, databaseFeedFollowToFeedFollow(feedFollow));
        }

    }

    respondWithJSON(w, 201, followedFeeds)
}
