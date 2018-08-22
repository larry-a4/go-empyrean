package ethobjectinterface

import "github.com/ShyftNetwork/go-empyrean/eth"

// getEthObject()Get var EthereumObject *eth.Ethereum
func GetEthObject() interface{} {
	return eth.EthereumObject
}
