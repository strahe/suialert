package discord

import (
	"fmt"
	"time"

	"github.com/strahe/suialert/service"

	"github.com/bwmarrin/discordgo"

	"github.com/samber/lo"

	"github.com/strahe/suialert/model"

	"github.com/strahe/suialert/types"

	"go.uber.org/zap"
)

func (b *Bot) handleAlert(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	opt, ok := lo.Find[*discordgo.ApplicationCommandInteractionDataOption](options,
		func(item *discordgo.ApplicationCommandInteractionDataOption) bool {
			if item.Name == "address" {
				return true
			}
			return false
		})
	if !ok {
		// address must be set
		return
	}

	u, err := b.findOrCreateUser(i)
	if err != nil {
		zap.S().Errorf("find user failed: %v", err)
		return
	}

	ao := types.HexToAddress(opt.StringValue())

	rule, err := b.ruleService.FindByAddress(u.ID, ao.Hex())
	//if err != nil && err != service.ErrNotFound {
	//	zap.S().Errorf("find rule failed: %v", err)
	//	return
	//}
	switch err {
	case nil:
		zap.S().Infof("rule already exists: %s", rule.Address)
	case service.ErrNotFound:
		// todo
	default:
		zap.S().Errorf("find rule failed: %v", err)
	}

	if rule != nil {
		zap.S().Infof("rule already exists: %s", rule.Address)
	}

	responses := map[discordgo.Locale]string{
		discordgo.ChineseCN: fmt.Sprintf("账户 %s 已成功添加到你的监控列表中", ao.Hex()),
	}
	response := fmt.Sprintf("The account %s is added successfully to your monitored list", ao.Hex())
	if r, ok := responses[i.Locale]; ok {
		response = r
	}

	placeholders := map[discordgo.Locale]string{
		discordgo.ChineseCN: "请选择你喜欢的告警类型.",
	}
	placeholder := "Please select the alerts you like to receive."
	if r, ok := placeholders[i.Locale]; ok {
		placeholder = r
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							// Select menu, as other components, must have a customID, so we set it to this value.
							CustomID:    "select-alert",
							Placeholder: placeholder,
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
