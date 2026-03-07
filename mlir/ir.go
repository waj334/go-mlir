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
		return Context(C.mlirContextCreateWithRegistry(config.Registry.raw(), C.bool(!config.DisableThreading)))
	} else if config.DisableThreading {
		return Context(C.mlirContextCreateWithThreading(C.bool(!config.DisableThreading)))
	}
	return Context(C.mlirContextCreate())
}

func (c Context) raw() C.MlirContext {
	return C.MlirContext(c)
}

func (c Context) IsNull() bool {
	return bool(C.mlirContextIsNull(c.raw()))
}

// Destroy takes an MLIR context owned by the caller and destroys it.
func (c Context) Destroy() {
	C.mlirContextDestroy(c.raw())
}

// Equal checks if two contexts are equal.
func (c Context) Equal(other Context) bool {
	return bool(C.mlirContextEqual(c.raw(), other.raw()))
}

// SetAllowUnregisteredDialects sets whether unregistered dialects are allowed in this context.
func (c Context) SetAllowUnregisteredDialects(allow bool) {
	C.mlirContextSetAllowUnregisteredDialects(c.raw(), C.bool(allow))
}

// UnregisteredDialectsAllowed returns whether the context allows unregistered dialects.
func (c Context) UnregisteredDialectsAllowed() bool {
	return bool(C.mlirContextGetAllowUnregisteredDialects(c.raw()))
}

// NumRegisteredDialects returns the number of dialects registered with the given context. A
// registered dialect will be loaded if needed by the parser.
func (c Context) NumRegisteredDialects() int {
	return int(C.mlirContextGetNumRegisteredDialects(c.raw()))
}

// AppendDialectRegistry append the contents of the given dialect registry to the registry associated
// with the context.
func (c Context) AppendDialectRegistry(registry DialectRegistry) {
	C.mlirContextAppendDialectRegistry(c.raw(), registry.raw())
}

// NumLoadedDialects returns the number of dialects loaded by the context.
func (c Context) NumLoadedDialects() int {
	return int(C.mlirContextGetNumLoadedDialects(c.raw()))
}

// GetOrLoadDialect gets the dialect instance owned by the given context using the dialect
// namespace to identify it, loads (i.e., constructs the instance of) the
// dialect if necessary. If the dialect is not registered with the context,
// returns null. Use mlirContextLoad<Name>Dialect to load an unregistered
// dialect.
func (c Context) GetOrLoadDialect(name string) Dialect {
	refName := NewStringRef(name)
	defer refName.Destroy()
	return Dialect(C.mlirContextGetOrLoadDialect(c.raw(), refName.raw()))
}

// EnableMultithreading sets the multithreading mode (must be set to false if using mlir-print-ir-after-all).
func (c Context) EnableMultithreading(enable bool) {
	C.mlirContextEnableMultithreading(c.raw(), C.bool(enable))
}

// LoadAllAvailableDialects Eagerly loads all available dialects registered with a context, making
// them available for use for IR construction.
func (c Context) LoadAllAvailableDialects() {
	C.mlirContextLoadAllAvailableDialects(c.raw())
}

func (c Context) SetThreadPool(pool LLVMThreadPool) {
	C.mlirContextSetThreadPool(c.raw(), pool.raw())
}

func (c Context) NumThreads() int {
	return int(C.mlirContextGetNumThreads(c.raw()))
}

func (c Context) ThreadPool() LLVMThreadPool {
	return LLVMThreadPool(C.mlirContextGetThreadPool(c.raw()))
}

//===----------------------------------------------------------------------===//
// Dialect API.
//===----------------------------------------------------------------------===//

type Dialect C.MlirDialect

func (d Dialect) raw() C.MlirDialect {
	return C.MlirDialect(d)
}

func (d Dialect) IsNull() bool {
	return bool(C.mlirDialectIsNull(d.raw()))
}

func (d Dialect) Equal(other Dialect) bool {
	return bool(C.mlirDialectEqual(d.raw(), other.raw()))
}

func (d Dialect) Namespace() string {
	ref := BorrowedStringRef(C.mlirDialectGetNamespace(d.raw()))
	return ref.String()
}

//===----------------------------------------------------------------------===//
// DialectHandle API.
//===----------------------------------------------------------------------===//

type DialectHandle C.MlirDialectHandle

func (d DialectHandle) raw() C.MlirDialectHandle {
	return C.MlirDialectHandle(d)
}

func (d DialectHandle) Namespace() string {
	ref := BorrowedStringRef(C.mlirDialectHandleGetNamespace(d.raw()))
	return ref.String()
}

func (d DialectHandle) InsertDialect(registry DialectRegistry) {
	C.mlirDialectHandleInsertDialect(d.raw(), registry.raw())
}

func (d DialectHandle) RegisterDialect(context Context) {
	C.mlirDialectHandleRegisterDialect(d.raw(), context.raw())
}

func (d DialectHandle) LoadDialect(context Context) {
	C.mlirDialectHandleLoadDialect(d.raw(), context.raw())
}

//===----------------------------------------------------------------------===//
// DialectRegistry API.
//===----------------------------------------------------------------------===//

type DialectRegistry C.MlirDialectRegistry

func NewDialectRegistry() DialectRegistry {
	return DialectRegistry(C.mlirDialectRegistryCreate())
}

