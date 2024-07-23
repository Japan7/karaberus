// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

#include "karaberus_tools.h"
#include <dakara_check.h>
#include <stddef.h>
#include <stdint.h>

void karaberus_dakara_check_avio(void *obj,
                                 int (*read_packet)(void *, uint8_t *, int),
                                 int64_t (*seek)(void *, int64_t, int),
                                 dakara_check_results *res,
                                 bool needs_duration) {
  dakara_check_avio(KARABERUS_BUFSIZE, obj, read_packet, seek, res);
  if (needs_duration) {
    res->report.errors.global_duration = false;
  } else {
    res->report.errors.no_duration = false;
  }
  dakara_check_print_results(res, "minio object");
}
