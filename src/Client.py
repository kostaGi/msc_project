import requests
import logging
import json
import os
import numpy as np
import fhelweSNARK
import asyncio
import oqs
import base64
import time
#import dilithium_sig as dilsig
from dilithium.generate_key import generate_keys
from dilithium.sign_message import dilithium_sign
from dilithium.verify_signature import verify
from updatable_encryption_python.key_generation import KeyGen
from updatable_encryption_python.utility_functions import encode
from updatable_encryption_python.encrypt_message import Encrypt
from updatable_encryption_python.token_generation import TokenGen

# Configure basic logger
logging.basicConfig(level=logging.INFO)
order = 73
q = 8929
t = 73
d = 57
delta = q // t

nk = 320

# Polynomial modulus
p_q = np.poly1d([1] + ([0] * (d - 1)) + [1])

def send_witness_data(url, data):
    """
    Send witness data to the specified off-chain computation service and return the proof.
    """
    try:
        response = requests.post(url, json=data)
        response.raise_for_status()  # Will raise an HTTPError for bad responses (400 or 500 level responses)
        return response.json()
    except requests.exceptions.HTTPError as http_err:
        logging.error(f"HTTP error occurred: {http_err} - Status code: {response.status_code} - Response content: {response.content}")
    except Exception as err:
        logging.error(f"An error occurred: {err}")
    return None

