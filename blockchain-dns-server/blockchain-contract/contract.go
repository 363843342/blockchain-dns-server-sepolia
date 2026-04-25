package contract

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	// 修正点：确保直接导入，不加 'contract' 别名
	"ayls/blockchain-dns-server/blockchain-contract/dnsrecord"
)

const (
	// ErrTransactionWait 当交易未确认时返回
	ErrTransactionWait = "if you've just started the application, wait a while for the network to confirm your transaction."
)

// NewContract 部署合约 [cite: 9]
func NewContract(session dnsrecord.DnsrecordSession, client *ethclient.Client) (dnsrecord.DnsrecordSession, string) {
	contractAddress, tx, instance, err := dnsrecord.DeployDnsrecord(&session.TransactOpts, client)
	if err != nil {
		log.Fatalf("could not deploy contract: %v\n", err)
	}
	fmt.Printf("Contract deployed! Wait for tx %s to be confirmed.\n", tx.Hash().Hex())

	session.Contract = instance
	return session, contractAddress.Hex()
}

// LoadContract 加载合约
func LoadContract(session dnsrecord.DnsrecordSession, client *ethclient.Client, address string) dnsrecord.DnsrecordSession {
	if address == "" {
		log.Println("could not find a contract address to load")
		return session
	}
	addr := common.HexToAddress(address)
	instance, err := dnsrecord.NewDnsrecord(addr, client)
	if err != nil {
		log.Fatalf("could not load contract: %v\n", err)
		log.Println(ErrTransactionWait)
	}
	session.Contract = instance
	return session
}

// NewSession 初始化并解决 Sepolia Chain ID 问题 
func NewSession(ctx context.Context, keystorePath string, keystorePass string) (session dnsrecord.DnsrecordSession) {
	keystore, err := os.Open(keystorePath)
	if err != nil {
		log.Fatalf("could not load keystore from location %s: %v\n", keystorePath, err)
	}
	defer keystore.Close()

	// 适配 Sepolia 网络 (Chain ID: 11155111) 
	chainID := big.NewInt(11155111)
	
	// 使用可直接接收文件流和密码的函数 
	auth, err := bind.NewTransactorWithChainID(keystore, keystorePass, chainID)
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v\n", err)
	}

	return dnsrecord.DnsrecordSession{
		TransactOpts: *auth,
		CallOpts: bind.CallOpts{
			From:    auth.From,
			Context: ctx,
		},
	}
}