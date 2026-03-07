package mlir

/*
#include <mlir-c/BuiltinTypes.h>
*/
import "C"
import "unsafe"

//===----------------------------------------------------------------------===//
// Integer types.
//===----------------------------------------------------------------------===//

type IntegerType struct {
	baseType
}

func NewIntegerType(ctx Context, width int) IntegerType {
	return IntegerType{
		baseType: baseType(C.mlirIntegerTypeGet(ctx.raw(), C.unsigned(width))),
	}
}

func NewSignedIntegerType(ctx Context, width int) IntegerType {
	return IntegerType{
		baseType: baseType(C.mlirIntegerTypeSignedGet(ctx.raw(), C.unsigned(width))),
	}
}

func NewUnsignedIntegerType(ctx Context, width int) IntegerType {
	return IntegerType{
		baseType: baseType(C.mlirIntegerTypeUnsignedGet(ctx.raw(), C.unsigned(width))),
	}
}

func IntegerTypeId() TypeId { return TypeId(C.mlirIntegerTypeGetTypeID()) }

func AsIntegerType(ty TypeLike) (IntegerType, bool) {
	if C.mlirTypeIsAInteger(ty.raw()) {
		return IntegerType{baseType: baseType(ty.raw())}, true
	}
	return IntegerType{}, false
}

func (i IntegerType) Width() int       { return int(C.mlirIntegerTypeGetWidth(i.raw())) }
func (i IntegerType) IsSignless() bool { return bool(C.mlirIntegerTypeIsSignless(i.raw())) }
func (i IntegerType) IsSigned() bool   { return bool(C.mlirIntegerTypeIsSigned(i.raw())) }
func (i IntegerType) IsUnsigned() bool { return bool(C.mlirIntegerTypeIsUnsigned(i.raw())) }

//===----------------------------------------------------------------------===//
// Index type.
//===----------------------------------------------------------------------===//

type IndexType struct {
	baseType
}

func NewIndexType(ctx Context) IndexType {
	return IndexType{baseType: baseType(C.mlirIndexTypeGet(ctx.raw()))}
}

func AsIndexType(ty TypeLike) (IndexType, bool) {
	if C.mlirTypeIsAIndex(ty.raw()) {
		return IndexType{baseType: baseType(ty.raw())}, true
	}
	return IndexType{}, false
}

func IndexTypeId() TypeId { return TypeId(C.mlirIndexTypeGetTypeID()) }
func IndexTypeName() string {
	ref := BorrowedStringRef(C.mlirIndexTypeGetName())
	return ref.String()
}

//===----------------------------------------------------------------------===//
// Floating-point types.
//===----------------------------------------------------------------------===//

type FloatType struct {
	baseType
}

type FloatFormat int

const (
	FloatUnknown FloatFormat = iota
	Float4E2M1FN
	Float6E2M3FN
	Float6E3M2FN
	Float8E5M2
	Float8E4M3
	Float8E4M3FN
	Float8E5M2FNUZ
	Float8E4M3FNUZ
	Float8E4M3B11FNUZ
	Float8E3M4
	Float8E8M0FNU
	BFloat16
	Float16
	Float32
	Float64
	FloatTF32
)

type floatInfo struct {
	get    func(C.MlirContext) C.MlirType
	typeID func() C.MlirTypeID
	name   func() string
}

