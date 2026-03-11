package mlir

/*
#include <mlir-c/Dialect/LLVM.h>
#include <mlir-c/Target/LLVMIR.h>
*/
import "C"

func TranslateModuleToLLVMIR(module Operation, ctx LLVMContextRef) (ref LLVMModuleRef) {
	ref.C = C.mlirTranslateModuleToLLVMIR(module.Raw(), ctx.C)
	return ref
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
