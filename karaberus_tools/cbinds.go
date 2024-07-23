//go:build cgo

package karaberus_tools

/*
#cgo pkg-config: karaberus_tools
#include <karaberus_tools.h>
#include <dakara_check.h>
#include <unistd.h>
#include <errno.h>
#include <stdint.h>
#include <string.h>

int AVIORead(void *obj, uint8_t *buf, int n);
int64_t AVIOSeek(void *obj, int64_t offset, int whence);

static inline void karaberus_dakara_check(void *obj, dakara_check_results *res, bool needs_duration) {
  karaberus_dakara_check_avio(obj, AVIORead, AVIOSeek, res, needs_duration);
}
*/
import "C"
import (
	"errors"
	"io"
	"strings"
	"unsafe"

	"github.com/mattn/go-pointer"
	"github.com/minio/minio-go/v7"
)

//export AVIORead
func AVIORead(opaque unsafe.Pointer, buf *C.uint8_t, n C.int) C.int {
	obj := pointer.Restore(opaque).(*minio.Object)
	rbuf := make([]byte, n)
	nread, err := obj.Read(rbuf)
	if errors.Is(err, io.EOF) {
		if nread == 0 {
			return C.AVERROR_EOF
		}
	} else if err != nil {
		panic(err)
	}
	C.memcpy(unsafe.Pointer(buf), C.CBytes(rbuf), C.size_t(nread))
	return C.int(nread)
}

//export AVIOSeek
func AVIOSeek(opaque unsafe.Pointer, offset C.int64_t, whence C.int) C.int64_t {
	obj := pointer.Restore(opaque).(*minio.Object)
	if whence == C.AVSEEK_SIZE {
		stat, err := obj.Stat()
		if err != nil {
			panic(err)
		}
		return C.int64_t(stat.Size)
	}
	pos, err := obj.Seek(int64(offset), int(whence))
	if err != nil {
		panic(err)
	}
	return C.int64_t(pos)
}

func DakaraCheckResults(obj *minio.Object) DakaraCheckResultsOutput {
	res := C.dakara_check_results{}
	stat, _ := obj.Stat()
	ftype, _, _ := strings.Cut(stat.Key, "/")
	needs_duration := ftype == "video"
	C.karaberus_dakara_check(pointer.Save(obj), &res, C.bool(needs_duration))

	passed := true
	for i := 0; i < len(res.report); i++ {
		passed = passed && (res.report[i] == 0)
	}
	out := DakaraCheckResultsOutput{Passed: passed}
	return out
}