func (d DialectRegistry) raw() C.MlirDialectRegistry {
	return C.MlirDialectRegistry(d)
}

func (d DialectRegistry) IsNull() bool {
	return bool(C.mlirDialectRegistryIsNull(d.raw()))
}

func (d DialectRegistry) Destroy() {
	C.mlirDialectRegistryDestroy(d.raw())
}

//===----------------------------------------------------------------------===//
// Location API.
//===----------------------------------------------------------------------===//

type LocationLike interface {
	raw() C.MlirLocation
	IsNull() bool
	Attribute() Attribute
	Context() Context
	Equal(other LocationLike) bool
	String() string
}

type baseLocation C.MlirLocation

func (l baseLocation) raw() C.MlirLocation  { return C.MlirLocation(l) }
func (l baseLocation) IsNull() bool         { return bool(C.mlirLocationIsNull(l.raw())) }
func (l baseLocation) Attribute() Attribute { return Attribute(C.mlirLocationGetAttribute(l.raw())) }
func (l baseLocation) Context() Context     { return Context(C.mlirLocationGetContext(l.raw())) }
func (l baseLocation) Equal(other LocationLike) bool {
	return bool(C.mlirLocationEqual(l.raw(), other.raw()))
}

func (l baseLocation) String() string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirLocationPrint(l.raw(), cb, ud)
	})
}

type Location struct {
	baseLocation
}

func wrapLocation(v C.MlirLocation) Location {
	return Location{baseLocation: baseLocation(v)}
}

func (l Location) IsAFileLineColRange() bool { return bool(C.mlirLocationIsAFileLineColRange(l.raw())) }
func (l Location) IsACallSite() bool         { return bool(C.mlirLocationIsACallSite(l.raw())) }
func (l Location) AsFileLineColRange() (FileLineColRangeLoc, bool) {
	if bool(C.mlirLocationIsAFileLineColRange(l.raw())) {
		return FileLineColRangeLoc{baseLocation: l.baseLocation}, true
	}
	return FileLineColRangeLoc{}, false
}

func (l Location) AsName() (NameLoc, bool) {
	if bool(C.mlirLocationIsAName(l.raw())) {
		return NameLoc{baseLocation: l.baseLocation}, true
	}
	return NameLoc{}, false
}

func (l Location) AsCallSite() (CallSiteLoc, bool) {
	if bool(C.mlirLocationIsACallSite(l.raw())) {
		return CallSiteLoc{baseLocation: l.baseLocation}, true
	}
	return CallSiteLoc{}, false
}

type FileLineColRangeLoc struct {
	baseLocation
}

func NewFileLineColRange(ctx Context, filename string, startLine, startColumn, endLine, endColumn int) FileLineColRangeLoc {
	refFilename := NewStringRef(filename)
	defer refFilename.Destroy()
	return FileLineColRangeLoc{
		baseLocation: baseLocation(C.mlirLocationFileLineColRangeGet(
			ctx.raw(),
			refFilename.raw(),
			C.unsigned(startLine),
			C.unsigned(startColumn),
			C.unsigned(endLine),
			C.unsigned(endColumn))),
	}
}

func (l FileLineColRangeLoc) TypeId() TypeId {
	return TypeId(C.mlirLocationFileLineColRangeGetTypeID())
}

func (l FileLineColRangeLoc) Filename() Identifier {
	return Identifier(C.mlirLocationFileLineColRangeGetFilename(l.raw()))
}

func (l FileLineColRangeLoc) StartLine() int {
	return int(C.mlirLocationFileLineColRangeGetStartLine(l.raw()))
}

func (l FileLineColRangeLoc) StartColumn() int {
	return int(C.mlirLocationFileLineColRangeGetStartColumn(l.raw()))
}

func (l FileLineColRangeLoc) EndLine() int {
	return int(C.mlirLocationFileLineColRangeGetEndLine(l.raw()))
}

func (l FileLineColRangeLoc) EndColumn() int {
	return int(C.mlirLocationFileLineColRangeGetEndColumn(l.raw()))
}

type CallSiteLoc struct {
	baseLocation
}

func NewCallSite(callee, caller Location) CallSiteLoc {
	return CallSiteLoc{
		baseLocation: baseLocation(C.mlirLocationCallSiteGet(callee.raw(), caller.raw())),
	}
}

func (l CallSiteLoc) TypeId() TypeId   { return TypeId(C.mlirLocationCallSiteGetTypeID()) }
func (l CallSiteLoc) Callee() Location { return wrapLocation(C.mlirLocationCallSiteGetCallee(l.raw())) }
func (l CallSiteLoc) Caller() Location { return wrapLocation(C.mlirLocationCallSiteGetCaller(l.raw())) }

type FusedLoc struct {
	baseLocation
}

func NewFusedLoc(ctx Context, locations []LocationLike, metadata Attribute) FusedLoc {
	raw := make([]C.MlirLocation, len(locations))
	for i, loc := range locations {
		raw[i] = loc.raw()
	}

	var ptr *C.MlirLocation
	if len(raw) != 0 {
		ptr = &raw[0]
	}

	return FusedLoc{
		baseLocation: baseLocation(C.mlirLocationFusedGet(
			ctx.raw(),
			C.intptr_t(len(raw)),
			ptr,
			metadata.raw())),
	}
}

