// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

#include "karaberus_tools.h"
#include <dakara_check.h>
#include <stddef.h>
#include <stdint.h>

struct dakara_check_results *
karaberus_dakara_check_avio(void *obj,
                            int (*read_packet)(void *, uint8_t *, int),
                            int64_t (*seek)(void *, int64_t, int)) {
  return dakara_check_avio(KARABERUS_BUFSIZE, obj, read_packet, seek);
}

void karaberus_dakara_check_results_free(struct dakara_check_results *res) {
  dakara_check_results_free(res);
}
