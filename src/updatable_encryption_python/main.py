#set PATH=d:\python311\;%PATH%;
#set PATH=c:\python311\;%PATH%;

#import numpy as np
#from scipy.linalg import qr
#from sympy import Matrix
#from scipy.stats import norm
#from sampledfunc import *

#G = generate_G(n, k, q)

from key_generation import *
from parameters import *
from encrypt_message import *
from utility_functions import Lemma1Check, EqualityCheck
from decrypt_message import *
from token_generation import *


def encode_nothing():
    return np.random.randint(0, 2, nk)















def main():

    #global e0_hidden, e1_hidden, e2_hidden

    msg1 = encode_nothing()
    pk1, sk1 = KeyGen(None)
    H_mu, b = Encrypt(pk1, encode(msg1), H_hardcoded, None, False)
    e0, e1, e2 = getError()
    Lemma1Check(sk1, np.concatenate((e0, e1)))
    EqualityCheck(pk1, sk1)
    msg2 = Decrypt(sk1, H_mu, b, pk1, True)
    assert np.array_equal(msg1, msg2)


    current_b = b
    current_pk = pk1
    current_Hmu = H_mu

    for counter1 in range(1, 11):
        print("-------------------- UPDATE PHASE", counter1, "--------------------")
        new_pk, _ = KeyGen(sk1)
        token = TokenGen(current_pk, sk1, new_pk, current_Hmu)
        b_prime, H_mu_prime = update(token, (current_Hmu, current_b))
        new_msg = Decrypt(sk1, H_mu_prime, b_prime, new_pk, False)
        assert np.array_equal(msg1, new_msg)

        current_b = b_prime
        current_pk = new_pk
        current_Hmu = H_mu_prime







    '''
    print("-------------------- UPDATE PHASE --------------------")
    b_curr = b
    pk_curr = pk1
    pk2, sk2 = KeyGen(sk1)
    token = TokenGen(pk_curr, sk1, pk2, H_mu)
    b_prime, H_mu_prime = update(token, (H_mu, b_curr))
    #b_prime = b_prime % q
    print("H_mu_prime", H_mu_prime.shape)
    msg3 = Decrypt(sk1, H_mu_prime, b_prime, pk2, False)
    assert np.array_equal(msg1, msg3)


    print("-------------------- UPDATE PHASE 2 --------------------")
    b_curr = b_prime
    pk_curr = pk2
    pk2, sk2 = KeyGen(sk1)
    token = TokenGen(pk_curr, sk1, pk2, H_mu_prime)
    b_prime, H_mu_prime = update(token, (H_mu_prime, b_curr))
    print("H_mu_prime", H_mu_prime.shape)
    msg4 = Decrypt(sk1, H_mu_prime, b_prime, pk2, False)
    assert np.array_equal(msg1, msg4)
    '''

    print("Success: msg1 == msg2 == msg3 == msg4")
    
if __name__=="__main__":
    main()


#to check
#EqualityCheck