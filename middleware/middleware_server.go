package main

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"net/http"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	//Org id
	mspID = "Org1MSP"
	// Path to crypto materials.
	cryptoPath = "/home/kosix/Desktop/hyperledge_fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com"
	// Path to user certificate directory.
	certPath = cryptoPath + "/users/User1@org1.example.com/msp/signcerts"
	// Path to user private key directory.
	keyPath = cryptoPath + "/users/User1@org1.example.com/msp/keystore"
	// Path to peer tls certificate.
	tlsCertPath = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	// Gateway peer endpoint.
	peerEndpoint = "localhost:7051"
	// Gateway peer SSL host name override.
	gatewayPeer = "peer0.org1.example.com"
)

type ResponseWithId struct {
	Id string `json:"id"`
}

type ResponseWithBool struct {
	Bool bool `json:"result"`
}

type ResponseWithPublicUE struct {
	Public_key_UE [][]int `json:"public_key_UE"`
}

type ResponseWithPublicD struct {
	Public_key_D string `json:"public_key_D"`
}

// secret
type ParseData struct {
	ID         string  `json:"ID"`
	Hmu        [][]int `json:"Hmu"`
	Ciphertext []int   `json:"Ciphertext"`
	Owner      string  `json:"Owner"`
	Size       int     `json:"Size"`
}

//var now = time.Now()
//var assetId = fmt.Sprintf("asset%d", now.Unix()*1e3+int64(now.Nanosecond())/1e6)

