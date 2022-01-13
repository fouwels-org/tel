// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package modbus

import (
	"context"
	"fmt"
	"tel/config"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

type Goose struct {
	device config.GooseDevice
	tagmap []gooseMap
	opc    *opcua.Client
}

type gooseMap struct {
	Goose  config.GooseTag
	Tag    config.TagListTag
	NodeID ua.NodeID
}

func NewGoose(tags []config.TagListTag, cfg config.GooseDriver, opc string) (*Goose, error) {

	mb := Goose{
		device: cfg.Device,
	}

	err := mb.tagLoad(tags, cfg.Tags)
	if err != nil {
		return nil, fmt.Errorf("failed to load tags: %w", err)
	}

	mb.opc = opcua.NewClient(opc)
	return &mb, nil
}

func (m *Goose) Run(ctx context.Context) error {

	err := m.opc.Connect(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect OPC: %w", err)
	}

	return fmt.Errorf("unexpected exit")
}
func (m *Goose) tagLoad(tags []config.TagListTag, mtags []config.GooseTag) error {

	for _, v := range mtags {

		tag := config.TagListTag{}

		for _, x := range tags {
			if v.Name == x.Name {
				tag = x
			}
		}

		if tag.Name == "" {
			return fmt.Errorf("goose tag %v was not found in global tag list", v)
		}

		nodeid, err := ua.ParseNodeID("ns=1;s=" + v.Name)
		if err != nil {
			return fmt.Errorf("node id could not be parsed for tag: %+v: %w", v, err)
		}

		//m.tagmap = append(m.tagmap, record)
		_ = nodeid
	}

	return nil
}
func (m *Goose) opcwrite() error {

	// for _, v := range m.tagmap {
	//
	//	var variant ua.Variant
	//
	// 	switch v.Modbus.Type {
	// 	case config.ModbusCoil:
	// 		continue
	// 	case config.ModbusDiscrete:
	// 		pvariant, err := ua.NewVariant(m.buffer.discretes[v.Modbus.Index])
	// 		if err != nil {
	// 			return fmt.Errorf("failed to encode value for %+v", v.Tag.Name)
	// 		}
	// 		variant = *pvariant
	// 	case config.ModbusHolding:
	// 		continue
	// 	case config.ModbusInput:
	// 		pvariant, err := ua.NewVariant(m.buffer.input[v.Modbus.Index])
	// 		if err != nil {
	// 			return fmt.Errorf("failed to encode value for %+v", v.Tag.Name)
	// 		}
	// 		variant = *pvariant
	// 	}

	// 	req := &ua.WriteRequest{
	// 		NodesToWrite: []*ua.WriteValue{
	// 			{
	// 				NodeID:      &v.NodeID,
	// 				AttributeID: ua.AttributeIDValue,
	// 				Value: &ua.DataValue{
	// 					EncodingMask: ua.DataValueValue,
	// 					Value:        &variant,
	// 				},
	// 			},
	// 		},
	// 	}

	// 	resp, err := m.opc.Write(req)
	// 	if err != nil {
	// 		return fmt.Errorf("write failed for %v (%v): %w", v.Tag.Name, v.NodeID.String(), err)
	// 	}
	// 	if len(resp.Results) < 1 {
	// 		return fmt.Errorf("no results returned for %v (%v)", v.Tag.Name, v.NodeID.String())
	// 	}
	// 	if resp.Results[0].Error() != ua.StatusOK.Error() {
	// 		return fmt.Errorf("write failed for %v (%v): %v", v.Tag.Name, v.NodeID.String(), resp.Results[0].Error())
	// 	}
	// }

	return nil
}

func (m *Goose) ioread() error {

	for _, v := range m.tagmap {

		_ = v

		// index := v.Modbus.Index
		// switch v.Modbus.Type {
		// case config.ModbusCoil:
		// 	continue
		// case config.ModbusDiscrete:
		// 	result, err := m.conn.ReadDiscreteInputs(index, 1)
		// 	if err != nil {
		// 		return fmt.Errorf("failed to read discrete %v: %w", index, err)
		// 	}
		// 	if result[0] == 0x0000 {
		// 		m.buffer.discretes[index] = false
		// 	} else {
		// 		m.buffer.discretes[index] = true
		// 	}
		// case config.ModbusHolding:
		// 	continue
		// case config.ModbusInput:
		// 	result, err := m.conn.ReadInputRegisters(index, 1)
		// 	if err != nil {
		// 		return fmt.Errorf("failed to read input reg %v: %w", index, err)
		// 	}
		// 	m.buffer.input[index] = binary.BigEndian.Uint16(result)
		// }
	}
	return nil
}
