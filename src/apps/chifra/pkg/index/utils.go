package index

import (
	"encoding/binary"
	"log"

	"github.com/ethereum/go-ethereum/common"
)

// WhichBits returns the five bits calculated from an address used to determine if the address is
// in the bloom filter. We get the five bits by cutting the 20-byte address into five equal four-byte
// parts, turning those four bytes into an 32-bit integer modulo the width of a bloom array item.
func WhichBits(addr common.Address) (bits [5]uint32) {
	slice := addr.Bytes()
	if len(slice) != 20 {
		log.Fatal("address is not 20 bytes long - should not happen")
	}

	cnt := 0
	for i := 0; i < len(slice); i += 4 {
		bytes := slice[i : i+4]
		bits[cnt] = (binary.BigEndian.Uint32(bytes) % uint32(BLOOM_WIDTH_IN_BITS))
		cnt++
	}

	return
}
