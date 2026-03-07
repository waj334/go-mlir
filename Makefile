ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

ifeq ($(OS),Windows_NT)
	EXECUTABLE_POSTFIX=.exe
	CGO_LDFLAGS += -lole32 -luuid -lpsapi -lshell32 -ladvapi32 -lntdll
	CMAKE_CXX_FLAGS += -pthread -femulated-tls
	CMAKE_CXX_STANDARD_LIBRARIES += -lpthread
	CLANG_TARGET := x86_64-pc-windows-gnu
	CC ?= clang
	CXX ?= clang++
	# NOTE: ld should be replaced with ld.lld directly on Windows since Go is dumb.
else
	CGO_LDFLAGS += -fuse-ld=lld -lrt -ldl -lpthread -lm -lz -ltinfo -lzstd
endif

LLVM_SRC_DIR=$(ROOT_DIR)/thirdparty/llvm-project
#LLVM_BUILD_DIR=$(ROOT_DIR)/build/$(CMAKE_BUILD_TYPE)/llvm-build
LLVM_BUILD_DIR=/home/waj334/Projects/sigo/build/Debug/llvm-build
LLVM_CMAKE_CACHE=$(LLVM_BUILD_DIR)/CMakeCache.txt
LLVM_CONFIG_EXECUTABLE=${LLVM_BUILD_DIR}/bin/llvm-config$(EXECUTABLE_POSTFIX)
LLVM_BUILD_COMPONENTS := ARM AVR RISCV
LLVM_COMPONENTS := ARM AVR RISCV passes

# Build a semicolon separated list that CMake can accept
CMAKE_LLVM_COMPONENTS :=
$(foreach item, $(LLVM_BUILD_COMPONENTS),$(if $(CMAKE_LLVM_COMPONENTS),$(eval CMAKE_LLVM_COMPONENTS := $(CMAKE_LLVM_COMPONENTS);))$(eval CMAKE_LLVM_COMPONENTS := $(CMAKE_LLVM_COMPONENTS)$(strip $(item))))

# Determine build flags required by LLVM
CGO_LDFLAGS += -Wl,--gc-sections $(shell ${LLVM_CONFIG_EXECUTABLE} --ldflags) -Wl,--start-group $(shell ${LLVM_CONFIG_EXECUTABLE} --libs ${LLVM_COMPONENTS}) -Wl,--end-group
CGO_CFLAGS += -fPIC -ffunction-sections -fdata-sections $(shell ${LLVM_CONFIG_EXECUTABLE} --cflags)

# Add LLVM includes
CGO_CFLAGS += -I$(LLVM_SRC_DIR)/llvm/include
CGO_CFLAGS += -I${LLVM_BUILD_DIR}/tools/mlir/include

# Add MLIR includes
CGO_CFLAGS += -I$(LLVM_SRC_DIR)/mlir/include

build:
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go build ./mlir

example-minimal:
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go build -gcflags "all=-N -l" -ldflags="-linkmode external" -o ./bin/example-minimal ./examples/minimal
