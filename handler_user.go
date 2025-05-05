package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"rssagg/internal/database"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (apiCfg *apiConfig) handlerCreateUser (w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Error parsing JSON: %v", err))
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 10)
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Error hashing password: %v", err))
        return
    }

    user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
        ID: uuid.New(),
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
        Username: params.Username,
        Password: string(hashedPassword),
    })
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't create user: %v", err))
        return
    }

    respondWithJSON(w, 201, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerLogin (w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Error parsing JSON: %v", err))
        return
    }

    user, err := apiCfg.DB.GetUser(r.Context(), params.Username)
    if err != nil {
        respondWithError(w,401, fmt.Sprintf("Username does not exist: %v", err))
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
    if err != nil {
        respondWithError(w,401, fmt.Sprintf("Invalid login information: %v", err))
        return
    }

    claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
        Issuer: user.ID.String(), 
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
    })

    token, err := claims.SignedString([]byte(os.Getenv("KEY")))

    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Couldn't login: %v", err))
        return
    }

    cookie := http.Cookie{
        Name: "jwt",
        Value: token,
        Expires: time.Now().Add(time.Hour * 24),
        HttpOnly: true,
        Path: "/",
    }

    http.SetCookie(w, &cookie)
    w.Header().Set("Access-Control-Allow-Credentials", "true")

    respondWithJSON(w, 200, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerLogout (w http.ResponseWriter, r *http.Request, user database.User) {
    cookie := http.Cookie{
        Name: "jwt",
        MaxAge: -1,
        HttpOnly: true,
        Path: "/",
    }

    http.SetCookie(w, &cookie)
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    respondWithJSON(w, 200, nil)
}

func (apiCfg *apiConfig) handlerGetUser (w http.ResponseWriter, r *http.Request, user database.User) {
    respondWithJSON(w, 200, databaseUserToUser(user))
}
 
func (apiCfg *apiConfig) handlerGetPostsForUser (w http.ResponseWriter, r *http.Request, user database.User) {
    posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
        UserID: user.ID,
        Limit: 10,
    })
    if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't get posts: %v", err))
        return
    }

    respondWithJSON(w, 200, databasePostsForUserToPosts(posts))
}
