#!/usr/bin/env python3
"""
Tests for file uploads/downloads.
We can't test it in the go test module because the mocked server panics on downloads.
Otherwise no reason to use this script rather than the go tests to test the API.
"""

import http.client
import json
import os
import pathlib
import secrets
import subprocess
import time
import unittest
import urllib.request as request
from typing import IO, ClassVar, TypedDict
from urllib.error import URLError


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


def json_body(data: KaraberusInputTypes) -> bytes:
    return json.dumps(data, separators=(",", ":")).encode()


class KaraberusInstance:
    def __init__(self) -> None:
        self.port = 8889
        self.proc: None | subprocess.Popen[bytes] = None
        self.token: None | str = None
        self.base_url = f"http://127.0.0.1:{self.port}"

    def launch_karaberus(self) -> None:
        karaberus_bin = os.environ["KARABERUS_BIN"]
        db = pathlib.Path(os.environ["KARABERUS_S3_TEST_DB_FILE"])
        db.unlink(missing_ok=True)
        os.environ["KARABERUS_DB_FILE"] = str(db)
        os.environ["KARABERUS_LISTEN_HOST"] = "127.0.0.1"
        os.environ["KARABERUS_LISTEN_PORT"] = str(self.port)
        env = {
            "KARABERUS_DB_FILE": str(db),
            "KARABERUS_S3_ENDPOINT": os.environ["KARABERUS_S3_ENDPOINT"],
            "KARABERUS_S3_KEYID": os.environ["KARABERUS_S3_KEYID"],
            "KARABERUS_S3_SECRET": os.environ["KARABERUS_S3_SECRET"],
            "KARABERUS_S3_SECURE": os.environ.get("KARABERUS_S3_SECURE", ""),
            "KARABERUS_S3_BUCKET_NAME": os.environ.get("KARABERUS_S3_BUCKET_NAME", ""),
            "KARABERUS_LISTEN_PORT": str(self.port),
            # dummy oidc issuer, shouldn't matter
            "KARABERUS_OIDC_ISSUER": "http://localhost:9998",
            "KARABERUS_OIDC_CLIENT_ID": "web",
            "KARABERUS_OIDC_CLIENT_SECRET": "secret",
            "KARABERUS_OIDC_GROUPS_CLAIM": "groups",
            "KARABERUS_OIDC_ADMIN_GROUP": "admin",
            "KARABERUS_OIDC_JWT_SIGN_KEY": secrets.token_hex(),
        }

        user = "testadmin"
        # doesn't matter if the user exists, so we don't check the return code
        subprocess.run(
            [karaberus_bin, "create-user", "--admin", user],
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

        self.proc = subprocess.Popen([karaberus_bin], env=env)

        self.wait_ready()

    def wait_ready(self):
        while True:
            try:
                return request.urlopen(f"{self.base_url}/readyz")
            except URLError:
                time.sleep(0.1)

    def stop_karaberus(self) -> None:
        if self.proc is not None:
            self.proc.kill()

    def get(self, path: str) -> http.client.HTTPResponse:
        url = f"{self.base_url}{path}"
        headers = {
            "Authorization": f"Bearer {self.token}",
        }
        req = request.Request(url, headers=headers, method="GET")
        return request.urlopen(req)

    def json_request(
        self, method: str, path: str, data: KaraberusInputTypes
    ) -> http.client.HTTPResponse:
        url = f"{self.base_url}{path}"
        headers = {
            "Authorization": f"Bearer {self.token}",
            "Content-Type": "application/json",
        }
        req = request.Request(url, data=json_body(data), headers=headers, method=method)
        return request.urlopen(req)

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
            return request.urlopen(req)


class KaraberusKaraDB(TypedDict):
    ID: int


class KaraberusKaraResponse(TypedDict):
    kara: KaraberusKaraDB


class DakaraCheckResults(TypedDict):
    passed: bool
    duration: int


class DakaraCheckSubResults(TypedDict):
    passed: bool
    lyrics: str


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

    @classmethod
    def setUpClass(cls) -> None:
        cls.karaberus = KaraberusInstance()
        cls.karaberus.launch_karaberus()

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
        resp = self.karaberus.upload_file("PUT", kara_upload_path, video_test_file)
        upload_data: UploadOutput = json.load(resp)

        video_check = upload_data["check_results"]["Video"]
        self.assertTrue(video_check["passed"])
        self.assertEqual(video_check["duration"], 30)

        # upload instrumental file
        kara_upload_path = f"/api/kara/{kara_data['kara']['ID']}/upload/inst"
        inst_test_file = generated_tests / "karaberus_test.opus"
        resp = self.karaberus.upload_file("PUT", kara_upload_path, inst_test_file)
        upload_data = json.load(resp)

        # shouldn't have changed
        video_check = upload_data["check_results"]["Video"]
        self.assertTrue(video_check["passed"])
        self.assertEqual(video_check["duration"], 30)

        audio_check = upload_data["check_results"]["Instrumental"]
        self.assertTrue(audio_check["passed"])
        # duration isn't really used

        # upload subtitles file
        kara_upload_path = f"/api/kara/{kara_data['kara']['ID']}/upload/sub"
        sub_test_file = tests_dir / "test.ass"
        resp = self.karaberus.upload_file("PUT", kara_upload_path, sub_test_file)
        upload_data = json.load(resp)

        # shouldn't have changed
        video_check = upload_data["check_results"]["Video"]
        self.assertTrue(video_check["passed"])
        self.assertEqual(video_check["duration"], 30)

        audio_check = upload_data["check_results"]["Instrumental"]
        self.assertTrue(audio_check["passed"])
        # duration isn't really used

        sub_check = upload_data["check_results"]["Subtitles"]
        self.assertTrue(sub_check["passed"])
        self.assertEqual(sub_check["lyrics"], "It's a small ASS.")

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

    def test_sub_upload(self) -> None:
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


if __name__ == "__main__":
    unittest.main()