func (l FusedLoc) TypeId() TypeId    { return TypeId(C.mlirLocationFusedGetTypeID()) }
func (l FusedLoc) NumLocations() int { return int(C.mlirLocationFusedGetNumLocations(l.raw())) }
func (l FusedLoc) Locations() []Location {
	n := l.NumLocations()
	if n == 0 {
		return nil
	}

	raw := make([]C.MlirLocation, n)
	C.mlirLocationFusedGetLocations(l.raw(), &raw[0])

	out := make([]Location, n)
	for i := range raw {
		out[i] = Location{baseLocation: baseLocation(raw[i])}
	}
	return out
}

func (l FusedLoc) Metadata() Attribute { return Attribute(C.mlirLocationFusedGetMetadata(l.raw())) }

type NameLoc struct {
	baseLocation
}

func NewName(ctx Context, name string, childLoc Location) NameLoc {
	refName := NewStringRef(name)
	defer refName.Destroy()
	return NameLoc{
		baseLocation: baseLocation(C.mlirLocationNameGet(ctx.raw(), refName.raw(), childLoc.raw())),
	}
}

func (l NameLoc) TypeId() TypeId     { return TypeId(C.mlirLocationNameGetTypeID()) }
func (l NameLoc) Name() Identifier   { return Identifier(C.mlirLocationNameGetName(l.raw())) }
func (l NameLoc) ChildLoc() Location { return wrapLocation(C.mlirLocationNameGetChildLoc(l.raw())) }

type UnknownLoc struct {
	baseLocation
}

func NewUnknownLoc(ctx Context) UnknownLoc {
	return UnknownLoc{
		baseLocation: baseLocation(C.mlirLocationUnknownGet(ctx.raw())),
	}
}

//===----------------------------------------------------------------------===//
// Module API.
//===----------------------------------------------------------------------===//

type Module C.MlirModule

func NewModule(loc LocationLike) Module {
	return Module(C.mlirModuleCreateEmpty(loc.raw()))
}

func NewModuleFromString(ctx Context, input string) Module {
	refInput := NewStringRef(input)
	defer refInput.Destroy()
	return Module(C.mlirModuleCreateParse(ctx.raw(), refInput.raw()))
}

func NewModuleFromFile(ctx Context, filename string) Module {
	refFilename := NewStringRef(filename)
	defer refFilename.Destroy()
	return Module(C.mlirModuleCreateParseFromFile(ctx.raw(), refFilename.raw()))
}

func NewModuleFromOperation(op Operation) Module {
	return Module(C.mlirModuleFromOperation(op.raw()))
}

func (m Module) raw() C.MlirModule       { return C.MlirModule(m) }
func (m Module) IsNull() bool            { return bool(C.mlirModuleIsNull(m.raw())) }
func (m Module) Destroy()                { C.mlirModuleDestroy(m.raw()) }
func (m Module) Context() Context        { return Context(C.mlirModuleGetContext(m.raw())) }
func (m Module) Body() Block             { return Block(C.mlirModuleGetBody(m.raw())) }
func (m Module) Operation() Operation    { return Operation(C.mlirModuleGetOperation(m.raw())) }
func (m Module) Equal(other Module) bool { return bool(C.mlirModuleEqual(m.raw(), other.raw())) }
func (m Module) HashValue() uintptr      { return uintptr(C.mlirModuleHashValue(m.raw())) }

//===----------------------------------------------------------------------===//
// Operation state.
//===----------------------------------------------------------------------===//

type OperationState C.MlirOperationState

func NewOperationState(name string, loc LocationLike) *OperationState {
	refName := NewStringRef(name)
	s := OperationState(C.mlirOperationStateGet(refName.raw(), loc.raw()))
	return &s
}

func (s *OperationState) raw() *C.MlirOperationState { return (*C.MlirOperationState)(s) }
func (s *OperationState) AddResults(results ...TypeLike) *OperationState {
	if len(results) > 0 {
		rawResults := unwrapTypeSlice(results)
		C.mlirOperationStateAddResults(s.raw(), C.intptr_t(len(results)), (*C.MlirType)(&rawResults[0]))
	}
	return s
}

func (s *OperationState) AddOperands(operands ...Value) *OperationState {
	if len(operands) > 0 {
		C.mlirOperationStateAddOperands(s.raw(), C.intptr_t(len(operands)), (*C.MlirValue)(&operands[0]))
	}
	return s
}

func (s *OperationState) AddOwnedRegions(regions ...Region) *OperationState {
	if len(regions) > 0 {
		C.mlirOperationStateAddOwnedRegions(s.raw(), C.intptr_t(len(regions)), (*C.MlirRegion)(&regions[0]))
	}
	return s
}

func (s *OperationState) AddSuccessors(successors ...Block) *OperationState {
	if len(successors) > 0 {
		C.mlirOperationStateAddSuccessors(s.raw(), C.intptr_t(len(successors)), (*C.MlirBlock)(&successors[0]))
	}
	return s
}

func (s *OperationState) AddAttributes(attributes ...NamedAttribute) *OperationState {
	if len(attributes) > 0 {
		C.mlirOperationStateAddAttributes(s.raw(), C.intptr_t(len(attributes)), (*C.MlirNamedAttribute)(&attributes[0]))
	}
	return s
}

