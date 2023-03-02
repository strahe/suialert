package model

import (
	"github.com/bwmarrin/discordgo"
	"gopkg.in/telebot.v3"
)

type User struct {
	ID           int64           `json:"id" pg:"id,pk,notnull"`
	Name         string          `json:"name" pg:"name,notnull"`
	DiscordID    *string         `json:"discord_id" pg:"discord_id,pk"`
	DiscordInfo  *discordgo.User `json:"discord_info" pg:"discord_info"`
	TelegramID   *int64          `json:"telegram_id" pg:"telegram_id,pk"`
	TelegramInfo *telebot.User   `json:"telegram_info" pg:"telegram_info"`
	CreatedAt    int64           `json:"created_at" pg:"created_at"`
}
