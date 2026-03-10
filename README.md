# go-mlir

`go-mlir` provides Go bindings for the [MLIR](https://mlir.llvm.org/) framework, a part of the LLVM project. It allows you to build and manipulate MLIR's intermediate representation (IR) using Go. These bindings are built on top of the official MLIR C API.

## Prerequisites

Before you begin, ensure you have the following dependencies installed on your system:

*   **Go**: Version 1.18 or later.
*   **C/C++ Toolchain**: A compatible C and C++ compiler, such as GCC or Clang.
*   **CMake**: Version 3.20 or later.
*   **Ninja**: A small build system with a focus on speed, used for building LLVM/MLIR.

## Import the package

```sh
go get pkg.si-go.dev/go-mlir
```

## Building

The build process involves two main steps: first, building the required MLIR C libraries from the included LLVM submodule, and second, building the Go package itself.

### 1. Clone the Repository

First, clone the `go-mlir` repository, making sure to initialize the `llvm-project` submodule:

```sh
git clone --recurse-submodules https://github.com/waj334/go-mlir.git
cd go-mlir
```

### 2. Build MLIR

Next, configure and build the MLIR libraries using CMake and Ninja.

```sh
# Create a build directory
mkdir -p thirdparty/llvm-project/build
cd thirdparty/llvm-project/build

# Configure the build with CMake
cmake -G Ninja ../llvm \
    -DLLVM_ENABLE_PROJECTS=mlir \
    -DLLVM_TARGETS_TO_BUILD="host" \
    -DCMAKE_BUILD_TYPE=Release \
    -DLLVM_ENABLE_ASSERTIONS=ON \
    -DMLIR_ENABLE_BINDINGS_C=ON \
    -DBUILD_SHARED_LIBS=ON \
    -DLLVM_ENABLE_RTTI=ON

# Build the MLIR libraries
ninja
```

This will compile the necessary MLIR components as shared libraries, which the Go bindings will link against.

### 3. Build the Go Package

Once MLIR is built, you can build the `go-mlir` package. The Go build process uses `cgo` to link against the MLIR libraries you just built.

Set the following environment variables to tell `cgo` where to find the MLIR headers and libraries:

```sh
export CGO_CFLAGS="-I$(pwd)/../thirdparty/llvm-project/build/include"
export CGO_LDFLAGS="-L$(pwd)/../thirdparty/llvm-project/build/lib"
export LD_LIBRARY_PATH="$(pwd)/../thirdparty/llvm-project/build/lib:$LD_LIBRARY_PATH"
```

Now, you can build or run your Go application:

```sh
go build ./...
go test ./...
```

## Linking

The `go-mlir` package relies on `cgo` to dynamically link against the MLIR shared libraries at build time. The necessary linker flags are already included in the source files (like `cgo.go`).

At runtime, the dynamic linker needs to be able to locate the MLIR `.so` files. You can ensure this by:

1.  **Setting `LD_LIBRARY_PATH`**: Add the path to your MLIR build's `lib` directory to the `LD_LIBRARY_PATH` environment variable before running your application.

    ```sh
    export LD_LIBRARY_PATH=/path/to/go-mlir/thirdparty/llvm-project/build/lib:$LD_LIBRARY_PATH
    go run your/app
    ```

2.  **Using `rpath`**: Alternatively, you can embed the library path into your executable at build time using an `rpath`.

    ```sh
    go build -ldflags='-r /path/to/go-mlir/thirdparty/llvm-project/build/lib' your/app
    ```

## Usage Example

Here is a simple example of how to create an MLIR context and a module:

```go
package main

import (
	"fmt"
	"pkg.si-go.dev/go-mlir/mlir"
)

func main() {
	// Create a new MLIR context.
	ctx := mlir.NewContext()
	defer ctx.Destroy()

	// Create a location for our operations.
	loc := mlir.NewUnknownLoc(ctx)

	// Create an empty module.
	module := mlir.NewModule(loc)
	defer module.Destroy()

	// Print the module's IR.
	fmt.Println(module.Operation().String())
}
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any bugs or feature requests.

## License

`go-mlir` is licensed under the Apache License v2.0 with LLVM Exceptions. See the `LICENSE` file for more details.
```