// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"tel/config"
	"tel/modbus"
)

type Driver interface {
	Run(ctx context.Context) error
}

func main() {

	log.SetFlags(log.Lmicroseconds | log.LUTC)

	ctx, ctxx := context.WithCancel(context.Background())
	err := run(ctx)
	if err != nil {
		log.Printf("%v", err)
	} else {
		log.Printf("exit without error?: %v", err)
	}
	ctxx()
	os.Exit(1)
}

func run(ctx context.Context) error {

	cTagList := os.Getenv("CONFIG_TAGLIST")
	cConfigDriver := os.Getenv("CONFIG_DRIVER")
	cDriver := os.Getenv("DRIVER")
	cOpc := os.Getenv("OPC")

	if cTagList == "" {
		return fmt.Errorf("CONFIG_TAGLIST is not set")
	}

	if cConfigDriver == "" {
		return fmt.Errorf("CONFIG_DRIVERS is not set")
	}

	if cDriver == "" {
		return fmt.Errorf("DRIVER is not set")
	}

	if cOpc == "" {
		return fmt.Errorf("OPC is not set")
	}

	configTags, err := config.LoadTagList(cTagList)
	if err != nil {
		return fmt.Errorf("failed to load tags: %w", err)
	}

	var driver Driver

	switch cDriver {
	case "modbus":

		configModbus, err := config.LoadModbus(cConfigDriver)
		if err != nil {
			return fmt.Errorf("failed to load modbus configuration: %w", err)
		}

		log.Printf("starting modbus as: %+v", configModbus.Modbus.Device)

		d, err := modbus.NewModbus(configTags.Tags, configModbus.Modbus, cOpc)
		if err != nil {
			return fmt.Errorf("failed to create modbus driver: %w", err)
		}
		driver = d
	default:
		return fmt.Errorf("driver %v not recognised", cDriver)
	}

	return driver.Run(ctx)
}
