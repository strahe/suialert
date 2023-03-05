package discord

import (
	"github.com/bwmarrin/discordgo"
)

var commands = []discordgo.ApplicationCommand{
	{
		Name:        "alert",
		Description: "Add a new address to the alert list",
		NameLocalizations: &map[discordgo.Locale]string{
			discordgo.ChineseCN: "监控",
		},
		DescriptionLocalizations: &map[discordgo.Locale]string{
			discordgo.ChineseCN: "添加一个地址到监控列表",
		},
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "address",
				Description: "Sui network address",
				Required:    true,
			},
		},
	},
}