var floatFormats = map[FloatFormat]floatInfo{
	Float4E2M1FN: {
		get:    func(c C.MlirContext) C.MlirType { return C.mlirFloat4E2M1FNTypeGet(c) },
		typeID: func() C.MlirTypeID { return C.mlirFloat4E2M1FNTypeGetTypeID() },
		name:   func() string { return BorrowedStringRef(C.mlirFloat4E2M1FNTypeGetName()).String() },
	},
	Float6E2M3FN: {
		get:    func(c C.MlirContext) C.MlirType { return C.mlirFloat6E2M3FNTypeGet(c) },
		typeID: func() C.MlirTypeID { return C.mlirFloat6E2M3FNTypeGetTypeID() },
		name:   func() string { return BorrowedStringRef(C.mlirFloat6E2M3FNTypeGetName()).String() },
	},
	Float6E3M2FN: {
		get:    func(c C.MlirContext) C.MlirType { return C.mlirFloat6E3M2FNTypeGet(c) },
		typeID: func() C.MlirTypeID { return C.mlirFloat6E3M2FNTypeGetTypeID() },
		name:   func() string { return BorrowedStringRef(C.mlirFloat6E3M2FNTypeGetName()).String() },
	},
	Float8E5M2: {
		get:    func(c C.MlirContext) C.MlirType { return C.mlirFloat8E5M2TypeGet(c) },
		typeID: func() C.MlirTypeID { return C.mlirFloat8E5M2TypeGetTypeID() },
		name:   func() string { return BorrowedStringRef(C.mlirFloat8E5M2TypeGetName()).String() },
	},
	Float8E4M3: {
		get:    func(c C.MlirContext) C.MlirType { return C.mlirFloat8E4M3TypeGet(c) },
		typeID: func() C.MlirTypeID { return C.mlirFloat8E4M3TypeGetTypeID() },
		name:   func() string { return BorrowedStringRef(C.mlirFloat8E4M3TypeGetName()).String() },
	},
	Float8E4M3FN: {
		get:    func(c C.MlirContext) C.MlirType { return C.mlirFloat8E4M3FNTypeGet(c) },
		typeID: func() C.MlirTypeID { return C.mlirFloat8E4M3FNTypeGetTypeID() },
		name:   func() string { return BorrowedStringRef(C.mlirFloat8E4M3FNTypeGetName()).String() },
	},
	Float8E5M2FNUZ: {
		get:    func(c C.MlirContext) C.MlirType { return C.mlirFloat8E5M2FNUZTypeGet(c) },
		typeID: func() C.MlirTypeID { return C.mlirFloat8E5M2FNUZTypeGetTypeID() },
		name:   func() string { return BorrowedStringRef(C.mlirFloat8E5M2FNUZTypeGetName()).String() },
	},
	Float8E4M3FNUZ: {
		get:    func(c C.MlirContext) C.MlirType { return C.mlirFloat8E4M3FNUZTypeGet(c) },
		typeID: func() C.MlirTypeID { return C.mlirFloat8E4M3FNUZTypeGetTypeID() },
		name:   func() string { return BorrowedStringRef(C.mlirFloat8E4M3FNUZTypeGetName()).String() },
	},
	Float8E4M3B11FNUZ: {
		get:    func(c C.MlirContext) C.MlirType { return C.mlirFloat8E4M3B11FNUZTypeGet(c) },
		typeID: func() C.MlirTypeID { return C.mlirFloat8E4M3B11FNUZTypeGetTypeID() },
		name:   func() string { return BorrowedStringRef(C.mlirFloat8E4M3B11FNUZTypeGetName()).String() },
	},
	Float8E3M4: {
		get:    func(c C.MlirContext) C.MlirType { return C.mlirFloat8E3M4TypeGet(c) },
		typeID: func() C.MlirTypeID { return C.mlirFloat8E3M4TypeGetTypeID() },
		name:   func() string { return BorrowedStringRef(C.mlirFloat8E3M4TypeGetName()).String() },
	},
	Float8E8M0FNU: {
		get:    func(c C.MlirContext) C.MlirType { return C.mlirFloat8E8M0FNUTypeGet(c) },
		typeID: func() C.MlirTypeID { return C.mlirFloat8E8M0FNUTypeGetTypeID() },
		name:   func() string { return BorrowedStringRef(C.mlirFloat8E8M0FNUTypeGetName()).String() },
	},
	BFloat16: {
		get:    func(c C.MlirContext) C.MlirType { return C.mlirBF16TypeGet(c) },
		typeID: func() C.MlirTypeID { return C.mlirBFloat16TypeGetTypeID() },
		name:   func() string { return BorrowedStringRef(C.mlirBF16TypeGetName()).String() },
	},
	Float16: {
		get:    func(c C.MlirContext) C.MlirType { return C.mlirF16TypeGet(c) },
		typeID: func() C.MlirTypeID { return C.mlirFloat16TypeGetTypeID() },
		name:   func() string { return BorrowedStringRef(C.mlirF16TypeGetName()).String() },
	},
	Float32: {
		get:    func(c C.MlirContext) C.MlirType { return C.mlirF32TypeGet(c) },
		typeID: func() C.MlirTypeID { return C.mlirFloat32TypeGetTypeID() },
		name:   func() string { return BorrowedStringRef(C.mlirF32TypeGetName()).String() },
	},
	Float64: {
		get:    func(c C.MlirContext) C.MlirType { return C.mlirF64TypeGet(c) },
		typeID: func() C.MlirTypeID { return C.mlirFloat64TypeGetTypeID() },
		name:   func() string { return BorrowedStringRef(C.mlirF64TypeGetName()).String() },
	},
	FloatTF32: {
		get:    func(c C.MlirContext) C.MlirType { return C.mlirTF32TypeGet(c) },
		typeID: func() C.MlirTypeID { return C.mlirFloatTF32TypeGetTypeID() },
		name:   func() string { return "tf32" },
	},
}

