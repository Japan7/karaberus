// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

#ifndef KARABERUS_TOOLS_H
#define KARABERUS_TOOLS_H
#include <dakara_check.h>
#include <libavformat/avio.h>
#include <libavutil/error.h>
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>

#define KARABERUS_BUFSIZE 1024 * 8

typedef struct dakara_check_diagnostic dakara_check_diagnostic;

typedef struct {
  int32_t n_reports;
  dakara_check_diagnostic *reports;
  int32_t duration;
  bool failed;
} karaberus_reports;

karaberus_reports karaberus_dakara_check_avio(
    void *obj, int (*read_packet)(void *, uint8_t *, int),
    int64_t (*seek)(void *, int64_t, int));

karaberus_reports karaberus_dakara_inst_check_avio(
    void *obj, int (*read_packet)(void *, uint8_t *, int),
    int64_t (*seek)(void *, int64_t, int));

void free_reports(karaberus_reports reports);

typedef struct {
  char *lyrics;
  bool io_error;
} karaberus_sub_reports;

karaberus_sub_reports *karaberus_check_sub(char *mem, size_t bufsize);

void karaberus_sub_reports_free(karaberus_sub_reports *res);

#endif