func (s *OperationState) EnableResultTypeInference() *OperationState {
	C.mlirOperationStateEnableResultTypeInference(s.raw())
	return s
}

func (s *OperationState) Create() Operation {
	return Operation(C.mlirOperationCreate(s.raw()))
}

//===----------------------------------------------------------------------===//
// AsmState API.
//===----------------------------------------------------------------------===//

type AsmState C.MlirAsmState

func NewAsmStateFromOperation(op Operation, flags OpPrintingFlags) AsmState {
	return AsmState(C.mlirAsmStateCreateForOperation(op.raw(), flags.raw()))
}

func NewAsmStateFromValue(value Value, flags OpPrintingFlags) AsmState {
	return AsmState(C.mlirAsmStateCreateForValue(value.raw(), flags.raw()))
}

func (s AsmState) raw() C.MlirAsmState { return C.MlirAsmState(s) }
func (s AsmState) Destroy()            { C.mlirAsmStateDestroy(s.raw()) }

//===----------------------------------------------------------------------===//
// Op Printing flags API.
//===----------------------------------------------------------------------===//

type OpPrintingFlags C.MlirOpPrintingFlags

func NewOpPrintingFlags() OpPrintingFlags {
	return OpPrintingFlags(C.mlirOpPrintingFlagsCreate())
}

func (f OpPrintingFlags) raw() C.MlirOpPrintingFlags { return C.MlirOpPrintingFlags(f) }
func (f OpPrintingFlags) Destroy()                   { C.mlirOpPrintingFlagsDestroy(f.raw()) }

func (f OpPrintingFlags) WithElidLargeElementsAttrs(limit int) OpPrintingFlags {
	C.mlirOpPrintingFlagsElideLargeElementsAttrs(f.raw(), C.intptr_t(limit))
	return f
}

func (f OpPrintingFlags) WithElideLargeResourceString(limit int) OpPrintingFlags {
	C.mlirOpPrintingFlagsElideLargeResourceString(f.raw(), C.intptr_t(limit))
	return f
}

func (f OpPrintingFlags) WithEnableDebugInfo(enable bool, pretty bool) OpPrintingFlags {
	C.mlirOpPrintingFlagsEnableDebugInfo(f.raw(), C.bool(enable), C.bool(pretty))
	return f
}

func (f OpPrintingFlags) WithGenericOpForm() OpPrintingFlags {
	C.mlirOpPrintingFlagsPrintGenericOpForm(f.raw())
	return f
}

func (f OpPrintingFlags) WithPrintNameLocAsPrefix() OpPrintingFlags {
	C.mlirOpPrintingFlagsPrintNameLocAsPrefix(f.raw())
	return f
}

func (f OpPrintingFlags) WithUseLocalScope() OpPrintingFlags {
	C.mlirOpPrintingFlagsUseLocalScope(f.raw())
	return f
}

func (f OpPrintingFlags) WithAssumeVerified() OpPrintingFlags {
	C.mlirOpPrintingFlagsAssumeVerified(f.raw())
	return f
}

func (f OpPrintingFlags) WithSkipRegions() OpPrintingFlags {
	C.mlirOpPrintingFlagsSkipRegions(f.raw())
	return f
}

//===----------------------------------------------------------------------===//
// Bytecode printing flags API.
//===----------------------------------------------------------------------===//

type BytecodeWriterConfig C.MlirBytecodeWriterConfig

func NewBytecodeWriterConfig() BytecodeWriterConfig {
	return BytecodeWriterConfig(C.mlirBytecodeWriterConfigCreate())
}

func (c BytecodeWriterConfig) raw() C.MlirBytecodeWriterConfig { return C.MlirBytecodeWriterConfig(c) }

func (c BytecodeWriterConfig) SetDesiredEmitVersion(version int) {
	C.mlirBytecodeWriterConfigDesiredEmitVersion(c.raw(), C.int64_t(version))
}

//===----------------------------------------------------------------------===//
// Operation API.
//===----------------------------------------------------------------------===//

type Operation C.MlirOperation

func NewOperation(state *OperationState) Operation {
	return Operation(C.mlirOperationCreate(state.raw()))
}

func NewOperationFromString(ctx Context, source string, sourceName string) Operation {
	refSource := NewStringRef(source)
	defer refSource.Destroy()

	refSourceName := NewStringRef(sourceName)
	defer refSourceName.Destroy()

	return Operation(C.mlirOperationCreateParse(ctx.raw(), refSource.raw(), refSourceName.raw()))
}

