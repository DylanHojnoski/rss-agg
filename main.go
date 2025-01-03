package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"rssagg/internal/database"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
    
    db := database.New(conn)
    apiCfg := apiConfig{
        DB: db,
    }

    go startScraping(db, 10, time.Minute)

    router := chi.NewRouter()

    router.Use(cors.Handler(cors.Options{
        AllowedOrigins: []string{"https://*", "http://*"},
        AllowedMethods: []string{"GET","POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"*"},
        ExposedHeaders: []string{"Link"},
        AllowCredentials: false,
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
    v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

    // feeds
    v1Router.Post("/feeds", apiCfg.handlerCreateFeed)
    v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

    // feed follows
    v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
    v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
    v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

    // posts
    v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))
    v1Router.Get("/posts/{feedID}", apiCfg.handlerGetFeedPosts)
    v1Router.Get("/posts/{feedID}/{date}", apiCfg.handlerGetFeedPostsBeforeDate)


    router.Mount("/v1", v1Router)


    log.Printf("Server starting on port %v", portString)
    err = srv.ListenAndServe()
    if err != nil {
        log.Fatal(err)
    }

}
