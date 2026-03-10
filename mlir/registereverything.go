package mlir

/*
#include <mlir-c/RegisterEverything.h>
*/
import "C"

func (d DialectRegistry) RegisterAllDialects() {
	C.mlirRegisterAllDialects(d.Raw())
}

func (c Context) RegisterAllLLVMTranslations() {
	C.mlirRegisterAllLLVMTranslations(c.Raw())
}

func RegisterAllPasses() {
	C.mlirRegisterAllPasses()
}
