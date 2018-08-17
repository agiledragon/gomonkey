package gomonkey

import (
    "reflect"
    "fmt"
    "unsafe"
    "syscall"
)

type Patches struct {
    originals map[reflect.Value][]byte
}

func ApplyFunc(target, double interface{}) *Patches {
    return New().ApplyFunc(target, double)
}

func ApplyMethod(target reflect.Type, methodName string, double interface{}) *Patches {
    return New().ApplyMethod(target, methodName, double)
}

func New() *Patches {
    return &Patches{make(map[reflect.Value][]byte)}
}

func (this *Patches) ApplyFunc(target, double interface{}) *Patches {
    t := reflect.ValueOf(target)
    d := reflect.ValueOf(double)
    this.applyCore(t, d)
    return this
}

func (this *Patches) ApplyMethod(target reflect.Type, methodName string, double interface{}) *Patches {
    m, ok := target.MethodByName(methodName);
    if !ok {
        panic("retrieve method by name failed")
    }
    d := reflect.ValueOf(double)
    this.applyCore(m.Func, d)
    return this
}

func (this *Patches) Reset() {
    for target, bytes := range this.originals {
        modifyBinary(*(*uintptr)(getPointer(target)), bytes)
        delete(this.originals, target)
    }
}

func (this *Patches) applyCore(target, double reflect.Value) {
    this.check(target, double)
    original := replace(*(*uintptr)(getPointer(target)), uintptr(getPointer(double)))
    this.originals[target] = original
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

    if _, ok := this.originals[target]; ok {
        panic("patch has been existed")
    }
}

func replace(target, double uintptr) []byte {
    code := buildJmpDirective(double)
    bytes := entryAddress(target, len(code))
    original := make([]byte, len(bytes))
    copy(original, bytes)
    modifyBinary(target, code)
    return original
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