package dilithium

//#define SHAKE128_RATE 168
const SHAKE128_RATE uint32 = 168

//#define SHAKE256_RATE 136
const SHAKE256_RATE uint32 = 136

//#define SHA3_256_RATE 136
const SHA3_256_RATE uint32 = 136

//#define SHA3_512_RATE 72
const SHA3_512_RATE uint32 = 72

const KECCAK_STATE_SIZE uint32 = 25

const STREAM128_BLOCKBYTES = SHAKE128_RATE
const STREAM256_BLOCKBYTES = SHAKE256_RATE

type Constants struct {
	// Global
	SEEDBYTES     uint16
	CRHBYTES      uint16
	TRBYTES       uint16
	RNDBYTES      uint16
	N             uint16
	Q             int32
	D             uint16
	ROOT_OF_UNITY uint64
	// different per mode
	K           uint16
	L           uint16
	ETA         uint16
	TAU         uint16
	BETA        uint16
	GAMMA1      int32
	GAMMA2      int32
	OMEGA       uint16
	CTILDEBYTES uint16
	//
	MONT int32
	QINV int64

	// other
	POLYT1_PACKEDBYTES   uint16
	POLYT0_PACKEDBYTES   uint16
	POLYVECH_PACKEDBYTES uint16
	POLYZ_PACKEDBYTES    uint16
	POLYW1_PACKEDBYTES   uint16
	POLYETA_PACKEDBYTES  uint16

	CRYPTO_PUBLICKEYBYTES uint32
	CRYPTO_SECRETKEYBYTES uint32
	CRYPTO_BYTES          uint32

	POLY_UNIFORM_NBLOCKS         uint32
	POLY_UNIFORM_ETA_NBLOCKS     uint32
	POLY_UNIFORM_GAMMA1_NBLOCKS  uint32
	DILITHIUM_RANDOMIZED_SIGNING uint32

	// PY ADD
	SAMPLE_PRIMITIVE     uint32
	SAMPLE_PRIMITIVE_INV uint32
	ROOT_POWERS          []int32
	ROOT_POWERS_INV      []int32
}

func (obj *Constants) InitConstants(strMode string) {
	obj.SEEDBYTES = 32
	obj.CRHBYTES = 64
	obj.TRBYTES = 64
	obj.RNDBYTES = 32
	obj.N = 256
	obj.Q = 8380417
	obj.D = 13
	obj.ROOT_OF_UNITY = 1753

	obj.MONT = -4186625 // 2^32 % Q
	obj.QINV = 58728449 // q^(-1) mod 2^32

	obj.ETA = 2
	obj.TAU = 39
	obj.BETA = obj.ETA * obj.TAU
	obj.GAMMA1 = (1 << 17)
	obj.GAMMA2 = ((obj.Q - 1) / 88)
	obj.OMEGA = 80
	obj.CTILDEBYTES = 32
	//
	obj.POLYZ_PACKEDBYTES = 576
	obj.POLYW1_PACKEDBYTES = 192
	obj.POLYETA_PACKEDBYTES = 96

	switch strMode {
	case "2py": // set here difference from above if any
		obj.K = 4
		obj.L = 4
		obj.TRBYTES = 48
	case "3": // set here difference from above if any
		obj.K = 6
		obj.L = 5
	default: // set here difference from above if any
		obj.K = 4
		obj.L = 4
	}
	obj.POLYT1_PACKEDBYTES = 320
	obj.POLYT0_PACKEDBYTES = 416
	obj.POLYVECH_PACKEDBYTES = obj.OMEGA + uint16(obj.K)

	obj.CRYPTO_PUBLICKEYBYTES = uint32(obj.SEEDBYTES + uint16(obj.K)*obj.POLYT1_PACKEDBYTES)
	obj.CRYPTO_SECRETKEYBYTES = uint32(2*obj.SEEDBYTES + obj.TRBYTES + uint16(obj.L)*obj.POLYETA_PACKEDBYTES + uint16(obj.K)*obj.POLYETA_PACKEDBYTES + uint16(obj.K)*obj.POLYT0_PACKEDBYTES)
	obj.CRYPTO_BYTES = uint32(obj.CTILDEBYTES + uint16(obj.L)*obj.POLYZ_PACKEDBYTES + obj.POLYVECH_PACKEDBYTES)

	/*
		obj.POLY_UNIFORM_NBLOCKS = ((768 + STREAM128_BLOCKBYTES - 1) / STREAM128_BLOCKBYTES)

		if obj.ETA == 2 {
			obj.POLY_UNIFORM_ETA_NBLOCKS = ((136 + STREAM256_BLOCKBYTES - 1) / STREAM256_BLOCKBYTES)
		} else if obj.ETA == 4 {
			obj.POLY_UNIFORM_ETA_NBLOCKS = ((227 + STREAM256_BLOCKBYTES - 1) / STREAM256_BLOCKBYTES)
		}

		obj.POLY_UNIFORM_GAMMA1_NBLOCKS = ((uint32(obj.POLYZ_PACKEDBYTES) + STREAM256_BLOCKBYTES - 1) / STREAM256_BLOCKBYTES)
	*/
	obj.DILITHIUM_RANDOMIZED_SIGNING = 0

	obj.SAMPLE_PRIMITIVE = find_primitive_root()
	obj.SAMPLE_PRIMITIVE_INV = find_primitive_root_inv()
	obj.ROOT_POWERS = make([]int32, C.N)
	obj.ROOT_POWERS_INV = make([]int32, C.N)
	build_root_powers(obj)
}

func powm(a, b, m int64) int64 {
	if b < 0 {
		panic("Negative power given")
	}

	var result int64 = 1
	var power int64 = a

	for b > 0 {
		if b&1 == 1 {
			result = result * power % m
		}
		power = power * power % m
		b >>= 1
	}

	return result
}

// Check if g is a primitive n-th root of unity modulo q
func is_primitive_root(g uint32) bool {
	var k uint32
	var p int64 = powm(int64(g), int64(C.N), int64(C.Q))

	if p != 1 {
		return false
	}
	// Check that g^k != 1 for 1 <= k < n
	for k = 1; k < uint32(C.N); k++ {
		p = powm(int64(g), int64(k), int64(C.Q))
		if p == 1 {
			return false
		}
	}
	return true
}

// Find a primitive n-th root of unity modulo q
func find_primitive_root() uint32 {
	// n must divide q-1
	var g uint32
	for g = 2; g < uint32(C.Q); g++ {
		if is_primitive_root(g) {
			return g
		}
	}
	panic("No primitive {n}-th root of unity found modulo {q}")
}

func find_primitive_root_inv() uint32 {
	var p uint32 = uint32(powm(int64(C.SAMPLE_PRIMITIVE), int64(C.Q-2), int64(C.Q)))
	return p
}

func build_root_powers(obj *Constants) {
	var i uint32
	var p int64
	obj.ROOT_POWERS[0] = 1
	obj.ROOT_POWERS_INV[0] = 1
	for i = 1; i < uint32(C.N); i++ {
		p = (int64(obj.ROOT_POWERS[i-1]) * int64(obj.SAMPLE_PRIMITIVE)) % int64(obj.Q)
		obj.ROOT_POWERS[i] = int32(p)
		p = (int64(obj.ROOT_POWERS_INV[i-1]) * int64(obj.SAMPLE_PRIMITIVE_INV)) % int64(obj.Q)
		obj.ROOT_POWERS_INV[i] = int32(p)
	}
}

var C Constants
