package dilithium

/*************************************************
* Name:        power2round
*
* Description: For finite field element a, compute a0, a1 such that
*              a mod^+ Q = a1*2^D + a0 with -2^{D-1} < a0 <= 2^{D-1}.
*              Assumes a to be standard representative.
*
* Arguments:   - int32_t a: input element
*              - int32_t *a0: pointer to output element a0
*
* Returns a1.
*************************************************
func power2round(a0 *int32, a int32) int32 {
	var a1 int32

	a1 = (a + (1 << (C.D - 1)) - 1) >> C.D
	*a0 = a - (a1 << C.D)
	return a1
}
*/

/*************************************************
* Name:        decompose
*
* Description: For finite field element a, compute high and low bits a0, a1 such
*              that a mod^+ Q = a1*ALPHA + a0 with -ALPHA/2 < a0 <= ALPHA/2 except
*              if a1 = (Q-1)/ALPHA where we set a1 = 0 and
*              -ALPHA/2 <= a0 = a mod^+ Q - Q < 0. Assumes a to be standard
*              representative.
*
* Arguments:   - int32_t a: input element
*              - int32_t *a0: pointer to output element a0
*
* Returns a1.
**************************************************/
func decompose(a0 *int32, a int32) int32 {
	var a1 int32

	a1 = (a + 127) >> 7
	if C.GAMMA2 == (C.Q-1)/32 {
		a1 = (a1*1025 + (1 << 21)) >> 22
		a1 &= 15
	} else if C.GAMMA2 == (C.Q-1)/88 {
		a1 = (a1*11275 + (1 << 23)) >> 24
		a1 ^= ((43 - a1) >> 31) & a1
	}

	*a0 = a - a1*2*C.GAMMA2
	*a0 -= (((C.Q-1)/2 - *a0) >> 31) & C.Q
	return a1
}

/*************************************************
* Name:        make_hint
*
* Description: Compute hint bit indicating whether the low bits of the
*              input element overflow into the high bits.
*
* Arguments:   - int32_t a0: low bits of input element
*              - int32_t a1: high bits of input element
*
* Returns 1 if overflow.
*************************************************
func make_hint(a0 int32, a1 int32) int32 {
	if a0 > C.GAMMA2 || a0 < -C.GAMMA2 || (a0 == -C.GAMMA2 && a1 != 0) {
		return 1
	}

	return 0
}
*/

/*************************************************
* Name:        use_hint
*
* Description: Correct high bits according to hint.
*
* Arguments:   - int32_t a: input element
*              - unsigned int hint: hint bit
*
* Returns corrected high bits.
**************************************************/
func use_hint(a int32, hint int32) int32 {
	var a0 int32
	var a1 int32 = decompose(&a0, a)

	if hint == 0 {
		return a1
	}

	if C.GAMMA2 == (C.Q-1)/32 {
		if a0 > 0 {
			return (a1 + 1) & 15
		} else {
			return (a1 - 1) & 15
		}
	} else if C.GAMMA2 == (C.Q-1)/88 {
		if a0 > 0 {
			if a1 == 43 {
				return 0
			} else {
				return a1 + 1
			}
		} else {
			if a1 == 0 {
				return 43
			} else {
				return a1 - 1
			}
		}
	}
	return 111110 // this should not be reached
}
