package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"rss-aggregator/internal/database"
	"rss-aggregator/model"
	"time"
)

func (config *ApiConfig) createFeed(w http.ResponseWriter, req *http.Request, user database.User) {
	now := time.Now()
	newFeed := model.Feed{
		CreatedAt: now,
		UpdatedAt: now,
	}
	parseJson(req, &newFeed)
	newFeed.UserId = user.ID
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      newFeed.Name,
		CreatedAt: newFeed.CreatedAt,
		UpdatedAt: newFeed.UpdatedAt,
		UserID:    newFeed.UserId,
		Url:       newFeed.Url,
	}
	createFeed, err := config.DB.CreateFeed(req.Context(), feedParams)
	if err != nil {
		panic(err)
	}
	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		FeedID:    createFeed.ID,
		UserID:    createFeed.UserID,
	}

	follow, err := config.DB.CreateFeedFollow(req.Context(), followParams)
	if err != nil {
		return
	}
	newFeed.Id = createFeed.ID
	respondWithJSON(w, 201,
		model.CreateFeedRS{
			Feed: newFeed,
			FeedFollow: model.FeedFollow{
				Id:        follow.ID,
				FeedId:    follow.FeedID,
				UserId:    follow.UserID,
				UpdatedAt: follow.UpdatedAt,
				CreatedAt: follow.CreatedAt,
			},
		})
}

func (config *ApiConfig) getAllFeeds() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		feeds, err := config.DB.GetALlFeeds(request.Context())
		if err != nil {
			panic(err)
		}
		respondWithJSON(writer, 200, feeds)
	}
}

func (config *ApiConfig) followFeed(w http.ResponseWriter, req *http.Request, user database.User) {
	now := time.Now()
	feedFollow := model.FeedFollow{
		Id:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserId:    user.ID,
	}
	parseJson(req, &feedFollow)
	params := database.CreateFeedFollowParams{
		ID:        feedFollow.Id,
		FeedID:    feedFollow.FeedId,
		UserID:    feedFollow.UserId,
		UpdatedAt: feedFollow.UpdatedAt,
		CreatedAt: feedFollow.CreatedAt,
	}
	createFeedFollow, err := config.DB.CreateFeedFollow(req.Context(), params)
	if err != nil {
		panic(err)
	}
	respondWithJSON(w, 201, createFeedFollow)
}

func (config *ApiConfig) getFeedFollowByUser(w http.ResponseWriter, req *http.Request, user database.User) {
	feedFollows, err := config.DB.GetAllFeedFollowByUser(req.Context(), user.ID)
	if err != nil {
		panic(err)
	}
	respondWithJSON(w, 200, feedFollows)
}

func (config *ApiConfig) unfollow(w http.ResponseWriter, req *http.Request, user database.User) {
	feedFollowId := chi.URLParam(req, "feedFollowId")
	if len(feedFollowId) == 0 {
		respondWithError(w, 400, "Bad request.")
	}

	followById, err := config.DB.GetByFeedFollowById(req.Context(), uuid.MustParse(feedFollowId))
	if err != nil {
		panic(err)
	}

	if user.ID != followById.ID {
		respondWithError(w, 401, "Not authorized")
	}
	err = config.DB.DeleteFeedFollow(req.Context(), uuid.MustParse(feedFollowId))
	if err != nil {
		panic(err)
	}
	w.WriteHeader(200)
	w.Write([]byte("OK"))

}
