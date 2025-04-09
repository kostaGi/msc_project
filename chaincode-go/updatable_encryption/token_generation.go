package updatable_encryption

import (
	"fmt"
	"math"
	"reflect"
)

func SampleDO(R1, Amu01, Hmu, A_part_prime [][]int, tau float64) ([][]int, error) {
	scalling_d := make([]int, C.n)
	for i := range C.n {
		scalling_d[i] = 1
	}

	x := make([][]int, len(A_part_prime[0]))

	for i := range x {
		x[i] = make([]int, len(Amu01[0]))
	}

	H_hardcoded_inv := inverse_matrixQ(Hmu, len(Hmu))
	A_part_prime_t := transpose_matrix(A_part_prime)

	for i := range len(A_part_prime[0]) {
		p, err0 := sample_perturbation_v2(R1, tau, len(Amu01[0]))

		if err0 != nil {
			return nil, fmt.Errorf("SampleDO - err0, %v", err0)
		}

		u_tmp, err1 := dot_product_MV(Amu01, p)

		if err1 != nil {
			return nil, fmt.Errorf("SampleDO - err1, %v", err1)
		}

		u_tmp, err2 := add_vectors(A_part_prime_t[i], negate_vector_values(u_tmp))

		if err2 != nil {
			return nil, fmt.Errorf("SampleDO - err1, %v", err2)
		}

		v, err3 := dot_product_MV(H_hardcoded_inv, u_tmp)

		if err3 != nil {
			return nil, fmt.Errorf("SampleDO - err1, %v", err3)
		}

		z := oracle_sampleD(v, C.q, C.k, 2, C.isexactpower)

		R_ext, err4 := concatenate_matrices_col(R1, identity_matrix(C.nk))

		if err4 != nil {
			return nil, fmt.Errorf("SampleDO - err1, %v", err4)
		}

		x_vec, err5 := dot_product_MV(R_ext, z)

		if err5 != nil {
			return nil, fmt.Errorf("SampleDO - err1, %v", err5)
		}

		x_vec, err6 := add_vectors(p, x_vec)

		if err6 != nil {
			return nil, fmt.Errorf("SampleDO - err1, %v", err6)
		}

		for j := range x_vec {
			x[i][j] = x_vec[j]
		}
	}
	x = transpose_matrix(x)

	return x, nil
}

func Token_gen(A0, A1, A2, R1, A0_prime, A1_prime, A2_prime, Hmu [][]int) ([][]int, []int, [][]int, error) {
	Hmu_prime := generate_invertible_matrixQ(C.n)
	Hmu_primeG, err1 := dot_product_MM(Hmu_prime, C.G)

	if err1 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err2, %v", err1)
	}

	helper1, err2 := add_matrix(A1_prime, Hmu_primeG)

	if err2 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err2, %v", err2)
	}

	A_mu_prime, err3 := concatenate_matrices_row(A0_prime, helper1)

	if err3 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err3, %v", err3)
	}

	A_mu_prime, err4 := concatenate_matrices_row(A_mu_prime, A2_prime)

	if err4 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err4, %v", err4)
	}

	HmuG, err5 := dot_product_MM(Hmu, C.G)

	if err5 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err5, %v", err5)
	}

	helper2, err6 := add_matrix(A1, HmuG)

	if err6 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err6, %v", err6)
	}

	Amu01, err7 := concatenate_matrices_row(A0, helper2)

	if err7 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err7, %v", err7)
	}

	X0, err8 := SampleDO(R1, Amu01, Hmu, A0_prime, C.tau_sample)

	if err8 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err8, %v", err8)
	}

	X1, err9 := SampleDO(R1, Amu01, Hmu, helper1, C.tau_sample*math.Sqrt(float64(C.m_tilda)/2))

	if err9 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err9, %v", err9)
	}

	helper_4, err_4 := add_matrix(A2_prime, negate_matrix_values(A2))

	if err_4 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err_4, %v", err_4)
	}

	X2, err10 := SampleDO(R1, Amu01, Hmu, helper_4, C.tau_sample*math.Sqrt(float64(C.m_tilda)/2))

	//X1 = SampleDO(sk, A01m, Hmu, A1_prime + Hmu_primeG, tau_sample*np.sqrt(m_tilda/2))
	//X2 = SampleDO(sk, A01m, Hmu, A2_prime-A2, tau_sample*np.sqrt(m_tilda/2))

	if err10 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err10, %v", err10)
	}

	V0, err11 := dot_product_MM(Amu01, X0)

	V0 = over_QM(V0)

	if err11 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err11, %v", err11)
	}

	V1, err12 := dot_product_MM(Amu01, X1)

	V1 = over_QM(V1)

	if err12 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err12, %v", err12)
	}

	V2, err13 := dot_product_MM(Amu01, X2)

	V2 = over_QM(V2)

	if err13 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err13, %v", err13)
	}

	if !reflect.DeepEqual(V0, A0_prime) {
		return nil, nil, nil, fmt.Errorf("token_gen - V0 is not equal to A0_prime")
	}

	/*
		helper3, err14 := concatenate_matrices_row(A1_prime, Hmu_primeG)

		if err14 != nil {
			return nil, nil, nil, fmt.Errorf("token_gen - err14, %v", err14)
		}
	*/

	helper1 = over_QM(helper1)

	if !reflect.DeepEqual(V1, helper1) {
		return nil, nil, nil, fmt.Errorf("token_gen - V1 is not equal to A1_prime + Hmu_primeG")
	}

	helper4, err15 := add_matrix(A2_prime, negate_matrix_values(A2))

	if err15 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err15, %v", err15)
	}

	helper4 = over_QM(helper4)

	if !reflect.DeepEqual(V2, helper4) {

		//fmt.Println("V2", V2)
		//fmt.Println("helper4", helper4)

		return nil, nil, nil, fmt.Errorf("token_gen - V2 is not equal to A2_prime - A2")
	}

	helper5, err16 := concatenate_matrices_row(X0, X1)

	if err16 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err16, %v", err16)
	}

	helper5, err17 := concatenate_matrices_row(helper5, X2)

	if err17 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err17, %v", err17)
	}

	zeroMatrix1 := make([][]int, C.nk)
	for i := range zeroMatrix1 {
		zeroMatrix1[i] = make([]int, C.m_tilda)
	}

	zeroMatrix2 := make([][]int, C.nk)
	for i := range zeroMatrix2 {
		zeroMatrix2[i] = make([]int, C.nk)
	}

	helper6, err18 := concatenate_matrices_row(zeroMatrix1, zeroMatrix2)

	if err18 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err18, %v", err18)
	}

	helper6, err19 := concatenate_matrices_row(helper6, identity_matrix(C.nk))

	if err19 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err19, %v", err19)
	}

	M, err20 := concatenate_matrices_col(helper5, helper6)

	if err20 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err20, %v", err20)
	}

	Amu, err21 := concatenate_matrices_row(Amu01, A2)

	if err21 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err21, %v", err21)
	}

	helper7, err22 := dot_product_MM(Amu, M)

	if err22 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err22, %v", err22)
	}

	//func encrypt(A0, A1, A2 [][]int, encoded_message []int, updateMatrix [][]int, updateWithMatrix bool) ([][]int, []int, error)

	if !reflect.DeepEqual(over_QM(helper7), over_QM(A_mu_prime)) {
		return nil, nil, nil, fmt.Errorf("token_gen - dot_product_MM(Amu, M) is not equal to A_mu_prime")
	}

	_, b_zero_message, err23 := Encrypt(A0_prime, A1_prime, A2_prime, make([]int, C.nk), Hmu_prime, M, true)

	if err23 != nil {
		return nil, nil, nil, fmt.Errorf("token_gen - err23, %v", err23)
	}

	return M, b_zero_message, Hmu_prime, nil
}

