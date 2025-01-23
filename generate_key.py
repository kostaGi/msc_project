#outside libraries
import numpy as np

#local files
from dilithium_constants import * 
from ntt import ntt, intt , ntt_matrix, ntt_vector
from dilithium_multi_usage_functions import shake256_hash, expand_matrix_A, mod_multiply_poly, reduce_mod
from pack_functions import pack_pk, pack_sk, unpack_pk, unpack_sk

# Generates a random 256-bit string
def random_256_bit_string():
    return np.random.bytes(32)  # 32 bytes = 256 bits

# Generates a fixed number of polynomials (total_number_of_poly), 
# each with specific length/degree (poly_length) 
# with cofficients bounded by some value (bounded).
# The 'seed' can be used to make the method deterministic
def sample_polynomials(seed, total_number_of_poly, bound, poly_length, modulus):
    # Set RNG seed deterministically
    np.random.seed(int.from_bytes(seed, 'big') % (2**32))  
    polynomials = []
    for _ in range(total_number_of_poly):
        coeffs = list(np.random.randint(-bound, bound + 1, size=poly_length))
        polynomials.append(coeffs) #% modulus)
    return tuple(polynomials)    

def power2_round(poly, d):
    factor = 1 << d
    high_array = []
    low_array = []
    for element in poly:
        tmp = element % q
        low = reduce_mod(tmp, factor)
        high_array.append((tmp - low) >> d)
        low_array.append(low)
    return high_array, low_array

# Generates public and private key
def generate_keys():

    # Step 01: Generate random seed zeta
    zeta = random_256_bit_string()

    # Step 02: Compute (ρ, ς, K) = H(zeta) using SHAKE-256
    # 256/8 * 3 = 96 bytes
    hashed_output = shake256_hash(zeta, 96)  
    rho, sigma, K = hashed_output[:32], hashed_output[32:64], hashed_output[64:]

    # Step 03: Generate secret key vectors s1 and s2
    s1 = sample_polynomials(sigma, l, eta, n, q)
    # Reuse sigma for simplicity (backwards)
    s2 = sample_polynomials(sigma[::-1], k, eta, n, q)  

    # Step 04: Expand matrix A from rho 
    # A = (k by l)
    A = expand_matrix_A(rho, k, l, n, q)

    # Step 05: Compute t = As1 + s2      
    # Transform A to NTT domain
    # !!! sample_primitive - generated in the dilithium_constansts.py 
    A_ntt = ntt_matrix(A, n, q, sample_primitive)
    
    # Transform s1 to NTT domain
    s1_ntt = ntt_vector(s1, n, q, sample_primitive)

    # As1 = product = (NTT^-1)(NTT(A) * NTT(s1))
    # t = As1 + s2 
    t = []
    for counter_1 in range(k):
        row_sum = np.zeros(n, dtype=int)
        for counter_2 in range(l):
            product = intt(mod_multiply_poly(A_ntt[counter_1][counter_2] , s1_ntt[counter_2], q), n, q, sample_primitive)
            row_sum = (row_sum + product) % q
        t.append((row_sum + s2[counter_1]) % q)


    # Step 06: Split t into (t1, t0) using Power2Round
    t1, t0 = zip(*[power2_round(poly, d) for poly in t])

    # Step 07-08: Compute tr = CRH(ρ || t1) ; Output public and secret keys
    # Pack up the bytes
    pk = pack_pk(rho, t1)
    
    # 384 bits = 48 bytes
    tr = shake256_hash(pk, 48) 
    sk = pack_sk(rho, K, tr, s1, s2, t0)

    # testing (un)pack
    rho1, t11 = unpack_pk(pk)
    #[print(type(c)) for c in t1]
    #print("t1: ", type(t1) , t1)
    #print("t11: ", type(t11), t11)
    assert rho == rho1, "rho pk FAILS"
    assert t1 == t11, "t1 pk FAILS"

    (rho2, K2, tr2, s12, s22, t02) = unpack_sk(sk)
    assert rho == rho2, "rho sk FAILS"
    assert K == K2, "K sk FAILS"
    assert tr == tr2, "tr sk FAILS"
    assert (s1 == s12), "s1 sk FAILS"
    assert (s2 == s22), "s2 sk FAILS"
    assert (t0 == t02), "t0 sk FAILS"

    return pk, sk