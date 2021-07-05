package modbus

import (
	"fmt"
	"tel/config"
	"tel/contrib/modbus"

	"github.com/gopcua/opcua"
)

type Modbus struct {
	c    config.DriverModbus
	conn modbus.Client
	opc  *opcua.Client
}

func NewModbus(c config.DriverModbus, opc string) (*Modbus, error) {

	mb := Modbus{
		c: c,
	}

	var handler modbus.ClientHandler

	switch c.Mode {
	case ModeTCP:
		handler = modbus.NewTCPClientHandler(c.Target)
	default:
		return nil, fmt.Errorf("Modbus mode %v is not supported, options are [%v]", c.Mode, ModeTCP)
	}

	mb.conn = modbus.NewClient(handler)
	mb.opc = opcua.NewClient(opc)
	return &mb, nil
}

func (m *Modbus) Run() error {

	return nil
}
