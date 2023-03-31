package tokenizer

// TODO packaging: how do we build the rust lib for distribution?

/*
#cgo LDFLAGS: ./lib/libtokenizer.a -ldl -lstdc++
#include <stdlib.h>
#include "./lib/tokenizer.h"
*/
import "C"

// NOTE: There should be NO space between the comments and the `import "C"` line.
import (
	"unsafe"
)

func Encode(str string) []uint32 {
	cStr := C.CString(str)
	defer C.free(unsafe.Pointer(cStr))
	var len C.uint
	res := C.encode(cStr, &len)
	defer C.free(unsafe.Pointer(res))
	slice := unsafe.Slice(res, len)

	tokenIDs := make([]uint32, len)
	for i, v := range slice {
		tokenIDs[i] = uint32(v)
	}
	return tokenIDs
}

func Decode(tokenIDs []uint32) string {
	len := C.uint(len(tokenIDs))
	res := C.decode((*C.uint)(unsafe.Pointer(&tokenIDs[0])), len)
	defer C.free(unsafe.Pointer(res))
	return C.GoString(res)
}