import numpy as np

#from parameters import *
from updatable_encryption_python.utility_functions import *
from updatable_encryption_python.parameters import *


def compact_error(e_vec):
    qhalf = round(q/2)
    for i in range(len(e_vec)):
        if e_vec[i] > qhalf:
            e_vec[i] = e_vec[i] - q


# generate error vectors
def generate_error_vector(m_tilda, nk, alpha, q, sigma, message_blocks):
    e0 = sample_normal_matrix(1, m_tilda, 0, edisp)[0]
    VerifyE(e0)
    #d = calculate_d(e0, m_tilda, alpha, q, sigma)
    e1 = sample_normal_matrix(1, nk, 0, edisp)[0]
    VerifyE(e1)
    e2 = sample_normal_matrix(1, message_blocks * nk, 0, edisp)[0]
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

# encrypted an encoded message
def Encrypt(pk, encoded_message, H_input, updateMatrix, updateWithMatrix, message_size):
    global s_hidden, e0_hidden, e1_hidden, e2_hidden
    A0, A1, A2 = pk

    # H_input is None for new encryption
    # H_input is not None when used in token_generation for encrypting the 0-message
    if not H_input is None:
        new_Hmu = H_input
    else:
        new_Hmu = generate_invertible_matrixQ(n)

    # Combined matrix
    A_mu = np.block([A0, (A1 + np.dot(new_Hmu, G)) % q])  

    # calculate the number of needed blocks
    message_blocks = message_size // nk
    if message_size % nk != 0:
        message_blocks+=1

    # add A2 matrixes depending on message length
    for _ in range(message_blocks):
        A_mu =  np.block([A_mu, A2])   

    # Random vector
    s = np.random.randint(0, q, n) 
    updateS(s)

    #TO DO choose value of alpha
    e0, e1, e2 = generate_error_vector(m_tilda, nk, alpha, q, edisp, message_blocks)
    


    #print("----------------------------------- e0=", e0)
    #print("----------------------------------- e1=", e1)
    #print("----------------------------------- e2=", e2)
    if updateWithMatrix:
        updateError2(e0, e1, e2, updateMatrix)
    else:
        updateError(e0, e1, e2)
    
    #print("e0=", e0)
    #print("e1=", e1)
    #print("e2=", e2)
    #print("encoded msg: ", encode(message))
    b = (s @ A_mu + np.concatenate((e0, e1, e2)) + np.concatenate((np.zeros(m_tilda), np.zeros(nk), encoded_message))).astype(int) % q 

    return (new_Hmu, b)
