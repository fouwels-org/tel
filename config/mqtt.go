// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package config

type MQTT struct {
	Meta ConfigMeta
	Mqtt MQTTDriver
}

type MQTTDriver struct {
	Device MQTTDevice
	Tags   []MQTTTag
}

type MQTTDevice struct {
	Label       string
	Target      string
	ClientID    string `yaml:"client_id"`
	Username    string
	Token       string
	KeepaliveMs int `yaml:"keepalive_ms"`
}

type MQTTTag struct {
	Name  string
	Topic string
}
