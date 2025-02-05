package dilithium

type poly struct {
	coeffs []int32 //[N]   slice
}

func (obj *poly) Init() {
	obj.coeffs = make([]int32, C.N)
}

/*************************************************
 * Name:        poly_use_hint
 *
 * Description: Use hint polynomial to correct the high bits of a polynomial.
 *
 * Arguments:   - poly *b: pointer to output polynomial with corrected high bits
 *              - const poly *a: pointer to input polynomial
 *              - const poly *h: pointer to input hint polynomial
 **************************************************/
func poly_use_hint(b *poly, a *poly, h *poly) {
	var i uint16
	//DBENCH_START();

	for i = 0; i < C.N; i++ {
		b.coeffs[i] = use_hint(a.coeffs[i], h.coeffs[i])
	}
	//for(i = 0; i < N; ++i)
	//  b->coeffs[i] = use_hint(a->coeffs[i], h->coeffs[i]);

	//DBENCH_STOP(*tround);
}

/*************************************************
 * Name:        poly_chknorm
 *
 * Description: Check infinity norm of polynomial against given bound.
 *              Assumes input coefficients were reduced by reduce32().
 *
 * Arguments:   - const poly *a: pointer to polynomial
 *              - int32_t B: norm bound
 *
 * Returns 0 if norm is strictly smaller than B <= (Q-1)/8 and 1 otherwise.
 **************************************************/
func poly_chknorm(a *poly, B int32) int {
	var i uint16
	var t int32
	//DBENCH_START();

	if B > (C.Q-1)/8 {
		return 1
	}

	/* It is ok to leak which coefficient violates the bound since
	   the probability for each coefficient is independent of secret
	   data but we must not leak the sign of the centralized representative. */
	for i = 0; i < C.N; i++ {
		// Absolute value
		t = a.coeffs[i] >> 31
		t = a.coeffs[i] - (t & 2 * a.coeffs[i])

		if t >= B {
			//DBENCH_STOP(*tsample);
			return 1
		}
	}

	/*
		for(i = 0; i < N; ++i) {
		  // Absolute value
		  t = a->coeffs[i] >> 31;
		  t = a->coeffs[i] - (t & 2*a->coeffs[i]);

		  if(t >= B) {
			//DBENCH_STOP(*tsample);
			return 1;
		  }
		}
	*/

	//DBENCH_STOP(*tsample);
	return 0
}
