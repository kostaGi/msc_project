#outside libraries
import numpy as np

# Parameters

# Security level parameters (current security level 2)
q = 8380417  # Prime modulus  q = 2**23 - 2**13 + 1
n = 256  # Polynomial degree
l = 4  # Parameter defining s1's dimension
k = 4  # Parameter defining s2's dimension
eta = 2  # Bound on coefficients of s1, s2
d = 13  # Precision for power-of-2 rounding
tau = 39 #
gamma1 = 131072 # gamma1 = 2**17
gamma2 = 95232  # gamma2 = (q-1)//88
beta = tau * eta #
omega = 80 #

'''
# Security level parameters (current security level 3)
q = 8380417  # Prime modulus  q = 2**23 - 2**13 + 1
n = 256  # Polynomial degree
l = 5  # Parameter defining s1's dimension
k = 6  # Parameter defining s2's dimension
eta = 4  # Bound on coefficients of s1, s2
d = 13  # Precision for power-of-2 rounding
tau = 39 #
gamma1 = 524288 # gamma1 = 2**17
gamma2 = 261888  # gamma2 = (q-1)//88
beta = tau * eta #
omega = 55 #
'''

'''
# Security level parameters (current security level 5)
q = 8380417  # Prime modulus  q = 2**23 - 2**13 + 1
n = 256  # Polynomial degree
l = 7  # Parameter defining s1's dimension
k = 8  # Parameter defining s2's dimension
eta = 2  # Bound on coefficients of s1, s2
d = 13  # Precision for power-of-2 rounding
tau = 60 # 
gamma1 = 524288 # gamma1 = 2**17
gamma2 = 261888  # gamma2 = (q-1)//88
beta = tau * eta #
omega = 75 #
'''


### Following variable is generated for the ntt computation ones


# Check if g is a primitive n-th root of unity modulo q
def is_primitive_root(g, n, q):
    if pow(g, n, q) != 1:
        return False
    # Check that g^k != 1 for 1 <= k < n
    for k in range(1, n):
        if pow(g, k, q) == 1:
            return False
    return True


# Find a primitive n-th root of unity modulo q
def find_primitive_root(n, q):
    # n must divide q-1
    assert (q - 1) % n == 0
    for g in range(2, q):
        if is_primitive_root(g, n, q):
            return g
    raise ValueError(f"No primitive {n}-th root of unity found modulo {q}")


# sample primitive for the ntt computation
sample_primitive = find_primitive_root(n, q)
