// go get -u golang.org/x/crypto/
// go get golang.org/x/crypto/sha3@v0.32.0
package main

import (
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/sha3"
)

type packet struct {
	a uint8
	b []byte
}

func main() {

	var i uint64
	var pkstring string = "a1ff67dfb26ecb1c4b1e3a7b0ecbbbdc7e06f733e8d529676169aec1ef2a91864556b6a4c920bf355cd8b55c9f3bfe548c3def4e31a0656c6906ef6d90ae3bfd26be06be6a749bf7a3d9377d387f39e7144b5839746163e09c0379740d0405a4778d9763bb703177239b79a731d0a77f7e8d8c8e4c18d9d5beb54d86c8ec63d7ae27314690c7e091100841d849b9541c415ceb5ed55c10f352b2f1627cf20560e25fc6e673fa7c8c3042af98a721ec3b1767d74bea3ca835b453eaf65454beb7d029963addecbd6ee1fbbadbb65d6cb503933aad5b0b73fa06e22605a12d6a3f0ee17aee48e8c6a17c4d3a2f90f8e6b66dd44996956fab25d94350de49a81cff324bd7b95b83dfe49976dcc5b491b0c2ec0ed3e0df5c268b45dffb56a23f37a38d41c24986d7e5449d91a1f749fc231f36172615ba6c1aa2bb77dbbafc74e9a96eddbeb21fc3be9c20dfa789148c120f9d3edac2e48f9dcd5a266b91367bb08f8c2797bf39c4bbcde409eb00e9a7dc2f6fd69e23c8c0d56923382716d1c0110175a0a8611375aa51d753897ecfd6e3ce8f284a6573b9a36b5cfb014315ad372c7d6d18ecfa43ef375be16f2aea1af08c6b056eed7ef0b6933445ba570fb2d4e10ba8be1f368ba49e1fe057b765af8a83bcaff27398f0c136a2866baa5f29afa4462ab1dd0328da0ba8ef003df625f508cb1d5566445886085526eaa94d38bc4cf6195d7ae82f8534c7b83a85b596d3f62c909e62a8908b9caa42af63c75f5633d2d1aa8799ccb2ba72261f61fd4b55c2c0046a787215e8c4c6125ce0bc2a00bb282df83f8d0d25883a6cf53dd2d3f58cb1cf078815cd6560cbe05f8322819d41e8bad0d5ce956b331269b74023bf3cbb1ba52208ec157d067446b8f667ccc852023796531b5b69ecefcac590ab1b6b7a88a1adafb15986b157346106415a1319597c3eae5f0d01e6fa2e6e993c81599402bf0b32d4f923b49d42876a2b07ef99dd38100c466226558e012ea92bddc6cbcaa7025df2e97a705c88dfb7f2db33ced0da26bd7cd1dd053f8423b95a42df737ae2412225f50d89743e33cb590005a3b33b8224979aadcca3135f903f1d8c0017c9e90afc90cea040377d3a1cd1b6ee0e047948040245a141da83a2307061647b19970c8d9d63638770576d152af07c97dc4915974c08777739c0be61a7f6d0d4580f1aa398e7edbbd0b4569cd558ba9a3e6477208513feb8ada0b80d6bc2a7a2305a208597f258f3f093ea8f4862f126733e9b591c609ed679dd27b8640f038b22766c145cadc9a937198fd132343c0b4dade8a0a28774cca73f8027bad25c52bdf587e58bc3b927709ee411c54ccbbc497c4dd22a27c107660e04951f36ae81f32afc581b74d2ba4bfd1c3c1f0f45f04eae4650b68e433fd509c6a03993963f270dd94ff5975107cfb8bbfdcf874250706658e280b76d792378d2790d2f1fb77780375682e0877f8333ffbd03df194a96d6ca64966a27f1f5e34a1a61c5ac79055e623b7b90b63445ec83124e73f94d49e647687f352021bb96548b5a21a901207b44ca4897c2aa9f903f1ef91f3e3708d54ce3d2e54e7023a4af6d3315b729816d8d1adb4e6cb8f0f1e5c63a0506e3f284fa2af4577311f395eda6d80ced605574d447d87e50815595a878240ef08407e9e5b3c0dae80d8642f5265a6c38cf9cb3c2138d0ca47d174f4b37b523ba43ff4ae50dd53f6dd58b02ae36c709d7dc83b575e244227ea83dc070d8b45477ebceb732d439edfed8fd50b7ee53847f526294c7a1b33cba35e39842e332ed355be0a501b5eb56ce9d612edb9b49d7b5ad005050f06d3155bc0bc06595f5c584"
	pk := make([]uint8, len(pkstring)/2)

	data, err := hex.DecodeString(pkstring)
	if err != nil {
		panic(err)
	}
	fmt.Println("pk len", len(data), len(pkstring)/2)
	for i = 0; i < uint64(len(data)); i++ {
		pk[i] = data[i]
	}

	// Example 1: Simple cshake

	tr := make([]uint8, 48)
	sha3.ShakeSum256(tr, pk)

	fmt.Printf("%02x", tr)

	/*
		var p = packet{}
		p.a = 1
		p.b = []byte("foo\x00\x00")
		buf := bytes.Buffer{}
		err := binary.Write(&buf, binary.BigEndian, p.a)
		if err != nil {
			fmt.Println(err)
		}
		_, err = buf.Write(p.b)
		if err != nil {
			fmt.Println(err)
		}
	*/

	/*
		h := sha256.New()
		//h.Write(buf.Bytes())
		wr, err := h.Write(pk)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("written=", wr)
		hash := h.Sum([]byte{})
		fmt.Printf("% x\n", hash)
		//hash = h.Sum([]byte{})
		//fmt.Printf("% x\n", hash)
	*/
}