func (o Operation) raw() C.MlirOperation         { return C.MlirOperation(o) }
func (o Operation) IsNull() bool                 { return bool(C.mlirOperationIsNull(o.raw())) }
func (o Operation) Context() Context             { return Context(C.mlirOperationGetContext(o.raw())) }
func (o Operation) Clone() Operation             { return Operation(C.mlirOperationClone(o.raw())) }
func (o Operation) Destroy()                     { C.mlirOperationDestroy(o.raw()) }
func (o Operation) RemoveFromParent()            { C.mlirOperationRemoveFromParent(o.raw()) }
func (o Operation) HashValue() uintptr           { return uintptr(C.mlirOperationHashValue(o.raw())) }
func (o Operation) Location() Location           { return wrapLocation(C.mlirOperationGetLocation(o.raw())) }
func (o Operation) SetLocation(loc LocationLike) { C.mlirOperationSetLocation(o.raw(), loc.raw()) }
func (o Operation) TypeId() TypeId               { return TypeId(C.mlirOperationGetTypeID(o.raw())) }
func (o Operation) Name() Identifier             { return Identifier(C.mlirOperationGetName(o.raw())) }
func (o Operation) Block() Block                 { return Block(C.mlirOperationGetBlock(o.raw())) }
func (o Operation) NumRegions() int              { return int(C.mlirOperationGetNumRegions(o.raw())) }
func (o Operation) NextInBlock() Operation       { return Operation(C.mlirOperationGetNextInBlock(o.raw())) }
func (o Operation) NumOperands() int             { return int(C.mlirOperationGetNumOperands(o.raw())) }
func (o Operation) NumResults() int              { return int(C.mlirOperationGetNumResults(o.raw())) }
func (o Operation) NumSuccessors() int           { return int(C.mlirOperationGetNumSuccessors(o.raw())) }
func (o Operation) NumAttributes() int           { return int(C.mlirOperationGetNumAttributes(o.raw())) }

func (o Operation) Attribute(pos int) NamedAttribute {
	return NamedAttribute(C.mlirOperationGetAttribute(o.raw(), C.intptr_t(pos)))
}

func (o Operation) AttributeByName(name string) Attribute {
	refName := NewStringRef(name)
	defer refName.Destroy()
	return Attribute(C.mlirOperationGetAttributeByName(o.raw(), refName.raw()))
}

func (o Operation) SetAttributeByName(name string, attr Attribute) {
	refName := NewStringRef(name)
	defer refName.Destroy()
	C.mlirOperationSetAttributeByName(o.raw(), refName.raw(), attr.raw())
}

func (o Operation) RemoveAttributeByName(name string) {
	refName := NewStringRef(name)
	defer refName.Destroy()
	C.mlirOperationRemoveAttributeByName(o.raw(), refName.raw())
}

func (o Operation) Successor(pos int) Block {
	return Block(C.mlirOperationGetSuccessor(o.raw(), C.intptr_t(pos)))
}

func (o Operation) SetSuccessor(pos int, block Block) {
	C.mlirOperationSetSuccessor(o.raw(), C.intptr_t(pos), block.raw())
}

func (o Operation) Result(pos int) Value {
	return Value(C.mlirOperationGetResult(o.raw(), C.intptr_t(pos)))
}

func (o Operation) Operand(pos int) Value {
	return Value(C.mlirOperationGetOperand(o.raw(), C.intptr_t(pos)))
}

func (o Operation) SetOperand(pos int, newValue Value) {
	C.mlirOperationSetOperand(o.raw(), C.intptr_t(pos), newValue.raw())
}

func (o Operation) SetOperands(operands ...Value) {
	if len(operands) > 0 {
		C.mlirOperationSetOperands(o.raw(), C.intptr_t(len(operands)), (*C.MlirValue)(&operands[0]))
	}
}

func (o Operation) Region(pos int) Region {
	return Region(C.mlirOperationGetRegion(o.raw(), C.intptr_t(pos)))
}

func (o Operation) ParentOperation() Operation {
	return Operation(C.mlirOperationGetParentOperation(o.raw()))
}

func (o Operation) Equal(other Operation) bool {
	return bool(C.mlirOperationEqual(o.raw(), other.raw()))
}

func (o Operation) HasInherentAttributeByName(name string) bool {
	refName := NewStringRef(name)
	defer refName.Destroy()
	return bool(C.mlirOperationHasInherentAttributeByName(o.raw(), refName.raw()))
}

func (o Operation) InherentAttributeByName(name string) Attribute {
	refName := NewStringRef(name)
	defer refName.Destroy()
	return Attribute(C.mlirOperationGetInherentAttributeByName(o.raw(), refName.raw()))
}

func (o Operation) SetInherentAttributeByName(name string, attr Attribute) {
	refName := NewStringRef(name)
	defer refName.Destroy()
	C.mlirOperationSetInherentAttributeByName(o.raw(), refName.raw(), attr.raw())
}

func (o Operation) NumDiscardableAttributes() int {
	return int(C.mlirOperationGetNumDiscardableAttributes(o.raw()))
}

func (o Operation) DiscardableAttribute(pos int) NamedAttribute {
	return NamedAttribute(C.mlirOperationGetDiscardableAttribute(o.raw(), C.intptr_t(pos)))
}

func (o Operation) DiscardableAttributeByName(name string) Attribute {
	refName := NewStringRef(name)
	defer refName.Destroy()
	return Attribute(C.mlirOperationGetDiscardableAttributeByName(o.raw(), refName.raw()))
}

func (o Operation) SetDiscardableAttributeByName(name string, attr Attribute) {
	refName := NewStringRef(name)
	defer refName.Destroy()
	C.mlirOperationSetDiscardableAttributeByName(o.raw(), refName.raw(), attr.raw())
}

func (o Operation) RemoveDiscardableAttributeByName(name string) {
	refName := NewStringRef(name)
	defer refName.Destroy()
	C.mlirOperationRemoveDiscardableAttributeByName(o.raw(), refName.raw())
}

func (o Operation) String() string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirOperationPrint(o.raw(), cb, ud)
	})
}

