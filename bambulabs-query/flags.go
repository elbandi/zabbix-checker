package main

import "github.com/urfave/cli/v2"

var (
	hostFlag = cli.StringFlag{
		Name:     "host",
		Usage:    "Bambulab host",
		Required: true,
	}

	accessCodeFlag = cli.StringFlag{
		Name:     "access-code",
		Usage:    "Bambulab Printer access code",
		Required: true,
	}

	serialNumberFlag = cli.StringFlag{
		Name:     "serial-number",
		Usage:    "Bambulab Printer serial number",
		Required: true,
	}
)
