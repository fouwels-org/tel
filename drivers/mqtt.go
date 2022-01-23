// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package drivers

import (
	"context"
	"encoding/json"
	"fmt"
	"tel/config"
	"time"

	pahmqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

type MQTT struct {
	device config.MQTTDevice
	tagmap []mqttMap
	opc    *opcua.Client
	mqc    pahmqtt.Client
	buffer map[string]string
}

type mqttMap struct {
	Mqtt config.MQTTTag
	Tag  config.TagListTag
}

type mqttMessage struct {
	Timestamp   time.Time `json:"timestamp"`
	Tag         string    `json:"tag"`
	StringValue string    `json:"string_value"`
	FloatValue  float64   `json:"float_value"`
}

func NewMQTT(tags []config.TagListTag, cfg config.MQTTDriver, opc string) (*MQTT, error) {

	mb := MQTT{
		device: cfg.Device,
		buffer: map[string]string{},
	}

	err := mb.tagLoad(tags, cfg.Tags)
	if err != nil {
		return nil, fmt.Errorf("failed to load tags: %w", err)
	}

	mqconfig := pahmqtt.NewClientOptions()
	mqconfig.AddBroker(cfg.Device.Target)
	mqconfig.SetClientID(cfg.Device.ClientID)
	mqconfig.SetUsername(cfg.Device.Username)
	mqconfig.SetPassword(cfg.Device.Token)
	mqconfig.SetKeepAlive(time.Duration(cfg.Device.KeepaliveMs) * time.Millisecond)

	mqc := pahmqtt.NewClient(mqconfig)

	mb.mqc = mqc
	mb.opc = opcua.NewClient(opc)
	return &mb, nil
}

func (m *MQTT) Run(ctx context.Context) error {

	err := m.opc.Connect(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect OPC: %w", err)
	}

	token := m.mqc.Connect()
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect: %w", token.Error())
	}

	ticker := time.NewTicker(10 * time.Millisecond)

	publisher := time.NewTicker(time.Duration(m.device.ScantimeMs) * time.Millisecond)

	for range ticker.C {
		select {
		case <-ctx.Done():
			return fmt.Errorf("ctx caught")
		case <-publisher.C:
			err := m.iotick()
			if err != nil {
				return fmt.Errorf("tick error: %w", err)
			}
		}
	}
	return fmt.Errorf("unexpected exit")
}
func (m *MQTT) tagLoad(tags []config.TagListTag, mtags []config.MQTTTag) error {

	for _, v := range mtags {

		tag := config.TagListTag{}

		for _, x := range tags {
			if v.Name == x.Name {
				tag = x
			}
		}

		if tag.Name == "" {
			return fmt.Errorf("mqtt tag %v was not found in global tag list", v)
		}

		record := mqttMap{
			Mqtt: v,
			Tag:  tag,
		}

		m.tagmap = append(m.tagmap, record)
	}

	return nil
}

func (m *MQTT) iotick() error {

	for _, v := range m.tagmap {

		nid, err := v.Tag.NodeID()
		if err != nil {
			return fmt.Errorf("failed to parse nodeID for: %v: %w", v, err)
		}

		req := &ua.ReadRequest{
			MaxAge:             0,
			NodesToRead:        []*ua.ReadValueID{{NodeID: &nid}},
			TimestampsToReturn: ua.TimestampsToReturnBoth,
		}

		resp, err := m.opc.Read(req)
		if err != nil {
			return fmt.Errorf("failed to read %v (%v): %w", v.Tag.Name, nid, err)
		}
		if len(resp.Results) < 1 {
			return fmt.Errorf("no results returned for %v (%v)", v.Tag.Name, nid)
		}
		if resp.Results[0].Status != ua.StatusOK {
			return fmt.Errorf("read failed for for %v (%v): %v", v.Tag.Name, nid, resp.Results[0].Status)
		}

		variant := resp.Results[0].Value
		tp := variant.Type()
		var value interface{}

		fval := 0.0

		switch tp {
		case ua.TypeIDBoolean:
			value = variant.Bool()
			b := variant.Bool()
			if !b {
				fval = float64(0)
			} else {
				fval = float64(1)
			}
		case ua.TypeIDInt16, ua.TypeIDInt32, ua.TypeIDInt64:
			value = variant.Int()
			fval = float64(variant.Int())
		case ua.TypeIDUint16, ua.TypeIDUint32, ua.TypeIDUint64:
			value = variant.Uint()
			fval = float64(variant.Uint())
		case ua.TypeIDFloat, ua.TypeIDDouble:
			value = variant.Float()
			fval = float64(variant.Float())
		default:
			return fmt.Errorf("unknown type for tag %v: %v", nid, variant)
		}

		strval := fmt.Sprintf("%v", value)
		id := nid.String()
		_, ok := m.buffer[id]
		if !ok {
			m.buffer[id] = strval
		}

		// if no change, skip
		if ok && m.buffer[id] == strval {
			continue
		}

		m.buffer[id] = strval

		p := mqttMessage{
			Timestamp:   time.Now(),
			StringValue: strval,
			FloatValue:  fval,
		}

		j, err := json.Marshal(p)
		if err != nil {
			return fmt.Errorf("failed to marshal: %v", err)
		}

		token := m.mqc.Publish(v.Mqtt.Topic+"/"+v.Tag.Name, 0x00, true, j)

		tout := token.WaitTimeout(time.Second * 1)
		if !tout {
			return fmt.Errorf("timed out")
		}
		err = token.Error()
		if err != nil {
			return fmt.Errorf("failed to publish: %w", err)
		}
	}

	return nil
}
