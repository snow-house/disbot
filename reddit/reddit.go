package reddit

import (
	"github.com/jzelinskie/geddit"
	// "fmt"
	"log"
	"os"
	// "strings"
	"math/rand"
)

var (
	session *geddit.OAuthSession
	subOpts geddit.ListingOptions

	REDDITCLIENTID string
	REDDITCLIENTSECRET string
	REDDITREFRESHTOKEN string
	REDDITACCESSTOKEN string

	REDDITUSERNAME string
	REDDITPWD string

)

func init() {

	// read env
	REDDITCLIENTID = os.Getenv("REDDITCLIENTID")
	REDDITCLIENTSECRET = os.Getenv("REDDITCLIENTSECRET")
	REDDITACCESSTOKEN = os.Getenv("REDDITACCESSTOKEN")
	REDDITREFRESHTOKEN = os.Getenv("REDDITREFRESHTOKEN")

	REDDITUSERNAME = os.Getenv("REDDITUSERNAME")
	REDDITPWD = os.Getenv("REDDITPWD")


	// init geddit session 
	session, err := geddit.NewOAuthSession(
		REDDITCLIENTID,
		REDDITCLIENTSECRET,
		"gedditAgent v2",
		"redirect url",
	)

	if err != nil {
		log.Println(err)
	}
	log.Println("geddit oauth session opened")

	// login using personal reddit account
	err = session.LoginAuth(REDDITUSERNAME, REDDITPWD)
	if err != nil {
		log.Println(err)
	}
	log.Println("geddit login succesful")

	subOpts = geddit.ListingOptions {
		Limit: 100,
	}

}

// fetch a post from a subreddit
func R(subreddit string, comment int) (status bool, title, url, desc, flair, comments string) {
	
	posts, err := session.SubredditSubmissions(subreddit, geddit.HotSubmissions, subOpts)

	if err != nil {
		log.Println("failed when trying to fetch post from " + subreddit)
		log.Println(err)
		return false, "something wrong", "url", "desc", "flair", "comments"
	}

	// pick one random post from hot
	idxRange := 100
	if len(posts) < idxRange {
		idxRange = len(posts)
	}

	idx := rand.Intn(idxRange)






	return true, posts[idx].Title , posts[idx].URL, posts[idx].Selftext, posts[idx].LinkFlairText, "empty"
}

// fetch a post and comments from r/askreddit
func Ask() (status bool, title, desc, comments string) {

	posts, err := session.SubredditSubmissions("askreddit", geddit.HotSubmissions, subOpts)
	if err != nil {
		log.Println("failed when trying to fetch post from askreddit")
		log.Println(err)
		return false, "something wrong :(", "desc", "comments"
	}

	// pick one random post from hot
	idxRange := 100
	if len(posts) < idxRange {
		idxRange = len(posts)
	}

	idx := rand.Intn(idxRange)

	return true, posts[idx].Title, posts[idx].Selftext, "comments"
}

// fetch a meme from r/dankmemes
func Random() (status bool, title, url string) {

	posts, err := session.SubredditSubmissions("dankmemes", geddit.HotSubmissions, subOpts)
	if err != nil {
		log.Println("failed when trying to fetch post from dankmemes")
		log.Println(err)
		return false, "something wrong :(", "url"
	}

	// pick one random post from hot
	idxRange := 100
	if len(posts) < idxRange {
		idxRange = len(posts)
	}

	idx := rand.Intn(idxRange)

	return true, posts[idx].Title, posts[idx].URL
}

