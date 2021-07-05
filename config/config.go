// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadConfig(taglist *os.File, driver *os.File) (Config, error) {

	c := Config{}

	y := yaml.NewDecoder(taglist)
	err := y.Decode(&c.TagList)
	if err != nil {
		return Config{}, fmt.Errorf("failed to load taglist: %w", err)
	}

	y = yaml.NewDecoder(driver)
	err = y.Decode(&c.Driver)
	if err != nil {
		return Config{}, fmt.Errorf("failed to load driver: %w", err)
	}

	return c, nil
}
