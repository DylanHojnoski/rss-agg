package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
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
        if err != nil {

        } else {
            newFeeds = append(newFeeds, newFeed)
        }
    }

    respondWithJSON(w, 201, newFeeds)
}
