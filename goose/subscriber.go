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
	"unsafe"
)

type Message struct {
	Valid                  uint32
	ErrorCode              uint32
	Timestamp              uint64
	StateNumber            uint32
	SequenceNumber         uint32
	ConfigurationReference uint32
	ApplicationID          uint32
	TTL                    uint32
	Dataset                string
	GoCBReference          string
	GoId                   string
	ValuesString           string
	Values                 MMSValue
}

//Initialize the driver
func Initialize(networkInterface string, destinationMac []byte, applicationId uint16, GoCBReference string) {

	cNetworkInterface := C.CString(networkInterface)
	cDestinationMac := C.CBytes(destinationMac)
	cGoCBReference := C.CString(GoCBReference)

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

	msg := Message{
		Valid:                  uint32(cmsg.valid),
		ErrorCode:              uint32(cmsg.error_code),
		Timestamp:              uint64(cmsg.timestamp),
		StateNumber:            uint32(cmsg.state_number),
		SequenceNumber:         uint32(cmsg.sequence_number),
		ConfigurationReference: uint32(cmsg.configuration_reference),
		ApplicationID:          uint32(cmsg.application_id),
		TTL:                    uint32(cmsg.ttl),
		Dataset:                C.GoString(cmsg.dataset),
		GoCBReference:          C.GoString(cmsg.goCb_reference),
		GoId:                   C.GoString(cmsg.go_id),
		ValuesString:           C.GoString(cmsg.values_string),
		Values:                 NewMMSValue(cmsg.values),
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
