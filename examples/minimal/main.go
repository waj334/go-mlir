package main

import "pkg.si-go.dev/go-mlir/mlir"

func main() {
	// Initialize MLIR.
	ctx := mlir.NewContext()
	defer ctx.Destroy()

	// Register dialects.
	registry := mlir.NewDialectRegistry()
	registry.RegisterAllDialects()
	ctx.AppendDialectRegistry(registry)
	registry.Destroy()

	// Load all available dialect in the current context.
	ctx.LoadAllAvailableDialects()

	// Create the new module.
	module := mlir.NewModule(mlir.NewUnknownLoc(ctx))
	defer module.Destroy()

	// Get the module body block.
	moduleBody := module.Body()
	loc := mlir.NewUnknownLoc(ctx)

	memrefType := mlir.NewTypeFromString(ctx, "memref<?xf32>")

	// Create a function.
	funcBodyRegion := mlir.NewRegion()
	funcBody := mlir.NewBlock([]mlir.Type{memrefType, memrefType}, []mlir.LocationLike{loc, loc})
	funcBodyRegion.AppendOwnedBlock(funcBody)
	funcTypeAttr := mlir.NewAttributeFromString(ctx, "(memref<?xf32>, memref<?xf32>) -> ()")
	funcNameAttr := mlir.NewAttributeFromString(ctx, "\"add\"")
	funcAttrs := []mlir.NamedAttribute{
		mlir.NewNamedAttribute("function_type", funcTypeAttr),
		mlir.NewNamedAttribute("sym_name", funcNameAttr),
	}

	funcOp := mlir.NewOperationState("func.func", loc).
		AddAttributes(funcAttrs...).
		AddOwnedRegions(funcBodyRegion).
		Create()

	// Add the function to the module.
	moduleBody.AppendOwnedOperation(funcOp)

	// Create constant op (const 0).
	indexType := mlir.NewTypeFromString(ctx, "index")
	indexZeroLiteral := mlir.NewAttributeFromString(ctx, "0 : index")
	indexZeroValueAttr := mlir.NewNamedAttribute("value", indexZeroLiteral)
	constZeroOp := mlir.NewOperationState("arith.constant", loc).
		AddResults(indexType).
		AddAttributes(indexZeroValueAttr).
		Create()
	funcBody.AppendOwnedOperation(constZeroOp)

	// Create dim op.
	funcArg0 := funcBody.Argument(0)
	constZeroValue := constZeroOp.Result(0)
	dimOp := mlir.NewOperationState("memref.dim", loc).
		AddOperands(funcArg0, constZeroValue).
		AddResults(indexType).
		Create()
	funcBody.AppendOwnedOperation(dimOp)

	loopBodyRegion := mlir.NewRegion()
	loopBody := mlir.NewBlock(nil, nil)
	loopBody.AddArgument(indexType, loc)
	loopBodyRegion.AppendOwnedBlock(loopBody)

	// Create constant op (const 1).
	indexOneLiteral := mlir.NewAttributeFromString(ctx, "1 : index")
	indexOneValueAttr := mlir.NewNamedAttribute("value", indexOneLiteral)
	constOneOp := mlir.NewOperationState("arith.constant", loc).
		AddResults(indexType).
		AddAttributes(indexOneValueAttr).
		Create()
	funcBody.AppendOwnedOperation(constOneOp)

	dimValue := dimOp.Result(0)
	constOneValue := constOneOp.Result(0)
	loopOp := mlir.NewOperationState("scf.for", loc).
		AddOperands(constZeroValue, dimValue, constOneValue).
		AddOwnedRegions(loopBodyRegion).
		Create()
	funcBody.AppendOwnedOperation(loopOp)

	retOp := mlir.NewOperationState("func.return", loc).Create()
	funcBody.AppendOwnedOperation(retOp)

	// Dump the module.
	module.Operation().Dump()
}
