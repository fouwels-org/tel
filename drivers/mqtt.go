// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package drivers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	Mqtt          config.MQTTTag
	Tag           config.TagListTag
	MonitorHandle uint32
}

type mqttMessage struct {
	Timestamp time.Time   `json:"timestamp"`
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
}

func NewMQTT(tags []config.TagListTag, cfg config.MQTTDriver, opc string) (*MQTT, error) {

	mb := MQTT{
		device: cfg.Device,
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
	mqconfig.SetKeepAlive(time.Duration(cfg.Device.KeepAliveMs) * time.Millisecond)

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
	defer m.opc.CloseSessionWithContext(ctx)

	token := m.mqc.Connect()
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect: %w", token.Error())
	}

	subChan := make(chan *opcua.PublishNotificationData)

	sub, err := m.opc.Subscribe(&opcua.SubscriptionParameters{
		Interval:                   time.Duration(m.device.SubscriptionMs) * time.Millisecond,
		LifetimeCount:              opcua.DefaultSubscriptionLifetimeCount,
		MaxKeepAliveCount:          opcua.DefaultSubscriptionMaxKeepAliveCount,
		MaxNotificationsPerPublish: opcua.DefaultSubscriptionMaxNotificationsPerPublish,
		Priority:                   opcua.DefaultSubscriptionPriority,
	}, subChan)

	defer sub.Cancel(ctx)

	if err != nil {
		return fmt.Errorf("failed to create subscription")
	}

	for i, v := range m.tagmap {
		nid, err := v.Tag.NodeID()
		if err != nil {
			return fmt.Errorf("failed to get nodeId for %v: %w", v, err)
		}

		monitorHandle := uint32(i + 1)
		m.tagmap[i].MonitorHandle = monitorHandle

		res, err := sub.Monitor(ua.TimestampsToReturnBoth, &ua.MonitoredItemCreateRequest{
			ItemToMonitor: &ua.ReadValueID{
				NodeID:       &nid,
				AttributeID:  ua.AttributeIDValue,
				DataEncoding: &ua.QualifiedName{},
			},
			MonitoringMode: ua.MonitoringModeReporting,
			RequestedParameters: &ua.MonitoringParameters{
				ClientHandle:     monitorHandle,
				DiscardOldest:    true,
				Filter:           nil,
				QueueSize:        10,
				SamplingInterval: 1.0,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to monitor %v: %w", v.Tag.Name, err)
		}
		if res.Results[0].StatusCode != ua.StatusOK {
			return fmt.Errorf("bad status code monitoring %v: %v", v.Tag.Name, res.Results[0].StatusCode)
		}
	}

	for {
		select {

		case <-ctx.Done():
			return fmt.Errorf("ctx done")

		case res := <-subChan:

			if res.Error != nil {
				log.Printf("error in sub: %v", res.Error)
				continue
			}

			switch x := res.Value.(type) {
			case *ua.DataChangeNotification:

				for _, item := range x.MonitoredItems {

					if item.ClientHandle == 0 {
						return fmt.Errorf("monitor item returned with ClientHandle 0?: %v", item)
					}

					written := false
					for _, v := range m.tagmap {
						if v.MonitorHandle == item.ClientHandle {
							val := item.Value.Value
							err := m.writeItem(v.Tag, v.Mqtt, val)
							if err != nil {
								return fmt.Errorf("failed to write %v: %w", v.Tag, err)
							}
							written = true
							break
						}
					}
					if !written {
						return fmt.Errorf("failed to identify map for handle %v: %v", item.ClientHandle, item)
					}
				}

			default:
				return fmt.Errorf("unknown change type returned: %T", x)
			}
		}
	}
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

func (m *MQTT) writeItem(tag config.TagListTag, mqtt config.MQTTTag, value *ua.Variant) error {

	nid, err := tag.NodeID()
	if err != nil {
		return fmt.Errorf("failed to parse nodeID for: %v: %w", tag, err)
	}

	req := &ua.ReadRequest{
		MaxAge:             0,
		NodesToRead:        []*ua.ReadValueID{{NodeID: &nid}},
		TimestampsToReturn: ua.TimestampsToReturnBoth,
	}

	resp, err := m.opc.Read(req)
	if err != nil {
		return fmt.Errorf("failed to read %v (%v): %w", tag.Name, nid, err)
	}
	if len(resp.Results) < 1 {
		return fmt.Errorf("no results returned for %v (%v)", tag.Name, nid)
	}
	if resp.Results[0].Status != ua.StatusOK {
		return fmt.Errorf("read failed for for %v (%v): %v", tag.Name, nid, resp.Results[0].Status)
	}

	variant := resp.Results[0].Value
	val := variant.Value()
	p := mqttMessage{
		Timestamp: time.Now(),
		Value:     val,
		Type:      fmt.Sprintf("%T", val),
	}

	j, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal: %v", err)
	}

	token := m.mqc.Publish(mqtt.Topic+"/"+mqtt.Name, 0x00, true, j)

	tout := token.WaitTimeout(time.Second * 1)
	if !tout {
		return fmt.Errorf("timed out")
	}
	err = token.Error()
	if err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}
	return nil
}
