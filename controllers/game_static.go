/**
 * This file contains methods which only retrive current state
 */
package controllers

import (
	"github.com/antsy/tournament/models"
	"github.com/antsy/tournament/utils"
	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/twinj/uuid"
)

func GetAllGames(c *routing.Context) error {
	type GameInfo struct {
		ID      uuid.Uuid `json:"id"`
		Name    string    `json:"name"`
		Players int       `json:"players"`
	}

	tournamentData := make([]GameInfo, 0)

	for _, game := range models.Games {
		tournamentData = append(tournamentData, GameInfo{ID: game.ID, Name: game.Name, Players: len(game.Players)})
	}

	return c.Write(tournamentData)
}

func GetGame(c *routing.Context) error {
	ID := c.Param("id")

	if utils.IsEmptyOrWhitespace(ID) {
		return routing.NewHTTPError(422, "Please specify id")
	}
	gameID, error := uuid.Parse(ID)
	if error != nil {
		return routing.NewHTTPError(422, "Bad tournament id")
	}

	for _, game := range models.Games {
		if uuid.Compare(game.ID, gameID) == 0 {
			return c.Write(game)
		}
	}

	return routing.NewHTTPError(500, "Tournament not found")
}

func GetPlayer(c *routing.Context) error {
	ID := c.Param("id")

	if utils.IsEmptyOrWhitespace(ID) {
		return routing.NewHTTPError(422, "Please specify player id")
	}

	playerID, error := uuid.Parse(ID)
	if error != nil {
		return routing.NewHTTPError(422, "Bad player id")
	}

	for _, game := range models.Games {
		for _, player := range game.Players {
			if uuid.Compare(player.ID, playerID) == 0 {
				return c.Write(player)
			}

		}
	}

	return routing.NewHTTPError(500, "Player not found")
}
