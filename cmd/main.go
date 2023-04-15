package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/labbs/github-exporter/bootstrap"
	"github.com/labbs/github-exporter/config"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var version = "development"

func main() {
	serverFlags := bootstrap.ServerFlags()

	app := cli.NewApp()
	app.Name = "github-exporter"
	app.Usage = "Github exporter"
	app.Version = version
	app.Compiled = time.Now()
	app.Commands = []*cli.Command{
		{
			Name:   "server",
			Usage:  "Start the server",
			Flags:  serverFlags,
			Before: altsrc.InitInputSourceWithContext(serverFlags, altsrc.NewJSONSourceFromFlagFunc("config")),
			Action: func(c *cli.Context) error {
				appBootstrap := bootstrap.App(version)

				if config.Debug {
					appBootstrap.Logger.Debug().Interface("fiber.routes", appBootstrap.Fiber.GetRoutes()).Msg("Fiber routes")
				}

				appBootstrap.Logger.Info().Msg("Starting server on port " + strconv.Itoa(config.Port))
				return appBootstrap.Fiber.Listen(":" + strconv.Itoa(config.Port))
			},
		},
		{
			Name:  "generate-config",
			Usage: "Generate a config file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "path",
					Aliases: []string{"p"},
					Usage:   "Path to save the config file",
				},
			},
			Action: func(c *cli.Context) error {
				config.GenerateTemplateConfigFile(c.String("path"))
				return nil
			},
		},
	}

	config.Version = version
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
