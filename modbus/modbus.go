// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package modbus

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"tel/config"
	"tel/contrib/modbus"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

type Modbus struct {
	device config.ModbusDevice
	tagmap []modbusMap
	conn   modbus.Client
	opc    *opcua.Client
	buffer registerTable
}

type modbusMap struct {
	Modbus config.ModbusTag
	Tag    config.TagListTag
	NodeID ua.NodeID
}

type registerTable struct {
	coils     [65536]bool
	discretes [65536]bool
	input     [65536]uint16
	holding   [65536]uint16
}

func NewModbus(tags []config.TagListTag, cfg config.ModbusDriver, opc string) (*Modbus, error) {

	mb := Modbus{
		device: cfg.Device,
		buffer: registerTable{
			coils:     [65536]bool{},
			discretes: [65536]bool{},
			input:     [65536]uint16{},
			holding:   [65536]uint16{},
		},
	}

	err := mb.tagLoad(tags, cfg.Tags)
	if err != nil {
		return nil, fmt.Errorf("failed to load tags: %w", err)
	}

	var handler modbus.ClientHandler

	if mb.device.TimeoutMs == 0 {
		return nil, fmt.Errorf("timeout cannot be 0")
	}
	if mb.device.Slave == 0 {
		log.Printf("Slave has been provided as 0 (broadcast), this will likely fail")
	}

	switch mb.device.Mode {
	case string(config.ModbusModeTCP):
		tcphandler := modbus.NewTCPClientHandler(mb.device.Target)
		tcphandler.Timeout = time.Duration(mb.device.TimeoutMs) * time.Millisecond
		tcphandler.SlaveId = mb.device.Slave
		handler = tcphandler
	default:
		return nil, fmt.Errorf("modbus mode %v is not supported, options are [%v]", mb.device.Mode, config.ModbusModeTCP)
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

	ticker := time.NewTicker(10 * time.Millisecond)
	ioread := time.NewTicker(time.Duration(m.device.ScantimeMs) * time.Millisecond)

	for range ticker.C {
		select {
		case <-ctx.Done():
			return fmt.Errorf("ctx caught")
		case <-ioread.C:

			err = m.opcread()
			if err != nil {
				return fmt.Errorf("opc read failed: %v", err)
			}

			err = m.iowrite()
			if err != nil {
				return fmt.Errorf("io write failed: %v", err)
			}

			err := m.ioread()
			if err != nil {
				return fmt.Errorf("io read file: %v", err)
			}

			err = m.opcwrite()
			if err != nil {
				return fmt.Errorf("opc write failed: %v", err)
			}

		}
	}
	return fmt.Errorf("unexpected exit")
}
func (m *Modbus) tagLoad(tags []config.TagListTag, mtags []config.ModbusTag) error {

	for _, v := range mtags {

		tag := config.TagListTag{}

		for _, x := range tags {
			if v.Name == x.Name {
				tag = x
			}
		}

		if tag.Name == "" {
			return fmt.Errorf("modbus tag %v was not found in global tag list", v)
		}

		nodeid, err := ua.ParseNodeID("ns=1;s=" + v.Name)
		if err != nil {
			return fmt.Errorf("node id could not be parsed for tag: %+v: %w", v, err)
		}

		record := modbusMap{
			Modbus: v,
			Tag:    tag,
			NodeID: *nodeid,
		}

		m.tagmap = append(m.tagmap, record)
	}

	return nil
}

func (m *Modbus) opcread() error {

	for _, v := range m.tagmap {

		req := &ua.ReadRequest{
			MaxAge:             0,
			NodesToRead:        []*ua.ReadValueID{{NodeID: &v.NodeID}},
			TimestampsToReturn: ua.TimestampsToReturnBoth,
		}

		resp, err := m.opc.Read(req)
		if err != nil {
			return fmt.Errorf("failed to read %v (%v): %w", v.Tag.Name, v.NodeID.String(), err)
		}
		if len(resp.Results) < 1 {
			return fmt.Errorf("no results returned for %v (%v)", v.Tag.Name, v.NodeID.String())
		}
		if resp.Results[0].Status != ua.StatusOK {
			return fmt.Errorf("read failed for for %v (%v): %v", v.Tag.Name, v.NodeID.String(), resp.Results[0].Status)
		}

		variant := resp.Results[0].Value

		switch v.Modbus.Type {
		case config.ModbusCoil:
			m.buffer.coils[v.Modbus.Index] = variant.Bool()
		case config.ModbusDiscrete:
			continue
		case config.ModbusHolding:
			m.buffer.holding[v.Modbus.Index] = uint16(variant.Uint())
		case config.ModbusInput:
			continue
		}

	}

	return nil
}
func (m *Modbus) opcwrite() error {

	for _, v := range m.tagmap {

		var variant ua.Variant

		switch v.Modbus.Type {
		case config.ModbusCoil:
			continue
		case config.ModbusDiscrete:
			pvariant, err := ua.NewVariant(m.buffer.discretes[v.Modbus.Index])
			if err != nil {
				return fmt.Errorf("failed to encode value for %+v", v.Tag.Name)
			}
			variant = *pvariant
		case config.ModbusHolding:
			continue
		case config.ModbusInput:
			pvariant, err := ua.NewVariant(m.buffer.input[v.Modbus.Index])
			if err != nil {
				return fmt.Errorf("failed to encode value for %+v", v.Tag.Name)
			}
			variant = *pvariant
		}

		req := &ua.WriteRequest{
			NodesToWrite: []*ua.WriteValue{
				{
					NodeID:      &v.NodeID,
					AttributeID: ua.AttributeIDValue,
					Value: &ua.DataValue{
						EncodingMask: ua.DataValueValue,
						Value:        &variant,
					},
				},
			},
		}

		resp, err := m.opc.Write(req)
		if err != nil {
			return fmt.Errorf("write failed for %v (%v): %w", v.Tag.Name, v.NodeID.String(), err)
		}
		if len(resp.Results) < 1 {
			return fmt.Errorf("no results returned for %v (%v)", v.Tag.Name, v.NodeID.String())
		}
		if resp.Results[0].Error() != ua.StatusOK.Error() {
			return fmt.Errorf("write failed for %v (%v): %v", v.Tag.Name, v.NodeID.String(), resp.Results[0].Error())
		}
	}

	return nil
}

func (m *Modbus) ioread() error {

	for _, v := range m.tagmap {

		index := v.Modbus.Index
		switch v.Modbus.Type {
		case config.ModbusCoil:
			continue
		case config.ModbusDiscrete:
			result, err := m.conn.ReadDiscreteInputs(index, 1)
			if err != nil {
				return fmt.Errorf("failed to read discrete %v: %w", index, err)
			}
			if result[0] == 0x0000 {
				m.buffer.discretes[index] = false
			} else {
				m.buffer.discretes[index] = true
			}
		case config.ModbusHolding:
			continue
		case config.ModbusInput:
			result, err := m.conn.ReadInputRegisters(index, 1)
			if err != nil {
				return fmt.Errorf("failed to read input reg %v: %w", index, err)
			}
			m.buffer.input[index] = binary.BigEndian.Uint16(result)
		}
	}
	return nil
}

func (m *Modbus) iowrite() error {

	for _, v := range m.tagmap {

		index := v.Modbus.Index
		switch v.Modbus.Type {
		case config.ModbusCoil:
			var i uint16
			if m.buffer.coils[index] {
				i = 0xFF00
			} else {
				i = 0x0000
			}
			_, err := m.conn.WriteSingleCoil(index, i)
			if err != nil {
				return fmt.Errorf("failed to write coil %v: %w", index, err)
			}
		case config.ModbusDiscrete:
			continue
		case config.ModbusHolding:
			_, err := m.conn.WriteSingleRegister(index, m.buffer.holding[index])
			if err != nil {
				return fmt.Errorf("failed to write holding register %v: %w", index, err)
			}
		case config.ModbusInput:
			continue
		}
	}
	return nil
}
