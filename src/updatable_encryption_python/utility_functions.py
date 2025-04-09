import numpy as np

from updatable_encryption_python.parameters import *

def sample_uniform_matrix(rows, cols, q):
    return np.random.randint(0, q, (rows, cols))

def sample_normal_matrix(rows, cols, mean, sigma):
    # Sample from a continuous Gaussian distribution
    samples = np.random.normal(loc=mean, scale=sigma, size=(rows, cols))
    # Round to the nearest integer to make it discrete
    discrete_samples = np.round(samples).astype(int)
    return discrete_samples
# returns a value in range -q/2 to q/2
def modQ(input):
    input = input % q
    if input > q // 2:
        return input - q
    return input

#print
def VerifyR(R):
    cnz = 0
    sum = 0
    second = 0
    rows = len(R)
    cols = len(R[0])
    for i in range(rows):
        for j in range(cols):
            if R[i][j] != 0:
                sum = sum + R[i][j]
                second = second + (R[i][j]*R[i][j])
                cnz = cnz + 1

    print("VR=", rows, "x", cols, "cnz=", cnz, "sum=", sum, "2nd=", second)

#print
def VerifyE(E):
    cnz = 0
    sum = 0
    second = 0
    cols = len(E)
    for j in range(cols):
        if E[j] != 0:
            sum = sum + E[j]
            second = second + (E[j]*E[j])
            cnz = cnz + 1

    print("VE=", "size=",cols, "cnz=", cnz, "sum=", sum, "2nd=", second)

# lemma check 1
def Lemma1Check(R, e):
    Rt = R.T
    #m,n = Rt.Shape
    print(len(Rt), len(Rt[0]), len(e))

    Rt_ext = np.block([ Rt, np.eye(len(e)-len(Rt[0]), dtype=int) ])
    Rte = np.dot(Rt_ext, e)
    #print("Rte=",Rte)
    #VerifyE(Rte)

    for r in Rte:
        if abs(r) > q/4:
            print("Lemma1Check failed =", r)
            assert 1 == 2

def EqualityCheck(pk, sk):

    A0, A1, A2 = pk

    # Combined matrix
    HmG = np.dot(H_hardcoded, G) % q

    A_mu_01 = np.block([A0, A1 + HmG])  

    R_padded = np.block([
        [sk],  
        [np.eye(nk, dtype=int)]
    ])
    left = np.dot(A_mu_01, R_padded) % q

    #print("left: ", np.dot(A_mu_01, R_padded) % q)
    #print("right: ", HmG % q)
    assert np.array_equal(left, HmG) == True
   
    H_hardcoded_inv_t = np.linalg.inv(H_hardcoded).astype(int).T
    H_I = np.dot(H_hardcoded, H_hardcoded_inv_t.T)
    I = np.block([
        [np.eye(n, dtype=int)]
    ])
    assert np.array_equal(H_I, I) == True

def decode(encoded_message):

    #print("encoded_message=",encoded_message)
    b1 = q / 4
    b2 = q * 3 / 4
    p = round(q/2)
    decoded_message = np.zeros(len(encoded_message), dtype=int)
    e2_calc = np.zeros(len(encoded_message), dtype=int)
    for i in range(len(decoded_message)):
        if encoded_message[i] >= b1 and encoded_message[i] <= b2:
            decoded_message[i] = 1
        e2_calc[i] = encoded_message[i] - p * decoded_message[i]

    #print("decoded_message=",decoded_message)
    return decoded_message, e2_calc


def encode(message):
    p = round(q/2)
    encoded = message * p
    encoded = encoded.astype(int) % q
    return encoded

def sample_perturbation_v2(R1, s, m):
    w = R1.shape[1]          # R is (m_bar x w)
    #n = H.shape[0]          # H is (n x n)

    # Step 1: Compute Sigma_p
    Sigma_G = np.eye(w)  # Assume identity for simplicity
    RI = np.block([[R1], [np.eye(w)]])
    #print(RI.shape)
    Sigma_p = np.eye(m) - RI @ (s * Sigma_G) @ RI.T
    #print(Sigma_G.shape, Sigma_p.shape)

    # Sample fresh perturbation p
    p = np.random.multivariate_normal(np.zeros(m), Sigma_p)  # Sampling from Gaussian
    p = p.astype(int)
    return p 

# Function to generate a random n x n matrix with determinant 1 and not an identity matrix
def generate_matrixQ(matrix_size):
    matrix = identity_matrix(matrix_size)

    for i in range(matrix_size*3):
        r1 = np.random.randint(0, matrix_size)
        r2 = np.random.randint(0, matrix_size)
        if r1 != r2 :
            factor = np.random.randint(0, q-1) + 1
            for j in range(matrix_size):
                matrix[r1][j] = modQ(matrix[r1][j] + factor*matrix[r2][j])
    return matrix

# Function to generate an identity matrix of size n
def identity_matrix(matrix_size):
    return np.eye(matrix_size)

# Function to compute the modular inverse of a number mod q
def mod_inverseQ(a):
    a = a % q
    for x in range(q):
        if a*x % q == 1:
            return x 

    #Assumes q is prime
    return 1
'''
def inverse_matrixQ(mat, matrix_size):
    inv = identity_matrix(matrix_size)
    augmented = [row + inv_row for row, inv_row in zip(mat, inv)]

    for i in range(matrix_size):
        pivot = augmented[i][i]
        invPivot = mod_inverseQ(pivot)
        
        for j in range(2 * matrix_size):
            augmented[i][j] = (augmented[i][j] * invPivot) % q
        
        for k in range(matrix_size):
            if k != i:
                factor = augmented[k][i]
                for j in range(2 * matrix_size):
                    augmented[k][j] = (augmented[k][j] - factor * augmented[i][j]) % q
                    if augmented[k][j] < 0:
                        augmented[k][j] += q
    
    return [row[matrix_size:] for row in augmented]
'''

def inverse_matrixQ(mat, matrix_size):
    inv = identity_matrix(matrix_size)
    augmented = []
    for i in range (len(mat)):
        new_array = []
        for j in range(len(mat[i])):
            new_array.append(mat[i][j])

        for j in range(len(inv[i])):
            new_array.append(inv[i][j])     
        augmented.append(new_array)

    for i in range(matrix_size):
        pivot = augmented[i][i]  
        invPivot = mod_inverseQ(pivot)
        for j in range(2*matrix_size):
            augmented[i][j] = modQ(augmented[i][j] * invPivot)
        for k in range(matrix_size):
            if k!=i:
                factor = augmented[k][i]
                for j in range(2*matrix_size):
                    augmented[k][j] = modQ(augmented[k][j] - factor*augmented[i][j])
                    if augmented[k][j] < 0 :
                        augmented[k][j] += q
					
    for i in range(matrix_size):
        inv[i] = augmented[i][matrix_size:]
    return inv

'''
 Function to compute the inverse of a matrix modulo q
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
'''

# Function to compare a matrix with the identity matrix
def is_identity_matrix(mat, matrix_size):
    identity = identity_matrix(matrix_size)
    for i in range(matrix_size):
        for j in range(matrix_size):
            if mat[i][j] != identity[i][j]:
                return False
    return True

# Function to generate an invertible matrix mod q
def generate_invertible_matrixQ(matrix_size):
    for _ in range(1000):
        matrix = generate_matrixQ(matrix_size)
        inv_matrix = inverse_matrixQ(matrix, matrix_size)
        product = np.dot(matrix, inv_matrix) % q
        if is_identity_matrix(product, matrix_size):
            return matrix
    return None