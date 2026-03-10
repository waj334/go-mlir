package mlir

/*
#include <stdlib.h>
#include <mlir-c/BuiltinAttributes.h>
*/
import "C"
import "unsafe"

func NewNullAttribute() Attribute {
	return WrapAttribute(C.mlirAttributeGetNull())
}

//===----------------------------------------------------------------------===//
// Location attribute.
//===----------------------------------------------------------------------===//

func AttributeIsALocation(attr AttributeLike) bool {
	return bool(C.mlirAttributeIsALocation(attr.Raw()))
}

//===----------------------------------------------------------------------===//
// Affine map attribute.
//===----------------------------------------------------------------------===//

type AffineMapAttr struct {
	baseAttribute
}

func NewAffineMapAttr(value AffineMap) AffineMapAttr {
	return AffineMapAttr{
		baseAttribute: baseAttribute(C.mlirAffineMapAttrGet(value.Raw())),
	}
}

func AffineMapAttrName() string {
	return BorrowedStringRef(C.mlirAffineMapAttrGetName()).String()
}

func AffineMapAttrTypeId() TypeId {
	return TypeId(C.mlirAffineMapAttrGetTypeID())
}

func (attr AffineMapAttr) Value() AffineMap {
	return AffineMap(C.mlirAffineMapAttrGetValue(attr.Raw()))
}

//===----------------------------------------------------------------------===//
// Array attribute.
//===----------------------------------------------------------------------===//

type ArrayAttr struct {
	baseAttribute
}

func NewArrayAttr[T AttributeLike](ctx Context, elements []T) ArrayAttr {
	var cElements *C.MlirAttribute
	if len(elements) > 0 {
		cElements = &UnwrapAttributeSlice(elements)[0]
	}
	return ArrayAttr{
		baseAttribute: baseAttribute(C.mlirArrayAttrGet(ctx.Raw(), C.intptr_t(len(elements)), cElements)),
	}
}

func ArrayAttrName() string {
	return BorrowedStringRef(C.mlirArrayAttrGetName()).String()
}

func ArrayAttrTypeId() TypeId {
	return TypeId(C.mlirArrayAttrGetTypeID())
}

func (attr ArrayAttr) NumElements() int {
	return int(C.mlirArrayAttrGetNumElements(attr.Raw()))
}

func (attr ArrayAttr) Element(pos int) Attribute {
	return WrapAttribute(C.mlirArrayAttrGetElement(attr.Raw(), C.intptr_t(pos)))
}

//===----------------------------------------------------------------------===//
// Dictionary attribute.
//===----------------------------------------------------------------------===//

type DictionaryAttr struct {
	baseAttribute
}

func NewDictionaryAttr(ctx Context, elements []NamedAttribute) DictionaryAttr {
	var cElements *C.MlirNamedAttribute
	if len(elements) > 0 {
		cElements = &UnwrapNamedAttributeSlice(elements)[0]
	}
	return DictionaryAttr{
		baseAttribute: baseAttribute(C.mlirDictionaryAttrGet(ctx.Raw(), C.intptr_t(len(elements)), cElements)),
	}
}

func DictionaryAttrName() string {
	return BorrowedStringRef(C.mlirDictionaryAttrGetName()).String()
}

func DictionaryAttrTypeId() TypeId {
	return TypeId(C.mlirDictionaryAttrGetTypeID())
}

func (attr DictionaryAttr) NumElements() int {
	return int(C.mlirDictionaryAttrGetNumElements(attr.Raw()))
}

func (attr DictionaryAttr) Element(pos int) NamedAttribute {
	return NamedAttribute(C.mlirDictionaryAttrGetElement(attr.Raw(), C.intptr_t(pos)))
}

func (attr DictionaryAttr) ElementByName(name string) Attribute {
	refName := NewStringRef(name)
	defer refName.Destroy()
	return WrapAttribute(C.mlirDictionaryAttrGetElementByName(attr.Raw(), refName.Raw()))
}

//===----------------------------------------------------------------------===//
// Floating point attribute.
//===----------------------------------------------------------------------===//

type FloatAttr struct {
	baseAttribute
}

func NewFloatAttr(ctx Context, type_ TypeLike, value float64) FloatAttr {
	return FloatAttr{
		baseAttribute: baseAttribute(C.mlirFloatAttrDoubleGet(ctx.Raw(), type_.Raw(), C.double(value))),
	}
}

