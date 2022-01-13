// SPDX-FileCopyrightText: 2022 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT
#include "error.h"

char error[4096];

char *GetError() {
  return error;
}

void setError(char *e) {
  strncpy(error, e, 4096);
}