func NewFloatType(ctx Context, format FloatFormat) FloatType {
	if info, ok := floatFormats[format]; ok {
		return FloatType{baseType: baseType(info.get(ctx.raw()))}
	}
	panic("unsupported float format")
}

func AsFloatType(ty TypeLike) (FloatType, bool) {
	if C.mlirTypeIsAFloat(ty.raw()) {
		return FloatType{baseType: baseType(ty.raw())}, true
	}
	return FloatType{}, false
}

func FloatTypeId(format FloatFormat) TypeId {
	if info, ok := floatFormats[format]; ok {
		return TypeId(info.typeID())
	}
	panic("unsupported float format")
}

func FloatTypeName(format FloatFormat) string {
	if info, ok := floatFormats[format]; ok {
		return info.name()
	}
	panic("unsupported float format")
}

func (f FloatType) Format() FloatFormat {
	typeID := C.mlirTypeGetTypeID(f.raw())
	for format, info := range floatFormats {
		if C.mlirTypeIDEqual(typeID, info.typeID()) {
			return format
		}
	}
	return FloatUnknown
}

//===----------------------------------------------------------------------===//
// None type.
//===----------------------------------------------------------------------===//

type NoneType struct {
	baseType
}

func NewNoneType(ctx Context) NoneType {
	return NoneType{baseType: baseType(C.mlirNoneTypeGet(ctx.raw()))}
}

func AsNoneType(ty TypeLike) (NoneType, bool) {
	if C.mlirTypeIsANone(ty.raw()) {
		return NoneType{baseType: baseType(ty.raw())}, true
	}
	return NoneType{}, false
}

func NoneTypeId() TypeId { return TypeId(C.mlirNoneTypeGetTypeID()) }
func NoneTypeName() string {
	ref := BorrowedStringRef(C.mlirNoneTypeGetName())
	return ref.String()
}

//===----------------------------------------------------------------------===//
// Complex type.
//===----------------------------------------------------------------------===//

type ComplexType struct {
	baseType
}

func NewComplexType(elementType TypeLike) ComplexType {
	return ComplexType{baseType: baseType(C.mlirComplexTypeGet(elementType.raw()))}
}

func AsComplexType(ty TypeLike) (ComplexType, bool) {
	if C.mlirTypeIsAComplex(ty.raw()) {
		return ComplexType{baseType: baseType(ty.raw())}, true
	}
	return ComplexType{}, false
}

func ComplexTypeId() TypeId { return TypeId(C.mlirComplexTypeGetTypeID()) }
func ComplexTypeName() string {
	ref := BorrowedStringRef(C.mlirComplexTypeGetName())
	return ref.String()
}

func (c ComplexType) ElementType() Type {
	return Type{baseType: baseType(C.mlirComplexTypeGetElementType(c.raw()))}
}

//===----------------------------------------------------------------------===//
// Shaped type.
//===----------------------------------------------------------------------===//

type ShapedType struct {
	baseType
}

func AsShapedType(ty TypeLike) (ShapedType, bool) {
	if C.mlirTypeIsAShaped(ty.raw()) {
		return ShapedType{baseType: baseType(ty.raw())}, true
	}
	return ShapedType{}, false
}