func NewFloatAttrChecked(loc LocationLike, type_ TypeLike, value float64) FloatAttr {
	return FloatAttr{
		baseAttribute: baseAttribute(C.mlirFloatAttrDoubleGetChecked(loc.Raw(), type_.Raw(), C.double(value))),
	}
}

func FloatAttrName() string {
	return BorrowedStringRef(C.mlirFloatAttrGetName()).String()
}

func FloatAttrTypeId() TypeId {
	return TypeId(C.mlirFloatAttrGetTypeID())
}

func (attr FloatAttr) Value() float64 {
	return float64(C.mlirFloatAttrGetValueDouble(attr.Raw()))
}

//===----------------------------------------------------------------------===//
// Integer attribute.
//===----------------------------------------------------------------------===//

type IntegerAttr struct {
	baseAttribute
}

func NewIntegerAttr(type_ TypeLike, value int64) IntegerAttr {
	return IntegerAttr{
		baseAttribute: baseAttribute(C.mlirIntegerAttrGet(type_.Raw(), C.int64_t(value))),
	}
}

func IntegerAttrName() string {
	return BorrowedStringRef(C.mlirIntegerAttrGetName()).String()
}

func IntegerAttrTypeId() TypeId {
	return TypeId(C.mlirIntegerAttrGetTypeID())
}

func (attr IntegerAttr) ValueInt() int64 {
	return int64(C.mlirIntegerAttrGetValueInt(attr.Raw()))
}

func (attr IntegerAttr) ValueSInt() int64 {
	return int64(C.mlirIntegerAttrGetValueSInt(attr.Raw()))
}

func (attr IntegerAttr) ValueUInt() uint64 {
	return uint64(C.mlirIntegerAttrGetValueUInt(attr.Raw()))
}

//===----------------------------------------------------------------------===//
// Bool attribute.
//===----------------------------------------------------------------------===//

type BoolAttr struct {
	baseAttribute
}

func NewBoolAttr(ctx Context, value bool) BoolAttr {
	var cValue C.int
	if value {
		cValue = 1
	}
	return BoolAttr{
		baseAttribute: baseAttribute(C.mlirBoolAttrGet(ctx.Raw(), cValue)),
	}
}

func (attr BoolAttr) Value() bool {
	return bool(C.mlirBoolAttrGetValue(attr.Raw()))
}

//===----------------------------------------------------------------------===//
// Integer set attribute.
//===----------------------------------------------------------------------===//

type IntegerSetAttr struct {
	baseAttribute
}

// TODO: Implement IntegerSet wrapper
// func NewIntegerSetAttr(set IntegerSet) IntegerSetAttr {
// 	return IntegerSetAttr{
// 		baseAttribute: baseAttribute(C.mlirIntegerSetAttrGet(set.Raw())),
// 	}
// }

func IntegerSetAttrName() string {
	return BorrowedStringRef(C.mlirIntegerSetAttrGetName()).String()
}

func IntegerSetAttrTypeId() TypeId {
	return TypeId(C.mlirIntegerSetAttrGetTypeID())
}

// func (attr IntegerSetAttr) Value() IntegerSet {
// 	return IntegerSet(C.mlirIntegerSetAttrGetValue(attr.Raw()))
// }

//===----------------------------------------------------------------------===//
// Opaque attribute.
//===----------------------------------------------------------------------===//

type OpaqueAttr struct {
	baseAttribute
}

func NewOpaqueAttr(ctx Context, dialectNamespace string, data string, type_ TypeLike) OpaqueAttr {
	cDialectNamespace := NewStringRef(dialectNamespace)
	defer cDialectNamespace.Destroy()
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))
	return OpaqueAttr{
		baseAttribute: baseAttribute(C.mlirOpaqueAttrGet(ctx.Raw(), cDialectNamespace.Raw(), C.intptr_t(len(data)), cData, type_.Raw())),
	}
}

func OpaqueAttrName() string {
	return BorrowedStringRef(C.mlirOpaqueAttrGetName()).String()
}

func OpaqueAttrTypeId() TypeId {
	return TypeId(C.mlirOpaqueAttrGetTypeID())
}

func (attr OpaqueAttr) DialectNamespace() string {
	return BorrowedStringRef(C.mlirOpaqueAttrGetDialectNamespace(attr.Raw())).String()
}

func (attr OpaqueAttr) Data() string {
	return BorrowedStringRef(C.mlirOpaqueAttrGetData(attr.Raw())).String()
}

//===----------------------------------------------------------------------===//
// String attribute.
//===----------------------------------------------------------------------===//

type StringAttr struct {
	baseAttribute
}

