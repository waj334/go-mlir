package mlir

/*
#include <mlir-c/AffineMap.h>
*/
import "C"
import "unsafe"

type AffineMap C.MlirAffineMap

func (a AffineMap) Raw() C.MlirAffineMap { return C.MlirAffineMap(a) }
func (a AffineMap) IsNull() bool         { return bool(C.mlirAffineMapIsNull(C.MlirAffineMap(a))) }
func (a AffineMap) Equal(b AffineMap) bool {
	return bool(C.mlirAffineMapEqual(C.MlirAffineMap(a), C.MlirAffineMap(b)))
}

func (a AffineMap) Print() {
	C.mlirAffineMapDump(C.MlirAffineMap(a))
}

func NewAffineMap(ctx Context, dimCount int, symbolCount int, affineExprs []AffineExpr) AffineMap {
	var cAffineExprs *C.MlirAffineExpr
	if len(affineExprs) > 0 {
		cAffineExprs = (*C.MlirAffineExpr)(unsafe.Pointer(&affineExprs[0]))
	}
	return AffineMap(C.mlirAffineMapGet(ctx.Raw(), C.intptr_t(dimCount), C.intptr_t(symbolCount), C.intptr_t(len(affineExprs)), cAffineExprs))
}

func NewEmptyAffineMap(ctx Context) AffineMap {
	return AffineMap(C.mlirAffineMapEmptyGet(ctx.Raw()))
}

func NewZeroResultAffineMap(ctx Context, dimCount int, symbolCount int) AffineMap {
	return AffineMap(C.mlirAffineMapZeroResultGet(ctx.Raw(), C.intptr_t(dimCount), C.intptr_t(symbolCount)))
}

func NewConstantAffineMap(ctx Context, val int64) AffineMap {
	return AffineMap(C.mlirAffineMapConstantGet(ctx.Raw(), C.int64_t(val)))
}

func NewMultiDimIdentityAffineMap(ctx Context, numDims int) AffineMap {
	return AffineMap(C.mlirAffineMapMultiDimIdentityGet(ctx.Raw(), C.intptr_t(numDims)))
}

func NewMinorIdentityAffineMap(ctx Context, dims int, results int) AffineMap {
	return AffineMap(C.mlirAffineMapMinorIdentityGet(ctx.Raw(), C.intptr_t(dims), C.intptr_t(results)))
}

func NewPermutationAffineMap(ctx Context, permutation []uint) AffineMap {
	var cPermutation *C.unsigned
	if len(permutation) > 0 {
		cPermutation = (*C.unsigned)(unsafe.Pointer(&permutation[0]))
	}
	return AffineMap(C.mlirAffineMapPermutationGet(ctx.Raw(), C.intptr_t(len(permutation)), cPermutation))
}

func (a AffineMap) IsIdentity() bool {
	return bool(C.mlirAffineMapIsIdentity(C.MlirAffineMap(a)))
}

func (a AffineMap) IsMinorIdentity() bool {
	return bool(C.mlirAffineMapIsMinorIdentity(C.MlirAffineMap(a)))
}

func (a AffineMap) IsEmpty() bool {
	return bool(C.mlirAffineMapIsEmpty(C.MlirAffineMap(a)))
}

func (a AffineMap) IsSingleConstant() bool {
	return bool(C.mlirAffineMapIsSingleConstant(C.MlirAffineMap(a)))
}

func (a AffineMap) SingleConstantResult() int64 {
	return int64(C.mlirAffineMapGetSingleConstantResult(C.MlirAffineMap(a)))
}

func (a AffineMap) NumDims() int {
	return int(C.mlirAffineMapGetNumDims(C.MlirAffineMap(a)))
}

func (a AffineMap) NumSymbols() int {
	return int(C.mlirAffineMapGetNumSymbols(C.MlirAffineMap(a)))
}

func (a AffineMap) NumResults() int {
	return int(C.mlirAffineMapGetNumResults(C.MlirAffineMap(a)))
}

func (a AffineMap) Result(pos int) AffineExpr {
	return AffineExpr(C.mlirAffineMapGetResult(C.MlirAffineMap(a), C.intptr_t(pos)))
}

func (a AffineMap) NumInputs() int {
	return int(C.mlirAffineMapGetNumInputs(C.MlirAffineMap(a)))
}

func (a AffineMap) IsProjectedPermutation() bool {
	return bool(C.mlirAffineMapIsProjectedPermutation(C.MlirAffineMap(a)))
}

func (a AffineMap) IsPermutation() bool {
	return bool(C.mlirAffineMapIsPermutation(C.MlirAffineMap(a)))
}

func (a AffineMap) SubMap(resultPos []int) AffineMap {
	var cResultPos *C.intptr_t
	if len(resultPos) > 0 {
		cResultPos = (*C.intptr_t)(unsafe.Pointer(&resultPos[0]))
	}
	return AffineMap(C.mlirAffineMapGetSubMap(C.MlirAffineMap(a), C.intptr_t(len(resultPos)), cResultPos))
}

func (a AffineMap) MajorSubMap(numResults int) AffineMap {
	return AffineMap(C.mlirAffineMapGetMajorSubMap(C.MlirAffineMap(a), C.intptr_t(numResults)))
}

func (a AffineMap) MinorSubMap(numResults int) AffineMap {
	return AffineMap(C.mlirAffineMapGetMinorSubMap(C.MlirAffineMap(a), C.intptr_t(numResults)))
}
