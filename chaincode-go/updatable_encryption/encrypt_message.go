package updatable_encryption

import (
	"fmt"
)

/*
func compact_error(e_vec []int) []int {
	qhalf := math.Round(float64(C.q) / 2)

	for i := range e_vec {
		if e_vec[i] > int(qhalf) {
			e_vec[i] = e_vec[i] - C.q
		}
	}

	return e_vec
}


func compact_error2(e_vec []int) []int {
	qhalf := math.Round(float64(C.q) / 2)

	for i := range e_vec {
		if e_vec[i] > int(qhalf) {
			e_vec[i] = e_vec[i] - C.q
		}
	}

	for i := range e_vec {
		if e_vec[i] < -int(qhalf) {
			e_vec[i] = e_vec[i] + C.q
		}
	}

	return e_vec
}
*/

func generate_error_vector() ([]int, []int, []int) {
	e0 := sample_normal_matrix(1, C.m_tilda, 0, C.edisp)[0]
	VerifyE(e0)
	//#d = calculate_d(e0, m_tilda, alpha, q, sigma)
	e1 := sample_normal_matrix(1, C.nk, 0, C.edisp)[0]
	VerifyE(e1)
	e2 := sample_normal_matrix(1, C.nk, 0, C.edisp)[0]
	VerifyE(e2)

	return e0, e1, e2
}

func generate_error_vector_multiple(size int) ([]int, []int, []int) {
	e0 := sample_normal_matrix(1, C.m_tilda, 0, C.edisp)[0]
	VerifyE(e0)
	//#d = calculate_d(e0, m_tilda, alpha, q, sigma)
	e1 := sample_normal_matrix(1, C.nk, 0, C.edisp)[0]
	VerifyE(e1)
	e2 := sample_normal_matrix(1, size, 0, C.edisp)[0]
	VerifyE(e2)

	return e0, e1, e2
}

/*
# calculate d
def calculate_d(e0, m_tilda, alpha, q, sigma):
    sum = 0
    for value in e0:
        sum+= value**2
    #print("sum=", sum, "m_tilda=", m_tilda*((alpha*q)**2), "sigma=", sigma**2)
    d = np.sqrt(sum + m_tilda*((alpha*q)**2) * (sigma**2))
    return d
*/

func Encrypt(A0, A1, A2 [][]int, encoded_message []int, H_mu_input [][]int, updateMatrix [][]int, updateWithMatrix bool) ([][]int, []int, error) {

	var H_mu [][]int

	if H_mu_input != nil {
		H_mu = H_mu_input
	} else {
		H_mu = generate_invertible_matrixQ(C.n)
	}

	if H_mu == nil {
		return nil, nil, fmt.Errorf("encrypt - H_mu")
	}

	HmG, err1 := dot_product_MM(H_mu, C.G)

	if err1 != nil {
		return nil, nil, fmt.Errorf("encrypt - err1")
	}

	helper1, err2 := add_matrix(A1, HmG)

	if err2 != nil {
		return nil, nil, fmt.Errorf("encrypt - err2")
	}

	A_mu, err3 := concatenate_matrices_row(A0, helper1)

	if err3 != nil {
		return nil, nil, fmt.Errorf("encrypt - err3")
	}

	A_mu, err4 := concatenate_matrices_row(A_mu, A2)

	if err4 != nil {
		return nil, nil, fmt.Errorf("encrypt - err4")
	}

	s := sample_random_vector(C.n, C.q)

	err_s := updateS(s)

	if err_s != nil {
		return nil, nil, fmt.Errorf("encrypt - err_s")
	}

	e0, e1, e2 := generate_error_vector()

	if updateWithMatrix {
		err_e1 := update_error2(e0, e1, e2, updateMatrix)
		fmt.Printf("We are here!")
		if err_e1 != nil {
			return nil, nil, fmt.Errorf("encrypt - errUpdate1")
		}
	} else {
		err_e2 := update_error(e0, e1, e2)

		if err_e2 != nil {
			return nil, nil, fmt.Errorf("encrypt - errUpdate2")
		}
	}

	helper2, err5 := dot_product_VM(s, A_mu)

	if err5 != nil {
		return nil, nil, fmt.Errorf("encrypt - err5")
	}

	helper3 := concatenate_vector(concatenate_vector(e0, e1), e2)
	helper4 := concatenate_vector(concatenate_vector(make([]int, C.m_tilda), make([]int, C.nk)), encoded_message)

	b, err6 := add_vectors(helper2, helper3)

	if err6 != nil {
		return nil, nil, fmt.Errorf("encrypt - err6")
	}

	b, err7 := add_vectors(b, helper4)

	if err7 != nil {
		return nil, nil, fmt.Errorf("encrypt - err6")
	}

	return H_mu, b, nil
}

func Encrypt_multiple(A0, A1, A2 [][]int, encoded_message []int, H_mu_input [][]int, size int) ([][]int, []int, error) {

	var H_mu [][]int

	if H_mu_input != nil {
		H_mu = H_mu_input
	} else {
		H_mu = generate_invertible_matrixQ(C.n)
	}

	if H_mu == nil {
		return nil, nil, fmt.Errorf("encrypt - H_mu")
	}

	HmG, err1 := dot_product_MM(H_mu, C.G)

	if err1 != nil {
		return nil, nil, fmt.Errorf("encrypt - err1")
	}

	helper1, err2 := add_matrix(A1, HmG)

	if err2 != nil {
		return nil, nil, fmt.Errorf("encrypt - err2")
	}

	A_mu, err3 := concatenate_matrices_row(A0, helper1)

	if err3 != nil {
		return nil, nil, fmt.Errorf("encrypt - err3")
	}

	full_blocks := size / C.nk
	if size%C.nk != 0 {
		full_blocks += 1
	}

	for range full_blocks {
		A_mu, _ = concatenate_matrices_row(A_mu, A2)

		/*
			if err4 != nil {
				return nil, nil, fmt.Errorf("encrypt - err4")
			}
		*/
	}

	s := sample_random_vector(C.n, C.q)

	err_s := updateS(s)

	if err_s != nil {
		return nil, nil, fmt.Errorf("encrypt - err_s")
	}

	e0, e1, e2 := generate_error_vector_multiple(size)

	helper2, err5 := dot_product_VM(s, A_mu)

	if err5 != nil {
		return nil, nil, fmt.Errorf("encrypt - err5")
	}

	helper3 := concatenate_vector(concatenate_vector(e0, e1), e2)
	helper4 := concatenate_vector(concatenate_vector(make([]int, C.m_tilda), make([]int, C.nk)), encoded_message)

	b, err6 := add_vectors(helper2, helper3)

	if err6 != nil {
		return nil, nil, fmt.Errorf("encrypt - err6")
	}

	b, err7 := add_vectors(b, helper4)

	if err7 != nil {
		return nil, nil, fmt.Errorf("encrypt - err6")
	}

	return H_mu, b, nil
}
