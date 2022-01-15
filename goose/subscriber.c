// SPDX-FileCopyrightText: 2022 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

#include "error.h"
#include "subscriber.h"

static int running = 1;
GooseSubscriber subscriber;
GooseReceiver receiver;

int Initialize(char* network_interface, uint8_t* destination_mac,  uint16_t application_id, char* goCb_reference) {

  subscriber = GooseSubscriber_create(goCb_reference, NULL);
  GooseSubscriber_setDstMac(subscriber, destination_mac);
  GooseSubscriber_setAppId(subscriber, application_id);
  GooseSubscriber_setListener(subscriber, listener, NULL);
  
  receiver = GooseReceiver_create();
  GooseReceiver_setInterfaceId(receiver, network_interface);
  GooseReceiver_addSubscriber(receiver, subscriber);

  return 0;
}

int Configure_SetObserver() {
  GooseSubscriber_setObserver(subscriber);
  return 0;
}

int Start() {
  GooseReceiver_startThreadless(receiver);

  signal(SIGINT, sigint_handler);

  if (GooseReceiver_isRunning(receiver) != 1) {
    GooseReceiver_stopThreadless(receiver);
    GooseReceiver_destroy(receiver);
    setError("failed to create iec61850 subscriber");
    return 1;
  }

  while (running){
    uint8_t received = GooseReceiver_tick(receiver);
    if (received != 1) {
      Thread_sleep(1);
    }
  } 

  GooseReceiver_stopThreadless(receiver);
  GooseReceiver_destroy(receiver);
  return 0;
}

static void listener(GooseSubscriber subscriber, void *parameter) {

  uint32_t valid = GooseSubscriber_isValid(subscriber);
  uint32_t error_code = GooseSubscriber_getParseError(subscriber);
  uint64_t timestamp = GooseSubscriber_getTimestamp(subscriber);
  uint32_t state_number = GooseSubscriber_getStNum(subscriber);
  uint32_t sequence_number = GooseSubscriber_getSqNum(subscriber);
  uint32_t configuration_reference = GooseSubscriber_getConfRev(subscriber);
  uint32_t application_id = GooseSubscriber_getAppId(subscriber);

  uint8_t src_mac[6];
  GooseSubscriber_getSrcMac(subscriber, src_mac);

  uint8_t dst_mac[6];
  GooseSubscriber_getDstMac(subscriber, dst_mac);

  char* dataset = GooseSubscriber_getDataSet(subscriber);
  char* goCb_reference = GooseSubscriber_getGoCbRef(subscriber);
  char* go_id = GooseSubscriber_getGoId(subscriber);

  MmsValue *values = GooseSubscriber_getDataSetValues(subscriber);

  char buffer[4096];
  MmsValue_printToBuffer(values, buffer, 4096);
  
  printf("%u.%u ", (uint32_t)(timestamp / 1000), (uint32_t)(timestamp % 1000));
  printf("stNum: %u ", state_number);
  printf("sqNum: %u ", sequence_number);
  printf("TTL: %u ", GooseSubscriber_getTimeAllowedToLive(subscriber));
  printf("valid: %u ", valid);
  printf("err: %u ", error_code);
  printf("dataset: %s ", dataset);
  printf("go_id: %s ", go_id);
  printf("goCb_reference: %s ", goCb_reference);
  printf("\n> message: %s ", buffer);
  printf("\n");
}

static void sigint_handler(int signalId) {
  running = 0;
}