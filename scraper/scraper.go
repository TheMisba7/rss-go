package scraper

import (
	"context"
	"encoding/xml"
	"fmt"
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
					err := db.MarkFeedFetched(context.Background(), feed_.ID)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("feed fetch", feed_.Name)
					fmt.Println("feed fetch: ", rss.Channel.Title)
				}
				wg.Done()
			}(feed)
		}
		wg.Wait()
	}
}
