package discord

import (
	"github.com/bwmarrin/discordgo"
)

var commands = []discordgo.ApplicationCommand{
	{
		Name:        "add-alert",
		Description: "Add a new address to the alert list",
	},
}
