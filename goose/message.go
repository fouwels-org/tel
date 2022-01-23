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
import "time"

type Header struct {
	Timestamp             time.Time
	Valid                 bool
	ErrorCode             GooseParseError
	ControlBlockReference string
	Dataset               string
	Id                    string
	StateNumber           uint32
	SequenceNumber        uint32
	ApplicationID         uint32
	ConfigurationRevision uint32
	TTL                   uint32
}

type Message struct {
	Header Header
	Value  MMSValue
}
