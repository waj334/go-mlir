package mlir

/*
#include <mlir-c/IR.h>

extern MlirWalkResult goMlirOperationWalkCallback(MlirOperation, void*);
static inline MlirOperationWalkCallback getGoMlirOperationWalkCallback(void) {
  return (MlirOperationWalkCallback)goMlirOperationWalkCallback;
}

typedef void (*MlirWalkSymbolTablesCallback)(MlirOperation, bool, void *userData);
extern void goMlirWalkSymbolTablesCallback(MlirOperation, bool, void*);
static inline MlirWalkSymbolTablesCallback getGoMlirWalkSymbolTablesCallback(void) {
  return (MlirWalkSymbolTablesCallback)goMlirWalkSymbolTablesCallback;
}

static inline bool mlirIRPrintingFlagsIsNull(MlirOpPrintingFlags flags) {
  return !flags.ptr;
}
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

//===----------------------------------------------------------------------===//
// Context API.
//===----------------------------------------------------------------------===//

type Context C.MlirContext

type ContextConfig struct {
	DisableThreading bool
	Registry         DialectRegistry
}

// NewContext creates an MLIR context and transfers its ownership to the caller.
func NewContext() Context {
	return Context(C.mlirContextCreate())
}

func NewContextWithConfig(config ContextConfig) Context {
	if !config.Registry.IsNull() {
		return Context(C.mlirContextCreateWithRegistry(config.Registry.Raw(), C.bool(!config.DisableThreading)))
	} else if config.DisableThreading {
		return Context(C.mlirContextCreateWithThreading(C.bool(!config.DisableThreading)))
	}
	return Context(C.mlirContextCreate())
}

func (c Context) Raw() C.MlirContext  { return C.MlirContext(c) }
func (c Context) Ptr() unsafe.Pointer { return unsafe.Pointer(c.Raw().ptr) }
func (c Context) IsNull() bool        { return bool(C.mlirContextIsNull(c.Raw())) }

// Destroy takes an MLIR context owned by the caller and destroys it.
func (c Context) Destroy() {
	C.mlirContextDestroy(c.Raw())
}

// Equal checks if two contexts are equal.
func (c Context) Equal(other Context) bool {
	return bool(C.mlirContextEqual(c.Raw(), other.Raw()))
}

// SetAllowUnregisteredDialects sets whether unregistered dialects are allowed in this context.
func (c Context) SetAllowUnregisteredDialects(allow bool) {
	C.mlirContextSetAllowUnregisteredDialects(c.Raw(), C.bool(allow))
}

// UnregisteredDialectsAllowed returns whether the context allows unregistered dialects.
func (c Context) UnregisteredDialectsAllowed() bool {
	return bool(C.mlirContextGetAllowUnregisteredDialects(c.Raw()))
}

// NumRegisteredDialects returns the number of dialects registered with the given context. A
// registered dialect will be loaded if needed by the parser.
func (c Context) NumRegisteredDialects() int {
	return int(C.mlirContextGetNumRegisteredDialects(c.Raw()))
}

// AppendDialectRegistry append the contents of the given dialect registry to the registry associated
// with the context.
func (c Context) AppendDialectRegistry(registry DialectRegistry) {
	C.mlirContextAppendDialectRegistry(c.Raw(), registry.Raw())
}

// NumLoadedDialects returns the number of dialects loaded by the context.
func (c Context) NumLoadedDialects() int {
	return int(C.mlirContextGetNumLoadedDialects(c.Raw()))
}

// GetOrLoadDialect gets the dialect instance owned by the given context using the dialect
// namespace to identify it, loads (i.e., constructs the instance of) the
// dialect if necessary. If the dialect is not registered with the context,
// returns null. Use mlirContextLoad<Name>Dialect to load an unregistered
// dialect.
func (c Context) GetOrLoadDialect(name string) Dialect {
	refName := NewStringRef(name)
	defer refName.Destroy()
	return Dialect(C.mlirContextGetOrLoadDialect(c.Raw(), refName.Raw()))
}

// EnableMultithreading sets the multithreading mode (must be set to false if using mlir-print-ir-after-all).
func (c Context) EnableMultithreading(enable bool) {
	C.mlirContextEnableMultithreading(c.Raw(), C.bool(enable))
}

// LoadAllAvailableDialects Eagerly loads all available dialects registered with a context, making
// them available for use for IR construction.
func (c Context) LoadAllAvailableDialects() {
	C.mlirContextLoadAllAvailableDialects(c.Raw())
}

func (c Context) SetThreadPool(pool LLVMThreadPool) {
	C.mlirContextSetThreadPool(c.Raw(), pool.Raw())
}

func (c Context) NumThreads() int {
	return int(C.mlirContextGetNumThreads(c.Raw()))
}

func (c Context) ThreadPool() LLVMThreadPool {
	return LLVMThreadPool(C.mlirContextGetThreadPool(c.Raw()))
}

//===----------------------------------------------------------------------===//
// Dialect API.
//===----------------------------------------------------------------------===//

type Dialect C.MlirDialect

func (d Dialect) Raw() C.MlirDialect {
	return C.MlirDialect(d)
}

func (d Dialect) IsNull() bool {
	return bool(C.mlirDialectIsNull(d.Raw()))
}

func (d Dialect) Equal(other Dialect) bool {
	return bool(C.mlirDialectEqual(d.Raw(), other.Raw()))
}

func (d Dialect) Namespace() string {
	ref := BorrowedStringRef(C.mlirDialectGetNamespace(d.Raw()))
	return ref.String()
}

//===----------------------------------------------------------------------===//
// DialectHandle API.
//===----------------------------------------------------------------------===//

type DialectHandle C.MlirDialectHandle

func WrapExternalDialectHandle(pointer unsafe.Pointer) DialectHandle {
	return DialectHandle(C.MlirDialectHandle{ptr: pointer})
}

func (d DialectHandle) Raw() C.MlirDialectHandle {
	return C.MlirDialectHandle(d)
}

func (d DialectHandle) Namespace() string {
	ref := BorrowedStringRef(C.mlirDialectHandleGetNamespace(d.Raw()))
	return ref.String()
}

func (d DialectHandle) InsertDialect(registry DialectRegistry) {
	C.mlirDialectHandleInsertDialect(d.Raw(), registry.Raw())
}

func (d DialectHandle) RegisterDialect(context Context) {
	C.mlirDialectHandleRegisterDialect(d.Raw(), context.Raw())
}

func (d DialectHandle) LoadDialect(context Context) {
	C.mlirDialectHandleLoadDialect(d.Raw(), context.Raw())
}

//===----------------------------------------------------------------------===//
// DialectRegistry API.
//===----------------------------------------------------------------------===//

type DialectRegistry C.MlirDialectRegistry

func NewDialectRegistry() DialectRegistry {
	return DialectRegistry(C.mlirDialectRegistryCreate())
}

func (d DialectRegistry) Raw() C.MlirDialectRegistry {
	return C.MlirDialectRegistry(d)
}

func (d DialectRegistry) IsNull() bool {
	return bool(C.mlirDialectRegistryIsNull(d.Raw()))
}

func (d DialectRegistry) Destroy() {
	C.mlirDialectRegistryDestroy(d.Raw())
}

//===----------------------------------------------------------------------===//
// Location API.
//===----------------------------------------------------------------------===//

type LocationLike interface {
	Raw() C.MlirLocation
	Ptr() unsafe.Pointer
	IsNull() bool
	Attribute() Attribute
	Context() Context
	Equal(other LocationLike) bool
	String() string
	AsLocation() Location
}

type baseLocation C.MlirLocation

func (l baseLocation) Raw() C.MlirLocation { return C.MlirLocation(l) }
func (l baseLocation) Ptr() unsafe.Pointer { return unsafe.Pointer(l.Raw().ptr) }
func (l baseLocation) IsNull() bool        { return bool(C.mlirLocationIsNull(l.Raw())) }
func (l baseLocation) Attribute() Attribute {
	return WrapAttribute(C.mlirLocationGetAttribute(l.Raw()))
}
func (l baseLocation) Context() Context { return Context(C.mlirLocationGetContext(l.Raw())) }
func (l baseLocation) Equal(other LocationLike) bool {
	return bool(C.mlirLocationEqual(l.Raw(), other.Raw()))
}

func (l baseLocation) String() string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirLocationPrint(l.Raw(), cb, ud)
	})
}

type Location struct {
	baseLocation
}

func wrapLocation(v C.MlirLocation) Location {
	return Location{baseLocation: baseLocation(v)}
}

func (l Location) AsLocation() Location      { return l }
func (l Location) IsAFileLineColRange() bool { return bool(C.mlirLocationIsAFileLineColRange(l.Raw())) }
func (l Location) IsACallSite() bool         { return bool(C.mlirLocationIsACallSite(l.Raw())) }
func (l Location) AsFileLineColRange() (FileLineColRangeLoc, bool) {
	if bool(C.mlirLocationIsAFileLineColRange(l.Raw())) {
		return FileLineColRangeLoc{baseLocation: l.baseLocation}, true
	}
	return FileLineColRangeLoc{}, false
}

func (l Location) AsName() (NameLoc, bool) {
	if bool(C.mlirLocationIsAName(l.Raw())) {
		return NameLoc{baseLocation: l.baseLocation}, true
	}
	return NameLoc{}, false
}

func (l Location) AsCallSite() (CallSiteLoc, bool) {
	if bool(C.mlirLocationIsACallSite(l.Raw())) {
		return CallSiteLoc{baseLocation: l.baseLocation}, true
	}
	return CallSiteLoc{}, false
}

type FileLineColLoc struct {
	baseLocation
}

func NewFileLineCol(ctx Context, filename string, line, column int) FileLineColLoc {
	refFilename := NewStringRef(filename)
	defer refFilename.Destroy()
	return FileLineColLoc{baseLocation(C.mlirLocationFileLineColGet(
		ctx.Raw(),
		refFilename.Raw(),
		C.unsigned(line),
		C.unsigned(column)))}
}

func (l FileLineColLoc) AsLocation() Location { return Location{baseLocation: l.baseLocation} }

type FileLineColRangeLoc struct {
	baseLocation
}

func NewFileLineColRange(ctx Context, filename string, startLine, startColumn, endLine, endColumn int) FileLineColRangeLoc {
	refFilename := NewStringRef(filename)
	defer refFilename.Destroy()
	return FileLineColRangeLoc{
		baseLocation: baseLocation(C.mlirLocationFileLineColRangeGet(
			ctx.Raw(),
			refFilename.Raw(),
			C.unsigned(startLine),
			C.unsigned(startColumn),
			C.unsigned(endLine),
			C.unsigned(endColumn))),
	}
}

func FileLineColRangeLocTypeId() TypeId {
	return TypeId(C.mlirLocationFileLineColRangeGetTypeID())
}

func (l FileLineColRangeLoc) AsLocation() Location { return Location{baseLocation: l.baseLocation} }

func (l FileLineColRangeLoc) Filename() Identifier {
	return Identifier(C.mlirLocationFileLineColRangeGetFilename(l.Raw()))
}

func (l FileLineColRangeLoc) StartLine() int {
	return int(C.mlirLocationFileLineColRangeGetStartLine(l.Raw()))
}
func (l FileLineColRangeLoc) StartColumn() int {
	return int(C.mlirLocationFileLineColRangeGetStartColumn(l.Raw()))
}
func (l FileLineColRangeLoc) EndLine() int {
	return int(C.mlirLocationFileLineColRangeGetEndLine(l.Raw()))
}
func (l FileLineColRangeLoc) EndColumn() int {
	return int(C.mlirLocationFileLineColRangeGetEndColumn(l.Raw()))
}

type CallSiteLoc struct {
	baseLocation
}

func NewCallSite(callee, caller Location) CallSiteLoc {
	return CallSiteLoc{
		baseLocation: baseLocation(C.mlirLocationCallSiteGet(callee.Raw(), caller.Raw())),
	}
}

func CallSiteLocTypeId() TypeId { return TypeId(C.mlirLocationCallSiteGetTypeID()) }

func (l CallSiteLoc) AsLocation() Location { return Location{baseLocation: l.baseLocation} }
func (l CallSiteLoc) Callee() Location     { return wrapLocation(C.mlirLocationCallSiteGetCallee(l.Raw())) }
func (l CallSiteLoc) Caller() Location     { return wrapLocation(C.mlirLocationCallSiteGetCaller(l.Raw())) }

type FusedLoc struct {
	baseLocation
}

func NewFusedLoc(ctx Context, locations []LocationLike, metadata AttributeLike) FusedLoc {
	raw := make([]C.MlirLocation, len(locations))
	for i, loc := range locations {
		raw[i] = loc.Raw()
	}

	var ptr *C.MlirLocation
	if len(raw) != 0 {
		ptr = &raw[0]
	}

	return FusedLoc{
		baseLocation: baseLocation(C.mlirLocationFusedGet(
			ctx.Raw(),
			C.intptr_t(len(raw)),
			ptr,
			metadata.Raw())),
	}
}

func FusedLocTypeId() TypeId { return TypeId(C.mlirLocationFusedGetTypeID()) }

func (l FusedLoc) AsLocation() Location { return Location{baseLocation: l.baseLocation} }
func (l FusedLoc) NumLocations() int    { return int(C.mlirLocationFusedGetNumLocations(l.Raw())) }
func (l FusedLoc) Locations() []Location {
	n := l.NumLocations()
	if n == 0 {
		return nil
	}

	raw := make([]C.MlirLocation, n)
	C.mlirLocationFusedGetLocations(l.Raw(), &raw[0])

	out := make([]Location, n)
	for i := range raw {
		out[i] = Location{baseLocation: baseLocation(raw[i])}
	}
	return out
}

func (l FusedLoc) Metadata() Attribute { return WrapAttribute(C.mlirLocationFusedGetMetadata(l.Raw())) }

type NameLoc struct {
	baseLocation
}

func NameLocTypeId() TypeId { return TypeId(C.mlirLocationNameGetTypeID()) }

func NewName(ctx Context, name string, childLoc LocationLike) NameLoc {
	refName := NewStringRef(name)
	defer refName.Destroy()
	return NameLoc{
		baseLocation: baseLocation(C.mlirLocationNameGet(ctx.Raw(), refName.Raw(), childLoc.Raw())),
	}
}

func (l NameLoc) AsLocation() Location { return Location{baseLocation: l.baseLocation} }
func (l NameLoc) Name() Identifier     { return Identifier(C.mlirLocationNameGetName(l.Raw())) }
func (l NameLoc) ChildLoc() Location   { return wrapLocation(C.mlirLocationNameGetChildLoc(l.Raw())) }

type UnknownLoc struct {
	baseLocation
}

func NewUnknownLoc(ctx Context) UnknownLoc {
	return UnknownLoc{
		baseLocation: baseLocation(C.mlirLocationUnknownGet(ctx.Raw())),
	}
}

func (l UnknownLoc) AsLocation() Location { return Location{baseLocation: l.baseLocation} }

//===----------------------------------------------------------------------===//
// Module API.
//===----------------------------------------------------------------------===//

type Module C.MlirModule

func NewModule(loc LocationLike) Module {
	return Module(C.mlirModuleCreateEmpty(loc.Raw()))
}

func NewModuleFromString(ctx Context, input string) Module {
	refInput := NewStringRef(input)
	defer refInput.Destroy()
	return Module(C.mlirModuleCreateParse(ctx.Raw(), refInput.Raw()))
}

func NewModuleFromFile(ctx Context, filename string) Module {
	refFilename := NewStringRef(filename)
	defer refFilename.Destroy()
	return Module(C.mlirModuleCreateParseFromFile(ctx.Raw(), refFilename.Raw()))
}

func NewModuleFromOperation(op Operation) Module {
	return Module(C.mlirModuleFromOperation(op.Raw()))
}

func (m Module) Raw() C.MlirModule       { return C.MlirModule(m) }
func (m Module) Ptr() unsafe.Pointer     { return unsafe.Pointer(m.Raw().ptr) }
func (m Module) IsNull() bool            { return bool(C.mlirModuleIsNull(m.Raw())) }
func (m Module) Destroy()                { C.mlirModuleDestroy(m.Raw()) }
func (m Module) Context() Context        { return Context(C.mlirModuleGetContext(m.Raw())) }
func (m Module) Body() Block             { return Block(C.mlirModuleGetBody(m.Raw())) }
func (m Module) Operation() Operation    { return Operation(C.mlirModuleGetOperation(m.Raw())) }
func (m Module) Equal(other Module) bool { return bool(C.mlirModuleEqual(m.Raw(), other.Raw())) }
func (m Module) HashValue() uintptr      { return uintptr(C.mlirModuleHashValue(m.Raw())) }

//===----------------------------------------------------------------------===//
// Operation state.
//===----------------------------------------------------------------------===//

type OperationState C.MlirOperationState

func NewOperationState(name string, loc LocationLike) *OperationState {
	refName := NewStringRef(name)
	s := OperationState(C.mlirOperationStateGet(refName.Raw(), loc.Raw()))
	return &s
}

func (s *OperationState) Raw() *C.MlirOperationState { return (*C.MlirOperationState)(s) }
func (s *OperationState) AddResults(results ...TypeLike) *OperationState {
	if len(results) > 0 {
		rawResults := UnwrapTypeSlice(results)
		C.mlirOperationStateAddResults(s.Raw(), C.intptr_t(len(results)), (*C.MlirType)(&rawResults[0]))
	}
	return s
}

func (s *OperationState) AddOperands(operands ...Value) *OperationState {
	if len(operands) > 0 {
		C.mlirOperationStateAddOperands(s.Raw(), C.intptr_t(len(operands)), (*C.MlirValue)(&operands[0]))
	}
	return s
}

func (s *OperationState) AddOwnedRegions(regions ...Region) *OperationState {
	if len(regions) > 0 {
		C.mlirOperationStateAddOwnedRegions(s.Raw(), C.intptr_t(len(regions)), (*C.MlirRegion)(&regions[0]))
	}
	return s
}

func (s *OperationState) AddSuccessors(successors ...Block) *OperationState {
	if len(successors) > 0 {
		C.mlirOperationStateAddSuccessors(s.Raw(), C.intptr_t(len(successors)), (*C.MlirBlock)(&successors[0]))
	}
	return s
}

func (s *OperationState) AddAttributes(attributes ...NamedAttribute) *OperationState {
	if len(attributes) > 0 {
		C.mlirOperationStateAddAttributes(s.Raw(), C.intptr_t(len(attributes)), (*C.MlirNamedAttribute)(&attributes[0]))
	}
	return s
}

func (s *OperationState) EnableResultTypeInference() *OperationState {
	C.mlirOperationStateEnableResultTypeInference(s.Raw())
	return s
}

func (s *OperationState) Create() Operation {
	return Operation(C.mlirOperationCreate(s.Raw()))
}

//===----------------------------------------------------------------------===//
// AsmState API.
//===----------------------------------------------------------------------===//

type AsmState C.MlirAsmState

func NewAsmStateFromOperation(op Operation, flags OpPrintingFlags) AsmState {
	return AsmState(C.mlirAsmStateCreateForOperation(op.Raw(), flags.Raw()))
}

func NewAsmStateFromValue(value Value, flags OpPrintingFlags) AsmState {
	return AsmState(C.mlirAsmStateCreateForValue(value.Raw(), flags.Raw()))
}

func (s AsmState) Raw() C.MlirAsmState { return C.MlirAsmState(s) }
func (s AsmState) Destroy()            { C.mlirAsmStateDestroy(s.Raw()) }

//===----------------------------------------------------------------------===//
// Op Printing flags API.
//===----------------------------------------------------------------------===//

type OpPrintingFlags C.MlirOpPrintingFlags

func NewOpPrintingFlags() OpPrintingFlags {
	return OpPrintingFlags(C.mlirOpPrintingFlagsCreate())
}

func (f OpPrintingFlags) Raw() C.MlirOpPrintingFlags { return C.MlirOpPrintingFlags(f) }
func (f OpPrintingFlags) Destroy()                   { C.mlirOpPrintingFlagsDestroy(f.Raw()) }
func (f OpPrintingFlags) IsNull() bool               { return bool(C.mlirIRPrintingFlagsIsNull(f.Raw())) }

func (f OpPrintingFlags) WithElidLargeElementsAttrs(limit int) OpPrintingFlags {
	C.mlirOpPrintingFlagsElideLargeElementsAttrs(f.Raw(), C.intptr_t(limit))
	return f
}

func (f OpPrintingFlags) WithElideLargeResourceString(limit int) OpPrintingFlags {
	C.mlirOpPrintingFlagsElideLargeResourceString(f.Raw(), C.intptr_t(limit))
	return f
}

func (f OpPrintingFlags) WithEnableDebugInfo(enable bool, pretty bool) OpPrintingFlags {
	C.mlirOpPrintingFlagsEnableDebugInfo(f.Raw(), C.bool(enable), C.bool(pretty))
	return f
}

func (f OpPrintingFlags) WithGenericOpForm() OpPrintingFlags {
	C.mlirOpPrintingFlagsPrintGenericOpForm(f.Raw())
	return f
}

func (f OpPrintingFlags) WithPrintNameLocAsPrefix() OpPrintingFlags {
	C.mlirOpPrintingFlagsPrintNameLocAsPrefix(f.Raw())
	return f
}

func (f OpPrintingFlags) WithUseLocalScope() OpPrintingFlags {
	C.mlirOpPrintingFlagsUseLocalScope(f.Raw())
	return f
}

func (f OpPrintingFlags) WithAssumeVerified() OpPrintingFlags {
	C.mlirOpPrintingFlagsAssumeVerified(f.Raw())
	return f
}

func (f OpPrintingFlags) WithSkipRegions() OpPrintingFlags {
	C.mlirOpPrintingFlagsSkipRegions(f.Raw())
	return f
}

//===----------------------------------------------------------------------===//
// Bytecode printing flags API.
//===----------------------------------------------------------------------===//

type BytecodeWriterConfig C.MlirBytecodeWriterConfig

func NewBytecodeWriterConfig() BytecodeWriterConfig {
	return BytecodeWriterConfig(C.mlirBytecodeWriterConfigCreate())
}

func (c BytecodeWriterConfig) Raw() C.MlirBytecodeWriterConfig { return C.MlirBytecodeWriterConfig(c) }

func (c BytecodeWriterConfig) SetDesiredEmitVersion(version int) {
	C.mlirBytecodeWriterConfigDesiredEmitVersion(c.Raw(), C.int64_t(version))
}

//===----------------------------------------------------------------------===//
// Operation API.
//===----------------------------------------------------------------------===//

type Operation C.MlirOperation

func NewOperation(state *OperationState) Operation {
	return Operation(C.mlirOperationCreate(state.Raw()))
}

func NewOperationFromString(ctx Context, source string, sourceName string) Operation {
	refSource := NewStringRef(source)
	defer refSource.Destroy()

	refSourceName := NewStringRef(sourceName)
	defer refSourceName.Destroy()

	return Operation(C.mlirOperationCreateParse(ctx.Raw(), refSource.Raw(), refSourceName.Raw()))
}

func WrapExternalOperation(pointer unsafe.Pointer) Operation {
	return Operation(C.MlirOperation{ptr: pointer})
}

func (o Operation) Raw() C.MlirOperation         { return C.MlirOperation(o) }
func (o Operation) Ptr() unsafe.Pointer          { return unsafe.Pointer(o.Raw().ptr) }
func (o Operation) IsNull() bool                 { return bool(C.mlirOperationIsNull(o.Raw())) }
func (o Operation) Context() Context             { return Context(C.mlirOperationGetContext(o.Raw())) }
func (o Operation) Clone() Operation             { return Operation(C.mlirOperationClone(o.Raw())) }
func (o Operation) Destroy()                     { C.mlirOperationDestroy(o.Raw()) }
func (o Operation) RemoveFromParent()            { C.mlirOperationRemoveFromParent(o.Raw()) }
func (o Operation) HashValue() uintptr           { return uintptr(C.mlirOperationHashValue(o.Raw())) }
func (o Operation) Location() Location           { return wrapLocation(C.mlirOperationGetLocation(o.Raw())) }
func (o Operation) SetLocation(loc LocationLike) { C.mlirOperationSetLocation(o.Raw(), loc.Raw()) }
func (o Operation) TypeId() TypeId               { return TypeId(C.mlirOperationGetTypeID(o.Raw())) }
func (o Operation) Name() Identifier             { return Identifier(C.mlirOperationGetName(o.Raw())) }
func (o Operation) Block() Block                 { return Block(C.mlirOperationGetBlock(o.Raw())) }
func (o Operation) NumRegions() int              { return int(C.mlirOperationGetNumRegions(o.Raw())) }
func (o Operation) NextInBlock() Operation       { return Operation(C.mlirOperationGetNextInBlock(o.Raw())) }
func (o Operation) NumOperands() int             { return int(C.mlirOperationGetNumOperands(o.Raw())) }
func (o Operation) NumResults() int              { return int(C.mlirOperationGetNumResults(o.Raw())) }
func (o Operation) NumSuccessors() int           { return int(C.mlirOperationGetNumSuccessors(o.Raw())) }
func (o Operation) NumAttributes() int           { return int(C.mlirOperationGetNumAttributes(o.Raw())) }

func (o Operation) Attribute(pos int) NamedAttribute {
	return NamedAttribute(C.mlirOperationGetAttribute(o.Raw(), C.intptr_t(pos)))
}

func (o Operation) AttributeByName(name string) Attribute {
	refName := NewStringRef(name)
	defer refName.Destroy()
	return WrapAttribute(C.mlirOperationGetAttributeByName(o.Raw(), refName.Raw()))
}

func (o Operation) SetAttributeByName(name string, attr AttributeLike) {
	refName := NewStringRef(name)
	defer refName.Destroy()
	C.mlirOperationSetAttributeByName(o.Raw(), refName.Raw(), attr.Raw())
}

func (o Operation) RemoveAttributeByName(name string) {
	refName := NewStringRef(name)
	defer refName.Destroy()
	C.mlirOperationRemoveAttributeByName(o.Raw(), refName.Raw())
}

func (o Operation) Successor(pos int) Block {
	return Block(C.mlirOperationGetSuccessor(o.Raw(), C.intptr_t(pos)))
}

func (o Operation) SetSuccessor(pos int, block Block) {
	C.mlirOperationSetSuccessor(o.Raw(), C.intptr_t(pos), block.Raw())
}

func (o Operation) Result(pos int) Result {
	return Result{Value: Value(C.mlirOperationGetResult(o.Raw(), C.intptr_t(pos)))}
}

func (o Operation) Operand(pos int) Value {
	return Value(C.mlirOperationGetOperand(o.Raw(), C.intptr_t(pos)))
}

func (o Operation) SetOperand(pos int, newValue Value) {
	C.mlirOperationSetOperand(o.Raw(), C.intptr_t(pos), newValue.Raw())
}

func (o Operation) SetOperands(operands ...Value) {
	if len(operands) > 0 {
		C.mlirOperationSetOperands(o.Raw(), C.intptr_t(len(operands)), (*C.MlirValue)(&operands[0]))
	}
}

func (o Operation) Region(pos int) Region {
	return Region(C.mlirOperationGetRegion(o.Raw(), C.intptr_t(pos)))
}

func (o Operation) ParentOperation() Operation {
	return Operation(C.mlirOperationGetParentOperation(o.Raw()))
}

func (o Operation) Equal(other Operation) bool {
	return bool(C.mlirOperationEqual(o.Raw(), other.Raw()))
}

func (o Operation) HasInherentAttributeByName(name string) bool {
	refName := NewStringRef(name)
	defer refName.Destroy()
	return bool(C.mlirOperationHasInherentAttributeByName(o.Raw(), refName.Raw()))
}

func (o Operation) InherentAttributeByName(name string) Attribute {
	refName := NewStringRef(name)
	defer refName.Destroy()
	return WrapAttribute(C.mlirOperationGetInherentAttributeByName(o.Raw(), refName.Raw()))
}

func (o Operation) SetInherentAttributeByName(name string, attr Attribute) {
	refName := NewStringRef(name)
	defer refName.Destroy()
	C.mlirOperationSetInherentAttributeByName(o.Raw(), refName.Raw(), attr.Raw())
}

func (o Operation) NumDiscardableAttributes() int {
	return int(C.mlirOperationGetNumDiscardableAttributes(o.Raw()))
}

func (o Operation) DiscardableAttribute(pos int) NamedAttribute {
	return NamedAttribute(C.mlirOperationGetDiscardableAttribute(o.Raw(), C.intptr_t(pos)))
}

func (o Operation) DiscardableAttributeByName(name string) Attribute {
	refName := NewStringRef(name)
	defer refName.Destroy()
	return WrapAttribute(C.mlirOperationGetDiscardableAttributeByName(o.Raw(), refName.Raw()))
}

func (o Operation) SetDiscardableAttributeByName(name string, attr Attribute) {
	refName := NewStringRef(name)
	defer refName.Destroy()
	C.mlirOperationSetDiscardableAttributeByName(o.Raw(), refName.Raw(), attr.Raw())
}

func (o Operation) RemoveDiscardableAttributeByName(name string) {
	refName := NewStringRef(name)
	defer refName.Destroy()
	C.mlirOperationRemoveDiscardableAttributeByName(o.Raw(), refName.Raw())
}

func (o Operation) String() string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirOperationPrint(o.Raw(), cb, ud)
	})
}

func (o Operation) StringWithFlags(flags OpPrintingFlags) string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirOperationPrintWithFlags(o.Raw(), flags.Raw(), cb, ud)
	})
}

func (o Operation) StringWithState(state AsmState) string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirOperationPrintWithState(o.Raw(), state.Raw(), cb, ud)
	})
}

func (o Operation) Bytecode() string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirOperationWriteBytecode(o.Raw(), cb, ud)
	})
}

func (o Operation) BytecodeWithConfig(config BytecodeWriterConfig) string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirOperationWriteBytecodeWithConfig(o.Raw(), config.Raw(), cb, ud)
	})
}

func (o Operation) Dump()                      { C.mlirOperationDump(o.Raw()) }
func (o Operation) Verify() bool               { return bool(C.mlirOperationVerify(o.Raw())) }
func (o Operation) MoveAfter(other Operation)  { C.mlirOperationMoveAfter(o.Raw(), other.Raw()) }
func (o Operation) MoveBefore(other Operation) { C.mlirOperationMoveBefore(o.Raw(), other.Raw()) }

func (o Operation) IsBeforeInBlock(other Operation) bool {
	return bool(C.mlirOperationIsBeforeInBlock(o.Raw(), other.Raw()))
}

type WalkResult C.MlirWalkResult

const (
	WalkResultAdvance   WalkResult = C.MlirWalkResultAdvance
	WalkResultInterrupt WalkResult = C.MlirWalkResultInterrupt
	WalkResultSkip      WalkResult = C.MlirWalkResultSkip
)

type WalkOrder C.MlirWalkOrder

const (
	WalkOrderPreOrder  WalkOrder = C.MlirWalkPreOrder
	WalkOrderPostOrder WalkOrder = C.MlirWalkPostOrder
)

type OperationWalkCallback func(Operation) WalkOrder

func NewOperationWalkCallback(fn func(Operation) WalkOrder) OperationWalkCallback {
	return OperationWalkCallback(fn)
}

func (cb OperationWalkCallback) Callback() (C.MlirOperationWalkCallback, unsafe.Pointer, func()) {
	handle := cgo.NewHandle(cb)
	return C.getGoMlirOperationWalkCallback(), unsafe.Pointer(uintptr(handle)), handle.Delete
}

//export goMlirOperationWalkCallback
func goMlirOperationWalkCallback(o C.MlirOperation, userdata unsafe.Pointer) C.MlirWalkResult {
	handle := cgo.Handle(uintptr(userdata))
	callback := handle.Value().(OperationWalkCallback)
	ret := callback(Operation(o))
	return C.MlirWalkResult(ret)
}

func (o Operation) Walk(walkOrder WalkOrder, fn OperationWalkCallback) {
	cb, userdata, cleanup := fn.Callback()
	defer cleanup()
	C.mlirOperationWalk(o.Raw(), cb, userdata, C.MlirWalkOrder(walkOrder))
}

func (o Operation) ReplaceUsesOfWith(of, with Value) {
	C.mlirOperationReplaceUsesOfWith(o.Raw(), of.Raw(), with.Raw())
}

type WalkSymbolTablesCallback func(Operation, bool)

func NewWalkSymbolTablesCallback(fn func(Operation, bool)) WalkSymbolTablesCallback {
	return WalkSymbolTablesCallback(fn)
}

func (cb WalkSymbolTablesCallback) Callback() (C.MlirWalkSymbolTablesCallback, unsafe.Pointer, func()) {
	handle := cgo.NewHandle(cb)
	return C.getGoMlirWalkSymbolTablesCallback(), unsafe.Pointer(uintptr(handle)), handle.Delete
}

//export goMlirWalkSymbolTablesCallback
func goMlirWalkSymbolTablesCallback(op C.MlirOperation, visible C.bool, userdata unsafe.Pointer) {
	handle := cgo.Handle(uintptr(userdata))
	callback := handle.Value().(WalkSymbolTablesCallback)
	callback(Operation(op), bool(visible))
}

func (o Operation) WalkSymbolTables(allSymbolUsesVisible bool, fn WalkSymbolTablesCallback) {
	cb, userdata, cleanup := fn.Callback()
	defer cleanup()
	C.mlirSymbolTableWalkSymbolTables(o.Raw(), C.bool(allSymbolUsesVisible), cb, userdata)
}

//===----------------------------------------------------------------------===//
// Region API.
//===----------------------------------------------------------------------===//

type Region C.MlirRegion

func NewRegion() Region {
	return Region(C.mlirRegionCreate())
}

func (r Region) Raw() C.MlirRegion            { return C.MlirRegion(r) }
func (r Region) Ptr() unsafe.Pointer          { return unsafe.Pointer(r.Raw().ptr) }
func (r Region) IsNull() bool                 { return bool(C.mlirRegionIsNull(r.Raw())) }
func (r Region) Destroy()                     { C.mlirRegionDestroy(r.Raw()) }
func (r Region) Equal(other Region) bool      { return bool(C.mlirRegionEqual(r.Raw(), other.Raw())) }
func (r Region) FirstBlock() Block            { return Block(C.mlirRegionGetFirstBlock(r.Raw())) }
func (r Region) AppendOwnedBlock(block Block) { C.mlirRegionAppendOwnedBlock(r.Raw(), block.Raw()) }
func (r Region) NextInOperation() Region      { return Region(C.mlirRegionGetNextInOperation(r.Raw())) }
func (r Region) TakeBody(source Region)       { C.mlirRegionTakeBody(r.Raw(), source.Raw()) }
func (r Region) InsertOwnedBlock(pos int, block Block) {
	C.mlirRegionInsertOwnedBlock(r.Raw(), C.intptr_t(pos), block.Raw())
}

func (r Region) InsertOwnedBlockAfter(reference, block Block) {
	C.mlirRegionInsertOwnedBlockAfter(r.Raw(), reference.Raw(), block.Raw())
}

func (r Region) InsertOwnedBlockBefore(reference, block Block) {
	C.mlirRegionInsertOwnedBlockBefore(r.Raw(), reference.Raw(), block.Raw())
}

//===----------------------------------------------------------------------===//
// Block API.
//===----------------------------------------------------------------------===//

type Block C.MlirBlock

func NewBlock(args []TypeLike, locs []LocationLike) Block {
	var rawArgsPtr *C.MlirType
	var rawLocsPtr *C.MlirLocation

	if len(args) != len(locs) {
		panic("number of arguments does match number of locations")
	}

	if len(args) > 0 {
		rawArgs := UnwrapTypeSlice(args)

		rawLocs := make([]C.MlirLocation, len(locs))
		for i, loc := range locs {
			rawLocs[i] = loc.Raw()
		}

		rawArgsPtr = (*C.MlirType)(&rawArgs[0])
		rawLocsPtr = (*C.MlirLocation)(&rawLocs[0])
	}

	return Block(C.mlirBlockCreate(C.intptr_t(len(args)), rawArgsPtr, rawLocsPtr))
}

func UnwrapBlockSlice(blocks []Block) []C.MlirBlock {
	rawTypes := make([]C.MlirBlock, len(blocks))
	for i, t := range blocks {
		rawTypes[i] = t.Raw()
	}
	return rawTypes
}

func (b Block) Raw() C.MlirBlock           { return C.MlirBlock(b) }
func (b Block) Ptr() unsafe.Pointer        { return unsafe.Pointer(b.Raw().ptr) }
func (b Block) IsNull() bool               { return bool(C.mlirBlockIsNull(b.Raw())) }
func (b Block) Destroy()                   { C.mlirBlockDestroy(b.Raw()) }
func (b Block) Equal(other Block) bool     { return bool(C.mlirBlockEqual(b.Raw(), other.Raw())) }
func (b Block) Detach()                    { C.mlirBlockDetach(b.Raw()) }
func (b Block) ParentOperation() Operation { return Operation(C.mlirBlockGetParentOperation(b.Raw())) }
func (b Block) ParentRegion() Region       { return Region(C.mlirBlockGetParentRegion(b.Raw())) }
func (b Block) NextInRegion() Block        { return Block(C.mlirBlockGetNextInRegion(b.Raw())) }
func (b Block) FirstOperation() Operation  { return Operation(C.mlirBlockGetFirstOperation(b.Raw())) }
func (b Block) Terminator() Operation      { return Operation(C.mlirBlockGetTerminator(b.Raw())) }

func (b Block) AppendOwnedOperation(operation Operation) {
	C.mlirBlockAppendOwnedOperation(b.Raw(), operation.Raw())
}

func (b Block) InsertOwnedOperation(pos int, operation Operation) {
	C.mlirBlockInsertOwnedOperation(b.Raw(), C.intptr_t(pos), operation.Raw())
}

func (b Block) InsertOwnedOperationAfter(reference, operation Operation) {
	C.mlirBlockInsertOwnedOperationAfter(b.Raw(), reference.Raw(), operation.Raw())
}

func (b Block) InsertOwnedOperationBefore(reference, operation Operation) {
	C.mlirBlockInsertOwnedOperationBefore(b.Raw(), reference.Raw(), operation.Raw())
}

func (b Block) NumArguments() int { return int(C.mlirBlockGetNumArguments(b.Raw())) }

func (b Block) AddArgument(ty TypeLike, loc LocationLike) {
	C.mlirBlockAddArgument(b.Raw(), ty.Raw(), loc.Raw())
}

func (b Block) EraseArgument(index int) {
	C.mlirBlockEraseArgument(b.Raw(), C.unsigned(index))
}

func (b Block) InsertArgument(index int, ty TypeLike, loc LocationLike) {
	C.mlirBlockInsertArgument(b.Raw(), C.intptr_t(index), ty.Raw(), loc.Raw())
}

func (b Block) Argument(pos int) BlockArg {
	return BlockArg{Value: Value(C.mlirBlockGetArgument(b.Raw(), C.intptr_t(pos)))}
}

func (b Block) String() string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirBlockPrint(b.Raw(), cb, ud)
	})
}

func (b Block) NumSuccessors() int {
	return int(C.mlirBlockGetNumSuccessors(b.Raw()))
}

func (b Block) Successor(pos int) Block {
	return Block(C.mlirBlockGetSuccessor(b.Raw(), C.intptr_t(pos)))
}

func (b Block) NumPredecessors() int {
	return int(C.mlirBlockGetNumPredecessors(b.Raw()))
}

func (b Block) Predecessor(pos int) Block {
	return Block(C.mlirBlockGetPredecessor(b.Raw(), C.intptr_t(pos)))
}

//===----------------------------------------------------------------------===//
// Value API.
//===----------------------------------------------------------------------===//

type ValueLike interface {
	Raw() C.MlirValue
	Ptr() unsafe.Pointer
	IsNull() bool
	Context() Context
	Location() Location
	Type() TypeLike
	AsValue() Value
}

func UnwrapValue(value ValueLike) C.MlirValue {
	if value == nil || value.IsNull() {
		return C.MlirValue{}
	}
	return value.Raw()
}

type Value C.MlirValue

func UnwrapValueSlice[T ValueLike](values []T) []C.MlirValue {
	rawAttrs := make([]C.MlirValue, len(values))
	for i, v := range values {
		rawAttrs[i] = v.Raw()
	}
	return rawAttrs
}

func (v Value) Raw() C.MlirValue       { return C.MlirValue(v) }
func (v Value) Ptr() unsafe.Pointer    { return unsafe.Pointer(v.Raw().ptr) }
func (v Value) IsNull() bool           { return bool(C.mlirValueIsNull(v.Raw())) }
func (v Value) Context() Context       { return Context(C.mlirValueGetContext(v.Raw())) }
func (v Value) Location() Location     { return wrapLocation(C.mlirValueGetLocation(v.Raw())) }
func (v Value) IsABlockArgument() bool { return bool(C.mlirValueIsABlockArgument(v.Raw())) }
func (v Value) IsAOpResult() bool      { return bool(C.mlirValueIsAOpResult(v.Raw())) }
func (v Value) Type() TypeLike         { return WrapType(C.mlirValueGetType(v.Raw())) }
func (v Value) SetType(ty TypeLike)    { C.mlirValueSetType(v.Raw(), ty.Raw()) }
func (v Value) Dump()                  { C.mlirValueDump(v.Raw()) }

func (v Value) String() string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirValuePrint(v.Raw(), cb, ud)
	})
}

func (v Value) OperandString(state AsmState) string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirValuePrintAsOperand(v.Raw(), state.Raw(), cb, ud)
	})
}

func (v Value) FirstUse() OpOperand           { return OpOperand(C.mlirValueGetFirstUse(v.Raw())) }
func (v Value) ReplaceAllUsesWith(with Value) { C.mlirValueReplaceAllUsesOfWith(v.Raw(), with.Raw()) }
func (v Value) ReplaceAllUsesExcept(with Value, except []Operation) {
	var array *C.MlirOperation
	if len(except) > 0 {
		array = (*C.MlirOperation)(&except[0])
	}
	C.mlirValueReplaceAllUsesExcept(v.Raw(), with.Raw(), C.intptr_t(len(except)), array)
}

func (v Value) AsValue() Value { return v }
func (v Value) AsBlockArg() (BlockArg, bool) {
	if !v.IsABlockArgument() {
		return BlockArg{}, false
	}
	return BlockArg{Value: v}, true
}

func (v Value) AsResult() (Result, bool) {
	if !v.IsAOpResult() {
		return Result{}, false
	}
	return Result{Value: v}, true
}

type BlockArg struct {
	Value
}

func (v BlockArg) OwningBlock() Block          { return Block(C.mlirBlockArgumentGetOwner(v.Raw())) }
func (v BlockArg) BlockArgNumber() int         { return int(C.mlirBlockArgumentGetArgNumber(v.Raw())) }
func (v BlockArg) SetBlockArgType(ty TypeLike) { C.mlirBlockArgumentSetType(v.Raw(), ty.Raw()) }
func (v BlockArg) SetBlockArgLocation(loc LocationLike) {
	C.mlirBlockArgumentSetLocation(v.Raw(), loc.Raw())
}

type Result struct {
	Value
}

func (v Result) OwningOperation() Operation { return Operation(C.mlirOpResultGetOwner(v.Raw())) }
func (v Result) ResultNumber() int {
	return int(C.mlirOpResultGetResultNumber(v.Raw()))
}

//===----------------------------------------------------------------------===//
// OpOperand API.
//===----------------------------------------------------------------------===//

type OpOperand C.MlirOpOperand

func (o OpOperand) Raw() C.MlirOpOperand { return C.MlirOpOperand(o) }
func (o OpOperand) IsNull() bool         { return bool(C.mlirOpOperandIsNull(o.Raw())) }
func (o OpOperand) Value() Value         { return Value(C.mlirOpOperandGetValue(o.Raw())) }
func (o OpOperand) Owner() Operation     { return Operation(C.mlirOpOperandGetOwner(o.Raw())) }
func (o OpOperand) OperandNumber() int   { return int(C.mlirOpOperandGetOperandNumber(o.Raw())) }
func (o OpOperand) NextUse() OpOperand   { return OpOperand(C.mlirOpOperandGetNextUse(o.Raw())) }

//===----------------------------------------------------------------------===//
// Type API.
//===----------------------------------------------------------------------===//

type TypeLike interface {
	Raw() C.MlirType
	Ptr() unsafe.Pointer
	IsNull() bool
	Context() Context
	Equal(other TypeLike) bool
	TypeId() TypeId
	Dialect() Dialect
	Dump()
	ToType() Type
	String() string
}

type baseType C.MlirType

func (t baseType) Raw() C.MlirType           { return C.MlirType(t) }
func (t baseType) Ptr() unsafe.Pointer       { return unsafe.Pointer(t.Raw().ptr) }
func (t baseType) IsNull() bool              { return bool(C.mlirTypeIsNull(t.Raw())) }
func (t baseType) Context() Context          { return Context(C.mlirTypeGetContext(t.Raw())) }
func (t baseType) Equal(other TypeLike) bool { return bool(C.mlirTypeEqual(t.Raw(), other.Raw())) }
func (t baseType) TypeId() TypeId            { return TypeId(C.mlirTypeGetTypeID(t.Raw())) }
func (t baseType) Dialect() Dialect          { return Dialect(C.mlirTypeGetDialect(t.Raw())) }
func (t baseType) Dump()                     { C.mlirTypeDump(t.Raw()) }
func (t baseType) ToType() Type              { return Type{t} }
func (t baseType) String() string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirTypePrint(t.Raw(), cb, ud)
	})
}

type Type struct {
	baseType
}

func WrapType(raw C.MlirType) Type {
	return Type{baseType: baseType(raw)}
}

func WrapExternalType(pointer unsafe.Pointer) Type {
	return Type{baseType: baseType(C.MlirType{ptr: pointer})}
}

func UnwrapType(typ TypeLike) C.MlirType {
	if typ == nil || typ.IsNull() {
		return C.MlirType{}
	}
	return typ.Raw()
}

func UnwrapTypeSlice[T TypeLike](types []T) []C.MlirType {
	rawTypes := make([]C.MlirType, len(types))
	for i, t := range types {
		rawTypes[i] = t.Raw()
	}
	return rawTypes
}

func NewTypeFromString(ctx Context, input string) Type {
	refInput := NewStringRef(input)
	defer refInput.Destroy()
	return WrapType(C.mlirTypeParseGet(ctx.Raw(), refInput.Raw()))
}

//===----------------------------------------------------------------------===//
// Attribute API.
//===----------------------------------------------------------------------===//

// NamedAttribute
//
// A named attribute is essentially a (name, attribute) pair where the name is
// a string.
type NamedAttribute C.MlirNamedAttribute

func NewNamedAttribute(name string, attribute AttributeLike) NamedAttribute {
	ident := NewIdentifier(attribute.Context(), name)
	return NamedAttribute(
		C.mlirNamedAttributeGet(ident.Raw(), attribute.Raw()))
}

func (a NamedAttribute) Raw() C.MlirNamedAttribute {
	return C.MlirNamedAttribute(a)
}

type AttributeLike interface {
	Raw() C.MlirAttribute
	Ptr() unsafe.Pointer
	IsNull() bool
	Context() Context
	TypeId() TypeId
	Dialect() Dialect
	Equal(other AttributeLike) bool
	String() string
	AsAttribute() Attribute
}

func UnwrapAttribute(attr AttributeLike) C.MlirAttribute {
	if attr == nil || attr.IsNull() {
		return NewNullAttribute().Raw()
	}
	return attr.Raw()
}

type baseAttribute C.MlirAttribute

func (a baseAttribute) Raw() C.MlirAttribute   { return C.MlirAttribute(a) }
func (a baseAttribute) Ptr() unsafe.Pointer    { return unsafe.Pointer(a.Raw().ptr) }
func (a baseAttribute) IsNull() bool           { return bool(C.mlirAttributeIsNull(a.Raw())) }
func (a baseAttribute) Context() Context       { return Context(C.mlirAttributeGetContext(a.Raw())) }
func (a baseAttribute) TypeId() TypeId         { return TypeId(C.mlirAttributeGetTypeID(a.Raw())) }
func (a baseAttribute) Dialect() Dialect       { return Dialect(C.mlirAttributeGetDialect(a.Raw())) }
func (a baseAttribute) Dump()                  { C.mlirAttributeDump(a.Raw()) }
func (a baseAttribute) AsAttribute() Attribute { return Attribute{a} }
func (a baseAttribute) Equal(other AttributeLike) bool {
	return bool(C.mlirAttributeEqual(a.Raw(), other.Raw()))
}

func (a baseAttribute) String() string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirAttributePrint(a.Raw(), cb, ud)
	})
}

type Attribute struct {
	baseAttribute
}

func (a Attribute) AsAttribute() Attribute { return a }

func WrapAttribute(raw C.MlirAttribute) Attribute {
	return Attribute{baseAttribute: baseAttribute(raw)}
}

func WrapExternalAttribute(pointer unsafe.Pointer) Attribute {
	return Attribute{baseAttribute: baseAttribute(C.MlirAttribute{ptr: pointer})}
}

func UnwrapAttributeSlice[T AttributeLike](attrs []T) []C.MlirAttribute {
	rawAttrs := make([]C.MlirAttribute, len(attrs))
	for i, a := range attrs {
		rawAttrs[i] = a.Raw()
	}
	return rawAttrs
}

func UnwrapNamedAttributeSlice(attrs []NamedAttribute) []C.MlirNamedAttribute {
	rawAttrs := make([]C.MlirNamedAttribute, len(attrs))
	for i, a := range attrs {
		rawAttrs[i] = a.Raw()
	}
	return rawAttrs
}

func NewAttributeFromString(ctx Context, input string) Attribute {
	refInput := NewStringRef(input)
	defer refInput.Destroy()
	return WrapAttribute(C.mlirAttributeParseGet(ctx.Raw(), refInput.Raw()))
}

//===----------------------------------------------------------------------===//
// Identifier API.
//===----------------------------------------------------------------------===//

type Identifier C.MlirIdentifier

func NewIdentifier(ctx Context, name string) Identifier {
	refName := NewStringRef(name)
	defer refName.Destroy()
	return Identifier(C.mlirIdentifierGet(ctx.Raw(), refName.Raw()))
}

func (i Identifier) Raw() C.MlirIdentifier {
	return C.MlirIdentifier(i)
}

func (i Identifier) Context() Context {
	return Context(C.mlirIdentifierGetContext(i.Raw()))
}

func (i Identifier) Equal(other Identifier) bool {
	return bool(C.mlirIdentifierEqual(i.Raw(), other.Raw()))
}

func (i Identifier) String() string {
	ref := BorrowedStringRef(C.mlirIdentifierStr(i.Raw()))
	return ref.String()
}

//===----------------------------------------------------------------------===//
// Symbol and SymbolTable API.
//===----------------------------------------------------------------------===//

func SymbolAttributeName() string {
	ref := BorrowedStringRef(C.mlirSymbolTableGetSymbolAttributeName())
	return ref.String()
}

func SymbolVisibilityAttributeName() string {
	ref := BorrowedStringRef(C.mlirSymbolTableGetVisibilityAttributeName())
	return ref.String()
}

func ReplaceAllSymbolUses(oldSymbol, newSymbol string, from Operation) LogicalResult {
	refOldSymbol := NewStringRef(oldSymbol)
	defer refOldSymbol.Destroy()

	refNewSymbol := NewStringRef(newSymbol)
	defer refNewSymbol.Destroy()

	return LogicalResult(C.mlirSymbolTableReplaceAllSymbolUses(refOldSymbol.Raw(), refNewSymbol.Raw(), from.Raw()))
}

type SymbolTable C.MlirSymbolTable

func NewSymbolTable(operation Operation) SymbolTable {
	return SymbolTable(C.mlirSymbolTableCreate(operation.Raw()))
}

func (s SymbolTable) Raw() C.MlirSymbolTable { return C.MlirSymbolTable(s) }
func (s SymbolTable) Destroy()               { C.mlirSymbolTableDestroy(s.Raw()) }
func (s SymbolTable) IsNull() bool           { return bool(C.mlirSymbolTableIsNull(s.Raw())) }
func (s SymbolTable) Lookup(name string) Operation {
	refName := NewStringRef(name)
	defer refName.Destroy()
	return Operation(C.mlirSymbolTableLookup(s.Raw(), refName.Raw()))
}

func (s SymbolTable) Insert(op Operation) Attribute {
	return WrapAttribute(C.mlirSymbolTableInsert(s.Raw(), op.Raw()))
}

func (s SymbolTable) Erase(op Operation) {
	C.mlirSymbolTableErase(s.Raw(), op.Raw())
}
