import numpy as np

from updatable_encryption_python.parameters import *
from updatable_encryption_python.oracle_functions import *
from updatable_encryption_python.utility_functions import *
from updatable_encryption_python.encrypt_message import *

def SampleDO(R1, A01m, H_mu, A_part_prime, tau):

    print("---------- SampleDO ----------")
    wwr,hhr = R1.shape
    wwa,hha = A01m.shape
    wwp,hhp = A_part_prime.shape
    print("R1, A, AP=", wwr, hhr, wwa, hha, wwp, hhp)

    #scalling_d = np.zeros(n)
    #for i in range(n):
    #    scalling_d[i] = 1 # Using uniform d for simplicity

    x = np.zeros((hhp, hha))
    #print("x=", x.shape)
    H_hardcoded_inv = inverse_matrixQ(H_mu, len(H_mu))
    
    #s = 5.0  # Example scaling factor
    #r = 3.0  # Rounding parameter

    # Define a simple covariance matrix for G (you may need to replace it with an actual one)
    #Sigma_G = np.eye(n * k) * r

    A_part_prime_t = A_part_prime.T


    for i in range(hhp):
        # Generate perturbation

        p = sample_perturbation_v2(R1, tau, A01m.shape[1])
        #p = sample_perturbation(R1, s, r, Sigma_G, n, k, m_tilda, q)
        #print(p, p.shape)
        u_tmp = np.dot(A01m, p) % q
        #print(u_tmp.shape, A_part_prime_t[i].shape)
        #print(A_part_prime_t[i])

        u_tmp = A_part_prime_t[i] - u_tmp.flatten()
        #print("u_tmp=", u_tmp, u_tmp.shape)
        v = np.dot(H_hardcoded_inv, u_tmp) % q
        #v = v.reshape(-1,1).T
        #print("v=", v)

        z = oracle_sampleD(v, q, k, 2, isexactpower)

        #print("z=",z)
        #s = 15.0 # q not exact power 3.0  # Example scaling factor
        #z = sampleGV(s, v, 2, q, k)
        #print("sampleGV=",z.shape)
        #z = np.zeros(nk)
        #for zi in range(n):
        #    t = sampleG(s, v[zi], 2, q, k) 
        #    for zj in range(k):
        #        z[zi*k+zj] = t[zj]

        R_ext = np.block([ [R1], [np.eye(nk, dtype=int)] ])
        #print(R_ext.shape, z.shape)
        x_vec = np.dot(R_ext, z) % q
        #print(i, x_vec.shape, p.shape)
        x_vec = p.flatten() + x_vec
        #print(i, x_vec.shape)
        x_vec = x_vec.astype(int)  % q
        #print(i, x_vec.shape)
        for j in range(len(x_vec)):
            x[i][j] = x_vec[j]
        
    #print("xshape=", x.shape)
    x = x.T
    x = x.astype(int) % q
    #print("xshape=", x.shape)
    return x


