/*
# requirements
# n ≥ 1
# q ≥ 2
# q = poly(λ) - means q should be big so its hard to break
# m ≥ nk ≥ n
# k = ⌈log2q⌉
# m = ¯m+2nk
*/

package updatable_encryption

import (
	"fmt"
	"math"
)

type Constants struct {
	n            int
	q            int
	k            int
	isexactpower bool
	nk           int
	m_tilda      int
	m            int
	alpha        float64
	G            [][]int

	mean       float64
	sigma_sk   float64
	edisp      float64
	tau_sample float64

	//testing
	s_hidden  []int
	e0_hidden []int
	e1_hidden []int
	e2_hidden []int
	e00       []int
	e01       []int
	e02       []int
}

func (obj *Constants) InitConstants() {
	obj.n = 16
	obj.q = 1048573 //8380417
	obj.k = int(math.Ceil(math.Log2(float64(obj.q))))
	obj.isexactpower = (obj.q == int(math.Pow(float64(2), float64(obj.k))))
	obj.nk = obj.n * obj.k
	obj.m_tilda = 32
	obj.m = obj.m_tilda + 2*obj.nk
	obj.alpha = float64(1) / float64(10)
	obj.G = generate_G(obj.n, obj.k, obj.q, obj.nk)

	obj.mean = 0
	obj.sigma_sk = 0.7
	obj.edisp = 0.7
	obj.tau_sample = 2 / (2 + 1)

	obj.s_hidden = make([]int, obj.n)
	obj.e0_hidden = make([]int, obj.m_tilda)
	obj.e1_hidden = make([]int, obj.nk)
	obj.e2_hidden = make([]int, obj.nk)
	obj.e00 = make([]int, obj.m_tilda)
	obj.e01 = make([]int, obj.nk)
	obj.e02 = make([]int, obj.nk)
}

func generate_G(n, k, q, nk int) [][]int {
	gT := make([]int, k)
	for i := 0; i < k; i++ {
		gT[i] = int(math.Pow(2, float64(i)))
	}

	G := make([][]int, n)
	for i := 0; i < n; i++ {
		G[i] = make([]int, nk)
		for j := 0; j < k; j++ {
			G[i][i*k+j] = gT[j] % q
		}
	}
	return G
}

var C Constants

func update_error(e0, e1, e2 []int) error {

	if len(C.e0_hidden) != len(e0) || len(C.e1_hidden) != len(e1) || len(C.e2_hidden) != len(e2) {
		return fmt.Errorf("update_error - missmatching sizes")
	}

	for i := range e0 {
		C.e0_hidden[i] = e0[i]
	}

	for i := range e1 {
		C.e1_hidden[i] = e1[i]
	}

	for i := range e2 {
		C.e2_hidden[i] = e2[i]
	}

	return nil
}

func update_error2(e0, e1, e2 []int, M [][]int) error {
	helpe1 := append(append(C.e0_hidden, C.e1_hidden...), C.e2_hidden...)

	helper2, err2 := dot_product_VM(helpe1, M)

	if err2 != nil {
		return fmt.Errorf("update_error2 - err2 is not nil")
	}

	if len(helper2) != C.m {
		return fmt.Errorf("update_error2 - missmatching sizes")
	}

	//helper2 = modQ(helper2)

	for i := range e0 {
		C.e00[i] = helper2[i]
	}

	for i := range e1 {
		C.e01[i] = helper2[i+C.m_tilda]
	}

	for i := range e2 {
		C.e02[i] = helper2[i+C.m_tilda+C.nk]
	}

	helper3, err3 := add_vectors(helper2[:C.m_tilda], e0)
	helper4, err4 := add_vectors(helper2[C.m_tilda:C.m_tilda+C.nk], e1)
	helper5, err5 := add_vectors(helper2[C.m_tilda+C.nk:], e2)

	if err3 != nil || err4 != nil || err5 != nil {
		return fmt.Errorf("update_error2 - err3 or err4 or err5 is not nil")
	}

	//helper3 = modQ(helper3)
	//helper4 = modQ(helper4)
	//helper5 = modQ(helper5)

	if len(C.e0_hidden) != len(helper3) || len(C.e1_hidden) != len(helper4) || len(C.e2_hidden) != len(helper5) {
		return fmt.Errorf("update_error2 - missmatching sizes 2")
	}

	for i := range C.e0_hidden {
		C.e0_hidden[i] = helper3[i]
	}

	for i := range C.e1_hidden {
		C.e1_hidden[i] = helper4[i]
	}

	for i := range C.e2_hidden {
		C.e2_hidden[i] = helper5[i]
	}

	return nil
}

func get_error() ([]int, []int, []int) {
	return C.e0_hidden, C.e1_hidden, C.e2_hidden
}

func get_error2() ([]int, []int, []int) {
	return C.e00, C.e01, C.e02
}

func updateS(s []int) error {

	if len(s) != len(C.s_hidden) {
		return fmt.Errorf("update_s - missmatching sizes")
	}
	for i := range C.s_hidden {
		C.s_hidden[i] = (C.s_hidden[i] + s[i]) % C.q
	}

	//C.s_hidden = s
	return nil
}

func getS() []int {
	return C.s_hidden
}
