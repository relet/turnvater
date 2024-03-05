package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite", "testing.sqlite3")
	if err != nil {
		fmt.Println("error opening database", err)
		return nil
	}

	return db
}

func TestTournament(t *testing.T) {
	db := InitDB()
	defer db.Close()

	if db == nil {
		t.Errorf("Error opening database")
	}
	// reset the tournament with a random id
	randName := fmt.Sprintf("test-%d", rand.Intn(1000))
	DBResetTournament(db, randName)

	name := DBGetTournamentName(db)
	if name != randName {
		t.Errorf("Expected tournament name %s, got %s", randName, name)
	}

	//register 17 participants
	for i := 0; i < 21; i++ {
		DBRegisterParticipant(db, fmt.Sprintf("user%d", i), fmt.Sprintf("ign%d", i))
	}

	// check if 17 participants are registered
	participants := DBGetParticipants(db, 0)
	if len(participants) != 21 {
		t.Errorf("Expected 21 participants, got %d", len(participants))
	}

	// start the tournament with group size 4
	DBStartTournament(db, 4, 3, 5)

	// check if the tournament is started
	status := DBGetTournamentStatus(db)
	if status != "status-started" {
		t.Errorf("Expected status-started, got %s", status)
	}

	// check if the tournament has 4 active groups, this tests group count reduction to 2^x
	groups := DBGetGroups(db)
	if len(groups) != 4 {
		t.Errorf("Expected 4 groups, got %d", len(groups))
	}

	// check if group 1 is complete
	advance, _, err := DBCheckGroupComplete(db, 1)
	if err != nil {
		t.Errorf("Error checking group 1: %s", err)
	}
	if advance != nil {
		t.Errorf("Expected group 1 to be incomplete")
	}

	// generate results for all matches
	for _, group := range groups {
		for _, p1 := range group.Participants {
			for _, p2 := range group.Participants {
				if p1 != p2 {
					DBCreateMatch(db, p1, p2, 3, 2)
				}
			}
		}
	}

	// check if group 1 is complete
	advance, _, err = DBCheckGroupComplete(db, 1)
	if err != nil {
		t.Errorf("Error checking group 1: %s", err)
	}
	if advance == nil {
		t.Errorf("Expected group 1 to be complete")
	}

	// complete the remaining groups
	for _, group := range groups {
		advance, _, err = DBCheckGroupComplete(db, group.Id)
		if err != nil {
			t.Errorf("Error checking group %d: %s", group.Id, err)
		}
		if group.Id == 1 && advance != nil {
			t.Errorf("Expected group 1 to be skipped")
		} else if group.Id > 1 && len(advance) < 2 {
			t.Errorf("Expected group %d to be complete and have at least 2 winners", group.Id)
		}
	}

	// group 5 is critical, because the sizes of the incoming groups mismatch, check it has two participants. we also call checkgroup twice
	groups = DBGetGroups(db)
	for _, group := range groups {
		if group.Id == 5 {
			if len(group.Participants) != 2 {
				t.Errorf("Expected group 5 to have 2 participants, got %d", len(group.Participants))
			}
		}
	}

	// check if the tournament is finished
	status = DBGetTournamentStatus(db)
	if status != "status-started" {
		t.Errorf("Expected status-started, got %s", status)
	}

	//complete the process for the remaining groups
	for {
		groups = DBGetGroups(db)
		if len(groups) == 1 {
			break
		}
		for _, group := range groups {
			for _, p1 := range group.Participants {
				for _, p2 := range group.Participants {
					if p1 != p2 {
						DBCreateMatch(db, p1, p2, 3, 2)
					}
				}
			}
			_, _, err = DBCheckGroupComplete(db, group.Id)
			if err != nil {
				t.Errorf("Error checking group %d: %s", group.Id, err)
			}
		}
	}

	// check if the tournament is finished
	status = DBGetTournamentStatus(db)
	if status != "status-started" {
		t.Errorf("Expected status-started, got %s", status)
	}

	// close the finals
	groups = DBGetGroups(db)
	if len(groups) != 1 {
		t.Errorf("Expected 1 group, got %d", len(groups))
	}
	if len(groups[0].Participants) != 2 {
		t.Errorf("Expected 2 participants, got %d", len(groups[0].Participants))
	}

	DBCreateMatch(db, groups[0].Participants[0], groups[0].Participants[1], 3, 2)
	advance, _, _ = DBCheckGroupComplete(db, groups[0].Id)

	if len(advance) != 1 {
		t.Errorf("Expected 1 winner, got %d", len(advance))
	}
	if advance[0].Group.Id != 0 {
		t.Errorf("Expected winner to be promoted to group 0, got %d", advance[0].Group.Id)
	}

	DBCloseTournament(db, advance[0].Player)

	// check if the tournament is finished
	status = DBGetTournamentStatus(db)
	if status != "status-finished" {
		t.Errorf("Expected status-finished, got %s", status)
	}

	// read winner
	winner := DBGetTournamentWinner(db)
	if winner != advance[0].Player {
		t.Errorf("Expected winner %s, got %s", advance[0].Player, winner)
	}

	//finally, close and delete the database file
	db.Close()
	os.Remove("testing.sqlite3")
}
