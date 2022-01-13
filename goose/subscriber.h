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

char *GetError();
int Start();

static void gooseListener(GooseSubscriber subscriber, void* parameter);
static void sigint_handler(int signalId);
