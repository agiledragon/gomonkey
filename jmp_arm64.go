package gomonkey

import "unsafe"

func movX(opc, shift int, val uintptr) []byte {
	var m uint32 = 27          // rd
	m |= uint32(val) << 5      // imm16
	m |= uint32(shift&3) << 21 // hw
	m |= 0b100101 << 23        // const
	m |= uint32(opc&0x3) << 29 // opc
	m |= 0b1 << 31             // sf

	res := make([]byte, 4)
	*(*uint32)(unsafe.Pointer(&res[0])) = m

	return res
}

func buildJmpDirective(targetAddr uintptr) []byte {
	targetFuncAddr := *(*uintptr)(unsafe.Pointer(targetAddr)) //func address
	res := make([]byte, 0, 12)
	res = append(res, movX(0b10, 0, targetFuncAddr&0xffff)...)     //movz x27, addr[16:]
	res = append(res, movX(0b11, 1, targetFuncAddr>>16&0xffff)...) //movk x27, addr[32:16]
	res = append(res, []byte{0x60, 0x03, 0x1f, 0xd6}...)           //br  x27
	return res
}
