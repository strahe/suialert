package discord

import (
	"time"

	"github.com/strahe/suialert/model"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func (b *Bot) handleAddAddress(s *discordgo.Session, i *discordgo.InteractionCreate) {
	responses := map[discordgo.Locale]string{
		discordgo.ChineseCN: "选择你喜欢的告警类型.",
	}
	response := "Select the alerts you like to receive."
	if r, ok := responses[i.Locale]; ok {
		response = r
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							// Select menu, as other components, must have a customID, so we set it to this value.
							CustomID:    "select-events",
							Placeholder: response,
							Options:     alertLevelOptions,
						},
					},
				},
			},
		}})
	if err != nil {
		zap.S().Error(err)
	}
}

func (b *Bot) handleSelectEventType(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var response *discordgo.InteractionResponse

	data := i.MessageComponentData()
	switch data.Values[0] {
	case "go":
		response = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This is the way.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		}
	default:
		response = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "It is not the way to go.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		}
	}
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second) // Doing that so user won't see instant response.
	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "Anyways, now when you know how to use single select menus, let's see how multi select menus work. " +
			"Try calling `/selects multi` command.",
		Flags: discordgo.MessageFlagsEphemeral,
	})
	if err != nil {
		panic(err)
	}
}

var alertLevelOptions = []discordgo.SelectMenuOption{
	{
		Label:       "None",
		Value:       string(model.AlertLevelNone),
		Description: "Do not send any alerts for now",
	},
	{
		Label:       "Low",
		Value:       string(model.AlertLevelLow),
		Description: "Alert only SUI transfers",
	},
	{
		Label:       "Medium",
		Value:       string(model.AlertLevelMedium),
		Description: "All, excluding popular dapps, and < 0.1 SUI",
	},
	{
		Label:       "High",
		Value:       string(model.AlertLevelHigh),
		Description: "All, excluding < 0.1 SUI",
	},
	{
		Label:       "All",
		Value:       string(model.AlertLevelAll),
		Description: "Send me all the alerts",
	},
}