func ShapedTypeDynamicSize() int { return int(C.mlirShapedTypeGetDynamicSize()) }
func IsDynamicStrideOfOffset(val int) bool {
	return bool(C.mlirShapedTypeIsDynamicStrideOrOffset(C.int64_t(val)))
}

func IsStaticStrideOfOffset(val int) bool {
	return bool(C.mlirShapedTypeIsStaticStrideOrOffset(C.int64_t(val)))
}

func (s ShapedType) ElementType() Type {
	return Type{baseType: baseType(C.mlirShapedTypeGetElementType(s.raw()))}
}

func (s ShapedType) HasRank() bool        { return bool(C.mlirShapedTypeHasRank(s.raw())) }
func (s ShapedType) Rank() int            { return int(C.mlirShapedTypeGetRank(s.raw())) }
func (s ShapedType) HasStaticShape() bool { return bool(C.mlirShapedTypeHasStaticShape(s.raw())) }
func (s ShapedType) IsDynamicDim(dim int) bool {
	return bool(C.mlirShapedTypeIsDynamicDim(s.raw(), C.intptr_t(dim)))
}

func (s ShapedType) IsStaticDim(dim int) bool {
	return bool(C.mlirShapedTypeIsStaticDim(s.raw(), C.intptr_t(dim)))
}

func (s ShapedType) DimSize(dim int) int {
	return int(C.mlirShapedTypeGetDimSize(s.raw(), C.intptr_t(dim)))
}

func (s ShapedType) IsDynamicSize(size int) bool {
	return bool(C.mlirShapedTypeIsDynamicSize(C.int64_t(size)))
}

func (s ShapedType) IsStaticSize(size int) bool {
	return bool(C.mlirShapedTypeIsStaticSize(C.int64_t(size)))
}

//===----------------------------------------------------------------------===//
// Vector type.
//===----------------------------------------------------------------------===//

type VectorType struct {
	baseType
}

func NewVectorType(shape []int64, elementType TypeLike) VectorType {
	var cShape *C.int64_t
	if len(shape) > 0 {
		cShape = (*C.int64_t)(unsafe.Pointer(&shape[0]))
	}
	return VectorType{
		baseType: baseType(C.mlirVectorTypeGet(C.intptr_t(len(shape)), cShape, elementType.raw())),
	}
}

func NewVectorTypeChecked(loc Location, shape []int64, elementType TypeLike) VectorType {
	var cShape *C.int64_t
	if len(shape) > 0 {
		cShape = (*C.int64_t)(unsafe.Pointer(&shape[0]))
	}
	return VectorType{
		baseType: baseType(C.mlirVectorTypeGetChecked(loc.raw(), C.intptr_t(len(shape)), cShape, elementType.raw())),
	}
}

func NewScalableVectorType(shape []int64, scalable []bool, elementType TypeLike) VectorType {
	var cShape *C.int64_t
	if len(shape) > 0 {
		cShape = (*C.int64_t)(unsafe.Pointer(&shape[0]))
	}
	var cScalable *C.bool
	if len(scalable) > 0 {
		cScalable = (*C.bool)(unsafe.Pointer(&scalable[0]))
	}
	return VectorType{
		baseType: baseType(C.mlirVectorTypeGetScalable(C.intptr_t(len(shape)), cShape, cScalable, elementType.raw())),
	}
}

func NewScalableVectorTypeChecked(loc Location, shape []int64, scalable []bool, elementType TypeLike) VectorType {
	var cShape *C.int64_t
	if len(shape) > 0 {
		cShape = (*C.int64_t)(unsafe.Pointer(&shape[0]))
	}
	var cScalable *C.bool
	if len(scalable) > 0 {
		cScalable = (*C.bool)(unsafe.Pointer(&scalable[0]))
	}
	return VectorType{
		baseType: baseType(C.mlirVectorTypeGetScalableChecked(loc.raw(), C.intptr_t(len(shape)), cShape, cScalable, elementType.raw())),
	}
}

func AsVectorType(ty TypeLike) (VectorType, bool) {
	if C.mlirTypeIsAVector(ty.raw()) {
		return VectorType{baseType: baseType(ty.raw())}, true
	}
	return VectorType{}, false
}

