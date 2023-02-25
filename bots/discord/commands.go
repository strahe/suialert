package discord

import (
	"github.com/bwmarrin/discordgo"
)

var commands = []discordgo.ApplicationCommand{
	{
		Name:        "start",
		Description: "Start command.",
		NameLocalizations: &map[discordgo.Locale]string{
			discordgo.ChineseCN: "开始",
		},
		DescriptionLocalizations: &map[discordgo.Locale]string{
			discordgo.ChineseCN: "开始命令",
		},
	},
	{
		Name:        "add-address",
		Description: "Add a new address.",
		NameLocalizations: &map[discordgo.Locale]string{
			discordgo.ChineseCN: "添加地址",
		},
		DescriptionLocalizations: &map[discordgo.Locale]string{
			discordgo.ChineseCN: "添加新地址",
		},
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "address",
				Description: "Sui network address",
			},
		},
	},
}
