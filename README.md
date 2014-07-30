Book Tweets
====================

Book Tweets is a simple application written in Go to search for books on Goodreads and then search for related tweets on each book

Setup
-----

Change the application ID "davidntest28" in app.yaml to your own application ID before you deploy:

Create a keys.go file with API keys for Goodreads and Twitter as below:


	package booktweets


	var twitterConsumerKey = "a9a9a9a9a9a9a9a9a9a9a9"
	var twitterConsumerSecret = "a9a9a9a9a9a9a9a9a9a9a9a9a9a9a9a9a9a9a9a9a9a9"
	var twitterAccessToken = "a9a9a9a9a9a9a9a9a9a9a9-a9a9a9a9a9a9a9a9a9a9a9a9a9a9a9"
	var twitterAccessTokenSecret = "a9a9a9a9a9a9a9a9a9a9a9"

	var goodReadsKey = "a9a9a9a9a9a9a9a9a9a9a9"
	var goodReadsSecret = "a9a9a9a9a9a9a9a9a9a9a9a9a9a9a9a9a9a9a9"
