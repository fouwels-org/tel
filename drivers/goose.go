// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package drivers

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"tel/config"
	"tel/goose"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

type Goose struct {
	device    config.GooseDevice
	endpoints []config.GooseEndpoint
	tagmap    []gooseMap
	opc       *opcua.Client
}

type gooseMap struct {
	Dataset string
	Tag     config.TagListTag
	gTag    config.GooseTag
}

func NewGoose(tags []config.TagListTag, cfg config.GooseDriver, opc string) (*Goose, error) {

	g := Goose{
		device:    cfg.Device,
		endpoints: cfg.Endpoints,
	}

	err := g.tagLoad(tags, cfg.Endpoints)
	if err != nil {
		return nil, fmt.Errorf("failed to load tags: %w", err)
	}

	g.opc = opcua.NewClient(opc)
	return &g, nil
}

func (m *Goose) tagLoad(tags []config.TagListTag, endpoints []config.GooseEndpoint) error {

	for _, e := range endpoints {

		for _, d := range e.Datasets {

			for _, t := range d.Tags {
				tag := config.TagListTag{}

				found := false
				compoundName := fmt.Sprintf("%v-%v", d.Name, t.Index)

				for _, x := range tags {

					if found {
						break
					}

					if compoundName == x.Name {
						tag = x
						found = true
					}
				}

				if !found {
					log.Printf("goose tag %v was not found in global tag list, skipped", compoundName)
				}

				record := gooseMap{
					Tag:     tag,
					Dataset: d.Name,
					gTag:    t,
				}

				m.tagmap = append(m.tagmap, record)
			}
		}
	}

	return nil
}

func (m *Goose) Run(ctx context.Context) error {

	err := m.opc.Connect(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect OPC: %w", err)
	}

	req := goose.NewReceiver(m.device.Interface)

	subs := []*goose.Subscriber{}

	for _, e := range m.endpoints {

		mac := strings.ReplaceAll(e.FilterMAC, "-", "")
		hmac, err := hex.DecodeString(mac)
		if err != nil {
			return fmt.Errorf("failed to decode configured filter_mac: %w", err)
		}

		sub := goose.NewSubscriber(hmac, e.ApplicationID, e.ControlBlockReference)

		if e.Observer {
			sub.Configure_SetObserver()
		}

		subs = append(subs, sub)
	}

	for _, s := range subs {
		req.RegisterSubscriber(s)
	}

	req.Start()
	defer req.StopAndDestroy()

	for {
		ticked := req.Tick()
		if !ticked {
			time.Sleep(1 * time.Millisecond)
		} else {
			for _, s := range subs {

				msg, err := s.GetCurrentMessage()
				if err != nil {
					log.Printf("failed to get message: %v", err)
					continue
				}
				if m.device.Log {
					log.Printf("%+v", msg.Header)
				}

				err = m.write(msg, msg.Value, 0)
				if err != nil {
					return fmt.Errorf("failed to write: %w", err)
				}

			}
		}
	}
}

func (m *Goose) write(message goose.Message, value goose.MMSValue, index int) error {

	record := message.Value.Read()
	switch record.(type) {
	case []goose.MMSValue:
		{
			v := record.([]goose.MMSValue)
			for index, i := range v {
				m.write(message, i, index)
			}
		}
	case bool, uint32, int32, float32, float64:

		mapper := config.TagListTag{}
		found := false
		for _, v := range m.tagmap {
			if found {
				break
			}
			if v.Dataset == message.Header.Dataset && v.gTag.Index == index {
				mapper = v.Tag
			}
		}
		if !found {
			return fmt.Errorf("matching tag not found within: %v, skipped", message)
		}
		nid, err := mapper.NodeID()
		if err != nil {
			return fmt.Errorf("failed to parse nodeID within %v: %w", message, err)
		}

		pvariant, err := ua.NewVariant(record)
		if err != nil {
			return fmt.Errorf("failed to encode value within %+v", message)
		}
		variant := *pvariant
		req := &ua.WriteRequest{
			NodesToWrite: []*ua.WriteValue{
				{
					NodeID:      &nid,
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
			return fmt.Errorf("write failed for %v: %w", nid.String(), err)
		}
		if len(resp.Results) < 1 {
			return fmt.Errorf("no results returned for %v", nid.String())
		}
		if resp.Results[0].Error() != ua.StatusOK.Error() {
			return fmt.Errorf("write failed for %v: %v", nid.String(), resp.Results[0].Error())
		}

	case error:
		return fmt.Errorf("error type within %v, skipped: %v", message, message.Value.Read())
	default:
		return fmt.Errorf("unnkown type within %v, skipped: %v", message, message.Value.Read())
	}

	return nil
}
