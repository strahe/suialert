package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/samber/lo"

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
		if err := b.returnError(s, i, map[discordgo.Locale]string{
			discordgo.EnglishUS: "cant find user",
		}); err != nil {
			zap.S().Error(err)
		}
		return
	}

	fmt.Println(u.Name)

	ao := types.HexToAddress(opt.StringValue())
	zap.L().Debug("cache address",
		zap.String("id", i.ID),
		zap.String("address", ao.Hex()))
	if err := b.cache.Set(i.ID, []byte(ao.Hex())); err != nil {
		zap.S().Errorf("failed to cache alert address: %v", err)
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
		Type: discordgo.InteractionResponseModal,
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
						},
					},
				},
			},
		}})
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
