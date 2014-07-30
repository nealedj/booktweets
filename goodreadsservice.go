package booktweets

import (
	"appengine/urlfetch"
	"encoding/xml"
	"io/ioutil"
	"net/url"
	"github.com/nealedj/anaconda"
)

type BooksResponse struct {
	BookList []Book `xml:"search>results>work"`
}

func (br BooksResponse) HasBooks() bool{
	return len(br.BookList) > 0
}

type Book struct {
	Title string `xml:"best_book>title" json:"Title"`
	ImageUrl string `xml:"best_book>image_url" json:"ImageUrl"`
	Tweets []anaconda.Tweet `json:"tweets"`
}

func (b Book) HasTweets() bool {
	return len(b.Tweets) > 0
}

func getBooksUrl(q string) string {
	booksUrl, err := url.Parse("https://www.goodreads.com/search.xml")

	if err != nil {
		panic(err)
	}

	qs := url.Values{}
	qs.Add("key", goodReadsKey)

	if q != "" {
		qs.Add("q", q)
	}
	booksUrl.RawQuery = qs.Encode()

	return booksUrl.String()
}

func getBooks(c RequestContext, q string) (BooksResponse, error) {
	client := urlfetch.Client(c)

	resp, err := client.Get(getBooksUrl(q))
	if err != nil {
		return BooksResponse{}, err
	}

	rawXml, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var booksResponse BooksResponse
	xml.Unmarshal(rawXml, &booksResponse)

	booksResponse.BookList = booksResponse.BookList[:4] // limit this to save twitter quota

	return booksResponse, nil
}