def TokenGen(pk, sk, pk_prime, Hmu, message_size):
    #generate new Hmu_prime
    Hmu_prime = generate_invertible_matrixQ(n)
    A0, A1, A2 = pk
    A0_prime, A1_prime, A2_prime = pk_prime
    Hmu_primeG = np.dot(Hmu_prime, G) % q
    A_mu_prime = np.block([A0_prime, A1_prime + Hmu_primeG, A2_prime]) 

    HmG = np.dot(Hmu, G) % q
    A01m = np.block([A0, A1 + HmG]) 
    X0 = SampleDO(sk, A01m, Hmu, A0_prime, tau_sample)
    X1 = SampleDO(sk, A01m, Hmu, A1_prime + Hmu_primeG, tau_sample*np.sqrt(m_tilda/2))
    X2 = SampleDO(sk, A01m, Hmu, A2_prime-A2, tau_sample*np.sqrt(m_tilda/2))

    #print(A1_prime + Hmu_primeG)

    # verify X0
    #print("X0, X1, X2, A01 =", X0.shape, X1.shape, X2.shape, A01m.shape)
    V0 = np.dot(A01m, X0) % q
    V1 = np.dot(A01m, X1) % q
    V2 = np.dot(A01m, X2) % q
    print("V0, V1, V2 =", V0.shape, V1.shape, V2.shape)

    #print("A0_prime=", A0_prime)
    #print("V0=", V0)
    #print("A1_prime=", A1_prime)
    #print("V1=", V1)

    #print("V0, V1, V2, A0P =", V0.shape, V1.shape, V2.shape, A0_prime.shape)
    #print(V0)
    #print(A0_prime)

    if not np.array_equal(V0 , A0_prime % q):
        print("Caclculated: V0", V0, "\n\n\n")
        print("Caclculated: Test1", A0_prime % q, "\n\n\n")

    assert np.array_equal(V0, A0_prime % q)


    if not np.array_equal(V2, (A2_prime - A2) % q):
        print("Caclculated: V2", V2, "\n\n\n")  
        print("Caclculated: Test3", (A2_prime - A2) % q, "\n\n\n")

    assert np.array_equal(V2, (A2_prime - A2) % q)



    if not np.array_equal(V1, (A1_prime + Hmu_primeG) % q):
        print("Caclculated: V1", V1, "\n\n\n")  
        print("Caclculated: Test2", (A1_prime + Hmu_primeG) % q, "\n\n\n")

    assert np.array_equal(V1, (A1_prime + Hmu_primeG) % q)

    

    #print(A1)
    #print(A1_prime)
    #assert np.array_equal(A1 % q, A1_prime % q)
    
    # returning matrix M , last row??? 
    Amu = np.block([A0, A1 + HmG, A2]) 
    M = np.block([[X0, X1, X2],[np.zeros((nk, m_tilda)), np.zeros((nk, nk)), np.eye((nk))]]).astype(int)

    assert np.array_equal(np.dot(Amu, M) % q, A_mu_prime % q)

    #print("np.dot(Amu, M): ", np.dot(Amu, M)% q)
    #print("A_mu_prime: ", A_mu_prime % q)

    e0h, e1h, e2h = getError()
    print("e0h: ", e0h)
    print("e1h: ", e1h)
    print("e2h: ", e2h)
    # incorrect H (added it as parameter to change)

    message_blocks = message_size // nk
    if message_size % nk != 0:
        message_blocks+=1

    # to put as true if you want to check values
    _, b_zero_message = Encrypt(pk_prime, np.zeros(message_blocks* nk), Hmu_prime, M, False, message_blocks * nk)
    #b_zero_message = 0
    #e0ha, e1ha, e2ha = getError()
    #Lemma1Check(sk, np.concatenate((e0ha, e1ha)))

    ## Update error moved to Encrypt

    #eb = np.dot(np.concatenate((e0h, e1h, e2h)), M) % q
    #compact_error(eb)
            
    #print("eb=", eb)
    #e0b = eb[:m_tilda]
    #e1b = eb[m_tilda:m_tilda+nk]
    #e2b = eb[m_tilda+nk:]
    #updateError(e0b, e1b, e2b)

    #print("e0b=", e0b)
    #print("e1b=", e1b)
    #Lemma1Check(sk, np.concatenate((e0b, e1b)))


    return X0, X1, X2, b_zero_message, Hmu_prime

    #return (M, b_zero_message, Hmu_prime)

def update(token, c):
    M, b_zero_message, Hmu_prime = token
    Hmu , b = c


    #print("update b=",b, b.shape)
    #print("update M=",M, M.shape)

    e00, e01, e02 = getError2()
    helper1 = np.concatenate((e00, e01, e02))
 
    b_prime = ((np.dot(b, M)) + b_zero_message) % q # - helper1 ) % q
    
    print("b_prime shape: ", b_prime.shape)

    return (b_prime, Hmu_prime)

