import numpy as np
from updatable_encryption_python.parameters import *
from updatable_encryption_python.utility_functions import *

def KeyGen(sk1):
    # uniform distribution with lower and upper bound q
    A0 = sample_uniform_matrix(n, m_tilda, q)

    # normal distribution
    if sk1 is None:
        R1 = sample_normal_matrix(m_tilda, nk, mean, sigma_sk)
        VerifyR(R1)
    else:
        R1 = sk1
    R2 = sample_normal_matrix(m_tilda, nk, mean, sigma_sk)
    VerifyR(R2)
    
    A1 = -np.dot(A0, R1) % q
    A2 = -np.dot(A0, R2) % q
    
    pk = [A0, A1, A2]  # Public Key
    #pk_new = np.block([A0, A1, A2])
    sk = R1

    # to be removed
    # verify trapdoor
    R = np.block([R1, R2])
    pk_local = np.block([A0, A1, A2])
    R_ext = np.block([[R], [np.eye(R.shape[1], dtype=int)] ])
    Res = np.dot(pk_local, R_ext) % q
    all_zeros = not np.any(Res)
    assert all_zeros

    #print("KeyGen sk=", sk)
    return pk, sk
