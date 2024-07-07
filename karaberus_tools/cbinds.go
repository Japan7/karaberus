package karaberus_tools

/*
#cgo pkg-config: karaberus_tools
#include <karaberus_tools.h>
#include <unistd.h>
#include <errno.h>
#include <stdint.h>
#include <string.h>

int AVIORead(void *obj, uint8_t *buf, int n);
int64_t AVIOSeek(void *obj, int64_t offset, int whence);

static inline struct dakara_check_results *karaberus_dakara_check(void *obj) {
  return karaberus_dakara_check_avio(obj, AVIORead, AVIOSeek);
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
	C.memcpy(C.CBytes(rbuf), unsafe.Pointer(buf), C.size_t(nread))
	return C.int(nread)
}

//export AVIOSeek
func AVIOSeek(opaque unsafe.Pointer, offset C.int64_t, whence C.int) C.int64_t {
	obj := pointer.Restore(opaque).(*minio.Object)
	pos, err := obj.Seek(int64(offset), int(whence))
	if err != nil {
		panic(err)
	}
	return C.int64_t(pos)
}

type DakaraCheckResultsOutput struct {
	Passed bool `json:"passed" example:"true" doc:"true if file passed all checks"`
}

func DakaraCheckResults(obj *minio.Object) DakaraCheckResultsOutput {
	res := C.karaberus_dakara_check(pointer.Save(obj))
	defer C.karaberus_dakara_check_results_free(res)

	out := DakaraCheckResultsOutput{Passed: bool(res.passed)}
	return out
}
