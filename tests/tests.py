#!/usr/bin/env python3
"""
Tests for file uploads/downloads.
We can't test it in the go test module because the mocked server panics on downloads.
Otherwise no reason to use this script rather than the go tests to test the API.
"""

from __future__ import annotations

import http.client
import json
import os
import pathlib
import secrets
import shlex
import subprocess
import time
import unittest
import urllib.request as request
import zlib
from typing import IO, ClassVar, TypedDict, final
from urllib.error import HTTPError, URLError


class KaraberusKara(TypedDict):
    title: str
    title_aliases: list[str]
    authors: list[int]
    artists: list[int]
    source_media: int
    song_order: int
    medias: list[int]
    audio_tags: list[str]
    video_tags: list[str]
    comment: str
    version: str
    language: str
    # should be NotRequired but CI runner uses python 3.10
    karaoke_creation_time: int | None
    is_hardsub: bool | None


KaraberusInputTypes = KaraberusKara


def exe_wrapper():
    exe_wrapper_str = os.environ.get("MESON_EXE_WRAPPER", "")
    return shlex.split(exe_wrapper_str)


def calculate_crc32(file: pathlib.Path) -> int:
    with file.open("rb") as fd:
        sum = 0
        while buf := fd.read(1024 * 8):
            sum = zlib.crc32(buf, sum)

        return sum


def json_body(data: KaraberusInputTypes) -> bytes:
    return json.dumps(data, separators=(",", ":")).encode()


@final
class KaraberusInstance:
    def __init__(self) -> None:
        self.port = 10201
        self.s3_port = 10202
        self.oidc_port = 10203
        self.proc: None | subprocess.Popen[bytes] = None
        self.oidc_server_proc: None | subprocess.Popen[bytes] = None
        self.gofakes3_proc: None | subprocess.Popen[bytes] = None
        self.token: None | str = None
        self.base_url = f"http://127.0.0.1:{self.port}"

    def launch_karaberus(self) -> None:
        if gofakes3_exe := os.environ.get("GOFAKES3_EXE"):
            self.gofakes3_proc = subprocess.Popen(
                [
                    gofakes3_exe,
                    "-backend",
                    "mem",
                    "-initialbucket",
                    "karaberus",
                    "-host",
                    f":{self.s3_port}",
                ]
            )

        if oidc_server_exe := os.environ.get("OIDC_SERVER_EXE"):
            # original env is needed on windows (possibly only SYSTEMROOT)
            # https://github.com/golang/go/issues/25513
            env = {
                **os.environ,
                "LISTEN_ADDR": "127.0.0.1",
                "LISTEN_PORT": str(self.oidc_port),
            }
            self.oidc_server_proc = subprocess.Popen([oidc_server_exe], env=env)

        karaberus_bin = os.environ["KARABERUS_BIN"]
        db = pathlib.Path(os.environ["KARABERUS_S3_TEST_DB_FILE"])
        db.unlink(missing_ok=True)

        # original env is needed on windows (possibly only SYSTEMROOT)
        # https://github.com/golang/go/issues/25513
        env = {
            **os.environ,
            "KARABERUS_DB_FILE": str(db),
            "KARABERUS_S3_ENDPOINT": f"127.0.0.1:{self.s3_port}",
            "KARABERUS_S3_KEYID": "keyid",
            "KARABERUS_S3_SECRET": "secret",
            "KARABERUS_S3_SECURE": os.environ.get("KARABERUS_S3_SECURE", ""),
            "KARABERUS_S3_BUCKET_NAME": os.environ.get("KARABERUS_S3_BUCKET_NAME", ""),
            "KARABERUS_LISTEN_PORT": str(self.port),
            "KARABERUS_OIDC_ISSUER": f"http://127.0.0.1:{self.oidc_port}",
            "KARABERUS_OIDC_CLIENT_ID": "web",
            "KARABERUS_OIDC_CLIENT_SECRET": "secret",
            "KARABERUS_OIDC_GROUPS_CLAIM": "groups",
            "KARABERUS_OIDC_ADMIN_GROUP": "admin",
            "KARABERUS_OIDC_JWT_SIGN_KEY": secrets.token_hex(),
        }

        user = "testadmin"
        # doesn't matter if the user exists, so we don't check the return code
        _ = subprocess.run(
            [*exe_wrapper(), karaberus_bin, "create-user", "--admin", user],
            env=env,
            stdout=subprocess.PIPE,
        )

        create_token = subprocess.run(
            [
                karaberus_bin,
                "create-token",
                user,
                "testtoken",
            ],
            env=env,
            stdout=subprocess.PIPE,
        )
        create_token.check_returncode()
        self.token = create_token.stdout.decode().strip()

        self.wait_oidc_ready()
        self.wait_s3_ready()

        self.proc = subprocess.Popen([karaberus_bin], env=env)

        self.wait_ready()

    def wait_oidc_ready(self):
        while True:
            try:
                request.urlopen(
                    f"http://127.0.0.1:{self.oidc_port}/.well-known/openid-configuration"
                )
                break
            except URLError:
                time.sleep(0.1)

    def wait_s3_ready(self):
        while True:
            try:
                request.urlopen(
                    f"http://127.0.0.1:{self.s3_port}",
                    timeout=2,
                )
                break
            except URLError:
                time.sleep(0.1)

    def wait_ready(self):
        while True:
            try:
                request.urlopen(f"{self.base_url}/readyz")
                break
            except URLError:
                time.sleep(0.1)

    def stop_karaberus(self) -> None:
        if self.proc is not None:
            self.proc.kill()
        if self.oidc_server_proc is not None:
            self.oidc_server_proc.kill()
        if self.gofakes3_proc is not None:
            self.gofakes3_proc.kill()

    def get(
        self,
        path: str,
        headers: dict[str, str] | None = None,
    ) -> http.client.HTTPResponse:
        url = f"{self.base_url}{path}"

        if headers is None:
            headers = {}

        headers["Authorization"] = f"Bearer {self.token}"

        req = request.Request(url, headers=headers, method="GET")
        return request.urlopen(req, timeout=5)

    def json_request(
        self, method: str, path: str, data: KaraberusInputTypes
    ) -> http.client.HTTPResponse:
        url = f"{self.base_url}{path}"
        headers = {
            "Authorization": f"Bearer {self.token}",
            "Content-Type": "application/json",
        }
        req = request.Request(url, data=json_body(data), headers=headers, method=method)
        return request.urlopen(req, timeout=5)

    def upload_file(
        self, method: str, path: str, file: pathlib.Path
    ) -> http.client.HTTPResponse:
        url = f"{self.base_url}{path}"
        headers = {
            "Authorization": f"Bearer {self.token}",
            "Content-Type": "application/octet-stream",
            "Filename": file.name,
        }

        with file.open("rb") as fd:
            req = request.Request(url, data=fd, headers=headers, method=method)
            return request.urlopen(req, timeout=5)