func VectorTypeId() TypeId { return TypeId(C.mlirVectorTypeGetTypeID()) }
func VectorTypeName() string {
	ref := BorrowedStringRef(C.mlirVectorTypeGetName())
	return ref.String()
}

func (v VectorType) IsScalable() bool { return bool(C.mlirVectorTypeIsScalable(v.raw())) }
func (v VectorType) IsDimScalable(dim int) bool {
	return bool(C.mlirVectorTypeIsDimScalable(v.raw(), C.intptr_t(dim)))
}

//===----------------------------------------------------------------------===//
// Ranked / Unranked Tensor type.
//===----------------------------------------------------------------------===//

type RankedTensorType struct {
	baseType
}

func NewRankedTensorType(shape []int64, elementType TypeLike, encoding Attribute) RankedTensorType {
	var cShape *C.int64_t
	if len(shape) > 0 {
		cShape = (*C.int64_t)(unsafe.Pointer(&shape[0]))
	}
	return RankedTensorType{
		baseType: baseType(C.mlirRankedTensorTypeGet(C.intptr_t(len(shape)), cShape, elementType.raw(), encoding.raw())),
	}
}

func NewRankedTensorTypeChecked(loc Location, shape []int64, elementType TypeLike, encoding Attribute) RankedTensorType {
	var cShape *C.int64_t
	if len(shape) > 0 {
		cShape = (*C.int64_t)(unsafe.Pointer(&shape[0]))
	}
	return RankedTensorType{
		baseType: baseType(C.mlirRankedTensorTypeGetChecked(loc.raw(), C.intptr_t(len(shape)), cShape, elementType.raw(), encoding.raw())),
	}
}

func AsRankedTensorType(ty TypeLike) (RankedTensorType, bool) {
	if C.mlirTypeIsARankedTensor(ty.raw()) {
		return RankedTensorType{baseType: baseType(ty.raw())}, true
	}
	return RankedTensorType{}, false
}

func RankedTensorTypeId() TypeId { return TypeId(C.mlirRankedTensorTypeGetTypeID()) }
func RankedTensorTypeName() string {
	ref := BorrowedStringRef(C.mlirRankedTensorTypeGetName())
	return ref.String()
}

func (t RankedTensorType) Encoding() Attribute {
	return Attribute(C.mlirRankedTensorTypeGetEncoding(t.raw()))
}

type UnrankedTensorType struct {
	baseType
}

func NewUnrankedTensorType(elementType TypeLike) UnrankedTensorType {
	return UnrankedTensorType{
		baseType: baseType(C.mlirUnrankedTensorTypeGet(elementType.raw())),
	}
}

func NewUnrankedTensorTypeChecked(loc Location, elementType TypeLike) UnrankedTensorType {
	return UnrankedTensorType{
		baseType: baseType(C.mlirUnrankedTensorTypeGetChecked(loc.raw(), elementType.raw())),
	}
}

func AsUnrankedTensorType(ty TypeLike) (UnrankedTensorType, bool) {
	if C.mlirTypeIsAUnrankedTensor(ty.raw()) {
		return UnrankedTensorType{baseType: baseType(ty.raw())}, true
	}
	return UnrankedTensorType{}, false
}

func UnrankedTensorTypeId() TypeId { return TypeId(C.mlirUnrankedTensorTypeGetTypeID()) }
func UnrankedTensorTypeName() string {
	ref := BorrowedStringRef(C.mlirUnrankedTensorTypeGetName())
	return ref.String()
}

//===----------------------------------------------------------------------===//
// Ranked / Unranked MemRef type.
//===----------------------------------------------------------------------===//

type MemRefType struct {
	baseType
}

func NewMemRefType(elementType TypeLike, shape []int64, layout Attribute, memorySpace Attribute) MemRefType {
	var cShape *C.int64_t
	if len(shape) > 0 {
		cShape = (*C.int64_t)(unsafe.Pointer(&shape[0]))
	}
	return MemRefType{
		baseType: baseType(C.mlirMemRefTypeGet(elementType.raw(), C.intptr_t(len(shape)), cShape, layout.raw(), memorySpace.raw())),
	}
}

