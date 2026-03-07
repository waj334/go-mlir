package mlir

import "C"
import (
	"runtime/cgo"
	"strings"
	"unsafe"
)

/*
#include <stdlib.h>

#include <mlir-c/Support.h>

extern void goMlirStringCallback(MlirStringRef, void*);
static inline MlirStringCallback getGoMlirStringCallback(void) {
  return (MlirStringCallback)goMlirStringCallback;
}
*/
import "C"

//===----------------------------------------------------------------------===//
// MlirStringRef.
//===----------------------------------------------------------------------===//

type StringRef struct {
	ref   C.MlirStringRef
	owned bool
}

func (s StringRef) raw() C.MlirStringRef {
	return s.ref
}

func (s StringRef) String() string {
	return C.GoStringN(s.ref.data, C.int(s.ref.length))
}

func (s StringRef) Destroy() {
	if !s.owned || s.ref.data == nil {
		return
	}
	C.free(unsafe.Pointer(s.ref.data))
}

func NewStringRef(value string) StringRef {
	ptr := C.CString(value)
	return StringRef{
		ref:   C.MlirStringRef{data: ptr, length: C.size_t(len(value))},
		owned: true,
	}
}

func BorrowedStringRef(v C.MlirStringRef) StringRef {
	return StringRef{ref: v, owned: false}
}

type StringCallback func(string)

func NewStringCallback(fn func(string)) StringCallback {
	return StringCallback(fn)
}

func (cb StringCallback) Callback() (C.MlirStringCallback, unsafe.Pointer, func()) {
	handle := cgo.NewHandle(cb)
	return C.getGoMlirStringCallback(), unsafe.Pointer(uintptr(handle)), handle.Delete
}

func collectString(fn func(cb C.MlirStringCallback, ud unsafe.Pointer)) string {
	var b strings.Builder
	cb, ud, cleanup := NewStringCallback(func(s string) {
		b.WriteString(s)
	}).Callback()
	defer cleanup()
	fn(cb, ud)
	return b.String()
}

//export goMlirStringCallback
func goMlirStringCallback(chunk C.MlirStringRef, userdata unsafe.Pointer) {
	handle := cgo.Handle(uintptr(userdata))
	callback := handle.Value().(StringCallback)
	callback(C.GoStringN(chunk.data, C.int(chunk.length)))
}

//===----------------------------------------------------------------------===//
// MlirLogicalResult.
//===----------------------------------------------------------------------===//

type LogicalResult C.MlirLogicalResult

func (result LogicalResult) raw() C.MlirLogicalResult {
	return C.MlirLogicalResult(result)
}

func (result LogicalResult) IsSuccess() bool {
	return bool(C.mlirLogicalResultIsSuccess(C.MlirLogicalResult(result)))
}

func (result LogicalResult) IsFailure() bool {
	return bool(C.mlirLogicalResultIsFailure(C.MlirLogicalResult(result)))
}

func Success() LogicalResult {
	return LogicalResult(C.mlirLogicalResultSuccess())
}

func Failure() LogicalResult {
	return LogicalResult(C.mlirLogicalResultFailure())
}

//===----------------------------------------------------------------------===//
// MlirLlvmThreadPool.
//===----------------------------------------------------------------------===//

type LLVMThreadPool C.MlirLlvmThreadPool

func (pool LLVMThreadPool) raw() C.MlirLlvmThreadPool {
	return C.MlirLlvmThreadPool(pool)
}

func NewLLVMThreadPool() LLVMThreadPool {
	return LLVMThreadPool(C.mlirLlvmThreadPoolCreate())
}

func (pool LLVMThreadPool) Destroy() {
	C.mlirLlvmThreadPoolDestroy(C.MlirLlvmThreadPool(pool))
}

//===----------------------------------------------------------------------===//
// TypeId API.
//===----------------------------------------------------------------------===//

type TypeId C.MlirTypeID

func NewTypeID(ptr unsafe.Pointer) TypeId {
	return TypeId(C.mlirTypeIDCreate(ptr))
}

func (id TypeId) raw() C.MlirTypeID {
	return C.MlirTypeID(id)
}

func (id TypeId) IsNull() bool {
	return bool(C.mlirTypeIDIsNull(C.MlirTypeID(id)))
}

func (id TypeId) Equal(other TypeId) bool {
	return bool(C.mlirTypeIDEqual(C.MlirTypeID(id), C.MlirTypeID(other)))
}

func (id TypeId) HashValue() uintptr {
	return uintptr(C.mlirTypeIDHashValue(C.MlirTypeID(id)))
}

//===----------------------------------------------------------------------===//
// TypeIdAllocator API.
//===----------------------------------------------------------------------===//

type TypeIdAllocator C.MlirTypeIDAllocator

func NewTypeIDAllocator() TypeIdAllocator {
	return TypeIdAllocator(C.mlirTypeIDAllocatorCreate())
}

func (allocator TypeIdAllocator) raw() C.MlirTypeIDAllocator {
	return C.MlirTypeIDAllocator(allocator)
}

func (allocator TypeIdAllocator) Destroy() {
	C.mlirTypeIDAllocatorDestroy(C.MlirTypeIDAllocator(allocator))
}

func (allocator TypeIdAllocator) AllocateTypeID() TypeId {
	return TypeId(C.mlirTypeIDAllocatorAllocateTypeID(C.MlirTypeIDAllocator(allocator)))
}
