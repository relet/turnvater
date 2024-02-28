package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func TurnResetHandler(dg *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if the user has the correct permissions

	if !HasPermission(dg, i.Member, i.GuildID, "ADMINISTRATOR") {
		Respond(dg, i, i18n[lang]["err-not-allowed"])
		return
	}
	// Reset the tournament
	name := i.ApplicationCommandData().Options[0].StringValue()
	err := DBResetTournament(backend, name)
	if err != nil {
		Respond(dg, i, i18n[lang]["err-reset"])
		return
	}
	Respond(dg, i, fmt.Sprintf(i18n[lang]["ok-reset"], name))
}

func TurnRegisterHandler(dg *discordgo.Session, i *discordgo.InteractionCreate) {
	// check if registration is open
	status := DBGetTournamentStatus(backend)
	if status != "status-open" {
		Respond(dg, i, i18n[lang]["err-started"])
		return
	}

	// get discord handle, name and IGN parameter
	ign := i.ApplicationCommandData().Options[0].StringValue()
	// if ign begins with ! reject
	if ign[0] == '!' {
		Respond(dg, i, i18n[lang]["err-register-name"])
		return
	}

	err := DBRegisterParticipant(backend, i.Member.User.ID, ign)
	if err != nil {
		Respond(dg, i, i18n[lang]["err-register"]+" "+err.Error())
		return
	}

	// Register a participant
	Respond(dg, i, fmt.Sprintf(i18n[lang]["welcome"], i.Member.User.Username, ign))
}

func TurnStatusHandler(dg *discordgo.Session, i *discordgo.InteractionCreate) {
	// Get the status of the tournament
	status := DBGetTournamentStatus(backend)
	name := DBGetTournamentName(backend)
	message := "**" + name + "**\n"
	message += i18n[lang]["turn-status"] + ": " + i18n[lang][status] + "\n\n"
	participants := DBGetParticipants(backend, 0)
	num := len(participants)
	message += fmt.Sprintf(i18n[lang]["status-number"], num) + "\n"

	if status == "status-open" {
		message += strings.Join(participants, ", ") + "\n"
		message += "\n" + i18n[lang]["info-register"]
		message += "\n"
		// show grouping info for sizes 2 to 6
		for i := 2; i <= 6; i++ {
			groups, rest := CalcGroups(num, i)
			if groups > 1 {
				if groups&(groups-1) != 0 {
					reduced := 1
					for reduced < groups {
						reduced *= 2
					}
					reduced /= 2
					rest += (groups - reduced) * i
					groups = reduced
				}
				if rest < groups {
					message += fmt.Sprintf(i18n[lang]["info-grouping"], i, groups, rest) + ".\n"
				}
			}
		}
	} else if status == "status-started" {
		// show grouping info
		groups := DBGetGroups(backend)
		for _, g := range groups { 
			message += fmt.Sprintf("%s: %s\n", g.Name, strings.Join(g.Participants, ", "))
			//print matches
			matches := DBGetMatches(backend, g.Id)
			for _, m := range matches {
				if m.Score1 > 0 || m.Score2 > 0 {
					message += fmt.Sprintf("\t%s vs %s: %d-%d\n", m.Player1, m.Player2, m.Score1, m.Score2)
				}
			}
		}
	} else if status == "status-finished" {
		winner := DBGetTournamentWinner(backend)
		message += "*" + i18n[lang]["tournament-winner"] + "*: " + winner + "\n"
	}
	Respond(dg, i, message)
}

func TurnStartHandler(dg *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if the user has the correct permissions
	if !HasPermission(dg, i.Member, i.GuildID, "ADMINISTRATOR") {
		Respond(dg, i, i18n[lang]["err-not-allowed"])
		return
	}
	// Check if the tournament is not yet started
	status := DBGetTournamentStatus(backend)
	if status != "status-open" {
		Respond(dg, i, i18n[lang]["err-start"])
		return
	}

	// Start the tournament
	groupsize := i.ApplicationCommandData().Options[0].IntValue()
	bestof := i.ApplicationCommandData().Options[1].IntValue()
	finals := i.ApplicationCommandData().Options[2].IntValue()
	err := DBStartTournament(backend, groupsize, bestof, finals)
	if err != nil {
		Respond(dg, i, i18n[lang]["err-start"]+" "+err.Error())
		return
	}
	Respond(dg, i, fmt.Sprintf(i18n[lang]["ok-start"], groupsize))

	ReRegister()
}

func TurnResultHandler(dg *discordgo.Session, i *discordgo.InteractionCreate) {
	status := DBGetTournamentStatus(backend)
	if status != "status-started" {
		Respond(dg, i, i18n[lang]["err-not-started"])
		return
	}
	p1 := i.ApplicationCommandData().Options[0].StringValue()
	score1 := i.ApplicationCommandData().Options[1].IntValue()
	p2 := i.ApplicationCommandData().Options[2].StringValue()
	score2 := i.ApplicationCommandData().Options[3].IntValue()

	if p1 == p2 {
		Respond(dg, i, i18n[lang]["err-no-match"])
		return
	}

	group, bestof := DBGetGroupAndBestOf(backend, p1, p2)
	if group.Id == 0 {
		Respond(dg, i, i18n[lang]["err-no-match"])
		return
	}

	// check if the scores add up to the best-of value
	if int(score1+score2) != bestof {
		Respond(dg, i, fmt.Sprintf(i18n[lang]["err-score-total"], bestof))
		return
	}

	err := DBCreateMatch(backend, p1, p2, score1, score2)
	if err != nil {
		Respond(dg, i, i18n[lang]["err-set-score"]+" "+err.Error())
		return
	}

	message := i18n[lang]["ok-set-score"] + " " + p1 + " vs " + p2 + ": " + fmt.Sprintf("%d-%d", score1, score2)

	// check if this concludes the group
	winners, err := DBCheckGroupComplete(backend, group.Id)
	if err != nil {
		Respond(dg, i, message+"\n\n"+i18n[lang]["err-group-complete"]+" "+err.Error())
		return
	}
	if winners != nil {
		// check if the tournament has been won
		first := winners[0]
		if first.Group.Id == 0 {
			DBCloseTournament(backend, first.Player)
			Respond(dg, i, message+"\n\n"+fmt.Sprintf(i18n[lang]["congratulate"], first.Player))
			return
		}
		// send a new message informing about the promotion
		message += "\n\n" + fmt.Sprintf(i18n[lang]["ok-group-winner"], first.Player, group.Name, first.Group.Name)
		if len(winners) > 1 {
			second := winners[1]
			message += "\n" + fmt.Sprintf(i18n[lang]["ok-group-second"], second.Player, group.Name, second.Group.Name)
		}
	}
	Respond(dg, i, message)
}

func TurnGamesHandler(dg *discordgo.Session, i *discordgo.InteractionCreate) {
	// check if the tournament is running
	status := DBGetTournamentStatus(backend)
	if status != "status-started" {
		Respond(dg, i, i18n[lang]["err-not-started"])
		return
	}
	// display all games ordered by group
	groups := DBGetAllGames(backend)
	var message string
	for _, group := range groups {
		message := "*" + fmt.Sprintf(i18n[lang]["summary-group"], group.Name) + "*\n"
		for _, match := range group.Matches {
			p1 := match.Player1
			p2 := match.Player2
			if p1[0] != '!' && p2[0] != '!' {
				message += fmt.Sprintf(i18n[lang]["summary-match"], match.Player1, match.Player2, match.Score1, match.Score2) + "\n"
			}
		}
		message += "\n"
	}
	Respond(dg, i, message)
}
