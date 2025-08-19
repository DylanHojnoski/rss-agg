package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"rssagg/internal/database"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

type apiConfig struct {
    DB *database.Queries
}

func main() {
     
    godotenv.Load(".env")

    portString := os.Getenv("PORT")
    if portString == "" {
        log.Fatal("PORT is not found in the environment")
    }

    dbURL := os.Getenv("DB_URL")
    if dbURL == "" {
        log.Fatal("DB_URL is not found in the environment")
    }

    conn, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatal("Cannot connect to database:", err)
    }

	if err := goose.Up(conn, "sql/schema"); err != nil {
		log.Fatal(err)
	}

    
    db := database.New(conn)
    apiCfg := apiConfig{
        DB: db,
    }

	scrapeConcurrencyString := os.Getenv("SCRAPE_CONCURRENCY") 
	scrapeConcurrency, err := strconv.Atoi(scrapeConcurrencyString)
	if err != nil || scrapeConcurrency <= 0 {
		scrapeConcurrency = 10
	}
	log.Printf("Scrape Concurrency: %d\n", scrapeConcurrency)

	scrapeFrequencyString := os.Getenv("SCRAPE_FREQUENCY") 
	scrapeFrequency, err := strconv.Atoi(scrapeFrequencyString)
	if err != nil || scrapeFrequency <= 0 {
		scrapeFrequency = 1
	}
	log.Printf("Scrape Frequency: %d minutes\n", scrapeFrequency)

    go startScraping(db, scrapeConcurrency, time.Duration(scrapeFrequency)*time.Minute)

	originsEnv := os.Getenv("ALLOWED_ORIGINS") 

    allowedOrigins := strings.Split(originsEnv, ",")
    for i := range allowedOrigins {
        allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i]) // cleanup
		log.Printf("ORIGIN: %s\n", allowedOrigins[i])
    }

    router := chi.NewRouter()

    router.Use(cors.Handler(cors.Options{
        AllowedOrigins: allowedOrigins,
        AllowedMethods: []string{"GET","POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
        ExposedHeaders: []string{"Link"},
        AllowCredentials: true,
        MaxAge: 300,
    }))


    srv := &http.Server{
        Handler: router, 
        Addr:   ":" + portString,
    }

    v1Router := chi.NewRouter()
    v1Router.Get("/healthz", handlerReadiness)
    v1Router.Get("/err", handlerErr)

    // users 
    v1Router.Post("/users", apiCfg.handlerCreateUser)
    v1Router.Post("/users/login", apiCfg.handlerLogin)
    v1Router.Post("/users/logout", apiCfg.middlewareAuth(apiCfg.handlerLogout))
    v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

    // feeds
    v1Router.Post("/feeds", apiCfg.handlerCreateFeed)
    v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
    v1Router.Get("/feeds/{feedID}", apiCfg.handlerGetFeedForID)
    v1Router.Get("/feeds/category/{categoryID}", apiCfg.handlerGetFeedsForCategory)
    v1Router.Get("/feeds/category", apiCfg.handlerGetFeedCategories)
    v1Router.Get("/feeds/search/{name}", apiCfg.handlerGetFeedsSearch)

    // OPML
    v1Router.Get("/feeds/opml",apiCfg.middlewareAuth(apiCfg.handlerGetOPMLUser))
    v1Router.Post("/feeds/opml", apiCfg.middlewareAuth(apiCfg.handlerImportOPMLUser))

    // feed follows
    v1Router.Post("/feeds/follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
    v1Router.Get("/feeds/follows", apiCfg.middlewareAuth(apiCfg.handlerGetFollowedFeeds))
    v1Router.Delete("/feeds/follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

    // posts
    v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))
    v1Router.Get("/posts/{feedID}", apiCfg.handlerGetFeedPosts)
    v1Router.Get("/posts/{feedID}/{date}", apiCfg.handlerGetFeedPostsDate)

    // posts views
    v1Router.Post("/posts/views",apiCfg.middlewareAuth(apiCfg.handlerCreatePostView))

    router.Mount("/v1", v1Router)

    log.Printf("Server starting on port %v", portString)
    err = srv.ListenAndServe()
    if err != nil {
        log.Fatal(err)
    }

}