func (o Operation) StringWithFlags(flags OpPrintingFlags) string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirOperationPrintWithFlags(o.raw(), flags.raw(), cb, ud)
	})
}

func (o Operation) StringWithState(state AsmState) string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirOperationPrintWithState(o.raw(), state.raw(), cb, ud)
	})
}

func (o Operation) Bytecode() string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirOperationWriteBytecode(o.raw(), cb, ud)
	})
}

func (o Operation) BytecodeWithConfig(config BytecodeWriterConfig) string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirOperationWriteBytecodeWithConfig(o.raw(), config.raw(), cb, ud)
	})
}

func (o Operation) Dump()                      { C.mlirOperationDump(o.raw()) }
func (o Operation) Verify()                    { C.mlirOperationVerify(o.raw()) }
func (o Operation) MoveAfter(other Operation)  { C.mlirOperationMoveAfter(o.raw(), other.raw()) }
func (o Operation) MoveBefore(other Operation) { C.mlirOperationMoveBefore(o.raw(), other.raw()) }

func (o Operation) IsBeforeInBlock(other Operation) bool {
	return bool(C.mlirOperationIsBeforeInBlock(o.raw(), other.raw()))
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
	C.mlirOperationWalk(o.raw(), cb, userdata, C.MlirWalkOrder(walkOrder))
}

func (o Operation) ReplaceUsesOfWith(of, with Value) {
	C.mlirOperationReplaceUsesOfWith(o.raw(), of.raw(), with.raw())
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
	C.mlirSymbolTableWalkSymbolTables(o.raw(), C.bool(allSymbolUsesVisible), cb, userdata)
}

//===----------------------------------------------------------------------===//
// Region API.
//===----------------------------------------------------------------------===//

type Region C.MlirRegion

func NewRegion() Region {
	return Region(C.mlirRegionCreate())
}

func (r Region) raw() C.MlirRegion            { return C.MlirRegion(r) }
func (r Region) IsNull() bool                 { return bool(C.mlirRegionIsNull(r.raw())) }
func (r Region) Destroy()                     { C.mlirRegionDestroy(r.raw()) }
func (r Region) Equal(other Region) bool      { return bool(C.mlirRegionEqual(r.raw(), other.raw())) }
func (r Region) FirstBlock() Block            { return Block(C.mlirRegionGetFirstBlock(r.raw())) }
func (r Region) AppendOwnedBlock(block Block) { C.mlirRegionAppendOwnedBlock(r.raw(), block.raw()) }
func (r Region) NextInOperation() Region      { return Region(C.mlirRegionGetNextInOperation(r.raw())) }
func (r Region) TakeBody(source Region)       { C.mlirRegionTakeBody(r.raw(), source.raw()) }
func (r Region) InsertOwnedBlock(pos int, block Block) {
	C.mlirRegionInsertOwnedBlock(r.raw(), C.intptr_t(pos), block.raw())
}

func (r Region) InsertOwnedBlockAfter(reference, block Block) {
	C.mlirRegionInsertOwnedBlockAfter(r.raw(), reference.raw(), block.raw())
}

func (r Region) InsertOwnedBlockBefore(reference, block Block) {
	C.mlirRegionInsertOwnedBlockBefore(r.raw(), reference.raw(), block.raw())
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
		rawArgs := unwrapTypeSlice(args)

		rawLocs := make([]C.MlirLocation, len(locs))
		for i, loc := range locs {
			rawLocs[i] = loc.raw()
		}

		rawArgsPtr = (*C.MlirType)(&rawArgs[0])
		rawLocsPtr = (*C.MlirLocation)(&rawLocs[0])
	}

	return Block(C.mlirBlockCreate(C.intptr_t(len(args)), rawArgsPtr, rawLocsPtr))
}

func (b Block) raw() C.MlirBlock           { return C.MlirBlock(b) }
func (b Block) IsNull() bool               { return bool(C.mlirBlockIsNull(b.raw())) }
func (b Block) Destroy()                   { C.mlirBlockDestroy(b.raw()) }
func (b Block) Equal(other Block) bool     { return bool(C.mlirBlockEqual(b.raw(), other.raw())) }
func (b Block) Detach()                    { C.mlirBlockDetach(b.raw()) }
func (b Block) ParentOperation() Operation { return Operation(C.mlirBlockGetParentOperation(b.raw())) }
func (b Block) ParentRegion() Region       { return Region(C.mlirBlockGetParentRegion(b.raw())) }
func (b Block) NextInRegion() Block        { return Block(C.mlirBlockGetNextInRegion(b.raw())) }
func (b Block) FirstOperation() Operation  { return Operation(C.mlirBlockGetFirstOperation(b.raw())) }
func (b Block) Terminator() Operation      { return Operation(C.mlirBlockGetTerminator(b.raw())) }

func (b Block) AppendOwnedOperation(operation Operation) {
	C.mlirBlockAppendOwnedOperation(b.raw(), operation.raw())
}

func (b Block) InsertOwnedOperation(pos int, operation Operation) {
	C.mlirBlockInsertOwnedOperation(b.raw(), C.intptr_t(pos), operation.raw())
}

func (b Block) InsertOwnedOperationAfter(reference, operation Operation) {
	C.mlirBlockInsertOwnedOperationAfter(b.raw(), reference.raw(), operation.raw())
}

