//go:build !debug

package dilithium

import (
	"fmt"
	"math/bits"
	"reflect"

	"golang.org/x/crypto/sha3"
)

// compute the Hamming weight of a polynomial (number of nonzero coefficients).
func weight(h *polyveck) int32 {
	var sum int32 = 0
	var k, i uint32

	for k = 0; k < uint32(C.K); k++ {
		for i = 0; i < uint32(C.N); i++ {
			sum += h.vec[k].coeffs[i]
		}
	}
	return sum
}

/*
// multiply 2 numbers modulo a number (without overflow)
func mod_multiply_q(l int32, r int32) int32 {
	var result int32 = 0
	l = l % C.Q
	for r > 0 {
		// If right_int is odd, add left_int to the result
		if r%2 == 1 {
			result = (result + l) % C.Q
		}

		// Double left_int and halve right_int
		l = (l * 2) % C.Q
		r /= 2
	}
	return result
}
*/

func sample_in_ball(c_tilde []uint8, p *poly) {

	var i, j, byte_index, current_sign_pos uint32
	random_bytes := make([]uint8, 320)
	sign_bits := make([]uint8, C.TAU)

	// Absorb ˜c into SHAKE-256 to create a random intitial stream of 32 bytes
	hl := sha3.NewShake256()
	hl.Write(c_tilde)
	hl.Read(random_bytes)

	// Extract τ sign bits from the first 8 bytes
	for i = 0; i < uint32(C.TAU); i++ {
		//sign_bits = [(random_bytes[i // 8] >> (i % 8)) & 1 for i in range(tau)]
		sign_bits[i] = random_bytes[i/8] >> (i % 8) & 1
	}

	//# Start reading after the first 8 bytes
	byte_index = 8
	for i = uint32(C.N) - uint32(C.TAU); i < uint32(C.N); i++ {
		current_sign_pos = 0
		for 2 > 1 {
			if byte_index >= uint32(len(random_bytes)) {
				// Extend the random stream if needed
				random_bytes_add := make([]uint8, 64)
				hl.Read(random_bytes_add)
				random_bytes = append(random_bytes, random_bytes_add...)
			}
			j = uint32(random_bytes[byte_index])
			byte_index += 1
			// Rejection sampling condition
			if j <= i {
				p.coeffs[i] = p.coeffs[j]
				if sign_bits[current_sign_pos] == 1 {
					p.coeffs[j] = 1
				} else {
					p.coeffs[j] = -1
				}
				current_sign_pos += 1
				break
			}
		}
	}
	//return poly
}

/*
func ntt_vector_k(v *polyveck) {
	var k uint32
	for k = 0; k < uint32(C.K); k++ {
		ntt_py(&v.vec[k], C.ROOT_POWERS)
	}
}
*/

func ntt_vector_l(v *polyvecl) {
	var l uint32
	for l = 0; l < uint32(C.L); l++ {
		ntt_py(&v.vec[l], C.ROOT_POWERS)
	}
}

func ntt_matrix(m []polyvecl) {
	var k, l uint64
	for k = 0; k < uint64(C.K); k++ {
		for l = 0; l < uint64(C.L); l++ {
			ntt_py(&m[k].vec[l], C.ROOT_POWERS)
		}
	}
}

