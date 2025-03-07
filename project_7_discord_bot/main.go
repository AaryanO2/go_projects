package main

import (
	"fmt"

	"github.com/AaryanO2/go_projects/project_7_discord_bot/bot"
	"github.com/AaryanO2/go_projects/project_7_discord_bot/config"
)

func main() {
	err := config.ReadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start()

	<-make(chan struct{})

}