// handle request
func handler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	// Override default values for chaincode and channel name as they may differ in testing contexts.
	chaincodeName := "basic"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}
	channelName := "mychannel"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	//Parse Json from Client
	var data map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		response := ResponseWithId{Id: "-1"}
		json.NewEncoder(w).Encode(response)
		fmt.Println("\n", err)
		return
	}

	request_type, isString := data["request_type"].(string)

	if !isString {
		response := ResponseWithId{Id: "-1"}
		json.NewEncoder(w).Encode(response)
		fmt.Println("\n request_type is not a string")
		return
	}

	//case 1 store peer's 2 public keys
	//case 2 get peer's public key UE (updatable encryption)
	//case 3 get peer's public key D (dilithium)
	//case 4 store secret
	//case 5 update secret
	//case 6 get secret
	switch request_type {
	case "1":
		owner, isString1 := data["owner"].(string)
		public_key_D, isString2 := data["public_key_D"].(string)

		if !isString1 || !isString2 {
			response := ResponseWithId{Id: "-1"}
			json.NewEncoder(w).Encode(response)
			fmt.Println("\n Case 1: At least 1 value is not strings")
			return
		}

		public_key_UE_string, err := parseJsonTo2DString(data["public_key_UE"])

		if err != nil {
			response := ResponseWithId{Id: "-1"}
			json.NewEncoder(w).Encode(response)
			fmt.Println("\n", err)
			return
		}

		return_id, err := storeOwnerData(contract, owner, public_key_D, public_key_UE_string)

		if err != nil {
			response := ResponseWithId{Id: "-1"}
			json.NewEncoder(w).Encode(response)
			fmt.Println("\n", err)
			return
		}

		response := ResponseWithId{Id: string(return_id)}

		json.NewEncoder(w).Encode(response)
		return

	case "2":
		owner, isString1 := data["owner"].(string)
		stored_keys_id, isString2 := data["stored_keys_id"].(string)

		if !isString1 || !isString2 {
			response := ResponseWithId{Id: "-1"}
			json.NewEncoder(w).Encode(response)
			fmt.Println("\n Case 2: At least 1 value is not strings")
			return
		}

		return_public_key_UE, err := readOwnerPublicUE(contract, stored_keys_id, owner)

		if err != nil {
			response := ResponseWithId{Id: "-1"}
			json.NewEncoder(w).Encode(response)
			fmt.Println("\n", err)
			return
		}

		response := ResponseWithPublicUE{Public_key_UE: (return_public_key_UE)}
		json.NewEncoder(w).Encode(response)
		return

	case "3":
		owner, isString1 := data["owner"].(string)
		stored_keys_id, isString2 := data["stored_keys_id"].(string)

		if !isString1 || !isString2 {
			response := ResponseWithId{Id: "-1"}
			json.NewEncoder(w).Encode(response)
			fmt.Println("\n Case 3: At least 1 value is not strings")
			return
		}

		return_public_key_D, err := readOwnerPublicD(contract, stored_keys_id, owner)

		if err != nil {
			response := ResponseWithId{Id: "-1"}
			json.NewEncoder(w).Encode(response)
			fmt.Println("\n", err)
			return
		}

		response := ResponseWithPublicD{Public_key_D: (return_public_key_D)}
		json.NewEncoder(w).Encode(response)
		return
	case "4":
		owner_id, isString1 := data["stored_keys_id"].(string)
		owner, isString2 := data["owner"].(string)
		signature, isString3 := data["signature"].(string)
		size, isString4 := data["size"].(string)

		if !isString1 || !isString2 || !isString3 || !isString4 {
			response := ResponseWithId{Id: "-1"}
			json.NewEncoder(w).Encode(response)
			fmt.Println("\n Case 4: At least 1 value is not strings")
			return
		}

		H_mu_string, err := parseJsonTo2DString(data["H_mu"])

		if err != nil {
			response := ResponseWithId{Id: "-1"}
			json.NewEncoder(w).Encode(response)
			fmt.Println("\n", err)
			return
		}

		ciphertext_string, err := parseJsonTo1DString(data["ciphertext"])

		if err != nil {
			response := ResponseWithId{Id: "-1"}
			json.NewEncoder(w).Encode(response)
			fmt.Println("\n", err)
			return
		}

		return_id, err := verifyAndStoreSecret(contract, owner_id, H_mu_string, ciphertext_string, owner, size, signature)

		if err != nil {
			response := ResponseWithId{Id: "-1"}
			json.NewEncoder(w).Encode(response)
			fmt.Println("\n", err)
			return
		}

		response := ResponseWithId{Id: string(return_id)}
		json.NewEncoder(w).Encode(response)
		return

	case "5":
		owner, isString1 := data["owner"].(string)
		id, isString2 := data["id"].(string)

		if !isString1 || !isString2 {
			response := ResponseWithId{Id: "-1"}
			json.NewEncoder(w).Encode(response)
			fmt.Println("\n Case 5: At least 1 value is not strings")
			return
		}

		return_secret := readSecretData(contract, id, owner)
		var parsedData ParseData
		err := json.Unmarshal([]byte(return_secret), &parsedData)

		if err != nil {
			response := ResponseWithId{Id: "-1"}
			json.NewEncoder(w).Encode(response)
			fmt.Println("\n", err)
			return
		}

		json.NewEncoder(w).Encode(parsedData)
		return

	case "6":
		stored_keys_id, isString1 := data["stored_keys_id"].(string)
		owner, isString2 := data["owner"].(string)
		signature, isString3 := data["signature"].(string)
		stored_secret_id, isString4 := data["stored_secret_id"].(string)

		if !isString1 || !isString2 || !isString3 || !isString4 {
			response := ResponseWithId{Id: "-1"}
			json.NewEncoder(w).Encode(response)
			fmt.Println("\n Case 6: At least 1 value is not strings")
			return
		}

		X0_string, err := parseJsonTo2DString(data["X0"])

		if err != nil {
			response := ResponseWithId{Id: "-1"}
			json.NewEncoder(w).Encode(response)
			fmt.Println("\n", err)
			return
		}

		X1_string, err := parseJsonTo2DString(data["X1"])

		if err != nil {
			response := ResponseWithId{Id: "-1"}
			json.NewEncoder(w).Encode(response)
			fmt.Println("\n", err)
			return
		}

		X2_string, err := parseJsonTo2DString(data["X2"])

		if err != nil {
			response := ResponseWithId{Id: "-1"}
			json.NewEncoder(w).Encode(response)
			fmt.Println("\n", err)
			return
		}

		Hmu_prime_string, err := parseJsonTo2DString(data["Hmu_prime"])

		if err != nil {
			response := ResponseWithId{Id: "-1"}
			json.NewEncoder(w).Encode(response)
			fmt.Println("\n", err)
			return
		}

		b_zero_message_string, err := parseJsonTo1DString(data["b_zero_message"])

		if err != nil {
			response := ResponseWithId{Id: "-1"}
			json.NewEncoder(w).Encode(response)
			fmt.Println("\n", err)
			return
		}

		return_id, err := verifyAndUpdateSecret(contract, stored_keys_id, stored_secret_id, X0_string, X1_string, X2_string, Hmu_prime_string, b_zero_message_string, owner, signature)

		if err != nil {
			response := ResponseWithId{Id: "-1"}
			json.NewEncoder(w).Encode(response)
			fmt.Println("\n", err)
			return
		}

		response := ResponseWithBool{Bool: return_id}
		json.NewEncoder(w).Encode(response)
		return

	default:
		response := ResponseWithId{Id: "-1"}
		json.NewEncoder(w).Encode(response)
		fmt.Println("\n Case default: Request is invalid ")
		return
	}
}

