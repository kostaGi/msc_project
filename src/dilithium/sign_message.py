#outside libraries
import numpy as np


#local files
from dilithium.dilithium_constants import *  
from dilithium.dilithium_multi_usage_functions import shake256_hash, expand_matrix_A, expand_mask, sample_in_ball, low_bits_vector_q, norm_inf_array, mod_multiply_poly, make_hint_q, weight, high_bits_vector_q
from dilithium.ntt import ntt, intt, ntt_matrix, ntt_vector, intt_vector
from dilithium.pack_functions import unpack_sk, pack_sig, unpack_sig, bit_pack_w_vector



# Inputs: Message M, public key (tr, K, s1, s2, t0, rho)
# Can substitue rho with matrix A and skip step 9
def dilithium_sign(M, sk):

    #unpack secret key from bytes
    rho, K, tr, s1, s2, t0 = unpack_sk(sk)

    # Step 09: Expand matrix A from rho  (can be substituted with input the method)
    # Transform A to NTT domain
    #A = expand_matrix_A(rho, k, l, n, q)
    #A_ntt = ntt_matrix(A, n, q, sample_primitive)
    A_ntt = ntt_matrix(expand_matrix_A(rho, k, l, n, q), n, q, sample_primitive)
    
    # Step 10: Compute µ from tr and M
    # CRH(tr || M), 384-bit output
    mu = shake256_hash((tr + M), 48)

    # Step 11: Initialize variables
    kappa = 0
    z, h = None, None

    # Step 12: Compute ρ′ (not randomized)
    rho_prime = shake256_hash((K + mu), 48)
    
    # Step 13.1: Precompute NTT forms
    s1_ntt = ntt_vector(s1, n, q, sample_primitive)
    s2_ntt = ntt_vector(s2, n, q, sample_primitive)
    t0_ntt = ntt_vector(t0, n, q, sample_primitive)
    
    # Step 13.2 (repeat until parameters satisfy the conditions in step 21)
    while z is None and h is None:
        # Step 14: Generate y from ExpandMask
        # Expand y deterministically
        y = expand_mask(rho_prime, kappa, l, n, gamma_1)  
        y_ntt = ntt_vector(y, n, q, sample_primitive)

        # Step 15: Compute w = A * y
        w = []
        for counter_1 in range(k):
            row_sum = np.zeros(n, dtype=int)
            for counter_2 in range(l):
                product = intt(mod_multiply_poly(A_ntt[counter_1][counter_2] , y_ntt[counter_2], q), n, q, sample_primitive)
                row_sum = (row_sum + product) % q  
            w.append(row_sum)

        # Step 16: Compute high bits w1
        w1 = high_bits_vector_q(w, 2 * gamma_2)

        # Step 17: Hash µ || w1 to generate challenge
        # Pack w1 into bytes
        w1_bytes = bit_pack_w_vector(w1)
        # Concatenate µ and w1
        data = mu + w1_bytes

        # 256-bit output
        c_tilde = shake256_hash(data, 32)  
        
        # Step 18: Calculate sample_in_ball with c_tilde
        c = sample_in_ball(c_tilde, n, tau)  

        # calculate the ntt of c 
        c_ntt = ntt(c, n, q, sample_primitive)

        # Step 19: Compute z = y + c * s1
        # z = y + (NTT^-1)(NTT(c) * NTT(s1)) 
        #ntt_c_s1 = (c_ntt * s1_ntt) % q
        ntt_c_s1 = []
        for counter_1 in range(l):
                product = mod_multiply_poly(c_ntt, s1_ntt[counter_1], q)
                ntt_c_s1.append((product+q) % q)
        cs1 = intt_vector(ntt_c_s1, n, q, sample_primitive) % q
        z = (y + cs1) % q

        # Step 20: Compute r0 = LowBits_q(w - c * s2, 2γ2)
        ntt_c_s2 = []
        for counter_1 in range(k):
                product = mod_multiply_poly(c_ntt, s2_ntt[counter_1], q)
                ntt_c_s2.append((product+q) % q)        
        cs2 = intt_vector(ntt_c_s2, n, q, sample_primitive) % q
        w_cs2 = (w - cs2)
        r0 = low_bits_vector_q(w_cs2, 2 * gamma_2)

        # Step 21: Check bounds for z and r0
        # Compute infinity norms
        norm_z = norm_inf_array(z, (gamma_1 - beta))
        norm_r0 = norm_inf_array(r0, (gamma_2 - beta))

        # Check bounds
        if norm_z or norm_r0:
            z, h = None, None  # Reset and retry
            kappa += 1
            continue
        
        # Step 22-23: Compute hint h
        # Compute ct0 = (NTT^-1)(NTT(c) * NTT(t0))
        t0_ntt = np.array([ntt(poly, n, q, sample_primitive) for poly in t0])
        ntt_c_t0 = []
        for counter_1 in range(k):
            product = mod_multiply_poly(c_ntt, t0_ntt[counter_1], q)
            ntt_c_t0.append((product) % q)    
        ct0 = intt_vector(ntt_c_t0, n, q, sample_primitive)        

        # Compute hint h 
        h = make_hint_q(-ct0, (w_cs2 + ct0)%q, 2 * gamma_2, n)
        
        # Step 23-24: Validate hint
        norm_ct0 = norm_inf_array(ct0, gamma_2)
        if norm_ct0 or weight(h) > omega:
            z, h = None, None  # Reset and retry
            kappa += 1
            continue

    z = z.tolist()

    # pack sigma
    sigma = pack_sig(c_tilde, z, h, gamma_1)

    # for testig
    '''
    c_tilde1, z1, h1 = unpack_sig(sigma, gamma_1)
    
    assert c_tilde == c_tilde1, "c_tilde1 FAILS"
    if not (z == [[(a2 + q)%q for a2 in a1 ]for a1 in z1]):
        print("z value: ", z)
        print("z1 value: ", z1)
    assert z == [[(a2 + q)%q for a2 in a1 ]for a1 in z1], "z1 FAILS"
    assert h == h1, "h1 FAILS"
    '''
    # end testing

    return sigma
