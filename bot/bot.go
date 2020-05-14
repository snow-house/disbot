package bot

import (
	"github.com/bwmarrin/discordgo"

	"fmt"
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

	nimRE *regexp.Regexp

	publicTagRE *regexp.Regexp
	guildTagRE *regexp.Regexp
	channelTagRE *regexp.Regexp
	addTagRE *regexp.Regexp
	listTagRE *regexp.Regexp
	deleteTagRE *regexp.Regexp
	infoTagRE *regexp.Regexp

	subredditRE *regexp.Regexp
	askredditRE *regexp.Regexp
	randomRE *regexp.Regexp
)

func init() {
	nimRE, _ = regexp.Compile("^/nim .*")

	publicTagRE, _ = regexp.Compile(`"([^"]+)"`)
	guildTagRE, _ = regexp.Compile(":([^:])+:")
	channelTagRE, _ = regexp.Compile(";([^;])+;")
	addTagRE, _ = regexp.Compile("/addtag .*")
	listTagRE, _ = regexp.Compile("^/taglist")
	deleteTagRE, _ = regexp.Compile("^/deletetag .*")
	infoTagRE, _ = regexp.Compile("^/taginfo .*")

	subredditRE, _ = regexp.Compile("^/r .*")
	askredditRE, _ = regexp.Compile("^/ask")
	randomRE, _ = regexp.Compile("^/random")
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
		fmt.Println("Error opening Discord session: ", err)
	}

	fmt.Println("Bot is running")
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

	if randomRE.MatchString(m.Content) {

		status, title, url := reddit.Random()
		if !status {
			s.ChannelMessageSend(m.ChannelID, "something wrong :(")
			return
		}

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


	if subredditRE.MatchString(m.Content) {
		args := strings.Split(m.Content, " ")
		commentNum := 0

		if len(args) > 2 { // extract comment num from command args
			if temp, err := strconv.Atoi(args[2]); err == nil {
				commentNum = temp
			}
		}

		status, title, url, desc, flair, comments := reddit.R(args[1], commentNum)

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

	if askredditRE.MatchString(m.Content) {

		status, title, desc, comments := reddit.Ask()

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

