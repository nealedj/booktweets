package booktweets

import (
		"testing"
		"appengine/aetest"
)

func TestGetTweets(t *testing.T) {
	aeContext, _ := aetest.NewContext(nil)
	defer aeContext.Close()

	c := RequestContext{ AppengineContext: aeContext, }
 
	bookList := []Book{
		Book {
			Title: "The Name of the wind",
		},
		Book {
			Title: "Wise man's fear",
		},
	}

	getTweets(c, bookList)

	for _, book := range bookList {
		if len(book.Tweets) == 0 {
			t.Error("No tweets for " + book.Title)
		}
	}

}

func TestGetTwitterKeyword(t *testing.T) {
	if getTwitterKeyword("The Name of the Wind") != "name of the wind" {
		t.Fail()
	}

	if getTwitterKeyword("Star Wars: The Empire Strikes Back") != "empire strikes back" {
		t.Fail()
	}

	if getTwitterKeyword("The Girl with the Dragon Tattoo (Volume #1)") != "girl with the dragon tattoo" {
		t.Fail()
	}
}