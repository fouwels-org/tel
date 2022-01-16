// SPDX-FileCopyrightText: 2022 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

#include "subscriber.h"

const uint64_t STRING_VALUE_BUFFER_SIZE=4096;
const uint64_t BER_VALUE_BUFFER_SIZE=4096;

static int running = 1;
GooseSubscriber subscriber;
GooseReceiver receiver;

struct Message currentMessage;
struct Message GetCurrentMessage() {
  return currentMessage;
}

void Initialize(char* network_interface, uint8_t* destination_mac,  uint16_t application_id, char* goCb_reference) {

  subscriber = GooseSubscriber_create(goCb_reference, NULL);
  GooseSubscriber_setDstMac(subscriber, destination_mac);
  GooseSubscriber_setAppId(subscriber, application_id);
  //GooseSubscriber_setListener(subscriber, listener, NULL);
  
  receiver = GooseReceiver_create();
  GooseReceiver_setInterfaceId(receiver, network_interface);
  GooseReceiver_addSubscriber(receiver, subscriber);

  currentMessage.value_string = (char *) malloc(STRING_VALUE_BUFFER_SIZE);
  currentMessage.value_ber = (char *) malloc(BER_VALUE_BUFFER_SIZE);
}

void Configure_SetObserver() {
  GooseSubscriber_setObserver(subscriber);
}

void Start() {
  GooseReceiver_startThreadless(receiver);
}

int Tick() {

    int result = GooseReceiver_tick(receiver);
    if (result == 0) {
      return 0;
    }

    currentMessage.valid = GooseSubscriber_isValid(subscriber);
    currentMessage.error_code = GooseSubscriber_getParseError(subscriber);
    currentMessage.timestamp = GooseSubscriber_getTimestamp(subscriber);
    currentMessage.state_number = GooseSubscriber_getStNum(subscriber);
    currentMessage.sequence_number = GooseSubscriber_getSqNum(subscriber);
    currentMessage.configuration_reference = GooseSubscriber_getConfRev(subscriber);
    currentMessage.application_id = GooseSubscriber_getAppId(subscriber);
    currentMessage.ttl = GooseSubscriber_getTimeAllowedToLive(subscriber);

    currentMessage.dataset = GooseSubscriber_getDataSet(subscriber);
    currentMessage.goCb_reference = GooseSubscriber_getGoCbRef(subscriber);
    currentMessage.go_id = GooseSubscriber_getGoId(subscriber);
    
    MmsValue* values = GooseSubscriber_getDataSetValues(subscriber);
    MmsValue_printToBuffer(values, currentMessage.value_string, 4096);

    if (values == NULL) {
      printf("nil values returned\n");
      currentMessage.valid = 0;
      currentMessage.value_ber_length = 0;
      return 1;
    }
    // Run with encode=0 to calculate max size
    uint64_t len = MmsValue_encodeMmsData(values, currentMessage.value_ber, 0, 0);
    if (len > (BER_VALUE_BUFFER_SIZE - 1)){
      printf("failed to encode MMS, size > BER_VALUE_BUFFER_SIZE\n");
      currentMessage.valid = 0;
      currentMessage.value_ber_length = 0;
      return 1;
    }

    // Run in anger
    len = MmsValue_encodeMmsData(values, currentMessage.value_ber, 0, 1);
    currentMessage.value_ber_length = len;

    return 1;
}

void StopAndDestroy() {
  GooseReceiver_stopThreadless(receiver);
  GooseReceiver_destroy(receiver);
  free(currentMessage.value_string);
}

  //MmsValue_printToBuffer(values, currentMessage.buffer, 4096);
  
  // printf("%u.%u ", (uint32_t)(currentMessage.timestamp / 1000), (uint32_t)(currentMessage.timestamp % 1000));
  // printf("stNum: %u ", currentMessage.state_number);
  // printf("sqNum: %u ", currentMessage.sequence_number);
  // printf("TTL: %u ", currentMessage.ttl);
  // printf("valid: %u ", currentMessage.valid);
  // printf("err: %u ", currentMessage.error_code);
  // printf("dataset: %s ", currentMessage.dataset);
  // printf("go_id: %s ", currentMessage.go_id);
  // printf("goCb_reference: %s ", currentMessage.goCb_reference);
  // printf("\n> message: %s ", currentMessage.buffer);
  // printf("\n");