class KaraberusKaraDB(TypedDict):
    ID: int
    VideoSize: int
    VideoCRC32: int
    InstrumentalSize: int
    InstrumentalCRC32: int
    SubtitlesSize: int
    SubtitlesCRC32: int


class KaraberusKaraResponse(TypedDict):
    kara: KaraberusKaraDB


class DakaraCheckResults(TypedDict):
    passed: bool
    duration: int
    size: int
    crc32: int


class DakaraCheckSubResults(TypedDict):
    passed: bool
    lyrics: str
    size: int
    crc32: int


class CheckKaraResults(TypedDict):
    Video: DakaraCheckResults
    Instrumental: DakaraCheckResults
    Subtitles: DakaraCheckSubResults


class UploadOutput(TypedDict):
    file_id: int
    check_results: CheckKaraResults


class Font(TypedDict):
    ID: int
    Name: str


class FontUpload(TypedDict):
    font: Font


class TestKaraberus(unittest.TestCase):
    karaberus: ClassVar[KaraberusInstance]

    # windows CI uses python 3.9 which doesn’t have typing.override
    # @override
    @classmethod
    def setUpClass(cls) -> None:
        cls.karaberus = KaraberusInstance()
        cls.karaberus.launch_karaberus()

    # @override
    @classmethod
    def tearDownClass(cls) -> None:
        cls.karaberus.stop_karaberus()

    def compare_files(self, orig: IO[bytes], download: IO[bytes]) -> None:
        chunk_size = 1024 * 1024
        while True:
            orig_bytes = orig.read(chunk_size)
            download_bytes = download.read(chunk_size)
            self.assertEqual(orig_bytes, download_bytes)

            if len(orig_bytes) == 0:
                break

    def test_upload(self) -> None:
        kara: KaraberusKara = {
            "title": "aaaa",
            "title_aliases": [],
            "authors": [],
            "artists": [],
            "source_media": 0,
            "song_order": 0,
            "medias": [],
            "audio_tags": [],
            "video_tags": [],
            "comment": "",
            "version": "",
            "language": "",
            "karaoke_creation_time": None,
            "is_hardsub": None,
        }

        resp = self.karaberus.json_request("POST", "/api/kara", kara)
        kara_data: KaraberusKaraResponse = json.load(resp)

        tests_dir = pathlib.Path(__file__).parent
        generated_tests = pathlib.Path(os.environ["KARABERUS_TEST_DIR_GENERATED"])

        # upload video file
        kara_upload_path = f"/api/kara/{kara_data['kara']['ID']}/upload/video"
        video_test_file = generated_tests / "karaberus_test.mkv"
        video_crc32 = calculate_crc32(video_test_file)
        video_size = video_test_file.stat().st_size

        resp = self.karaberus.upload_file("PUT", kara_upload_path, video_test_file)
        upload_data: UploadOutput = json.load(resp)

        video_duration = 0 if os.environ.get("NO_NATIVE_DEPS") else 30

        video_check = upload_data["check_results"]["Video"]
        self.assertTrue(video_check["passed"])
        self.assertEqual(video_check["duration"], video_duration)

        # upload instrumental file
        kara_upload_path = f"/api/kara/{kara_data['kara']['ID']}/upload/inst"
        inst_test_file = generated_tests / "karaberus_test.opus"
        inst_crc32 = calculate_crc32(inst_test_file)
        inst_size = inst_test_file.stat().st_size

        resp = self.karaberus.upload_file("PUT", kara_upload_path, inst_test_file)
        upload_data = json.load(resp)

        # shouldn't have changed
        video_check = upload_data["check_results"]["Video"]
        self.assertTrue(video_check["passed"])
        self.assertEqual(video_check["duration"], video_duration)

        audio_check = upload_data["check_results"]["Instrumental"]
        self.assertTrue(audio_check["passed"])
        # duration isn't really used

        # upload subtitles file
        kara_upload_path = f"/api/kara/{kara_data['kara']['ID']}/upload/sub"
        sub_test_file = tests_dir / "test.ass"
        sub_crc32 = calculate_crc32(sub_test_file)
        sub_size = sub_test_file.stat().st_size

        resp = self.karaberus.upload_file("PUT", kara_upload_path, sub_test_file)
        upload_data = json.load(resp)

        # shouldn't have changed
        video_check = upload_data["check_results"]["Video"]
        self.assertTrue(video_check["passed"])
        self.assertEqual(video_check["duration"], video_duration)

        audio_check = upload_data["check_results"]["Instrumental"]
        self.assertTrue(audio_check["passed"])
        # duration isn't really used

        lyrics = "" if os.environ.get("NO_NATIVE_DEPS") else "It's a small ASS."

        sub_check = upload_data["check_results"]["Subtitles"]
        self.assertTrue(sub_check["passed"])
        self.assertEqual(sub_check["lyrics"], lyrics)

        kara_info = f"/api/kara/{kara_data['kara']['ID']}"
        resp = self.karaberus.get(kara_info)
        kara_data = json.load(resp)

        self.assertEqual(kara_data["kara"]["VideoSize"], video_size)
        self.assertEqual(kara_data["kara"]["VideoCRC32"], video_crc32)
        self.assertEqual(kara_data["kara"]["InstrumentalSize"], inst_size)
        self.assertEqual(kara_data["kara"]["InstrumentalCRC32"], inst_crc32)
        self.assertEqual(kara_data["kara"]["SubtitlesSize"], sub_size)
        self.assertEqual(kara_data["kara"]["SubtitlesCRC32"], sub_crc32)

        # compare local files with uploaded files
        video_download = f"/api/kara/{kara_data['kara']['ID']}/download/video"
        resp = self.karaberus.get(video_download)
        with video_test_file.open("rb") as fd:
            self.compare_files(fd, resp)

        inst_download = f"/api/kara/{kara_data['kara']['ID']}/download/inst"
        resp = self.karaberus.get(inst_download)
        with inst_test_file.open("rb") as fd:
            self.compare_files(fd, resp)

        sub_download = f"/api/kara/{kara_data['kara']['ID']}/download/sub"
        resp = self.karaberus.get(sub_download)
        with sub_test_file.open("rb") as fd:
            self.compare_files(fd, resp)

    def test_font_upload(self) -> None:
        tests_dir = pathlib.Path(__file__).parent
        font_file = tests_dir / "KaraberusTestFont.ttf"
        resp = self.karaberus.upload_file("POST", "/api/font", font_file)
        font_data: FontUpload = json.load(resp)

        font = font_data["font"]
        self.assertEqual(font["Name"], font_file.name)

        font_download = f"/api/font/{font['ID']}/download"
        resp = self.karaberus.get(font_download)
        with font_file.open("rb") as fd:
            self.compare_files(fd, resp)

    def test_etag(self) -> None:
        out = self.karaberus.get("/api/kara")
        etag = out.headers["ETag"]
        self.assertEqual(out.status, 200)

        try:
            out = self.karaberus.get("/api/kara", {"If-None-Match": etag})
            self.assertEqual(out.status, 304)
        except HTTPError as e:
            self.assertEqual(e.status, 304)


if __name__ == "__main__":
    _ = unittest.main()
