package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

// Requirements:
//
//  - Load a configuration file

func main() {
	var config string
	var name string
	var baseUrl string
	var ownerName string
	var ownerEmail string
	var cachePath string
	var timeout time.Duration
	var themePath string
	var outputPath string
	var itemsPerPage int

	app := &cli.App{
		Name:  "mercury",
		Usage: "A feed aggregator in the style of Planet",

		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:        "config",
				Usage:       "Load configuration from `FILE`",
				Value:       "./mercury.json",
				Destination: &config,
			},
			&cli.StringFlag{
				Name:        "name",
				Value:       "My Planet!",
				Usage:       "Your planet's `name`",
				Destination: &name,
			},
			&cli.StringFlag{
				Name:        "base-url",
				Value:       "http://localhost/",
				Usage:       "`URL` of the main page",
				Destination: &baseUrl,
			},
			&cli.StringFlag{
				Name:        "owner-name",
				Usage:       "Your `name`",
				Destination: &ownerName,
			},
			&cli.StringFlag{
				Name:        "owner-email",
				Usage:       "Your `email` address",
				Destination: &ownerEmail,
			},
			&cli.PathFlag{
				Name:        "cache",
				Usage:       "Path to where cached feeds are stored",
				Value:       "./cache",
				Destination: &cachePath,
			},
			&cli.DurationFlag{
				Name:        "timeout",
				Usage:       "Number of `seconds` to wait on a given feed",
				Value:       time.Second * 20,
				Destination: &timeout,
			},
			&cli.PathFlag{
				Name:        "theme",
				Usage:       "Path to your theme's templates",
				Value:       "./theme",
				Destination: &themePath,
			},
			&cli.PathFlag{
				Name:        "output",
				Usage:       "Directory to write the site files",
				Value:       "./output",
				Destination: &outputPath,
			},
			&cli.IntFlag{
				Name:        "items-per-page",
				Value:       60,
				Usage:       "How many items to put on each page",
				Destination: &itemsPerPage,
			},
		},

		Action: func(c *cli.Context) error {
			fmt.Println("Hello, friend!")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
