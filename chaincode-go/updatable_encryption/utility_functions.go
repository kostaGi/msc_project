package updatable_encryption

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"reflect"
	"time"

	"gonum.org/v1/gonum/mat"
)

func modQ(vector []int) []int {
	for i := range vector {
		vector[i] %= C.q
		for vector[i] > C.q/2 {
			vector[i] -= C.q
		}
		for vector[i] < -C.q/2 {
			vector[i] += C.q
		}
	}
	return vector
}

func negate_matrix_values(matrix [][]int) [][]int {

	negative_matrix := make([][]int, len(matrix))

	for i := range negative_matrix {
		negative_matrix[i] = make([]int, len(matrix[0]))
		for j := range matrix[0] {
			negative_matrix[i][j] = -matrix[i][j]
		}
	}

	return negative_matrix
}

func negate_vector_values(vector []int) []int {

	negative_vector := make([]int, len(vector))

	for i := range negative_vector {
		negative_vector[i] = -vector[i]
	}

	return negative_vector
}

func add_vectors(vector_left []int, vector_right []int) ([]int, error) {

	if len(vector_left) != len(vector_right) {
		return nil, fmt.Errorf("add_vectors - missmatching sizes %v vs %v", len(vector_left), len(vector_right))
	}

	result := make([]int, len(vector_left))

	for i := range vector_left {
		result[i] = vector_left[i] + vector_right[i]
		result[i] %= C.q
	}

	return result, nil
}

func add_matrix(matrix_left [][]int, matrix_right [][]int) ([][]int, error) {

	if len(matrix_left) != len(matrix_right) || len(matrix_left[0]) != len(matrix_right[0]) {
		return nil, fmt.Errorf("add_matrix - missmatching sizes %v vs %v or %v vs %v", len(matrix_left), len(matrix_right), len(matrix_left[0]), len(matrix_right[0]))
	}

	result := make([][]int, len(matrix_left))

	for i := range matrix_left {
		result[i] = make([]int, len(matrix_left[0]))

		for j := range matrix_left[0] {
			result[i][j] = matrix_left[i][j] + matrix_right[i][j]
			result[i][j] %= C.q
		}
	}

	return result, nil
}

func dot_product_MV(matrix [][]int, vector []int) ([]int, error) {
	if len(matrix[0]) != len(vector) {
		return nil, fmt.Errorf("dot_product_MV - missmatching sizes %v vs %v", len(matrix[0]), len(vector))
	}
	result := make([]int, len(matrix))
	for i := range matrix {
		for j := range matrix[i] {
			result[i] += matrix[i][j] * vector[j]
		}
		result[i] %= C.q
	}
	return result, nil
}

func dot_product_VM(vector []int, matrix [][]int) ([]int, error) {
	if len(matrix) != len(vector) {
		return nil, fmt.Errorf("dot_product_VM - missmatching sizes %v vs %v", len(vector), len(matrix))
	}
	result := make([]int, len(matrix[0]))
	for i := range matrix[0] {
		for j := range matrix {
			result[i] += vector[j] * matrix[j][i]
		}
		result[i] %= C.q
	}
	return result, nil
}

func dot_product_MM(matrix_left [][]int, matrix_right [][]int) ([][]int, error) {

	if len(matrix_left[0]) != len(matrix_right) {
		return nil, fmt.Errorf("dot_product_MM - missmatching sizes %v vs %v", len(matrix_left[0]), len(matrix_right))
	}

	result := make([][]int, len(matrix_left))
	for i := range result {
		result[i] = make([]int, len(matrix_right[0]))
	}

	for i := range matrix_left {
		for j := range matrix_right[0] {
			for k := range matrix_right {
				result[i][j] += matrix_left[i][k] * matrix_right[k][j]
			}
			result[i][j] %= C.q
		}
	}
	return result, nil
}

func transpose_matrix(matrix [][]int) [][]int {
	result := make([][]int, len(matrix[0]))

	for i := range result {
		result[i] = make([]int, len(matrix))
		for j := range len(matrix) {
			result[i][j] = matrix[j][i]
		}
	}

	return result
}