func ntt_py(p *poly, root_powers []int32) {

	//var bits = 8
	var polyreverse poly
	polyreverse.Init()
	var ir uint8
	var i, j, w uint32
	var u, v int32
	var plen, half, root uint32

	// Ensure Python integers
	for i = 0; i < uint32(C.N); i++ {
		ir = bits.Reverse8(uint8(i))
		polyreverse.coeffs[i] = p.coeffs[ir]
	}

	// 0x007FE001
	plen = 2
	for plen <= uint32(C.N) {
		half = plen / 2
		root = uint32(root_powers[uint32(C.N)/plen])
		i = 0
		for i < uint32(C.N) {
			w = 1
			for j = 0; j < half; j++ {
				u = polyreverse.coeffs[i+j]
				v = int32(int64(int64(polyreverse.coeffs[i+j+half])*int64(w)) % int64(C.Q))
				if v < 0 {
					v += C.Q
				}
				polyreverse.coeffs[i+j] = (u + v) % C.Q
				if polyreverse.coeffs[i+j] < 0 {
					polyreverse.coeffs[i+j] += C.Q
				}

				polyreverse.coeffs[i+j+half] = (u - v) % C.Q
				if polyreverse.coeffs[i+j+half] < 0 {
					polyreverse.coeffs[i+j+half] += C.Q
				}
				w = uint32(int64(int64(int64(w)*int64(root)) % int64(C.Q)))
			}
			i += plen
		}
		plen *= 2
	}

	for i = 0; i < uint32(C.N); i++ {
		p.coeffs[i] = polyreverse.coeffs[i]
	}
}

func intt_py(p *poly) {

	var i uint32
	var poly_inv uint32 = uint32(powm(int64(C.N), int64(C.Q-2), int64(C.Q)))
	ntt_py(p, C.ROOT_POWERS_INV)

	//# Normalize by multiplying by poly_inv
	for i = 0; i < uint32(C.N); i++ {
		p.coeffs[i] = int32(int64(int64(p.coeffs[i])*int64(poly_inv)) % int64(C.Q))
	}
}

func expand_matrix_A(m []polyvecl, seed []uint8) {
	var k, l, n, offset uint64
	var vv uint32
	//var current_hash int32 = 0
	var hash []uint8 = make([]uint8, 3*C.N+2)
	var added []uint8 = make([]uint8, 2)
	for k = 0; k < uint64(C.K); k++ {
		for l = 0; l < uint64(C.L); l++ {
			h := sha3.NewShake128()
			added[0] = uint8(l)
			added[1] = uint8(k)
			h.Write(seed)
			h.Write(added)
			h.Read(hash)

			//big endian
			/*
				offset = uint64(len(hash)) - 1
				for n = 0; n < uint64(C.N); n++ {

					vv = uint32(hash[offset]) + (uint32(hash[offset-1]) << 8) + (uint32(hash[offset-2]) << 16)
					vv &= 0x007FFFFF
					m[k].vec[l].coeffs[n] = int32(vv)
					offset -= 3
				}
			*/
			offset = uint64(0)
			for n = 0; n < uint64(C.N); n++ {

				vv = uint32(hash[offset]) + (uint32(hash[offset+1]) << 8) + (uint32(hash[offset+2]) << 16)
				vv &= 0x007FFFFF
				m[k].vec[l].coeffs[n] = int32(vv)
				offset += 3
			}

		}
	}
}

