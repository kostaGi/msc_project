import numpy as np

#from parameters import *
from utility_functions import *
import parameters


def compact_error(e_vec):
    qhalf = round(q/2)
    for i in range(len(e_vec)):
        if e_vec[i] > qhalf:
            e_vec[i] = e_vec[i] - q


# generate error vectors
def generate_error_vector(m_tilda, nk, alpha, q, sigma):
    e0 = sample_normal_matrix(1, m_tilda, 0, edisp)[0]
    VerifyE(e0)
    #d = calculate_d(e0, m_tilda, alpha, q, sigma)
    e1 = sample_normal_matrix(1, nk, 0, edisp)[0]
    VerifyE(e1)
    e2 = sample_normal_matrix(1, nk, 0, edisp)[0]
    #print("e2: " , e2)
    VerifyE(e2)
    #parameters.updateError(e0, e1, e2)
    return (e0, e1, e2)
    #return (np.zeros(m_tilda), np.zeros(nk), np.zeros(nk))

# calculate d    
def calculate_d(e0, m_tilda, alpha, q, sigma):
    sum = 0
    for value in e0:
        sum+= value**2
    #print("sum=", sum, "m_tilda=", m_tilda*((alpha*q)**2), "sigma=", sigma**2)
    d = np.sqrt(sum + m_tilda*((alpha*q)**2) * (sigma**2))    
    return d





def Encrypt(pk, encoded_message, H_input, updateMatrix, updateWithMatrix):
    global s_hidden, e0_hidden, e1_hidden, e2_hidden
    A0, A1, A2 = pk
    # Invertable matrix hardcoded for now
    #H_mu = generate_invertible_matrix(n, q)


   
    # Combined matrix
    HmG = np.dot(H_input, G) % q
    #print("HmG=", HmG, len(HmG), len(HmG[0]))

    A_mu = np.block([A0, A1 + HmG, A2])  
    #A_m_0_1 = np.block([A0, A1 + HmG])  
    #print(A_mu, len(A_mu), len(A_mu[0]))

    # Rndom vector
    s = np.random.randint(0, q, n) 
    updateS(s)

    #TO DO choose value of alpha
    e0, e1, e2 = generate_error_vector(m_tilda, nk, alpha, q, edisp)
    


    print("----------------------------------- e0=", e0)
    print("----------------------------------- e1=", e1)
    print("----------------------------------- e2=", e2)
    if updateWithMatrix:
        updateError2(e0, e1, e2, updateMatrix)
    else:
        updateError(e0, e1, e2)
    
    #print("e0=", e0)
    #print("e1=", e1)
    #print("e2=", e2)
    #print("encoded msg: ", encode(message))
    b = (s @ A_mu + np.concatenate((e0, e1, e2)) + np.concatenate((np.zeros(m_tilda), np.zeros(nk), encoded_message))).astype(int) % q 

    return (H_input, b)
