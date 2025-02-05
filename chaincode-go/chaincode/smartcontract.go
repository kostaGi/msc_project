package chaincode

import (
	"encoding/hex"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-go/dilithium"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Verify dilithium signiture
func (s *SmartContract) Verify(ctx contractapi.TransactionContextInterface, pk_string string, M_string string, sigma_string string) (bool, error) {

	var i uint32
	str := string(pk_string)
	data, err := hex.DecodeString(str)
	if err != nil {
		return false, err
	}
	pk := make([]uint8, len(data))
	for i = 0; i < uint32(len(data)); i++ {
		pk[i] = data[i]
	}
	//-------------------------------------------------------------------------------------------------------
	str = string(sigma_string)
	data, err = hex.DecodeString(str)
	if err != nil {
		return false, err
	}
	sigma := make([]uint8, len(data))
	for i = 0; i < uint32(len(data)); i++ {
		sigma[i] = data[i]
	}
	//-------------------------------------------------------------------------------------------------------
	str = string(M_string)
	data, err = hex.DecodeString(str)
	if err != nil {
		return false, err
	}
	msg := make([]uint8, len(data))
	for i = 0; i < uint32(len(data)); i++ {
		msg[i] = data[i]
	}

	dilithium.C.InitConstants("2py")
	res := dilithium.Verify(pk, msg, sigma)

	return res, nil
}
