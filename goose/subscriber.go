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

func Start() int {
	return int(C.Start())
}
