package main

import (
  "crypto/md5"
  "encoding/hex"
  "fmt"
  "github.com/ChimeraCoder/anaconda"
  "github.com/dukex/squeue"
  rss "github.com/haarts/go-pkg-rss"
  "log"
  "net/http"
  "net/url"
  "os"
  "time"
)

const timeout = 50

var (
  first  = map[string]bool{}
  TWEETS map[string]string
  FEEDS  []string
  QUEUE  *squeue.Queue
)

func main() {
  FEEDS = []string{
    "http://www.mises.org.br/RSSArticles.aspx?type=3&culture=pt",
    "http://www.mises.org.br/RSSArticles.aspx?type=2&culture=pt",
    "http://www.mises.org.br/RSSArticles.aspx?type=1&culture=pt",
    "http://feeds.feedburner.com/BrunoGarschagen?format=xml",
    "http://maovisivel.blogspot.com/feeds/posts/default",
    "http://www.liberdade.cc/feed",
    "http://ordemlivre.org/feed/blogs",
    "http://epl.org.br/feed/",
    "http://feeds.feedburner.com/org/hetj?format=xml",
    "http://ordemlivre.org/feed/artigos",
    "http://www.libertarianismo.org/index.php/category/artigos/feed/",
    "http://www.institutoliberal.org.br/blog/feed/",
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

func Exists(name string) bool {
  if _, err := os.Stat(name); err != nil {
    if os.IsNotExist(err) {
      return false
    }
  }
  return true
}

func GetMD5Hash(text string) string {
  hasher := md5.New()
  hasher.Write([]byte(text))
  return hex.EncodeToString(hasher.Sum(nil))
}

func itemHandler(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
  f := func(item *rss.Item) {
    file := GetMD5Hash("data/" + item.Key())
    root_on_heroku := "/apps/data/"

    if !Exists(root_on_heroku + file) {
      fo, err := os.Create(root_on_heroku + file)
      if err != nil {
        fmt.Println("CREATE ERROR", err)
      } else {
        defer fo.Close()
        buf := make([]byte, 1024)

        if _, err = fo.Write(buf[:]); err != nil {
          fmt.Println("WRITE ERROR:", err)
        } else {
          short_title := item.Title
          if len(short_title) > 100 {
            short_title = short_title[:99] + "…"
          }

          QUEUE.Push(func() {
            err = PostTweet(short_title + " " + item.Links[0].Href)
            if err != nil {
              os.Remove(root_on_heroku + file)
            }
          })
        }
      }
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

func PostTweet(tweet string) error {
  anaconda.SetConsumerKey(ReadConsumerKey())
  anaconda.SetConsumerSecret(ReadConsumerSecret())
  api := anaconda.NewTwitterApi(ReadAccessToken(), ReadAccessTokenSecret())

  v := url.Values{}
  _, err := api.PostTweet(tweet, v)
  if err != nil {
    log.Printf("Error posting tweet: %s", err)
  }
  return err
}
