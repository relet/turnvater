package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type TurnvaterBot struct {
	Token        string
	AppId        string
	GuildId      string
	Participants []string
	Restart      bool
}

var commands = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
	"turn-reset":    TurnResetHandler,
	"turn-register": TurnRegisterHandler,
	"turn-status":   TurnStatusHandler,
	"turn-start":    TurnStartHandler,
	"turn-result":   TurnResultHandler,
}

func GenChoices(choices []string) []*discordgo.ApplicationCommandOptionChoice {
	var result []*discordgo.ApplicationCommandOptionChoice
	for _, choice := range choices {
		result = append(result, &discordgo.ApplicationCommandOptionChoice{
			Name:  choice,
			Value: choice,
		})
	}
	return result
}

func NewBot(token string, appId string, guildId string, participants []string) TurnvaterBot {
	bot := TurnvaterBot{
		Token:        token,
		AppId:        appId,
		GuildId:      guildId,
		Participants: participants,
		Restart:      false,
	}

	return bot
}

func (bot *TurnvaterBot) ReRegisterCommands() error {
	dg, err := discordgo.New("Bot " + bot.Token)
	if err != nil {
		return fmt.Errorf("error creating discord session: %w", err)
	}
	defer dg.Close()

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Printf("Logged in as: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
	})
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		fmt.Println("Interaction", i.Type)
		if i.Type == discordgo.InteractionApplicationCommand {
			if handler, ok := commands[i.ApplicationCommandData().Name]; ok {
				handler(s, i)
			} else {
				fmt.Println("Unknown command", i.ApplicationCommandData().Name)
			}
		}
	})

	// wait until ready
	err = dg.Open()
	if err != nil {
		return fmt.Errorf("error opening discord connection: %w", err)
	}

	// Register slash commands and their handlers
	// /turn-reset

	_, err = dg.ApplicationCommandCreate(bot.AppId, bot.GuildId, &discordgo.ApplicationCommand{
		Name:                     "turn-reset",
		Description:              i18n[lang]["turn-reset"],
		DefaultMemberPermissions: &permAdmin,
		DMPermission:             &allow,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: i18n[lang]["opt-name"],
				Required:    true,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error creating command: %w", err)
	}

	// /turn-register
	_, err = dg.ApplicationCommandCreate(bot.AppId, bot.GuildId, &discordgo.ApplicationCommand{
		Name:         "turn-register",
		Description:  i18n[lang]["turn-register"],
		DMPermission: &deny,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "ign",
				Description: i18n[lang]["opt-ign"],
				Required:    true,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error creating command: %w", err)
	}

	// /turn-status
	_, err = dg.ApplicationCommandCreate(bot.AppId, bot.GuildId, &discordgo.ApplicationCommand{
		Name:         "turn-status",
		Description:  i18n[lang]["turn-status"],
		DMPermission: &allow,
	})
	if err != nil {
		return fmt.Errorf("error creating command: %w", err)
	}

	// /turn-start
	_, err = dg.ApplicationCommandCreate(bot.AppId, bot.GuildId, &discordgo.ApplicationCommand{
		Name:                     "turn-start",
		Description:              i18n[lang]["turn-start"],
		DMPermission:             &allow,
		DefaultMemberPermissions: &permAdmin,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "groupsize",
				Description: i18n[lang]["opt-groupsize"],
				Required:    true,
				Choices:     GenChoices([]string{"2", "3", "4", "5", "6"}),
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "bestof",
				Description: i18n[lang]["opt-bestof"],
				Required:    true,
				Choices:     GenChoices([]string{"1", "3", "5", "7", "9"}),
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "finals-bestof",
				Description: i18n[lang]["opt-finals-bestof"],
				Required:    true,
				Choices:     GenChoices([]string{"1", "3", "5", "7", "9"}),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error creating command: %w", err)
	}

	// /turn-result
	_, err = dg.ApplicationCommandCreate(bot.AppId, bot.GuildId, &discordgo.ApplicationCommand{
		Name:         "turn-result",
		Description:  i18n[lang]["turn-result"],
		DMPermission: &deny,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "p1",
				Description: i18n[lang]["opt-p1"],
				Required:    true,
				Choices:     GenChoices(bot.Participants),
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "score1",
				Description: i18n[lang]["opt-score1"],
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "p2",
				Description: i18n[lang]["opt-p2"],
				Required:    true,
				Choices:     GenChoices(bot.Participants),
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "score2",
				Description: i18n[lang]["opt-score2"],
				Required:    true,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error creating command: %w", err)
	}

	fmt.Println("Commands registered.")

	return nil
}

func (bot *TurnvaterBot) Run() {
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	for {
		if bot.Restart {
			bot.Restart = false
			bot.ReRegisterCommands()
		}
	}
}
