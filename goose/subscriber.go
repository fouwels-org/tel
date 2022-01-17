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
	"time"
	"unsafe"
)

type Message struct {
	Timestamp             time.Time
	Valid                 uint32
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

//Initialize the driver
func Initialize(networkInterface string, destinationMac []byte, applicationId uint16, ControlBlockReference string) {

	cNetworkInterface := C.CString(networkInterface)
	cDestinationMac := C.CBytes(destinationMac)
	cGoCBReference := C.CString(ControlBlockReference)

	defer C.free(unsafe.Pointer(cNetworkInterface))
	defer C.free(unsafe.Pointer(cDestinationMac))
	defer C.free(unsafe.Pointer(cGoCBReference))

	C.Initialize(cNetworkInterface, (*C.uint8_t)(cDestinationMac), C.ushort(applicationId), cGoCBReference)
}

//Start the driver
func Start() {
	C.Start()
}

//Tick the driver
func Tick() bool {
	return int(C.Tick()) == 1
}

//Get current message
func GetCurrentMessage() Message {
	cmsg := C.GetCurrentMessage()

	datetime := time.Unix(int64(uint64(cmsg.timestamp))/1000, 0)
	msg := Message{
		Valid:                 uint32(cmsg.valid),
		ErrorCode:             uint32(cmsg.error_code),
		Timestamp:             datetime,
		StateNumber:           uint32(cmsg.state_number),
		SequenceNumber:        uint32(cmsg.sequence_number),
		ConfigurationRevision: uint32(cmsg.configuration_reference),
		ApplicationID:         uint32(cmsg.application_id),
		TTL:                   uint32(cmsg.ttl),
		Dataset:               C.GoString(cmsg.dataset),
		ControlBlockReference: C.GoString(cmsg.goCb_reference),
		Id:                    C.GoString(cmsg.go_id),
		Values:                NewMMSValue(cmsg.values),
	}
	return msg
}

//Stop and Destroy the drivr
func StopAndDestroy() {
	C.StopAndDestroy()
}

//Configure_SetObserver Set the observer flag to configure the driver to listen to any and all recieved messages.
func Configure_SetObserver() {
	C.Configure_SetObserver()
}
