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
	"tel/drivers"
)

func main() {

	log.SetFlags(0)

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

	var driver drivers.Driver

	switch cDriver {
	case "modbus":

		configModbus, err := config.LoadModbus(cConfigDriver)
		if err != nil {
			return fmt.Errorf("failed to load modbus configuration: %w", err)
		}

		log.Printf("starting modbus as: %+v", configModbus.Modbus.Device)

		d, err := drivers.NewModbus(configTags.Tags, configModbus.Modbus, cOpc)
		if err != nil {
			return fmt.Errorf("failed to create modbus driver: %w", err)
		}
		driver = d

	case "mqtt":

		configMqtt, err := config.LoadMqtt(cConfigDriver)
		if err != nil {
			return fmt.Errorf("failed to load mqtt configuration: %w", err)
		}

		log.Printf("starting mqtt as: %+v", configMqtt.Mqtt.Device.Target)

		d, err := drivers.NewMQTT(configTags.Tags, configMqtt.Mqtt, cOpc)
		if err != nil {
			return fmt.Errorf("failed to create mqtt driver: %w", err)
		}
		driver = d

	case "goose":

		configGoose, err := config.LoadGoose(cConfigDriver)
		if err != nil {
			return fmt.Errorf("failed to load goose configuration: %w", err)
		}

		log.Printf("starting goose as: %+v", configGoose.Goose.Device)

		d, err := drivers.NewGoose(configTags.Tags, configGoose.Goose, cOpc)
		if err != nil {
			return fmt.Errorf("failed to create goose driver: %w", err)
		}
		driver = d

	default:
		return fmt.Errorf("driver %v not recognised", cDriver)
	}

	return driver.Run(ctx)
}
