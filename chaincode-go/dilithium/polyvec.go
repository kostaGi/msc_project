package dilithium

/* Vectors of polynomials of length L */
type polyvecl struct {
	//poly vec[L]
	vec []poly
}

func (obj *polyvecl) Init() {
	obj.vec = make([]poly, C.L)
	for i := 0; i < int(C.L); i++ {
		obj.vec[i].Init()
	}
}

/* Vectors of polynomials of length K */
type polyveck struct {
	//poly vec[K]
	vec []poly
}

func (obj *polyveck) Init() {
	obj.vec = make([]poly, C.K)
	for i := 0; i < int(C.K); i++ {
		obj.vec[i].Init()
	}
}

/**************************************************************/
/************ Vectors of polynomials of length L **************/
/**************************************************************/

/*************************************************
* Name:        polyvecl_chknorm
*
* Description: Check infinity norm of polynomials in vector of length L.
*              Assumes input polyvecl to be reduced by polyvecl_reduce().
*
* Arguments:   - const polyvecl *v: pointer to vector
*              - int32_t B: norm bound
*
* Returns 0 if norm of all polynomials is strictly smaller than B <= (Q-1)/8
* and 1 otherwise.
**************************************************/
func polyvecl_chknorm(v *polyvecl, bound int32) int {
	var i uint16

	for i = 0; i < C.L; i++ {
		if poly_chknorm(&v.vec[i], bound) == 1 {
			return 1
		}
	}

	/*
	  for(i = 0; i < L; ++i)
	    if(poly_chknorm(&v->vec[i], bound))
	      return 1;
	*/

	return 0
}

/**************************************************************/
/************ Vectors of polynomials of length K **************/
/**************************************************************/

/*************************************************
* Name:        polyveck_use_hint
*
* Description: Use hint vector to correct the high bits of input vector.
*
* Arguments:   - polyveck *w: pointer to output vector of polynomials with
*                             corrected high bits
*              - const polyveck *u: pointer to input vector
*              - const polyveck *h: pointer to input hint vector
**************************************************/
func polyveck_use_hint(w *polyveck, u *polyveck, h *polyveck) {
	var i uint16

	for i = 0; i < C.K; i++ {
		poly_use_hint(&w.vec[i], &u.vec[i], &h.vec[i])
	}
	/*
	  for(i = 0; i < K; ++i)
	    poly_use_hint(&w->vec[i], &u->vec[i], &h->vec[i]);
	*/
}