func Update(M [][]int, b_zero_message []int, Hmu_prime [][]int, b []int) ([]int, [][]int, error) {
	//e00, e01, e02 := get_error2()

	//helper1 := concatenate_vector(concatenate_vector(e00, e01), e02)

	b_prime, err1 := dot_product_VM(b, M)

	if err1 != nil {
		return nil, nil, fmt.Errorf("update - err1, %v", err1)
	}

	b_prime, err2 := add_vectors(b_prime, b_zero_message)

	if err2 != nil {
		return nil, nil, fmt.Errorf("update - err2, %v", err2)
	}

	/*
		b_prime, err3 := add_vectors(b_prime, negate_vector_values(helper1))

		if err3 != nil {
			return nil, nil, fmt.Errorf("update - err3, %v", err3)
		}
	*/

	return b_prime, Hmu_prime, nil
}

func Update2(X0, X1, X2 [][]int, b_zero_message []int, b []int, message_size int) ([]int, error) {

	full_blocks_count := message_size / C.nk
	if message_size%C.nk != 0 {
		full_blocks_count += 1
	}

	M, err1 := CreateM(X0, X1, X2, full_blocks_count)

	if err1 != nil {
		return nil, fmt.Errorf("Update2 - err1, %v", err1)
	}

	b_prime, err2 := dot_product_VM(b, M)

	if err2 != nil {
		return nil, fmt.Errorf("Update2 - err2, %v", err2)
	}

	b_prime, err3 := add_vectors(b_prime, b_zero_message)

	if err3 != nil {
		return nil, fmt.Errorf("Update2 - err3, %v", err3)
	}

	return b_prime, nil
}

func CreateM(X0, X1, X2 [][]int, full_blocks_count int) ([][]int, error) {
	M, err1 := concatenate_matrices_row(X0, X1)

	if err1 != nil {
		return nil, fmt.Errorf("CreateM - err1, %v", err1)
	}

	for range full_blocks_count {
		helper1, err2 := concatenate_matrices_row(M, X2)

		if err2 != nil {
			return nil, fmt.Errorf("CreateM - err2, %v", err2)
		}

		M = helper1
	}

	zeroMatrix1 := make([][]int, C.nk)
	for i := range zeroMatrix1 {
		zeroMatrix1[i] = make([]int, C.m_tilda)
	}

	zeroMatrix2 := make([][]int, C.nk)
	for i := range zeroMatrix2 {
		zeroMatrix2[i] = make([]int, C.nk)
	}

	zero_matrix_left := 0
	for range full_blocks_count {
		helper2, _ := concatenate_matrices_row(zeroMatrix1, zeroMatrix2)

		for range zero_matrix_left {
			zeroMatrix_3 := make([][]int, C.nk)
			for i := range zeroMatrix_3 {
				zeroMatrix_3[i] = make([]int, C.nk)
			}
			helper2, _ = concatenate_matrices_row(helper2, zeroMatrix_3)
		}

		helper2, _ = concatenate_matrices_row(helper2, identity_matrix(C.nk))

		for range full_blocks_count - zero_matrix_left - 1 {

			zeroMatrix_4 := make([][]int, C.nk)
			for i := range zeroMatrix_4 {
				zeroMatrix_4[i] = make([]int, C.nk)
			}
			helper2, _ = concatenate_matrices_row(helper2, zeroMatrix_4)
		}

		zero_matrix_left += 1
		M, _ = concatenate_matrices_col(M, helper2)
	}

	return M, nil
}
