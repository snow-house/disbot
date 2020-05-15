package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jzelinskie/geddit"

	"fmt"
	"log"
	"regexp"
	"os"
	"time"
	"strings"
	"strconv"

	"../nim"
	"../tag"
	"../reddit"
)

var (
	Token string
	BotID string
	bot *discordgo.Session

	deleteTagRE *regexp.Regexp
	infoTagRE *regexp.Regexp


	redditSession *geddit.OAuthSession
	// subOpts geddit.ListingOptions

	REDDITCLIENTID string
	REDDITCLIENTSECRET string
	REDDITREFRESHTOKEN string
	REDDITACCESSTOKEN string

	REDDITUSERNAME string
	REDDITPWD string

	
)

func init() {
	
	
	deleteTagRE, _ = regexp.Compile("^/deletetag .*")
	infoTagRE, _ = regexp.Compile("^/taginfo .*")

	// read env
	REDDITCLIENTID = os.Getenv("REDDITCLIENTID")
	REDDITCLIENTSECRET = os.Getenv("REDDITCLIENTSECRET")
	REDDITACCESSTOKEN = os.Getenv("REDDITACCESSTOKEN")
	REDDITREFRESHTOKEN = os.Getenv("REDDITREFRESHTOKEN")

	REDDITUSERNAME = os.Getenv("REDDITUSERNAME")
	REDDITPWD = os.Getenv("REDDITPWD")

	log.Println("REDDITCLIENTID: " + REDDITCLIENTID)
	log.Println("REDDITCLIENTSECRET: " + REDDITCLIENTSECRET)


	// init geddit session 
	redditSession, err := geddit.NewOAuthSession(
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
	err = redditSession.LoginAuth(REDDITUSERNAME, REDDITPWD)
	if err != nil {
		log.Println(err)
	}
	log.Println("geddit login successful")
	log.Printf("geddit session: %v", redditSession)
	// log.Println(redditSession)
		
}

func Start() {
	Token = os.Getenv("DISBOTTOKEN")

	// create new discord session using provided bot token
	bot, err := discordgo.New("Bot "+Token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// add handlers
	bot.AddHandler(fuckHandler)

	// nim
	bot.AddHandler(nimHandler)

	// tag
	bot.AddHandler(getTagHandler)
	bot.AddHandler(addTagHandler)
	bot.AddHandler(deleteTagHandler)
	bot.AddHandler(listTagHandler)
	bot.AddHandler(infoTagHandler)

	// reddit
	bot.AddHandler(rHandler)
	bot.AddHandler(askHandler)
	bot.AddHandler(randomHandler)



	// open the websocket and begin listening
	err = bot.Open()
	if err != nil {
		log.Println("Error opening Discord session: ", err)
	}

	log.Println("Bot is running")
}

// simple pings
func fuckHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	if m.Content == "fuck" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "yeah")
	}
}

// nim command handler
func nimHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	nimRE, _ := regexp.Compile("^/nim .*")
	// if message match with nim command
	if nimRE.MatchString(m.Content) {
		
		// extract query
		query := strings.Replace(nimRE.FindString(m.Content), "/nim ", "", -1)
		
		// find nim or name
		result := nim.Find(query)

		// send reply
		_, _ = s.ChannelMessageSend(m.ChannelID, result)
	}	
}

// get tag handler
func getTagHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	var name string
	var scope int

	if m.Author.Bot {
		return
	}

	publicTagRE, _ := regexp.Compile(`"([^"]+)"`)
	guildTagRE, _ := regexp.Compile(":([^:])+:")
	channelTagRE, _ := regexp.Compile(";([^;])+;")

	if publicTagRE.MatchString(m.Content) {
		name = strings.Replace(publicTagRE.FindString(m.Content), `"`, "", -1)
		scope = 2
		
	} else if guildTagRE.MatchString(m.Content) {
		name = strings.Replace(guildTagRE.FindString(m.Content), ":", "", -1)
		scope = 1
		
	} else if channelTagRE.MatchString(m.Content) {
		name = strings.Replace(channelTagRE.FindString(m.Content), ";", "", -1)
		scope = 0
		
	} else {
		return
	}

	status, url := tag.Get(name, m.ChannelID, m.GuildID, scope)

	if !status {
		s.ChannelMessageSend(m.ChannelID, "tag "+ name +" doesn't exist")
		return
	}

	embed := &discordgo.MessageEmbed{
	    Author:      &discordgo.MessageEmbedAuthor{},
	    Color:       0x00ff00, // Green
	    // Description: "This is a discordgo embed",
	    // Fields: []*discordgo.MessageEmbedField{
	    //     &discordgo.MessageEmbedField{
	    //         Name:   "I am a field",
	    //         Value:  "I am a value",
	    //         Inline: true,
	    //     },
	    //     &discordgo.MessageEmbedField{
	    //         Name:   "I am a second field",
	    //         Value:  "I am a value",
	    //         Inline: true,
	    //     },
	    // },
	    Image: &discordgo.MessageEmbedImage{
	        URL: url,
	    },
	    // Thumbnail: &discordgo.MessageEmbedThumbnail{
	    //     URL: url,
	    // },
	    Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
	    Title:     name,
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return
}