func NewStringAttr(ctx Context, str string) StringAttr {
	cStr := NewStringRef(str)
	defer cStr.Destroy()
	return StringAttr{
		baseAttribute: baseAttribute(C.mlirStringAttrGet(ctx.Raw(), cStr.Raw())),
	}
}

func NewStringAttrTyped(type_ TypeLike, str string) StringAttr {
	cStr := NewStringRef(str)
	defer cStr.Destroy()
	return StringAttr{
		baseAttribute: baseAttribute(C.mlirStringAttrTypedGet(type_.Raw(), cStr.Raw())),
	}
}

func StringAttrName() string {
	return BorrowedStringRef(C.mlirStringAttrGetName()).String()
}

func StringAttrTypeId() TypeId {
	return TypeId(C.mlirStringAttrGetTypeID())
}

func (attr StringAttr) Value() string {
	return BorrowedStringRef(C.mlirStringAttrGetValue(attr.Raw())).String()
}

//===----------------------------------------------------------------------===//
// SymbolRef attribute.
//===----------------------------------------------------------------------===//

type SymbolRefAttr struct {
	baseAttribute
}

func NewSymbolRefAttr(ctx Context, symbol string, references []Attribute) SymbolRefAttr {
	cSymbol := NewStringRef(symbol)
	defer cSymbol.Destroy()
	var cReferences *C.MlirAttribute
	if len(references) > 0 {
		cReferences = &UnwrapAttributeSlice(references)[0]
	}
	return SymbolRefAttr{
		baseAttribute: baseAttribute(C.mlirSymbolRefAttrGet(ctx.Raw(), cSymbol.Raw(), C.intptr_t(len(references)), cReferences)),
	}
}

func SymbolRefAttrName() string {
	return BorrowedStringRef(C.mlirSymbolRefAttrGetName()).String()
}

func SymbolRefAttrTypeId() TypeId {
	return TypeId(C.mlirSymbolRefAttrGetTypeID())
}

func (attr SymbolRefAttr) RootReference() string {
	return BorrowedStringRef(C.mlirSymbolRefAttrGetRootReference(attr.Raw())).String()
}

func (attr SymbolRefAttr) LeafReference() string {
	return BorrowedStringRef(C.mlirSymbolRefAttrGetLeafReference(attr.Raw())).String()
}

func (attr SymbolRefAttr) NumNestedReferences() int {
	return int(C.mlirSymbolRefAttrGetNumNestedReferences(attr.Raw()))
}

func (attr SymbolRefAttr) NestedReference(pos int) Attribute {
	return WrapAttribute(C.mlirSymbolRefAttrGetNestedReference(attr.Raw(), C.intptr_t(pos)))
}

//===----------------------------------------------------------------------===//
// Flat SymbolRef attribute.
//===----------------------------------------------------------------------===//

type FlatSymbolRefAttr struct {
	baseAttribute
}

func NewFlatSymbolRefAttr(ctx Context, symbol string) FlatSymbolRefAttr {
	cSymbol := NewStringRef(symbol)
	defer cSymbol.Destroy()
	return FlatSymbolRefAttr{
		baseAttribute: baseAttribute(C.mlirFlatSymbolRefAttrGet(ctx.Raw(), cSymbol.Raw())),
	}
}

func FlatSymbolRefAttrName() string {
	return BorrowedStringRef(C.mlirFlatSymbolRefAttrGetName()).String()
}

func (attr FlatSymbolRefAttr) Value() string {
	return BorrowedStringRef(C.mlirFlatSymbolRefAttrGetValue(attr.Raw())).String()
}

//===----------------------------------------------------------------------===//
// Type attribute.
//===----------------------------------------------------------------------===//

type TypeAttr struct {
	baseAttribute
}

func NewTypeAttr(type_ TypeLike) TypeAttr {
	return TypeAttr{
		baseAttribute: baseAttribute(C.mlirTypeAttrGet(type_.Raw())),
	}
}

func TypeAttrName() string {
	return BorrowedStringRef(C.mlirTypeAttrGetName()).String()
}

func TypeAttrTypeId() TypeId {
	return TypeId(C.mlirTypeAttrGetTypeID())
}

func (attr TypeAttr) Value() Type {
	return WrapType(C.mlirTypeAttrGetValue(attr.Raw()))
}

//===----------------------------------------------------------------------===//
// Unit attribute.
//===----------------------------------------------------------------------===//

type UnitAttr struct {
	baseAttribute
}

func NewUnitAttr(ctx Context) UnitAttr {
	return UnitAttr{
		baseAttribute: baseAttribute(C.mlirUnitAttrGet(ctx.Raw())),
	}
}

