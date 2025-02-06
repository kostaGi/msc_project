#local files
from dilithium.dilithium_constants import * 

#pack and unpack different parameters into bytes

# utility method
def sub_mod_q(x, y):
    return (x - y) % q

# generic pack function used by each pack function
def bit_pack(coeffs, n_bits, n_bytes):
    r = 0
    for c in reversed(coeffs):
        r <<= n_bits
        r |= int(c)
    return r.to_bytes(n_bytes, "little")

# generic unpack function used by each unpack function
def bit_unpack(input_bytes, n_bits):
    if (len(input_bytes) * n_bits) % 8 != 0:
        raise ValueError(
            "Input bytes do not have a length compatible with the bit length"
        )
    r = int.from_bytes(input_bytes, "little")
    mask = (1 << n_bits) - 1
    return [(r >> n_bits * i) & mask for i in range(n)]

# pack public key from rho and t1
def pack_pk(rho, t1):
    return rho + bit_pack_t1_vector(t1)   

# unpack public key into rho and t1
def unpack_pk(pk_bytes):
    rho, t1_bytes = pk_bytes[:32], pk_bytes[32:]
    t1 = bit_unpack_t1_vector(t1_bytes, k, 1)
    return rho, t1

# pack secret key from rho, K, tr, s1, s2, t0
def pack_sk(rho, K, tr, s1, s2, t0):
    s1_bytes = bit_pack_s_vector(s1)
    s2_bytes = bit_pack_s_vector(s2)
    t0_bytes = bit_pack_t0_vector(t0)
    return rho + K + tr + s1_bytes + s2_bytes + t0_bytes

# unpack secret key into rho, K, tr, s1, s2, t0
def unpack_sk(sk_bytes):
    if eta == 2:
        s_bytes = 96
    else:
        s_bytes = 128
    s1_len = s_bytes * l
    s2_len = s_bytes * k
    t0_len = 416 * k
    if len(sk_bytes) != 2 * 32 + 48 + s1_len + s2_len + t0_len:
        print("My length: ", len(sk_bytes))
        print("Expected length: ", 3 * 32 + s1_len + s2_len + t0_len)
        raise ValueError("SK packed bytes is of the wrong length")
    # Split bytes between seeds and vectors
    sk_seed_bytes, sk_vec_bytes = sk_bytes[:112], sk_bytes[112:]
    # Unpack seed bytes
    rho, K, tr = (
        sk_seed_bytes[:32],
        sk_seed_bytes[32:64],
        sk_seed_bytes[64:112],
    )
    # Unpack vector bytes
    s1_bytes = sk_vec_bytes[:s1_len]
    s2_bytes = sk_vec_bytes[s1_len : s1_len + s2_len]
    t0_bytes = sk_vec_bytes[-t0_len:]
    # Unpack bytes to vectors
    s1 = bit_unpack_s_vector(s1_bytes, l, 1, eta)
    s2 = bit_unpack_s_vector(s2_bytes, k, 1, eta)
    t0 = bit_unpack_t0_vector(t0_bytes, k, 1)
    return rho, K, tr, s1, s2, t0

# pack t1 into bytes
def bit_pack_t1_vector(vector):
    return b"".join(bit_pack_t1(poly) for poly in vector)

# pack each poly of t1 into bytes 
def bit_pack_t1(coeffs):
    # 320 = 256 * 10 // 8
    return bit_pack(coeffs, 10, 320)

# unpack t1
def bit_unpack_t1_vector(input_bytes, m, n):
    packed_len = 320
    poly_bytes = [
        input_bytes[i : i + packed_len]
        for i in range(0, len(input_bytes), packed_len)
    ]
    matrix = [tuple([bit_unpack_t1(poly_bytes[n * i + j]) for j in range(n)]) for i in range(m)]
    matrix = tuple(tuple([x[0] for x in matrix]))
    return matrix

# unpack each poly of t1
def bit_unpack_t1(input_bytes):
    coefficients = bit_unpack(input_bytes, 10)
    return coefficients

# pack s into bytes
def bit_pack_s_vector(vector):
    return b"".join(bit_pack_s(poly) for poly in vector)

# pack each poly of s into bytes   
def bit_pack_s(coeffs):
    altered_coeffs = [sub_mod_q(eta, c) for c in coeffs]
    # Level 2 and 5 parameter set
    if eta == 2:
        return bit_pack(altered_coeffs, 3, 96)
    # Level 3 parameter set
    assert eta == 4, f"Expected eta to be either 2 or 4, got {eta = }"
    return bit_pack(altered_coeffs, 4, 128)  

# unpack s
def bit_unpack_s_vector(input_bytes, m, n, eta):
    # Level 2 and 5 parameter set
    if eta == 2:
        packed_len = 96
    # Level 3 parameter set
    elif eta == 4:
        packed_len = 128
    else:
        raise ValueError("Expected eta to be either 2 or 4")
    poly_bytes = [
        input_bytes[i : i + packed_len]
        for i in range(0, len(input_bytes), packed_len)
    ]
    matrix = [
        [bit_unpack_s(poly_bytes[n * i + j]) for j in range(n)] for i in range(m)
    ]
    matrix = (tuple([x[0] for x in matrix]))
    return matrix

