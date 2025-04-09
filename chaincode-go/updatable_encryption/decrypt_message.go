package updatable_encryption

import (
	"fmt"
	"reflect"
)

func invertO(R, Amu01 [][]int, b01 []int, H_2 [][]int) ([]int, []int, error) {

	R_padded, err1 := concatenate_matrices_col(R, identity_matrix(C.nk))

	if err1 != nil {
		return nil, nil, fmt.Errorf("invertO - err1, %v", err1)
	}

	b01_hat, err2 := dot_product_VM(b01, R_padded)

	if err2 != nil {
		return nil, nil, fmt.Errorf("invertO - err2, %v", err2)
	}

	//print(n, k, nk, "b_0_1=",len(b_0_1), ", R=", len(R), "x", len(R[0]), ", R_padded=", len(R_padded), "x", len(R_padded[0]), ", b_0_1_hat=", len(b_0_1_hat))

	s_m := make([]int, C.n)

	for i := range s_m {
		s_m[i] = oracle_sampleO(b01_hat[i*C.k : C.k*(i+1)])
	}
	//print("s_m=", s_m)

	H_2_inv := inverse_matrixQ(H_2, len(H_2))

	s_calc, err3 := dot_product_VM(s_m, H_2_inv)

	if err3 != nil {
		return nil, nil, fmt.Errorf("invertO - err3, %v", err3)
	}

	fmt.Println(C.m, C.nk, C.m_tilda, C.n, C.k, len(Amu01), len(Amu01[0]))
	helper1, err4 := dot_product_VM(s_calc, Amu01)

	if err4 != nil {
		return nil, nil, fmt.Errorf("invertO - err4, %v", err4)
	}

	e_calc, err5 := add_vectors(b01, negate_vector_values(helper1))

	if err5 != nil {
		return nil, nil, fmt.Errorf("invertO - err5, %v", err5)
	}

	fmt.Printf("s_calc: %v\n", s_calc)
	fmt.Printf("s_hidd: %v\n", getS())

	if !reflect.DeepEqual(getS(), s_calc) {
		return nil, nil, fmt.Errorf("invertO - s_hidden is not equal to s_calc")
	}

	e0_hidden, e1_hidden, _ := get_error()
	modQ(e_calc)
	e0_calc := (e_calc[:C.m_tilda])

	if !reflect.DeepEqual(e0_hidden, e0_calc) {
		return nil, nil, fmt.Errorf("invertO - e0_hidden is not equal to e0_calc")
	}

	e1_calc := e_calc[C.m_tilda:]

	if !reflect.DeepEqual(e1_hidden, e1_calc) {
		return nil, nil, fmt.Errorf("invertO - e1_hidden is not equal to e1_calc")
	}

	return s_calc, e_calc, nil
}

func Decrypt(R1, Hmu [][]int, b []int, A0, A1, A2 [][]int) ([]int, error) {

	_, _, e2_hidden := get_error()

	helper1, err1 := dot_product_MM(Hmu, C.G)

	if err1 != nil {
		return nil, fmt.Errorf("decrypt - err1, %v", err1)
	}

	helper2, err2 := add_matrix(A1, helper1)

	if err2 != nil {
		return nil, fmt.Errorf("decrypt - err2, %v", err2)
	}

	helper3, err3 := concatenate_matrices_row(A0, helper2)

	//fmt.Println("H=", len(helper3), len(helper3[0]), len(helper1), len(helper1[0]))
	if err3 != nil {
		return nil, fmt.Errorf("decrypt - err3, %v", err3)
	}

	//e01_calc, err10
	s_calc, _, err10 := invertO(R1, helper3, b[:C.m-C.nk], Hmu)

	if err10 != nil {
		return nil, fmt.Errorf("decrypt - err10, %v", err10)
	}

	helper4, err4 := dot_product_VM(s_calc, A2)

	if err4 != nil {
		return nil, fmt.Errorf("decrypt - err4, %v", err4)
	}

	helper5, err5 := add_vectors(b[C.m_tilda+C.nk:], negate_vector_values(helper4))

	if err5 != nil {
		return nil, fmt.Errorf("decrypt - err5, %v", err5)
	}
	//e2_calc
	decoded_message, e2_calc := decode(helper5)

	e2_calc = modQ(e2_calc)

	if !reflect.DeepEqual(e2_hidden, e2_calc) {
		return nil, fmt.Errorf("invertO - e1_hidden is not equal to e2_calc")
	}

	return decoded_message, nil

}

func Decrypt_multiple(R1, Hmu [][]int, b []int, A0, A1, A2 [][]int, size int) ([]int, error) {

	//_, _, e2_hidden := get_error()

	helper1, err1 := dot_product_MM(Hmu, C.G)

	if err1 != nil {
		return nil, fmt.Errorf("decrypt - err1, %v", err1)
	}

	helper2, err2 := add_matrix(A1, helper1)

	if err2 != nil {
		return nil, fmt.Errorf("decrypt - err2, %v", err2)
	}

	helper3, err3 := concatenate_matrices_row(A0, helper2)

	//fmt.Println("H=", len(helper3), len(helper3[0]), len(helper1), len(helper1[0]))
	if err3 != nil {
		return nil, fmt.Errorf("decrypt - err3, %v", err3)
	}

	//e01_calc, err10
	s_calc, _, err10 := invertO(R1, helper3, b[:C.m_tilda+C.nk], Hmu)

	if err10 != nil {
		return nil, fmt.Errorf("decrypt - err10, %v", err10)
	}

	helper4, err4 := dot_product_VM(s_calc, A2)

	if err4 != nil {
		return nil, fmt.Errorf("decrypt - err4, %v", err4)
	}

	decoded_message := make([]int, size)
	full_blocks := size / C.nk

	for counter1 := range full_blocks {
		helper5, err5 := add_vectors(b[C.m_tilda+C.nk+counter1*C.nk:C.m_tilda+C.nk+(counter1+1)*C.nk], negate_vector_values(helper4))
		if err5 != nil {
			return nil, fmt.Errorf("decrypt - err5, %v", err5)
		}
		decoded_message_block, _ := decode(helper5)

		if counter1+1 == full_blocks {
			for counter2 := range size % C.nk {
				decoded_message_block[counter1*C.nk+counter2] = decoded_message_block[counter2]
			}
		} else {
			for counter2 := range C.nk {
				decoded_message_block[counter1*C.nk+counter2] = decoded_message_block[counter2]
			}

		}

		/*
			e2_calc = modQ(e2_calc)

			if !reflect.DeepEqual(e2_hidden, e2_calc) {
				return nil, fmt.Errorf("invertO - e1_hidden is not equal to e2_calc")
			}
		*/
	}

	return decoded_message, nil
}
