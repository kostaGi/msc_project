import numpy as np
from numpy.linalg import norm

# Set the modulus q
q = 11  # Example modulus, can be changed
bound_value = 900  # Reasonable bound for the norm check
d = 64
# Function to simulate discrete Gaussian sampling
def discrete_gaussian_sample(mean, sigma, dimension, q):
    sample = np.random.normal(mean, sigma, dimension).round().astype(int)
    return sample % q  # Ensure the sample is modulo q

# Function to generate a deterministic challenge
def fiat_shamir_challenge(data, C):
    hash_int = sum(data) % len(C)  # Simple sum-based hash function
    return C[hash_int % len(C)]

# Function to flatten a 2D array
def flatten_2d_array(array):
    return [item for sublist in array for item in sublist]

# Function to sign the message
def sign(m, A1, A2, Bext, s1, s2, C, q):
    n = len(s1)  # Dimension of the secret vectors

    # Prover's side: Sampling y1, y2
    y1 = discrete_gaussian_sample(0, 1, n, q)
    y2 = discrete_gaussian_sample(0, 1, n, q)

    # Compute w
    w = (np.dot(A1, y1) + np.dot(A2, y2)) % q

    # Convert w and message to a simple data representation for hashing
    m_flat = flatten_2d_array(m)  # Flatten the 2D message
    w_data = list(w) + m_flat  # Concatenate w and flattened m into a list

    # Fiat-Shamir challenge
    c = fiat_shamir_challenge(w_data, C)

    # Prover's response to challenge
    z1 = (y1 + c * s1) % q
    z2 = (y2 + c * s2) % q

    # Prover constructs t_A and t_B
    t_A = (A1.dot(s1) + A2.dot(s2)) % q
    t_B = (Bext.dot(s2) + m.sum(axis=0)) % q  # Sum the columns of the 2D message

    # Prover sends (w, z1, z2) as the signature
    signature = (w, z1, z2)
    return signature, t_A, t_B

# Function to verify the signature
def verify_signature(signature, m, t_A, C, A1, A2, Bext, q):
    w, z1, z2 = signature

    # Convert w and message to a simple data representation for hashing
    m_flat = flatten_2d_array(m)  # Flatten the 2D message
    w_data = list(w) + m_flat  # Concatenate w and flattened m into a list

    # Fiat-Shamir challenge
    c = fiat_shamir_challenge(w_data, C)
    print("z1 norm: ", norm(z1))
    print("z2 norm: ", norm(z2))

    # Verification
    if norm(z1) <= 1.5 * bound_value and norm(z1) >= bound_value and norm(z2) <= 1.5 * bound_value and norm(z2) > bound_value:  # Check the bounds for z1 and z2
        # Compute w_prime
        w_prime = (A1.dot(z1) + A2.dot(z2) - c * t_A) % q
        w_prime = np.mod(w_prime, q)  # Ensure result is within the correct range
        print("w: ", w)
        print("w_prime:", w_prime)

        # Verification check
        return np.array_equal(w, w_prime)
    else:
        return False

'''
# Example usage
n = 6  # Dimension of the LWE problem
A1 = np.random.randint(0, q, size=(d, d)) % q
A2 = np.random.randint(0, q, size=(d, d)) % q  # Example for A2, adjust as needed
Bext = np.random.randint(0, q, size=(d, d)) % q  # Example for Bext, adjust as needed

s1 = np.random.randint(0, q, size=d) % q
s2 = np.random.randint(0, q, size=d) % q

m = np.random.randint(0, q, size=(d, d))  # The message m as a 2D array
C = np.array([1, -1])

signature, t_A, t_B = sign(m, A1, A2, Bext, s1, s2, C, q)
print("Signature:", signature)

verification_passed = verify_signature(signature, m, t_A, C, A1, A2, Bext, q)
print("Verification passed:", verification_passed)
'''