package model

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/telebot.v3"
)

type User struct {
	ID           int64           `json:"id" pg:"id,pk,notnull"`
	Name         string          `json:"name" pg:"name,notnull"`
	DiscordID    *string         `json:"discord_id" pg:"discord_id"`
	DiscordInfo  *discordgo.User `json:"discord_info" pg:"discord_info"`
	TelegramID   *int64          `json:"telegram_id" pg:"telegram_id"`
	TelegramInfo *telebot.User   `json:"telegram_info" pg:"telegram_info"`
	CreatedAt    int64           `json:"created_at" pg:"created_at"`
}

func (e *User) Persist(ctx context.Context, s StorageBatch) error {
	return s.PersistModel(ctx, e)
}
