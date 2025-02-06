import requests
import logging
import json
import os
import numpy as np
import fhelweSNARK
import asyncio
import oqs
import base64
#import dilithium_sig as dilsig
from dilithium.generate_key import generate_keys
from dilithium.sign_message import dilithium_sign

# Configure basic logger
logging.basicConfig(level=logging.INFO)
order = 73
q = 8929
t = 73
d = 57
delta = q // t

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
    off_chain_service_url = "http://localhost:8000/compute-proof"
    
    # Assuming fhelweSNARK.setup() returns these values correctly
    alpha, sk, a2, e, pk, u, e1, e2 = fhelweSNARK.setup()

    #secret that prover knows
    w = [1, 5, 8, 1, 8, 4, 7, 1, 0, 6, 4, 9, 1, 0, 8, 0, 1, 5, 0, 3, 2, 1, 0, 0, 1, 1, 1, 0, 1, 4, 0, 0, 0, 6, 0, 0, 1, 0, 0, 0, 1, 5, 0, 0, 2, 4, 4, 4, 6, 6, 7, 0, 0, 1, 5, 5, 7]

    # prove that secret is known by prover 
    proof_hex, proof_shape, proof = fhelweSNARK.prover(pk, u, e1, alpha, w)

    proof_response = [p.tolist() for p in proof]
    print("proof_response: ", proof_response)

    url = "http://localhost:2800/checkProof"

    
    # generate dilithium keys (byte format)
    d_pk, d_sk = generate_keys()

    # transform proof to bytes
    proof_bytes = bytes(proof_hex, 'utf-8')

    # sign dilithium message - proof_bytes
    d_sigma = dilithium_sign(proof_bytes, d_sk)
        
    payload_dict = {
        "proof": proof_bytes.hex(),
        "public_key": d_pk.hex(),
        "signature": d_sigma.hex()
    }	
    payload = json.dumps(payload_dict)
    headers = {
        'Content-Type': 'application/json'
    }
    response = requests.post(url, headers=headers, data=payload)
    logging.info(f"Check Proof response status: {response.status_code}")
    logging.info(f"Check Proof response content: {response.content}")

if __name__ == "__main__":
    asyncio.run(main())
