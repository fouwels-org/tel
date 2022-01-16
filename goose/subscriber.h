// SPDX-FileCopyrightText: 2022 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

#include <signal.h>
#include <stdio.h>
#include <stdlib.h>

#include "libiec61850/goose_receiver.h"
#include "libiec61850/goose_subscriber.h"
#include "libiec61850/hal_thread.h"
#include "libiec61850/linked_list.h"

struct Message {
  uint32_t valid;
  uint32_t error_code;
  uint64_t timestamp;
  uint32_t state_number;
  uint32_t sequence_number;
  uint32_t configuration_reference;
  uint32_t application_id;
  uint32_t ttl;
  char* dataset;
  char* goCb_reference;
  char* go_id;
  uint8_t* value_ber;
  uint64_t value_ber_length;
  char* value_string;
};

char *GetError();
void Initialize(char* network_interface, uint8_t* destination_mac, uint16_t application_id,  char* goCb_reference);
void Start();
int Tick();
void StopAndDestroy();
void Configure_SetObserver();
struct Message GetCurrentMessage();

static void listener(GooseSubscriber subscriber, void* parameter);
static void sigint_handler(int signalId);
