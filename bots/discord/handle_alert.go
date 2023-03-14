package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/strahe/suialert/model"

	"github.com/samber/lo"

	"github.com/strahe/suialert/types"

	"go.uber.org/zap"
)

func (b *Bot) handleAddAlert(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			CustomID: "add-alert-start",
			Flags:    discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    "selected-event",
							Placeholder: "Which type of event would you like to monitor?",
							Options:     buildEventOptions(),
						},
					},
				},
			},
		},
	}, b.options()...)
	if err != nil {
		zap.S().Error(err)
	}
}

func (b *Bot) handSelectedEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()
	event := data.Values[0]
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "add-alert-for-" + event,
			Title:    "Add an alarm of type " + event,
			Content:  "hello",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "address",
							Label:     "What is the SUI address you like to monitor?",
							Style:     discordgo.TextInputShort,
							Required:  true,
							MaxLength: 42,
							MinLength: 10,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "rules",
							Label:     "Please enter the monitoring rules!",
							Style:     discordgo.TextInputParagraph,
							Required:  true,
							MaxLength: 200,
						},
					},
				},
			},
		},
	}, b.options()...)
	if err != nil {
		panic(err)
	}
}

func (b *Bot) handAddAlertFormSubmitted(s *discordgo.Session, i *discordgo.InteractionCreate) {
	md := i.ModalSubmitData()
	if !strings.HasPrefix(md.CustomID, "add-alert-for-") {
		// todo
		return
	}
	u, err := b.findOrCreateUser(i)
	if err != nil {
		// todo
		return
	}
	// todo: check condition invalid
	event := md.CustomID[len("add-alert-for-"):]
	addr := md.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	rule := md.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	err = b.ruleService.Create(&model.Rule{
		Address:   types.HexToAddress(addr),
		Event:     types.EventType(event),
		User:      *u,
		Condition: rule,
	})
	if err != nil {
		zap.S().Infof("failed to create rule: %v", err)
		// todo
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Alert added",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}, b.options()...)
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
		}}, b.options()...)
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
			Label:       string(e),
			Value:       string(e),
			Description: e.Description(),
			Emoji: discordgo.ComponentEmoji{
				Name: e.Emoji(),
			},
		})
	}
	return options
}
