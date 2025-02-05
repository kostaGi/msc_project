package main

//go run -tags debug v2.go

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-go/dilithium"
)

func test_verify_false(strMode string) bool {
	dilithium.C.InitConstants(strMode)

	var i uint32
	//-------------------------------------------------------------------------------------------------------
	b, err := os.ReadFile("./data/pk.pem")
	if err != nil {
		panic(err)
	}
	str := string(b)
	data, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	pk := make([]uint8, len(data))
	for i = 0; i < uint32(len(data)); i++ {
		pk[i] = data[i]
	}
	//-------------------------------------------------------------------------------------------------------
	b, err = os.ReadFile("./data/sigma.tst")
	if err != nil {
		panic(err)
	}
	str = string(b)
	data, err = hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	sigma := make([]uint8, len(data))
	for i = 0; i < uint32(len(data)); i++ {
		sigma[i] = data[i]
	}
	//-------------------------------------------------------------------------------------------------------
	b, err = os.ReadFile("./data/msg2.tst")
	if err != nil {
		fmt.Print(err)
	}
	str = string(b)
	data, err = hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	b2 := make([]uint8, len(data))
	for i = 0; i < uint32(len(data)); i++ {
		b2[i] = data[i]
	}
	if !dilithium.Verify(pk, b2, sigma) {
		fmt.Println("Verification failed")
		return false
	}

	fmt.Println("Verification success")
	return true
}

func main() {
	test_verify_false("2py")
}
