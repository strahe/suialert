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
		if err := b.returnError(s, i, map[discordgo.Locale]string{
			discordgo.EnglishUS: "cant find user",
		}); err != nil {
			zap.S().Error(err)
		}
		return
	}

	ao := types.HexToAddress(opt.StringValue())
	zap.L().Debug("cache address",
		zap.String("id", i.ID),
		zap.String("address", ao.Hex()))
	if err := b.cache.Set(i.ID, []byte(ao.Hex())); err != nil {
		zap.S().Errorf("failed to cache alert address: %v", err)
	}

	rule, err := b.ruleService.FindByAddress(u.ID, ao.Hex())
	switch err {
	case nil:
		zap.S().Infof("rule already exists: %s", rule.Address)
	case service.ErrNotFound:
		if err := b.ruleService.Create(&model.Rule{
			Address:    ao.Hex(),
			UserID:     u.ID,
			AlertLevel: model.AlertLevelNone,
			CreatedAt:  time.Now().Unix(),
			UpdatedAt:  time.Now().Unix(),
		}); err != nil {
			err = b.returnError(s, i, map[discordgo.Locale]string{
				discordgo.EnglishUS: "cant create rule",
			})
			zap.S().Errorf("cant create rule: %s", err)
		}
	default:
		if err := b.returnError(s, i, map[discordgo.Locale]string{
			discordgo.EnglishUS: "address already exists",
		}); err != nil {
			zap.S().Error(err)
		}
		return
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
	u, err := b.findOrCreateUser(i)
	if err != nil {
		zap.S().Errorf("send error msg: %s", b.returnError(s, i, map[discordgo.Locale]string{
			discordgo.EnglishUS: "cant find user",
		}))
		return
	}

	var content string
	data := i.MessageComponentData()
	for _, v := range alertLevelOptions {
		if data.Values[0] == v.Value {
			content = fmt.Sprintf("You have selected %s\n"+
				"%s", v.Label, v.Description)
			break
		}
	}
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}

	v, err := b.cache.Get(i.Message.Interaction.ID)
	if err != nil {
		zap.S().Error(err)
	}
	zap.L().Debug("get cache address",
		zap.String("ID", i.Interaction.ID),
		zap.String("address", string(v)),
	)

	rule, err := b.ruleService.FindByAddress(u.ID, string(v))
	if err != nil {
		zap.S().Errorf("cant find rule: %s", err)
		b.returnError(s, i, map[discordgo.Locale]string{
			discordgo.EnglishUS: "cant find rule",
		})
		return
	}
	rule.AlertLevel = model.AlertLevel(data.Values[0])
	rule.UpdatedAt = time.Now().Unix()
	if err := b.ruleService.Update(rule); err != nil {
		zap.S().Errorf("cant update rule: %s", err)
		b.returnError(s, i, map[discordgo.Locale]string{
			discordgo.EnglishUS: "cant update rule",
		})
		return
	}

	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		zap.S().Error(err)
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
