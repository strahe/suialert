package model

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/telebot.v3"
)

type User struct {
	ID           uint            `json:"id" gorm:"primaryKey"`
	Name         string          `json:"name"`
	DiscordID    *string         `json:"discord_id" gorm:"index"`
	DiscordInfo  *discordgo.User `json:"discord_info" gorm:"serializer:json"`
	TelegramID   *int64          `json:"telegram_id" gorm:"index"`
	TelegramInfo *telebot.User   `json:"telegram_info" gorm:"serializer:json"`
	RuleCount    int             `json:"rule_count" gorm:"->"` // read only
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

func (*User) TableName() string {
	return "users"
}
