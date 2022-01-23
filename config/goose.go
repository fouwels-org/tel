// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package config

type Goose struct {
	Meta  ConfigMeta
	Goose GooseDriver
}

type GooseDriver struct {
	Device    GooseDevice
	Endpoints []GooseEndpoint
}

type GooseDevice struct {
	Label     string `yaml:"label"`
	Interface string `yaml:"interface"`
	LogHeader bool   `yaml:"log_header"`
	LogValues bool   `yaml:"log_values"`
}

type GooseEndpoint struct {
	ControlBlockReference string         `yaml:"control_block_reference"`
	ApplicationID         uint16         `yaml:"application_id"`
	FilterMAC             string         `yaml:"filter_mac"`
	Observer              bool           `yaml:"observer"`
	Datasets              []GooseDataset `yaml:"datasets"`
}

type GooseDataset struct {
	Name string `yaml:"name"`
	Tags int    `yaml:"tags"`
}
