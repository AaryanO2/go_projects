package bot

import (
	"fmt"
	"strings"

	"github.com/AaryanO2/go_projects/project_7_discord_bot/config"
	"github.com/bwmarrin/discordgo"
)

var BotId string
var goBot *discordgo.Session

func Start() {
	var err error
	goBot, err = discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	BotId = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()

	if err != nil {
		return
	}
	fmt.Println("Bot is running.....")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotId {
		return
	}
	fmt.Printf("Message: %v\n", m.Content)
	if strings.ToLower(m.Content) == "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong")
	} else {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Try again")
	}
}
