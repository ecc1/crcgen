package main

//go:generate ../crcgen/crcgen

func crc8(msg []byte) byte {
	res := byte(0)
	for _, b := range msg {
		res = crc8Table[res^b]
	}
	return res
}