# unpack each poly of s
def bit_unpack_s(input_bytes):
    # Level 2 and 5 parameter set
    if eta == 2:
        altered_coeffs = bit_unpack(input_bytes, 3)
    # Level 3 parameter set
    else:
        assert eta == 4, f"Expected eta to be either 2 or 4, got {eta = }"
        altered_coeffs = bit_unpack(input_bytes, 4)
    coefficients = [eta - c for c in altered_coeffs]
    return coefficients

# pack t0 into bytes
def bit_pack_t0_vector(vector):
    return b"".join(bit_pack_t0(poly) for poly in vector)

# pack each poly of t0 into bytes
def bit_pack_t0(coeffs):
    # 416 = 256 * 13 // 8
    altered_coeffs = [(1 << 12) - c for c in coeffs]
    return bit_pack(altered_coeffs, 13, 416)

# unpack t0
def bit_unpack_t0_vector(input_bytes, m, n):
    packed_len = 416
    poly_bytes = [
        input_bytes[i : i + packed_len]
        for i in range(0, len(input_bytes), packed_len)
    ]
    matrix = [
        [bit_unpack_t0(poly_bytes[n * i + j]) for j in range(n)] for i in range(m)
    ]
    matrix = (tuple([x[0] for x in matrix]))
    return matrix

# unpack each poly of t0
def bit_unpack_t0(input_bytes):
    altered_coeffs = bit_unpack(input_bytes, 13)
    coefficients = [(1 << 12) - c for c in altered_coeffs]
    return coefficients

# pack sigma from c_tilde, z , h
def pack_sig(c_tilde, z, h, gamma_1):
    return c_tilde + bit_pack_z_vector(z, gamma_1) + pack_h(h)

# unpack sigma into c_tilde, z , h
def unpack_sig(sig_bytes, gamma_1):
    c_tilde = sig_bytes[:32]
    z_bytes = sig_bytes[32 : -(k + omega)]
    h_bytes = sig_bytes[-(k + omega) :]
    z = bit_unpack_z_vector(z_bytes, l, 1, gamma_1)
    h = unpack_h(h_bytes)
    return c_tilde, z, tuple(h)

# pack z into bytes
def bit_pack_z_vector(vector, gamma_1):
    return b"".join(bit_pack_z(poly, gamma_1) for poly in vector)

# pack each poly of z into bytes
def bit_pack_z(coeffs, gamma_1):
    altered_coeffs = [sub_mod_q(gamma_1, c) for c in coeffs]
    # Level 2 parameter set
    if gamma_1 == (1 << 17):
        return bit_pack(altered_coeffs, 18, 576)
    # Level 3 and 5 parameter set
    assert gamma_1 == (
        1 << 19
    ), f"Expected gamma_1 to be either 2^17 or 2^19, got: {gamma_1 = }"
    return bit_pack(altered_coeffs, 20, 640)

# unpack z
def bit_unpack_z_vector(input_bytes, m, n, gamma_1):
    # Level 2 parameter set
    if gamma_1 == (1 << 17):
        packed_len = 576
    # Level 3 and 5 parameter set
    elif gamma_1 == (1 << 19):
        packed_len = 640
    else:
        raise ValueError("Expected gamma_1 to be either 2^17 or 2^19")
    poly_bytes = [
        input_bytes[i : i + packed_len]
        for i in range(0, len(input_bytes), packed_len)
    ]
    matrix = [
        [bit_unpack_z(poly_bytes[n * i + j], gamma_1) for j in range(n)] for i in range(m)
    ]
    matrix = [x[0] for x in matrix]
    return matrix

# unpack each poly of z 
def bit_unpack_z(input_bytes, gamma_1):
    # Level 2 parameter set
    if gamma_1 == (1 << 17):
        altered_coeffs = bit_unpack(input_bytes, 18)
    # Level 3 and 5 parameter set
    else:
        assert gamma_1 == (
            1 << 19
        ), f"Expected gamma_1 to be either 2^17 or 2^19, got {gamma_1 = }"
        altered_coeffs = bit_unpack(input_bytes, 20)
    coefficients = [gamma_1 - c for c in altered_coeffs]
    return coefficients

# pack h into bytes
def pack_h(h):
    non_zero_positions = [
        [i for i, c in enumerate(poly) if c == 1]
        for poly in h
    ]
    packed = []
    offsets = []
    for positions in non_zero_positions:
        packed.extend(positions)
        offsets.append(len(packed))
    padding_len = omega - offsets[-1]
    packed.extend([0 for _ in range(padding_len)])
    return bytes(packed + offsets)    

# unpack h
def unpack_h(h_bytes):
    offsets = [0] + list(h_bytes[-k :])
    non_zero_positions = [
        list(h_bytes[offsets[i] : offsets[i + 1]]) for i in range(k)
    ]
    matrix = []
    for poly_non_zero in non_zero_positions:
        coeffs = [0 for _ in range(256)]
        for non_zero in poly_non_zero:
            coeffs[non_zero] = 1
        matrix.append(coeffs)
    return matrix

# pack high parts of w1 into bytes
def bit_pack_w_vector(vector):
    return b"".join(bit_pack_w(poly) for poly in vector)

# pack each poly of w1
def bit_pack_w(w):
    # Level 2 parameter set
    if gamma_2 == 95232:
        return bit_pack(w, 6, 192)
    # Level 3 and 5 parameter set
    assert (
        gamma_2 == 261888
    ), f"Expected gamma_2 to be either (q-1)/88 or (q-1)/32, got {gamma_2 = }"
    return bit_pack(w, 4, 128)