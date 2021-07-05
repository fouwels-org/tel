package modbus

import "tel/config"

type Modbus struct {
	c config.DriverModbus
}

func NewModbus(c config.DriverModbus) *Modbus {
	return &Modbus{
		c: c,
	}
}
