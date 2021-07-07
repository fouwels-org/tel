// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package config

type ModbusMode string

const (
	ModbusModeTCP ModbusMode = "tcp"
)

type ModbusRegister string

const (
	ModbusCoil     = "coil"
	ModbusDiscrete = "discrete"
	ModbusInput    = "input"
	ModbusHolding  = "holding"
)

type Modbus struct {
	Meta   TagListMeta
	Modbus ModbusDriver
}

type ModbusDriver struct {
	Device ModbusDevice
	Tags   []ModbusTag
}

type ModbusDevice struct {
	Label      string
	Mode       string
	Target     string
	ScantimeMs int   `yaml:"scantime_ms"`
	TimeoutMs  int   `yaml:"timeout_ms"`
	Slave      uint8 `yaml:"slave_id"`
}

type ModbusTag struct {
	Name  string
	Type  ModbusRegister
	Index uint16
}