// add tag handler
func addTagHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	
	if m.Author.ID == BotID {
		return
	}

	addTagRE, _ := regexp.Compile("/addtag .*")
	if addTagRE.MatchString(m.Content) {
		args := strings.Split(m.Content, " ")

		var status bool

		// handle guild scope
		if len(args) == 3 {
			status = tag.Add(args[1], args[2], m.Author.ID, m.ChannelID, m.GuildID, 1)
		} else if len(args) == 4 { // handle channel and guild scope
			if args[3] == "-c" {
				status = tag.Add(args[1], args[2], m.Author.ID, m.ChannelID, m.GuildID, 0)
			} else if args[3] == "-p" {
				status = tag.Add(args[1], args[2], m.Author.ID, m.ChannelID, m.GuildID, 2)
			} else {
				s.ChannelMessageSend(m.ChannelID, "unknown flag, you might need some /help")
				return
			}
		} else { // send error message
			s.ChannelMessageSend(m.ChannelID, "wrong arguments, you might need some /help")
			return
		}

		var reply string
		if status {
			reply = "tag " + args[1] + " added"
		} else {
			reply = "failed to add tag"
		}
		s.ChannelMessageSend(m.ChannelID, reply)
	}
}

// delete tag handler
func deleteTagHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	
}

// list tag handler
func listTagHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == BotID {
		return
	}

	listTagRE, _ := regexp.Compile("^/taglist")
	if listTagRE.MatchString(m.Content) {
		result :=  tag.List(m.ChannelID, m.GuildID)

		_, _ = s.ChannelMessageSend(m.ChannelID, result)
	}
}

// info tag handler
func infoTagHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	
}

// random meme handler
func randomHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == BotID {
		return
	}

	randomRE, _ := regexp.Compile("^/random")
	if randomRE.MatchString(m.Content) {

		log.Printf("reddit session: %v", redditSession)
		// log.Println(redditSession)
		status, title, url := reddit.Random(redditSession)
		if !status {
			s.ChannelMessageSend(m.ChannelID, "something wrong :(")
			return
		}

		log.Println("fetched post from dankmemes")
		embed := &discordgo.MessageEmbed{
		    Author:      &discordgo.MessageEmbedAuthor{},
		    Color:       0xFF5700, // reddit orange
		    // Description: desc,
		    // Fields: []*discordgo.MessageEmbedField{
		    //     &discordgo.MessageEmbedField{
		    //         Name:   "I am a field",
		    //         Value:  "I am a value",
		    //         Inline: true,
		    //     },
		    //     &discordgo.MessageEmbedField{
		    //         Name:   "I am a second field",
		    //         Value:  "I am a value",
		    //         Inline: true,t
		    //     },
		    // },
		    Image: &discordgo.MessageEmbedImage{
		        URL: url,
		    },
		    // Thumbnail: &discordgo.MessageEmbedThumbnail{
		    //     URL: url,
		    // },
		    Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		    Title:     title,
		}

		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return


	}
}

// /r handler
func rHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == BotID {
		return 
	}

	subredditRE, _ := regexp.Compile("^/r .*")
	if subredditRE.MatchString(m.Content) {
		args := strings.Split(m.Content, " ")
		commentNum := 0

		if len(args) > 2 { // extract comment num from command args
			if temp, err := strconv.Atoi(args[2]); err == nil {
				commentNum = temp
			}
		}

		
		log.Printf("reddit session: %v", redditSession)
		// log.Println(redditSession)
		status, title, url, desc, flair, comments := reddit.R(redditSession, args[1], commentNum)

		if !status {
			s.ChannelMessageSend(m.ChannelID, "something wrong :(")
			return 
		}

		embed := &discordgo.MessageEmbed{
		    Author:      &discordgo.MessageEmbedAuthor{},
		    Color:       0xFF5700, // reddit orange
		    Description: desc,
		    // Fields: []*discordgo.MessageEmbedField{
		    //     &discordgo.MessageEmbedField{
		    //         Name:   "I am a field",
		    //         Value:  "I am a value",
		    //         Inline: true,
		    //     },
		    //     &discordgo.MessageEmbedField{
		    //         Name:   "I am a second field",
		    //         Value:  "I am a value",
		    //         Inline: true,t
		    //     },
		    // },
		    Image: &discordgo.MessageEmbedImage{
		        URL: url,
		    },
		    // Thumbnail: &discordgo.MessageEmbedThumbnail{
		    //     URL: url,
		    // },
		    Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		    Title:     title,
		}

		s.ChannelMessageSendEmbed(m.ChannelID, embed)

		if flair != "flair" {
			s.ChannelMessageSend(m.ChannelID, flair)
		}

		if desc != "desc" {
			s.ChannelMessageSend(m.ChannelID, desc)
		}

		if comments != "empty" {
			s.ChannelMessageSend(m.ChannelID, comments)
		}
		return
	}
}

// /ask handler
func askHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == BotID {
		return 
	}

	askredditRE, _ := regexp.Compile("^/ask")
	if askredditRE.MatchString(m.Content) {

		log.Printf("reddit session: %v", redditSession)
		// log.Println(redditSession)
		status, title, desc, comments := reddit.Ask(redditSession)

		if !status {
			s.ChannelMessageSend(m.ChannelID, "something wrong :(")
			return
		}


		s.ChannelMessageSend(m.ChannelID, title)
		s.ChannelMessageSend(m.ChannelID, desc)
		s.ChannelMessageSend(m.ChannelID, comments)
		return
	}

}

