package main

import (
	"github.com/urfave/cli/v2"
)

func testCommand() *cli.Command {

	return &cli.Command{
		Name:  "test",
		Usage: "Initializes the http server that handle http requests to this application",
		Action: func(c *cli.Context) error {
			return nil
		},
	}

}