func concatenate_vector(vector_left, vector_right []int) []int {
	return append(vector_left, vector_right...)
}

func concatenate_matrices_row(matrix_left, matrix_right [][]int) ([][]int, error) {
	if len(matrix_left) != len(matrix_right) {
		return nil, fmt.Errorf("concatenate_matrices_row - missmatching sizes %v vs %v", len(matrix_left), len(matrix_right))
	}

	result := make([][]int, len(matrix_left))
	for i := range result {
		result[i] = append(matrix_left[i], matrix_right[i]...)
	}
	return result, nil
}

func concatenate_matrices_col(matrix_left, matrix_right [][]int) ([][]int, error) {
	if len(matrix_left[0]) != len(matrix_right[0]) {
		return nil, fmt.Errorf("concatenate_matrices_col - missmatching sizes %v vs %v", len(matrix_left[0]), len(matrix_right[0]))
	}

	result := make([][]int, len(matrix_left)+len(matrix_right))
	for i := range matrix_left {
		result[i] = matrix_left[i]
	}

	for i := range matrix_right {
		result[len(matrix_left)+i] = matrix_right[i]
	}

	return result, nil
}

// Function to generate a random n x n matrix with determinant 1 and not an identity matrix
func generate_matrixQ(matrix_size int) [][]int {
	var factor int
	rand.Seed(time.Now().UnixNano())
	matrix := identity_matrix(matrix_size)

	// Apply random row operations to maintain determinant = 1
	for i := 0; i < matrix_size*3; i++ {
		r1 := rand.Intn(matrix_size)
		r2 := rand.Intn(matrix_size)
		if r1 != r2 {
			factor = int(rand.Intn(int(C.q)-1) + 1) // Ensure factor is non-zero
			for j := 0; j < matrix_size; j++ {
				matrix[r1][j] = (matrix[r1][j] + factor*matrix[r2][j]) % C.q
			}
		}
	}

	return matrix
}

// Function to generate an identity matrix of size n
func identity_matrix(matrix_size int) [][]int {
	matrix := make([][]int, matrix_size)
	for i := range matrix {
		matrix[i] = make([]int, matrix_size)
		matrix[i][i] = 1
	}
	return matrix
}

// Function to compute the modular inverse of a number mod q
func mod_inverseQ(a int) int {
	var x int
	a = a % C.q
	for x = 1; x < C.q; x++ {
		if (a*x)%C.q == 1 {
			return x
		}
	}
	return 1 // Assumes q is prime
}

// Function to compute the inverse of a matrix modulo q
func inverse_matrixQ(mat [][]int, matrix_size int) [][]int {
	inv := identity_matrix(matrix_size)
	augmented := make([][]int, matrix_size)
	for i := range mat {
		augmented[i] = append(mat[i], inv[i]...)
	}

	for i := 0; i < matrix_size; i++ {
		pivot := augmented[i][i]
		invPivot := mod_inverseQ(pivot)
		for j := 0; j < 2*matrix_size; j++ {
			augmented[i][j] = (augmented[i][j] * invPivot) % C.q
		}
		for k := 0; k < matrix_size; k++ {
			if k != i {
				factor := augmented[k][i]
				for j := 0; j < 2*matrix_size; j++ {
					augmented[k][j] = (augmented[k][j] - factor*augmented[i][j]) % C.q
					if augmented[k][j] < 0 {
						augmented[k][j] += C.q
					}
				}
			}
		}
	}

	for i := 0; i < matrix_size; i++ {
		inv[i] = augmented[i][matrix_size:]
	}
	return inv
}

// Function to multiply two matrices mod q
func multiply_matricesQ(a, b [][]int, matrix_size int) [][]int {
	result := make([][]int, matrix_size)
	for i := range result {
		result[i] = make([]int, matrix_size)
		for j := 0; j < matrix_size; j++ {
			for k := 0; k < matrix_size; k++ {
				result[i][j] = (result[i][j] + a[i][k]*b[k][j]) % C.q
			}
		}
	}
	return result
}

