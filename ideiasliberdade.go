package main

import (
	"fmt"
	// "github.com/ChimeraCoder/anaconda"
	rss "github.com/haarts/go-pkg-rss"
	"log"
	//"net/url"
	"os"
	//"regexp"
	"time"
)

const timeout = 50

var first = map[string]bool{}

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
		go PollFeed(feed, itemHandler)
	}
	PollFeed("http://feeds.feedburner.com/BrunoGarschagen?format=xml", itemHandler)
}

func itemHandler(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
	f := func(item *rss.Item) {
		short_title := item.Title
		if len(short_title) > 100 {
			short_title = short_title[:99] + "…"
		}
		PostTweet(short_title + " " + item.Links[0].Href + " " + TWEETS[ch.Links[0].Href])
	}

	genericItemHandler(feed, ch, newItems, f)
}

func PollFeed(uri string, itemHandler rss.ItemHandler) {
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
	//noop
}

func genericItemHandler(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item, individualItemHandler func(*rss.Item)) {
	log.Printf("%d new item(s) in %s\n", len(newItems), feed.Url)
	for _, item := range newItems {
		individualItemHandler(item)
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
}
