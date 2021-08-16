package gomonkey

import "unsafe"

func buildJmpDirective(double uintptr) []byte {
	res := make([]byte, 0, 24)
	d1d2 := double & 0xFFFF
	d3d4 := double >> 16 & 0xFFFF
	d5d6 := double >> 32 & 0xFFFF
	d7d8 := double >> 48 & 0xFFFF

	res = append(res, movImm(0B10, 0, d1d2)...) // MOVZ x26, double[16:0]
	res = append(res, movImm(0B11, 1, d3d4)...) // MOVK x26, double[32:16]
	res = append(res, movImm(0B11, 2, d5d6)...) // MOVK x26, double[48:32]
	res = append(res, movImm(0B11, 3, d7d8)...) // MOVK x26, double[64:48]
	res = append(res, []byte{0x4A, 0x03, 0x40, 0xF9}...)  // LDR x10, [x26]
	res = append(res, []byte{0x40, 0x01, 0x1F, 0xD6}...)  // BR x10

	return res
}

func movImm(opc, shift int, val uintptr) []byte {
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
