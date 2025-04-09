package chaincode

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-go/dilithium"
	"github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-go/updatable_encryption"
)

// index of data on ledger
var id_counter int = 0

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// public keys of owner
type OwnerData struct {
	ID                           string  `json:"ID"`
	Owner                        string  `json:"Owner"`
	PublicKeyDilithium           string  `json:"PublicKeyDilithium"`
	PublicKeyUpdatableEncryption [][]int `json:"PublicKeyUpdatableEncryption"`
}

// secret stored by owner
type SecretData struct {
	ID         string  `json:"ID"`
	Hmu        [][]int `json:"Hmu"`
	Ciphertext []int   `json:"Ciphertext"`
	Owner      string  `json:"Owner"`
	Size       int     `json:"Size"`
}

// Store both public keys on ledger -> returns id (index on the ledger) (callable by middleware)
func (s *SmartContract) StoreOwnerPublicKeyMiddleware(ctx contractapi.TransactionContextInterface, owner, public_key_D, public_key_UE string) (string, error) {

	public_key_UE_transform, err := convertStringTo2DIntArray(public_key_UE)

	if err != nil {
		return "", err
	}

	return s.StoreOwnerPublicKey(ctx, owner, public_key_D, public_key_UE_transform)
}

// Store both public keys on ledger -> returns id (index on the ledger)
func (s *SmartContract) StoreOwnerPublicKey(ctx contractapi.TransactionContextInterface, owner, public_key_D string, public_key_UE [][]int) (string, error) {

	id := strconv.Itoa(id_counter)
	id_counter += 1

	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("the asset %s already exists", id)
	}

	asset := OwnerData{
		ID:                           id,
		Owner:                        owner,
		PublicKeyDilithium:           public_key_D,
		PublicKeyUpdatableEncryption: public_key_UE,
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return "", err
	}

	return id, ctx.GetStub().PutState(id, assetJSON)
}

// Return Updatable encryption key on ledger (callable by middleware) -> returns updatable encryption key
func (s *SmartContract) ReadOwnerPublicKeyUE(ctx contractapi.TransactionContextInterface, id string, owner string) ([][]int, error) {

	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("the asset %s does not exists", id)
	}

	asset, err := s.ReadOwnerData(ctx, id, owner)

	if err != nil {
		return nil, err
	}

	return asset.PublicKeyUpdatableEncryption, nil
}

// Return Dilithium key on ledger (callable by middleware) -> returns dilithium key
func (s *SmartContract) ReadOwnerPublicKeyD(ctx contractapi.TransactionContextInterface, id string, owner string) (string, error) {

	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("the asset %s does not exists", id)
	}

	asset, err := s.ReadOwnerData(ctx, id, owner)

	if err != nil {
		return "", err
	}

	return asset.PublicKeyDilithium, nil
}

// Return OwnerData by id and owner -> returns OwnerData
func (s *SmartContract) ReadOwnerData(ctx contractapi.TransactionContextInterface, id string, owner string) (*OwnerData, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset OwnerData
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	if asset.Owner != owner {
		return nil, fmt.Errorf("the owner %s does not match: %s", owner, asset.Owner)
	}

	return &asset, nil
}

