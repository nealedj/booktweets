package booktweets

import (
	"strings"
	"sync"
	"net/url"
	"github.com/nealedj/anaconda"
	"appengine/urlfetch"
	"encoding/json"
)

type searchResponse struct {
	Statuses []anaconda.Tweet
}

func getTwitterKeyword(bookTitle string) string {
	bookTitle = strings.ToLower(bookTitle)

	// get everything before brackets
	bookTitle = strings.SplitN(bookTitle, "(", 2)[0]

	if theIdx := strings.Index(bookTitle, "the"); theIdx != -1 {
		// if "the" is in title then read everything after it
		bookTitle = bookTitle[theIdx+3:]
	}
	return strings.TrimSpace(bookTitle)
}

func getTweets(c RequestContext, bookList []Book) {
	anaconda.SetConsumerKey(twitterConsumerKey)
	anaconda.SetConsumerSecret(twitterConsumerSecret)
	api := anaconda.NewTwitterApi(twitterAccessToken, twitterAccessTokenSecret)
	api.HttpClient = urlfetch.Client(c)
	defer api.Close()

	var wg sync.WaitGroup

	wg.Add(len(bookList))
	for i := 0; i < len(bookList); i++ {
		go func(book *Book) {
			term := getTwitterKeyword(book.Title)
			c.Infof("Getting tweets for %s", term)

			v := url.Values{}
			v.Add("lang", "en")

			tweets, err := api.GetSearchNoQueue(term, v)
			defer wg.Done()

			if err != nil {
				c.Errorf("Twitter call failed for %s: %s", term, err)
				return
			}

			book.Tweets = tweets
			c.Infof("Got %d tweets for %s", len(book.Tweets), term)

		}(&bookList[i])
	}

	wg.Wait()

	PrintAsJson(c, bookList)
}

func getTweets__WaitWithChannel(c RequestContext, bookList []Book) {
	anaconda.SetConsumerKey(twitterConsumerKey)
	anaconda.SetConsumerSecret(twitterConsumerSecret)
	api := anaconda.NewTwitterApi(twitterAccessToken, twitterAccessTokenSecret)
	api.HttpClient = urlfetch.Client(c)
	defer api.Close()

	twitterc := make(chan error)

	for i := 0; i < len(bookList); i++ {
		go func(book *Book) {
			term := getTwitterKeyword(book.Title)
			c.Infof("Getting tweets for %s", term)

			v := url.Values{}
			v.Add("lang", "en")

			tweets, err := api.GetSearchNoQueue(term, v)

			if err != nil {
				c.Errorf("Twitter call failed for %s: %s", term, err)
			} else {
				book.Tweets = tweets
				c.Infof("Got %d tweets for %s", len(book.Tweets), term)
			}
			twitterc <- err

		}(&bookList[i])
	}

	for i := 0; i < len(bookList); i++ {
		<-twitterc
	}

	PrintAsJson(c, bookList)
}

func getTweetsSync(c RequestContext, bookList []Book){
	anaconda.SetConsumerKey(twitterConsumerKey)
	anaconda.SetConsumerSecret(twitterConsumerSecret)
	api := anaconda.NewTwitterApi(twitterAccessToken, twitterAccessTokenSecret)
	api.HttpClient = urlfetch.Client(c)
	defer api.Close()

	for i := 0; i < len(bookList); i++ {
		book := &bookList[i]
		term := getTwitterKeyword(book.Title)
		c.Infof("Getting tweets for %s", term)
		
		v := url.Values{}
		v.Add("lang", "en")
		tweets, err := api.GetSearch(term, nil)

		if err != nil {
			c.Errorf("Twitter call failed for %s: %s", term, err)
			return
		}
		book.Tweets = tweets
		c.Infof("Got %d tweets for %s", len(book.Tweets), term)
	}

	PrintAsJson(c, bookList)
}

func PrintAsJson(c RequestContext, data interface{}) {
	raw, _ := json.Marshal(data)
	c.Infof(string(raw[:]))
}

