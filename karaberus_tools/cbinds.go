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
	"unsafe"

	"github.com/mattn/go-pointer"
	"github.com/minio/minio-go/v7"
)

//export AVIORead
func AVIORead(opaque unsafe.Pointer, buf *C.uint8_t, n C.int) C.int {
	objbuf := pointer.Restore(opaque).(ObjectBuf)
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
	objbuf := pointer.Restore(opaque).(ObjectBuf)
	if whence == C.AVSEEK_SIZE {
		stat, err := objbuf.Object.Stat()
		if err != nil {
			panic(err)
		}
		return C.int64_t(stat.Size)
	}
	pos, err := objbuf.Object.Seek(int64(offset), int(whence))
	if err != nil {
		panic(err)
	}
	return C.int64_t(pos)
}

type ObjectBuf struct {
	Object *minio.Object
	Buffer []byte
}

func NewObjectBuf(obj *minio.Object) ObjectBuf {
	return ObjectBuf{
		Object: obj,
		Buffer: make([]byte, C.KARABERUS_BUFSIZE),
	}
}

func DakaraCheckResults(obj *minio.Object, ftype string) DakaraCheckResultsOutput {
	video_stream := ftype == "video"
	object_buf := NewObjectBuf(obj)
	ptr := pointer.Save(object_buf)
	defer pointer.Unref(ptr)
	res := C.karaberus_dakara_check(ptr, C.bool(video_stream))
	defer C.free_reports(res)

	passed := !bool(res.failed)
	out := DakaraCheckResultsOutput{
		Passed:   passed,
		Duration: int32(res.duration),
	}
	return out
}
