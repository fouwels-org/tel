// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"tel/config"
	"tel/modbus"
)

type Driver interface {
	Run(ctx context.Context) error
}

func main() {
	ctx, ctxx := context.WithCancel(context.Background())
	err := run(ctx)
	if err != nil {
		log.Printf("err: %v", err)
	} else {
		log.Printf("exit without error?: %v", err)
	}
	ctxx()
	os.Exit(1)
}

func run(ctx context.Context) error {

	cTagList := os.Getenv("CONFIG_TAGLIST")
	cConfigDrivers := os.Getenv("CONFIG_DRIVERS")
	cDriver := os.Getenv("DRIVER")
	cOpc := os.Getenv("OPC")

	if cTagList == "" {
		return fmt.Errorf("CONFIG_TAGLIST is not set")
	}

	if cDriver == "" {
		return fmt.Errorf("DRIVER is not set")
	}

	if cConfigDrivers == "" {
		return fmt.Errorf("CONFIG_DRIVERS is not set")
	}

	if cOpc == "" {
		return fmt.Errorf("OPC is not set")
	}

	f, err := os.Open(filepath.Clean(cTagList))
	if err != nil {
		return fmt.Errorf("failed to open %v: %w", cTagList, err)
	}

	f2, err := os.Open(filepath.Clean(cConfigDrivers))
	if err != nil {
		return fmt.Errorf("failed to open %v: %w", cTagList, err)
	}

	c, err := config.LoadConfig(f, f2)
	e := f.Close()
	e2 := f2.Close()
	if e2 != nil && e != nil {
		return fmt.Errorf("failed to close config file: %w", err)
	}
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	var driver Driver

	switch cDriver {
	case "modbus":
		d, err := modbus.NewModbus(c.Driver.Modbus, c.TagList, cOpc)
		if err != nil {
			return fmt.Errorf("failed to create modbus driver: %w", err)
		}
		driver = d
	default:
		return fmt.Errorf("driver %v not recognised", cDriver)
	}

	return driver.Run(ctx)
}
