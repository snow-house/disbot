package reddit

import (
	"github.com/jzelinskie/geddit"
	// "fmt"
	"log"
	"os"
	"sort"

	// "strings"
	"math/rand"
	"regexp"
)

var (
	// session *geddit.OAuthSession
	session *geddit.Session
	subOpts geddit.ListingOptions

	REDDITCLIENTID     string
	REDDITCLIENTSECRET string
	REDDITREFRESHTOKEN string
	REDDITACCESSTOKEN  string

	REDDITUSERNAME string
	REDDITPWD      string

	refreshRE *regexp.Regexp
)

func init() {

	// // read env
	REDDITCLIENTID = os.Getenv("REDDITCLIENTID")
	REDDITCLIENTSECRET = os.Getenv("REDDITCLIENTSECRET")
	REDDITACCESSTOKEN = os.Getenv("REDDITACCESSTOKEN")
	REDDITREFRESHTOKEN = os.Getenv("REDDITREFRESHTOKEN")

	REDDITUSERNAME = os.Getenv("REDDITUSERNAME")
	REDDITPWD = os.Getenv("REDDITPWD")

	// var err error
	// // init geddit session
	// session, err = geddit.NewOAuthSession(
	// 	REDDITCLIENTID,
	// 	REDDITCLIENTSECRET,
	// 	"gedditAgent v2",
	// 	"redirect url",
	// )

	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println("geddit oauth session opened")

	// // // login using personal reddit account
	// err = session.LoginAuth(REDDITUSERNAME, REDDITPWD)
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println("geddit login succesful")

	// scopes := strings.Split("identity,edit,flair,history,modconfig,modflair,modlog,modposts,modwiki,mysubreddits,privatemessages,read,repost,save,submit,subscibe,vote,wikiedit,wikiread", ",")
	// scopes := strings.Split("read", ",")
	// code := session.AuthCodeURL("state", scopes)
	// log.Println(code)

	// try using normal session instead of oauth, since we are not going
	// to do anything beside fetching post like posting and upvoting
	session = geddit.NewSession("aryuuu")

	subOpts = geddit.ListingOptions{
		Limit: 100,
	}

	// var e error
	refreshRE = regexp.MustCompile("(?i)refresh")

}

// fetch a post from a subreddit
func R(subreddit string, comment int) (status bool, title, url, desc, flair, comments string) {

	log.Println("fetching post from " + subreddit)
	posts, err := session.SubredditSubmissions(subreddit, geddit.HotSubmissions, subOpts)

	if err != nil {
		// log.Println("failed when trying to fetch post from " + subreddit)
		log.Println(err)
		if refreshRE.MatchString(err.Error()) {
			// session.LoginAuth(REDDITUSERNAME, REDDITPWD)
			session = geddit.NewSession("aryuuu")
			return R(subreddit, comment)
		} else {
			return false, "something wrong", "url", "desc", "flair", "comments"
		}
	}

	// pick one random post from hot
	idxRange := 100
	if len(posts) < idxRange {
		idxRange = len(posts)
	}

	idx := rand.Intn(idxRange)

	// log.Println(posts[idx].Title)
	// log.Println(posts[idx].URL)
	var coms []*geddit.Comment
	c := ""
	if comment > 0 {
		coms, err = session.Comments(posts[idx])

		if err == nil {

			sort.SliceStable(coms, func(i, j int) bool {
				return coms[i].Score > coms[j].Score
			})

			commentRange := comment
			if len(coms) < commentRange {
				commentRange = len(comments)
			}

			for i := 0; i < commentRange; i++ {
				c += ">> " + coms[i].Body + "\n"
			}

		}
	}

	return true, posts[idx].Title, posts[idx].URL, posts[idx].Selftext, posts[idx].LinkFlairText, c
}

// fetch a post and comments from r/askreddit
func Ask() (status bool, title, desc, comments string) {

	log.Println("fetching post from askreddit")
	posts, err := session.SubredditSubmissions("askreddit", geddit.HotSubmissions, subOpts)
	if err != nil {
		log.Println("failed when trying to fetch post from askreddit")
		log.Println(err)
		if refreshRE.MatchString(err.Error()) {
			// session.LoginAuth(REDDITUSERNAME, REDDITPWD)
			session = geddit.NewSession("aryuuu")
			return Ask()
		} else {
			return false, "something wrong :(", "desc", "comments"
		}
	}

	// pick one random post from hot
	idxRange := 100
	if len(posts) < idxRange {
		idxRange = len(posts)
	}

	idx := rand.Intn(idxRange)

	s := geddit.NewSession("aryuuu")

	coms, err := s.Comments(posts[idx])
	c := ""
	if err == nil {

		// sort comments descending
		sort.SliceStable(coms, func(i, j int) bool {
			return coms[i].Score > coms[j].Score
		})

		for i := 0; i < 3; i++ {
			c += ">> " + coms[i].Body + "\n"
		}
	}

	// log.Println(posts[idx].Title)
	// log.Println(posts[idx].Selftext)

	return true, posts[idx].Title, posts[idx].Selftext, c
}

// fetch a meme from r/dankmemes
func Random() (status bool, title, url string) {

	log.Println("fetching post from dankmemes")
	posts, err := session.SubredditSubmissions("dankmemes", geddit.HotSubmissions, subOpts)
	if err != nil {
		log.Println("failed when trying to fetch post from dankmemes")
		log.Println(err)

		if refreshRE.MatchString(err.Error()) { // if re login already attempted
			// session.LoginAuth(REDDITUSERNAME, REDDITPWD)
			session = geddit.NewSession("aryuuu")
			return Random()
		} else {
			return false, "something wrong :(", "url"
		}
	}

	// pick one random post from hot
	idxRange := 100
	if len(posts) < idxRange {
		idxRange = len(posts)
	}

	idx := rand.Intn(idxRange)

	// log.Println(posts[idx].Title)
	// log.Println(posts[idx].URL)

	return true, posts[idx].Title, posts[idx].URL
}