func (b Block) InsertOwnedOperationBefore(reference, operation Operation) {
	C.mlirBlockInsertOwnedOperationBefore(b.raw(), reference.raw(), operation.raw())
}

func (b Block) NumArguments() int { return int(C.mlirBlockGetNumArguments(b.raw())) }

func (b Block) AddArgument(ty TypeLike, loc LocationLike) {
	C.mlirBlockAddArgument(b.raw(), ty.raw(), loc.raw())
}

func (b Block) EraseArgument(index int) {
	C.mlirBlockEraseArgument(b.raw(), C.unsigned(index))
}

func (b Block) InsertArgument(index int, ty TypeLike, loc LocationLike) {
	C.mlirBlockInsertArgument(b.raw(), C.intptr_t(index), ty.raw(), loc.raw())
}

func (b Block) Argument(pos int) Value {
	return Value(C.mlirBlockGetArgument(b.raw(), C.intptr_t(pos)))
}

func (b Block) String() string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirBlockPrint(b.raw(), cb, ud)
	})
}

func (b Block) NumSuccessors() int {
	return int(C.mlirBlockGetNumSuccessors(b.raw()))
}

func (b Block) Successor(pos int) Block {
	return Block(C.mlirBlockGetSuccessor(b.raw(), C.intptr_t(pos)))
}

func (b Block) NumPredecessors() int {
	return int(C.mlirBlockGetNumPredecessors(b.raw()))
}

func (b Block) Predecessor(pos int) Block {
	return Block(C.mlirBlockGetPredecessor(b.raw(), C.intptr_t(pos)))
}

//===----------------------------------------------------------------------===//
// Value API.
//===----------------------------------------------------------------------===//

type Value C.MlirValue

func (v Value) raw() C.MlirValue            { return C.MlirValue(v) }
func (v Value) IsNull() bool                { return bool(C.mlirValueIsNull(v.raw())) }
func (v Value) Context() Context            { return Context(C.mlirValueGetContext(v.raw())) }
func (v Value) Location() Location          { return wrapLocation(C.mlirValueGetLocation(v.raw())) }
func (v Value) IsABlockArgument() bool      { return bool(C.mlirValueIsABlockArgument(v.raw())) }
func (v Value) IsAOpResult() bool           { return bool(C.mlirValueIsAOpResult(v.raw())) }
func (v Value) OwningBlock() Block          { return Block(C.mlirBlockArgumentGetOwner(v.raw())) }
func (v Value) BlockArgNumber() int         { return int(C.mlirBlockArgumentGetArgNumber(v.raw())) }
func (v Value) SetBlockArgType(ty TypeLike) { C.mlirBlockArgumentSetType(v.raw(), ty.raw()) }
func (v Value) SetBlockArgLocation(loc LocationLike) {
	C.mlirBlockArgumentSetLocation(v.raw(), loc.raw())
}

func (v Value) OwningOperation() Operation { return Operation(C.mlirOpResultGetOwner(v.raw())) }
func (v Value) ResultNumber() int          { return int(C.mlirOpResultGetResultNumber(v.raw())) }
func (v Value) Type() Type                 { return wrapType(C.mlirValueGetType(v.raw())) }
func (v Value) SetType(ty TypeLike)        { C.mlirValueSetType(v.raw(), ty.raw()) }
func (v Value) Dump()                      { C.mlirValueDump(v.raw()) }

func (v Value) String() string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirValuePrint(v.raw(), cb, ud)
	})
}

func (v Value) OperandString(state AsmState) string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirValuePrintAsOperand(v.raw(), state.raw(), cb, ud)
	})
}

func (v Value) FirstUse() OpOperand           { return OpOperand(C.mlirValueGetFirstUse(v.raw())) }
func (v Value) ReplaceAllUsesWith(with Value) { C.mlirValueReplaceAllUsesOfWith(v.raw(), with.raw()) }
func (v Value) ReplaceAllUsesExcept(with Value, except []Operation) {
	var array *C.MlirOperation
	if len(except) > 0 {
		array = (*C.MlirOperation)(&except[0])
	}
	C.mlirValueReplaceAllUsesExcept(v.raw(), with.raw(), C.intptr_t(len(except)), array)
}

//===----------------------------------------------------------------------===//
// OpOperand API.
//===----------------------------------------------------------------------===//

type OpOperand C.MlirOpOperand

func (o OpOperand) raw() C.MlirOpOperand { return C.MlirOpOperand(o) }
func (o OpOperand) IsNull() bool         { return bool(C.mlirOpOperandIsNull(o.raw())) }
func (o OpOperand) Value() Value         { return Value(C.mlirOpOperandGetValue(o.raw())) }
func (o OpOperand) Owner() Operation     { return Operation(C.mlirOpOperandGetOwner(o.raw())) }
func (o OpOperand) OperandNumber() int   { return int(C.mlirOpOperandGetOperandNumber(o.raw())) }
func (o OpOperand) NextUse() OpOperand   { return OpOperand(C.mlirOpOperandGetNextUse(o.raw())) }

//===----------------------------------------------------------------------===//
// Type API.
//===----------------------------------------------------------------------===//

type TypeLike interface {
	raw() C.MlirType
	IsNull() bool
	Context() Context
	Equal(other TypeLike) bool
	TypeId() TypeId
	Dialect() Dialect
	Dump()
	String() string
}

