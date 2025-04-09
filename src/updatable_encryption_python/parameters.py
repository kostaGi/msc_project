import numpy as np

# requirements
# n ≥ 1
# q ≥ 2
# q = poly(λ) - means q should be big so its hard to break
# m ≥ nk ≥ n
# k = ⌈log2q⌉
# m = ¯m+2nk

n = 16
q = 1048573 #8380417 #2647 #2048 #1979  #2647   #1979
k = int(np.ceil(np.log2(q)))
print("k = ", k)
isexactpower = (q == 2**k)
nk = n*k
m_tilda = 32 #nk #* 2 #10
m = m_tilda + 2*nk
alpha = 1/10
s_hidden = np.zeros(n, dtype=int)  #used for testing
e0_hidden = np.zeros(m_tilda, dtype=int)  #used for testing
e1_hidden = np.zeros(nk, dtype=int)  #used for testing
e2_hidden = np.zeros(nk, dtype=int)  #used for testing

e00 = np.zeros(m_tilda, dtype=int)  #used for testing
e01 = np.zeros(nk, dtype=int)  #used for testing
e02 = np.zeros(nk, dtype=int)  #used for testing

def mod_q_vector(vector):
    for current_element in range(len(vector)):

        vector[current_element] = vector[current_element] % q
        while vector[current_element] > q/2:
            vector[current_element] -= q
        while vector[current_element] < -q/2:
            vector[current_element] += q   
    return vector

def updateError(e0, e1, e2):
    global e0_hidden, e1_hidden, e2_hidden
    e0_hidden = e0
    e1_hidden = e1
    e2_hidden = e2
        

def updateError2(e0, e1, e2, M):
    global e0_hidden, e1_hidden, e2_hidden, e00, e01,e02
    helpe1 = np.concatenate((e0_hidden, e1_hidden, e2_hidden))
    helpe1 = np.dot(helpe1, M % q)
    e00 = mod_q_vector(helpe1[:m_tilda])
    e01 = mod_q_vector(helpe1[m_tilda:m_tilda+nk])
    e02 = mod_q_vector(helpe1[m_tilda+nk:])
    e0_hidden = mod_q_vector(helpe1[:m_tilda] + e0)
    e1_hidden = mod_q_vector(helpe1[m_tilda:m_tilda+nk] + e1)
    e2_hidden = mod_q_vector(helpe1[m_tilda+nk:] + e2)



def getError():
    global e0_hidden, e1_hidden, e2_hidden
    return (e0_hidden, e1_hidden, e2_hidden) 

def getError2():
    global e00, e01,e02
    return (e00, e01, e02) 

def updateS(s):
    global s_hidden
    s_hidden = (s + s_hidden) % q
    print("s_hidden: ", s_hidden)

def getS():
    global s_hidden
    return s_hidden    

# normal distribution with following params
mean = 0
sigma_sk =  0.7 # 1 * np.sqrt(np.log(n))  # R dispersion, included in Rte(<= q/4) affecting correct sampling and decoding
edisp = 0.7 # error dispersion , included in Rte(<= q/4) affecting correct sampling and decoding
tau_sample = 2 / (2+1)  # s / (b+1)  recommended by LWE algorithm , affecting correct sampling and decoding

def generate_G(n, k, q):
    # Construct the vector g^t = [1, 2, 4, ..., 2^(k-1)]
    g_t = np.array([2**i for i in range(k)])  # Shape: (k,)
    # Compute G using the Kronecker product G = I_n ⊗ g^t
    I_n = np.eye(n, dtype=int)  # Identity matrix of size (n, n)
    G = np.kron(I_n, g_t)  # Shape: (n, nk)
    #print("Matrix G:\n", (G))
    #type = <class 'numpy.ndarray'>
    return G % q

G = generate_G(n, k, q)