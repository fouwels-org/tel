package modbus

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
	"strings"
	"tel/config"
	"tel/contrib/modbus"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

type Modbus struct {
	c            config.DriverModbus
	t            config.TagList
	conn         modbus.Client
	opc          *opcua.Client
	buffer       registerTable
	doubleBuffer registerTable
}

type registerTable struct {
	coils     [65536]byte
	discretes [65536]byte
	input     [65536]uint16
	holding   [65536]uint16
}

func NewModbus(c config.DriverModbus, t config.TagList, opc string) (*Modbus, error) {

	mb := Modbus{
		c: c,
		t: t,
		buffer: registerTable{
			coils:     [65536]byte{},
			discretes: [65536]byte{},
			input:     [65536]uint16{},
			holding:   [65536]uint16{},
		},
		doubleBuffer: registerTable{
			coils:     [65536]byte{},
			discretes: [65536]byte{},
			input:     [65536]uint16{},
			holding:   [65536]uint16{},
		},
	}

	var handler modbus.ClientHandler

	if c.TimeoutMs == 0 {
		return nil, fmt.Errorf("timeout cannot be 0")
	}
	if c.Slave == 0 {
		log.Printf("Slave has been provided as 0 (broadcast), this will likely fail")
	}

	switch c.Mode {
	case ModeTCP:
		tcphandler := modbus.NewTCPClientHandler(c.Target)
		tcphandler.Timeout = time.Duration(c.TimeoutMs) * time.Millisecond
		tcphandler.SlaveId = c.Slave
		handler = tcphandler
	default:
		return nil, fmt.Errorf("modbus mode %v is not supported, options are [%v]", c.Mode, ModeTCP)
	}

	mb.conn = modbus.NewClient(handler)
	mb.opc = opcua.NewClient(opc)
	return &mb, nil
}

func (m *Modbus) Run(ctx context.Context) error {

	err := m.opc.Connect(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect OPC: %w", err)
	}

	chopc := make(chan *opcua.PublishNotificationData)

	sub, err := m.opc.Subscribe(&opcua.SubscriptionParameters{
		Interval:                   100 * time.Millisecond,
		Priority:                   0,
		LifetimeCount:              10000,
		MaxKeepAliveCount:          3000,
		MaxNotificationsPerPublish: 10000,
	}, chopc)

	for i, v := range m.t.Tags {

		handleID := uint32(i)

		log.Printf("monitoring %+v", v)

		nodeID, err := ua.ParseNodeID("ns=1;s=" + v.Name)
		if err != nil {
			return fmt.Errorf("failed to parse node ID: %w", err)
		}

		miCreateRequest := opcua.NewMonitoredItemCreateRequestWithDefaults(nodeID, ua.AttributeIDValue, handleID)

		res, err := sub.Monitor(ua.TimestampsToReturnBoth, miCreateRequest)
		if err != nil || res.Results[0].StatusCode != ua.StatusOK {
			return fmt.Errorf("failed to monitor %v: %v %v", nodeID, res.Results[0].StatusCode, err)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to create opc subscription: %w", err)
	}

	go func() {
		sub.Run(ctx)
		panic("opc subscription returned unexpectedly") //todo, this is not optimal...
	}()

	err = m.eventLoop(ctx, chopc)
	return fmt.Errorf("event loop exit: %w", err)
}

func (m *Modbus) eventLoop(ctx context.Context, opchan chan *opcua.PublishNotificationData) error {

	ticker := time.NewTicker(10 * time.Millisecond)
	ioread := time.NewTicker(time.Duration(m.c.ScantimeMs) * time.Millisecond)

	for range ticker.C {
		select {
		case <-ctx.Done():
			return fmt.Errorf("eventLoop: context cancellation caught")
		case <-ioread.C:
			err := m.ioread()
			if err != nil {
				return fmt.Errorf("failed to read io: %w", err)
			}
		case o := <-opchan:
			log.Printf("opc sub: %+v", o)
		}
	}

	return nil
}

func (m *Modbus) ioread() error {
	for _, v := range m.t.Tags {
		id := strings.Split(v.Index, ".")
		if len(id) != 2 {
			return fmt.Errorf("index for %v could not be parsed: len[:] == %v", v, len(id))
		}

		objectType := id[0]
		pindex, err := strconv.ParseUint(id[1], 10, 16)
		if err != nil {
			return fmt.Errorf("index for %v could not be parsed: %w", v, err)
		}

		if pindex > 65535 {
			return fmt.Errorf("index exceeds 65535, will not process: %v", pindex)
		}
		index := uint16(pindex)

		switch strings.ToLower(objectType) {
		case "c":
			result, err := m.conn.ReadCoils(index, 1)
			if err != nil {
				return fmt.Errorf("failed to read coil: %w", err)
			}
			m.buffer.coils[index] = result[0]
		case "d":
			result, err := m.conn.ReadDiscreteInputs(index, 1)
			if err != nil {
				return fmt.Errorf("failed to read discrete: %w", err)
			}
			m.buffer.discretes[index] = result[0]
		case "i":
			result, err := m.conn.ReadInputRegisters(index, 1)
			if err != nil {
				return fmt.Errorf("failed to read input reg: %w", err)
			}
			m.buffer.input[index] = binary.BigEndian.Uint16(result)
		case "h":
			result, err := m.conn.ReadHoldingRegisters(index, 1)
			if err != nil {
				return fmt.Errorf("failed to read discrete: %w", err)
			}
			m.buffer.holding[index] = binary.BigEndian.Uint16(result)
		}

		// compare buffers
		for k := range m.buffer.coils {
			if m.doubleBuffer.coils[k] != m.buffer.coils[k] {
				log.Printf("CHANGE: coil %v = %v", k, m.buffer.coils[k])
			}
		}

		for k := range m.buffer.discretes {
			if m.doubleBuffer.discretes[k] != m.buffer.discretes[k] {
				log.Printf("CHANGE: discrete %v = %v", k, m.buffer.discretes[k])
			}
		}

		for k := range m.buffer.input {
			if m.doubleBuffer.input[k] != m.buffer.input[k] {
				log.Printf("CHANGE: input %v = %v", k, m.buffer.input[k])
			}
		}

		for k := range m.buffer.holding {
			if m.doubleBuffer.holding[k] != m.buffer.holding[k] {
				log.Printf("CHANGE: holding %v = %v", k, m.buffer.holding[k])
			}
		}

		//swap buffers
		m.doubleBuffer = m.buffer
	}

	return nil
}
