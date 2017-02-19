package models

import (
	"encoding/json"

	"log"

	uuid "github.com/twinj/uuid"
)

type Tournament struct {
	ID      uuid.Uuid `json:"id"`
	Name    string    `json:"name"`
	Players []Player  `json:"players"`
}

type Player struct {
	ID       uuid.Uuid `json:"id"`
	Name     string    `json:"name"`
	Initials string    `json:"initials"`
	Scores   []Score   `json:"scores"`
}

type Score struct {
	Opponent uuid.Uuid `json:"opponent"`
	Wins     int       `json:"wins"`
	Losses   int       `json:"losses"`
}

// Games in memory
var Games []Tournament

/**
 * Creates new tournament entity
 */
func NewTournament(name string) Tournament {
	var tournament *Tournament
	tournament = new(Tournament)
	tournament.ID = uuid.NewV4()
	tournament.Name = name
	tournament.Players = make([]Player, 0)

	log.Printf("New tournament %s (%s) was created", name, tournament.ID)

	return *tournament
}

/**
 * Creates new tournament player
 */
func NewPlayer(name string, initials string) Player {
	var player *Player
	player = new(Player)
	player.ID = uuid.NewV4()
	player.Name = name
	player.Initials = initials

	return *player
}

/**
 * Add new player to the tournament
 */
func AddPlayer(tournament Tournament, player Player) Tournament {
	player.Scores = ConstructEmptyPlayerScores(player, tournament)
	// Add player to other player's scoresheets
	for index, opponent := range tournament.Players {
		var score Score
		score.Opponent = player.ID
		score.Wins = 0
		score.Losses = 0
		tournament.Players[index].Scores = append(opponent.Scores, score)
	}
	tournament.Players = append(tournament.Players, player)

	log.Printf("New player %s (%s) was added to tournament %s", player.Name, player.ID, tournament.Name)

	return tournament
}

/**
 * Create the initial scores array for new player
 */
func ConstructEmptyPlayerScores(player Player, tournament Tournament) []Score {
	scores := make([]Score, len(tournament.Players))
	for index, opponent := range tournament.Players {
		var score Score
		score.Opponent = opponent.ID
		score.Wins = 0
		score.Losses = 0
		scores[index] = score
	}
	return scores
}

func RemovePlayerFromTournament(tournament Tournament, id uuid.Uuid) Tournament {

	var playerName string

	removePlayer := func(players []Player, index int) []Player {
		// Here's a neat trick, we replace the element to be removed with the last element of the slice
		// and return the slice minus the last element
		// we don't get to preserve the ordering but the operation is very fast
		length := len(players) - 1
		players[index] = players[length]
		return players[:length]
	}

	removeScore := func(scores []Score, index int) []Score {
		length := len(scores) - 1
		scores[index] = scores[length]
		return scores[:length]
	}

	for pindex, player := range tournament.Players {
		if uuid.Compare(player.ID, id) == 0 {
			playerName = player.Name
			tournament.Players = removePlayer(tournament.Players, pindex)
		}
	}
	for pindex, player := range tournament.Players {
		// Also remove the player from every other player's scoresheets
		for sindex, opponent := range player.Scores {
			if uuid.Compare(opponent.Opponent, id) == 0 {
				tournament.Players[pindex].Scores = removeScore(tournament.Players[pindex].Scores, sindex)
			}
		}
	}

	log.Printf("Player %s (%s) was removed from tournament %s", playerName, id, tournament.Name)

	return tournament
}

func TournamentToJSON(tournament Tournament) string {
	tournamentAsJson, _ := json.Marshal(tournament)
	return string(tournamentAsJson)
}

func PlayerExistsInTournament(tournament Tournament, id uuid.Uuid) bool {
	for _, player := range tournament.Players {
		if uuid.Compare(player.ID, id) == 0 {
			return true
		}
	}
	return false
}

func PlayerExistsInScoreSheet(scores []Score, id uuid.Uuid) bool {
	for _, score := range scores {
		if uuid.Compare(score.Opponent, id) == 0 {
			return true
		}
	}
	return false
}

/**
 * Find the right tournament from games array
 */
func GetTournamentByID(tournamentID uuid.Uuid) Tournament {
	for _, tournament := range Games {
		if uuid.Compare(tournament.ID, tournamentID) == 0 {
			// TODO: Sort players and their scores
			return tournament
		}
	}
	return *new(Tournament)
}

func GetPlayerByID(tournament Tournament, id uuid.Uuid) Player {
	for _, player := range tournament.Players {
		if uuid.Compare(player.ID, id) == 0 {
			return player
		}
	}
	return *new(Player)
}

/**
 * Update tournament state
 */
func UpdateTournament(tournament Tournament) {
	tournamentExists := false
	for index, game := range Games {
		if uuid.Compare(game.ID, tournament.ID) == 0 {
			Games[index] = tournament
			tournamentExists = true
		}
	}
	if !tournamentExists {
		Games = append(Games, tournament)
	}

	// Trigger datastore update
	err := Persist()
	if err != nil {
		log.Print("Error: Unable to save tournament data")
	}
}
