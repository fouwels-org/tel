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

	gconfig, err := config.LoadGoose("../config/goose.yml")
	if err != nil {
		t.Fatalf("failed to load taglist: %v", err)
	}

	d, err := NewGoose(tags.Tags, gconfig.Goose, _opc)
	if err != nil {
		t.Fatalf("failed to create goose driver: %v", err)
	}

	err = d.Run(ctx)
	if err != nil {
		t.Fatalf("%v", err)
	}

}
