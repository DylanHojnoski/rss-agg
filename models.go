package main

import (
	"encoding/json"
	"rssagg/internal/database"
	"time"

	"github.com/google/uuid"
)

type User struct {
    ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username      string    `json:"name"`
    // APIKey    string    `json:"api_key"`
}

func databaseUserToUser(dbUser database.User) User {
    return User {
        ID: dbUser.ID,
        CreatedAt: dbUser.CreatedAt,
        UpdatedAt: dbUser.UpdatedAt,
        Username: dbUser.Username,
        // APIKey: dbUser.ApiKey,
    }
}

type Feed struct {
    ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string `json:"name"`
    Description string `json:"description"`
	Url       string `json:"url"`
    Image     string `json:"image"`
    Categories []Category `json:"categories"`
}

func databaseFeedToFeed(dbFeed database.Feed) Feed{
    return Feed {
        ID: dbFeed.ID,
        CreatedAt: dbFeed.CreatedAt,
        UpdatedAt: dbFeed.UpdatedAt,
        Name: dbFeed.Name,
        Description: dbFeed.Description.String,
        Url: dbFeed.Url,
        Image: dbFeed.Image.String,
    }
}

func databaseFeedsToFeeds(dbFeed []database.Feed) []Feed{
    feeds := []Feed{}
    for _, dbFeed := range dbFeed {
        feeds = append(feeds, databaseFeedToFeed(dbFeed))
    }
    return feeds
}

type CategoryTuple struct {
    ID    uuid.UUID `json:"f1"`
    Title string `json:"f2"`
}

func databaseFeedForIDToFeed(dbFeed database.GetFeedForIDRow) Feed {
    var tuples []CategoryTuple
    err := json.Unmarshal(dbFeed.Categories, &tuples)
    if err != nil || tuples[0].Title == "" {
        return Feed {
            ID: dbFeed.ID,
            Name: dbFeed.Name,
            Description: dbFeed.Description.String,
            Url: dbFeed.Url,
            Image: dbFeed.Image.String,
            Categories: []Category{},
        }
    }

    var categories []Category

    for _, tuple := range tuples {
        categories = append(categories, Category{
            ID: tuple.ID,
            Title: tuple.Title,
        })

    }

    return Feed {
        ID: dbFeed.ID,
        Name: dbFeed.Name,
        Description: dbFeed.Description.String,
        Url: dbFeed.Url,
        Image: dbFeed.Image.String,
        Categories: categories,
    }
}

func databaseFeedRowToFeed(dbFeed database.GetFeedsRow) Feed {
    var tuples []CategoryTuple
    err := json.Unmarshal(dbFeed.Categories, &tuples)
    if err != nil || tuples[0].Title == "" {
        return Feed {
            ID: dbFeed.ID,
            Name: dbFeed.Name,
            Description: dbFeed.Description.String,
            Url: dbFeed.Url,
            Image: dbFeed.Image.String,
            Categories: []Category{},
        }
    }

    var categories []Category

    for _, tuple := range tuples {
        categories = append(categories, Category{
            ID: tuple.ID,
            Title: tuple.Title,
        })

    }

    return Feed {
        ID: dbFeed.ID,
        Name: dbFeed.Name,
        Description: dbFeed.Description.String,
        Url: dbFeed.Url,
        Image: dbFeed.Image.String,
        Categories: categories,
    }
}

func databaseFeedsForCategoryRowToFeeds(dbFeed []database.GetFeedsForCategoryRow) []Feed{
    feeds := []Feed{}
    for _, dbFeed := range dbFeed {
        feeds = append(feeds, databaseFeedForCategoryRowToFeed(dbFeed))
    }
    return feeds
}

func databaseFeedForCategoryRowToFeed(dbFeed database.GetFeedsForCategoryRow) Feed {
    var tuples []CategoryTuple
    err := json.Unmarshal(dbFeed.Categories, &tuples)
    if err != nil || tuples[0].Title == "" {
        return Feed {
            ID: dbFeed.ID,
            Name: dbFeed.Name,
            Description: dbFeed.Description.String,
            Url: dbFeed.Url,
            Image: dbFeed.Image.String,
            Categories: []Category{},
        }
    }

    var categories []Category

    for _, tuple := range tuples {
        categories = append(categories, Category{
            ID: tuple.ID,
            Title: tuple.Title,
        })

    }

    return Feed {
        ID: dbFeed.ID,
        Name: dbFeed.Name,
        Description: dbFeed.Description.String,
        Url: dbFeed.Url,
        Image: dbFeed.Image.String,
        Categories: categories,
    }
}