// Function to compare a matrix with the identity matrix
func is_identity_matrix(mat [][]int, matrix_size int) bool {
	identity := identity_matrix(matrix_size)
	for i := 0; i < matrix_size; i++ {
		for j := 0; j < matrix_size; j++ {
			if mat[i][j] != identity[i][j] {
				return false
			}
		}
	}
	return true
}

func generate_invertible_matrixQ(matrix_size int) [][]int {
	var count int16
	for count = 0; count < 1000; count++ {
		matrix := generate_matrixQ(matrix_size)
		invMatrix := inverse_matrixQ(matrix, matrix_size)
		product := multiply_matricesQ(matrix, invMatrix, matrix_size)

		if is_identity_matrix(product, matrix_size) {
			return matrix
		}
	}
	return nil
}

func sample_random_vector(size, max int) []int {
	vector := make([]int, size)
	for i := range vector {
		vector[i] = rand.Intn(max)
	}
	return vector
}

func sample_uniform_matrix(rows, cols int) [][]int {
	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
		for j := range matrix[i] {
			matrix[i][j] = rand.Intn(C.q)
		}
	}
	return matrix
}

func sample_normal_matrix(rows, cols int, mean, sigma float64) [][]int {

	//fmt.Println("Hello", rows)

	matrix := make([][]int, rows)
	//fmt.Println("Hello", len(matrix))

	for i := range matrix {
		matrix[i] = make([]int, cols)
		//fmt.Println(len(matrix[0]))
		for j := range matrix[i] {
			matrix[i][j] = int(math.Round(rand.NormFloat64()*sigma + mean))
		}
	}

	//fmt.Println(matrix)
	return matrix
}

func VerifyR(matrix [][]int) {

	//fmt.Println(len(matrix))
	//fmt.Println(len(matrix[0]))
	//fmt.Println(matrix)

	rows := len(matrix)
	cols := len(matrix[0])
	cnz, sum, second := 0, 0, 0
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			val := matrix[i][j]
			if val != 0 {
				sum += val
				second += val * val
				cnz++
			}
		}
	}
	fmt.Printf("VR= %dx%d cnz=%d sum=%d 2nd=%d\n", rows, cols, cnz, sum, second)
}

