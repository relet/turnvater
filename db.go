package main

import (
	"database/sql"
	"fmt"
	"math/rand"
)

type Match struct {
	Player1 string
	Player2 string
	Score1  int
	Score2  int
}

type Group struct {
	Id           int
	Name         string
	Participants []string
}

func DBResetTournament(db *sql.DB, name string) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS options (
		id INTEGER PRIMARY KEY,
		key TEXT NOT NULL,
		value TEXT NOT NULL
	)`)
	if err != nil {
		fmt.Println("error creating options table:", err)
		return err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS participants (
		discord_id TEXT PRIMARY KEY NOT NULL,
		ign TEXT UNIQUE NOT NULL,
		group_id INTEGER DEFAULT 0
	);`)
	if err != nil {
		fmt.Println("error creating participants table:", err)
		return err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS groups (id INTEGER PRIMARY KEY, name TEXT NOT NULL, complete INTEGER DEFAULT 0)")
	if err != nil {
		fmt.Println("error creating groups table:", err)
		return err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS matches (id INTEGER PRIMARY KEY, group_id INTEGER NOT NULL, bestof INTEGER NOT NULL, player1 TEXT NOT NULL, player2 TEXT NOT NULL, score1 INTEGER DEFAULT 0, score2 INTEGER DEFAULT 0)")
	if err != nil {
		fmt.Println("error creating matches table:", err)
		return err
	}

	// reset
	_, err = db.Exec("DELETE FROM participants")
	if err != nil {
		fmt.Println("error deleting participants:", err)
		return err
	}
	_, err = db.Exec("DELETE FROM options")
	if err != nil {
		fmt.Println("error deleting options:", err)
		return err
	}
	_, err = db.Exec("DELETE FROM groups")
	if err != nil {
		fmt.Println("error deleting groups:", err)
		return err
	}
	_, err = db.Exec("DELETE FROM matches")
	if err != nil {
		fmt.Println("error deleting matches:", err)
		return err
	}

	// set name
	_, err = db.Exec("INSERT INTO options (key, value) VALUES ('name', ?)", name)
	if err != nil {
		fmt.Println("error setting name:", err)
		return err
	}
	_, err = db.Exec("INSERT INTO options (key, value) VALUES ('status', ?)", "status-open")
	if err != nil {
		fmt.Println("error setting status:", err)
		return err
	}
	return nil
}

func DBGetOption(db *sql.DB, key string) string {
	var value string
	err := db.QueryRow("SELECT value FROM options WHERE key = ?", key).Scan(&value)
	if err != nil {
		return "error"
	}
	return value
}

func DBRegisterParticipant(db *sql.DB, discordID, ign string) error {
	_, err := db.Exec("INSERT OR REPLACE INTO participants (discord_id, ign) VALUES (?, ?)", discordID, ign)
	if err != nil {
		return err
	}
	return nil
}

func DBGetParticipants(db *sql.DB, groupId int) []string {
	var rows *sql.Rows
	var err error
	if groupId > 0 {
		rows, err = db.Query("SELECT ign FROM participants WHERE group_id = ?", groupId)
	} else {
		rows, err = db.Query("SELECT ign FROM participants")
	}
	if err != nil {
		return nil
	}
	defer rows.Close()
	var participants []string
	for rows.Next() {
		var ign string
		err = rows.Scan(&ign)
		if err != nil {
			return nil
		}
		participants = append(participants, ign)
	}
	return participants
}

func DBStartTournament(db *sql.DB, groupsize, bestof, finals int64) error {
	_, err := db.Exec("UPDATE options SET value = ? WHERE key = 'groupsize'", groupsize)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE options SET value = ? WHERE key = 'bestof'", bestof)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE options SET value = ? WHERE key = 'finals-bestof'", finals)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE options SET value = ? WHERE key = 'status'", "status-started")
	if err != nil {
		return err
	}

	// group participants randomly
	participants := DBGetParticipants(db, 0)
	numGroups := len(participants) / int(groupsize)
	// reduce number of groups to 2^x
	if numGroups&(numGroups-1) != 0 {
		reduced := 1
		for reduced < numGroups {
			reduced *= 2
		}
		numGroups = reduced / 2
	}

	// shuffle
	for i := range participants {
		j := i + rand.Intn(len(participants)-i)
		participants[i], participants[j] = participants[j], participants[i]
	}
	// assign to groups
	group := 1
	groupSizes := make(map[int]int)
	for _, p := range participants {
		_, err = db.Exec("UPDATE participants SET group_id = ? WHERE ign = ?", group, p)
		if err != nil {
			return err
		}
		groupSizes[group]++
		group++
		if group > numGroups {
			group = 1
		}
	}
	// name groups alphabetically
	for i := 1; i <= numGroups; i++ {
		_, err = db.Exec("INSERT INTO groups (name) VALUES (?)", fmt.Sprintf("Gruppe %c", 'A'+i-1))
		if err != nil {
			return err
		}
	}
	// populate matches: create one match per pairing in each group
	for i := 1; i <= numGroups; i++ {
		participants := DBGetParticipants(db, i)
		for j := 0; j < len(participants); j++ {
			for k := j + 1; k < len(participants); k++ {
				_, err = db.Exec("INSERT INTO matches (group_id, bestof, player1, player2) VALUES (?, ?, ?, ?)", i, bestof, participants[j], participants[k])
				if err != nil {
					return err
				}
			}
		}
	}

	start := 1
	stop := numGroups
	for {
		last := stop
		// if the declared group size is 4 or more the first two winners advance. If the group size is 3 or less, only the winner advances.
		if groupsize > 3 && groupSizes[start] > 2 {
			for i := start; i <= stop; i += 2 {
				last += 2
				_, err = db.Exec("INSERT INTO groups (name) VALUES (?)", fmt.Sprintf(i18n[lang]["winner-second"], 'A'+last-2, 'A'+i-1, 'A'+i))
				if err != nil {
					return err
				}
				_, err = db.Exec("INSERT INTO matches (group_id, bestof, player1, player2) VALUES (?, ?, ?, ?)", last-1, finals, fmt.Sprintf("!G%d", i), fmt.Sprintf("!G%d.2", i+1))
				if err != nil {
					return err
				}
				_, err = db.Exec("INSERT INTO groups (name) VALUES (?)", fmt.Sprintf(i18n[lang]["winner-second"], 'A'+last-1, 'A'+i, 'A'+i-1))
				if err != nil {
					return err
				}
				_, err = db.Exec("INSERT INTO matches (group_id, bestof, player1, player2) VALUES (?, ?, ?, ?)", last, finals, fmt.Sprintf("!G%d.2", i), fmt.Sprintf("!G%d", i+1))
				if err != nil {
					return err
				}
			}
		} else {
			for i := start; i <= stop; i += 2 {
				last += 1
				_, err = db.Exec("INSERT INTO groups (name) VALUES (?)", fmt.Sprintf(i18n[lang]["winner-groups"], 'A'+last-1, 'A'+i-1, 'A'+i))
				if err != nil {
					return err
				}
				_, err = db.Exec("INSERT INTO matches (group_id, bestof, player1, player2) VALUES (?, ?, ?, ?)", last, finals, fmt.Sprintf("!G%d", i), fmt.Sprintf("!G%d", i+1))
				if err != nil {
					return err
				}
			}
		}
		start = stop + 1
		if last == stop+1 {
			// if we only added one group, that's the finals
			break
		}
		stop = last
	}
	return nil
}

