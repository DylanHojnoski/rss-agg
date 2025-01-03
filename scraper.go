package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"rssagg/internal/database"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/araddon/dateparse"
	"github.com/google/uuid"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
    log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)

    ticker := time.NewTicker(timeBetweenRequest)
    for ; ; <-ticker.C {
        feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
        if err != nil {
            log.Println("error fetching feeds:", err)
            continue
        }

        wg := &sync.WaitGroup{}
        for _, feed := range feeds {
            wg.Add(1)

            go scrapeFeed(db, wg, feed)
        }
        wg.Wait()
    }
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
    defer wg.Done()

    _, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
    if err != nil {
        log.Println("Error marking feeda as fetched:", err)
        return
    }

    rssFeed, err := urlToFeed(feed.Url)
    if err != nil {
        log.Println("Error fetching feed:", err)
        return
    }

    for _, category := range rssFeed.Channel.Categories {
        if category.Text != "" {
            var dbCategory, err = db.CreateCategory(context.Background(), database.CreateCategoryParams{
                ID: uuid.New(),
                CreatedAt: time.Now().UTC(),
                UpdatedAt: time.Now().UTC(),
                Title: category.Text,
            })

            if err != nil && strings.Contains(err.Error(), "category_title_key") {
                dbCategory, err = db.GetCategoryByName(context.Background(), category.Text);
                if (err != nil) {
                    log.Println("Error getting category:", err)
                    continue;
                }
            }            

            _, err = db.CreateFeedCategory(context.Background(), database.CreateFeedCategoryParams{
                ID: uuid.New(),
                CreatedAt: time.Now().UTC(),
                UpdatedAt: time.Now().UTC(),
                FeedID: feed.ID,
                CategoryID: dbCategory.ID,
            })

            if err != nil {
                log.Println("Error creating feed category", err)
            }
        }
    }

    for _, item := range rssFeed.Channel.Item {
        description := sql.NullString{}
        if item.Description != "" {
            description.String = item.Description
            description.Valid = true
        }

        //pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
        pubAt, err := dateparse.ParseAny(item.PubDate)
        if err != nil {
            log.Printf("couldn't parse date %v with err %v", item.PubDate, err)
            continue
        }

        duration := sql.NullString{}
        if item.Duration != "" {
            if strings.Contains(item.Duration, ":") {
                duration.String = item.Duration
                duration.Valid = true
            } else {
                seconds, err := strconv.Atoi(item.Duration)
                if err == nil {
                    hours := seconds / 3600
                    minutes := (seconds % 3600) / 60
                    remaingingSeconds := seconds % 60

                    duration.String = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, remaingingSeconds)
                    duration.Valid = true
                }
            }
        }


        _, err = db.CreatePost(context.Background(), database.CreatePostParams {
            ID: uuid.New(),
            CreatedAt: time.Now().UTC(),
            UpdatedAt: time.Now().UTC(),
            Title: item.Title,
            Description: description,
            PublishedAt: pubAt,
            Audio: item.Audio.Url,
            Duration: duration,
            FeedID: feed.ID,
        })
        if err != nil {
            log.Println("failed to create post:", err)
            if strings.Contains(err.Error(), "posts_audio_key") {
                break;
            }
        }
    }
    log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
