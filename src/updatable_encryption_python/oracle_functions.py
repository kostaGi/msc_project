import numpy as np
from parameters import *

def base_b_decomposition(v, b, k):
    digits = []
    for _ in range(k):
        digits.append(v % b)  # Extract least significant digit
        v //= b  # Integer division to remove extracted digit
    return digits  # Lower indices correspond to lower-order digits

#OracleMeNotExactPower
def oracle_sampleO_not_exact_power(b_hat):

    # https://escholarship.org/uc/item/8b40w7r8 - Algorithm 11:DECODEG(v,b,r[q]kb)

    q_bits = np.zeros(k, dtype=int)
    q_local = q
    for i in range(k):
        q_bits[i] = q_local & 1
        q_local = q_local >> 1

    #print("q_bits=", q_bits)

    v = b_hat
    for i in range(k-1):
        v[i] = 2 * v[i] - v[i+1]

    v[k-1] = 2 * v[k-1]
    #print("B_HAT_L=", v)

    x = np.zeros(k, dtype=int)
    reg = 0

    for i in range(k-1):
        x[i] = round(v[i] / q)
        reg = (reg / 2) + ((2**(k-1)) * q_bits[i])
        v[k-1] = v[k-1] + x[i] * reg

    x[k-1] = round(v[k-1] / (2**k))
    s = x[k-1]
    reg = 0
    i = k-2
    #print("x=", x)
    while i >= 0:
        reg = 2*reg + q_bits[i+1]
        s = s+x[i]*reg
        i = i - 1

    #rint("B HAT Local=",  v, k, s)
    #print("S:", s)
    return s


# OracleMeExactPower
def oracle_sampleO_exact_power(b_hat):

    br1 = q/4
    br2 = (3*q)/4
    #brm = q/2
    blen = len(b_hat)
    #print("B HAT=",  b_hat, blen, "Glen=", len(G), "x", len(G[0]))
    s = 0
    s_bit = 0
    i = blen-1

    while i >= 0:
        bb = (b_hat[i] - s* 2**i) %q
        if bb >= br1 and bb <= br2:
            s_bit = 1
        else:
            s_bit = 0
        #if (b_hat[i]>= brm):
        #    s_bit = 1
        #else:
        #    s_bit = 0

        s = (s + 2**(blen-1-i) * s_bit) % q
        #print("s=", s, s_bit, bb, (blen-1-i), i, 2**i % q)
        i = i - 1

    return s

#OracleMe
def oracle_sampleO(b_hat):

    if not isexactpower:
        return oracle_sampleO_not_exact_power(b_hat)

    return oracle_sampleO_exact_power(b_hat)


# SampleOracleMeExactPower
def oracle_sampleD_exact_power(u, k, b):

    # https://escholarship.org/uc/item/8b40w7r8 - Algorithm 9: g−1(u) for q = b^k

    x = np.zeros(k, dtype=int)
    u_local = u
    for i in range(k):
        y = u_local % b
        if y == 0:
            x[i] = 0
        else:
            if y/b > 1/2:
                x[i] = y - b
            else:
                x[i] = y
        u_local = (u_local - x[i]) / b

    return x

# SampleOracleMeNotExactPower
def oracle_sampleD_not_exact_power(q, q_bits, u, u_bits, k, b):

    # https://escholarship.org/uc/item/8b40w7r8 - Algorithm 10: g−1(u) for q != b^k

    x = np.zeros(k, dtype=int)
    y = np.zeros(k, dtype=int)

    #if 2*u <= q:
    #    x[k-1] = 0
    #else:
    #    x[k-1] = -1

    if np.random.rand() < (q - u) / q:
        x[k-1] = 0
    else:
        x[k-1] = -1

    u_local = u
    q_local = q
    i = k-2
    while i >= 0:
        u_local = u_local - u_bits[i+1] * (b ** (i+1))
        q_local = q_local - q_bits[i+1] * (b ** (i+1))
        c = -(u_local + (x[k-1]*q_local))
        if c < 0:
            p = c + ((b ** (i+1)))
            z = -1
        else:
            p = c 
            z = 0

        #border = p / (b ** (i+1))
        #if (border >= 1/2):
        if np.random.rand() < p / (b ** (i+1)):
            x[i] = z + 1
        else:
            x[i] = z
        i = i - 1

    for i in range(k-1):
        if i == 0:
            y[i] = b * x[i] + x[k-1] * q_bits[i] + u_bits[i]
        else:
            y[i] = b * x[i] - x[i-1] + x[k-1] * q_bits[i] + u_bits[i]
    y[k-1] = -x[k-2] + x[k-1]*q_bits[k-1] + u_bits[k-1] 

    #print(y)
    return y


# SampleOracleV
def oracle_sampleD(UV, q, k, b, ixp):
    q_bits = base_b_decomposition(q, b, k)
    #print(q_bits)

    yarr = np.zeros(k * len(UV), dtype=int) # return array

    for ui in range(len(UV)):
        uvi = UV[ui] % q
        if ixp:
            y = oracle_sampleD_exact_power(uvi, k, b)
        else:
            u_bits = base_b_decomposition(uvi, b, k)
            y = oracle_sampleD_not_exact_power(q, q_bits, uvi, u_bits, k, b)

        for i in range(k):
            yarr[ui*k+i] = y[i]

    #print(tarr)
    return yarr




