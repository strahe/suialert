package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/strahe/suialert/config"
	"go.uber.org/zap"
)

type Bot struct {
	cfg     config.DiscordBotConfig
	session *discordgo.Session

	cmdIDs map[string]string
}

func New(cfg config.DiscordBotConfig) (*Bot, error) {
	ss, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %s", err)
	}
	bot := &Bot{
		cfg:     cfg,
		session: ss,
		cmdIDs:  map[string]string{},
	}
	bot.addHandlers()
	return bot, nil
}

func (b *Bot) Run() error {
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("failed to open session: %s", err)
	}
	return b.createCommands()
}

// Close closes the bot.
func (b *Bot) Close() error {
	if b.session != nil {
		return b.session.Close()
	}

	for id, name := range b.cmdIDs {
		err := b.session.ApplicationCommandDelete(b.cfg.AppID, "", id)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", name, err)
		}
	}

	return nil
}

func (b *Bot) addHandlers() {
	b.session.AddHandler(b.handleReady)
	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand, discordgo.InteractionApplicationCommandAutocomplete:
			switch i.ApplicationCommandData().Name {
			case "start":
				b.handleStart(s, i)
			case "add-address":
				b.handleAddAddress(s, i)
			default:
				zap.S().Errorf("Unknown slash command: %s", i.ApplicationCommandData().Name)
			}
		case discordgo.InteractionMessageComponent:
			switch i.MessageComponentData().CustomID {
			case "select-events":
				b.handleSelectEventType(s, i)
			}
		case discordgo.InteractionModalSubmit:
			data := i.ModalSubmitData()
			zap.S().Infof("Modal submit: %s", data.CustomID)
		default:
			zap.S().Errorf("Unknown slash command: %s", i.Type)
		}
	})
}

func (b *Bot) createCommands() error {
	for _, cmd := range commands {
		rc, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", &cmd)
		if err != nil {
			return fmt.Errorf("failed to create slash command: %s", err)
		}
		b.cmdIDs[rc.ID] = rc.Name
	}
	return nil
}
