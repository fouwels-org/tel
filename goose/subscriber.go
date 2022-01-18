// SPDX-FileCopyrightText: 2022 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package goose

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -liec61850
#include "subscriber.h"
#include "libiec61850/mms_value.h"

*/
import "C"
import (
	"log"
	"time"
	"unsafe"
)

//export callback_listener
func callback_listener(sub C.GooseSubscriber, v unsafe.Pointer) {
	log.Printf("tick!")
}

var Subscriber subscriber = subscriber{}

type subscriber struct {
	subscriber C.GooseSubscriber
	receiver   C.GooseReceiver
}

//Initialize the driver
func Initialize(networkInterface string, destinationMac []byte, applicationId uint16, ControlBlockReference string) {

	s := subscriber{}

	cNetworkInterface := C.CString(networkInterface)
	cDestinationMac := C.CBytes(destinationMac)
	cGoCBReference := C.CString(ControlBlockReference)

	defer C.free(unsafe.Pointer(cNetworkInterface))
	defer C.free(unsafe.Pointer(cDestinationMac))
	defer C.free(unsafe.Pointer(cGoCBReference))

	s.subscriber = C.GooseSubscriber_create(cGoCBReference, nil)
	C.GooseSubscriber_setDstMac(s.subscriber, (*C.uint8_t)(cDestinationMac))
	C.GooseSubscriber_setAppId(s.subscriber, C.ushort(applicationId))

	s.receiver = C.GooseReceiver_create()
	C.GooseReceiver_setInterfaceId(s.receiver, cNetworkInterface)
	C.GooseReceiver_addSubscriber(s.receiver, s.subscriber)

	Subscriber = s

}

func (s *subscriber) SetListener() {
	C.register_listener(s.subscriber, nil)
}

//Start the driver
func (s *subscriber) Start() {
	C.GooseReceiver_startThreadless(s.receiver)
}

func (s *subscriber) Configure_SetObserver() {
	C.GooseSubscriber_setObserver(s.subscriber)
}

//Tick the driver
func (s *subscriber) Tick() bool {
	result := bool(C.GooseReceiver_tick(s.receiver))
	return result
}

//Get current message
func (s *subscriber) GetCurrentMessage() Message {

	datetime := time.Unix(int64(uint64(C.GooseSubscriber_getTimestamp(s.subscriber)))/1000, 0)
	msg := Message{
		Valid:                 bool(C.GooseSubscriber_isValid(s.subscriber)),
		ErrorCode:             uint32(C.GooseSubscriber_getParseError(s.subscriber)),
		Timestamp:             datetime,
		StateNumber:           uint32(C.GooseSubscriber_getStNum(s.subscriber)),
		SequenceNumber:        uint32(C.GooseSubscriber_getSqNum(s.subscriber)),
		ConfigurationRevision: uint32(C.GooseSubscriber_getConfRev(s.subscriber)),
		ApplicationID:         uint32(C.GooseSubscriber_getAppId(s.subscriber)),
		TTL:                   uint32(C.GooseSubscriber_getTimeAllowedToLive(s.subscriber)),
		Dataset:               C.GoString(C.GooseSubscriber_getDataSet(s.subscriber)),
		ControlBlockReference: C.GoString(C.GooseSubscriber_getGoCbRef(s.subscriber)),
		Id:                    C.GoString(C.GooseSubscriber_getGoId(s.subscriber)),
		Values:                NewMMSValue(C.GooseSubscriber_getDataSetValues(s.subscriber)),
	}
	return msg
}

//Stop and Destroy the drivr
func (s *subscriber) StopAndDestroy() {
	C.GooseReceiver_stopThreadless(s.receiver)
	C.GooseReceiver_destroy(s.receiver)
}

type Message struct {
	Timestamp             time.Time
	Valid                 bool
	ErrorCode             uint32
	Dataset               string
	ControlBlockReference string
	Id                    string
	StateNumber           uint32
	SequenceNumber        uint32
	ApplicationID         uint32
	ConfigurationRevision uint32
	TTL                   uint32
	Values                MMSValue
}
