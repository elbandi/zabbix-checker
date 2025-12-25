package main

import (
	"encoding/json"
	"github.com/torbenconto/bambulabs_api"
	"github.com/urfave/cli/v2"
	"log"
)

var statusCommand = cli.Command{
	Name:  "status",
	Usage: "status",
	Flags: []cli.Flag{
		&hostFlag,
		&accessCodeFlag,
		&serialNumberFlag,
	},
	Action: cmdStatus,
}

func cmdStatus(ctx *cli.Context) error {
	log.SetOutput(ctx.App.ErrWriter)
	config := &bambulabs_api.PrinterConfig{
		Host:         ctx.String(hostFlag.Name),
		AccessCode:   ctx.String(accessCodeFlag.Name),
		SerialNumber: ctx.String(serialNumberFlag.Name),
	}

	// Create printer object
	printer := bambulabs_api.NewPrinter(config)

	// Connect to printer via MQTT
	err := printer.Connect()
	if err != nil {
		return err
	}
	defer printer.Disconnect()
	data, err := printer.Data()
	if err != nil {
		return err
	}
	if !data.IsEmpty() {
		d, err := json.Marshal(data)
		if err != nil {
			return err
		}
		println(string(d))
	}
	return nil
}
