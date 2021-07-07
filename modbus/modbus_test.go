// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package modbus

import (
	"context"
	"tel/config"
	"testing"
)

func TestModbus(t *testing.T) {

	ctx := context.Background()
	_opc := "opc.tcp://localhost:4840"

	tags, err := config.LoadTagList("../config/taglist.yml")
	if err != nil {
		t.Fatalf("failed to load taglist: %v", err)
	}

	mconfig, err := config.LoadModbus("../config/modbus.yml")
	if err != nil {
		t.Fatalf("failed to load taglist: %v", err)
	}

	d, err := NewModbus(tags.Tags, mconfig.Modbus, _opc)
	if err != nil {
		t.Fatalf("failed to create modbus driver: %v", err)
	}

	err = d.Run(ctx)
	if err != nil {
		t.Fatalf("%v", err)
	}

}
