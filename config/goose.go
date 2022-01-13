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
	Label string
}

type GooseTag struct {
	Name string
}
