#outside libraries
import numpy as np

#local files
from dilithium_constants import * 

# Bit-reversal permutation of input (a bit number is re-writen from right to left - example -> 001 -> 100 or 010 -> 010)
# Example : if we input number 1 and use 3 bits => 1 is 001 (in bits) , which transforms into 100 , which is equal to 4
def bit_reverse(input_number, bits):
    output_number = 0
    for bit_position in range(bits):
        if input_number & (1 << bit_position):
            output_number |= 1 << (bits - 1 - bit_position)
    return output_number

#Number Theoretic Transform (NTT)
# poly - polynomial to be transformed
# poly_length - length of the polynomial, must be a power of 2
# prime_modulus - prime modulus
# primitive - primitive n-th root of unity modulo prime_modulus
def ntt(poly, poly_length, prime_modulus, primitive):
    #The input polynomial must have exactly n coefficients
    assert len(poly) == poly_length , "Polynomial does not match poly_length - ntt"

    # Ensure Python integers
    poly = [int(x) for x in poly]  
    
    # Precompute powers of the root of unity g
    root_powers = [1]
    for _ in range(1, poly_length):
        root_powers.append((root_powers[-1] * primitive) % prime_modulus)
        
    bits = poly_length.bit_length() - 1
    poly = [poly[bit_reverse(Counter_1, bits)] for Counter_1 in range(poly_length)]

    # Iterative Cooley-Tukey NTT
    length = 2
    while length <= poly_length:
        half = length // 2
        root = root_powers[poly_length // length]
        for Counter_1 in range(0, poly_length, length):
            w = 1
            for Counter_2 in range(half):
                u = poly[Counter_1 + Counter_2]
                v = (poly[Counter_1 + Counter_2 + half] * w) % prime_modulus
                poly[Counter_1 + Counter_2] = (u + v) % prime_modulus
                poly[Counter_1 + Counter_2 + half] = (u - v) % prime_modulus
                w = (w * root) % prime_modulus
        length *= 2

    return poly


# Inverse Number Theoretic Transform (INTT)
# poly - polynomial to be reversed transformed
# poly_length - length of the polynomial, must be a power of 2
# prime_modulus - prime modulus
# primitive - primitive n-th root of unity modulo prime_modulus
def intt(poly, poly_length, prime_modulus, primitive):
    
    #The input polynomial must have exactly n coefficients
    assert len(poly) == poly_length , 'Polynomial does not match poly_length - intt'

    # Compute modular inverse of n
    poly_inv = pow(poly_length, prime_modulus - 2, prime_modulus)

    # Compute inverse of the primitive root
    primitive_inv = pow(primitive, prime_modulus - 2, prime_modulus)

    # Apply NTT using primitive_inv
    poly = ntt(poly, poly_length, prime_modulus, primitive_inv)
  
    # Normalize by multiplying by poly_inv
    poly = [(x * poly_inv) % prime_modulus for x in poly]

    return poly

# compute ntt for vector (a x 1) where each element is polynomial
def ntt_vector(vector, poly_length, prime_modulus, primitive):
    return np.array([ntt(poly, poly_length, prime_modulus, sample_primitive) for poly in vector])
    
# compute ntt for matrix (a x b) where each element is polynomial
def ntt_matrix(matrix, poly_length, prime_modulus, primitive):
    return np.array([[ntt(poly, n, q, sample_primitive) for poly in vector] for vector in matrix])


# compute intt for vector (a x 1) where each element is polynomial
def intt_vector(vector, poly_length, prime_modulus, primitive):
    return np.array([intt(poly, poly_length, prime_modulus, sample_primitive) for poly in vector])
    
# compute intt for matrix (a x b) where each element is polynomial
def intt_matrix(matrix, poly_length, prime_modulus, primitive):
    return np.array([[intt(poly, n, q, sample_primitive) for poly in vector] for vector in matrix])


