import numpy as np

from parameters import *

def sample_uniform_matrix(rows, cols, q):
    return np.random.randint(0, q, (rows, cols))

def sample_normal_matrix(rows, cols, mean, sigma):
    # Sample from a continuous Gaussian distribution
    samples = np.random.normal(loc=mean, scale=sigma, size=(rows, cols))
    # Round to the nearest integer to make it discrete
    discrete_samples = np.round(samples).astype(int)
    return discrete_samples



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
    print("Rte=",Rte)
    VerifyE(Rte)

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
'''
#NEW TO FIX
B_enc = np.random.randint(-q//2, q//2, size=(n * k, n * k))

def encode(message):
    print("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", n * k)
    """
    Encode a binary message vector into a lattice point.
    
    :param m: Binary message of length n*k (1D numpy array of 0s and 1s)
    :return: Encoded lattice point in Z^(n*k)
    """
    if len(message) != n * k:
        raise ValueError(f"Message must be of length {n*k}")
    
    return (B_enc @ message)  # Linear transformation into the lattice

def decode(encoded):
    """
    Decode a lattice point back into the original binary message.
    
    :param v: Lattice vector in Z^(n*k)
    :return: Decoded binary message
    """
    # Solve Bm ≈ v using the pseudoinverse (assuming B is full rank)
    m_approx = np.linalg.pinv(B_enc) @ encoded

    # Round to nearest binary values (0 or 1)
    m_decoded = np.round(m_approx).astype(int) % 2


    #print("m_decoded=", m_decoded)
    return m_decoded
'''

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