package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/samber/lo"

	"github.com/strahe/suialert/types"

	"go.uber.org/zap"
)

func (b *Bot) handleAlert(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			CustomID: "add-alert",
			Title:    "Add an alert",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    "select-event",
							Placeholder: "Which type of event would you like to monitor?",
							Options:     buildEventOptions(),
						},
					},
				},
			},
		},
	})
	if err != nil {
		zap.S().Error(err)
	}
}

func (b *Bot) returnError(s *discordgo.Session, i *discordgo.InteractionCreate, msgs map[discordgo.Locale]string) error {
	if len(msgs) == 0 {
		return fmt.Errorf("no messages")
	}
	var content string
	msg, ok := msgs[i.Locale]
	if ok {
		content = msg
	} else if msg, ok = msgs[discordgo.EnglishUS]; ok {
		content = msg
	} else {
		content = lo.Values[discordgo.Locale, string](msgs)[0]
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Title:   "Error",
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		}})
}

func buildEventOptions() []discordgo.SelectMenuOption {
	var (
		events = []types.EventType{
			types.EventTypeMove,
			types.EventTypePublish,
			types.EventTypeCoinBalanceChange,
			types.EventTypeTransferObject,
			types.EventTypeNewObject,
			types.EventTypeDeleteObject,
			types.EventTypeMutateObject,
		}
		options []discordgo.SelectMenuOption
	)

	for _, e := range events {
		options = append(options, discordgo.SelectMenuOption{
			Label:       e.Name(),
			Value:       e.Name(),
			Description: e.Description(),
			Emoji: discordgo.ComponentEmoji{
				Name: e.Emoji(),
			},
		})
	}
	return options
}
