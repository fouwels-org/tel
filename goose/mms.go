// SPDX-FileCopyrightText: 2022 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package goose

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

func DecodeBER([]byte) string {

	return ""
}
