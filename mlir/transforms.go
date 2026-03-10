package mlir

/*
#include <mlir-c/Transforms.h>
*/
import "C"

func RegisterTransformsPasses() {
	C.mlirRegisterTransformsPasses()
}

func RegisterTransformsBubbleDownMemorySpaceCasts() {
	C.mlirRegisterTransformsBubbleDownMemorySpaceCasts()
}

func NewTransformsBubbleDownMemorySpaceCasts() Pass {
	return Pass(C.mlirCreateTransformsBubbleDownMemorySpaceCasts())
}

func RegisterTransformsCSE() {
	C.mlirRegisterTransformsCSE()
}

func NewTransformsCSE() Pass {
	return Pass(C.mlirCreateTransformsCSE())
}

func RegisterTransformsCanonicalizer() {
	C.mlirRegisterTransformsCanonicalizer()
}

func NewTransformsCanonicalizer() Pass {
	return Pass(C.mlirCreateTransformsCanonicalizer())
}

func RegisterTransformsCompositeFixedPointPass() {
	C.mlirRegisterTransformsCompositeFixedPointPass()
}

func NewTransformsCompositeFixedPointPass() Pass {
	return Pass(C.mlirCreateTransformsCompositeFixedPointPass())
}

func RegisterTransformsControlFlowSink() {
	C.mlirRegisterTransformsControlFlowSink()
}

func NewTransformsControlFlowSink() Pass {
	return Pass(C.mlirCreateTransformsControlFlowSink())
}

func RegisterTransformsGenerateRuntimeVerification() {
	C.mlirRegisterTransformsGenerateRuntimeVerification()
}

func NewTransformsGenerateRuntimeVerification() Pass {
	return Pass(C.mlirCreateTransformsGenerateRuntimeVerification())
}

func RegisterTransformsInliner() {
	C.mlirRegisterTransformsInliner()
}

func NewTransformsInliner() Pass {
	return Pass(C.mlirCreateTransformsInliner())
}

func RegisterTransformsLocationSnapshot() {
	C.mlirRegisterTransformsLocationSnapshot()
}

func NewTransformsLocationSnapshot() Pass {
	return Pass(C.mlirCreateTransformsLocationSnapshot())
}

func RegisterTransformsLoopInvariantCodeMotion() {
	C.mlirRegisterTransformsLoopInvariantCodeMotion()
}

func NewTransformsLoopInvariantCodeMotion() Pass {
	return Pass(C.mlirCreateTransformsLoopInvariantCodeMotion())
}

func RegisterTransformsLoopInvariantSubsetHoisting() {
	C.mlirRegisterTransformsLoopInvariantSubsetHoisting()
}

func NewTransformsLoopInvariantSubsetHoisting() Pass {
	return Pass(C.mlirCreateTransformsLoopInvariantSubsetHoisting())
}

func RegisterTransformsMem2Reg() {
	C.mlirRegisterTransformsMem2Reg()
}

func NewTransformsMem2Reg() Pass {
	return Pass(C.mlirCreateTransformsMem2Reg())
}

func RegisterTransformsPrintIRPass() {
	C.mlirRegisterTransformsPrintIRPass()
}

func NewTransformsPrintIRPass() Pass {
	return Pass(C.mlirCreateTransformsPrintIRPass())
}

func RegisterTransformsPrintOpStats() {
	C.mlirRegisterTransformsPrintOpStats()
}

func NewTransformsPrintOpStats() Pass {
	return Pass(C.mlirCreateTransformsPrintOpStats())
}

func RegisterTransformsRemoveDeadValues() {
	C.mlirRegisterTransformsRemoveDeadValues()
}

func NewTransformsRemoveDeadValues() Pass {
	return Pass(C.mlirCreateTransformsRemoveDeadValues())
}

func RegisterTransformsSCCP() {
	C.mlirRegisterTransformsSCCP()
}

func NewTransformsSCCP() Pass {
	return Pass(C.mlirCreateTransformsSCCP())
}

func RegisterTransformsSROA() {
	C.mlirRegisterTransformsSROA()
}

func NewTransformsSROA() Pass {
	return Pass(C.mlirCreateTransformsSROA())
}

func RegisterTransformsStripDebugInfo() {
	C.mlirRegisterTransformsStripDebugInfo()
}

func NewTransformsStripDebugInfo() Pass {
	return Pass(C.mlirCreateTransformsStripDebugInfo())
}

func RegisterTransformsSymbolDCE() {
	C.mlirRegisterTransformsSymbolDCE()
}

func NewTransformsSymbolDCE() Pass {
	return Pass(C.mlirCreateTransformsSymbolDCE())
}

func RegisterTransformsSymbolPrivatize() {
	C.mlirRegisterTransformsSymbolPrivatize()
}

func NewTransformsSymbolPrivatize() Pass {
	return Pass(C.mlirCreateTransformsSymbolPrivatize())
}

func RegisterTransformsTopologicalSort() {
	C.mlirRegisterTransformsTopologicalSort()
}

func NewTransformsTopologicalSort() Pass {
	return Pass(C.mlirCreateTransformsTopologicalSort())
}

func RegisterTransformsViewOpGraph() {
	C.mlirRegisterTransformsViewOpGraph()
}

func NewTransformsViewOpGraph() Pass {
	return Pass(C.mlirCreateTransformsViewOpGraph())
}
