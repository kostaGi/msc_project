package main

import (
	"math"
	"math/rand"
)

func base_b_decomposition(v, b, k int) []int {
	digits := make([]int, k)

	for i := range digits {
		digits[i] = v & (b - 1)
		v = int(v / b)
	}

	return digits
}

// https://escholarship.org/uc/item/8b40w7r8 - Algorithm 11:DECODEG(v,b,r[q]kb)
func oracle_sampleO_not_exact_power(b_hat []int) int {
	q_bits := base_b_decomposition(C.q, 2, C.k)

	/*
		q_bits := make([]int, C.k)
		q_local := C.q
		for i := range C.k {
			q_bits[i] = q_local & 1
			q_local = q_local >> 1
		}
	*/

	v := b_hat

	for i := range C.k - 1 {
		v[i] = 2*v[i] - v[i+1]
	}
	v[C.k-1] = 2 * v[C.k-1]

	x := make([]int, C.k)

	reg := 0

	for i := range C.k - 1 {
		x[i] = int(math.Round(float64(v[i]) / float64(C.q)))
		reg = (reg / 2) + (int(math.Pow(float64(2), float64((C.k-1)))) * q_bits[i])
		v[C.k-1] = v[C.k-1] + x[i]*reg
	}

	x[C.k-1] = int(math.Round(float64(v[C.k-1]) / math.Pow(float64(2), float64((C.k)))))
	s := x[C.k-1]
	reg = 0
	i := C.k - 2

	for i >= 0 {
		reg = 2*reg + q_bits[i+1]
		s = s + x[i]*reg
		i = i - 1
	}

	if s < 0 {
		s = s + C.q
	}

	return s
}

func oracle_sampleO_exact_power(b_hat []int) int {
	br1 := C.q / 4
	br2 := (3 * C.q) / 4
	blen := len(b_hat)
	s := 0
	s_bit := 0
	i := blen - 1

	for i >= 0 {

		//int(math.Pow(float64(x), float64(y))))
		bb := (b_hat[i] - s*int(math.Pow(float64(2), float64(i)))) % C.q
		if bb >= br1 && bb <= br2 {
			s_bit = 1
		} else {
			s_bit = 0
		}

		s = (s + int(math.Pow(float64(2), float64((blen-1-i))))*s_bit) % C.q
		i = i - 1
	}

	return s
}

func oracle_sampleO(b_hat []int) int {
	if !C.isexactpower {
		return oracle_sampleO_not_exact_power(b_hat)
	}
	return oracle_sampleO_exact_power(b_hat)
}

// # https://escholarship.org/uc/item/8b40w7r8 - Algorithm 9: g−1(u) for q = b^k
func oracle_sampleD_exact_power(u, k, b int) []int {
	x := make([]int, k)

	u_local := u

	for i := range x {
		y := u_local % b

		if y == 0 {
			x[i] = 0
		} else {
			if float64(y)/float64(b) > 0.5 {
				x[i] = y - b
			} else {
				x[i] = y
			}
		}

		// u_local = (u_local - x[i]) / b ???
		u_local = (u_local - x[i]) / b
	}

	return x
}

func oracle_sampleD_not_exact_power(q int, q_bits []int, u int, u_bits []int, k, b int) []int {
	x := make([]int, k)
	y := make([]int, k)

	if rand.Float64() < float64(q-u)/float64(q) {
		x[k-1] = 0
	} else {
		x[k-1] = -1
	}

	u_local := u
	q_local := q
	i := k - 2

	for i >= 0 {
		u_local = u_local - u_bits[i+1]*int(math.Pow(float64(b), float64(i+1)))
		q_local = q_local - q_bits[i+1]*int(math.Pow(float64(b), float64(i+1)))
		c := -(u_local + (x[k-1] * q_local))

		var p, z int

		if c < 0 {
			p = c + int(math.Pow(float64(b), float64(i+1)))
			z = -1
		} else {
			p = c
			z = 0
		}

		if rand.Float64() < float64(p)/math.Pow(float64(b), float64(i+1)) {
			x[i] = z + 1

		} else {
			x[i] = z
		}
		i = i - 1
	}

	for i := range k - 1 {
		if i == 0 {
			y[i] = b*x[i] + x[k-1]*q_bits[i] + u_bits[i]
		} else {
			y[i] = b*x[i] - x[i-1] + x[k-1]*q_bits[i] + u_bits[i]

		}
	}
	y[k-1] = -x[k-2] + x[k-1]*q_bits[k-1] + u_bits[k-1]

	return y
}

func oracle_sampleD(UV []int, q, k, b int, ixp bool) []int {
	q_bits := base_b_decomposition(q, b, k)
	yarr := make([]int, k*len(UV))

	for ui := range len(UV) {
		uvi := UV[ui] % q

		if uvi < 0 {
			uvi += q
		}

		var y []int
		if ixp {
			y = oracle_sampleD_exact_power(uvi, k, b)
		} else {
			u_bits := base_b_decomposition(uvi, b, k)
			y = oracle_sampleD_not_exact_power(q, q_bits, uvi, u_bits, k, b)
		}

		for i := range k {
			yarr[ui*k+i] = y[i]
		}
	}

	return yarr
}
