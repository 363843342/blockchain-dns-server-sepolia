// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package dnsrecord

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// DnsrecordMetaData contains all meta data concerning the Dnsrecord contract.
var DnsrecordMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"internalType\":\"uint16\",\"name\":\"recType\",\"type\":\"uint16\"},{\"internalType\":\"string\",\"name\":\"recValue\",\"type\":\"string\"}],\"name\":\"addRecord\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"internalType\":\"uint16\",\"name\":\"recType\",\"type\":\"uint16\"}],\"name\":\"getRecord\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b5061074b8061001c5f395ff3fe608060405234801561000f575f5ffd5b5060043610610034575f3560e01c80636719bce114610038578063fb0d787d14610068575b5f5ffd5b610052600480360381019061004d9190610239565b610084565b60405161005f9190610306565b60405180910390f35b610082600480360381019061007d9190610326565b61014c565b005b60605f84846040516100979291906103f3565b90815260200160405180910390205f8361ffff1661ffff1681526020019081526020015f2080546100c790610438565b80601f01602080910402602001604051908101604052809291908181526020018280546100f390610438565b801561013e5780601f106101155761010080835404028352916020019161013e565b820191905f5260205f20905b81548152906001019060200180831161012157829003601f168201915b505050505090509392505050565b81815f878760405161015f9291906103f3565b90815260200160405180910390205f8661ffff1661ffff1681526020019081526020015f209182610191929190610648565b505050505050565b5f5ffd5b5f5ffd5b5f5ffd5b5f5ffd5b5f5ffd5b5f5f83601f8401126101c2576101c16101a1565b5b8235905067ffffffffffffffff8111156101df576101de6101a5565b5b6020830191508360018202830111156101fb576101fa6101a9565b5b9250929050565b5f61ffff82169050919050565b61021881610202565b8114610222575f5ffd5b50565b5f813590506102338161020f565b92915050565b5f5f5f604084860312156102505761024f610199565b5b5f84013567ffffffffffffffff81111561026d5761026c61019d565b5b610279868287016101ad565b9350935050602061028c86828701610225565b9150509250925092565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f6102d882610296565b6102e281856102a0565b93506102f28185602086016102b0565b6102fb816102be565b840191505092915050565b5f6020820190508181035f83015261031e81846102ce565b905092915050565b5f5f5f5f5f6060868803121561033f5761033e610199565b5b5f86013567ffffffffffffffff81111561035c5761035b61019d565b5b610368888289016101ad565b9550955050602061037b88828901610225565b935050604086013567ffffffffffffffff81111561039c5761039b61019d565b5b6103a8888289016101ad565b92509250509295509295909350565b5f81905092915050565b828183375f83830152505050565b5f6103da83856103b7565b93506103e78385846103c1565b82840190509392505050565b5f6103ff8284866103cf565b91508190509392505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061044f57607f821691505b6020821081036104625761046161040b565b5b50919050565b5f82905092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026104fb7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826104c0565b61050586836104c0565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f61054961054461053f8461051d565b610526565b61051d565b9050919050565b5f819050919050565b6105628361052f565b61057661056e82610550565b8484546104cc565b825550505050565b5f5f905090565b61058d61057e565b610598818484610559565b505050565b5b818110156105bb576105b05f82610585565b60018101905061059e565b5050565b601f821115610600576105d18161049f565b6105da846104b1565b810160208510156105e9578190505b6105fd6105f5856104b1565b83018261059d565b50505b505050565b5f82821c905092915050565b5f6106205f1984600802610605565b1980831691505092915050565b5f6106388383610611565b9150826002028217905092915050565b6106528383610468565b67ffffffffffffffff81111561066b5761066a610472565b5b6106758254610438565b6106808282856105bf565b5f601f8311600181146106ad575f841561069b578287013590505b6106a5858261062d565b86555061070c565b601f1984166106bb8661049f565b5f5b828110156106e2578489013582556001820191506020850194506020810190506106bd565b868310156106ff57848901356106fb601f891682610611565b8355505b6001600288020188555050505b5050505050505056fea264697066735822122060c1d0312199bcde8b8cefc299984e669222c62533d7a5d20dcd25f4c47dd9db64736f6c634300081e0033",
}

// DnsrecordABI is the input ABI used to generate the binding from.
// Deprecated: Use DnsrecordMetaData.ABI instead.
var DnsrecordABI = DnsrecordMetaData.ABI

// DnsrecordBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DnsrecordMetaData.Bin instead.
var DnsrecordBin = DnsrecordMetaData.Bin

// DeployDnsrecord deploys a new Ethereum contract, binding an instance of Dnsrecord to it.
func DeployDnsrecord(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Dnsrecord, error) {
	parsed, err := DnsrecordMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DnsrecordBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Dnsrecord{DnsrecordCaller: DnsrecordCaller{contract: contract}, DnsrecordTransactor: DnsrecordTransactor{contract: contract}, DnsrecordFilterer: DnsrecordFilterer{contract: contract}}, nil
}

// Dnsrecord is an auto generated Go binding around an Ethereum contract.
type Dnsrecord struct {
	DnsrecordCaller     // Read-only binding to the contract
	DnsrecordTransactor // Write-only binding to the contract
	DnsrecordFilterer   // Log filterer for contract events
}

