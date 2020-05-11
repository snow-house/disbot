package main

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"os"
)

var (
	Token string
	bot *discordgo.Session
	BotID string
)

func main() {
	Token = os.Getenv("DISBOTTOKEN")
	bot, err := discordgo.New("Bot " + Token)

	if err != nil {
		fmt.Println("error creating discord session,", err)
		return
	}

	bot.AddHandler(messageHandler)

	err = bot.Open()

	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running....")

	<-make(chan struct{})


	// close discord session
	bot.Close()

}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	if m.Content == "fuck" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "yeah")
	}
}




