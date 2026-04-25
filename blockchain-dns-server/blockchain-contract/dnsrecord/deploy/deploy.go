package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	contract "ayls/blockchain-dns-server/blockchain-contract/dnsrecord"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	rpc := "https://ethereum-sepolia-rpc.publicnode.com"

	client, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal(err)
	}

	privateKeyHex := "478c0d7de72320aa8ae4701b72852c0ba04fdc829da6976d64f057dc68ae5f63"

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	fmt.Println("Deploying from:", fromAddress.Hex())

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	chainID := big.NewInt(11155111)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth.GasPrice = gasPrice

	address, tx, _, err :=
		contract.DeployDnsrecord(auth, client)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Contract address:", address.Hex())
	fmt.Println("Transaction hash:", tx.Hash().Hex())
}
