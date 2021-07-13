// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package config

type TagList struct {
	Meta ConfigMeta
	Tags []TagListTag
}

type TagListTag struct {
	Name         string
	Namespace    string
	Description  string
	Type         string
	DefaultValue float64 `yaml:"default_value"`
}
