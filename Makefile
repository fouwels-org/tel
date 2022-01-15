# SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
#
# SPDX-License-Identifier: MIT

COMPOSE=docker compose
BUILDFILE=compose.yml
DOCKER=docker

.PHONY: modbus mqtt goose

build: Dockerfile
	$(COMPOSE) -f $(BUILDFILE) build
push:
	$(COMPOSE) -f $(BUILDFILE) push
up:
	$(COMPOSE) -f $(BUILDFILE) up
up-d:
	$(COMPOSE) -f $(BUILDFILE) up-d
down:
	$(COMPOSE) -f $(BUILDFILE) down

SHELL := /bin/bash
modbus:
	OPC=opc.tcp://localhost:4840 \
	DRIVER=modbus \
	CONFIG_TAGLIST=config/taglist.yml \
	CONFIG_DRIVER=config/modbus.yml \
	go run .

mqtt:
	OPC=opc.tcp://localhost:4840 \
	DRIVER=mqtt \
	CONFIG_TAGLIST=config/taglist.yml \
	CONFIG_DRIVER=config/mqtt.yml \
	go run .

goose:
	OPC=opc.tcp://localhost:4840 \
	DRIVER=goose \
	CONFIG_TAGLIST=config/taglist.yml \
	CONFIG_DRIVER=config/goose.yml \
	go run .