func databaseFeedsRowToFeeds(dbFeed []database.GetFeedsRow) []Feed{
    feeds := []Feed{}
    for _, dbFeed := range dbFeed {
        feeds = append(feeds, databaseFeedRowToFeed(dbFeed))
    }
    return feeds
}

func databaseFollwedFeedRowToFeed(dbFeed database.GetFollowedFeedsRow) Feed {
    var tuples []CategoryTuple
    err := json.Unmarshal(dbFeed.Categories, &tuples)
    if err != nil || tuples[0].Title == "" {
        return Feed {
            ID: dbFeed.ID,
            Name: dbFeed.Name,
            Description: dbFeed.Description.String,
            Url: dbFeed.Url,
            Image: dbFeed.Image.String,
            Categories: []Category{},
        }
    }

    var categories []Category

    for _, tuple := range tuples {
        categories = append(categories, Category{
            ID: tuple.ID,
            Title: tuple.Title,
        })

    }

    return Feed {
        ID: dbFeed.ID,
        Name: dbFeed.Name,
        Description: dbFeed.Description.String,
        Url: dbFeed.Url,
        Image: dbFeed.Image.String,
        Categories: categories,
    }
}

func databaseFollwedFeedsRowToFeeds(dbFeed []database.GetFollowedFeedsRow) []Feed{
    feeds := []Feed{}
    for _, dbFeed := range dbFeed {
        feeds = append(feeds, databaseFollwedFeedRowToFeed(dbFeed))
    }
    return feeds
}

type FeedFollow struct {
    ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func databaseFeedFollowToFeedFollow(dbFeedFollow database.FeedFollow) FeedFollow {
    return FeedFollow {
        ID: dbFeedFollow.ID,
        CreatedAt: dbFeedFollow.CreatedAt,
        UpdatedAt: dbFeedFollow.UpdatedAt,
        UserID: dbFeedFollow.UserID,
        FeedID: dbFeedFollow.FeedID,
    }
}

func databaseFeedFollowsToFeedFollows(dbFeedFollows []database.FeedFollow) []FeedFollow{
    feedFollows := []FeedFollow{}
    for _, dbFeedFollow := range dbFeedFollows {
        feedFollows = append(feedFollows, databaseFeedFollowToFeedFollow(dbFeedFollow))
    }
    return feedFollows
}

type Post struct {
    ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string `json:"title"`
	Description *string `json:"description"`
	PublishedAt time.Time `json:"published_at"`
    Audio       string `json:"audio"`
    Duration    string `json:"duration"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func databasePostToPost(dbPost database.Post) Post {
    var description *string
    if dbPost.Description.Valid {
        description = &dbPost.Description.String
    }

    return Post{
        ID: dbPost.ID,
        CreatedAt: dbPost.CreatedAt,
        UpdatedAt: dbPost.UpdatedAt,
        Title: dbPost.Title,
        Description: description,
        PublishedAt: dbPost.PublishedAt,
        Audio: dbPost.Audio,
        Duration: dbPost.Duration.String,
        FeedID: dbPost.FeedID,
    }
}

func databasePostsToPosts(dbPosts []database.Post) []Post {
    posts := []Post{}
    for _, dbPost := range dbPosts {
        posts = append(posts, databasePostToPost(dbPost))
    }
    return posts
}

type Category struct {
    ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Title     string    `json:"title"`
}

func databaseCategoryToCategory(dbCategory database.Category) Category {
    return Category {
        ID: dbCategory.ID,
        CreatedAt: dbCategory.CreatedAt,
        UpdatedAt: dbCategory.UpdatedAt,
        Title: dbCategory.Title,
    }
}

func databaseCategoriesToCategories(dbCategories []database.Category) []Category {
    categories := []Category{}
    for _, dbCategory := range dbCategories {
        categories = append(categories, databaseCategoryToCategory(dbCategory))
    }
    return categories
}

