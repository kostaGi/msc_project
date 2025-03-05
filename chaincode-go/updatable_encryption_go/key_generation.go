package main

import "fmt"

func key_gen(R1 [][]int) ([][]int, [][]int, [][]int, [][]int, error) {
	A0 := sample_uniform_matrix(C.n, C.m_tilda)
	var R1_new [][]int
	if R1 == nil {
		R1_new = sample_normal_matrix(C.m_tilda, C.nk, C.mean, C.sigma_sk)
		VerifyR(R1_new)
	} else {
		R1_new = R1
	}

	R2 := sample_normal_matrix(C.m_tilda, C.nk, C.mean, C.sigma_sk)
	VerifyR(R2)

	A1, err1 := dot_product_MM(A0, R1_new)

	if err1 != nil {
		return nil, nil, nil, nil, fmt.Errorf("key_gen - err1 %v", err1)
	}

	A1 = negate_matrix_values(A1)

	A2, err2 := dot_product_MM(A0, R2)

	if err2 != nil {
		return nil, nil, nil, nil, fmt.Errorf("key_gen - err2 %v", err2)
	}

	A2 = negate_matrix_values(A2)

	/*
			# to be removed
		    # verify trapdoor
		    R = np.block([R1, R2])
		    pk_local = np.block([A0, A1, A2])
		    R_ext = np.block([[R], [np.eye(R.shape[1], dtype=int)] ])
		    Res = np.dot(pk_local, R_ext) % q
		    all_zeros = not np.any(Res)
		    assert all_zeros
	*/

	return A0, A1, A2, R1_new, nil
}
