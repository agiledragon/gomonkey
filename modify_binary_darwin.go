package gomonkey

import (
	"reflect"
	"syscall"
	"unsafe"
)

func modifyBinary(target uintptr, bytes []byte) {
	function := entryAddress(target, len(bytes))

	page := entryAddress(pageStart(target), syscall.Getpagesize())

	machVmProtect(page)
	copy(function, bytes)
}

func machVmProtect(page []byte) {
	err := syscall.Mprotect(page, syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC)
	if err != nil {
		panic(err)
	}

	ret := MachVMProtect(MachTaskSelf(), pageStart((*reflect.SliceHeader)(unsafe.Pointer(&page)).Data), uint64(syscall.Getpagesize()), uint(0), 0x17)
	if ret != 0 {
		panic("machVmProtect failed")
	}
}

//go:cgo_import_dynamic libsystem_mach_task_self mach_task_self "/usr/lib/libSystem.B.dylib"
func mach_task_self_trampoline()

func MachTaskSelf() uint {
	args := struct {
		ret uint
	}{}
	libcCall(unsafe.Pointer(reflect.ValueOf(mach_task_self_trampoline).Pointer()), unsafe.Pointer(&args))
	return args.ret
}

//go:cgo_import_dynamic libsystem_mach_vm_protect mach_vm_protect "/usr/lib/libSystem.B.dylib"
func mach_vm_protect_trampoline()

func MachVMProtect(targetTask uint, address uintptr, size uint64, setMaximum uint, newProtection int) int {
	args := struct {
		targetTask    uint
		address       uintptr
		size          uint64
		setMaximum    uint
		newProtection int
		ret           int
	}{
		targetTask:    targetTask,
		address:       address,
		size:          size,
		setMaximum:    setMaximum,
		newProtection: newProtection,
		ret:           0,
	}
	libcCall(unsafe.Pointer(reflect.ValueOf(mach_vm_protect_trampoline).Pointer()), unsafe.Pointer(&args))
	return args.ret
}

//go:linkname libcCall runtime.libcCall
func libcCall(fn, arg unsafe.Pointer) int32
