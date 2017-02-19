package main

import (
	"log"
	"net/http"

	"fmt"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/go-ozzo/ozzo-routing/slash"
	//"github.com/go-ozzo/ozzo-routing/auth"
	"github.com/antsy/tournament/controllers"
	"github.com/antsy/tournament/models"
	"github.com/antsy/tournament/utils"
	//"github.com/markbates/goth"
)

func main() {

	router := routing.New()
	router.Use(
		access.Logger(log.Printf),
		slash.Remover(http.StatusMovedPermanently),
		fault.Recovery(log.Printf),
	)

	LoadRoutes(router)
	err := models.Retrive()
	if err != nil {
		log.Printf("Unable to retrieve from data store!")
		log.Printf("If this is first time you're running the software, this is normal")
		log.Printf("But if the datafile is corrupted you should remove it now and try restarting the tournament server")
		log.Printf("Datafile location: %s", utils.StoragePath)
		log.Printf(err.Error())
	}

	log.Print(fmt.Sprintf("Tournament server started (version %s.%s / %s)", utils.Version, utils.Buildnumber, utils.Buildtime))

	log.Fatal(http.ListenAndServe(":8014", router))
}

func LoadRoutes(router *routing.Router) {
	router.Use(
		// JSON JSON everywhere
		content.TypeNegotiator(content.JSON),
	)

	router.Get("/echo/<message>", controllers.EchoHandler)
	router.Get("/version", controllers.VersionHandler)

	gameApi := router.Group("/game")
	gameApi.Use(
		content.TypeNegotiator(content.JSON),
	)

	gameApi.Post("/create", controllers.CreateGameHandler)
	gameApi.Get("/list", controllers.GetAllGames)
	gameApi.Post("/player", controllers.AddPlayer)
	gameApi.Get("/player/<id>", controllers.GetPlayer)
	gameApi.Delete("/player/<id>", controllers.RemovePlayer)
	gameApi.Get("/<id>", controllers.GetGame)
}

/*
Requirements for version 1:

- add JWT based authentication
- rights check for all game_active functions
- scores input function
- scores listing
- calculate tiebreakers

Nice to have:

- move settings to a file
- API documentation
- Minimalistic UI
- websocket communication to push notifications whenever tournament state changes

*/

/*
import (
  "errors"
  "fmt"
  "net/http"
  "github.com/dgrijalva/jwt-go"
  "github.com/go-ozzo/ozzo-routing"
  "github.com/go-ozzo/ozzo-routing/auth"
)
func main() {
  signingKey := "secret-key"
  r := routing.New()

  r.Get("/login", func(c *routing.Context) error {
    id, err := authenticate(c)
    if err != nil {
      return err
    }
    token, err := auth.NewJWT(jwt.MapClaims{
      "id": id
    }, signingKey)
    if err != nil {
      return err
    }
    return c.Write(token)
  })

  r.Use(auth.JWT(signingKey))
  r.Get("/restricted", func(c *routing.Context) error {
    claims := c.Get("JWT").(*jwt.Token).Claims.(jwt.MapClaims)
    return c.Write(fmt.Sprint("Welcome, %v!", claims["id"]))
  })
}


import "golang.org/x/net/websocket"

func SendStateChangedEvent(where string) {
	websocket.Event.Send({"StateChanged": where})
}
*/
