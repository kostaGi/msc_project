package dilithium

func unpack_sig(c []uint8, //[CTILDEBYTES]
	z *polyvecl,
	h *polyveck,
	sig []uint8) int { // [CRYPTO_BYTES]

	var i uint16
	var sigoffset uint64
	sigoffset = 0
	//unsigned int i, j, k;

	// c_tilda + z + h
	var bytes_to_check uint32 = uint32(C.CTILDEBYTES) + uint32(C.POLYZ_PACKEDBYTES)*uint32(C.L)

	//safety check
	if len(sig) <= int(bytes_to_check) {
		return 1
	}

	for i = 0; i < C.CTILDEBYTES; i++ {
		c[i] = sig[i]
	}
	//sig += C.CTILDEBYTES
	sigoffset += uint64(C.CTILDEBYTES)

	for i = 0; i < C.L; i++ {
		unpack_z(&z.vec[i], sig, int(sigoffset)+int(i)*int(C.POLYZ_PACKEDBYTES))
	}
	//sig += C.L * C.POLYZ_PACKEDBYTES
	sigoffset += uint64(C.L) * uint64(C.POLYZ_PACKEDBYTES)
	//fmt.Println("z=", z)

	/* Decode h */
	unpack_h(sig, sigoffset, h)
	return 0
}

func unpack_h(sig []uint8, offset uint64, h *polyveck) {

	//fmt.Println(offset, len(sig))
	var pos uint64 = offset
	//var hlen uint8 = 0
	var i uint64
	var j uint8

	i = 0
	for i < uint64(C.K) {
		//fmt.Println("i=", i, "pos=", pos)
		for sig[pos] < sig[pos+1] {
			j = sig[pos]
			h.vec[i].coeffs[j] = 1
			pos++
		}
		if i < uint64(C.K) {
			// not last
			j = sig[pos]
			h.vec[i].coeffs[j] = 1
			pos++
		}
		i++
	}
	//fmt.Println("h=", h)
}

func unpack_pk(rho []uint8, //SEEDBYTES
	t1 *polyveck,
	pk []uint8) {
	var i uint64
	var pkoffset uint64
	pkoffset = 0

	for i = 0; i < uint64(C.SEEDBYTES); i++ {
		rho[i] = pk[pkoffset+i]
	}
	pkoffset += uint64(C.SEEDBYTES)

	for i = 0; i < uint64(C.K); i++ {
		unpack_t1(&t1.vec[i], pk, pkoffset+i*uint64(C.POLYT1_PACKEDBYTES))
	}
}

func unpack_t1(r *poly, a []uint8, offset uint64) {
	var i uint64

	for i = 0; i < uint64(C.N)/4; i++ {
		r.coeffs[4*i+0] = int32((uint32(a[offset+5*i+0]>>0) | ((uint32)(a[offset+5*i+1]) << 8)) & 0x3FF)
		r.coeffs[4*i+1] = int32((uint32(a[offset+5*i+1]>>2) | ((uint32)(a[offset+5*i+2]) << 6)) & 0x3FF)
		r.coeffs[4*i+2] = int32((uint32(a[offset+5*i+2]>>4) | ((uint32)(a[offset+5*i+3]) << 4)) & 0x3FF)
		r.coeffs[4*i+3] = int32((uint32(a[offset+5*i+3]>>6) | ((uint32)(a[offset+5*i+4]) << 2)) & 0x3FF)
	}
}