func UnitAttrName() string {
	return BorrowedStringRef(C.mlirUnitAttrGetName()).String()
}

func UnitAttrTypeId() TypeId {
	return TypeId(C.mlirUnitAttrGetTypeID())
}

//===----------------------------------------------------------------------===//
// Elements attributes.
//===----------------------------------------------------------------------===//

type ElementsAttr struct {
	baseAttribute
}

func (attr ElementsAttr) Value(rank int, idxs []uint64) Attribute {
	var cIdxs *C.uint64_t
	if len(idxs) > 0 {
		cIdxs = (*C.uint64_t)(unsafe.Pointer(&idxs[0]))
	}
	return WrapAttribute(C.mlirElementsAttrGetValue(attr.Raw(), C.intptr_t(rank), cIdxs))
}

func (attr ElementsAttr) IsValidIndex(rank int, idxs []uint64) bool {
	var cIdxs *C.uint64_t
	if len(idxs) > 0 {
		cIdxs = (*C.uint64_t)(unsafe.Pointer(&idxs[0]))
	}
	return bool(C.mlirElementsAttrIsValidIndex(attr.Raw(), C.intptr_t(rank), cIdxs))
}

func (attr ElementsAttr) NumElements() int64 {
	return int64(C.mlirElementsAttrGetNumElements(attr.Raw()))
}

//===----------------------------------------------------------------------===//
// Dense array attribute.
//===----------------------------------------------------------------------===//

type DenseArrayAttr struct {
	baseAttribute
}

func DenseArrayAttrTypeId() TypeId {
	return TypeId(C.mlirDenseArrayAttrGetTypeID())
}

func NewDenseBoolArrayAttr(ctx Context, values []bool) DenseArrayAttr {
	var cValues *C.int
	if len(values) > 0 {
		// Convert []bool to []int (0 or 1)
		intValues := make([]C.int, len(values))
		for i, v := range values {
			if v {
				intValues[i] = 1
			} else {
				intValues[i] = 0
			}
		}
		cValues = &intValues[0]
	}
	return DenseArrayAttr{
		baseAttribute: baseAttribute(C.mlirDenseBoolArrayGet(ctx.Raw(), C.intptr_t(len(values)), cValues)),
	}
}

func NewDenseI8ArrayAttr(ctx Context, values []int8) DenseArrayAttr {
	var cValues *C.int8_t
	if len(values) > 0 {
		cValues = (*C.int8_t)(unsafe.Pointer(&values[0]))
	}
	return DenseArrayAttr{
		baseAttribute: baseAttribute(C.mlirDenseI8ArrayGet(ctx.Raw(), C.intptr_t(len(values)), cValues)),
	}
}

func NewDenseI16ArrayAttr(ctx Context, values []int16) DenseArrayAttr {
	var cValues *C.int16_t
	if len(values) > 0 {
		cValues = (*C.int16_t)(unsafe.Pointer(&values[0]))
	}
	return DenseArrayAttr{
		baseAttribute: baseAttribute(C.mlirDenseI16ArrayGet(ctx.Raw(), C.intptr_t(len(values)), cValues)),
	}
}

func NewDenseI32ArrayAttr(ctx Context, values []int32) DenseArrayAttr {
	var cValues *C.int32_t
	if len(values) > 0 {
		cValues = (*C.int32_t)(unsafe.Pointer(&values[0]))
	}
	return DenseArrayAttr{
		baseAttribute: baseAttribute(C.mlirDenseI32ArrayGet(ctx.Raw(), C.intptr_t(len(values)), cValues)),
	}
}

func NewDenseI64ArrayAttr(ctx Context, values []int64) DenseArrayAttr {
	var cValues *C.int64_t
	if len(values) > 0 {
		cValues = (*C.int64_t)(unsafe.Pointer(&values[0]))
	}
	return DenseArrayAttr{
		baseAttribute: baseAttribute(C.mlirDenseI64ArrayGet(ctx.Raw(), C.intptr_t(len(values)), cValues)),
	}
}

func NewDenseF32ArrayAttr(ctx Context, values []float32) DenseArrayAttr {
	var cValues *C.float
	if len(values) > 0 {
		cValues = (*C.float)(unsafe.Pointer(&values[0]))
	}
	return DenseArrayAttr{
		baseAttribute: baseAttribute(C.mlirDenseF32ArrayGet(ctx.Raw(), C.intptr_t(len(values)), cValues)),
	}
}

