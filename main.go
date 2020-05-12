package main

import (
	"./bot"
)


func main() {

	// start bot
	bot.Start()

	<-make(chan struct{})
}




