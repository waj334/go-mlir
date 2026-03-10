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
		baseType: baseType(C.mlirIntegerTypeGet(ctx.Raw(), C.unsigned(width))),
	}
}

func NewSignedIntegerType(ctx Context, width int) IntegerType {
	return IntegerType{
		baseType: baseType(C.mlirIntegerTypeSignedGet(ctx.Raw(), C.unsigned(width))),
	}
}

func NewUnsignedIntegerType(ctx Context, width int) IntegerType {
	return IntegerType{
		baseType: baseType(C.mlirIntegerTypeUnsignedGet(ctx.Raw(), C.unsigned(width))),
	}
}

func IntegerTypeId() TypeId { return TypeId(C.mlirIntegerTypeGetTypeID()) }

func AsIntegerType(ty TypeLike) (IntegerType, bool) {
	if C.mlirTypeIsAInteger(ty.Raw()) {
		return IntegerType{baseType: baseType(ty.Raw())}, true
	}
	return IntegerType{}, false
}

func (i IntegerType) Width() int       { return int(C.mlirIntegerTypeGetWidth(i.Raw())) }
func (i IntegerType) IsSignless() bool { return bool(C.mlirIntegerTypeIsSignless(i.Raw())) }
func (i IntegerType) IsSigned() bool   { return bool(C.mlirIntegerTypeIsSigned(i.Raw())) }
func (i IntegerType) IsUnsigned() bool { return bool(C.mlirIntegerTypeIsUnsigned(i.Raw())) }

//===----------------------------------------------------------------------===//
// Index type.
//===----------------------------------------------------------------------===//

type IndexType struct {
	baseType
}

func NewIndexType(ctx Context) IndexType {
	return IndexType{baseType: baseType(C.mlirIndexTypeGet(ctx.Raw()))}
}

func AsIndexType(ty TypeLike) (IndexType, bool) {
	if C.mlirTypeIsAIndex(ty.Raw()) {
		return IndexType{baseType: baseType(ty.Raw())}, true
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
		return FloatType{baseType: baseType(info.get(ctx.Raw()))}
	}
	panic("unsupported float format")
}

