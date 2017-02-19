/**
 * This file contains methods which can modify game state
 * and therefore require authentication
 */
package controllers

import (
	"fmt"

	"github.com/antsy/tournament/models"
	"github.com/antsy/tournament/utils"
	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/twinj/uuid"
)

func CreateGameHandler(c *routing.Context) error {
	data := &struct {
		Name string
	}{}
	if err := c.Read(&data); err != nil {
		return err
	}

	if utils.IsEmptyOrWhitespace(data.Name) {
		return routing.NewHTTPError(422, "Name is missing")
	}

	for _, game := range models.Games {
		if game.Name == data.Name {
			return routing.NewHTTPError(422, "Tournament with given name already exists")
		}
	}

	tournament := models.NewTournament(data.Name)
	models.UpdateTournament(tournament)

	return c.Write(tournament)
}

func AddPlayer(c *routing.Context) error {
	data := &struct {
		Name     string
		Initials string
		// Tournament ID to join player to
		Game string
	}{}
	if err := c.Read(&data); err != nil {
		return err
	}

	gameID, error := uuid.Parse(data.Game)
	if error != nil {
		return routing.NewHTTPError(422, fmt.Sprint("Bad tournament id (%s)", gameID))
	}

	var tournament models.Tournament
	var player models.Player
	for _, game := range models.Games {
		if uuid.Compare(game.ID, gameID) == 0 {
			tournament = models.GetTournamentByID(gameID)
			player = models.NewPlayer(data.Name, data.Initials)

			tournament = models.AddPlayer(tournament, player)
			models.UpdateTournament(tournament)
		}
	}

	player = models.GetPlayerByID(tournament, player.ID)

	return c.Write(player)
}

func RemovePlayer(c *routing.Context) error {
	//RemovePlayerFromTournament(tournament, id)
	return c.Write("Not implemented!")
}
