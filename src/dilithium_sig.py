import subprocess
import json

def sign_proof(proof):
    proof_str = json.dumps(proof)  # Convert tuple of arrays to JSON string
    result = subprocess.run(
        ['node', 'dilithium.js', 'signProof', proof_str],
        capture_output=True,
        text=True
    )
    if result.returncode != 0:
        print("STDOUT:", result.stdout)
        print("STDERR:", result.stderr)
        raise Exception(f"Error signing proof: {result.stderr}")
    signed_proof = json.loads(result.stdout)
    return signed_proof

def verify_signature(public_key, proof, signature):
    proof_str = json.dumps(proof)  # Convert tuple of arrays to JSON string
    result = subprocess.run(
        ['node', 'dilithium.js', 'verifySignature', public_key, proof_str, signature],
        capture_output=True,
        text=True
    )
    if result.returncode != 0:
        print("STDOUT:", result.stdout)
        print("STDERR:", result.stderr)
        raise Exception(f"Error verifying signature: {result.stderr}")
    is_valid = result.stdout.strip() == "true"
    return is_valid

if __name__ == "__main__":
    # Example proof tuple of arrays
    proof = ([1, 2, 3], [4, 5, 6])

    # Signing the proof
    signed_proof = sign_proof(proof)
    print("Signed Proof:", signed_proof)

    # Extracting the required values from the signed proof
    public_key = signed_proof["publicKey"]
    signature = signed_proof["signature"]

    # Verifying the signature
    is_valid = verify_signature(public_key, proof, signature)
    print("Signature is valid:", is_valid)
