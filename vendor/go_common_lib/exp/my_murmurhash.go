package multi_exp

import ()

const m uint32 = 0x5bd1e995
const r uint32 = 24

func MyMurmurhash2(data []byte, seed uint32) uint32 {
	h := seed ^ uint32(len(data))
	var k uint32
	for l := len(data); l >= 4; l -= 4 {
		k = uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24
		k *= m
		k ^= k >> r
		k *= m

		h *= m
		h ^= k
		data = data[4:]
	}
	// Handle the last few bytes of the input array
	switch len(data) {
	case 3:
		h ^= uint32(data[2]) << 16
		fallthrough
	case 2:
		h ^= uint32(data[1]) << 8
		fallthrough
	case 1:
		h ^= uint32(data[0])
		h *= m
	}

	// Do a few final mixes of the hash to ensure the last few bytes are well incorporated
	h ^= h >> 13
	h *= m
	h ^= h >> 15
	return h
}
