#The following functions are used in at least 2 of the 3 parts of the Dilithium algorithm - generate_key, sign_message, verify_signature

#outside libraries
import numpy as np
from hashlib import shake_256, shake_128

#local files
from dilithium_constants import * 


# multiply 2 numbers modulo a number (without overflow)
def mod_multiply(left_int, right_int, modulo):
    result = 0
    left_int = left_int % modulo
    while right_int > 0:
        # If right_int is odd, add left_int to the result
        if right_int % 2 == 1:  
            result = (result + left_int) % modulo

        # Double left_int and halve right_int    
        left_int = (left_int * 2) % modulo  
        right_int //= 2
    return result

# multiply 2 polynomials modulo a number (without overflow) (ntt format)
def mod_multiply_poly(poly_1, poly_2, modulo):
    assert len(poly_1) == len(poly_2) , "Invalid length - mod_multiply_poly"
    result = np.zeros(len(poly_1), dtype=int)
    for counter_1 in range(len(poly_1)):
        result[counter_1] = mod_multiply(poly_1[counter_1], poly_2[counter_1], modulo)
    return result

# SHAKE-128 hash function with variable-length output
def shake128_hash(data, output_length):
    return shake_128(data).digest(output_length)  

# SHAKE-256 hash function with variable-length output
def shake256_hash(data, output_length):
    return shake_256(data).digest(output_length)   

# Expand matrix A deterministically from seed,
# matrix has size k by l,
# each element is a polynomial with length poly_length
def expand_matrix_A(seed, k, l, poly_length, modulus): 
    matrix = []
    for counter_1 in range(k):
        row = []
        for counter_2 in range(l):
            random_bytes = shake128_hash(seed + (counter_1*256+counter_2).to_bytes(2, 'little'), 3*260)
            current_hash = int.from_bytes(random_bytes, byteorder='big')
            poly = []
            for _ in range(n):
                poly.append(current_hash & ((1 << 23) - 1))
                current_hash = current_hash >> 24
            row.append(poly)
        matrix.append(row)
    return matrix   


