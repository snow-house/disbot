package bot

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"regexp"
	"os"
	"time"
	"strings"
	"../nim"
	"../tag"
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
)

func init() {
	nimRE, _ = regexp.Compile("^/nim .*")

	publicTagRE, _ = regexp.Compile("#([^#]+)#")
	guildTagRE, _ = regexp.Compile("\\|([^$])+\\|")
	channelTagRE, _ = regexp.Compile(";([^;])+;")
	addTagRE, _ = regexp.Compile("/addtag .*")
	listTagRE, _ = regexp.Compile("^/taglist")
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

	bot.AddHandler(nimHandler)

	bot.AddHandler(getTag)
	bot.AddHandler(addTag)
	bot.AddHandler(deleteTag)
	bot.AddHandler(listTag)
	bot.AddHandler(infoTag)


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
func getTag(s *discordgo.Session, m *discordgo.MessageCreate) {

	var name string
	var scope int

	if m.Author.Bot {
		return
	}

	if publicTagRE.MatchString(m.Content) {
		name = strings.Replace(publicTagRE.FindString(m.Content), "#", "", -1)
		scope = 2
		
	} else if guildTagRE.MatchString(m.Content) {
		name = strings.Replace(guildTagRE.FindString(m.Content), "|", "", -1)
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
func addTag(s *discordgo.Session, m *discordgo.MessageCreate) {
	
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
func deleteTag(s *discordgo.Session, m *discordgo.MessageCreate) {
	
}

// list tag handler
func listTag(s *discordgo.Session, m *discordgo.MessageCreate) {

	if listTagRE.MatchString(m.Content) {
		result :=  tag.List(m.ChannelID, m.GuildID)

		_, _ = s.ChannelMessageSend(m.ChannelID, result)
	}
}

// info tag handler
func infoTag(s *discordgo.Session, m *discordgo.MessageCreate) {
	
}