func NewDenseF64ArrayAttr(ctx Context, values []float64) DenseArrayAttr {
	var cValues *C.double
	if len(values) > 0 {
		cValues = (*C.double)(unsafe.Pointer(&values[0]))
	}
	return DenseArrayAttr{
		baseAttribute: baseAttribute(C.mlirDenseF64ArrayGet(ctx.Raw(), C.intptr_t(len(values)), cValues)),
	}
}

func (attr DenseArrayAttr) NumElements() int {
	return int(C.mlirDenseArrayGetNumElements(attr.Raw()))
}

func (attr DenseArrayAttr) ElementBool(pos int) bool {
	return bool(C.mlirDenseBoolArrayGetElement(attr.Raw(), C.intptr_t(pos)))
}

func (attr DenseArrayAttr) ElementI8(pos int) int8 {
	return int8(C.mlirDenseI8ArrayGetElement(attr.Raw(), C.intptr_t(pos)))
}

func (attr DenseArrayAttr) ElementI16(pos int) int16 {
	return int16(C.mlirDenseI16ArrayGetElement(attr.Raw(), C.intptr_t(pos)))
}

func (attr DenseArrayAttr) ElementI32(pos int) int32 {
	return int32(C.mlirDenseI32ArrayGetElement(attr.Raw(), C.intptr_t(pos)))
}

func (attr DenseArrayAttr) ElementI64(pos int) int64 {
	return int64(C.mlirDenseI64ArrayGetElement(attr.Raw(), C.intptr_t(pos)))
}

func (attr DenseArrayAttr) ElementF32(pos int) float32 {
	return float32(C.mlirDenseF32ArrayGetElement(attr.Raw(), C.intptr_t(pos)))
}

func (attr DenseArrayAttr) ElementF64(pos int) float64 {
	return float64(C.mlirDenseF64ArrayGetElement(attr.Raw(), C.intptr_t(pos)))
}

//===----------------------------------------------------------------------===//
// Dense elements attribute.
//===----------------------------------------------------------------------===//

type DenseElementsAttr struct {
	baseAttribute
}

func DenseIntOrFPElementsAttrTypeId() TypeId {
	return TypeId(C.mlirDenseIntOrFPElementsAttrGetTypeID())
}

func NewDenseElementsAttr(shapedType TypeLike, elements []Attribute) DenseElementsAttr {
	var cElements *C.MlirAttribute
	if len(elements) > 0 {
		cElements = &UnwrapAttributeSlice(elements)[0]
	}
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrGet(shapedType.Raw(), C.intptr_t(len(elements)), cElements)),
	}
}

func NewDenseElementsAttrRawBuffer(shapedType TypeLike, rawBuffer []byte) DenseElementsAttr {
	var cRawBuffer unsafe.Pointer
	if len(rawBuffer) > 0 {
		cRawBuffer = unsafe.Pointer(&rawBuffer[0])
	}
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrRawBufferGet(shapedType.Raw(), C.size_t(len(rawBuffer)), cRawBuffer)),
	}
}

func NewDenseElementsAttrSplat(shapedType TypeLike, element Attribute) DenseElementsAttr {
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrSplatGet(shapedType.Raw(), element.Raw())),
	}
}

func NewDenseElementsAttrBoolSplat(shapedType TypeLike, element bool) DenseElementsAttr {
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrBoolSplatGet(shapedType.Raw(), C.bool(element))),
	}
}

func NewDenseElementsAttrUInt8Splat(shapedType TypeLike, element uint8) DenseElementsAttr {
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrUInt8SplatGet(shapedType.Raw(), C.uint8_t(element))),
	}
}

func NewDenseElementsAttrInt8Splat(shapedType TypeLike, element int8) DenseElementsAttr {
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrInt8SplatGet(shapedType.Raw(), C.int8_t(element))),
	}
}

func NewDenseElementsAttrUInt32Splat(shapedType TypeLike, element uint32) DenseElementsAttr {
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrUInt32SplatGet(shapedType.Raw(), C.uint32_t(element))),
	}
}

func NewDenseElementsAttrInt32Splat(shapedType TypeLike, element int32) DenseElementsAttr {
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrInt32SplatGet(shapedType.Raw(), C.int32_t(element))),
	}
}

func NewDenseElementsAttrUInt64Splat(shapedType TypeLike, element uint64) DenseElementsAttr {
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrUInt64SplatGet(shapedType.Raw(), C.uint64_t(element))),
	}
}

func NewDenseElementsAttrInt64Splat(shapedType TypeLike, element int64) DenseElementsAttr {
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrInt64SplatGet(shapedType.Raw(), C.int64_t(element))),
	}
}