# ExpandMask generates a polynomial with coefficients in [-gamma_1, gamma_1) 
# rho_prime: Seed for pseudorandom generation (bytes)
# kappa: Counter to ensure unique output (integer)
# n: Number of coefficients in the polynomial
# gamma_1: Bound on the coefficients
def expand_mask(rho_prime, kappa, l, n, gamma_1):
    gamma_power = 17

    if gamma_1 == (1 << 17):
        gamma_power = 17
    else:
        gamma_power = 19

    # (256 * 18) // 8 = 576 
    # (256 * 20) // 8 = 640 
    bytes_to_generate = (n * (gamma_power+1) // 8)
    polynomials = []
    for counter_1 in range(l):
        random_bytes = shake_256(rho_prime + (counter_1+kappa).to_bytes(2, 'little')).digest(bytes_to_generate)
        random_bytes = int.from_bytes(random_bytes, byteorder='big')
        poly = []
        AND_mask = (1 << (gamma_power + 1)) - 1
        for counter_2 in range(n):
            shift_amount = (n - 1 - counter_2) * 18
            coeff = (gamma_1 -1) - ((random_bytes >> shift_amount) & AND_mask)
            poly.append(coeff)   
        polynomials.append(poly)
    return polynomials

# Maps a binary hash to a sparse polynomial with tau nonzero coefficients in {+1, -1}
# c_tilde: Input hash as bytes.
# n: Degree of the polynomial.
# tau: Number of nonzero coefficients.
def sample_in_ball(c_tilde, n, tau):
    
    # initilize poly with all 0
    poly = [0] * n

    # Absorb ˜c into SHAKE-256 to create a random intitial stream of 32 bytes
    shake = shake_256()
    shake.update(c_tilde)
    random_bytes = shake.digest(32)  

    # Extract τ sign bits from the first 8 bytes
    sign_bits = [(random_bytes[i // 8] >> (i % 8)) & 1 for i in range(tau)]

    # Start reading after the first 8 bytes
    byte_index = 8  
    for i in range(n-tau, n):
        current_sign_pos = 0
        while True:
            if byte_index >= len(random_bytes):
                # Extend the random stream if needed
                random_bytes += shake.digest(64)
            j = random_bytes[byte_index]
            byte_index += 1
            # Rejection sampling condition
            if j <= i:  
                poly[i] = poly[j]
                if sign_bits[current_sign_pos] == 1:
                    poly[j] = 1
                else:
                    poly[j] = -1
                current_sign_pos += 1
                break
    return poly
    

# Computes a modulo operation with an optimization to pass the mandatory checks when signing a message 
# After modulo if (element) is bigger than : (modolu // 2), return : (element - modolu), else return: (element)
def reduce_mod(element,modolu):
    element = element % modolu
    if element > (modolu >> 1):
        element -= modolu
    return element

# Decompose an element into high (element // alpha) and low (element % alpha)
def decompose_element(element, alpha):
    element = element % q
    low = reduce_mod(element, alpha)
    high = 0
    if element - low == q - 1:
        #high = 0
        low -= 1
    else:
        high = (element - low) // alpha
    return high, low

# Apply decompose to an n x 1 array
def decompose_poly(poly, alpha):
    high_array = []
    low_array = []
    for current_index in range(len(poly)):
        high, low = decompose_element(poly[current_index], alpha)
        high_array.append(high)
        low_array.append(low)
    return high_array, low_array

# Apply decompose to an l x n array
def decompose_vector(vector, alpha):
    high_array = []
    low_array = []
    for current_index in range(len(vector)):
        high, low = decompose_poly(vector[current_index], alpha)
        high_array.append(high)
        low_array.append(low)
    return high_array, low_array

# Apply decompose to an k x l x n array
def decompose_matrix(matrix, alpha):
    high_array = []
    low_array = []
    for current_index in range(len(matrix)):
        high, low = decompose_vector(matrix[current_index], alpha)
        high_array.append(high)
        low_array.append(low)
    return high_array, low_array


# Returns high bits of poly
def high_bits_poly_q(poly, alpha):
    high, _ = decompose_poly(poly, alpha)
    return high

# Returns low bits of poly
def low_bits_poly_q(poly, alpha):
    _, low = decompose_poly(poly, alpha)
    return low

# Returns high bits of vector
def high_bits_vector_q(vector, alpha):
    high, _ = decompose_vector(vector, alpha)
    return high

# Returns low bits of vector
def low_bits_vector_q(vector, alpha):
    _, low = decompose_vector(vector, alpha)
    return low

# Returns high bits of matrix
def high_bits_matrix_q(matrix, alpha):
    high, _ = decompose_matrix(matrix, alpha)
    return high

# Returns low bits of matrix
def low_bits_matrix_q(matrix, alpha):
    _, low = decompose_matrix(matrix, alpha)
    return low    

'''
Generic methods from documentation
# Returns high bits of poly
def high_bits_q(poly, alpha):
    high, _ = decompose_poly(poly, alpha)
    return high

# Returns low bits of poly
def low_bits_q(poly, alpha):
    _, low = decompose_poly(poly, alpha)
    return low
'''
# Make a hint for a vector
def make_hint_q(zeta, vector, alpha):
    assert len(zeta) == len(vector) , "Vectors does not match length - make_hint_q"
    high_vector = high_bits_vector_q(vector, alpha)
    high_vector_check = high_bits_vector_q( (vector+zeta), alpha)
    return_array = []
    for counter_1 in range(k):
        row_value = [0] * n
        for counter_2 in range(n):
            if high_vector[counter_1][counter_2] != high_vector_check[counter_1][counter_2]:
                row_value[counter_2] = 1   
        return_array.append(row_value)
    return tuple(return_array)

# Use a hint for a vector     
def use_hint_q(h, vector, alpha):

    assert len(h) == len(vector) , "Vectors does not match length - use_hint_q"
    m = (q - 1) // alpha
    high_vector, low_vector = decompose_vector(vector, alpha)
    return_vector = []
    for counter_1 in range(len(h)):
        current_row = [0] * n
        for counter_2 in range(n):
            if h[counter_1][counter_2] == 1:
                if low_vector[counter_1][counter_2] > 0:
                    current_row[counter_2] = (high_vector[counter_1][counter_2] + 1) % m
                else:
                    current_row[counter_2] = (high_vector[counter_1][counter_2] - 1) % m
            else:
                current_row[counter_2] = high_vector[counter_1][counter_2]
        return_vector.append(current_row)
    return np.array(return_vector)   

# check if matrix has true value
def norm_inf_array(poly_array, max_value):
    return any(any(check_norm_bound(coeff, max_value, q) for coeff in poly) for poly in poly_array)    

# check whether element is in boundary 
def check_norm_bound(element, max_value, modulo):
    """
    Norm bound is checked in the following four steps:
    x ∈ {0,        ...,                    ...,     q-1}
    x ∈ {-(q-1)/2, ...,       -1,       0, ..., (q-1)/2}
    x ∈ { (q-3)/2, ...,        0,       0, ..., (q-1)/2}
    x ∈ {0, 1,     ...,  (q-1)/2, (q-1)/2, ...,       1}
    """
    x = element % modulo
    x = ((modulo - 1) >> 1) - x
    x = x ^ (x >> 31)
    x = ((modulo - 1) >> 1) - x
    return x >= max_value


# compute the Hamming weight of a polynomial (number of nonzero coefficients).
def weight(h):
    #return sum([[element for element in poly] for poly in h]
    sum = 0
    for counter_1 in range(k):
        for counter_2 in range(n):    
            sum += h[counter_1][counter_2]
    return sum        


     





    