async def main():

    # Default params

    off_chain_service_url = "http://localhost:8000/compute-proof"
    url = "https://localhost:2800"
    middleware_certificate = "middleware_certificate/cert.pem"
    peer_name = "Gosho"
    
    # Step 1 Generate secret
    
    # Assuming fhelweSNARK.setup() returns these values correctly
    alpha, sk, a2, e, pk, u, e1, e2 = fhelweSNARK.setup()

    #secret that prover knows
    w = [1, 5, 8, 1, 8, 4, 7, 1, 0, 6, 4, 9, 1, 0, 8, 0, 1, 5, 0, 3, 2, 1, 0, 0, 1, 1, 1, 0, 1, 4, 0, 0, 0, 6, 0, 0, 1, 0, 0, 0, 1, 5, 0, 0, 2, 4, 4, 4, 6, 6, 7, 0, 0, 1, 5, 5, 7]

    # prove that secret is known by prover 
    proof_hex, proof_shape, proof = fhelweSNARK.prover(pk, u, e1, alpha, w)

    # flatten to array of 2 equal arrays => [[a1,b2,c3 ..], [a1, b2, c3 ...]]
    proof_response = [p.tolist() for p in proof]

    # Step 2 Generate Dilithium and Updatable encryption keys
    
    # generate dilithium keys (byte format)
    d_pk, d_sk = generate_keys()

    # generate UE keys matrix format
    ue_pk1, ue_sk1 = KeyGen(None)
    A0, A1, A2  = ue_pk1
    ue_pk_send = np.block([A0, A1, A2])

    # Step 3 Store the public keys on the hyperledger (D and UE)
    payload_dict = {
        "request_type": "1",
        "owner": peer_name,
        "public_key_D":  d_pk.hex(),
        "public_key_UE": ue_pk_send.tolist()
    }	
    payload = json.dumps(payload_dict)
    response = requests.post(url, data=payload, verify = 'middleware_certificate/cert.pem')
    #logging.info(f"Check Proof response status: {response.status_code}")
    #logging.info(f"Check Proof response content: {response.content}")
    #logging.info(f"Check Proof response content: {response}")
    stored_keys_id = json.loads(str(response.content, 'utf-8'))["id"]

    
    time.sleep(5)

    # Step 4 Read Dilithium public key from ledger and compare
    
    payload_dict = {
        "request_type": "2",
        "owner": peer_name,
        "stored_keys_id":  str(stored_keys_id)
    }	
    payload = json.dumps(payload_dict)
    response = requests.post(url, data=payload, verify = 'middleware_certificate/cert.pem')
    #logging.info(f"Check Proof response status: {response.status_code}")
    #logging.info(f"Check Proof response content: {response.content}")
   
    stored_dilithium_value = json.loads(str(response.content, 'utf-8'))["public_key_UE"]

    # assert that they are equal
    assert ue_pk_send.tolist() == stored_dilithium_value

    time.sleep(5)

    # Step 5 Read Updatable encryption public key from ledger and compare
    
    payload_dict = {
        "request_type": "3",
        "owner": peer_name,
        "stored_keys_id":  str(stored_keys_id)
    }	
    payload = json.dumps(payload_dict)
    response = requests.post(url, data=payload, verify = 'middleware_certificate/cert.pem')
    #logging.info(f"Check Proof response status: {response.status_code}")
    #logging.info(f"Check Proof response content: {response.content}")
   
    stored_updatable_encryption_value = json.loads(str(response.content, 'utf-8'))["public_key_D"]

    # assert that they are equal
    assert d_pk.hex() == stored_updatable_encryption_value


    # Step 6 Encode the secret into bits and encrypt with updatable encryption

    # array to bits
    proof_response_bits = array_to_bits(proof_response)
    # size of bits
    proof_response_bits_size = len(proof_response_bits)

    # pad a message to be size which can be UE
    padded_message = proof_response_bits
    # pad with random bits
    if proof_response_bits_size % nk != 0:
        ending = encode_random(nk - proof_response_bits_size % nk)
        padded_message = np.concatenate((padded_message, ending), axis=0)

    # encode the bits (UE method) 
    encoded_msg = encode(padded_message)

    # encrypt the encoded bits (UE method)
    H_mu, b = Encrypt(ue_pk1, encoded_msg, None, None, False, proof_response_bits_size) 

    # Step 7 Apply dilithium sign to the encryption

    # transform the components to string and concatenate them
    dilithium_input = convert_2d_int_array_to_string(H_mu.astype(int).tolist()) + convert_1d_int_array_to_string(b.tolist()) + peer_name + str(proof_response_bits_size)
    # transform the string to bytes utf-8 and sign 
    dilithium_signature = dilithium_sign(bytes(dilithium_input, 'utf-8'), d_sk)
    #print("\n dilithium_input: ", dilithium_input, "\n\n\n")

    # assert that it was signed Client side
    assert verify(d_pk, bytes(dilithium_input, 'utf-8'), dilithium_signature)

    time.sleep(5)
    
    # Step 8 Send the encryption and signature to the hyperledge farbric to be verified and stored

    #Get request
    payload_dict = {
        "request_type": "4",
        "stored_keys_id":  str(stored_keys_id),
        "H_mu": H_mu.astype(int).tolist(),
        "ciphertext": b.astype(int).tolist(),
        "owner": peer_name,
        "size": str(proof_response_bits_size),
        "signature": dilithium_signature.hex()
    }	
    payload = json.dumps(payload_dict)
    #print("\n", payload)
    #print("\n\n\n", dilithium_signature)
    #print("\n", b.tolist())
    #print("\n", str(proof_response_bits_size))
    #print("\n", (dilithium_signature))
    response = requests.post(url, data=payload, verify = 'middleware_certificate/cert.pem')
    #logging.info(f"Check Proof response status: {response.status_code}")
    #logging.info(f"Check Proof response content: {response.content}")

    stored_secret_id = json.loads(str(response.content, 'utf-8'))["id"]

    time.sleep(5)
    
    # Step 9 Read Stored encryption values and assert equality

    payload_dict = {
        "request_type": "5",
        "owner": peer_name,
        "id":  str(stored_secret_id) 
    }	
    payload = json.dumps(payload_dict)
    response = requests.post(url, data=payload, verify = 'middleware_certificate/cert.pem')
    #logging.info(f"Check Proof response status: {response.status_code}")
    #logging.info(f"Check Proof response content: {response.content}")
    Hmu_check = json.loads(response.content)["Hmu"]
    size_check = json.loads(response.content)["Size"]
    ciphertext_check = json.loads(response.content)["Ciphertext"]

    #print("Hmu_check: ", Hmu_check)
    #print("H_mu.astype(int).tolist(): ", H_mu.astype(int).tolist())

    assert Hmu_check == H_mu.astype(int).tolist()
    assert size_check == proof_response_bits_size
    assert ciphertext_check == b.astype(int).tolist() 


    # Step 10 Generate an updatable encryption token

    ue_pk2, ue_sk2 = KeyGen(ue_sk1)
    A0_prime, A1_prime, A2_prime  = ue_pk2
    ue_pk_send2 = np.block([A0_prime, A1_prime, A2_prime])    
    X0, X1, X2, b_zero_message, Hmu_prime = TokenGen(ue_pk1, ue_sk1, ue_pk2, H_mu, proof_response_bits_size)


    # Step 11 Apply dilithium sign to the encryption

    dilithium_input = convert_2d_int_array_to_string(X0.tolist())  + convert_2d_int_array_to_string(X1.tolist()) + convert_2d_int_array_to_string(X2.tolist()) + \
    convert_2d_int_array_to_string(Hmu_prime.astype(int).tolist()) + convert_1d_int_array_to_string(b_zero_message.tolist()) + peer_name
    
    dilithium_signature = dilithium_sign(bytes(dilithium_input, 'utf-8'), d_sk)

    # assert that it was signed Client side
    assert verify(d_pk, bytes(dilithium_input, 'utf-8'), dilithium_signature)

    payload_dict = {
        "request_type": "6",
        "stored_keys_id" : str(stored_keys_id),
        "stored_secret_id":  str(stored_secret_id), 
        "X0": X0.tolist(),
        "X1": X1.tolist(),
        "X2": X2.tolist(),
        "Hmu_prime": Hmu_prime.astype(int).tolist(),
        "b_zero_message": b_zero_message.tolist(),
        "owner": peer_name,
        "signature": dilithium_signature.hex(),
    }	

    payload = json.dumps(payload_dict)

    response = requests.post(url, data=payload, verify = 'middleware_certificate/cert.pem')
    #logging.info(f"Check Proof response status: {response.status_code}")
    #logging.info(f"Check Proof response content: {response.content}")

    isUpdated_check = json.loads(response.content)["result"]

    assert isUpdated_check

    # Step 12 Apply dilithium sign to the encryption
    
    


def convert_1d_int_array_to_string(arr):
    try:
        return json.dumps(arr).replace(' ', '')  # Convert list to JSON string
    except Exception as e:
        return str(e)  # Handle errors

def convert_2d_int_array_to_string(arr):
    try:
        return json.dumps(arr).replace(' ', '')  # Convert 2D list to JSON string
    except Exception as e:
        return str(e)  # Handle errors

def array_to_bits(arr):
    arr = np.array(arr, dtype=np.uint8) 
    bit_strings = [format(num, '08b') for row in arr for num in row]
    bit_list = [int(bit) for bit_string in bit_strings for bit in bit_string]
    return np.array(bit_list, dtype=np.uint8)



def encode_random(size):
    return np.random.randint(0, 2, size)

if __name__ == "__main__":
    asyncio.run(main())
