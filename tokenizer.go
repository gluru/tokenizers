package tokenizers

// TODO packaging: how do we build the rust lib for distribution?
//    Since we pre-compile rust maybe we can bundle different versions.
// #cgo LDFLAGS: ${SRCDIR}/libtokenizers-${GOOS}-${GOARCH}.a -ldl -lstdc++

/*
#cgo LDFLAGS: ${SRCDIR}/libtokenizers-linux-amd64.a -ldl -lstdc++
#include <stdlib.h>
#include "tokenizers.h"
*/
import "C"

// NOTE: There should be NO space between the comments and the `import "C"` line.
import (
	"io"
	"unsafe"
)

type Tokenizer struct {
	tokenizer unsafe.Pointer
}

type TruncationDirection int

const (
	TruncationDirectionLeft TruncationDirection = iota
	TruncationDirectionRight
)

var _ io.Closer = (*Tokenizer)(nil)

func FromBytes(data []byte) (*Tokenizer, error) {
	tokenizer := C.from_bytes((*C.uchar)(unsafe.Pointer(&data[0])), C.uint(len(data)))
	return &Tokenizer{tokenizer: tokenizer}, nil
}

func FromBytesWithTruncation(data []byte, maxLen uint32, dir TruncationDirection) (*Tokenizer, error) {
	tokenizer := C.from_bytes_with_truncation((*C.uchar)(unsafe.Pointer(&data[0])), C.uint(len(data)), C.uint(maxLen), C.uchar(dir))
	return &Tokenizer{tokenizer: tokenizer}, nil
}

func FromFile(path string) (*Tokenizer, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	tokenizer, err := C.from_file(cPath)
	if err != nil {
		return nil, err
	}
	return &Tokenizer{tokenizer: tokenizer}, nil
}

func (t *Tokenizer) Close() error {
	C.free_tokenizer(t.tokenizer)
	t.tokenizer = nil
	return nil
}

func (t *Tokenizer) Encode(str string, addSpecialTokens bool) []uint32 {
	cStr := C.CString(str)
	defer C.free(unsafe.Pointer(cStr))
	var len C.uint
	res := C.encode(t.tokenizer, cStr, &len, C.bool(addSpecialTokens))
	if len > 0 {
		// can't dealloc nil
		defer C.free(unsafe.Pointer(res))
	}
	slice := unsafe.Slice(res, len)

	tokenIDs := make([]uint32, len)
	for i, v := range slice {
		tokenIDs[i] = uint32(v)
	}
	return tokenIDs
}

func (t *Tokenizer) Decode(tokenIDs []uint32, skipSpecialTokens bool) string {
	if len(tokenIDs) == 0 {
		return ""
	}
	len := C.uint(len(tokenIDs))
	res := C.decode(t.tokenizer, (*C.uint)(unsafe.Pointer(&tokenIDs[0])), len, C.bool(skipSpecialTokens))
	defer C.free(unsafe.Pointer(res))
	return C.GoString(res)
}

func (t *Tokenizer) VocabSize() uint32 {
	return uint32(C.vocab_size(t.tokenizer))
}
