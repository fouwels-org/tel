// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package config

import (
	"fmt"

	"github.com/gopcua/opcua/ua"
)

type TagList struct {
	Meta ConfigMeta   `yaml:"meta"`
	Tags []TagListTag `yaml:"tags"`
}

type TagListTag struct {
	Name         string  `yaml:"name"`
	Namespace    string  `yaml:"namespace"`
	Description  string  `yaml:"description"`
	Type         string  `yaml:"type"`
	DefaultValue float64 `yaml:"default_value"`
}

func (t TagListTag) NodeID() (ua.NodeID, error) {

	id, err := ua.ParseNodeID("ns=1;s=" + t.Name)
	if err != nil {
		return ua.NodeID{}, fmt.Errorf("node id could not be parsed: %v", err)
	}

	return *id, nil
}
