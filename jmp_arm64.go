package gomonkey

import "unsafe"

func movX(opc, shift int, val uintptr) []byte {
	var m uint32 = 26          // rd
	m |= uint32(val) << 5      // imm16
	m |= uint32(shift&3) << 21 // hw
	m |= 0b100101 << 23        // const
	m |= uint32(opc&0x3) << 29 // opc
	m |= 0b1 << 31             // sf

	res := make([]byte, 4)
	*(*uint32)(unsafe.Pointer(&res[0])) = m

	return res
}

func buildJmpDirective(targetFuncAddr uintptr) []byte {
	res := make([]byte, 0, 24)
	res = append(res, movX(0b10, 0, targetFuncAddr&0xffff)...)     //movz x26, addr[16:]
	res = append(res, movX(0b11, 1, targetFuncAddr>>16&0xffff)...) //movk x26, addr[32:16]
	res = append(res, movX(0b11, 2, targetFuncAddr>>32&0xffff)...) //movk x26, addr[48:32]
	res = append(res, movX(0b11, 3, targetFuncAddr>>48&0xffff)...) //movk x26, addr[64:48]
	res = append(res, []byte{0x4a, 0x03, 0x40, 0xf9}...)           //ldr  x10, [x26]
	res = append(res, []byte{0x40, 0x01, 0x1f, 0xd6}...)           //br   x10
	return res
}
