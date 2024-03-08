package scraper

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"rss-aggregator/internal/database"
	"rss-aggregator/model"
	"sync"
	"time"
)

func fetchFeed(url string) (model.Rss, error) {
	rss := model.Rss{}
	resp, err := http.Get(url)
	if err != nil {
		return rss, err
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		decoder := xml.NewDecoder(resp.Body)
		err := decoder.Decode(&rss)
		if err != nil {
			return rss, err
		}

		return rss, nil
	}
	return rss, fmt.Errorf("status response: %v, url-> %v", resp.StatusCode, url)
}

func StartScrapper(nbrFeeds int, intervalSeconds int, db *database.Queries) {
	var wg sync.WaitGroup
	layout := "Mon, 02 Jan 2006 15:04:05 MST"
	tick := time.Tick(time.Second * time.Duration(intervalSeconds))

	for _ = range tick {
		fmt.Println("fetching feeds...")
		feeds, _ := db.GetNextFeedsToFetch(context.Background(), int32(nbrFeeds))
		for _, feed := range feeds {
			wg.Add(1)

			go func(feed_ database.Feed) {
				rss, err := fetchFeed(feed_.Url)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("feed fetch", feed_.Name)
					fmt.Println("feed fetch: ", rss.Channel.Title)
					for _, item := range rss.Channel.Item {
						date, err := time.Parse(layout, item.PubDate)
						if err != nil {
							fmt.Println(err)
						}
						post := database.CreatePostParams{
							ID:          uuid.New(),
							CreatedAt:   time.Now(),
							Url:         item.Link,
							UpdatedAt:   time.Now(),
							Title:       item.Title,
							FeedID:      feed_.ID,
							PublishedAt: sql.NullTime{Time: date, Valid: true},
							Description: sql.NullString{String: item.Description, Valid: true},
						}
						_, err = db.CreatePost(context.Background(), post)
						if err != nil {
							fmt.Println(err)
						}
					}
					err := db.MarkFeedFetched(context.Background(), feed_.ID)
					if err != nil {
						fmt.Println(err)
					}
				}
				wg.Done()
			}(feed)
		}
		wg.Wait()
	}
}