// DnsrecordCaller is an auto generated read-only Go binding around an Ethereum contract.
type DnsrecordCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DnsrecordTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DnsrecordTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DnsrecordFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DnsrecordFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DnsrecordSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DnsrecordSession struct {
	Contract     *Dnsrecord        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DnsrecordCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DnsrecordCallerSession struct {
	Contract *DnsrecordCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// DnsrecordTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DnsrecordTransactorSession struct {
	Contract     *DnsrecordTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// DnsrecordRaw is an auto generated low-level Go binding around an Ethereum contract.
type DnsrecordRaw struct {
	Contract *Dnsrecord // Generic contract binding to access the raw methods on
}

// DnsrecordCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DnsrecordCallerRaw struct {
	Contract *DnsrecordCaller // Generic read-only contract binding to access the raw methods on
}

// DnsrecordTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DnsrecordTransactorRaw struct {
	Contract *DnsrecordTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDnsrecord creates a new instance of Dnsrecord, bound to a specific deployed contract.
func NewDnsrecord(address common.Address, backend bind.ContractBackend) (*Dnsrecord, error) {
	contract, err := bindDnsrecord(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Dnsrecord{DnsrecordCaller: DnsrecordCaller{contract: contract}, DnsrecordTransactor: DnsrecordTransactor{contract: contract}, DnsrecordFilterer: DnsrecordFilterer{contract: contract}}, nil
}

// NewDnsrecordCaller creates a new read-only instance of Dnsrecord, bound to a specific deployed contract.
func NewDnsrecordCaller(address common.Address, caller bind.ContractCaller) (*DnsrecordCaller, error) {
	contract, err := bindDnsrecord(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DnsrecordCaller{contract: contract}, nil
}

// NewDnsrecordTransactor creates a new write-only instance of Dnsrecord, bound to a specific deployed contract.
func NewDnsrecordTransactor(address common.Address, transactor bind.ContractTransactor) (*DnsrecordTransactor, error) {
	contract, err := bindDnsrecord(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DnsrecordTransactor{contract: contract}, nil
}

// NewDnsrecordFilterer creates a new log filterer instance of Dnsrecord, bound to a specific deployed contract.
func NewDnsrecordFilterer(address common.Address, filterer bind.ContractFilterer) (*DnsrecordFilterer, error) {
	contract, err := bindDnsrecord(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DnsrecordFilterer{contract: contract}, nil
}

// bindDnsrecord binds a generic wrapper to an already deployed contract.
func bindDnsrecord(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DnsrecordMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Dnsrecord *DnsrecordRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Dnsrecord.Contract.DnsrecordCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Dnsrecord *DnsrecordRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Dnsrecord.Contract.DnsrecordTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Dnsrecord *DnsrecordRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Dnsrecord.Contract.DnsrecordTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Dnsrecord *DnsrecordCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Dnsrecord.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Dnsrecord *DnsrecordTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Dnsrecord.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Dnsrecord *DnsrecordTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Dnsrecord.Contract.contract.Transact(opts, method, params...)
}

// GetRecord is a free data retrieval call binding the contract method 0x6719bce1.
//
// Solidity: function getRecord(string key, uint16 recType) view returns(string)
func (_Dnsrecord *DnsrecordCaller) GetRecord(opts *bind.CallOpts, key string, recType uint16) (string, error) {
	var out []interface{}
	err := _Dnsrecord.contract.Call(opts, &out, "getRecord", key, recType)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetRecord is a free data retrieval call binding the contract method 0x6719bce1.
//
// Solidity: function getRecord(string key, uint16 recType) view returns(string)
func (_Dnsrecord *DnsrecordSession) GetRecord(key string, recType uint16) (string, error) {
	return _Dnsrecord.Contract.GetRecord(&_Dnsrecord.CallOpts, key, recType)
}

// GetRecord is a free data retrieval call binding the contract method 0x6719bce1.
//
// Solidity: function getRecord(string key, uint16 recType) view returns(string)
func (_Dnsrecord *DnsrecordCallerSession) GetRecord(key string, recType uint16) (string, error) {
	return _Dnsrecord.Contract.GetRecord(&_Dnsrecord.CallOpts, key, recType)
}

// AddRecord is a paid mutator transaction binding the contract method 0xfb0d787d.
//
// Solidity: function addRecord(string key, uint16 recType, string recValue) returns()
func (_Dnsrecord *DnsrecordTransactor) AddRecord(opts *bind.TransactOpts, key string, recType uint16, recValue string) (*types.Transaction, error) {
	return _Dnsrecord.contract.Transact(opts, "addRecord", key, recType, recValue)
}

// AddRecord is a paid mutator transaction binding the contract method 0xfb0d787d.
//
// Solidity: function addRecord(string key, uint16 recType, string recValue) returns()
func (_Dnsrecord *DnsrecordSession) AddRecord(key string, recType uint16, recValue string) (*types.Transaction, error) {
	return _Dnsrecord.Contract.AddRecord(&_Dnsrecord.TransactOpts, key, recType, recValue)
}

// AddRecord is a paid mutator transaction binding the contract method 0xfb0d787d.
//
// Solidity: function addRecord(string key, uint16 recType, string recValue) returns()
func (_Dnsrecord *DnsrecordTransactorSession) AddRecord(key string, recType uint16, recValue string) (*types.Transaction, error) {
	return _Dnsrecord.Contract.AddRecord(&_Dnsrecord.TransactOpts, key, recType, recValue)
}
