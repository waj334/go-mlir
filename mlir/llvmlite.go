package mlir

/*
#include <stdlib.h>
#include <llvm-c/Analysis.h>
#include <llvm-c/Core.h>
#include <llvm-c/DebugInfo.h>
#include <llvm-c/Error.h>
#include <llvm-c/TargetMachine.h>
#include <llvm-c/Target.h>
#include <llvm-c/Transforms/PassBuilder.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

const (
	LLVMDWARFSourceLanguageC89 = iota + 1
	LLVMDWARFSourceLanguageC
	LLVMDWARFSourceLanguageAda83
	LLVMDWARFSourceLanguageCPlusPlus
	LLVMDWARFSourceLanguageCobol74
	LLVMDWARFSourceLanguageCobol85
	LLVMDWARFSourceLanguageFortran77
	LLVMDWARFSourceLanguageFortran90
	LLVMDWARFSourceLanguagePascal83
	LLVMDWARFSourceLanguageModula2
	LLVMDWARFSourceLanguageJava
	LLVMDWARFSourceLanguageC99
	LLVMDWARFSourceLanguageAda95
	LLVMDWARFSourceLanguageFortran95
	LLVMDWARFSourceLanguagePLI
	LLVMDWARFSourceLanguageObjC
	LLVMDWARFSourceLanguageObjCPlusPlus
	LLVMDWARFSourceLanguageUPC
	LLVMDWARFSourceLanguageD
	LLVMDWARFSourceLanguagePython
	LLVMDWARFSourceLanguageOpenCL
	LLVMDWARFSourceLanguageGo
	LLVMDWARFSourceLanguageModula3
	LLVMDWARFSourceLanguageHaskell
	LLVMDWARFSourceLanguageCPlusPlus03
	LLVMDWARFSourceLanguageCPlusPlus11
	LLVMDWARFSourceLanguageOCaml
	LLVMDWARFSourceLanguageRust
	LLVMDWARFSourceLanguageC11
	LLVMDWARFSourceLanguageSwift
	LLVMDWARFSourceLanguageJulia
	LLVMDWARFSourceLanguageDylan
	LLVMDWARFSourceLanguageCPlusPlus14
	LLVMDWARFSourceLanguageFortran03
	LLVMDWARFSourceLanguageFortran08
	LLVMDWARFSourceLanguageRenderScript
	LLVMDWARFSourceLanguageBLISS
	LLVMDWARFSourceLanguageKotlin
	LLVMDWARFSourceLanguageZig
	LLVMDWARFSourceLanguageCrystal
	LLVMDWARFSourceLanguageCPlusPlus17
	LLVMDWARFSourceLanguageCPlusPlus20
	LLVMDWARFSourceLanguageC17
	LLVMDWARFSourceLanguageFortran18
	LLVMDWARFSourceLanguageAda2005
	LLVMDWARFSourceLanguageAda2012
	LLVMDWARFSourceLanguageHIP
	LLVMDWARFSourceLanguageAssembly
	LLVMDWARFSourceLanguageCSharp
	LLVMDWARFSourceLanguageMojo
	LLVMDWARFSourceLanguageGLSL
	LLVMDWARFSourceLanguageGLSLES
	LLVMDWARFSourceLanguageHLSL
	LLVMDWARFSourceLanguageOpenCLCPP
	LLVMDWARFSourceLanguageCPPForOpenCL
	LLVMDWARFSourceLanguageSYCL
	LLVMDWARFSourceLanguageRuby
	LLVMDWARFSourceLanguageMove
	LLVMDWARFSourceLanguageHylo
	LLVMDWARFSourceLanguageMetal
	LLVMDWARFSourceLanguageMipsAssembler
	LLVMDWARFSourceLanguageGOOGLERenderScript
	LLVMDWARFSourceLanguageBORLANDDelphi
)

type LLVMCodeGenOptLevel C.LLVMCodeGenOptLevel

const (
	LLVMCodeGenLevelNone LLVMCodeGenOptLevel = iota
	LLVMCodeGenLevelLess
	LLVMCodeGenLevelDefault
	LLVMCodeGenLevelAggressive
)

type LLVMRelocMode C.LLVMRelocMode

const (
	LLVMRelocDefault LLVMRelocMode = iota
	LLVMRelocStatic
	LLVMRelocPIC
	LLVMRelocDynamicNoPic
	LLVMRelocROPI
	LLVMRelocRWPI
	LLVMRelocROPIRWPI
)

type LLVMCodeModel C.LLVMCodeModel

const (
	LLVMCodeModelDefault LLVMCodeModel = iota
	LLVMCodeModelJITDefault
	LLVMCodeModelTiny
	LLVMCodeModelSmall
	LLVMCodeModelKernel
	LLVMCodeModelMedium
	LLVMCodeModelLarge
)

type LLVMVerifierFailureAction C.LLVMVerifierFailureAction

const (
	LLVMAbortProcessAction LLVMVerifierFailureAction = iota
	LLVMPrintMessageAction
	LLVMReturnStatusAction
)

type LLVMCodeGenFileType C.LLVMCodeGenFileType

const (
	LLVMAssemblyFile LLVMCodeGenFileType = iota
	LLVMObjectFile
)

type LLVMGlobalISelAbortMode C.LLVMGlobalISelAbortMode

const (
	LLVMGlobalISelAbortEnable LLVMGlobalISelAbortMode = iota
	LLVMGlobalISelAbortDisable
	LLVMGlobalISelAbortDisableWithDiag
)

type RefT[T any] struct {
	C T
}

func (r RefT[T]) Ptr() unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&r.C))
}

type LLVMModuleRef struct {
	RefT[C.LLVMModuleRef]
}

func WrapExternalLLVMModuleRef(pointer unsafe.Pointer) (ref LLVMModuleRef) {
	ref.C = C.LLVMModuleRef(pointer)
	return ref
}

func (l LLVMModuleRef) AddGlobal(typ LLVMTypeRef, name string) (ref LLVMValueRef) {
	strName := C.CString(name)
	defer C.free(unsafe.Pointer(strName))

	ref.C = C.LLVMAddGlobal(l.C, typ.C, strName)
	return ref
}

func (l LLVMModuleRef) Context() (ref LLVMContextRef) {
	ref.C = C.LLVMGetModuleContext(l.C)
	return ref
}

func (l LLVMModuleRef) FirstGlobal() (ref LLVMValueRef) {
	ref.C = C.LLVMGetFirstGlobal(l.C)
	return ref
}

func (l LLVMModuleRef) String() string {
	str := C.LLVMPrintModuleToString(l.C)
	defer disposeMessage(str)
	return C.GoString(str)
}

func (l LLVMModuleRef) StripModuleDebugInfo() bool {
	return C.LLVMStripModuleDebugInfo(l.C) != 0
}

func (l LLVMModuleRef) Verify(action LLVMVerifierFailureAction) error {
	var strOutMessage *C.char
	ok := C.LLVMVerifyModule(l.C, C.LLVMVerifierFailureAction(action), &strOutMessage) == 0
	if !ok {
		if strOutMessage != nil {
			defer disposeMessage(strOutMessage)
			return errors.New(C.GoString(strOutMessage))
		}
	}
	return nil
}

type LLVMContextRef struct {
	RefT[C.LLVMContextRef]
}

func NewLLVMContext() (ref LLVMContextRef) {
	ref.C = C.LLVMContextCreate()
	return ref
}

func WrapExternalLLVMContextRef(pointer unsafe.Pointer) (ref LLVMContextRef) {
	ref.C = C.LLVMContextRef(pointer)
	return ref
}

func (l LLVMContextRef) Int1Type() (ref LLVMTypeRef) {
	ref.C = C.LLVMInt1TypeInContext(l.C)
	return ref
}

func (l LLVMContextRef) IntPtrType(dataRef LLVMTargetDataRef) (ref LLVMTypeRef) {
	ref.C = C.LLVMIntPtrTypeInContext(l.C, dataRef.C)
	return ref
}

func (l LLVMContextRef) PtrType(dataRef LLVMTargetDataRef) (ref LLVMTypeRef) {
	ref.C = C.LLVMIntPtrTypeInContext(l.C, dataRef.C)
	return ref
}

type LLVMPassBuilderOptionsRef struct {
	RefT[C.LLVMPassBuilderOptionsRef]
}

func NewPassBuilderOptions() (ref LLVMPassBuilderOptionsRef) {
	ref.C = C.LLVMCreatePassBuilderOptions()
	return ref
}

func (l LLVMPassBuilderOptionsRef) Dispose() {
	C.LLVMDisposePassBuilderOptions(l.C)
}

func LLVMRunPasses(module LLVMModuleRef,
	passes string,
	targetMachine LLVMTargetMachineRef,
	options LLVMPassBuilderOptionsRef,
) error {
	strPasses := C.CString(passes)
	defer C.free(unsafe.Pointer(strPasses))

	errRef := C.LLVMRunPasses(module.C, strPasses, targetMachine.C, options.C)
	if errRef != nil {
		strErrMessage := C.LLVMGetErrorMessage(errRef)
		defer C.LLVMDisposeErrorMessage(strErrMessage)
		return errors.New(C.GoString(strErrMessage))
	}
	return nil
}

type LLVMTargetDataRef struct {
	RefT[C.LLVMTargetDataRef]
}

func NewTargetDataLayout(machine LLVMTargetMachineRef) (ref LLVMTargetDataRef) {
	ref.C = C.LLVMCreateTargetDataLayout(machine.C)
	return ref
}

func WrapExternalLLVMTargetDataRef(pointer unsafe.Pointer) (ref LLVMTargetDataRef) {
	ref.C = C.LLVMTargetDataRef(pointer)
	return ref
}

func (l LLVMTargetDataRef) PointerSize() int {
	return int(C.LLVMPointerSize(l.C))
}

func (l LLVMTargetDataRef) PreferredAlignmentOfGlobal(value LLVMValueRef) uint {
	return uint(C.LLVMPreferredAlignmentOfGlobal(l.C, value.C))
}

type LLVMTargetMachineRef struct {
	RefT[C.LLVMTargetMachineRef]
}

func NewTargetMachine(
	target LLVMTargetRef,
	triple string,
	cpu string,
	features string,
	optLevel LLVMCodeGenOptLevel,
	relocMode LLVMRelocMode,
	codeModel LLVMCodeModel,
) (ref LLVMTargetMachineRef) {
	strTriple := C.CString(triple)
	defer C.free(unsafe.Pointer(strTriple))

	strCPU := C.CString(cpu)
	defer C.free(unsafe.Pointer(strCPU))

	strFeatures := C.CString(features)
	defer C.free(unsafe.Pointer(strFeatures))

	ref.C = C.LLVMCreateTargetMachine(
		target.C,
		strTriple,
		strCPU,
		strFeatures,
		C.LLVMCodeGenOptLevel(optLevel),
		C.LLVMRelocMode(relocMode),
		C.LLVMCodeModel(codeModel),
	)

	return ref
}

func (l LLVMTargetMachineRef) EmitToFile(
	module LLVMModuleRef,
	filename string,
	codegen LLVMCodeGenFileType,
) error {
	strFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(strFilename))

	var strErrOut *C.char
	ok := C.LLVMTargetMachineEmitToFile(
		l.C,
		module.C,
		strFilename,
		C.LLVMCodeGenFileType(codegen),
		&strErrOut,
	) == 0

	if !ok {
		if strErrOut != nil {
			defer disposeMessage(strErrOut)
			return errors.New(C.GoString(strErrOut))
		}
		return errors.New("EmitToFile failed")
	}

	return nil
}

type LLVMTargetRef struct {
	RefT[C.LLVMTargetRef]
}

func NewTargetFromTriple(triple string) (ref LLVMTargetRef, err error) {
	strTriple := C.CString(triple)
	defer C.free(unsafe.Pointer(strTriple))

	var strErr *C.char
	var c C.LLVMTargetRef
	ok := C.LLVMGetTargetFromTriple(strTriple, &c, &strErr) == 0
	if !ok {
		if strErr != nil {
			defer disposeMessage(strErr)
			return ref, errors.New(C.GoString(strErr))
		}
		return ref, errors.New("NewTargetFromTriple failed")
	}

	ref.C = c
	return ref, nil
}

type LLVMTypeRef struct {
	RefT[C.LLVMTypeRef]
}

type LLVMValueRef struct {
	RefT[C.LLVMValueRef]
}

func NewConstInt(typ LLVMTypeRef, value uint64, signExtend bool) (ref LLVMValueRef) {
	var cSignExtend C.LLVMBool
	if signExtend {
		cSignExtend = 1
	}
	ref.C = C.LLVMConstInt(typ.C, C.ulonglong(value), cSignExtend)
	return ref
}

func (l LLVMValueRef) SetAlignment(bytes uint) {
	C.LLVMSetAlignment(l.C, C.unsigned(bytes))
}

func (l LLVMValueRef) SetInitializer(value LLVMValueRef) {
	C.LLVMSetInitializer(l.C, value.C)
}

func (l LLVMValueRef) SetLinkage(linkage LLVMLinkage) {
	var c C.LLVMLinkage
	switch linkage {
	case LLVMLinkageExternal:
		c = C.LLVMExternalLinkage
	case LLVMLinkageAvailableExternally:
		c = C.LLVMAvailableExternallyLinkage
	case LLVMLinkageLinkonce:
		c = C.LLVMLinkOnceAnyLinkage
	case LLVMLinkageLinkonceODR:
		c = C.LLVMLinkOnceODRLinkage
	case LLVMLinkageWeak:
		c = C.LLVMWeakAnyLinkage
	case LLVMLinkageWeakODR:
		c = C.LLVMWeakODRLinkage
	case LLVMLinkageAppending:
		c = C.LLVMAppendingLinkage
	case LLVMLinkageInternal:
		c = C.LLVMInternalLinkage
	case LLVMLinkagePrivate:
		c = C.LLVMPrivateLinkage
	case LLVMLinkageExternWeak:
		c = C.LLVMExternalWeakLinkage
	case LLVMLinkageCommon:
		c = C.LLVMCommonLinkage
	}
	C.LLVMSetLinkage(l.C, C.LLVMLinkage(c))
}

func (l LLVMValueRef) SetGlobalConstant(value bool) {
	var cValue C.LLVMBool
	if value {
		cValue = 1
	}
	C.LLVMSetGlobalConstant(l.C, cValue)
}

func (l LLVMValueRef) IsNull() bool {
	return C.LLVMIsNull(l.C) != 0
}

func (l LLVMValueRef) Name() string {
	str := C.LLVMGetValueName(l.C)
	return C.GoString(str)
}

func (l LLVMValueRef) NextGlobal() (ref LLVMValueRef) {
	ref.C = C.LLVMGetNextGlobal(l.C)
	return ref
}

func (l LLVMValueRef) GlobalValueType() (ref LLVMTypeRef) {
	ref.C = C.LLVMGlobalGetValueType(l.C)
	return ref
}

func disposeMessage(message *C.char) {
	C.LLVMDisposeMessage(message)
}

func LLVMInitializeAllTargets() {
	C.LLVMInitializeAllTargets()
}

func LLVMInitializeAllTargetInfos() {
	C.LLVMInitializeAllTargetInfos()
}

func LLVMInitializeAllTargetMCs() {
	C.LLVMInitializeAllTargetMCs()
}

func LLVMInitializeAllAsmParsers() {
	C.LLVMInitializeAllAsmParsers()
}

func LLVMInitializeAllAsmPrinters() {
	C.LLVMInitializeAllAsmPrinters()
}

func LLVMInitializeAllDisassemblers() {
	C.LLVMInitializeAllDisassemblers()
}

func LLVMInitializeNativeTarget() bool {
	return C.LLVMInitializeNativeTarget() != 0
}
