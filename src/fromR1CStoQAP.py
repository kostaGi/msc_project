import galois
import numpy as np
import fromLWEtoR1CS as r1



print("Initializing a large field...")
order = 701
GF = galois.GF(order)
q = 2**8
print("q: ", q)
t = order
d = 64
delta = q // t
# Polynomial modulus
p_q = np.poly1d([1] + ([0] * (d - 1)) + [1])
print("Field initialized")
# s = GF(np.array([1, 22, 38, 21, 14, 37, 19, 1, 0, 1, 1, 21, 0, 37, 0, 21, 37]))
# s = GF(np.array([1, 2, 4, 1, 2, 3, 2, 1, 0, 1, 1, 1, 0, 3, 0, 1, 3]))
# s = GF(np.array([1, 6, 2, 1, 2, 3, 7, 0, 0, 1, 1, 2, 0, 3, 0, 1, 8]))


def interpolate_column(col, nb):
    #print("nb: ", nb)
    xs = GF(np.arange(1, nb+1))
    # print(col)
    return galois.lagrange_poly(xs, col)

def get_polys_of_matrix(matrix):
    polys = []
    nb_of_rows = len(matrix)
    nb_of_columns = len(matrix[0])
    print("nb_of_rows: ", nb_of_rows)
    print("nb_of_columns: ", nb_of_columns)
    # for each column
    for col_id in range(nb_of_columns):
        #print("col_id: ", col_id)
        column = []
        for row in range(nb_of_rows):
            column.append(matrix[row][col_id])
        polys.append(interpolate_column(GF(np.array(column)), nb_of_rows))
    return np.array(polys)

def polySum(s):
    A, B, C = r1.LWEToR1CS_transform()
    # witness = np.array([1, 2, 4, 1, 2, 3, 2, 1, 0, 1, 1, 1, 0, 3, 0, 1, 3])

    ## commented out 
    #print(np.matmul(C, np.array(s)) == (np.matmul(A, np.array(s)) * np.matmul(B, np.array(s))))

    # print("Extra A check: ", A)
    ## computes all interpolated polynomials for L, R, and O
    U_polys = get_polys_of_matrix(A)
    V_polys = get_polys_of_matrix(B)
    W_polys = get_polys_of_matrix(C)

    # Summing all the polynamials of the collection into one polynomial
    U = galois.Poly([0], field=GF)
    V = galois.Poly([0], field=GF)
    W = galois.Poly([0], field=GF)
    # print("U first check: ", U_polys)
    # print("V first check: ", V_polys)
    # print("W first check: ", W_polys)
    for i in range(len(U_polys)):
        U += U_polys[i]
        V += V_polys[i]
        W += W_polys[i]
    # A_polys = get_polys_of_S_matrix(s)

    Ua = galois.Poly([0], field=GF)
    for i in range(len(s)):
        Ua += U_polys[i] * s[i]

    Va = galois.Poly([0], field=GF)
    for i in range(len(s)):
        Va += V_polys[i] * s[i]

    Wa = galois.Poly([0], field=GF)
    for i in range(len(s)):
        Wa += W_polys[i] * s[i]

    # print("U: ", U)
    # print("V: ", V)
    # print("W: ", W)
    #
    # print("Ua: ", Ua)
    # print("Va: ", Va)
    # print("Wa: ", Wa)
    #
    # print("U degree: ", U.degree)
    # print("V degree: ", V.degree)
    # print("W degree: ", W.degree)
    #
    # print("Ua degree: ", Ua.degree)
    # print("Va degree: ", Va.degree)
    # print("Wa degree: ", Wa.degree)

    return U, V, W, Ua, Va, Wa


def get_polys_of_S_matrix(matrix):
    polys = []
    nb_of_columns = len(matrix)
    # for each column
    for col_id in range(nb_of_columns):
        polys.append(galois.Poly([matrix[0]], field=GF))
    return np.array(polys)

