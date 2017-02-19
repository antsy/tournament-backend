package controllers

import (
	"fmt"

	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/antsy/tournament/utils"
)

func VersionHandler(c *routing.Context) error {
	versionInformation := map[string]string{"Version": fmt.Sprintf("%s.%s", utils.Version, utils.Buildnumber), "Buildtime": utils.Buildtime}
	return c.Write(versionInformation)
}

func EchoHandler(c *routing.Context) error {
	msg := c.Param("message")
	return c.Write(fmt.Sprintf("%v", msg))
}
