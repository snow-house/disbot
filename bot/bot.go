package bot

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"os"
)

var (
	Token string
	BotID string
	bot *discordgo.Session
)

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
	bot.AddHandler(messageHandler)

	// open the websocket and begin listening
	err = bot.Open()

	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	fmt.Println("Bot is running")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	if m.Content == "fuck" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "yeah")
	}
}