func VerifyE(vector []int) {
	cols := len(vector)
	cnz, sum, second := 0, 0, 0
	for j := 0; j < cols; j++ {
		val := vector[j]
		if val != 0 {
			sum += val
			second += val * val
			cnz++
		}
	}
	fmt.Printf("VE= size=%d cnz=%d sum=%d 2nd=%d\n", cols, cnz, sum, second)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func lemma1Check(R [][]int, e []int) error {

	Rt := transpose_matrix(R)

	fmt.Println(len(Rt), "x", len(Rt[0]), len(e))

	identitySize := len(e) - len(Rt[0])

	RtExt, err1 := concatenate_matrices_row(Rt, identity_matrix(identitySize))

	if err1 != nil {
		return fmt.Errorf("lemma1Check - err1")
	}

	Rte, err2 := dot_product_MV(RtExt, e)

	if err2 != nil {
		return fmt.Errorf("lemma1Check - err2")
	}

	fmt.Println("Rte =", Rte)
	VerifyE(Rte)

	for _, r := range Rte {
		if abs(r) > C.q/4 {
			fmt.Println("Lemma1Check failed =", r)
			log.Fatal("Lemma1Check assertion failed")
			return fmt.Errorf("Lemma1Check - assertion failed")
		}
	}

	return nil
}

// pk, sk
func equalityCheck(A0, A1, R, H [][]int) error {
	HmG, err1 := dot_product_MM(H, C.G)

	if err1 != nil {
		return fmt.Errorf("EqualityCheck - err1")
	}

	helper1, err2 := add_matrix(A1, HmG)

	if err2 != nil {
		return fmt.Errorf("EqualityCheck - err2")
	}

	A_mu_01, err3 := concatenate_matrices_row(A0, helper1)

	if err3 != nil {
		return fmt.Errorf("EqualityCheck - err3")
	}

	R_padded, err4 := concatenate_matrices_col(R, identity_matrix(C.nk))

	if err4 != nil {
		return fmt.Errorf("EqualityCheck - err4")
	}

	left, err5 := dot_product_MM(A_mu_01, R_padded)

	if err5 != nil {
		return fmt.Errorf("EqualityCheck - err5")
	}

	if !reflect.DeepEqual(over_QM(left), over_QM(HmG)) {
		return fmt.Errorf("EqualityCheck - failed")
	}

	return nil
}

func encode(message []int) []int {
	p := math.Round(float64(C.q) / 2)
	encoded_message := make([]int, len(message))

	for i := range message {
		encoded_message[i] = int(float64(message[i])*p) % C.q
	}

	return encoded_message
}

func decode(encoded_message []int) ([]int, []int) {
	b1 := float64(C.q) / 4
	b2 := float64(C.q) * 3 / 4
	p := int(math.Round(float64(C.q) / 2))

	decoded_message := make([]int, len(encoded_message))
	e2_calc := make([]int, len(encoded_message))

	for i := range encoded_message {

		if encoded_message[i] < 0 {
			encoded_message[i] += C.q
		}

		if float64(encoded_message[i]) > b1 && float64(encoded_message[i]) < b2 {
			decoded_message[i] = 1
		}
		e2_calc[i] = encoded_message[i] - p*decoded_message[i]
	}

	return decoded_message, e2_calc
}

/*
func sample_gaussian_vector(size int) []int {
	vector := make([]int, size)
	for i := 0; i < size; i++ {
		vector[i] = int(math.Round(rand.NormFloat64()))
	}
	return vector
}

func sample_perturbation_v2(R1 [][]int, s int, m int) ([]int, error) {
	w := len(R1[0]) // R1 is (m_bar x w)

	// Step 1: Compute Sigma_p
	Sigma_G := identity_matrix(w)
	identityW := identity_matrix(w)

	// Block matrix RI = [R1; I_w]
	RI, err1 := concatenate_matrices_col(R1, identityW)

	if err1 != nil {
		return nil, fmt.Errorf("samplePerturbationV2 - err1, %v", err1)
	}

	fmt.Printf("len R1 %v\n", len(R1))
	fmt.Printf("len RI %v\n", len(RI))
	fmt.Printf("len M%v\n", C.m)

	// Compute Sigma_p = I_m - RI * (s * Sigma_G) * RI^T

	Sigma_G_s := make([][]int, len(Sigma_G))

	for i := range Sigma_G_s {
		Sigma_G_s[i] = make([]int, len(Sigma_G[0]))
		for j := range Sigma_G_s[0] {
			Sigma_G_s[i][j] = Sigma_G[i][j] * s
		}
	}

	helper1, err2 := dot_product_MM(RI, Sigma_G_s)

	if err2 != nil {
		return nil, fmt.Errorf("samplePerturbationV2 - err2, %v", err2)
	}

	//fmt.Printf("aaaaaaaa %v", helper1)

	helper1, err3 := dot_product_MM(helper1, transpose_matrix(RI))

	if err3 != nil {
		return nil, fmt.Errorf("samplePerturbationV2 - err3, %v", err3)
	}

	Sigma_p, err4 := add_matrix(identity_matrix(m), negate_matrix_values(helper1))

	if err4 != nil {
		return nil, fmt.Errorf("samplePerturbationV2 - err4, %v", err4)
	}

	for i := 0; i < m; i++ {
		for j := 0; j < m; j++ {
			Sigma_p[i][j] -= s * RI[i][j] // Simplified element-wise computation
		}
	}

	// Step 2: Sample fresh perturbation p from integer-valued Gaussian
	p := sample_gaussian_vector(m)
	return p, nil
}
*/

func over_QV(v []int) []int {
	for i := range v {
		if v[i] < 0 {
			v[i] += C.q
		}
	}

	return v
}

func over_QM(M [][]int) [][]int {
	rows := len(M)
	cols := len(M[0])
	for i := range rows {
		for j := range cols {
			if M[i][j] < 0 {
				M[i][j] += C.q
			}
		}
	}

	return M
}

// GaussianSample generates a sample from a discrete Gaussian distribution with mean 0 and given standard deviation
func GaussianSample(stddev float64, psize int) []float64 {
	rand.Seed(time.Now().UnixNano())
	samples := make([]float64, psize)
	for i := 0; i < psize; i++ {
		samples[i] = rand.NormFloat64() * stddev
	}
	return samples
}

// GenerateP generates p based on the given equation
func GenerateP(Sigma, R, SigmaG *mat.Dense, psize int) *mat.VecDense {
	// Compute Rt * R
	var RtR mat.Dense
	Rt := mat.DenseCopyOf(R.T())
	RtR.Mul(Rt, R)

	// Compute Sigma_p = Sigma - R * Sigma_G * Rt
	var RSigmaG mat.Dense
	RSigmaG.Mul(R, SigmaG)
	var RSigmaGRt mat.Dense
	RSigmaGRt.Mul(&RSigmaG, Rt)

	var SigmaP mat.Dense
	/*
		fmt.Println(SigmaG.Dims())
		fmt.Println(R.Dims())
		fmt.Println(Rt.Dims())
		fmt.Println(Sigma.Dims())
		fmt.Println(RSigmaGRt.Dims())
	*/
	SigmaP.Sub(Sigma, &RSigmaGRt)

	// Ensure Sigma_p is at least 2 * (R * Rt)
	var SigmaPLowerBound mat.Dense
	SigmaPLowerBound.Scale(2, &RtR)

	// Element-wise max operation: Sigma_p = max(Sigma_p, 2 * R * Rt)
	r, c := SigmaP.Dims()
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			SigmaP.Set(i, j, math.Max(SigmaP.At(i, j), SigmaPLowerBound.At(i, j)))
		}
	}

	// Compute r / sqrt(Sigma_p) for Gaussian sampling
	var sqrtSigmaP mat.Dense = *mat.NewDense(r, c, nil)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			sqrtSigmaP.Set(i, j, math.Sqrt(SigmaP.At(i, j)))
		}
	}

	// Sample from Gaussian with standard deviation r / sqrt(Sigma_p)
	stddev := 1.0 / sqrtSigmaP.At(0, 0) // Assuming isotropic case
	pSamples := GaussianSample(stddev, psize)

	return mat.NewVecDense(psize, pSamples)
}

