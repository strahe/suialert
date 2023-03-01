package discord

import (
	"time"

	"github.com/strahe/suialert/types"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func (b *Bot) handleAddAddress(s *discordgo.Session, i *discordgo.InteractionCreate) {
	responses := map[discordgo.Locale]string{
		discordgo.ChineseCN: "选择你想要接收通知的事件类型",
	}
	response := "Choose the type of event you want to add: "
	if r, ok := responses[i.Locale]; ok {
		response = r
	}
	minValue := 1
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							// Select menu, as other components, must have a customID, so we set it to this value.
							CustomID:    "select-events",
							Placeholder: response,
							MinValues:   &minValue,
							MaxValues:   len(eventSelectOptions),
							Options:     eventSelectOptions,
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

var eventSelectOptions = []discordgo.SelectMenuOption{
	{
		Label:       types.EventTypeMoveEvent.Name(),
		Value:       types.EventTypeMoveEvent.Name(),
		Description: types.EventTypeMoveEvent.Description(),
	},
	{
		Label:       types.EventTypePublish.Name(),
		Value:       types.EventTypePublish.Name(),
		Description: types.EventTypePublish.Description(),
	},
	{
		Label:       types.EventTypeCoinBalanceChange.Name(),
		Value:       types.EventTypeCoinBalanceChange.Name(),
		Description: types.EventTypeCoinBalanceChange.Description(),
	},
	{
		Label:       types.EventTypeTransferObject.Name(),
		Value:       types.EventTypeTransferObject.Name(),
		Description: types.EventTypeTransferObject.Description(),
	},
	{
		Label:       types.EventTypeNewObject.Name(),
		Value:       types.EventTypeNewObject.Name(),
		Description: types.EventTypeNewObject.Description(),
	},
	{
		Label:       types.EventTypeMutateObject.Name(),
		Value:       types.EventTypeMutateObject.Name(),
		Description: types.EventTypeMutateObject.Description(),
	},
	{
		Label:       types.EventTypeDeleteObject.Name(),
		Value:       types.EventTypeDeleteObject.Name(),
		Description: types.EventTypeDeleteObject.Description(),
	},
}
