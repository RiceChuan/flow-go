package contracts

import (
	_ "embed"
	"encoding/hex"
)

//go:embed test_bytes.hex
var testContractBytesInHex string

var TestContractBytes, _ = hex.DecodeString(testContractBytesInHex)

//go:embed test_abi.json
var TestContractABIJSON string

//go:embed dummy_kitty_bytes.hex
var dummyKittyContractBytesInHex string

var DummyKittyContractBytes, _ = hex.DecodeString(dummyKittyContractBytesInHex)

//go:embed dummy_kitty_abi.json
var DummyKittyContractABIJSON string

//go:embed factory_bytes.hex
var factoryContractBytesInHex string

var FactoryContractBytes, _ = hex.DecodeString(factoryContractBytesInHex)

//go:embed factory_abi.json
var FactoryContractABIJSON string
