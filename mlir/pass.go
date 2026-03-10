package mlir

/*
#include <mlir-c/Pass.h>
*/
import "C"
import (
	"io"
	"unsafe"
)

//===----------------------------------------------------------------------===//
// Pass Type.
//===----------------------------------------------------------------------===//

type Pass C.MlirPass

func (p Pass) Raw() C.MlirPass {
	return C.MlirPass(p)
}

func WrapExternalPass(ptr unsafe.Pointer) Pass {
	return Pass{ptr: ptr}
}

//===----------------------------------------------------------------------===//
// PassManager/OpPassManager APIs.
//===----------------------------------------------------------------------===//

type PassManager C.MlirPassManager

// NewPassManager creates a new top-level PassManager with the default anchor.
func NewPassManager(ctx Context) PassManager {
	return PassManager(C.mlirPassManagerCreate(ctx.Raw()))
}

// NewPassManagerOnOperation creates a new top-level PassManager anchored on `anchorOp`.
func NewPassManagerOnOperation(ctx Context, anchorOp string) PassManager {
	refAnchorOp := NewStringRef(anchorOp)
	defer refAnchorOp.Destroy()
	return PassManager(C.mlirPassManagerCreateOnOperation(ctx.Raw(), refAnchorOp.Raw()))
}

func (pm PassManager) Raw() C.MlirPassManager { return C.MlirPassManager(pm) }
func (pm PassManager) IsNull() bool           { return bool(C.mlirPassManagerIsNull(pm.Raw())) }
func (pm PassManager) Destroy()               { C.mlirPassManagerDestroy(pm.Raw()) }
func (pm PassManager) AsOpPassManager() OpPassManager {
	return OpPassManager(C.mlirPassManagerGetAsOpPassManager(pm.Raw()))
}

func (pm PassManager) Run(op Operation) LogicalResult {
	return LogicalResult(C.mlirPassManagerRunOnOp(pm.Raw(), op.Raw()))
}

type IRPrinterConfig struct {
	PrintBeforeAll          bool
	PrintAfterAll           bool
	PrintModuleScope        bool
	PrintAfterOnlyOnChange  bool
	PrintAfterOnlyOnFailure bool
	Flags                   OpPrintingFlags
	TreePrintingPath        string
}

func (pm PassManager) EnableIRPrinting(config IRPrinterConfig) PassManager {
	if config.Flags.IsNull() {
		// Use default flags.
		config.Flags = NewOpPrintingFlags()
		defer config.Flags.Destroy()
	}

	refTreePrintingPath := NewStringRef(config.TreePrintingPath)
	defer refTreePrintingPath.Destroy()

	C.mlirPassManagerEnableIRPrinting(
		pm.Raw(),
		C.bool(config.PrintBeforeAll),
		C.bool(config.PrintAfterAll),
		C.bool(config.PrintModuleScope),
		C.bool(config.PrintAfterOnlyOnChange),
		C.bool(config.PrintAfterOnlyOnFailure),
		config.Flags.Raw(),
		refTreePrintingPath.Raw())
	return pm
}

func (pm PassManager) EnableVerifier(enable bool) PassManager {
	C.mlirPassManagerEnableVerifier(pm.Raw(), C.bool(enable))
	return pm
}

func (pm PassManager) EnableTiming() PassManager {
	C.mlirPassManagerEnableTiming(pm.Raw())
	return pm
}

type PassDisplayMode C.MlirPassDisplayMode

const (
	PassDisplayModeList     PassDisplayMode = C.MLIR_PASS_DISPLAY_MODE_LIST
	PassDisplayModePipeline PassDisplayMode = C.MLIR_PASS_DISPLAY_MODE_PIPELINE
)

func (pm PassManager) EnableStatistics(displayMode PassDisplayMode) PassManager {
	C.mlirPassManagerEnableStatistics(pm.Raw(), C.MlirPassDisplayMode(displayMode))
	return pm
}

func (pm PassManager) NestedUnder(opName string) OpPassManager {
	refOpName := NewStringRef(opName)
	defer refOpName.Destroy()
	return OpPassManager(C.mlirPassManagerGetNestedUnder(pm.Raw(), refOpName.Raw()))
}

func (pm PassManager) AddOwnedPass(pass Pass) PassManager {
	C.mlirPassManagerAddOwnedPass(pm.Raw(), pass.Raw())
	return pm
}

type OpPassManager C.MlirOpPassManager

func (opm OpPassManager) Raw() C.MlirOpPassManager { return C.MlirOpPassManager(opm) }
func (opm OpPassManager) NestedUnder(opName string) OpPassManager {
	refOpName := NewStringRef(opName)
	defer refOpName.Destroy()
	return OpPassManager(C.mlirOpPassManagerGetNestedUnder(opm.Raw(), refOpName.Raw()))
}

func (opm OpPassManager) AddOwnedPass(pass Pass) OpPassManager {
	C.mlirOpPassManagerAddOwnedPass(opm.Raw(), pass.Raw())
	return opm
}

func (opm OpPassManager) AddPipeline(pipelineElements string, out io.Writer) OpPassManager {
	refPipelineElements := NewStringRef(pipelineElements)
	defer refPipelineElements.Destroy()

	writeString(out, func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirOpPassManagerAddPipeline(opm.Raw(), refPipelineElements.Raw(), cb, ud)
	})

	return opm
}

func (opm OpPassManager) Parse(pipeline string, out io.Writer) OpPassManager {
	refPipeline := NewStringRef(pipeline)
	defer refPipeline.Destroy()

	writeString(out, func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirParsePassPipeline(opm.Raw(), refPipeline.Raw(), cb, ud)
	})

	return opm
}

func (opm OpPassManager) String() string {
	return collectString(func(cb C.MlirStringCallback, ud unsafe.Pointer) {
		C.mlirPrintPassPipeline(opm.Raw(), cb, ud)
	})
}

//===----------------------------------------------------------------------===//
// External Pass API.
//===----------------------------------------------------------------------===//

// TODO: Implement external pass such that cgo.Handle holding a callback is deleted at the appropriate time.