// get all groups and their participants
func DBGetGroups(db *sql.DB) []Group {
	rows, err := db.Query("SELECT g.id, g.name, p.ign FROM groups g LEFT JOIN participants p ON g.id = p.group_id WHERE p.ign IS NOT NULL AND g.complete = 0 ORDER BY g.id, p.ign")
	if err != nil {
		return nil
	}
	defer rows.Close()
	var groups = make(map[int]*Group)
	for rows.Next() {
		var id int
		var name, ign string
		err = rows.Scan(&id, &name, &ign)
		if err != nil {
			return nil
		}
		if _, ok := groups[id]; !ok {
			groups[id] = &Group{Id: id, Name: name}
		}
		groups[id].Participants = append(groups[id].Participants, ign)
	}
	var result []Group
	for _, g := range groups {
		result = append(result, *g)
	}
	return result
}

func DBGetGroupAndBestOf(db *sql.DB, p1, p2 string) (Group, int) {
	var groupId, bestof int
	var groupName string
	err := db.QueryRow("SELECT g.id, g.name, m.bestof FROM matches m LEFT JOIN groups g ON m.group_id = g.id WHERE g.complete = 0 AND (m.player1 = ? AND m.player2 = ?) OR (m.player1 = ? AND m.player2 = ?)", p1, p2, p2, p1).Scan(&groupId, &groupName, &bestof)
	if err != nil {
		return Group{}, 0
	}
	return Group{Id: groupId, Name: groupName}, bestof
}

