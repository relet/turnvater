package main

var i18n = map[string]map[string]string{
	"de": {
		"turn-reset":           "Turnier zurücksetzen",
		"turn-register":        "Anmelden",
		"turn-status":          "Status des Turniers",
		"turn-start":           "Turnier starten",
		"turn-result":          "Ergebnis eintragen",
		"status-number":        "Anzahl Teilnehmer: %d",
		"info-register":        "Anmeldung mit `/turn-register`",
		"info-grouping":        "Gruppengrösse %d: %d Gruppen, %d Spieler übrig",
		"fail-rest":            "Spieler werden zusätzlich gruppiert",
		"fail-power":           "Zu viele Gruppen. Anzahl wird reduziert auf %d. Überzählige Spieler werden zusätzlich gruppiert.",
		"grouping-ok":          "Gruppeneinteilung erfolgreich",
		"status-open":          "Anmeldung offen",
		"status-started":       "Turnier gestartet",
		"status-finished":      "Turnier beendet",
		"err-not-allowed":      "Du darfst das nicht.",
		"err-reset":            "Fehler beim Anlegen des Turniers.",
		"err-no-match":         "Paarung nicht gefunden.",
		"err-score-total":      "Die Summe der Punkte ist nicht korrekt. Wir spielen Best of %d.",
		"err-set-score":        "Fehler beim Setzen des Ergebnisses.",
		"ok-set-score":         "Ergebnis wurde gespeichert.",
		"ok-reset":             "Neues Turnier '%s' wurde initialisiert.",
		"err-start":            "Fehler beim Starten des Turniers.",
		"ok-start":             "Turnier wurde mit Gruppengrösse %d gestartet.",
		"err-register":         "Fehler bei der Anmeldung:",
		"err-register-started": "Fehler bei der Anmeldung: Das Turnier hat bereits begonnen.",
		"welcome":              "Willkommen beim Turnier, %s. Du bist jetzt mit dem Nick %s angemeldet!",
		"desc-reset":           "Startet ein neues Turnier",
		"desc-register":        "Meldet dich für das Turnier an",
		"desc-status":          "Zeigt den Status des Turniers an",
		"opt-name":             "Name des Turniers",
		"opt-ign":              "In-Game Name (Nick)",
		"opt-groupsize":        "Gruppengrösse",
		"opt-bestof":           "Best of X",
		"opt-finals-bestof":    "Best of X (Finale)",
		"opt-p1":               "Spieler 1",
		"opt-p2":               "Spieler 2",
		"opt-score1":           "Punkte Spieler 1",
		"opt-score2":           "Punkte Spieler 2",
		"winner-groups":        "Gruppe %c: Sieger %c-%c",
		"winner-second":        "Gruppe %c: Sieger %c - Zweiter %c",
		"ok-group-winner":      "%s hat sich in Gruppe %c durchgesetzt und steigt auf zu %c",
		"ok-group-second":      "%s ist in Gruppe %c Zweiter und steigt auf zu %c",
		"perfect-draw":         "Unentschieden. Um das Unentschieden aufzulösen, spielt bitte weitere Spiele und ernennt einen Sieger, indem ihr ein(!) Ergebnis aktualisiert.",
		"tournament-winner":    "Turniersieger",
		"err-group-complete":   "Fehler beim Überprüfen der Gruppe: ",
		"congratulate":         "Herzlichen Glückwunsch, %s! Du bist Turniersieger!",
	},
	"en": {
		"turn-reset":           "Reset tournament",
		"turn-register":        "Register",
		"turn-status":          "Status of the tournament",
		"turn-start":           "Start tournament",
		"turn-result":          "Enter result",
		"status-number":        "Number of participants: %d",
		"info-register":        "Register with `/turn-register`",
		"info-grouping":        "Group size %d: %d groups, %d players left",
		"fail-rest":            "Players are grouped additionally",
		"fail-power":           "Too many groups, will be reduced to %d. Surplus players are grouped additionally.",
		"grouping-ok":          "Grouping successful",
		"status-open":          "Registration open",
		"status-started":       "Tournament started",
		"status-finished":      "Tournament finished",
		"err-not-allowed":      "You are not allowed to do that.",
		"err-reset":            "Error creating the tournament.",
		"err-no-match":         "Match not found.",
		"err-score-total":      "The sum of the scores is not correct. We play Best of %d.",
		"err-set-score":        "Error setting the score.",
		"ok-set-score":         "Score has been saved.",
		"ok-reset":             "New tournament '%s' has been initialized.",
		"err-start":            "Error starting the tournament.",
		"ok-start":             "Tournament started with group size %d.",
		"err-register":         "Error registering:",
		"err-register-started": "Error registering: The tournament has already started.",
		"welcome":              "Welcome to the tournament, %s. You are now registered with the nick %s!",
		"desc-reset":           "Starts a new tournament",
		"desc-register":        "Registers you for the tournament",
		"desc-status":          "Shows the status of the tournament",
		"opt-name":             "Name of the tournament",
		"opt-ign":              "In-Game Name (Nick)",
		"opt-groupsize":        "Group size",
		"opt-bestof":           "Best of X",
		"opt-finals-bestof":    "Best of X (Finals)",
		"opt-p1":               "Player 1",
		"opt-p2":               "Player 2",
		"opt-score1":           "Score Player 1",
		"opt-score2":           "Score Player 2",
		"winner-groups":        "Group %c: Winner %c-%c",
		"winner-second":        "Group %c: Winner %c - Second %c",
		"ok-group-winner":      "%s has won group %c and advances to %c",
		"ok-group-second":      "%s is second in group %c and advances to %c",
		"perfect-draw":         "Draw. To resolve the draw, please play more games and declare a winner by updating one(!) single result.",
		"tournament-winner":    "Tournament winner",
		"err-group-complete":   "Error checking group: ",
		"congratulate":         "Congratulations, %s! You are the tournament winner!",
	},
}