func NewDenseElementsAttrFloatSplat(shapedType TypeLike, element float32) DenseElementsAttr {
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrFloatSplatGet(shapedType.Raw(), C.float(element))),
	}
}

func NewDenseElementsAttrDoubleSplat(shapedType TypeLike, element float64) DenseElementsAttr {
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrDoubleSplatGet(shapedType.Raw(), C.double(element))),
	}
}

func NewDenseElementsAttrBool(shapedType TypeLike, elements []bool) DenseElementsAttr {
	var cElements *C.int
	if len(elements) > 0 {
		intElements := make([]C.int, len(elements))
		for i, v := range elements {
			if v {
				intElements[i] = 1
			} else {
				intElements[i] = 0
			}
		}
		cElements = &intElements[0]
	}
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrBoolGet(shapedType.Raw(), C.intptr_t(len(elements)), cElements)),
	}
}

func NewDenseElementsAttrUInt8(shapedType TypeLike, elements []uint8) DenseElementsAttr {
	var cElements *C.uint8_t
	if len(elements) > 0 {
		cElements = (*C.uint8_t)(unsafe.Pointer(&elements[0]))
	}
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrUInt8Get(shapedType.Raw(), C.intptr_t(len(elements)), cElements)),
	}
}

func NewDenseElementsAttrInt8(shapedType TypeLike, elements []int8) DenseElementsAttr {
	var cElements *C.int8_t
	if len(elements) > 0 {
		cElements = (*C.int8_t)(unsafe.Pointer(&elements[0]))
	}
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrInt8Get(shapedType.Raw(), C.intptr_t(len(elements)), cElements)),
	}
}

func NewDenseElementsAttrUInt16(shapedType TypeLike, elements []uint16) DenseElementsAttr {
	var cElements *C.uint16_t
	if len(elements) > 0 {
		cElements = (*C.uint16_t)(unsafe.Pointer(&elements[0]))
	}
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrUInt16Get(shapedType.Raw(), C.intptr_t(len(elements)), cElements)),
	}
}

func NewDenseElementsAttrInt16(shapedType TypeLike, elements []int16) DenseElementsAttr {
	var cElements *C.int16_t
	if len(elements) > 0 {
		cElements = (*C.int16_t)(unsafe.Pointer(&elements[0]))
	}
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrInt16Get(shapedType.Raw(), C.intptr_t(len(elements)), cElements)),
	}
}

func NewDenseElementsAttrUInt32(shapedType TypeLike, elements []uint32) DenseElementsAttr {
	var cElements *C.uint32_t
	if len(elements) > 0 {
		cElements = (*C.uint32_t)(unsafe.Pointer(&elements[0]))
	}
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrUInt32Get(shapedType.Raw(), C.intptr_t(len(elements)), cElements)),
	}
}

func NewDenseElementsAttrInt32(shapedType TypeLike, elements []int32) DenseElementsAttr {
	var cElements *C.int32_t
	if len(elements) > 0 {
		cElements = (*C.int32_t)(unsafe.Pointer(&elements[0]))
	}
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrInt32Get(shapedType.Raw(), C.intptr_t(len(elements)), cElements)),
	}
}

func NewDenseElementsAttrUInt64(shapedType TypeLike, elements []uint64) DenseElementsAttr {
	var cElements *C.uint64_t
	if len(elements) > 0 {
		cElements = (*C.uint64_t)(unsafe.Pointer(&elements[0]))
	}
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrUInt64Get(shapedType.Raw(), C.intptr_t(len(elements)), cElements)),
	}
}

func NewDenseElementsAttrInt64(shapedType TypeLike, elements []int64) DenseElementsAttr {
	var cElements *C.int64_t
	if len(elements) > 0 {
		cElements = (*C.int64_t)(unsafe.Pointer(&elements[0]))
	}
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrInt64Get(shapedType.Raw(), C.intptr_t(len(elements)), cElements)),
	}
}

func NewDenseElementsAttrFloat(shapedType TypeLike, elements []float32) DenseElementsAttr {
	var cElements *C.float
	if len(elements) > 0 {
		cElements = (*C.float)(unsafe.Pointer(&elements[0]))
	}
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrFloatGet(shapedType.Raw(), C.intptr_t(len(elements)), cElements)),
	}
}

func NewDenseElementsAttrDouble(shapedType TypeLike, elements []float64) DenseElementsAttr {
	var cElements *C.double
	if len(elements) > 0 {
		cElements = (*C.double)(unsafe.Pointer(&elements[0]))
	}
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrDoubleGet(shapedType.Raw(), C.intptr_t(len(elements)), cElements)),
	}
}

