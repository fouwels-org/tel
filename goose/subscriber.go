// SPDX-FileCopyrightText: 2022 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package goose

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -liec61850
#include "headers.h"
*/
import "C"
import (
	"time"
	"unsafe"
)

type Subscriber struct {
	subscriber C.GooseSubscriber
}

func NewSubscriber(destinationMac []byte, applicationId uint16, ControlBlockReference string) *Subscriber {

	cDestinationMac := C.CBytes(destinationMac)
	cGoCBReference := C.CString(ControlBlockReference)

	defer C.free(unsafe.Pointer(cDestinationMac))
	defer C.free(unsafe.Pointer(cGoCBReference))

	s := C.GooseSubscriber_create(cGoCBReference, nil)
	C.GooseSubscriber_setDstMac(s, (*C.uint8_t)(cDestinationMac))
	C.GooseSubscriber_setAppId(s, C.ushort(applicationId))

	return &Subscriber{
		subscriber: s,
	}
}

//Get current message
func (s *Subscriber) GetCurrentMessage() Message {

	datetime := time.Unix(int64(uint64(C.GooseSubscriber_getTimestamp(s.subscriber)))/1000, 0)
	msg := Message{
		Header: Header{
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
		},
		Values: NewMMSValue(C.GooseSubscriber_getDataSetValues(s.subscriber)),
	}
	return msg
}

func (s *Subscriber) Configure_SetObserver() {
	C.GooseSubscriber_setObserver(s.subscriber)
}
