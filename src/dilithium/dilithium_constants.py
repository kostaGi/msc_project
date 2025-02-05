#outside libraries
import numpy as np

# Global Parameters

# Default Security level parameters (level 2)
q = 8380417  # Prime modulus  q = 2**23 - 2**13 + 1
n = 256  # Polynomial degree
l = 4  # Parameter defining s1's dimension
k = 4  # Parameter defining s2's dimension
eta = 2  # Bound on coefficients of s1, s2
d = 13  # Precision for power-of-2 rounding
tau = 39 # Numbers of coefficients that are either −1 or 1 in sample_in_ball
gamma_1 = 131072 # Bound on the coefficients in expand_mask gamma_1= 2**17
gamma_2 = 95232  # Bound for high bits gamma_2 = (q-1)//88
beta = tau * eta # Parameter defining offset for bounds
omega = 80 # Maximum allowed weight of the hint h




### Following methods generated variable for the ntt computation


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

sample_primitive = find_primitive_root(n, q)

### Following function sets security level

# Set security level
def set_dilithium_security_level(security_level):
    if security_level == 2:
        q = 8380417  
        n = 256  
        l = 4  
        k = 4  
        eta = 2 
        d = 13  
        tau = 39 
        gamma_1 = 131072
        gamma_2 = 95232  
        beta = tau * eta 
        omega = 80 
    elif security_level == 3:
        q = 8380417 
        n = 256  
        l = 5  
        k = 6  
        eta = 4 
        d = 13 
        tau = 39 
        gamma_1 = 524288
        gamma_2 = 261888
        beta = tau * eta 
        omega = 55
    elif security_level == 5:
        q = 8380417  
        n = 256  
        l = 7 
        k = 8  
        eta = 2  
        d = 13  
        tau = 60
        gamma_1 = 524288 
        gamma_2 = 261888
        beta = tau * eta 
        omega = 75
    # sample primitive for the ntt computation
    sample_primitive = find_primitive_root(n, q)
    return 1



