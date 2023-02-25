package discord

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func (b *Bot) handleReady(s *discordgo.Session, _ *discordgo.InteractionCreate) {
	zap.S().Infof("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
}
