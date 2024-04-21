// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

#ifndef KARABERUS_TOOLS_H
#define KARABERUS_TOOLS_H

#include <stdint.h>
#include <stdlib.h>

const size_t BUFSIZE = 1024*4;

struct fdpipe {
  int fdr;
  int fdw;
};

struct fdpipe *create_pipe(void);

int read_piped(void *opaque, uint8_t *buf, int n);

struct dakara_check_results *karaberus_dakara_check(int fdr);

#endif
