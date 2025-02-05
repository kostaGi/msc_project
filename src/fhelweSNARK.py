import galois
import numpy as np
import fromLWEtoR1CS as r1
import fromR1CStoQAP as qap
import signing
import timeit
import Witnessfinder as wf
# import pickle

#parameters
n = 6
order = 701
GF = galois.GF(order)
q = 2 ** 8
t = order

num_var = 3*n**2 + 2*n + 1
d = 64
delta = q // t
p_q = np.poly1d([1] + ([0] * (d - 1)) + [1])

#signing keys
A1 = np.random.randint(0, q, size=(d, d)) % q
A2 = np.random.randint(0, q, size=(d, d)) % q  # Example for A2, adjust as needed
Bext = np.random.randint(0, q, size=(d, d)) % q  # Example for Bext, adjust as needed

s1 = np.random.randint(0, q, size=d) % q
s2 = np.random.randint(0, q, size=d) % q
# print("origional: ", len(origional[0]))
C = np.array([1, -1])

def polymod(poly, poly_mod, coeff_mod):
    return np.poly1d(np.floor(np.polydiv(poly, poly_mod)[1]) % coeff_mod)


def addition(poly_mod, coeff_mod):
    return lambda a, b: np.poly1d(polymod(np.polyadd(a, b), poly_mod, coeff_mod))


def multiplication(poly_mod, coeff_mod):
    return lambda a, b: np.poly1d(polymod(np.polymul(a, b), poly_mod, coeff_mod))


add = addition(p_q, q)
mul = multiplication(p_q, q)


def setup():
    alpha = np.random.randint(0, t, 1)[0]
    sk = np.poly1d(np.random.randint(0, 2, d))
    a2 = np.poly1d(np.random.randint(0, q, d, dtype=np.int64) % q)
    e = np.poly1d(np.random.normal(0, 2, d).astype(int) % q)
    pk_0 = add(-mul(a2, sk), e)
    pk_1 = a2
    pk = (pk_0, pk_1)
    u = np.poly1d(np.random.randint(0, 2, d))
    e1 = np.poly1d(np.random.normal(0, 2, d).astype(int) % q)
    e2 = np.poly1d(np.random.normal(0, 2, d).astype(int) % q)
    return alpha, sk, a2, e, pk, u, e1, e2


def prover(pk, u, e1, alpha, s):
    A, B, C = r1.LWEToR1CS_transform()
    print("r1cs completed")
    U, V, W, Ua, Va, Wa = qap.polySum(GF(s))

    print("polysum done")
    T = galois.Poly([1, order - 1], field=GF)
    for i in range(2, len(A) + 1):
        T = T * galois.Poly([1, order - i], field=GF)

    H = (Ua * Va - Wa) // T

    left = Ua(alpha) * Va(alpha) - Wa(alpha)
    right = H(alpha) * T(alpha)

    c1 = add(add(mul(pk[0], u), e1), mul(delta, left))
    c2 = add(add(mul(pk[0], u), e1), mul(delta, right))
    proof_package = [c1, c2]
    proof = np.array([p.coeffs.tolist() for p in proof_package]).astype(int)
    proof_flat = proof.flatten()

    # Convert to hexadecimal string
    hex_string = ''.join(format(x, '02x') for x in proof_flat)
    return hex_string, proof.shape, proof


def hex_to_array(hex_str, shape):
    # Each integer is represented by 2 hex characters (assuming 8-bit integers)
    num_elements = len(hex_str) // 2
    # Convert the hex string back to a list of integers
    int_list = [int(hex_str[i:i + 2], 16) for i in range(0, len(hex_str), 2)]
    # Convert the list of integers to a NumPy array and reshape to the original shape
    return np.array(int_list).reshape(shape)


def verifier(proof, shape):
    proof_array_reconstructed = hex_to_array(proof, shape)

    return np.all(proof_array_reconstructed[0] == proof_array_reconstructed[1])


def main():
    total_prover_time = 0
    total_verifier_time = 0
    total_sign_time = 0
    total_sig_verify_time = 0
    num_runs = 100

    for _ in range(num_runs):
        s = wf.witnessGenerator()
        alpha, sk, a2, e, pk, u, e1, e2 = setup()

        #generate proof
        start_prover = timeit.default_timer()
        proof, shape, origional = prover(pk, u, e1, alpha, s)
        end_prover = timeit.default_timer()
        prover_time = end_prover - start_prover
        total_prover_time += prover_time

        #sign the proof
        start_signer = timeit.default_timer()
        signature, t_A, t_B = signing.sign(origional, A1, A2, Bext, s1, s2, C, q)
        end_signer = timeit.default_timer()
        print("sig: ", signature)
        sign_time = end_signer - start_signer
        total_sign_time += sign_time

        #verify signature
        start_sig_verify = timeit.default_timer()
        verification_passed = signing.verify_signature(signature, origional, t_A, C, A1, A2, Bext, q)
        end_sig_verify = timeit.default_timer()
        sig_verify_time = end_sig_verify - start_sig_verify
        total_sig_verify_time += sig_verify_time

        #verify proof
        start_verifier = timeit.default_timer()
        verification = verifier(proof, shape)
        end_verifier = timeit.default_timer()
        verifier_time = end_verifier - start_verifier
        total_verifier_time += verifier_time



    # print(f"Verification result for run {_ + 1}: {verification}")
    print("signature_verification: ", verification_passed)
    print("verification: ", verification)
    average_prover_time = total_prover_time / num_runs
    average_verifier_time = total_verifier_time / num_runs
    average_sign_time = total_sign_time / num_runs
    average_sig_verify_time = total_sig_verify_time / num_runs

    print(f"Average proof generation runtime over {num_runs} runs: {average_prover_time} s")
    print(f"Average verification time over {num_runs} runs: {average_verifier_time * 1000} ms")
    print(f"Average proof signing time over {num_runs} runs: {average_sign_time * 1000} ms")
    print(f"Average signature verification time over {num_runs} runs: {average_sig_verify_time * 1000} ms")

    # Calculate the size in bytes of the hexadecimal string
    hex_size_bytes = len(proof) // 2
    print(f"The proof size is: {hex_size_bytes} bytes")

    # serialized_proof = pickle.dumps(proof)
    # # Calculate the size in bytes
    # proof_size_bytes = len(serialized_proof)
    # print(f"The proof size in bytes is: {proof_size_bytes} bytes")


if __name__ == "__main__":
    main()
