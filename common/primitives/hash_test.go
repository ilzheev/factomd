// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package primitives_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/FactomProject/ed25519"
	"github.com/FactomProject/factomd/common/constants"
	. "github.com/FactomProject/factomd/common/primitives"
	"math/rand"
	"testing"
)

var _ = fmt.Printf
var _ = ed25519.Sign
var _ = rand.New

// A hash
var hash = [constants.ADDRESS_LENGTH]byte{
	0x61, 0xe3, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72, 0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
	0x93, 0x1e, 0x83, 0x65, 0xe1, 0x5a, 0x08, 0x9c, 0x68, 0xd6, 0x19, 0x00, 0x00, 0x00, 0x00, 0x00,
}

func Test_HashEquals(test *testing.T) {
	h1 := new(Hash)
	h2 := new(Hash)

	if h1.IsEqual(h2) != nil { // Out of the box, hashes should be equal
		PrtStk()
		test.Fail()
	}

	h1.SetBytes(hash[:])

	if h1.IsEqual(h2) == nil { // Now they should not be equal
		PrtStk()
		test.Fail()
	}

	h2.SetBytes(hash[:])

	if h1.IsEqual(h2) != nil { // Back to equality!
		PrtStk()
		test.Fail()
	}
}

//Test vectors: http://www.di-mgt.com.au/sha_testvectors.html

func TestHash(t *testing.T) {
	fmt.Println("\nTest hash===========================================================================")

	h := new(Hash)
	err := h.SetBytes(constants.EC_CHAINID)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	bytes1, err := h.MarshalBinary()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("bytes1: %v\n", bytes1)

	h2 := new(Hash)
	err = h2.UnmarshalBinary(bytes1)
	t.Logf("h2.bytes: %v\n", h2.Bytes)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	bytes2, err := h2.MarshalBinary()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("bytes2: %v\n", bytes2)

	if bytes.Compare(bytes1, bytes2) != 0 {
		t.Errorf("Invalid output")
	}
}

func TestSha(t *testing.T) {
	testVector := map[string]string{}
	testVector["abc"] = "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	testVector[""] = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	testVector["abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq"] = "248d6a61d20638b8e5c026930c3e6039a33ce45964ff2167f6ecedd419db06c1"
	testVector["abcdefghbcdefghicdefghijdefghijkefghijklfghijklmghijklmnhijklmnoijklmnopjklmnopqklmnopqrlmnopqrsmnopqrstnopqrstu"] = "cf5b16a778af8380036ce59e7b0492370b249b11e8f07a51afac45037afee9d1"

	for k, v := range testVector {
		answer, err := DecodeBinary(v)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		hash := Sha([]byte(k))

		if bytes.Compare(hash.Bytes(), answer) != 0 {
			t.Errorf("Wrong SHA hash for %v", k)
		}
		if hash.String() != v {
			t.Errorf("Wrong SHA hash string for %v", k)
		}
	}
}

func TestSha512Half(t *testing.T) {
	testVector := map[string]string{}
	testVector["abc"] = "ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a"
	testVector[""] = "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce"
	testVector["abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq"] = "204a8fc6dda82f0a0ced7beb8e08a41657c16ef468b228a8279be331a703c335"
	testVector["abcdefghbcdefghicdefghijdefghijkefghijklfghijklmghijklmnhijklmnoijklmnopjklmnopqklmnopqrlmnopqrsmnopqrstnopqrstu"] = "8e959b75dae313da8cf4f72814fc143f8f7779c6eb9f7fa17299aeadb6889018"

	for k, v := range testVector {
		answer, err := DecodeBinary(v)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		hash := Sha512Half([]byte(k))

		if bytes.Compare(hash.Bytes(), answer) != 0 {
			t.Errorf("Wrong SHA512Half hash for %v", k)
		}
		if hash.String() != v {
			t.Errorf("Wrong SHA512Half hash string for %v", k)
		}
	}
}

func TestStrings(t *testing.T) {
	base := "ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a"
	hash, err := HexToHash(base)
	if err != nil {
		t.Error(err)
	}
	if hash.String() != base {
		t.Error("Invalid conversion to string")
	}
}

func TestIsSameAs(t *testing.T) {
	base := "ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a"
	hash, err := HexToHash(base)
	if err != nil {
		t.Error(err)
	}
	hex, err := DecodeBinary(base)
	if err != nil {
		t.Error(err)
	}
	hash2, err := NewShaHash(hex)
	if err != nil {
		t.Error(err)
	}
	if hash.IsSameAs(hash2) == false {
		t.Error("Identical hashes not recognized as such")
	}
}

func TestHashMisc(t *testing.T) {
	base := "4040404040404040404040404040404040404040404040404040404040404040"
	hash, err := HexToHash(base)
	if err != nil {
		t.Error(err)
	}
	if hash.String() != base {
		t.Error("Error in String")
	}

	hash2, err := NewShaHashFromStr(base)
	if err != nil {
		t.Error(err)
	}

	if hash2.ByteString() != "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@" {
		t.Errorf("Error in ByteString - received %v", hash2.ByteString())
	}

	h, err := hex.DecodeString(base)
	if err != nil {
		t.Error(err)
	}
	hash = NewHash(h)
	if hash.String() != base {
		t.Error("Error in NewHash")
	}

	//***********************

	if hash.IsSameAs(nil) != false {
		t.Error("Error in IsSameAs")
	}

	//***********************

	minuteHash, err := HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
	if err != nil {
		t.Error(err)
	}
	if minuteHash.IsMinuteMarker() == false {
		t.Error("Error in IsMinuteMarker")
	}

	hash = NewZeroHash()
	if hash.String() != "0000000000000000000000000000000000000000000000000000000000000000" {
		t.Error("Error in NewZeroHash")
	}
}

func TestStringUnmarshaller(t *testing.T) {
	base := "ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a"
	hash, err := HexToHash(base)
	if err != nil {
		t.Error(err)
	}

	h2 := new(Hash)
	err = h2.UnmarshalText([]byte(base))
	if err != nil {
		t.Error(err)
	}
	if hash.IsSameAs(h2) == false {
		t.Errorf("Hash from UnmarshalText is incorrect - %v vs %v", hash, h2)
	}

	h3 := new(Hash)
	err = json.Unmarshal([]byte("\""+base+"\""), h3)
	if err != nil {
		t.Error(err)
	}
	if hash.IsSameAs(h3) == false {
		t.Errorf("Hash from json.Unmarshal is incorrect - %v vs %v", hash, h3)
	}
}
