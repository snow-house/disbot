package bot

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"regexp"
	"os"
	"strings"
	"../nim"
)

var (
	Token string
	BotID string
	bot *discordgo.Session
	nimRE *regexp.Regexp

)

func init() {
	nimRE, _ = regexp.Compile("^/nim .*")
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
	//basic ping 
	bot.AddHandler(fuckHandler)
	bot.AddHandler(nimHandler)

	// open the websocket and begin listening
	err = bot.Open()

	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	fmt.Println("Bot is running")
}

func fuckHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	if m.Content == "fuck" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "yeah")
	}
}

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