func DBCreateMatch(db *sql.DB, p1, p2 string, score1, score2 int64) error {
	_, err := db.Exec("UPDATE matches SET score1 = ?, score2 = ? WHERE (player1 = ? AND player2 = ?)", score1, score2, p1, p2, p2, p1)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE matches SET score1 = ?, score2 = ? WHERE (player1 = ? AND player2 = ?)", score2, score1, p2, p1)
	if err != nil {
		return err
	}
	return nil
}

func DBGetMatches(db *sql.DB, groupId int) []Match {
	rows, err := db.Query("SELECT player1, player2, score1, score2 FROM matches WHERE group_id = ?", groupId)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var matches []Match
	for rows.Next() {
		var m Match
		err = rows.Scan(&m.Player1, &m.Player2, &m.Score1, &m.Score2)
		if err != nil {
			return nil
		}
		matches = append(matches, m)
	}
	return matches
}

type Score struct {
	Wins   int
	Points int
}

func DBGetScores(db *sql.DB, groupId int) (map[string]Score, error) {
	scores := make(map[string]Score)

	rows, err := db.Query("SELECT player1, score1, player2, score2 FROM matches WHERE group_id = ?", groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p1, p2 string
		var s1, s2 int
		err := rows.Scan(&p1, &s1, &p2, &s2)
		if err != nil {
			return nil, err
		}
		if _, ok := scores[p1]; !ok {
			scores[p1] = Score{}
		}
		if _, ok := scores[p2]; !ok {
			scores[p2] = Score{}
		}
		if s1 > s2 {
			scores[p1] = Score{Wins: scores[p1].Wins + 1, Points: scores[p1].Points + s1}
			scores[p2] = Score{Wins: scores[p2].Wins, Points: scores[p2].Points + s2}
		} else {
			scores[p1] = Score{Wins: scores[p1].Wins, Points: scores[p1].Points + s1}
			scores[p2] = Score{Wins: scores[p2].Wins + 1, Points: scores[p2].Points + s2}
		}
	}
	return scores, nil
}

type Advance struct {
	Player string
	Group  Group
}