// Store a secret on ledger -> returns id (index on the ledger)
func (s *SmartContract) StoreSecret(ctx contractapi.TransactionContextInterface, id string, H_mu [][]int, ciphertext []int, owner string, size int) error {

	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}

	// does not reaveal infomation as you need both id and owner to get output
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := SecretData{
		ID:         id,
		Hmu:        H_mu,
		Ciphertext: ciphertext,
		Owner:      owner,
		Size:       size,
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// Updates a secret on the ledger
func (s *SmartContract) UpdateSecret(ctx contractapi.TransactionContextInterface, id, owner string, X0, X1, X2, new_Hmu [][]int, b_zero_message []int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	bytes, err := s.ReadSecretData(ctx, id, owner)

	if err != nil {
		return err
	}

	var old_asset SecretData

	err1 := json.Unmarshal([]byte(bytes), &old_asset)
	if err1 != nil {
		return err1
	}

	updatable_encryption.C.InitConstants()
	new_ciphertext, err := updatable_encryption.Update2(X0, X1, X2, b_zero_message, old_asset.Ciphertext, old_asset.Size)

	if err != nil {
		return fmt.Errorf("error is here %v", err)
	}

	new_asset := SecretData{
		ID:         old_asset.ID,
		Hmu:        new_Hmu,
		Ciphertext: new_ciphertext,
		Owner:      old_asset.Owner,
		Size:       old_asset.Size,
	}

	assetJSON, err := json.Marshal(new_asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// Read a secret on ledger -> returns Json string of(SecretData)
func (s *SmartContract) ReadSecretData(ctx contractapi.TransactionContextInterface, id string, owner string) (string, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return "", fmt.Errorf("the asset %s does not exist", id)
	}

	return string(assetJSON), nil
}

// Check if asset exist on ledger -> returns bool
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// Verify dilithium signature and store secret on ledger (callable by middleware) -> return id (index on ledger)
func (s *SmartContract) VerifyAndStoreSecretMiddleware(ctx contractapi.TransactionContextInterface, owner_data_id, H_mu, ciphertext, owner, size, signature string) (string, error) {

	owner_data, err := s.ReadOwnerData(ctx, owner_data_id, owner)

	if err != nil {
		return "", err
	}

	var total_message strings.Builder
	total_message.WriteString(H_mu)
	total_message.WriteString(ciphertext)
	total_message.WriteString(owner)
	total_message.WriteString(size)

	//converts string to utf-8 - done to client side as well
	message_array := []byte(total_message.String())

	//converts hex to unsigned int
	data1, err := hex.DecodeString(owner_data.PublicKeyDilithium)
	if err != nil {
		return "", err
	}

	public_key_array := make([]uint8, len(data1))
	for i := uint32(0); i < uint32(len(data1)); i++ {
		public_key_array[i] = data1[i]
	}

	//converts hex to unsigned int
	data2, err := hex.DecodeString(signature)
	if err != nil {
		return "", err
	}

	signature_array := make([]uint8, len(data2))
	for i := uint32(0); i < uint32(len(data2)); i++ {
		signature_array[i] = data2[i]
	}

	isVerified, _ := verifyDilithiumSignature(public_key_array, message_array, signature_array)

	if !isVerified {
		return "", fmt.Errorf("VerifyAndStoreSecret - method isVerified is false")
	}

	id := strconv.Itoa(id_counter)
	id_counter += 1

	H_mu_arr, err := convertStringTo2DIntArray(H_mu)
	if err != nil {
		return "", nil
	}
	ciphertext_arr, err := convertStringTo1DIntArray(ciphertext)
	if err != nil {
		return "", err
	}
	size_int, err := convertStringToInt(size)
	if err != nil {
		return "", err
	}

	err2 := s.StoreSecret(ctx, id, H_mu_arr, ciphertext_arr, owner, size_int)

	if err2 != nil {
		return "", err2
	}

	return id, nil
}

// Verify dilithium signature and updates secret on ledger (callable by middleware) -> return bool
func (s *SmartContract) VerifyAndUpdateSecret(ctx contractapi.TransactionContextInterface, owner_data_id, secret_id, X0, X1, X2, new_Hmu, b_zero_message, owner, signature string) (bool, error) {

	owner_data, err := s.ReadOwnerData(ctx, owner_data_id, owner)

	if err != nil {
		return false, err
	}

	var total_message strings.Builder

	total_message.WriteString(X0)
	total_message.WriteString(X1)
	total_message.WriteString(X2)
	total_message.WriteString(new_Hmu)
	total_message.WriteString(b_zero_message)
	total_message.WriteString(owner)

	//converts string to utf-8 - done to client side as well
	message_array := []byte(total_message.String())

	//converts hex to unsigned int
	data1, err := hex.DecodeString(owner_data.PublicKeyDilithium)
	if err != nil {
		return false, err
	}

	public_key_array := make([]uint8, len(data1))
	for i := uint32(0); i < uint32(len(data1)); i++ {
		public_key_array[i] = data1[i]
	}

	//converts hex to unsigned int
	data2, err := hex.DecodeString(signature)
	if err != nil {
		return false, err
	}

	signature_array := make([]uint8, len(data2))
	for i := uint32(0); i < uint32(len(data2)); i++ {
		signature_array[i] = data2[i]
	}

	isVerified, _ := verifyDilithiumSignature(public_key_array, message_array, signature_array)

	if !isVerified {
		return false, fmt.Errorf("VerifyAndUpdateSecret - method isVerified is false")
	}

	X0_arr, err := convertStringTo2DIntArray(X0)
	if err != nil {
		return false, err
	}

	X1_arr, err := convertStringTo2DIntArray(X1)
	if err != nil {
		return false, err
	}

	X2_arr, err := convertStringTo2DIntArray(X2)
	if err != nil {
		return false, err
	}

	new_Hmu_arr, err := convertStringTo2DIntArray(new_Hmu)
	if err != nil {
		return false, err
	}

	b_zero_message_arr, err := convertStringTo1DIntArray(b_zero_message)
	if err != nil {
		return false, err
	}

	err2 := s.UpdateSecret(ctx, secret_id, owner, X0_arr, X1_arr, X2_arr, new_Hmu_arr, b_zero_message_arr)

	if err2 != nil {
		return false, err2
	}

	return true, nil
}

// TO DO (add error for dilithium golang - good practice)
// Verify dilithium signature
func verifyDilithiumSignature(public_key, input_message, signature []uint8) (bool, error) {
	dilithium.C.InitConstants("2py")
	return dilithium.Verify(public_key, input_message, signature), nil
}

// Convert JSON string back to 1D int array
func convertStringTo1DIntArray(inputstr string) ([]int, error) {
	var arr []int
	err := json.Unmarshal([]byte(inputstr), &arr)
	if err != nil {
		return nil, err
	}
	return arr, nil
}

// Convert JSON string back to 2D int array
func convertStringTo2DIntArray(inputstr string) ([][]int, error) {
	var arr [][]int
	err := json.Unmarshal([]byte(inputstr), &arr)
	if err != nil {
		return nil, err
	}
	return arr, nil
}

// Convert string to int  "100" -> 100
func convertStringToInt(inputstr string) (int, error) {
	num, err := strconv.ParseInt(inputstr, 10, 64) // Parse as base 10 and store as int64

	if err != nil {
		return -1, err
	}
	return int(num), nil
}
