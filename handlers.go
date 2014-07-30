package booktweets

import (
	"os"
	"strings"
	"appengine"
	"html/template"
	"net/http"
	"encoding/json"
)

type RequestContext struct {
	appengine.Context // anonymous field so can we pass this object to GAE API clients
	HttpRequest *http.Request
}

type HomeViewModel struct {
	BooksResponse // anonymous field
	Query         string
	Bootstrapped  bool
	SyncTwitterCalls bool
}

var homeTmpl = template.Must(template.ParseFiles("templates/home.html"))

func HomeHandler(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	context := RequestContext{c,r}

	viewModel := HomeViewModel{
		SyncTwitterCalls: r.FormValue("sync") == "1",
		Query: r.FormValue("q"),
		Bootstrapped: r.FormValue("bs") == "1",
	}

	if viewModel.Query == "" {
		homeTmpl.Execute(w, viewModel)
		return
	}

	if viewModel.Bootstrapped {
		getBootstrappedData(&viewModel)
		homeTmpl.Execute(w, viewModel)
		return
	}

	booksResponse, err := getBooks(context, viewModel.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	viewModel.BooksResponse = booksResponse

	if viewModel.SyncTwitterCalls {
		getTweetsSync(context, viewModel.BookList)
	} else {
		getTweets(context, viewModel.BookList)
	}

	homeTmpl.Execute(w, viewModel)
}


func getBootstrappedData(viewModel *HomeViewModel) {

	if viewModel.Query == "" {panic("viewModel.Query empty")}
	fileName := "bootstrappeddata/" + strings.ToLower(viewModel.Query) + ".json"

	fi, err := os.Open(fileName)
	
	if err != nil { 
		return
	}
	defer fi.Close()
 
	var bookList []Book
	err = json.NewDecoder(fi).Decode(&bookList)

	if err != nil { panic(err) }

	viewModel.BookList = bookList
}
