// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package config

type Goose struct {
	Meta  ConfigMeta
	Goose GooseDriver
}

type GooseDriver struct {
	Device GooseDevice
	Tags   []GooseTag
}

type GooseDevice struct {
	Label         string `yaml:"label"`
	Interface     string `yaml:"interface"`
	ApplicationID uint16 `yaml:"application_id"`
	GoCbReference string `yaml:"gocb_reference"`
	FilterMac     string `yaml:"filter_mac"`
	Observer      bool   `yaml:"observer"`
}

type GooseTag struct {
	Name string
}
