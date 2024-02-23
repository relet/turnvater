package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var commands = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
	"turn-reset":    TurnResetHandler,
	"turn-register": TurnRegisterHandler,
	"turn-status":   TurnStatusHandler,
	"turn-start":    TurnStartHandler,
	"turn-result":   TurnResultHandler,
}

func RunBot(token string, appId string, guildId string) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session", err)
		return
	}
	defer dg.Close()

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Printf("Logged in as: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
	})
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			if handler, ok := commands[i.ApplicationCommandData().Name]; ok {
				handler(s, i)
			}
		}
	})

	// wait until ready
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection", err)
		return
	}

	// Register slash commands and their handlers
	// /turn-reset
	_, err = dg.ApplicationCommandCreate(appId, guildId, &discordgo.ApplicationCommand{
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
		fmt.Println("error creating command:", err)
	}

	// /turn-register
	_, err = dg.ApplicationCommandCreate(appId, guildId, &discordgo.ApplicationCommand{
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
		fmt.Println("error creating command:", err)
	}

	// /turn-status
	_, err = dg.ApplicationCommandCreate(appId, guildId, &discordgo.ApplicationCommand{
		Name:         "turn-status",
		Description:  i18n[lang]["turn-status"],
		DMPermission: &allow,
	})
	if err != nil {
		fmt.Println("error creating command:", err)
	}

	// /turn-start
	_, err = dg.ApplicationCommandCreate(appId, guildId, &discordgo.ApplicationCommand{
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
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "bestof",
				Description: i18n[lang]["opt-bestof"],
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "finals-bestof",
				Description: i18n[lang]["opt-finals-bestof"],
				Required:    true,
			},
		},
	})
	if err != nil {
		fmt.Println("error creating command:", err)
	}

	// /turn-result
	_, err = dg.ApplicationCommandCreate(appId, guildId, &discordgo.ApplicationCommand{
		Name:         "turn-result",
		Description:  i18n[lang]["turn-result"],
		DMPermission: &deny,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "p1",
				Description: i18n[lang]["opt-p1"],
				Required:    true,
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
		fmt.Println("error creating command:", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	select {}
}
