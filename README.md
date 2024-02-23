# Turnvater

Is a simple tournament bot for discord. 

It is currently fixed to run one round of qualification groups, and a tree of knock-out finals until you have a winner. 
If at least four participants are in each group, the first two winners advance, otherwise only the winner advances.

It supports i18n in de and en

## setup

Create a JSON file ".settings" with the following values:

``` application/json
{
	"token"   : "...",
	"appId"   : "...",
	"guildId" : "...",
	"state"   : "./state.sqlite3"
}
```

The bot will use token, appId and guildId to log onto discord. 
The last option is the sqlite3 save file for the tournament data.

## slash commands

### /turn-reset

(Admin permissions required)

Starts a new tournament with a given name. All other settings will be erased.

### /turn-register

Registers a particpant with a nickname

### /turn-start

(Admin permissions required)

Starts the tournament with the following settings

* Group size: Participants will be assigned to N groups, with extra players being assigned to the first groups. 
* Bestof: Scores have to add up to this number in the qualifications
* Bestof-finals: Scores have to add up to this number in the other rounds.

### /turn-status

Prints a summary of where we're at

### /turn-result

Allows to register a result. Accepts to player names, and two scores. 