func main() {

	http.HandleFunc("/", handler)

	// Start HTTPS server
	log.Println("Starting HTTPS server on https://localhost:2800")
	err := http.ListenAndServeTLS(":2800", "middleware_certificate/cert.pem", "middleware_certificate/key.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServeTLS error: ", err)
	}

}

func readFirstFile(dirPath string) ([]byte, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}

	fileNames, err := dir.Readdirnames(1)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(path.Join(dirPath, fileNames[0]))
}

// Submit data for transaction -> returns an id (index in hyperledger)
func storeOwnerData(contract *client.Contract, owner, public_keyD, public_keyUE string) (string, error) {
	fmt.Printf("\n--> Submit Transaction: StoreOwnerPublicKey function\n")

	return_id, err := contract.SubmitTransaction("StoreOwnerPublicKeyMiddleware", owner, public_keyD, public_keyUE)
	if err != nil {
		fmt.Printf("*** storeOwnerData error %v", err)
		return "", err
	}

	fmt.Printf("*** Transaction committed successfully\n")
	return string(return_id), nil
}

// Evaluate a transaction by assetID and owner to query ledger state.
func readOwnerPublicUE(contract *client.Contract, owner_id, owner string) ([][]int, error) {
	fmt.Printf("\n--> Evaluate Transaction: readOwnerPublicUE function\n")

	evaluateResult, err := contract.EvaluateTransaction("ReadOwnerPublicKeyUE", owner_id, owner)
	if err != nil {
		fmt.Printf("*** readOwnerPublicUE error %v", err)
		return nil, err
	}

	result, err2 := ConvertBytesTo2DIntArray(evaluateResult)

	if err2 != nil {
		fmt.Printf("*** ConvertBytesTo2DIntArray error %v", err2)
		return nil, err2
	}

	fmt.Printf("*** Transaction evaluated successfully\n")
	return result, nil
}

// Evaluate a transaction by assetID and owner to query ledger state.
func readOwnerPublicD(contract *client.Contract, owner_id, owner string) (string, error) {
	fmt.Printf("\n--> Evaluate Transaction: readOwnerPublicD function\n")

	evaluateResult, err := contract.EvaluateTransaction("ReadOwnerPublicKeyD", owner_id, owner)
	if err != nil {
		fmt.Printf("*** readOwnerPublicD error %v", err)
		return "", err
	}

	fmt.Printf("*** Transaction evaluated successfully\n")
	return string(evaluateResult), nil
}

// Submit data to be verified and stored -> returns an id if verified (index in hyperledger)
func verifyAndStoreSecret(contract *client.Contract, owner_data_id, H_mu, ciphertext, owner, size, signature string) (string, error) {
	fmt.Printf("\n--> Submit Transaction: verifyAndStoreSecret, function\n")

	return_id, err := contract.SubmitTransaction("VerifyAndStoreSecretMiddleware", owner_data_id, H_mu, ciphertext, owner, size, signature)
	if err != nil {
		fmt.Printf("*** verifyAndStoreSecret error %v", err)
		return "", err
	}

	fmt.Printf("*** Transaction committed successfully\n")
	return string(return_id), nil
}