func NewDenseElementsAttrBFloat16(shapedType TypeLike, elements []uint16) DenseElementsAttr {
	var cElements *C.uint16_t
	if len(elements) > 0 {
		cElements = (*C.uint16_t)(unsafe.Pointer(&elements[0]))
	}
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrBFloat16Get(shapedType.Raw(), C.intptr_t(len(elements)), cElements)),
	}
}

func NewDenseElementsAttrFloat16(shapedType TypeLike, elements []uint16) DenseElementsAttr {
	var cElements *C.uint16_t
	if len(elements) > 0 {
		cElements = (*C.uint16_t)(unsafe.Pointer(&elements[0]))
	}
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrFloat16Get(shapedType.Raw(), C.intptr_t(len(elements)), cElements)),
	}
}

func NewDenseElementsAttrString(shapedType TypeLike, elements []string) DenseElementsAttr {
	var cElements []C.MlirStringRef
	var refs []StringRef
	if len(elements) > 0 {
		cElements = make([]C.MlirStringRef, len(elements))
		refs = make([]StringRef, len(elements))
		for i, s := range elements {
			refs[i] = NewStringRef(s)
			cElements[i] = refs[i].Raw()
		}
	}
	defer func() {
		for _, ref := range refs {
			ref.Destroy()
		}
	}()

	var ptr *C.MlirStringRef
	if len(elements) > 0 {
		ptr = &cElements[0]
	}

	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrStringGet(shapedType.Raw(), C.intptr_t(len(elements)), ptr)),
	}
}

func NewDenseElementsAttrReshape(attr Attribute, shapedType TypeLike) DenseElementsAttr {
	return DenseElementsAttr{
		baseAttribute: baseAttribute(C.mlirDenseElementsAttrReshapeGet(attr.Raw(), shapedType.Raw())),
	}
}

func (attr DenseElementsAttr) IsSplat() bool {
	return bool(C.mlirDenseElementsAttrIsSplat(attr.Raw()))
}

func (attr DenseElementsAttr) SplatValue() Attribute {
	return WrapAttribute(C.mlirDenseElementsAttrGetSplatValue(attr.Raw()))
}

func (attr DenseElementsAttr) BoolSplatValue() bool {
	return int(C.mlirDenseElementsAttrGetBoolSplatValue(attr.Raw())) != 0
}

func (attr DenseElementsAttr) Int8SplatValue() int8 {
	return int8(C.mlirDenseElementsAttrGetInt8SplatValue(attr.Raw()))
}

func (attr DenseElementsAttr) UInt8SplatValue() uint8 {
	return uint8(C.mlirDenseElementsAttrGetUInt8SplatValue(attr.Raw()))
}

func (attr DenseElementsAttr) Int32SplatValue() int32 {
	return int32(C.mlirDenseElementsAttrGetInt32SplatValue(attr.Raw()))
}

func (attr DenseElementsAttr) UInt32SplatValue() uint32 {
	return uint32(C.mlirDenseElementsAttrGetUInt32SplatValue(attr.Raw()))
}

func (attr DenseElementsAttr) Int64SplatValue() int64 {
	return int64(C.mlirDenseElementsAttrGetInt64SplatValue(attr.Raw()))
}

func (attr DenseElementsAttr) UInt64SplatValue() uint64 {
	return uint64(C.mlirDenseElementsAttrGetUInt64SplatValue(attr.Raw()))
}

func (attr DenseElementsAttr) FloatSplatValue() float32 {
	return float32(C.mlirDenseElementsAttrGetFloatSplatValue(attr.Raw()))
}

func (attr DenseElementsAttr) DoubleSplatValue() float64 {
	return float64(C.mlirDenseElementsAttrGetDoubleSplatValue(attr.Raw()))
}

func (attr DenseElementsAttr) StringSplatValue() string {
	return BorrowedStringRef(C.mlirDenseElementsAttrGetStringSplatValue(attr.Raw())).String()
}

func (attr DenseElementsAttr) BoolValue(pos int) bool {
	return bool(C.mlirDenseElementsAttrGetBoolValue(attr.Raw(), C.intptr_t(pos)))
}

func (attr DenseElementsAttr) Int8Value(pos int) int8 {
	return int8(C.mlirDenseElementsAttrGetInt8Value(attr.Raw(), C.intptr_t(pos)))
}

