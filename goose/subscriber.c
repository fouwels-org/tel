// SPDX-FileCopyrightText: 2022 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

#include "error.h"
#include "subscriber.h"

static int running = 1;

int Start() {

  uint8_t destination_mac[6] = {0x00, 0x00, 0x00, 0x00, 0x00, 0x00};
  uint16_t application_id = 1000;
  char target[] = "ioname/LLN0$GO$gcbAnalogValues";
  char interface[] = "eth0";

  GooseSubscriber subscriber = GooseSubscriber_create(target, NULL);
  GooseSubscriber_setDstMac(subscriber, destination_mac);
  GooseSubscriber_setAppId(subscriber, application_id);
  GooseSubscriber_setListener(subscriber, gooseListener, NULL);

  GooseReceiver receiver = GooseReceiver_create();
  GooseReceiver_setInterfaceId(receiver, interface);
  GooseReceiver_addSubscriber(receiver, subscriber);
  GooseReceiver_start(receiver);

  // Subscribe to SigInt
  signal(SIGINT, sigint_handler);

  if (GooseReceiver_isRunning(receiver)) {

    while (running) {
      Thread_sleep(100);
    }

    GooseReceiver_stop(receiver);
    GooseReceiver_destroy(receiver);

  } else {

    GooseReceiver_stop(receiver);
    GooseReceiver_destroy(receiver);

    setError("failed to create subscriber");
    return 1;
  }

  return 0;
}

static void gooseListener(GooseSubscriber subscriber, void *parameter) {

  printf("** GOOSE event **\n");
  printf("Station Number: %u Sequence Number: %u\n", GooseSubscriber_getStNum(subscriber), GooseSubscriber_getSqNum(subscriber));
  printf("TTL: %u\n", GooseSubscriber_getTimeAllowedToLive(subscriber));

  uint64_t timestamp = GooseSubscriber_getTimestamp(subscriber);

  printf("Tmestamp: %u.%u\n", (uint32_t)(timestamp / 1000), (uint32_t)(timestamp % 1000));
  printf("Valid: %s\n", GooseSubscriber_isValid(subscriber) ? "yes" : "no");

  MmsValue *values = GooseSubscriber_getDataSetValues(subscriber);

  const int _buffer_size = 4096;
  char buffer[_buffer_size];

  MmsValue_printToBuffer(values, buffer, _buffer_size);

  printf("Message: %s\n", buffer);
}

static void sigint_handler(int signalId) {
  running = 0;
}