func unpack_z(r *poly, a []uint8, offset int) {
	var i uint16

	if C.GAMMA1 == (1 << 17) {
		for i = 0; i < C.N/4; i++ {
			r.coeffs[4*i+0] = int32(a[offset+9*int(i)+0])
			r.coeffs[4*i+0] |= (int32)(a[offset+9*int(i)+1]) << 8
			r.coeffs[4*i+0] |= (int32)(a[offset+9*int(i)+2]) << 16
			r.coeffs[4*i+0] &= 0x3FFFF

			r.coeffs[4*i+1] = int32(a[offset+9*int(i)+2] >> 2)
			r.coeffs[4*i+1] |= (int32)(a[offset+9*int(i)+3]) << 6
			r.coeffs[4*i+1] |= (int32)(a[offset+9*int(i)+4]) << 14
			r.coeffs[4*i+1] &= 0x3FFFF

			r.coeffs[4*i+2] = int32(a[offset+9*int(i)+4] >> 4)
			r.coeffs[4*i+2] |= (int32)(a[offset+9*int(i)+5]) << 4
			r.coeffs[4*i+2] |= (int32)(a[offset+9*int(i)+6]) << 12
			r.coeffs[4*i+2] &= 0x3FFFF

			r.coeffs[4*i+3] = int32(a[offset+9*int(i)+6] >> 6)
			r.coeffs[4*i+3] |= (int32)(a[offset+9*int(i)+7]) << 2
			r.coeffs[4*i+3] |= (int32)(a[offset+9*int(i)+8]) << 10
			r.coeffs[4*i+3] &= 0x3FFFF

			r.coeffs[4*i+0] = C.GAMMA1 - r.coeffs[4*i+0]
			r.coeffs[4*i+1] = C.GAMMA1 - r.coeffs[4*i+1]
			r.coeffs[4*i+2] = C.GAMMA1 - r.coeffs[4*i+2]
			r.coeffs[4*i+3] = C.GAMMA1 - r.coeffs[4*i+3]
		}
	} else if C.GAMMA1 == (1 << 19) {
		for i = 0; i < C.N/2; i++ {
			r.coeffs[2*i+0] = int32(a[offset+5*int(i)+0])
			r.coeffs[2*i+0] |= (int32)(a[offset+5*int(i)+1]) << 8
			r.coeffs[2*i+0] |= (int32)(a[offset+5*int(i)+2]) << 16
			r.coeffs[2*i+0] &= 0xFFFFF

			r.coeffs[2*i+1] = int32(a[offset+5*int(i)+2] >> 4)
			r.coeffs[2*i+1] |= (int32)(a[offset+5*int(i)+3]) << 4
			r.coeffs[2*i+1] |= (int32)(a[offset+5*int(i)+4]) << 12
			/* r->coeffs[2*i+1] &= 0xFFFFF; */ /* No effect, since we're anyway at 20 bits */

			r.coeffs[2*i+0] = C.GAMMA1 - r.coeffs[2*i+0]
			r.coeffs[2*i+1] = C.GAMMA1 - r.coeffs[2*i+1]
		}
	}
}

func pack_w1_vec(r []uint8, offset int, a *poly) {
	var i uint16
	if C.GAMMA2 == (C.Q-1)/88 {
		for i = 0; i < C.N/4; i++ {
			r[offset+3*int(i)+0] = uint8(a.coeffs[4*i+0])
			r[offset+3*int(i)+0] |= uint8(a.coeffs[4*i+1] << 6)
			r[offset+3*int(i)+1] = uint8(a.coeffs[4*i+1] >> 2)
			r[offset+3*int(i)+1] |= uint8(a.coeffs[4*i+2] << 4)
			r[offset+3*int(i)+2] = uint8(a.coeffs[4*i+2] >> 4)
			r[offset+3*int(i)+2] |= uint8(a.coeffs[4*i+3] << 2)
		}
	} else if C.GAMMA2 == (C.Q-1)/32 {
		for i = 0; i < C.N/2; i++ {
			r[offset+int(i)] = uint8(a.coeffs[2*i+0] | (a.coeffs[2*i+1] << 4))
		}
	}
}

func pack_w1(r []uint8, w1 *polyveck) {
	var i uint16

	for i = 0; i < C.K; i++ {
		pack_w1_vec(r, int(i)*int(C.POLYW1_PACKEDBYTES), &w1.vec[i])
	}
}
