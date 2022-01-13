// SPDX-FileCopyrightText: 2022 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package goose

/*
#cgo CFLAGS: -I/usr/local/include
#include "error.h"
*/
import "C"
import "fmt"

func GetError() error {
	err := C.GoString(C.GetError())
	if err != "" {
		return fmt.Errorf("%v", err)
	}
	return nil
}
