// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func LoadTagList(path string) (TagList, error) {

	c := TagList{}

	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return TagList{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	y := yaml.NewDecoder(f)
	y.SetStrict(true)
	err = y.Decode(&c)
	if err != nil {
		return TagList{}, fmt.Errorf("failed to load taglist: %w", err)
	}
	return c, nil
}

func LoadModbus(path string) (Modbus, error) {

	c := Modbus{}

	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return Modbus{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	y := yaml.NewDecoder(f)
	y.SetStrict(true)

	err = y.Decode(&c)
	if err != nil {
		return Modbus{}, fmt.Errorf("failed to load modbus: %w", err)
	}

	for _, v := range c.Modbus.Tags {
		switch v.Type {
		case ModbusCoil, ModbusDiscrete, ModbusHolding, ModbusInput:
			continue
		default:
			return Modbus{}, fmt.Errorf("invalid type, expected one of [%v, %v, %v, %v] for: %+v", ModbusCoil, ModbusDiscrete, ModbusHolding, ModbusInput, v)
		}
	}
	return c, nil
}

func LoadMqtt(path string) (MQTT, error) {

	c := MQTT{}

	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return MQTT{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	y := yaml.NewDecoder(f)
	y.SetStrict(true)

	err = y.Decode(&c)
	if err != nil {
		return MQTT{}, fmt.Errorf("failed to load mqtt: %w", err)
	}

	return c, nil
}
