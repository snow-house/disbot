package bot

import (
	"github.com/bwmarrin/discordgo"

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
	"../help"
)

var (

	Token string
	// BotID string
	bot *discordgo.Session

	deleteTagRE *regexp.Regexp
	infoTagRE *regexp.Regexp

)

func init() {
		
	deleteTagRE, _ = regexp.Compile("^/deletetag .*")
	infoTagRE, _ = regexp.Compile("^/taginfo .*")
		
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

	// help
	bot.AddHandler(helpHandler)

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
	if m.Author.Bot {
		return
	}

	if m.Content == "fuck" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "yeah")
	}
}

func helpHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	
	if m.Author.Bot {
		return
	}

	helpRE, _ := regexp.Compile("^/help(\\s.+)?")
	if helpRE.MatchString(m.Content) {

		command := "default"

		args := strings.Split(m.Content, " ")
		if len(args) > 1 {
			command = args[1]
		}

		reply := help.Show(command)
		// s.ChannelMessageSend(m.ChannelID, reply)
		embed := &discordgo.MessageEmbed{
		    Author:      &discordgo.MessageEmbedAuthor{},
		    Color:       0x0088de, // blue cola
		    Description: reply,
		   
		    Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		    Title:     "help "+ command,
		}

		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}
}

// nim command handler
func nimHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.Bot {
		return
	}

	nimRE, _ := regexp.Compile("^/nim .*")
	// if message match with nim command
	if nimRE.MatchString(m.Content) {
		
		// extract query
		query := strings.Replace(nimRE.FindString(m.Content), "/nim ", "", -1)
		
		// find nim or name
		name, tpb, s1, major := nim.Find(query)

		embed := &discordgo.MessageEmbed{
		    Author:      &discordgo.MessageEmbedAuthor{},
		    Color:       0x0088de, // blue cola
		    // Description: desc,
		    Fields: []*discordgo.MessageEmbedField{
		        &discordgo.MessageEmbedField{
		            Name:   "Nama",
		            Value:  name,
		            Inline: true,
		        },
		        &discordgo.MessageEmbedField{
		            Name:   "TPB",
		            Value:  tpb,
		            Inline: true,
		        },
		        &discordgo.MessageEmbedField{
		            Name:   "S1",
		            Value:  s1,
		            Inline: true,
		        },
		        &discordgo.MessageEmbedField{
		            Name:   "Jurusan",
		            Value:  major,
		            Inline: true,
		        },
		    },

		    Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		    Title:     "result for "+ query,
		}

		s.ChannelMessageSendEmbed(m.ChannelID, embed)
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
	    
	    Image: &discordgo.MessageEmbedImage{
	        URL: url,
	    },
	    Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
	    Title:     name,
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return
}

// add tag handler
func addTagHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	
	if m.Author.Bot {
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

	if m.Author.Bot {
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

	if m.Author.Bot {
		return
	}

	randomRE, _ := regexp.Compile("^/random")
	if randomRE.MatchString(m.Content) {

		// log.Printf("[randomHandler] reddit session: %v", redditSession)
		// log.Println(redditSession)
		status, title, url := reddit.Random()
		if !status {
			s.ChannelMessageSend(m.ChannelID, "something wrong :(")
			return
		}

		log.Println("fetched post from dankmemes")
		embed := &discordgo.MessageEmbed{
		    Author:      &discordgo.MessageEmbedAuthor{},
		    Color:       0xFF5700, // reddit orange
		    Image: &discordgo.MessageEmbedImage{
		        URL: url,
		    },
		    Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		    Title:     title,
		}

		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return


	}
}

// /r handler
func rHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.Bot {
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
		    //         Name:   flair,
		    //         // Value:  "",
		    //         Inline: false,
		    //     },
		    // },
		    Image: &discordgo.MessageEmbedImage{
		        URL: url,
		    },
		    Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		    Title:     title,
		}

		s.ChannelMessageSendEmbed(m.ChannelID, embed)

		if flair != "flair" {
			s.ChannelMessageSend(m.ChannelID, flair)
		}

		if comments != "empty" {
			s.ChannelMessageSend(m.ChannelID, comments)
		}
		return
	}
}

// /ask handler
func askHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.Bot {
		return 
	}

	askredditRE, _ := regexp.Compile("^/ask")
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



