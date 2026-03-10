package mlir

/*
#include <llvm-c/Target.h>
#include <mlir-c/Dialect/LLVM.h>
#include <mlir-c/Target/LLVMIR.h>
*/
import "C"
import "unsafe"

type LLVMModuleRef struct {
	C C.LLVMModuleRef
}

func WrapExternalLLVMModuleRef(pointer unsafe.Pointer) LLVMModuleRef {
	return LLVMModuleRef{
		C: C.LLVMModuleRef(pointer),
	}
}

func (l LLVMModuleRef) Ptr() unsafe.Pointer {
	return unsafe.Pointer(l.C)
}

type LLVMContextRef struct {
	C C.LLVMContextRef
}

func WrapExternalLLVMContextRef(pointer unsafe.Pointer) LLVMContextRef {
	return LLVMContextRef{
		C: C.LLVMContextRef(pointer),
	}
}

func (l LLVMContextRef) Ptr() unsafe.Pointer {
	return unsafe.Pointer(l.C)
}

func WrapExternalLLVMTargetDataRef(pointer unsafe.Pointer) LLVMTargetDataRef {
	return LLVMTargetDataRef{
		C: C.LLVMTargetDataRef(pointer),
	}
}

func (l LLVMTargetDataRef) Ptr() unsafe.Pointer {
	return unsafe.Pointer(l.C)
}

type LLVMTargetDataRef struct {
	C C.LLVMTargetDataRef
}

func TranslateModuleToLLVMIR(module Operation, ctx LLVMContextRef) LLVMModuleRef {
	return LLVMModuleRef{
		C: C.mlirTranslateModuleToLLVMIR(module.Raw(), ctx.C),
	}
}

type LLVMLinkage C.MlirLLVMLinkage

const (
	LLVMLinkageExternal            LLVMLinkage = C.MlirLLVMLinkageExternal
	LLVMLinkageAvailableExternally LLVMLinkage = C.MlirLLVMLinkageAvailableExternally
	LLVMLinkageLinkonce            LLVMLinkage = C.MlirLLVMLinkageLinkonce
	LLVMLinkageLinkonceODR         LLVMLinkage = C.MlirLLVMLinkageLinkonceODR
	LLVMLinkageWeak                LLVMLinkage = C.MlirLLVMLinkageWeak
	LLVMLinkageWeakODR             LLVMLinkage = C.MlirLLVMLinkageWeakODR
	LLVMLinkageAppending           LLVMLinkage = C.MlirLLVMLinkageAppending
	LLVMLinkageInternal            LLVMLinkage = C.MlirLLVMLinkageInternal
	LLVMLinkagePrivate             LLVMLinkage = C.MlirLLVMLinkagePrivate
	LLVMLinkageExternWeak          LLVMLinkage = C.MlirLLVMLinkageExternWeak
	LLVMLinkageCommon              LLVMLinkage = C.MlirLLVMLinkageCommon
)

type LLVMLinkageAttr struct {
	baseAttribute
}

func NewLLVMLinkageAttr(ctx Context, linkage LLVMLinkage) LLVMLinkageAttr {
	return LLVMLinkageAttr{
		baseAttribute: baseAttribute(C.mlirLLVMLinkageAttrGet(ctx.Raw(), C.MlirLLVMLinkage(linkage))),
	}
}

func LLVMLinkageAttrName() string {
	return BorrowedStringRef(C.mlirLLVMLinkageAttrGetName()).String()
}

type LLVMDIFileAttr struct {
	baseAttribute
}

func NewLLVMDIFileAttr(ctx Context, name StringAttr, directory StringAttr) LLVMDIFileAttr {
	return LLVMDIFileAttr{
		baseAttribute: baseAttribute(C.mlirLLVMDIFileAttrGet(ctx.Raw(), name.Raw(), directory.Raw())),
	}
}

func LLVMDIFileAttrName() string {
	return BorrowedStringRef(C.mlirLLVMDIFileAttrGetName()).String()
}

type LLVMDIEmissionKind = C.MlirLLVMDIEmissionKind

const (
	LLVMDIEmissionKindNone                LLVMDIEmissionKind = C.MlirLLVMDIEmissionKindNone
	LLVMDIEmissionKindFull                LLVMDIEmissionKind = C.MlirLLVMDIEmissionKindFull
	LLVMDIEmissionKindLineTablesOnly      LLVMDIEmissionKind = C.MlirLLVMDIEmissionKindLineTablesOnly
	LLVMDIEmissionKindDebugDirectivesOnly LLVMDIEmissionKind = C.MlirLLVMDIEmissionKindDebugDirectivesOnly
)

type LLVMDINameTableKind = C.MlirLLVMDINameTableKind

const (
	LLVMDINameTableKindDefault LLVMDINameTableKind = C.MlirLLVMDINameTableKindDefault
	LLVMDINameTableKindGNU     LLVMDINameTableKind = C.MlirLLVMDINameTableKindGNU
	LLVMDINameTableKindNone    LLVMDINameTableKind = C.MlirLLVMDINameTableKindNone
	LLVMDINameTableKindApple   LLVMDINameTableKind = C.MlirLLVMDINameTableKindApple
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

type LLVMDICompileUnitAttr struct {
	baseAttribute
}

func NewLLVMDICompileUnitAttr(
	ctx Context,
	id AttributeLike,
	sourceLanguage uint32,
	file LLVMDIFileAttr,
	producer StringAttr,
	isOptimized bool,
	emissionKind LLVMDIEmissionKind,
	nameTableKind LLVMDINameTableKind,
	splitDebugFilename StringAttr,
) LLVMDICompileUnitAttr {
	return LLVMDICompileUnitAttr{baseAttribute: baseAttribute(C.mlirLLVMDICompileUnitAttrGet(
		ctx.Raw(),
		id.Raw(),
		C.uint32_t(sourceLanguage),
		file.Raw(),
		producer.Raw(),
		C.bool(isOptimized),
		C.MlirLLVMDIEmissionKind(emissionKind),
		C.MlirLLVMDINameTableKind(nameTableKind),
		splitDebugFilename.Raw(),
	))}
}

func LLVMDICompileUnitAttrName() string {
	return BorrowedStringRef(C.mlirLLVMDICompileUnitAttrGetName()).String()
}

func NewLLVMLegalizeForExportPass() Pass {
	return Pass(C.mlirCreateLLVMLLVMLegalizeForExportPass())
}
