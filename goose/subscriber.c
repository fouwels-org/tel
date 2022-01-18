// SPDX-FileCopyrightText: 2022 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

#include "libiec61850/goose_subscriber.h"
#include "subscriber.h"

//void callback_listener(GooseSubscriber sub, void* p);
//void register_listener(GooseSubscriber sub, void* s);

void register_listener(GooseSubscriber sub, void* p) {
	GooseSubscriber_setListener(sub, callback_listener, p);
}