package booktweets

import (
	"net/http"
	"github.com/mjibson/appstats"
)

func init() {
	http.Handle("/", appstats.NewHandler(HomeHandler))
}