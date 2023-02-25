package discord

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func (b *Bot) handleStart(s *discordgo.Session, i *discordgo.InteractionCreate) {
	responses := map[discordgo.Locale]string{
		discordgo.ChineseCN: "你好！ 这是一个本地化的命令",
	}
	response := "Hi! This is a localized message"
	if r, ok := responses[i.Locale]; ok {
		response = r
	}
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
	zap.S().Error(err)
}