func sample_perturbation_v2(R [][]int, s float64, psize int) ([]int, error) {

	rows := len(R)
	cols := len(R[0])

	Sigma := mat.NewDense(rows, rows, nil)
	for i := 0; i < rows; i++ {
		for j := 0; j < rows; j++ {
			Sigma.Set(i, j, rand.Float64()*s)
		}
	}

	Sigma_G := mat.NewDense(cols, cols, nil)
	for i := 0; i < cols; i++ {
		Sigma_G.Set(i, i, 1.0) // Identity matrix
	}

	// Construct RI = [R1; I_w]
	RI := mat.NewDense(rows, cols, nil)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if i < rows {
				RI.Set(i, j, float64(R[i][j])) // Copy R1
			} else if i-cols == j {
				RI.Set(i, j, 1.0) // Identity part
			}
		}
	}

	//fmt.Printf("Sigma=%v\n", Sigma)
	//fmt.Printf("Sigma_G=%v\n", Sigma_G)
	//fmt.Printf("RI=%v\n", RI)
	pmat := GenerateP(Sigma, RI, Sigma_G, psize)

	//fmt.Printf("pmat=%v\n", pmat)

	p := make([]int, psize)

	for i := range psize {
		p[i] = int(pmat.AtVec(i) * 10)
	}
	//fmt.Printf("p=%v\n", p)
	return p, nil
}
