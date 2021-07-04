package main

import (
	"github.com/aryuuu/disbot/bot"
)

func main() {

	// start bot
	bot.Start()

	<-make(chan struct{})
}
