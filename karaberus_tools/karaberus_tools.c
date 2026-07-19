// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

#include "karaberus_tools.h"
#include <dakara_check.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int karaberus_add_diagnostic(karaberus_reports *reports,
                             struct dakara_check_diagnostic report) {
  if (report.error_level == DC_ERROR)
    reports->failed = true;

  reports->reports =
      realloc(reports->reports, sizeof(report) * (reports->n_reports + 1));

  if (reports->reports == NULL) {
    perror("failed to allocate memory for new report");
    return -1;
  }

  reports->reports[reports->n_reports] = report;
  reports->n_reports++;
  return 0;
}

karaberus_reports
karaberus_dakara_check_avio(void *obj,
                            int (*read_packet)(void *, uint8_t *, int),
                            int64_t (*seek)(void *, int64_t, int)) {
  dakara_check_results res;
  karaberus_reports reports;
  reports.n_reports = 0;
  reports.reports = NULL;
  reports.duration = 0;
  reports.failed = false;

  dakara_check_avio(KARABERUS_BUFSIZE, obj, read_packet, seek, &res);

  // print reports for now so they are at least readable somewhere
  dakara_check_print_diagnostics(res.report, "minio object");

  if (!res.report.no_duration) {
    reports.duration = res.duration;
  }

  struct dakara_check_diagnostic diagnostic;
  while ((diagnostic = dakara_check_get_diagnostic(&res.report)).report_id !=
         DC_DONE) {
    karaberus_add_diagnostic(&reports, diagnostic);
  }

  return reports;
}

karaberus_reports
karaberus_dakara_inst_check_avio(void *obj,
                                 int (*read_packet)(void *, uint8_t *, int),
                                 int64_t (*seek)(void *, int64_t, int)) {
  dakara_check_results res;
  karaberus_reports reports;
  reports.n_reports = 0;
  reports.reports = NULL;
  reports.duration = 0;
  reports.failed = false;

  dakara_check_inst_avio(KARABERUS_BUFSIZE, obj, read_packet, seek, &res);

  // print reports for now so they are at least readable somewhere
  dakara_check_print_diagnostics(res.report, "minio object");

  if (!res.report.no_duration) {
    reports.duration = res.duration;
  }

  struct dakara_check_diagnostic diagnostic;
  while ((diagnostic = dakara_check_get_diagnostic(&res.report)).report_id !=
         DC_DONE) {
    karaberus_add_diagnostic(&reports, diagnostic);
  }

  return reports;
}

karaberus_reports
karaberus_dakara_audio_check_avio(void *obj,
                                  int (*read_packet)(void *, uint8_t *, int),
                                  int64_t (*seek)(void *, int64_t, int)) {
  dakara_check_results res;
  karaberus_reports reports;
  reports.n_reports = 0;
  reports.reports = NULL;
  reports.duration = 0;
  reports.failed = false;

  dakara_check_audio_avio(KARABERUS_BUFSIZE, obj, read_packet, seek, &res);

  if (!res.report.no_duration) {
    reports.duration = res.duration;
  }

  struct dakara_check_diagnostic diagnostic;
  while ((diagnostic = dakara_check_get_diagnostic(&res.report)).report_id !=
         DC_DONE) {
    karaberus_add_diagnostic(&reports, diagnostic);
  }

  return reports;
}

void free_reports(karaberus_reports reports) { free(reports.reports); }

karaberus_sub_reports *karaberus_check_sub(char *mem, size_t bufsize) {
  dakara_check_sub_results *res = dakara_check_subtitle_memory(mem, bufsize);
  if (res == NULL) {
    return NULL;
  }

  karaberus_sub_reports *kres = malloc(sizeof(karaberus_sub_reports));

  kres->io_error = res->report.io_error;
  if (!res->report.io_error) {
    kres->lyrics = strdup(res->lyrics);
  }

  dakara_check_sub_results_free(res);

  return kres;
}

void karaberus_sub_reports_free(karaberus_sub_reports *res) {
  if (res != NULL)
    free(res->lyrics);
  free(res);
}
