package gomonkey

import (
    "syscall"
    "unsafe"
)

func modifyBinary(target uintptr, bytes []byte) {
    function := entryAddress(target, len(bytes))
    
    proc := syscall.NewLazyDLL("kernel32.dll").NewProc("VirtualProtect")
    PROT_READ_WRITE := 0x40
    var old uint32
    result := proc.Call(target, len(bytes), PROT_READ_WRITE, unsafe.Pointer(&old))
    if result == 0 {
        panic(result)
    }
    copy(function, bytes)
    
    var ignore uint32
    result = proc.Call(target, len(bytes), old, unsafe.Pointer(&ignore))
    if result == 0 {
        panic(result)
    }
}