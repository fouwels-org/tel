// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package mqtt

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
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
}

type mqttMap struct {
	Mqtt   config.MQTTTag
	Tag    config.TagListTag
	NodeID ua.NodeID
}

func NewMQTT(tags []config.TagListTag, cfg config.MQTTDriver, opc string) (*MQTT, error) {

	mb := MQTT{
		device: cfg.Device,
	}

	ml := log.New(os.Stdout, "[paho] ", log.Lmicroseconds|log.LUTC|log.Lmsgprefix)

	pahmqtt.ERROR = ml
	pahmqtt.CRITICAL = ml
	pahmqtt.WARN = ml

	err := mb.tagLoad(tags, cfg.Tags)
	if err != nil {
		return nil, fmt.Errorf("failed to load tags: %w", err)
	}

	mqaddr, err := url.Parse(cfg.Device.Target)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url %v: %w", mqaddr, err)
	}

	t := tls.Config{ // use host root CAs
	}

	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: time.Duration(cfg.Device.TimeoutMs) * time.Millisecond}, "tcp", mqaddr.Host, &t)
	if err != nil {
		return nil, fmt.Errorf("failed test dial to %v: %w", mqaddr, err)
	}
	_ = conn.Close()

	mqconfig := pahmqtt.ClientOptions{
		Servers: []*url.URL{
			mqaddr,
		},
		KeepAlive:            int64(cfg.Device.KeepaliveMs * 1000),
		ConnectTimeout:       time.Duration(cfg.Device.TimeoutMs) * time.Millisecond,
		MaxReconnectInterval: 1,
		AutoReconnect:        true,
		ClientID:             cfg.Device.ClientID,
		Username:             cfg.Device.Username,
		Password:             cfg.Device.Token,
		TLSConfig:            &t,
	}

	mqc := pahmqtt.NewClient(&mqconfig)

	mb.mqc = mqc
	mb.opc = opcua.NewClient(opc)
	return &mb, nil
}

func (m *MQTT) Run(ctx context.Context) error {

	err := m.opc.Connect(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect OPC: %w", err)
	}

	mqerr := make(chan error)
	token := m.mqc.Connect()

	go func(e chan error) {
		token.Wait()
		err := token.Error()
		if err != nil {
			e <- err
		}

	}(mqerr)

	ticker := time.NewTicker(10 * time.Millisecond)

	for range ticker.C {
		select {
		case e := <-mqerr:
			return fmt.Errorf("mqtt error: %w", e)

		case <-ctx.Done():
			return fmt.Errorf("ctx caught")

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

		nodeid, err := ua.ParseNodeID("ns=1;s=" + v.Name)
		if err != nil {
			return fmt.Errorf("node id could not be parsed for tag: %+v: %w", v, err)
		}

		record := mqttMap{
			Mqtt:   v,
			Tag:    tag,
			NodeID: *nodeid,
		}

		m.tagmap = append(m.tagmap, record)
	}

	return nil
}

func (m *MQTT) opcread() error {

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

		log.Printf("read %v: %v", v, variant.String())
	}

	return nil
}
func (m *MQTT) opcwrite() error {

	// for _, v := range m.tagmap {

	// 	var variant ua.Variant

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

	// return nil

	return fmt.Errorf("not implemented")
}