package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dukex/buffer"
	"github.com/dukex/squeue"
	rss "github.com/haarts/go-pkg-rss"
	"github.com/jinzhu/gorm"
	pq "github.com/lib/pq"
)

const timeout = 50

var (
	database gorm.DB
	first    = map[string]bool{}
	TWEETS   map[string]string
	FEEDS    []string
	QUEUE    *squeue.Queue
	BUFFER   *buffer.Client
	PROFILES []string
)

type Post struct {
	Id    int64
	Key   string `sql:"not null;unique"`
	Title string `sql:"not null"`
	Link  string
}

func (p *Post) AfterSave() (err error) {
	QUEUE.Push(func() {
		text := p.Title + " " + p.Link
		BUFFER.CreateUpdate(text, PROFILES, map[string]interface{}{
			"now": true,
		})
	})
	return
}

func main() {
	FEEDS = []string{
		"http://www.mises.org.br/RSSArticles.aspx?type=3&culture=pt",
		"http://www.mises.org.br/RSSArticles.aspx?type=2&culture=pt",
		"http://www.mises.org.br/RSSArticles.aspx?type=1&culture=pt",
		"http://feeds.feedburner.com/BrunoGarschagen?format=xml",
		"http://maovisivel.blogspot.com/feeds/posts/default?alt=rss",
		"http://www.liberdade.cc/feed",
		"http://ordemlivre.org/feed/blogs",
		"http://epl.org.br/feed/",
		"http://feeds.feedburner.com/org/hetj?format=xml",
		"http://ordemlivre.org/feed/artigos",
		"http://www.libertarianismo.org/index.php/category/artigos/feed/",
		"http://www.institutoliberal.org.br/blog/feed/",
		"http://liberzone.com.br/feed/",
		"http://gdata.youtube.com/feeds/api/users/UCou8RLI69bMShhe_ziija-Q/uploads", // Talk Show com Evandro Sinotti
		"http://www.clubemissrand.com/1/feed",
		"http://spotniks.com/feed/",
	}

	databaseUrl, _ := pq.ParseURL(os.Getenv("DATABASE_URL"))
	database, _ = gorm.Open("postgres", databaseUrl)
	database.LogMode(os.Getenv("DEBUG") == "true")

	database.AutoMigrate(Post{})

	BUFFER = buffer.NewClient(os.Getenv("BUFFER_TOKEN"))
	PROFILES = make([]string, 0)
	for _, profile := range BUFFER.Profiles() {
		PROFILES = append(PROFILES, profile.Id)
	}

	QUEUE = squeue.NewQueue(1 * time.Minute)
	QUEUE.Run()

	http.HandleFunc("/", HomeHandler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	for _, feed := range FEEDS {
		go PollFeed(feed, itemHandler)
	}

	fmt.Fprintf(w, "Oi!")
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func itemHandler(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
	f := func(item *rss.Item) {
		key := GetMD5Hash("item-" + item.Key())
		short_title := item.Title
		if len(short_title) > 100 {
			short_title = short_title[:99] + "…"
		}

		err := database.Table("posts").Where("key = ?", key).Scan(&Post{}).Error
		if err == gorm.RecordNotFound {
			var post Post
			post.Title = short_title
			post.Key = key
			post.Link = item.Links[0].Href
			database.Save(&post)
		}
	}

	genericItemHandler(feed, ch, newItems, f)
}

func PollFeed(uri string, itemHandler rss.ItemHandler) {
	feed := rss.New(timeout, true, chanHandler, itemHandler)

	if err := feed.Fetch(uri, nil); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %s: %s", uri, err)
		return
	}

	<-time.After(time.Duration(feed.SecondsTillUpdate() * 1e9))

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
