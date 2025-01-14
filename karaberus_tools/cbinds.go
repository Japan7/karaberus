//go:build cgo

package karaberus_tools

/*
#cgo pkg-config: dakara_check
#include "karaberus_tools.h"
#include <dakara_check.h>
#include <unistd.h>
#include <errno.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

int AVIORead(void *obj, uint8_t *buf, int n);
int64_t AVIOSeek(void *obj, int64_t offset, int whence);

static inline karaberus_reports karaberus_dakara_check(void *obj, bool video_stream) {
  return karaberus_dakara_check_avio(obj, AVIORead, AVIOSeek, video_stream);
}
*/
import "C"
import (
	"errors"
	"io"
	"runtime/cgo"
	"unsafe"
)

//export AVIORead
func AVIORead(opaque unsafe.Pointer, buf *C.uint8_t, n C.int) C.int {
	h := *(*cgo.Handle)(opaque)
	objbuf := h.Value().(ObjectBuf)

	var rbuf []byte
	if int(n) < len(objbuf.Buffer) {
		rbuf = objbuf.Buffer[:n]
	} else {
		rbuf = objbuf.Buffer
	}
	nread, err := objbuf.Object.Read(rbuf)
	if errors.Is(err, io.EOF) {
		if nread == 0 {
			return C.AVERROR_EOF
		}
	} else if err != nil {
		panic(err)
	}
	c_rbuf := C.CBytes(rbuf)
	defer C.free(c_rbuf)
	C.memcpy(unsafe.Pointer(buf), c_rbuf, C.size_t(nread))
	return C.int(nread)
}

//export AVIOSeek
func AVIOSeek(opaque unsafe.Pointer, offset C.int64_t, whence C.int) C.int64_t {
	h := *(*cgo.Handle)(opaque)
	objbuf := h.Value().(ObjectBuf)

	if whence == C.AVSEEK_SIZE {
		return C.int64_t(objbuf.Size)
	}
	pos, err := objbuf.Object.Seek(int64(offset), int(whence))
	if err != nil {
		panic(err)
	}
	return C.int64_t(pos)
}

type ObjectBuf struct {
	Object io.ReadSeeker
	Buffer []byte
	Size   int64
}

func NewObjectBuf(obj io.ReadSeeker, size int64) ObjectBuf {
	return ObjectBuf{
		Object: obj,
		Buffer: make([]byte, C.KARABERUS_BUFSIZE),
		Size:   size,
	}
}

func stringForReportLevel(report C.karaberus_report) string {
	switch report.error_level {
	case C.K_ERROR:
		return "error: "
	case C.K_WARNING:
		return "warning: "
	case C.K_INFO:
		return "info: "
	default:
		return ""
	}
}

func stringForReport(report C.karaberus_report) string {
	return stringForReportLevel(report) + C.GoString(report.message)
}

func DakaraCheckResults(obj io.ReadSeeker, ftype string, size int64) DakaraCheckResultsOutput {
	video_stream := ftype == "video"
	object_buf := NewObjectBuf(obj, size)
	handle := cgo.NewHandle(object_buf)
	defer handle.Delete()
	res := C.karaberus_dakara_check(unsafe.Pointer(&handle), C.bool(video_stream))
	defer C.free_reports(res)

	messages := make([]string, res.n_reports)
	reports := unsafe.Slice(res.reports, res.n_reports)
	for i, report := range reports {
		messages[i] = stringForReport(report)
	}

	passed := !bool(res.failed)
	out := DakaraCheckResultsOutput{
		Passed:   passed,
		Duration: int32(res.duration),
		Messages: messages,
	}
	return out
}

func DakaraCheckSub(obj io.ReadSeeker, size int64) (DakaraCheckSubResultsOutput, error) {
	out := DakaraCheckSubResultsOutput{
		Lyrics: "",
		Passed: false,
	}

	buf := make([]byte, size)
	_, err := obj.Seek(0, 0)
	if err != nil {
		return out, err
	}
	_, err = io.ReadFull(obj, buf)
	if err != nil {
		return out, err
	}

	cbuf := C.CString(string(buf))
	defer C.free(unsafe.Pointer(cbuf))
	res := C.karaberus_check_sub(cbuf, C.size_t(len(buf)))
	defer C.karaberus_sub_reports_free(res)
	if res == nil {
		return out, errors.New("failed to parse subtitle file")
	}
	if res.io_error {
		return out, errors.New("IO error while reading sub file")
	}
	out.Lyrics = C.GoString(res.lyrics)
	out.Passed = true
	return out, nil
}