func AsFloatType(ty TypeLike) (FloatType, bool) {
	if C.mlirTypeIsAFloat(ty.Raw()) {
		return FloatType{baseType: baseType(ty.Raw())}, true
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
	typeID := C.mlirTypeGetTypeID(f.Raw())
	for format, info := range floatFormats {
		if C.mlirTypeIDEqual(typeID, info.typeID()) {
			return format
		}
	}
	return FloatUnknown
}

func (f FloatType) Width() int {
	return int(C.mlirFloatTypeGetWidth(f.Raw()))
}

//===----------------------------------------------------------------------===//
// None type.
//===----------------------------------------------------------------------===//

type NoneType struct {
	baseType
}

func NewNoneType(ctx Context) NoneType {
	return NoneType{baseType: baseType(C.mlirNoneTypeGet(ctx.Raw()))}
}

func AsNoneType(ty TypeLike) (NoneType, bool) {
	if C.mlirTypeIsANone(ty.Raw()) {
		return NoneType{baseType: baseType(ty.Raw())}, true
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

func NewComplexType(elementType FloatType) ComplexType {
	return ComplexType{baseType: baseType(C.mlirComplexTypeGet(elementType.Raw()))}
}

func AsComplexType(ty TypeLike) (ComplexType, bool) {
	if C.mlirTypeIsAComplex(ty.Raw()) {
		return ComplexType{baseType: baseType(ty.Raw())}, true
	}
	return ComplexType{}, false
}

func ComplexTypeId() TypeId { return TypeId(C.mlirComplexTypeGetTypeID()) }
func ComplexTypeName() string {
	ref := BorrowedStringRef(C.mlirComplexTypeGetName())
	return ref.String()
}

func (c ComplexType) ElementType() FloatType {
	return FloatType{baseType: baseType(C.mlirComplexTypeGetElementType(c.Raw()))}
}

//===----------------------------------------------------------------------===//
// Shaped type.
//===----------------------------------------------------------------------===//

type ShapedType struct {
	baseType
}

func AsShapedType(ty TypeLike) (ShapedType, bool) {
	if C.mlirTypeIsAShaped(ty.Raw()) {
		return ShapedType{baseType: baseType(ty.Raw())}, true
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
	return Type{baseType: baseType(C.mlirShapedTypeGetElementType(s.Raw()))}
}

func (s ShapedType) HasRank() bool        { return bool(C.mlirShapedTypeHasRank(s.Raw())) }
func (s ShapedType) Rank() int            { return int(C.mlirShapedTypeGetRank(s.Raw())) }
func (s ShapedType) HasStaticShape() bool { return bool(C.mlirShapedTypeHasStaticShape(s.Raw())) }
func (s ShapedType) IsDynamicDim(dim int) bool {
	return bool(C.mlirShapedTypeIsDynamicDim(s.Raw(), C.intptr_t(dim)))
}

func (s ShapedType) IsStaticDim(dim int) bool {
	return bool(C.mlirShapedTypeIsStaticDim(s.Raw(), C.intptr_t(dim)))
}

func (s ShapedType) DimSize(dim int) int {
	return int(C.mlirShapedTypeGetDimSize(s.Raw(), C.intptr_t(dim)))
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
		baseType: baseType(C.mlirVectorTypeGet(C.intptr_t(len(shape)), cShape, elementType.Raw())),
	}
}

func NewVectorTypeChecked(loc LocationLike, shape []int64, elementType TypeLike) VectorType {
	var cShape *C.int64_t
	if len(shape) > 0 {
		cShape = (*C.int64_t)(unsafe.Pointer(&shape[0]))
	}
	return VectorType{
		baseType: baseType(C.mlirVectorTypeGetChecked(loc.Raw(), C.intptr_t(len(shape)), cShape, elementType.Raw())),
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
		baseType: baseType(C.mlirVectorTypeGetScalable(C.intptr_t(len(shape)), cShape, cScalable, elementType.Raw())),
	}
}

func NewScalableVectorTypeChecked(loc LocationLike, shape []int64, scalable []bool, elementType TypeLike) VectorType {
	var cShape *C.int64_t
	if len(shape) > 0 {
		cShape = (*C.int64_t)(unsafe.Pointer(&shape[0]))
	}
	var cScalable *C.bool
	if len(scalable) > 0 {
		cScalable = (*C.bool)(unsafe.Pointer(&scalable[0]))
	}
	return VectorType{
		baseType: baseType(C.mlirVectorTypeGetScalableChecked(loc.Raw(), C.intptr_t(len(shape)), cShape, cScalable, elementType.Raw())),
	}
}

func AsVectorType(ty TypeLike) (VectorType, bool) {
	if C.mlirTypeIsAVector(ty.Raw()) {
		return VectorType{baseType: baseType(ty.Raw())}, true
	}
	return VectorType{}, false
}

func VectorTypeId() TypeId { return TypeId(C.mlirVectorTypeGetTypeID()) }
func VectorTypeName() string {
	ref := BorrowedStringRef(C.mlirVectorTypeGetName())
	return ref.String()
}

func (v VectorType) IsScalable() bool { return bool(C.mlirVectorTypeIsScalable(v.Raw())) }
func (v VectorType) IsDimScalable(dim int) bool {
	return bool(C.mlirVectorTypeIsDimScalable(v.Raw(), C.intptr_t(dim)))
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
		baseType: baseType(C.mlirRankedTensorTypeGet(C.intptr_t(len(shape)), cShape, elementType.Raw(), encoding.Raw())),
	}
}

func NewRankedTensorTypeChecked(loc LocationLike, shape []int64, elementType TypeLike, encoding Attribute) RankedTensorType {
	var cShape *C.int64_t
	if len(shape) > 0 {
		cShape = (*C.int64_t)(unsafe.Pointer(&shape[0]))
	}
	return RankedTensorType{
		baseType: baseType(C.mlirRankedTensorTypeGetChecked(loc.Raw(), C.intptr_t(len(shape)), cShape, elementType.Raw(), encoding.Raw())),
	}
}

func AsRankedTensorType(ty TypeLike) (RankedTensorType, bool) {
	if C.mlirTypeIsARankedTensor(ty.Raw()) {
		return RankedTensorType{baseType: baseType(ty.Raw())}, true
	}
	return RankedTensorType{}, false
}

func RankedTensorTypeId() TypeId { return TypeId(C.mlirRankedTensorTypeGetTypeID()) }
func RankedTensorTypeName() string {
	ref := BorrowedStringRef(C.mlirRankedTensorTypeGetName())
	return ref.String()
}

func (t RankedTensorType) Encoding() Attribute {
	return WrapAttribute(C.mlirRankedTensorTypeGetEncoding(t.Raw()))
}

type UnrankedTensorType struct {
	baseType
}

func NewUnrankedTensorType(elementType TypeLike) UnrankedTensorType {
	return UnrankedTensorType{
		baseType: baseType(C.mlirUnrankedTensorTypeGet(elementType.Raw())),
	}
}

func NewUnrankedTensorTypeChecked(loc LocationLike, elementType TypeLike) UnrankedTensorType {
	return UnrankedTensorType{
		baseType: baseType(C.mlirUnrankedTensorTypeGetChecked(loc.Raw(), elementType.Raw())),
	}
}

func AsUnrankedTensorType(ty TypeLike) (UnrankedTensorType, bool) {
	if C.mlirTypeIsAUnrankedTensor(ty.Raw()) {
		return UnrankedTensorType{baseType: baseType(ty.Raw())}, true
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
		baseType: baseType(C.mlirMemRefTypeGet(elementType.Raw(), C.intptr_t(len(shape)), cShape, layout.Raw(), memorySpace.Raw())),
	}
}

func NewMemRefTypeChecked(loc LocationLike, elementType TypeLike, shape []int64, layout Attribute, memorySpace Attribute) MemRefType {
	var cShape *C.int64_t
	if len(shape) > 0 {
		cShape = (*C.int64_t)(unsafe.Pointer(&shape[0]))
	}
	return MemRefType{
		baseType: baseType(C.mlirMemRefTypeGetChecked(loc.Raw(), elementType.Raw(), C.intptr_t(len(shape)), cShape, layout.Raw(), memorySpace.Raw())),
	}
}

func NewContiguousMemRefType(elementType TypeLike, shape []int64, memorySpace Attribute) MemRefType {
	var cShape *C.int64_t
	if len(shape) > 0 {
		cShape = (*C.int64_t)(unsafe.Pointer(&shape[0]))
	}
	return MemRefType{
		baseType: baseType(C.mlirMemRefTypeContiguousGet(elementType.Raw(), C.intptr_t(len(shape)), cShape, memorySpace.Raw())),
	}
}

func NewContiguousMemRefTypeChecked(loc LocationLike, elementType TypeLike, shape []int64, memorySpace Attribute) MemRefType {
	var cShape *C.int64_t
	if len(shape) > 0 {
		cShape = (*C.int64_t)(unsafe.Pointer(&shape[0]))
	}
	return MemRefType{
		baseType: baseType(C.mlirMemRefTypeContiguousGetChecked(loc.Raw(), elementType.Raw(), C.intptr_t(len(shape)), cShape, memorySpace.Raw())),
	}
}

func AsMemRefType(ty TypeLike) (MemRefType, bool) {
	if C.mlirTypeIsAMemRef(ty.Raw()) {
		return MemRefType{baseType: baseType(ty.Raw())}, true
	}
	return MemRefType{}, false
}

func MemRefTypeId() TypeId { return TypeId(C.mlirMemRefTypeGetTypeID()) }
func MemRefTypeName() string {
	ref := BorrowedStringRef(C.mlirMemRefTypeGetName())
	return ref.String()
}

func (m MemRefType) Layout() Attribute {
	return WrapAttribute(C.mlirMemRefTypeGetLayout(m.Raw()))
}

func (m MemRefType) AffineMap() AffineMap {
	return AffineMap(C.mlirMemRefTypeGetAffineMap(m.Raw()))
}

func (m MemRefType) MemorySpace() Attribute {
	return WrapAttribute(C.mlirMemRefTypeGetMemorySpace(m.Raw()))
}

type UnrankedMemRefType struct {
	baseType
}

func NewUnrankedMemRefType(elementType TypeLike, memorySpace Attribute) UnrankedMemRefType {
	return UnrankedMemRefType{
		baseType: baseType(C.mlirUnrankedMemRefTypeGet(elementType.Raw(), memorySpace.Raw())),
	}
}

func NewUnrankedMemRefTypeChecked(loc LocationLike, elementType TypeLike, memorySpace Attribute) UnrankedMemRefType {
	return UnrankedMemRefType{
		baseType: baseType(C.mlirUnrankedMemRefTypeGetChecked(loc.Raw(), elementType.Raw(), memorySpace.Raw())),
	}
}

func AsUnrankedMemRefType(ty TypeLike) (UnrankedMemRefType, bool) {
	if C.mlirTypeIsAUnrankedMemRef(ty.Raw()) {
		return UnrankedMemRefType{baseType: baseType(ty.Raw())}, true
	}
	return UnrankedMemRefType{}, false
}

func UnrankedMemRefTypeId() TypeId { return TypeId(C.mlirUnrankedMemRefTypeGetTypeID()) }
func UnrankedMemRefTypeName() string {
	ref := BorrowedStringRef(C.mlirUnrankedMemRefTypeGetName())
	return ref.String()
}

func (m UnrankedMemRefType) MemorySpace() Attribute {
	return WrapAttribute(C.mlirUnrankedMemrefGetMemorySpace(m.Raw()))
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
		cElements = &UnwrapTypeSlice(elements)[0]
	}
	return TupleType{
		baseType: baseType(C.mlirTupleTypeGet(ctx.Raw(), C.intptr_t(len(elements)), cElements)),
	}
}

func AsTupleType(ty TypeLike) (TupleType, bool) {
	if C.mlirTypeIsATuple(ty.Raw()) {
		return TupleType{baseType: baseType(ty.Raw())}, true
	}
	return TupleType{}, false
}

func TupleTypeId() TypeId { return TypeId(C.mlirTupleTypeGetTypeID()) }
func TupleTypeName() string {
	ref := BorrowedStringRef(C.mlirTupleTypeGetName())
	return ref.String()
}

func (t TupleType) NumTypes() int { return int(C.mlirTupleTypeGetNumTypes(t.Raw())) }
func (t TupleType) Type(pos int) Type {
	return Type{baseType: baseType(C.mlirTupleTypeGetType(t.Raw(), C.intptr_t(pos)))}
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
		cInputs = &UnwrapTypeSlice(inputs)[0]
	}
	var cResults *C.MlirType
	if len(results) > 0 {
		cResults = &UnwrapTypeSlice(results)[0]
	}

	return FunctionType{
		baseType: baseType(C.mlirFunctionTypeGet(ctx.Raw(), C.intptr_t(len(inputs)), cInputs, C.intptr_t(len(results)), cResults)),
	}
}

