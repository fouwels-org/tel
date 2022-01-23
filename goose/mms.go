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
	"fmt"
)

type MMSType int

const (
	MMSTypeArray           MMSType = 0
	MMSTypeStructure       MMSType = 1
	MMSTypeBoolean         MMSType = 2
	MMSTypeBitstring       MMSType = 3
	MMSTypeInteger         MMSType = 4
	MMSTypeUnsigned        MMSType = 5
	MMSTypeFloat           MMSType = 6
	MMSTypeOctetString     MMSType = 7
	MMSTypeVisibleString   MMSType = 8
	MMSTypeGeneralizedTime MMSType = 9
	MMSTypeBinaryTime      MMSType = 10
	MMSTypeBCD             MMSType = 11
	MMSTypeObjectId        MMSType = 12
	MMSTypeString          MMSType = 13
	MMSTypeUTCTime         MMSType = 14
	MMSTypeDataAccessError MMSType = 15
	MMSTypeNull            MMSType = 255
)

func (m MMSType) String() string {
	switch m {
	case MMSTypeArray:
		return "MMS_ARRAY"
	case MMSTypeStructure:
		return "MMS_STRUCTURE"
	case MMSTypeBoolean:
		return "MMS_BOOLEAN"
	case MMSTypeBitstring:
		return "MMS_BIT_STRING"
	case MMSTypeInteger:
		return "MMS_INTEGER"
	case MMSTypeUnsigned:
		return "MMS_UNSIGNED"
	case MMSTypeFloat:
		return "MMS_FLOAT"
	case MMSTypeOctetString:
		return "MMS_OCTET_STRING"
	case MMSTypeVisibleString:
		return "MMS_VISIBLE_STRING"
	case MMSTypeGeneralizedTime:
		return "MMS_GENERALIZED_TIME"
	case MMSTypeBinaryTime:
		return "MMS_BINARY_TIME"
	case MMSTypeBCD:
		return "MMS_BCD"
	case MMSTypeObjectId:
		return "MMS_OBJ_ID"
	case MMSTypeString:
		return "MMS_STRING"
	case MMSTypeUTCTime:
		return "MMS_UTC_TIME"
	case MMSTypeDataAccessError:
		return "MMS_DATA_ACCESS_ERROR"
	case MMSTypeNull:
		return "NULL"
	default:
		return "UNKNOWN"
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

func (m MMSValue) IsNull() bool {
	return m.value == nil
}

func (m MMSValue) String() string {
	if m.IsNull() {
		return fmt.Sprintf("%v", "nil")
	}
	return fmt.Sprintf("%v", m.Read())
}

func (m MMSValue) Type() MMSType {
	if m.IsNull() {
		return MMSTypeNull
	}
	return MMSType(int(C.MmsValue_getType(m.value)))
}

func (m MMSValue) Read() interface{} {

	switch m.Type() {
	case MMSTypeArray:
		mmsArray := []MMSValue{}
		len := int(C.MmsValue_getArraySize(m.value))
		for i := 0; i < len; i++ {
			v := C.MmsValue_getElement(m.value, C.int(i))
			gmms := NewMMSValue(v)

			mmsArray = append(mmsArray, gmms)
		}
		return mmsArray
	case MMSTypeBoolean:
		return bool(C.MmsValue_getBoolean(m.value))
	case MMSTypeBitstring:
		return uint32(C.MmsValue_getBitStringAsInteger(m.value))
	case MMSTypeInteger:
		return int32(C.MmsValue_toInt32(m.value))
	case MMSTypeUnsigned:
		return uint32(C.MmsValue_toUint32(m.value))
	case MMSTypeFloat:
		return float64(C.MmsValue_toDouble(m.value))
	default:
		return fmt.Errorf("unsupported type: %v", m.Type())
	}
}