func NewMemRefTypeChecked(loc Location, elementType TypeLike, shape []int64, layout Attribute, memorySpace Attribute) MemRefType {
	var cShape *C.int64_t
	if len(shape) > 0 {
		cShape = (*C.int64_t)(unsafe.Pointer(&shape[0]))
	}
	return MemRefType{
		baseType: baseType(C.mlirMemRefTypeGetChecked(loc.raw(), elementType.raw(), C.intptr_t(len(shape)), cShape, layout.raw(), memorySpace.raw())),
	}
}

func NewContiguousMemRefType(elementType TypeLike, shape []int64, memorySpace Attribute) MemRefType {
	var cShape *C.int64_t
	if len(shape) > 0 {
		cShape = (*C.int64_t)(unsafe.Pointer(&shape[0]))
	}
	return MemRefType{
		baseType: baseType(C.mlirMemRefTypeContiguousGet(elementType.raw(), C.intptr_t(len(shape)), cShape, memorySpace.raw())),
	}
}

func NewContiguousMemRefTypeChecked(loc Location, elementType TypeLike, shape []int64, memorySpace Attribute) MemRefType {
	var cShape *C.int64_t
	if len(shape) > 0 {
		cShape = (*C.int64_t)(unsafe.Pointer(&shape[0]))
	}
	return MemRefType{
		baseType: baseType(C.mlirMemRefTypeContiguousGetChecked(loc.raw(), elementType.raw(), C.intptr_t(len(shape)), cShape, memorySpace.raw())),
	}
}

func AsMemRefType(ty TypeLike) (MemRefType, bool) {
	if C.mlirTypeIsAMemRef(ty.raw()) {
		return MemRefType{baseType: baseType(ty.raw())}, true
	}
	return MemRefType{}, false
}

func MemRefTypeId() TypeId { return TypeId(C.mlirMemRefTypeGetTypeID()) }
func MemRefTypeName() string {
	ref := BorrowedStringRef(C.mlirMemRefTypeGetName())
	return ref.String()
}

func (m MemRefType) Layout() Attribute {
	return Attribute(C.mlirMemRefTypeGetLayout(m.raw()))
}

func (m MemRefType) AffineMap() AffineMap {
	return AffineMap(C.mlirMemRefTypeGetAffineMap(m.raw()))
}

func (m MemRefType) MemorySpace() Attribute {
	return Attribute(C.mlirMemRefTypeGetMemorySpace(m.raw()))
}

type UnrankedMemRefType struct {
	baseType
}

func NewUnrankedMemRefType(elementType TypeLike, memorySpace Attribute) UnrankedMemRefType {
	return UnrankedMemRefType{
		baseType: baseType(C.mlirUnrankedMemRefTypeGet(elementType.raw(), memorySpace.raw())),
	}
}

func NewUnrankedMemRefTypeChecked(loc Location, elementType TypeLike, memorySpace Attribute) UnrankedMemRefType {
	return UnrankedMemRefType{
		baseType: baseType(C.mlirUnrankedMemRefTypeGetChecked(loc.raw(), elementType.raw(), memorySpace.raw())),
	}
}

func AsUnrankedMemRefType(ty TypeLike) (UnrankedMemRefType, bool) {
	if C.mlirTypeIsAUnrankedMemRef(ty.raw()) {
		return UnrankedMemRefType{baseType: baseType(ty.raw())}, true
	}
	return UnrankedMemRefType{}, false
}

func UnrankedMemRefTypeId() TypeId { return TypeId(C.mlirUnrankedMemRefTypeGetTypeID()) }
func UnrankedMemRefTypeName() string {
	ref := BorrowedStringRef(C.mlirUnrankedMemRefTypeGetName())
	return ref.String()
}

func (m UnrankedMemRefType) MemorySpace() Attribute {
	return Attribute(C.mlirUnrankedMemrefGetMemorySpace(m.raw()))
}

//===----------------------------------------------------------------------===//
// Tuple type.
//===----------------------------------------------------------------------===//

type TupleType struct {
	baseType
}

func NewTupleType(ctx Context, elements []TypeLike) TupleType {
	var cElements *C.MlirType
	if len(elements) > 0 {
		cElements = &unwrapTypeSlice(elements)[0]
	}
	return TupleType{
		baseType: baseType(C.mlirTupleTypeGet(ctx.raw(), C.intptr_t(len(elements)), cElements)),
	}
}

