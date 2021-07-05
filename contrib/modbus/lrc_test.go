// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// SPDX-FileCopyrightText: 2014 (c) Quoc-Viet Nguyen
//
// SPDX-License-Identifier: BSD-3-Clause

package modbus

import (
	"testing"
)

func TestLRC(t *testing.T) {
	var lrc lrc
	lrc.reset().pushByte(0x01).pushByte(0x03)
	lrc.pushBytes([]byte{0x01, 0x0A})

	if lrc.value() != 0xF1 {
		t.Fatalf("lrc expected %v, actual %v", 0xF1, lrc.value())
	}
}
