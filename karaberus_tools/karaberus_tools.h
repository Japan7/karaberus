// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

#ifndef KARABERUS_TOOLS_H
#define KARABERUS_TOOLS_H
#include <dakara_check.h>
#include <libavformat/avio.h>
#include <libavutil/error.h>
#include <stddef.h>
#include <stdint.h>

#define KARABERUS_BUFSIZE 1024 * 1024

void karaberus_dakara_check_avio(void *obj,
                                 int (*read_packet)(void *, uint8_t *, int),
                                 int64_t (*seek)(void *, int64_t, int),
                                 dakara_check_results *res, bool video_stream);

#endif
