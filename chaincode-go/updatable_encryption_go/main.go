package main

import (
	"fmt"
	"reflect"
)

func main() {
	C.InitConstants()
	msg1 := sample_random_vector(C.nk, 2)

	A0, A1, A2, R1, err1 := key_gen(nil)

	if err1 != nil {
		fmt.Printf("error: %v\n", err1)
		panic("main - err1")
	}

	Hmu, b, err2 := encrypt(A0, A1, A2, encode(msg1), nil, nil, false)

	if err2 != nil {
		fmt.Printf("error: %v\n", err2)
		panic("main - err2")
	}

	e0, e1, _ := get_error()
	err_lema1 := lemma1Check(R1, concatenate_vector(e0, e1))

	if err_lema1 != nil {
		fmt.Printf("error: %v\n", err_lema1)
		panic("main - err err_lema1")
	}
	err_equal1 := equalityCheck(A0, A1, R1, Hmu)

	if err_equal1 != nil {
		fmt.Printf("error: %v\n", err_equal1)
		panic("main - err err_equal1")
	}

	msg2, err3 := decrypt(R1, Hmu, b, A0, A1, A2)

	if err3 != nil {
		fmt.Printf("error: %v\n", err3)
		panic("main - err3")
	}

	if !reflect.DeepEqual(msg1, msg2) {
		panic("main - msg1 is not decrypted correctly as msg2")
	}

	current_b := b
	cA0, cA1, cA2 := A0, A1, A2
	current_Hmu := Hmu

	for counter1 := range 10 {

		print("-------------------- UPDATE PHASE", counter1, "--------------------")
		nA0, nA1, nA2, _, err4 := key_gen(R1)

		if err4 != nil {
			fmt.Printf("error: %v\n", err4)
			panic("main - err4")
		}

		tM, tb_zero_message, tHmu_prime, err5 := token_gen(cA0, cA1, cA2, R1, nA0, nA1, nA2, current_Hmu)

		if err5 != nil {
			fmt.Printf("error: %v\n", err5)
			panic("main - err5")
		}

		b_prime, H_mu_prime, err6 := update(tM, tb_zero_message, tHmu_prime, current_b)

		e0, e1, _ := get_error()
		err_lema1 := lemma1Check(R1, concatenate_vector(e0, e1))

		if err_lema1 != nil {
			fmt.Printf("error: %v\n", err_lema1)
			panic("main - err err_lema recursive")
		}

		if err6 != nil {
			fmt.Printf("error: %v\n", err6)
			panic("main - err6")
		}

		new_msg, err7 := decrypt(R1, H_mu_prime, b_prime, nA0, nA1, nA2)

		if err7 != nil {
			fmt.Printf("error: %v\n", err7)
			panic("main - err7")
		}

		if !reflect.DeepEqual(msg1, new_msg) {
			panic("main - msg1 is not decrypted correctly as new_msg")
		}

		current_b = b_prime
		cA0, cA1, cA2 = nA0, nA1, nA2
		current_Hmu = H_mu_prime
	}
}
