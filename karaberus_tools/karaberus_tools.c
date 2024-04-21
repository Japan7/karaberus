// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

#include <stdlib.h>
#include <unistd.h>
#include <dakara_check.h>
#include "karaberus_tools.h"

struct fdpipe *create_pipe(void) {
  int pipefd[2];
  if (pipe(pipefd) < 0) {
    perror("failed to create pipe");
    return NULL;
  }
  struct fdpipe *fdpipe = malloc(sizeof(struct fdpipe));
  if (fdpipe == NULL)
    return NULL;

  fdpipe->fdr = pipefd[0];
  fdpipe->fdw = pipefd[1];
  return fdpipe;
}

int read_piped(void *opaque, uint8_t *buf, int n) {
  int *fd = (int*) opaque;
  return read(*fd, buf, n);
}

struct dakara_check_results *karaberus_dakara_check(int fdr) {
	return dakara_check_avio(BUFSIZE, &fdr, read_piped, NULL);
}
