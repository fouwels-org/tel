package modbus

import (
	"context"
	"os"
	"tel/config"
	"testing"
)

func TestModbus(t *testing.T) {

	ctx := context.Background()
	_opc := "opc.tcp://localhost:4840"

	f, err := os.Open("../config/taglist.yml")
	if err != nil {
		t.Fatalf("failed to open: %v", err)
	}
	defer f.Close()

	f2, err := os.Open("../config/driver.yml")
	if err != nil {
		t.Fatalf("failed to open: %v", err)
	}
	defer f2.Close()

	c, err := config.LoadConfig(f, f2)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	d, err := NewModbus(c.Driver.Modbus, c.TagList, _opc)
	if err != nil {
		t.Fatalf("failed to create modbus driver: %v", err)
	}

	err = d.Run(ctx)
	if err != nil {
		t.Fatalf("%v", err)
	}

}
