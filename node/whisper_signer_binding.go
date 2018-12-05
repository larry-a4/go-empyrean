// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package node

import (
	"math/big"
	"strings"

	ethereum "github.com/ShyftNetwork/go-empyrean"
	"github.com/ShyftNetwork/go-empyrean/accounts/abi"
	"github.com/ShyftNetwork/go-empyrean/accounts/abi/bind"
	"github.com/ShyftNetwork/go-empyrean/common"
	"github.com/ShyftNetwork/go-empyrean/core/types"
	"github.com/ShyftNetwork/go-empyrean/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// SignerABI is the input ABI used to generate the binding from.
const SignerABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"_signer\",\"type\":\"address\"}],\"name\":\"isValidSigner\",\"outputs\":[{\"name\":\"result\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Signer is an auto generated Go binding around an Ethereum contract.
type Signer struct {
	SignerCaller     // Read-only binding to the contract
	SignerTransactor // Write-only binding to the contract
	SignerFilterer   // Log filterer for contract events
}

// SignerCaller is an auto generated read-only Go binding around an Ethereum contract.
type SignerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SignerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SignerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SignerSession struct {
	Contract     *Signer           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SignerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SignerCallerSession struct {
	Contract *SignerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// SignerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SignerTransactorSession struct {
	Contract     *SignerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SignerRaw is an auto generated low-level Go binding around an Ethereum contract.
type SignerRaw struct {
	Contract *Signer // Generic contract binding to access the raw methods on
}

// SignerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SignerCallerRaw struct {
	Contract *SignerCaller // Generic read-only contract binding to access the raw methods on
}

// SignerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SignerTransactorRaw struct {
	Contract *SignerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSigner creates a new instance of Signer, bound to a specific deployed contract.
func NewSigner(address common.Address, backend bind.ContractBackend) (*Signer, error) {
	contract, err := bindSigner(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Signer{SignerCaller: SignerCaller{contract: contract}, SignerTransactor: SignerTransactor{contract: contract}, SignerFilterer: SignerFilterer{contract: contract}}, nil
}

// NewSignerCaller creates a new read-only instance of Signer, bound to a specific deployed contract.
func NewSignerCaller(address common.Address, caller bind.ContractCaller) (*SignerCaller, error) {
	contract, err := bindSigner(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SignerCaller{contract: contract}, nil
}

// NewSignerTransactor creates a new write-only instance of Signer, bound to a specific deployed contract.
func NewSignerTransactor(address common.Address, transactor bind.ContractTransactor) (*SignerTransactor, error) {
	contract, err := bindSigner(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SignerTransactor{contract: contract}, nil
}

// NewSignerFilterer creates a new log filterer instance of Signer, bound to a specific deployed contract.
func NewSignerFilterer(address common.Address, filterer bind.ContractFilterer) (*SignerFilterer, error) {
	contract, err := bindSigner(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SignerFilterer{contract: contract}, nil
}

// bindSigner binds a generic wrapper to an already deployed contract.
func bindSigner(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SignerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Signer *SignerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Signer.Contract.SignerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Signer *SignerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Signer.Contract.SignerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Signer *SignerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Signer.Contract.SignerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Signer *SignerCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Signer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Signer *SignerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Signer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Signer *SignerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Signer.Contract.contract.Transact(opts, method, params...)
}

// IsValidSigner is a free data retrieval call binding the contract method 0xd5f50582.
//
// Solidity: function isValidSigner(_signer address) constant returns(result bool)
func (_Signer *SignerCaller) IsValidSigner(opts *bind.CallOpts, _signer common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Signer.contract.Call(opts, out, "isValidSigner", _signer)
	return *ret0, err
}

// IsValidSigner is a free data retrieval call binding the contract method 0xd5f50582.
//
// Solidity: function isValidSigner(_signer address) constant returns(result bool)
func (_Signer *SignerSession) IsValidSigner(_signer common.Address) (bool, error) {
	return _Signer.Contract.IsValidSigner(&_Signer.CallOpts, _signer)
}

// IsValidSigner is a free data retrieval call binding the contract method 0xd5f50582.
//
// Solidity: function isValidSigner(_signer address) constant returns(result bool)
func (_Signer *SignerCallerSession) IsValidSigner(_signer common.Address) (bool, error) {
	return _Signer.Contract.IsValidSigner(&_Signer.CallOpts, _signer)
}