type baseType C.MlirType

func (t baseType) raw() C.MlirType           { return C.MlirType(t) }
func (t baseType) IsNull() bool              { return bool(C.mlirTypeIsNull(t.raw())) }
func (t baseType) Context() Context          { return Context(C.mlirTypeGetContext(t.raw())) }
func (t baseType) Equal(other TypeLike) bool { return bool(C.mlirTypeEqual(t.raw(), other.raw())) }
func (t baseType) TypeId() TypeId            { return TypeId(C.mlirTypeGetTypeID(t.raw())) }
func (t baseType) Dialect() Dialect          { return Dialect(C.mlirTypeGetDialect(t.raw())) }
func (t baseType) Dump()                     { C.mlirTypeDump(t.raw()) }
func (t baseType) String() string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirTypePrint(t.raw(), cb, ud)
	})
}

type Type struct {
	baseType
}

func wrapType(raw C.MlirType) Type {
	return Type{baseType: baseType(raw)}
}
func unwrapTypeSlice(types []TypeLike) []C.MlirType {
	rawTypes := make([]C.MlirType, len(types))
	for i, t := range types {
		rawTypes[i] = t.raw()
	}
	return rawTypes
}

func NewTypeFromString(ctx Context, input string) Type {
	refInput := NewStringRef(input)
	defer refInput.Destroy()
	return wrapType(C.mlirTypeParseGet(ctx.raw(), refInput.raw()))
}

//===----------------------------------------------------------------------===//
// Attribute API.
//===----------------------------------------------------------------------===//

// NamedAttribute
//
// A named attribute is essentially a (name, attribute) pair where the name is
// a string.
type NamedAttribute C.MlirNamedAttribute

func NewNamedAttribute(name string, attribute Attribute) NamedAttribute {
	ident := NewIdentifier(attribute.Context(), name)
	return NamedAttribute(
		C.mlirNamedAttributeGet(ident.raw(), attribute.raw()))
}

func (a NamedAttribute) raw() C.MlirNamedAttribute {
	return C.MlirNamedAttribute(a)
}

type Attribute C.MlirAttribute

func NewAttributeFromString(ctx Context, input string) Attribute {
	refInput := NewStringRef(input)
	defer refInput.Destroy()
	return Attribute(C.mlirAttributeParseGet(ctx.raw(), refInput.raw()))
}

func (a Attribute) raw() C.MlirAttribute { return C.MlirAttribute(a) }
func (a Attribute) IsNull() bool         { return bool(C.mlirAttributeIsNull(a.raw())) }
func (a Attribute) Context() Context     { return Context(C.mlirAttributeGetContext(a.raw())) }
func (a Attribute) TypeId() TypeId       { return TypeId(C.mlirAttributeGetTypeID(a.raw())) }
func (a Attribute) Dialect() Dialect     { return Dialect(C.mlirAttributeGetDialect(a.raw())) }
func (a Attribute) Dump()                { C.mlirAttributeDump(a.raw()) }
func (a Attribute) Equal(other Attribute) bool {
	return bool(C.mlirAttributeEqual(a.raw(), other.raw()))
}

func (a Attribute) String() string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirAttributePrint(a.raw(), cb, ud)
	})
}

//===----------------------------------------------------------------------===//
// Identifier API.
//===----------------------------------------------------------------------===//

type Identifier C.MlirIdentifier

func NewIdentifier(ctx Context, name string) Identifier {
	refName := NewStringRef(name)
	defer refName.Destroy()
	return Identifier(C.mlirIdentifierGet(ctx.raw(), refName.raw()))
}

func (i Identifier) raw() C.MlirIdentifier {
	return C.MlirIdentifier(i)
}

func (i Identifier) Context() Context {
	return Context(C.mlirIdentifierGetContext(i.raw()))
}

func (i Identifier) Equal(other Identifier) bool {
	return bool(C.mlirIdentifierEqual(i.raw(), other.raw()))
}

func (i Identifier) String() string {
	ref := BorrowedStringRef(C.mlirIdentifierStr(i.raw()))
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

	return LogicalResult(C.mlirSymbolTableReplaceAllSymbolUses(refOldSymbol.raw(), refNewSymbol.raw(), from.raw()))
}

type SymbolTable C.MlirSymbolTable

func NewSymbolTable(operation Operation) SymbolTable {
	return SymbolTable(C.mlirSymbolTableCreate(operation.raw()))
}

func (s SymbolTable) raw() C.MlirSymbolTable { return C.MlirSymbolTable(s) }
func (s SymbolTable) Destroy()               { C.mlirSymbolTableDestroy(s.raw()) }
func (s SymbolTable) IsNull() bool           { return bool(C.mlirSymbolTableIsNull(s.raw())) }
func (s SymbolTable) Lookup(name string) Operation {
	refName := NewStringRef(name)
	defer refName.Destroy()
	return Operation(C.mlirSymbolTableLookup(s.raw(), refName.raw()))
}

func (s SymbolTable) Insert(op Operation) Attribute {
	return Attribute(C.mlirSymbolTableInsert(s.raw(), op.raw()))
}

func (s SymbolTable) Erase(op Operation) {
	C.mlirSymbolTableErase(s.raw(), op.raw())
}
