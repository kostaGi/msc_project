import numpy as np

# requirements
# n ≥ 1
# q ≥ 2
# q = poly(λ) - means q should be big so its hard to break
# m ≥ nk ≥ n
# k = ⌈log2q⌉
# m = ¯m+2nk

n = 8
q = 8380417 #2647 #2048 #1979  #2647   #1979
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

def mod_q(vector):
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
    e00 = mod_q(helpe1[:m_tilda])
    e01 = mod_q(helpe1[m_tilda:m_tilda+nk])
    e02 = mod_q(helpe1[m_tilda+nk:])
    e0_hidden = mod_q(helpe1[:m_tilda] + e0)
    e1_hidden = mod_q(helpe1[m_tilda:m_tilda+nk] + e1)
    e2_hidden = mod_q(helpe1[m_tilda+nk:] + e2)



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

'''
# q= 23
H_hardcoded = np.array([[21 , 0 , 3 ,21 , 0 , 1 , 2 ,20],
[ 1 , 2 , 3 , 1 , 0 , 0 , 1 , 2],
 [ 0 , 0, 22 , 0 , 0 , 0 , 0 , 1],
 [22 ,14 , 7 ,19 , 1 , 1 , 2 ,18],
 [22 ,19, 20 ,19 , 1 , 1 ,22, 18],
 [ 2,  2 , 5 , 1 , 0 , 0 , 0 , 1],
 [ 0 , 0 , 0 , 0 , 0 , 0 , 1 , 0],
 [ 0 , 2 , 5 , 0 , 0 , 0, 21,  0]])
'''

'''
# q = 32
H_hardcoded = np.array([[ 0,  0, 26,  0,  0,  0,  0,  1],
 [ 0,  0,  1,  0,  0,  0,  0,  0],
 [ 1,  2,  3,  0,  0,  0,  0, 31],
 [ 9,  3, 12,  0,  1, 29,  0,  2],
 [ 1,  2, 21,  1,  0, 31,  0,  1],
 [ 6,  9, 14,  1,  0, 29,  1,  0],
 [ 6,  2,  8,  0,  1, 30,  0,  2],
 [ 0,  1, 30,  0,  0,  0,  0,  0]])
'''

'''
#q = any
H_hardcoded = np.array([
 [ 1,  0,  0,  0,  0,  0,  0,  0],
 [ 0,  1,  0,  0,  0,  0,  0,  0],
 [ 0,  0,  1,  0,  0,  0,  0,  0],
 [ 0,  0,  0,  1,  0,  0,  0,  0],
 [ 0,  0,  0,  0,  1,  0,  0,  0],
 [ 0,  0,  0,  0,  0,  1,  0,  0],
 [ 0,  0,  0,  0,  0,  0,  1,  0],
 [ 0,  0,  0,  0,  0,  0,  0,  1],
 ])

H_hardcoded_prime = np.array([
 [ 1,  0,  0,  0,  0,  0,  0,  0],
 [ 0,  1,  0,  0,  0,  0,  0,  0],
 [ 0,  0,  1,  0,  0,  0,  0,  0],
 [ 0,  0,  0,  1,  0,  0,  0,  0],
 [ 0,  0,  0,  0,  1,  0,  0,  0],
 [ 0,  0,  0,  0,  0,  1,  0,  0],
 [ 0,  0,  0,  0,  0,  0,  1,  0],
 [ 0,  0,  0,  0,  0,  0,  0,  1],
 ])
'''


'''
# q = 97
H_hardcoded = np.array([
 [ 0, 93,  7,  5,  6, 29,  2, 58],
 [87, 85, 94,  9,  1, 51,  0, 93],
 [ 4,  3,  0, 94,  0, 83,  0, 71],
 [ 7,  5,  0, 92,  0, 72,  0, 50],
 [ 0,  0,  0,  0,  0,  0,  0,  1],
 [ 6,  3,  3, 94,  3, 84,  1, 74],
 [ 8,  3,  6, 95,  6, 89,  2, 87],
 [ 0,  0,  0,  0,  0,  1,  0,  2],
 ])
 '''

'''
# q = 128
H_hardcoded = np.array([
 [  1,   0,   0,   1,   0,   0,   0,   0],
 [125,   0,   0, 127,   7,   0,   0,   1],
 [ 20, 126,   1,   5, 103,   0, 127, 119],
 [  9,   1,   0,   6,   6, 122,   0,   0],
 [124,   0,   0, 125, 126,   1,   0,   0],
 [100,   2,   0, 119,  35,   1,   1,   9],
 [125,   0,   0, 126, 124,   2,   0,   0],
 [ 81,   4,   0, 113,  50,   2,   2,  15],
])
'''

'''
# q = 2048
H_hardcoded = np.array([
 [   0,    0,    2,    0, 2046,    0,    0,    1],
 [   0,    1, 2046,    0,    2, 2047,    0, 2046],
 [   0,    0,    5,    0, 2043,    0,    0,    2],
 [   0,    1, 2047,    1, 2047,    1,    0,    2],
 [   0,    0,    0,    0,    1,    0,    0,    0],
 [   1,    3,    0,    0,    0,    0,    3,    0],
 [   0,    0,    2,    0, 2046,    1,    0,    2],
 [   1,    2,    6, 2046,    1,    2,    4,    3],
])
'''

'''
# q = 1237
H_hardcoded = np.array([
 [  12,    0,    2,   12,   35,   11,    0,   20],
 [1210,    0, 1220, 1215, 1139, 1207,    1, 1185],
 [   0,    1,    3,    1, 1231,    0,    0, 1236],
 [   0,    1,    4,    1, 1231,    0,    0, 1236],
 [  30,    0,    5,   30,   86,   27,    0,   50],
 [1225,    0, 1232, 1226, 1196, 1224,    0, 1216],
 [   1,    0,    0,    1,    0,    0,    0,    2],
 [   0,    0,    2,    0,    0,    0,    0,    1],
])
'''


# q = 2647
H_hardcoded = np.array([
 [   0, 2646,    0,    0,    0,    0, 2646,    1],
 [   1, 2637,    5,    1, 2645, 2644, 2633,    0],
 [   0,    0,    1,    0, 2645, 2645,    0,    0],
 [   0,    1,    0,    0,    0,    0,    1,    0],
 [   1, 2641,    1,    1,    5,    3, 2638,    0],
 [   0, 2646, 2646,    0,    0,    2, 2645,    0],
 [   0,    0,    0,    0,    1,    1,    0,    0],
 [   0,    3, 2646,    1,    0,    0,    3,    0],
])
H_hardcoded_prime = np.array([
 [ 574,   61,  386,  157, 1190, 1136,    6,    2],
 [ 240,   27,  175,   66,  552,  535,    3,    0],
 [ 272,   30,  199,   77,  624,  609,    3,    1],
 [2640,    0,   22,    1,   87,  109,    0,    0],
 [  75,   10,   51,   25,  151,  131,    1,    0],
 [   1,    0,    2,    0,    8,   10,    0,    0],
 [  15,    0,   15,    0,   57,   71,    0,    0],
 [ 946,   93,  596,  240, 1830, 1745,    9,    3],
])