func DBCheckGroupComplete(db *sql.DB, groupId int) ([]Advance, error) {
	// check if the group is closed
	var complete int
	err := db.QueryRow("SELECT complete FROM groups WHERE id = ?", groupId).Scan(&complete)
	if err != nil {
		return nil, err
	}
	if complete == 1 {
		return nil, nil
	}
	// list open matches
	var openMatches int
	err = db.QueryRow("SELECT count(*) FROM matches WHERE group_id = ? AND score1 = 0 AND score2 = 0", groupId).Scan(&openMatches)
	if err != nil {
		return nil, err
	}
	if openMatches > 0 {
		return nil, nil
	}
	// read participant count
	var participantCount int
	err = db.QueryRow("SELECT count(*) FROM participants WHERE group_id = ?", groupId).Scan(&participantCount)
	if err != nil {
		return nil, err
	}

	// group is complete, identify the successor(s)
	var nextGroupA Group
	var nextGroupB Group
	rows, err := db.Query("SELECT m.group_id, g.name FROM matches m LEFT JOIN groups g ON m.group_id = g.id WHERE player1 = ? OR player2 = ?", fmt.Sprintf("!G%d", groupId), fmt.Sprintf("!G%d", groupId))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// if no rows returned
	if !rows.Next() {
		nextGroupA.Id = 0
		nextGroupA.Name = i18n[lang]["tournament-winner"]
	} else {
		err = rows.Scan(&nextGroupA.Id, &nextGroupA.Name)
		if err != nil {
			return nil, err
		}
	}
	rows.Close()

	if participantCount > 3 {
		rows, err = db.Query("SELECT m.group_id, g.name FROM matches m LEFT JOIN groups g ON m.group_id = g.id WHERE player1 = ? OR player2 = ?", fmt.Sprintf("!G%d.2", groupId), fmt.Sprintf("!G%d.2", groupId))
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		if !rows.Next() {
			return nil, fmt.Errorf(i18n[lang]["err-group-complete"] + "no second group found")
		}
		err = rows.Scan(&nextGroupB.Id, &nextGroupB.Name)
		if err != nil {
			return nil, err
		}
		rows.Close()
	}

	// mark group as complete
	_, err = db.Exec("UPDATE groups SET complete = 1 WHERE id = ?", groupId)
	if err != nil {
		return nil, err
	}
	// calculate winner
	scores, err := DBGetScores(db, groupId)
	if err != nil {
		return nil, err
	}
	// find winner
	var winner string
	var second string
	count := 0
	var maxScore int
	var maxWins int
	for _, s := range scores {
		if s.Points > maxScore {
			maxScore = s.Points
		}
		if s.Wins > maxWins {
			maxWins = s.Wins
		}
	}
	for p, s := range scores {
		if s.Wins == maxWins {
			count++
			winner = p
		}
	}
	if count != 1 {
		for p, s := range scores {
			if s.Points == maxScore {
				count++
				winner = p
			}
		}
	}
	maxWins = 0
	maxScore = 0
	tie := false
	for p, s := range scores {
		if p != winner {
			if s.Wins > maxWins {
				maxWins = s.Wins
				second = p
				tie = false
			} else if s.Wins == maxWins {
				tie = true
			}
		}
	}
	if tie {
		for p, s := range scores {
			if p != winner {
				if s.Points > maxScore {
					maxScore = s.Points
					second = p
					tie = false
				} else if s.Points == maxScore {
					tie = true
				}
			}
		}
	}

	if count == 0 {
		return nil, fmt.Errorf(i18n[lang]["err-group-complete"] + "no winner found")
	}
	if count > 1 || tie {
		return nil, fmt.Errorf(i18n[lang]["err-group-complete"] + i18n[lang]["perfect-draw"])
	}

	// advance player to next group
	_, err = db.Exec("UPDATE participants SET group_id = ? WHERE ign = ?", nextGroupA.Id, winner)
	if err != nil {
		return nil, err
	}
	if nextGroupB.Id > 0 {
		_, err = db.Exec("UPDATE participants SET group_id = ? WHERE ign = ?", nextGroupB.Id, second)
		if err != nil {
			return nil, err
		}
	}
	// update player id in match
	_, err = db.Exec("UPDATE matches SET player1 = ? WHERE player1 = ?", winner, fmt.Sprintf("!G%d", groupId))
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("UPDATE matches SET player2 = ? WHERE player2 = ?", winner, fmt.Sprintf("!G%d", groupId))
	if err != nil {
		return nil, err
	}
	if nextGroupB.Id > 0 {
		_, err = db.Exec("UPDATE matches SET player1 = ? WHERE player1 = ?", second, fmt.Sprintf("!G%d.2", groupId))
		if err != nil {
			return nil, err
		}
		_, err = db.Exec("UPDATE matches SET player2 = ? WHERE player2 = ?", second, fmt.Sprintf("!G%d.2", groupId))
		if err != nil {
			return nil, err
		}
	}
	winners := []Advance{{Player: winner, Group: nextGroupA}}
	if nextGroupB.Id > 0 {
		winners = append(winners, Advance{Player: second, Group: nextGroupB})
	}
	return winners, nil
}

func DBCloseTournament(db *sql.DB, winner string) error {
	_, err := db.Exec("UPDATE options SET value = ? WHERE key = 'status'", "status-finished")
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO options (key, value) VALUES ('winner', ?)", winner)
	if err != nil {
		return err
	}
	return nil
}

func DBGetTournamentName(db *sql.DB) string {
	return DBGetOption(db, "name")
}

func DBGetTournamentWinner(db *sql.DB) string {
	return DBGetOption(db, "winner")
}

func DBGetTournamentStatus(db *sql.DB) string {
	return DBGetOption(db, "status")
}
