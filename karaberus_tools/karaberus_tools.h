// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

#ifndef KARABERUS_TOOLS_H
#define KARABERUS_TOOLS_H
#include <stddef.h>
#include <stdint.h>
#include <dakara_check.h>
#include <libavutil/error.h>

#define KARABERUS_BUFSIZE 1024 * 1024

struct dakara_check_results *
karaberus_dakara_check_avio(void *obj,
                            int (*read_packet)(void *, uint8_t *, int),
                            int64_t (*seek)(void *, int64_t, int));

void karaberus_dakara_check_results_free(struct dakara_check_results *res);

#endif