// Evaluate a transaction by assetID and owner to query ledger state.
func readSecretData(contract *client.Contract, owner_id, owner string) string {
	fmt.Printf("\n--> Evaluate Transaction: readSecretData, function\n")

	evaluateResult, err := contract.EvaluateTransaction("ReadSecretData", owner_id, owner)
	if err != nil {
		fmt.Printf("*** readSecretData error %v", err)
		return ""
	}

	fmt.Printf("*** Transaction evaluated successfully\n")
	return string(evaluateResult)
}

// Submit data to be verified and update data -> returns boolean (did the update occur or not)
func verifyAndUpdateSecret(contract *client.Contract, owner_data_id, secret_id, X0, X1, X2, new_Hmu, b_zero_message, owner, signature string) (bool, error) {
	fmt.Printf("\n--> Submit Transaction: verifyAndUpdateSecret function\n")

	return_value, err := contract.SubmitTransaction("VerifyAndUpdateSecret", owner_data_id, secret_id, X0, X1, X2, new_Hmu, b_zero_message, owner, signature)
	if err != nil {

		fmt.Printf("*** verifyAndUpdateSecret error %v", err)
		return false, err
	}

	//try to convert bytes to string and compare to true
	if string(return_value) != "true" {
		return false, nil
	}

	fmt.Printf("*** Transaction committed successfully\n")
	return true, nil
}

// From Client Json to int 1d array
func parseJsonTo1DArray(inputInterface interface{}) ([]int, error) {

	array, isArray := inputInterface.([]interface{})
	if !isArray {
		return nil, fmt.Errorf("parseJsonToArray - inputInterface is not an array  ")
	}

	outputArray := make([]int, len(array))

	for i := range array {
		outputArray[i] = int(array[i].(float64))
	}

	return outputArray, nil
}

// From Client Json to int 2d array
func parseJsonTo2DArray(inputInterface interface{}) ([][]int, error) {

	array1, isArray1 := inputInterface.([]interface{})
	if !isArray1 {
		return nil, fmt.Errorf("parseJsonToArray - inputInterface is not an array1  ")
	}

	outputArray := make([][]int, len(array1))

	for i := range array1 {
		array2, isArray2 := array1[i].([]interface{})
		if !isArray2 {
			return nil, fmt.Errorf("parseJsonToArray - inputInterface is not an array2  ")
		}

		outputArray[i] = make([]int, len(array2))

		for j := range array2 {
			outputArray[i][j] = int(array2[j].(float64))
		}
	}

	return outputArray, nil
}

// Convert 1D int slice to a JSON string
func Convert1DIntArrayToString(arr []int) (string, error) {
	bytes, err := json.Marshal(arr)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Convert 2D int slice to a JSON string
func Convert2DIntArrayToString(arr [][]int) (string, error) {
	bytes, err := json.Marshal(arr)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Convert JSON []byte to 2D int array
func ConvertBytesTo2DIntArray(data []byte) ([][]int, error) {
	var arr [][]int
	err := json.Unmarshal(data, &arr)
	if err != nil {
		return nil, err
	}
	return arr, nil
}

// From Client Json to string of 2d array
func parseJsonTo1DString(inputInterface interface{}) (string, error) {
	helper1, err1 := parseJsonTo1DArray(inputInterface)

	if err1 != nil {
		return "", err1
	}

	helper2, err2 := Convert1DIntArrayToString(helper1)

	if err2 != nil {
		return "", err2
	}

	return helper2, nil
}

// From Client Json to string of 2d array
func parseJsonTo2DString(inputInterface interface{}) (string, error) {

	helper1, err1 := parseJsonTo2DArray(inputInterface)

	if err1 != nil {
		return "", err1
	}

	helper2, err2 := Convert2DIntArrayToString(helper1)

	if err2 != nil {
		return "", err2
	}

	return helper2, nil
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection() *grpc.ClientConn {
	certificatePEM, err := os.ReadFile(tlsCertPath)
	if err != nil {
		panic(fmt.Errorf("failed to read TLS certifcate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.NewClient(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity() *identity.X509Identity {
	certificatePEM, err := readFirstFile(certPath)
	if err != nil {
		panic(fmt.Errorf("failed to read certificate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign() identity.Sign {
	privateKeyPEM, err := readFirstFile(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}
