// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package config

import (
	"log"
	"testing"
)

func TestLoadConfig(t *testing.T) {

	tags, err := LoadTagList("taglist.yml")
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	mods, err := LoadModbus("modbus.yml")
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	mq, err := LoadMqtt("mqtt.yml")
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	for _, v := range tags.Tags {
		log.Printf("tags: %+v", v)
	}

	for _, v := range mods.Modbus.Tags {
		log.Printf("modbus: %+v", v)
	}

	for _, v := range mq.Mqtt.Tags {
		log.Printf("modbus: %+v", v)
	}

}
