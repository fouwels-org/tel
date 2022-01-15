// SPDX-FileCopyrightText: 2022 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package goose

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -liec61850
#include "subscriber.h"
*/
import "C"
import "unsafe"

//Initialize the driver
func Initialize(networkInterface string, destinationMac []byte, applicationId uint16, GoCBReference string) int {

	cNetworkInterface := C.CString(networkInterface)
	cDestinationMac := C.CBytes(destinationMac)
	cGoCBReference := C.CString(GoCBReference)

	defer C.free(unsafe.Pointer(cNetworkInterface))
	defer C.free(unsafe.Pointer(cDestinationMac))
	defer C.free(unsafe.Pointer(cGoCBReference))

	return int(C.Initialize(cNetworkInterface, (*C.uint8_t)(cDestinationMac), C.ushort(applicationId), cGoCBReference))
}

//Start the driver
func Start() int {
	return int(C.Start())
}

//Configure_SetObserver Set the observer flag to configure the driver to listen to any and all recieved messages.
func Configure_SetObserver() int {
	return int(C.Configure_SetObserver())
}
