import numpy as np

n = 6
max_value = 12


def witnessGenerator():
    # Generate matrix A with values between 0 and 72
    A = np.random.randint(0, max_value + 1, (n, n))

    # Generate vector s with values between 0 and 72
    s = np.random.randint(0, max_value + 1, n)

    # Generate vector e with values between 0 and 72
    e = np.random.randint(0, max_value + 1, n)

    a = np.zeros((n, n))
    c = np.zeros((n, n - 1))
    out = np.zeros(n)

    # Compute a values
    for i in range(n):
        for j in range(n):
            a[i][j] = A[i][j] * s[j]

    # Compute c values and final output values
    for i in range(n):
        c[i][0] = a[i][0] + a[i][1]
        for j in range(1, n - 1):
            c[i][j] = c[i][j - 1] + a[i][j + 1]
        out[i] = c[i][n - 2] + e[i]

    A_flat = A.flatten()
    a_flat = a.flatten()
    c_flat = c.flatten()

    # Initialize s_example with the constant 1
    s_example = [1]

    # Add output values
    s_example.extend(out.astype(int))

    # Add matrix A elements
    s_example.extend(A_flat.astype(int))

    # Add vector s elements
    s_example.extend(s.astype(int))

    # Add vector e elements
    s_example.extend(e.astype(int))

    # Add intermediate a values
    s_example.extend(a_flat.astype(int))

    # Add intermediate c values
    s_example.extend(c_flat.astype(int))

    # # Ensure the length is exactly 1009
    # assert len(s_example) == 1009

    # Display the first 100 values to check
    print(s_example)
    print(len(s_example))
    return s_example