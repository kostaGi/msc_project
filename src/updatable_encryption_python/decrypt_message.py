import numpy as np
from parameters import *
from oracle_functions import *
from utility_functions import *
from encrypt_message import compact_error



def InvertO(R, A_0_1, b_0_1, H_2, doerrch):
    print("--- InvertO ---")


    R_padded = np.block([[R], [np.eye(nk, dtype=int)]])


    b_0_1_hat = np.dot(b_0_1, R_padded).astype(int) % q

    print(n, k, nk, "b_0_1=",len(b_0_1), ", R=", len(R), "x", len(R[0]), ", R_padded=", len(R_padded), "x", len(R_padded[0]), ", b_0_1_hat=", len(b_0_1_hat))
    s_m = np.zeros(n, dtype=int)
    for i in range(n):
        s_m[i] = oracle_sampleO(b_0_1_hat[i*k:k*(i+1)])
    print("s_m=", s_m)

    H_2_inv = np.linalg.inv(H_2).astype(int)
    s_calc = np.dot(s_m, H_2_inv) % q
    e_calc = b_0_1 - np.dot(s_calc, A_0_1).astype(int)    

    e0_hidden, e1_hidden, e2_hidden = getError()
    
    #if assertError:
    print("s_calc: ", s_calc)
    print("s_hidd: ", getS())
    assert np.array_equal(getS(), s_calc)

    #print("e0_hidden=", e0_hidden %q)
    #print("e_calc[:m_tilda]=", e_calc[:m_tilda] % q)
    e0_calc = e_calc[:m_tilda] % q
    compact_error(e0_calc)
    e1_calc = e_calc[m_tilda:] % q
    compact_error(e1_calc)


    if doerrch:
        print("InvertO=", e0_hidden)
        print("InvertO=", e0_calc)
        assert np.array_equal(e0_hidden, e0_calc)
        assert np.array_equal(e1_hidden, e1_calc)

    return s_calc, e_calc



# Decrypt
def Decrypt(sk, H_mu, b, pk, doerrch):

    e0_hidden, e1_hidden, e2_hidden = getError()

    A0, A1, A2 = pk  # Extract public key components
    # calculate s and e01
    s_calc, e01_calc = InvertO(sk, np.block([A0, A1 + np.dot(H_mu, G)]), b[:m-nk], H_mu, doerrch)

    # old to be changed
    helper1 = ( b[m_tilda+nk:] - np.dot(s_calc, A2) % q ) % q
    decoded_message, e2_calc = decode(helper1)

    #new Todo
    #decoded_message = decode(helper1)
    #print("decoded_message: ", decoded_message)
    #e2_calc = ( b[m_tilda+nk:] - np.dot(s_calc, A2) - encode(decoded_message)) % q


    if doerrch:
        e0h, e1h, e2h = getError()
        compact_error(e2_calc)
        print("e2Hidden: ", e2h)
        print("e2_calc: ", e2_calc)
        assert np.array_equal(e2h, e2_calc)
        
    return decoded_message