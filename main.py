from generate_key import generate_keys 
from sign_message import dilithium_sign
from verify_signature import verify

from dilithium_constants import * 



# tested on python 3.11
def main():

    print("Security parameters for Dilithium:")
    print("q: ", q)
    print("n: ", n)
    print("k: ", k)
    print("l: ", l)
    print("eta: ", eta)
    print("d: ", d)
    print("tau: ", tau)
    print("gamma1: ", gamma1)
    print("gamma2: ", gamma2)
    print("beta: ", beta)
    print("omega: ", omega)
    print("\n\n")

    # generate keys (byte format)
    pk, sk = generate_keys()

    msg1 = "It_is_working"
    msg2 = "It_is_workinG"

    msg1_bytes = bytes(msg1, 'utf-8')
    msg2_bytes = bytes(msg2, 'utf-8')

    # Sign msg1
    sigma = dilithium_sign(msg1_bytes, sk)

    print("Test Dilithium verify:")
    print("\n")
    
    print("Assert that:", msg1, "is verified as:", msg1, ":", verify(pk, msg1_bytes, sigma))
    print("\n\n")
    print("Assert that:", msg2, "is verified as:", msg1, ":", verify(pk, msg2_bytes, sigma))

if __name__ == "__main__":
    main()   