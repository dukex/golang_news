package main

import (
	"fmt"
	rss "github.com/jteeuwen/go-pkg-rss"
	"os"
	"time"
)

var TWEETS map[string]string

func main() {
	FEEDS := []string{
		"http://www.mises.org.br/RSSArticles.aspx?type=3&culture=pt",
		"http://www.mises.org.br/RSSArticles.aspx?type=2&culture=pt",
		"http://www.mises.org.br/RSSArticles.aspx?type=1&culture=pt",
	}

	TWEETS = map[string]string{
		"http://www.brunogarschagen.com/":      "@BrunoGarschagen",
		"http://www.mises.org.br/Default.aspx": "@mises_brasil",
	}

	for _, feed := range FEEDS {
		go PollFeed(feed, 5)
	}
	PollFeed("http://feeds.feedburner.com/BrunoGarschagen?format=xml", 5)
}

func PollFeed(uri string, timeout int) {
	feed := rss.New(timeout, true, chanHandler, itemHandler)

	for {
		if err := feed.Fetch(uri, nil); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %s: %s", uri, err)
			return
		}

		<-time.After(time.Duration(feed.SecondsTillUpdate() * 1e9))
	}
}

func chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
}

func itemHandler(feed *rss.Feed, ch *rss.Channel, items []*rss.Item) {
	for _, item := range items {
		short_title := item.Title
		if len(short_title) > 100 {
			short_title = short_title[:99] + "…"
		}
		PostTweet(short_title + " " + item.Links[0].Href + " " + TWEETS[ch.Links[0].Href])
	}
}

func PostTweet(tweet string) {
	//anaconda.SetConsumerKey(ReadConsumerKey())
	//anaconda.SetConsumerSecret(ReadConsumerSecret())
	//api := anaconda.NewTwitterApi(ReadAccessToken(), ReadAccessTokenSecret())

	//v := url.Values{}
	//_, err := api.PostTweet(tweet, v)
	//if err != nil {
	//  log.Printf("Error posting tweet: %s", err)
	//}
	fmt.Println(tweet)
	fmt.Println("")
}
