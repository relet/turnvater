package main

/* A discord bot to manage tournament brackets and participants */

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	_ "modernc.org/sqlite"
)

var (
	allow = true
	deny  = false

	backend   *sql.DB
	turnvater *TurnvaterBot
	permAdmin int64 = discordgo.PermissionAdministrator

	lang = "de"
)

func HasPermission(dg *discordgo.Session, member *discordgo.Member, guildID, permission string) bool {
	for _, roleID := range member.Roles {
		role, err := dg.State.Role(guildID, roleID)
		if err != nil {
			return false
		}
		if role.Permissions&permAdmin == permAdmin {
			return true
		}
	}
	return false
}

func Respond(dg *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

func CalcGroups(num, groupsize int) (int, int) {
	groups := num / groupsize
	rest := num % groupsize
	return groups, rest
}

/* === Main Loop === */

func ReRegister() {
	turnvater.Participants = DBGetParticipants(backend, 0)
	turnvater.Restart = true
}

func main() {
	// read lang from environment or command line
	if os.Getenv("LANG") == "de" || os.Getenv("LANG") == "en" {
		lang = os.Getenv("LANG")
	} else if len(os.Args) > 1 {
		lang = os.Args[1]
	}

	// read settings from .settings file
	settingsJson, err := os.ReadFile(".settings")
	if err != nil {
		fmt.Println("error reading settings", err)
		return
	}
	settings := make(map[string]string)
	err = json.Unmarshal(settingsJson, &settings)
	if err != nil {
		fmt.Println("error parsing settings", err)
		return
	}
	token, ok := settings["token"]
	if !ok {
		fmt.Println("token not found in settings")
		return
	}
	appId, ok := settings["appId"]
	if !ok {
		fmt.Println("appId not found in settings")
		return
	}
	guildId, ok := settings["guildId"]
	if !ok {
		fmt.Println("guildId not found in settings")
		return
	}
	state, ok := settings["state"]
	if !ok {
		fmt.Println("state not found in settings")
		return
	}
	pRoleId, ok := settings["participantRoleId"]
	if !ok {
		fmt.Println("participantRoleId not found in settings")
		return
	}

	db, err := sql.Open("sqlite", state)
	if err != nil {
		fmt.Println("error opening database", err)
		return
	}
	defer db.Close()

	backend = db

	participants := DBGetParticipants(backend, 0)

	bot, err := NewBot(token, appId, guildId, participants, pRoleId)
	if err != nil {
		fmt.Println("error running bot", err)
		return
	}
	turnvater = &bot

	bot.Run()
}