func AsFunctionType(ty TypeLike) (FunctionType, bool) {
	if C.mlirTypeIsAFunction(ty.Raw()) {
		return FunctionType{baseType: baseType(ty.Raw())}, true
	}
	return FunctionType{}, false
}

func FunctionTypeId() TypeId { return TypeId(C.mlirFunctionTypeGetTypeID()) }
func FunctionTypeName() string {
	ref := BorrowedStringRef(C.mlirFunctionTypeGetName())
	return ref.String()
}

func (f FunctionType) NumInputs() int  { return int(C.mlirFunctionTypeGetNumInputs(f.Raw())) }
func (f FunctionType) NumResults() int { return int(C.mlirFunctionTypeGetNumResults(f.Raw())) }
func (f FunctionType) Input(pos int) Type {
	return Type{baseType: baseType(C.mlirFunctionTypeGetInput(f.Raw(), C.intptr_t(pos)))}
}
func (f FunctionType) Result(pos int) Type {
	return Type{baseType: baseType(C.mlirFunctionTypeGetResult(f.Raw(), C.intptr_t(pos)))}
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
		baseType: baseType(C.mlirOpaqueTypeGet(ctx.Raw(), cDialectNamespace, cTypeData)),
	}
}

func AsOpaqueType(ty TypeLike) (OpaqueType, bool) {
	if C.mlirTypeIsAOpaque(ty.Raw()) {
		return OpaqueType{baseType: baseType(ty.Raw())}, true
	}
	return OpaqueType{}, false
}

func OpaqueTypeId() TypeId { return TypeId(C.mlirOpaqueTypeGetTypeID()) }
func OpaqueTypeName() string {
	ref := BorrowedStringRef(C.mlirOpaqueTypeGetName())
	return ref.String()
}

func (o OpaqueType) DialectNamespace() string {
	ref := BorrowedStringRef(C.mlirOpaqueTypeGetDialectNamespace(o.Raw()))
	return ref.String()
}

func (o OpaqueType) Data() string {
	ref := BorrowedStringRef(C.mlirOpaqueTypeGetData(o.Raw()))
	return ref.String()
}
