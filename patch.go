package gomonkey

import (
    "reflect"
    "fmt"
    "unsafe"
    "syscall"
)

type Patch struct {
    targetBytes []byte
    double      *reflect.Value
}

type Patches struct {
    patches map[reflect.Value]Patch
}

func New() *Patches {
    return &Patches{make(map[reflect.Value]Patch)}
}

func (this *Patches) ApplyFunc(target, double interface{}) *Patches {
    t := reflect.ValueOf(target)
    d := reflect.ValueOf(double)
    this.check(t, d)
    bytes := replace(*(*uintptr)(getPointer(t)), uintptr(getPointer(d)))
    this.patches[t] = Patch{bytes, &d}
    return this
}

func (this *Patches) Reset() {
    for target, patch := range this.patches {
        modifyBinary(*(*uintptr)(getPointer(target)), patch.targetBytes)
        delete(this.patches, target)
    }

}

func (this *Patches) check(target, double reflect.Value) {
    if target.Kind() != reflect.Func {
        panic("target is not a func")
    }

    if double.Kind() != reflect.Func {
        panic("double is not a func")
    }

    if target.Type() != double.Type() {
        panic(fmt.Sprintf("target type(%s) and double type(%s) are different", target.Type(), double.Type()))
    }

    if _, ok := this.patches[target]; ok {
        panic("patch has been existed")
    }
}

func ApplyFunc(target, double interface{}) *Patches {
    return New().ApplyFunc(target, double)
}

func replace(target, double uintptr) []byte {
    data := jmpPrepare(double)
    bytes := entryAddress(target, len(data))
    modifyBinary(target, data)
    return bytes
}

type value struct {
    _ uintptr
    p unsafe.Pointer
}

func getPointer(v reflect.Value) unsafe.Pointer {
    return (*value)(unsafe.Pointer(&v)).p
}

func entryAddress(p uintptr, l int) []byte {
    return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: p, Len: l, Cap: l}))
}

func modifyBinary(target uintptr, bytes []byte) {
    function := entryAddress(target, len(bytes))

    page := entryAddress(pageStart(target), syscall.Getpagesize())
    err := syscall.Mprotect(page, syscall.PROT_READ | syscall.PROT_WRITE | syscall.PROT_EXEC)
    if err != nil {
        panic(err)
    }
    copy(function, bytes)

    err = syscall.Mprotect(page, syscall.PROT_READ | syscall.PROT_EXEC)
    if err != nil {
        panic(err)
    }
}

func pageStart(ptr uintptr) uintptr {
    return ptr & ^(uintptr(syscall.Getpagesize() - 1))
}