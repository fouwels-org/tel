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
	"unsafe"
)

type Receiver struct {
	receiver C.GooseReceiver
}

func NewReceiver(networkInterface string) *Receiver {

	cNetworkInterface := C.CString(networkInterface)
	defer C.free(unsafe.Pointer(cNetworkInterface))

	r := C.GooseReceiver_create()
	C.GooseReceiver_setInterfaceId(r, cNetworkInterface)

	return &Receiver{
		receiver: r,
	}
}

func (r Receiver) RegisterSubscriber(s Subscriber) {
	C.GooseReceiver_addSubscriber(r.receiver, s.subscriber)
}

//Start the driver
func (r Receiver) Start() {
	C.GooseReceiver_startThreadless(r.receiver)
}

//Tick the driver
func (r Receiver) Tick() bool {
	result := bool(C.GooseReceiver_tick(r.receiver))
	return result
}

//Stop and Destroy the drivr
func (r Receiver) StopAndDestroy() {
	C.GooseReceiver_stopThreadless(r.receiver)
	C.GooseReceiver_destroy(r.receiver)
}
