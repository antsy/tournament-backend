package models

import (
	"math/rand"
	"testing"
	"time"

	"fmt"

	models "github.com/antsy/tournament/models"
)

func TestCreateNewTournament(t *testing.T) {
	tournamentName := "foobar"
	tournament := models.NewTournament(tournamentName)
	if tournament.Name != tournamentName {
		t.Errorf("game was not created")
	}
}

func TestCreatingScoresForNewPlayer(t *testing.T) {
	numOfPlayers := 4
	tournament := GetTestTournament(numOfPlayers)
	ValidateTournamentPlayerCount(t, numOfPlayers, tournament)
}

func TestRemovingExistingPlayerFromTournament(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	numOfPlayers := 5
	tournament := GetTestTournament(numOfPlayers)
	ValidateTournamentPlayerCount(t, numOfPlayers, tournament)

	// Remove all players one by one in a random order
	for len(tournament.Players) > 0 {
		expectedPlayerCountAfterRemoval := len(tournament.Players) - 1

		randomIndex := rand.Intn(len(tournament.Players))
		playerName := tournament.Players[randomIndex].Name
		playerID := tournament.Players[randomIndex].ID

		tournament = models.RemovePlayerFromTournament(tournament, playerID)

		if models.PlayerExistsInTournament(tournament, playerID) {
			t.Errorf("%s was not removed from the tournament (%v)", playerName, tournament)
		}

		for _, opponent := range tournament.Players {
			if models.PlayerExistsInScoreSheet(opponent.Scores, playerID) {
				t.Errorf("%s was not removed from %s's scores", playerName, opponent.Name)
			}
		}

		ValidateTournamentPlayerCount(t, expectedPlayerCountAfterRemoval, tournament)
	}
}

func GetTestTournament(numOfPlayers int) models.Tournament {
	tournament := models.NewTournament("test tournament")
	for i := 0; i < numOfPlayers; i++ {
		alphabet := rune('A' + i)
		name := fmt.Sprintf("Player %c", alphabet)
		initials := fmt.Sprintf("%c%c%c", alphabet, alphabet, alphabet)
		player := models.NewPlayer(name, initials)
		tournament = models.AddPlayer(tournament, player)
	}

	return tournament
}

func ValidateTournamentPlayerCount(t *testing.T, expectedPlayerCount int, tournament models.Tournament) {
	if len(tournament.Players) != expectedPlayerCount {
		t.Errorf("Test tournament (%s) had wrong amount of players, expected %d, got %d", tournament.Name, expectedPlayerCount, len(tournament.Players))
	}

	for _, player := range tournament.Players {
		if len(player.Scores) != (expectedPlayerCount - 1) {
			t.Errorf("%s had wrong amount of opponents, expected %d, got %d (%v)", player.Name, (expectedPlayerCount - 1), len(player.Scores), player.Scores)
		}
	}
}