// Verify for Python mode 2,3,5 --> 2py, 3py, 5py
func Verify(pk []uint8, msg []uint8, sigma []uint8) bool {

	var i, k, l uint64
	var pr, sum int32
	c_tilde := make([]uint8, C.CTILDEBYTES)
	c_tilde_check := make([]uint8, C.CTILDEBYTES)
	tr := make([]uint8, C.TRBYTES)
	mu := make([]uint8, C.TRBYTES)
	w1_bytes := make([]uint8, C.K*C.POLYW1_PACKEDBYTES)
	var c poly
	c.Init()
	var z polyvecl
	z.Init()
	var h polyveck
	h.Init()
	rho := make([]uint8, C.SEEDBYTES)
	var t1 polyveck
	t1.Init()
	var t1_scaled polyveck
	t1_scaled.Init()
	var A []polyvecl = make([]polyvecl, C.K)
	for i = 0; i < uint64(C.K); i++ {
		A[i].Init()
	}
	var Az polyveck
	Az.Init()
	var w polyveck
	w.Init()
	var w1 polyveck
	w1.Init()

	if unpack_sig(c_tilde, &z, &h, sigma) > 0 {
		fmt.Println("Problem with unpacking sigma")
		return false
	}

	unpack_pk(rho, &t1, pk)

	// Step 31: Verify conditions
	// (1) ∥z∥∞ < γ1 - β
	// NB: We are doing the check here in advance, because later structure z will bew manilupated for calculations
	if polyvecl_chknorm(&z, int32(C.GAMMA1)-int32(C.BETA)) > 0 {
		fmt.Println("Problem with z")
		return false
	}

	expand_matrix_A(A, rho)
	ntt_matrix(A)
	//# 384 bits = 48 bytes

	// Step 28: Compute µ = CRH(CRH(ρ || t1) || M)
	//Inner CRH - concatenate rho and t1
	sha3.ShakeSum256(tr, pk) // 384 bits = 48 bytes

	// 384 bits = 48 bytes
	// Outer CRH (384 bits)
	hl := sha3.NewShake256()
	hl.Write(tr)
	hl.Write(msg)
	hl.Read(mu)

	// Step 29: Compute c from SampleInBall
	sample_in_ball(c_tilde, &c)

	// Step 30: Compute w'1
	// Transform z into NTT
	ntt_vector_l(&z)

	// Transform c into NTT
	ntt_py(&c, C.ROOT_POWERS)

	// Transform t1·2^d into NTT
	for k = 0; k < uint64(C.K); k++ {
		for i = 0; i < uint64(C.N); i++ {
			t1_scaled.vec[k].coeffs[i] = int32(uint64(uint64(t1.vec[k].coeffs[i])*uint64((1<<C.D))) % uint64(C.Q))
		}
	}

	for k = 0; k < uint64(C.K); k++ {
		ntt_py(&t1_scaled.vec[k], C.ROOT_POWERS)
	}

	// Compute Az in NTT
	var plocal poly
	plocal.Init()
	for i = 0; i < uint64(C.N); i++ {
		for k = 0; k < uint64(C.K); k++ {
			sum = 0
			for l = 0; l < uint64(C.L); l++ {
				pr = int32(int64(A[k].vec[l].coeffs[i]) * int64(z.vec[l].coeffs[i]) % int64(C.Q))
				sum = (sum + pr) % C.Q
			}
			Az.vec[k].coeffs[i] = sum
		}
	}

	// Compute w = Az - c·t1 in NTT and transform back
	for k = 0; k < uint64(C.K); k++ {
		for i = 0; i < uint64(C.N); i++ {
			pr = int32(int64(c.coeffs[i]) * int64(t1_scaled.vec[k].coeffs[i]) % int64(C.Q))
			w.vec[k].coeffs[i] = int32((Az.vec[k].coeffs[i] - pr) % C.Q)
			if w.vec[k].coeffs[i] < 0 {
				w.vec[k].coeffs[i] += C.Q
			}
		}
	}

	// transforming w_ntt to n
	for k = 0; k < uint64(C.K); k++ {
		intt_py(&w.vec[k])
	}

	// UseHint to compute w'1
	polyveck_use_hint(&w1, &w, &h)

	// (2) Verify that c_tilde = H(µ || w'1)
	pack_w1(w1_bytes, &w1)

	hw := sha3.NewShake256()
	hw.Write(mu)
	hw.Write(w1_bytes)
	hw.Read(c_tilde_check)

	//c_tilde_check ?= c_tilde
	if !reflect.DeepEqual(c_tilde, c_tilde_check) {

		/*
			fmt.Print("c_tilde=")
			fmt.Printf("%02x", c_tilde)
			fmt.Println()
			fmt.Print("c_tilde_check=")
			fmt.Printf("%02x", c_tilde_check)
			fmt.Println()
		*/
		fmt.Println("c_tilde false")
		return false
	}

	// (3) Verify that the number of 1's in h ≤ ω
	if weight(&h) > int32(C.OMEGA) {
		fmt.Println("Weight false")
		return false
	}

	// All checks passed msg is verified
	return true
}
