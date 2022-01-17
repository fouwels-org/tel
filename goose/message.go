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
	"fmt"
	"reflect"
)

// Supported
type MmsArray []MMSValue
type MmsBoolean bool
type MmsBitString uint32
type MmsInteger int32
type MmsUnsigned uint32
type MmsFloat float64

// Unsupported
type MmsStructure error
type MmsOctetString error
type MmsVisibleString error
type MmsGeneralizedTime error
type MmsBinaryTime error
type MmsBCD error
type MmsObjId error
type MmsString error
type MmsUTCTime error
type MMSDataAccessError error

type MMSType int

const (
	MMS_ARRAY             MMSType = 0
	MMS_STRUCTURE         MMSType = 1
	MMS_BOOLEAN           MMSType = 2
	MMS_BIT_STRING        MMSType = 3
	MMS_INTEGER           MMSType = 4
	MMS_UNSIGNED          MMSType = 5
	MMS_FLOAT             MMSType = 6
	MMS_OCTET_STRING      MMSType = 7
	MMS_VISIBLE_STRING    MMSType = 8
	MMS_GENERALIZED_TIME  MMSType = 9
	MMS_BINARY_TIME       MMSType = 10
	MMS_BCD               MMSType = 11
	MMS_OBJ_ID            MMSType = 12
	MMS_STRING            MMSType = 13
	MMS_UTC_TIME          MMSType = 14
	MMS_DATA_ACCESS_ERROR MMSType = 15
)

func (m MMSType) String() string {
	switch m {
	case MMS_ARRAY:
		return "MMS_ARRAY"
	case MMS_STRUCTURE:
		return "MMS_STRUCTURE"
	case MMS_BOOLEAN:
		return "MMS_BOOLEAN"
	case MMS_BIT_STRING:
		return "MMS_BIT_STRING"
	case MMS_INTEGER:
		return "MMS_INTEGER"
	case MMS_UNSIGNED:
		return "MMS_UNSIGNED"
	case MMS_FLOAT:
		return "MMS_FLOAT"
	case MMS_OCTET_STRING:
		return "MMS_OCTET_STRING"
	case MMS_VISIBLE_STRING:
		return "MMS_VISIBLE_STRING"
	case MMS_GENERALIZED_TIME:
		return "MMS_GENERALIZED_TIME"
	case MMS_BINARY_TIME:
		return "MMS_BINARY_TIME"
	case MMS_BCD:
		return "MMS_BCD"
	case MMS_OBJ_ID:
		return "MMS_OBJ_ID"
	case MMS_STRING:
		return "MMS_STRING"
	case MMS_UTC_TIME:
		return "MMS_UTC_TIME"
	case MMS_DATA_ACCESS_ERROR:
		return "MMS_DATA_ACCESS_ERROR"
	default:
		return "unknown"
	}
}

type MMSValue struct {
	value *C.MmsValue
}

func NewMMSValue(value *C.MmsValue) MMSValue {
	return MMSValue{
		value: value,
	}
}

func (m MMSValue) String() string {
	return fmt.Sprintf("%v (%v)", m.Value(), reflect.TypeOf(m.Value()))
}

func (m MMSValue) Value() interface{} {

	vtype := MMSType(int(C.MmsValue_getType(m.value)))

	switch vtype {
	case MMS_ARRAY:
		mmsArray := MmsArray{}
		len := int(C.MmsValue_getArraySize(m.value))
		for i := 0; i < len; i++ {
			v := C.MmsValue_getElement(m.value, C.int(i))
			mmsArray = append(mmsArray, NewMMSValue(v))
		}
		return mmsArray
	case MMS_BOOLEAN:
		return MmsBoolean(C.MmsValue_getBoolean(m.value))
	case MMS_BIT_STRING:
		return MmsBitString(C.MmsValue_getBitStringAsInteger(m.value))
	case MMS_INTEGER:
		return MmsInteger(C.MmsValue_toInt32(m.value))
	case MMS_UNSIGNED:
		return MmsUnsigned(C.MmsValue_toUint32(m.value))
	case MMS_FLOAT:
		return MmsFloat(C.MmsValue_toDouble(m.value))
	case MMS_DATA_ACCESS_ERROR:
		return fmt.Errorf("MMS_DATA_ACCESS_ERROR")
	default:
		return MmsStructure(fmt.Errorf("will not parse type %v", vtype))
	}
}
