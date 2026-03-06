package mlir

/*
#include <mlir-c/RegisterEverything.h>
*/
import "C"

func (d DialectRegistry) RegisterAllDialects() {
	C.mlirRegisterAllDialects(d.raw())
}

func (c Context) RegisterAllLLVMTranslations() {
	C.mlirRegisterAllLLVMTranslations(c.raw())
}

func RegisterAllPasses() {
	C.mlirRegisterAllPasses()
}
