#outside libraries
import numpy as np

#local files
from dilithium_constants import *  
from dilithium_multi_usage_functions import shake256_hash, expand_matrix_A, sample_in_ball, mod_multiply_poly, use_hint_q, weight, norm_inf_array
from ntt import ntt, intt, ntt_matrix, ntt_vector
from pack_functions import pack_pk, unpack_pk, unpack_sig, bit_pack_w_vector

def verify(pk, M, sigma):

    # Unpack public key and signature
    c_tilde, z, h = unpack_sig(sigma, gamma_1)
    rho, t1 = unpack_pk(pk)

    # Step 27: Expand A from ρ and store in NTT representation
    # Transform A to NTT domain
    #A = expand_matrix_A(rho, k, l, n, q)
    #A_ntt = ntt_matrix(A, n, q, sample_primitive)
    A_ntt = ntt_matrix(expand_matrix_A(rho, k, l, n, q), n, q, sample_primitive)

    # Step 28: Compute µ = CRH(CRH(ρ || t1) || M)
    tr = shake256_hash(pk, 48)  # 384 bits = 48 bytes
    # 384 bits = 48 bytes
    #tr = shake256_hash(tr_input, 48)  
    # Outer CRH (384 bits)
    mu = shake256_hash((tr + M), 48)   

    # Step 29: Compute c from SampleInBall
    c = sample_in_ball(c_tilde, n, tau)

    # Step 30: Compute w'1
    # Transform z into NTT
    z_ntt = ntt_vector(z, n, q, sample_primitive)

    # Transform c into NTT
    c_ntt = ntt(c, n, q, sample_primitive)

    # Transform t1·2^d into NTT
    t1_scaled = [[(coeff * (1 << d)) % q for coeff in poly] for poly in t1]
    t1_ntt = [ntt(poly, n, q, sample_primitive) for poly in t1_scaled]

    # Compute Az in NTT
    Az_ntt = []
    for counter_1 in range(k):
        row_sum = np.zeros(n, dtype=int)
        for counter_2 in range(l):
            product = mod_multiply_poly(A_ntt[counter_1][counter_2] , z_ntt[counter_2], q)
            row_sum = (row_sum + product) % q     
        Az_ntt.append((row_sum+q) % q)

    # Compute w = Az - c·t1 in NTT and transform back
    w_ntt = []
    for counter_1 in range(k):
        row_sum = np.zeros(n, dtype=int)
        row_sum = row_sum + Az_ntt[counter_1] - mod_multiply_poly(c_ntt, t1_ntt[counter_1], q)
        w_ntt.append((row_sum+q) % q)
    w = [intt(poly, n, q, sample_primitive) for poly in w_ntt]    

    # UseHint to compute w'1
    w1 = use_hint_q(h, w, 2 * gamma_2, n)

    # Step 31: Verify conditions
    # (1) ∥z∥∞ < γ1 - β
    norm_z = norm_inf_array(z, (gamma_1 - beta))
    if norm_z:
        print("Norm Z check failed")
        return False

    # (2) Verify that c_tilde = H(µ || w'1)
    w1_bytes = bit_pack_w_vector(w1)
    c_tilde_check = shake256_hash((mu + w1_bytes), 32)  

    if c_tilde_check != c_tilde:
        print("c_tilde check failed")
        return False

    # (3) Verify that the number of 1's in h ≤ ω
    if weight(h) > omega:
        print("Weight check failed")
        return False

    # All checks passed msg is verified
    return True