func (attr DenseElementsAttr) UInt8Value(pos int) uint8 {
	return uint8(C.mlirDenseElementsAttrGetUInt8Value(attr.Raw(), C.intptr_t(pos)))
}

func (attr DenseElementsAttr) Int16Value(pos int) int16 {
	return int16(C.mlirDenseElementsAttrGetInt16Value(attr.Raw(), C.intptr_t(pos)))
}

func (attr DenseElementsAttr) UInt16Value(pos int) uint16 {
	return uint16(C.mlirDenseElementsAttrGetUInt16Value(attr.Raw(), C.intptr_t(pos)))
}

func (attr DenseElementsAttr) Int32Value(pos int) int32 {
	return int32(C.mlirDenseElementsAttrGetInt32Value(attr.Raw(), C.intptr_t(pos)))
}

func (attr DenseElementsAttr) UInt32Value(pos int) uint32 {
	return uint32(C.mlirDenseElementsAttrGetUInt32Value(attr.Raw(), C.intptr_t(pos)))
}

func (attr DenseElementsAttr) Int64Value(pos int) int64 {
	return int64(C.mlirDenseElementsAttrGetInt64Value(attr.Raw(), C.intptr_t(pos)))
}

func (attr DenseElementsAttr) UInt64Value(pos int) uint64 {
	return uint64(C.mlirDenseElementsAttrGetUInt64Value(attr.Raw(), C.intptr_t(pos)))
}

func (attr DenseElementsAttr) FloatValue(pos int) float32 {
	return float32(C.mlirDenseElementsAttrGetFloatValue(attr.Raw(), C.intptr_t(pos)))
}

func (attr DenseElementsAttr) DoubleValue(pos int) float64 {
	return float64(C.mlirDenseElementsAttrGetDoubleValue(attr.Raw(), C.intptr_t(pos)))
}

func (attr DenseElementsAttr) StringValue(pos int) string {
	return BorrowedStringRef(C.mlirDenseElementsAttrGetStringValue(attr.Raw(), C.intptr_t(pos))).String()
}

func (attr DenseElementsAttr) RawData() unsafe.Pointer {
	return C.mlirDenseElementsAttrGetRawData(attr.Raw())
}

//===----------------------------------------------------------------------===//
// Sparse elements attribute.
//===----------------------------------------------------------------------===//

type SparseElementsAttr struct {
	baseAttribute
}

func NewSparseElementsAttr(shapedType TypeLike, denseIndices Attribute, denseValues Attribute) SparseElementsAttr {
	return SparseElementsAttr{
		baseAttribute: baseAttribute(C.mlirSparseElementsAttribute(shapedType.Raw(), denseIndices.Raw(), denseValues.Raw())),
	}
}

func SparseElementsAttrTypeId() TypeId {
	return TypeId(C.mlirSparseElementsAttrGetTypeID())
}

func (attr SparseElementsAttr) Indices() Attribute {
	return WrapAttribute(C.mlirSparseElementsAttrGetIndices(attr.Raw()))
}

func (attr SparseElementsAttr) Values() Attribute {
	return WrapAttribute(C.mlirSparseElementsAttrGetValues(attr.Raw()))
}

//===----------------------------------------------------------------------===//
// Strided layout attribute.
//===----------------------------------------------------------------------===//

type StridedLayoutAttr struct {
	baseAttribute
}

func NewStridedLayoutAttr(ctx Context, offset int64, strides []int64) StridedLayoutAttr {
	var cStrides *C.int64_t
	if len(strides) > 0 {
		cStrides = (*C.int64_t)(unsafe.Pointer(&strides[0]))
	}
	return StridedLayoutAttr{
		baseAttribute: baseAttribute(C.mlirStridedLayoutAttrGet(ctx.Raw(), C.int64_t(offset), C.intptr_t(len(strides)), cStrides)),
	}
}

func StridedLayoutAttrName() string {
	return BorrowedStringRef(C.mlirStridedLayoutAttrGetName()).String()
}

func StridedLayoutAttrTypeId() TypeId {
	return TypeId(C.mlirStridedLayoutAttrGetTypeID())
}

func (attr StridedLayoutAttr) Offset() int64 {
	return int64(C.mlirStridedLayoutAttrGetOffset(attr.Raw()))
}

func (attr StridedLayoutAttr) NumStrides() int {
	return int(C.mlirStridedLayoutAttrGetNumStrides(attr.Raw()))
}

func (attr StridedLayoutAttr) Stride(pos int) int64 {
	return int64(C.mlirStridedLayoutAttrGetStride(attr.Raw(), C.intptr_t(pos)))
}
