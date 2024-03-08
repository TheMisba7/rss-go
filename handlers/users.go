package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"rss-aggregator/internal/database"
	"rss-aggregator/model"
	"time"
)

func (config *ApiConfig) getByAPIKey() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		apiKey := request.Header.Get("Authorization")
		userByAPIkey, err := config.DB.GetUserByAPIkey(request.Context(), apiKey)
		if err != nil {
			panic(err)
		}
		respondWithJSON(writer, 200, userByAPIkey)
	}
}
func (config *ApiConfig) newUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		now := time.Now()
		user := model.User{
			Id:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		}
		decoder := json.NewDecoder(request.Body)
		decoder.Decode(&user)
		userParams := database.CreateUserParams{Name: user.Name, ID: user.Id, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}

		createUser, err := config.DB.CreateUser(request.Context(), userParams)
		if err != nil {
			fmt.Println("my err")
			panic(err)
		}
		respondWithJSON(writer, 200, createUser)
	}
}