func AsTupleType(ty TypeLike) (TupleType, bool) {
	if C.mlirTypeIsATuple(ty.raw()) {
		return TupleType{baseType: baseType(ty.raw())}, true
	}
	return TupleType{}, false
}

func TupleTypeId() TypeId { return TypeId(C.mlirTupleTypeGetTypeID()) }
func TupleTypeName() string {
	ref := BorrowedStringRef(C.mlirTupleTypeGetName())
	return ref.String()
}

func (t TupleType) NumTypes() int { return int(C.mlirTupleTypeGetNumTypes(t.raw())) }
func (t TupleType) Type(pos int) Type {
	return Type{baseType: baseType(C.mlirTupleTypeGetType(t.raw(), C.intptr_t(pos)))}
}

//===----------------------------------------------------------------------===//
// Function type.
//===----------------------------------------------------------------------===//

type FunctionType struct {
	baseType
}

func NewFunctionType(ctx Context, inputs []TypeLike, results []TypeLike) FunctionType {
	var cInputs *C.MlirType
	if len(inputs) > 0 {
		cInputs = &unwrapTypeSlice(inputs)[0]
	}
	var cResults *C.MlirType
	if len(results) > 0 {
		cResults = &unwrapTypeSlice(results)[0]
	}

	return FunctionType{
		baseType: baseType(C.mlirFunctionTypeGet(ctx.raw(), C.intptr_t(len(inputs)), cInputs, C.intptr_t(len(results)), cResults)),
	}
}

func AsFunctionType(ty TypeLike) (FunctionType, bool) {
	if C.mlirTypeIsAFunction(ty.raw()) {
		return FunctionType{baseType: baseType(ty.raw())}, true
	}
	return FunctionType{}, false
}

func FunctionTypeId() TypeId { return TypeId(C.mlirFunctionTypeGetTypeID()) }
func FunctionTypeName() string {
	ref := BorrowedStringRef(C.mlirFunctionTypeGetName())
	return ref.String()
}

func (f FunctionType) NumInputs() int  { return int(C.mlirFunctionTypeGetNumInputs(f.raw())) }
func (f FunctionType) NumResults() int { return int(C.mlirFunctionTypeGetNumResults(f.raw())) }
func (f FunctionType) Input(pos int) Type {
	return Type{baseType: baseType(C.mlirFunctionTypeGetInput(f.raw(), C.intptr_t(pos)))}
}
func (f FunctionType) Result(pos int) Type {
	return Type{baseType: baseType(C.mlirFunctionTypeGetResult(f.raw(), C.intptr_t(pos)))}
}

//===----------------------------------------------------------------------===//
// Opaque type.
//===----------------------------------------------------------------------===//

type OpaqueType struct {
	baseType
}

func NewOpaqueType(ctx Context, dialectNamespace string, typeData string) OpaqueType {
	cDialectNamespace := C.mlirStringRefCreate(C.CString(dialectNamespace), C.size_t(len(dialectNamespace)))
	cTypeData := C.mlirStringRefCreate(C.CString(typeData), C.size_t(len(typeData)))
	return OpaqueType{
		baseType: baseType(C.mlirOpaqueTypeGet(ctx.raw(), cDialectNamespace, cTypeData)),
	}
}

func AsOpaqueType(ty TypeLike) (OpaqueType, bool) {
	if C.mlirTypeIsAOpaque(ty.raw()) {
		return OpaqueType{baseType: baseType(ty.raw())}, true
	}
	return OpaqueType{}, false
}

func OpaqueTypeId() TypeId { return TypeId(C.mlirOpaqueTypeGetTypeID()) }
func OpaqueTypeName() string {
	ref := BorrowedStringRef(C.mlirOpaqueTypeGetName())
	return ref.String()
}

func (o OpaqueType) DialectNamespace() string {
	ref := BorrowedStringRef(C.mlirOpaqueTypeGetDialectNamespace(o.raw()))
	return ref.String()
}

func (o OpaqueType) Data() string {
	ref := BorrowedStringRef(C.mlirOpaqueTypeGetData(o.raw()))
	return ref.String()
}
