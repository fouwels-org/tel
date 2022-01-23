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

type GooseParseError int

const (
	GooseParseErrorNone           GooseParseError = 0
	GooseParseErrorUnknownTag     GooseParseError = 1
	GooseParseErrorTagDecode      GooseParseError = 2
	GooseParseErrorSubLevel       GooseParseError = 3
	GooseParseErrorOverflow       GooseParseError = 4
	GooseParseErrorUnderflow      GooseParseError = 5
	GooseParseErrorTypeMisMatch   GooseParseError = 6
	GooseParseErrorLengthMisMatch GooseParseError = 7
)

func (m GooseParseError) String() string {
	switch m {
	case GooseParseErrorNone:
		return "GOOSE_PARSE_ERROR_NO_ERROR"
	case GooseParseErrorUnknownTag:
		return "GOOSE_PARSE_ERROR_UNKNOWN_TAG"
	case GooseParseErrorTagDecode:
		return "GOOSE_PARSE_ERROR_TAGDECODE"
	case GooseParseErrorSubLevel:
		return "GOOSE_PARSE_ERROR_SUBLEVEL"
	case GooseParseErrorOverflow:
		return "GOOSE_PARSE_ERROR_OVERFLOW"
	case GooseParseErrorUnderflow:
		return "GOOSE_PARSE_ERROR_UNDERFLOW"
	case GooseParseErrorTypeMisMatch:
		return "GOOSE_PARSE_ERROR_TYPE_MISMATCH"
	case GooseParseErrorLengthMisMatch:
		return "GOOSE_PARSE_ERROR_LENGTH_MISMATCH"
	default:
		return "unknown"
	}
}
