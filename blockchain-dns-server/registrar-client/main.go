package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"

	"ayls/blockchain-dns-server/blockchain-contract"
	"ayls/blockchain-dns-server/blockchain-contract/dnsrecord"
)

var myenv map[string]string

const (
	envLoc             = "../config/.env"
	ErrTransactionWait = "If you've just started the application, wait a while for the network to confirm your transaction."
)

func loadEnv() {
	var err error
	if myenv, err = godotenv.Read(envLoc); err != nil {
		log.Printf("Could not load env from %s: %v", envLoc, err)
	}
}

func updateEnvFile(k string, val string) {
	myenv[k] = val
	err := godotenv.Write(myenv, envLoc)
	if err != nil {
		log.Printf("Failed to update %s: %v\n", envLoc, err)
	}
}

func main() {
	loadEnv()

	ctx := context.Background()

	client, err := ethclient.Dial(myenv["GATEWAY"])
	if err != nil {
		log.Fatalf("Could not connect to Ethereum gateway: %v\n", err)
	}
	defer client.Close()

	session := NewSession(ctx)

	if myenv["CONTRACTADDR"] == "" {
		session = NewContract(session, client)
	}

	if myenv["CONTRACTADDR"] != "" {
		session = LoadContract(session, client)
	}

	for {
		fmt.Printf(
			"\nPick an option:\n" +
				"1. Register Domain (注册域名)\n" +
				"2. Set Record (设置解析记录)\n" +
				"3. Show Record (查询解析记录)\n" +
				"4. Show Owner (查询域名持有者)\n" +
				"5. Unregister Domain (注销域名)\n" +
				"6. Exit (退出)\n" +
				"7. Reset and exit (重置并退出)\n",
		)

		switch readStringStdin() {
		case "1":
			fmt.Println("Enter domain name to register:")
			domain := readStringStdin()
			registerDomain(session, domain)

		case "2":
			fmt.Println("Type in the record type (A, AAAA, CNAME, TXT, MX):")
			recType, err := parseRecordType(readStringStdin())
			if err == nil {
				fmt.Println("Type in the record name (domain):")
				recName := readStringStdin()
				fmt.Println("Type in the record value:")
				recValue := readStringStdin()
				setRecord(session, recType, recName, recValue)
			} else {
				log.Printf("%v\n", err)
			}

		case "3":
			fmt.Println("Type in the record type (A, AAAA, CNAME, TXT, MX):")
			recType, err := parseRecordType(readStringStdin())
			if err == nil {
				fmt.Println("Type in the record name (domain):")
				recName := readStringStdin()
				showRecord(session, recType, recName)
			} else {
				log.Printf("%v\n", err)
			}

		case "4":
			fmt.Println("Enter domain name to check owner:")
			domain := readStringStdin()
			showOwner(session, domain)

		case "5":
			fmt.Println("Enter domain name to unregister:")
			domain := readStringStdin()
			unregisterDomain(session, domain)

		case "6":
			fmt.Println("Bye!")
			return

		case "7":
			fmt.Println("Cleared contract address. Bye!")
			updateEnvFile("CONTRACTADDR", "")
			return

		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

// --- 核心功能补全 ---

// registerDomain 对应合约的 registerDomain
func registerDomain(session dnsrecord.DnsrecordSession, domain string) {
	tx, err := session.RegisterDomain(domain)
	if err != nil {
		log.Printf("Could not register domain: %v\n", err)
		return
	}
	fmt.Printf("Registration transaction sent! Hash: %s\n", tx.Hash().Hex())
}

// unregisterDomain 对应合约的 unregisterDomain
func unregisterDomain(session dnsrecord.DnsrecordSession, domain string) {
	tx, err := session.UnregisterDomain(domain)
	if err != nil {
		log.Printf("Could not unregister domain: %v\n", err)
		return
	}
	fmt.Printf("Unregistration transaction sent! Hash: %s\n", tx.Hash().Hex())
}

// showOwner 对应合约的 getOwner
func showOwner(session dnsrecord.DnsrecordSession, domain string) {
	owner, err := session.GetOwner(domain)
	if err != nil {
		log.Printf("Could not get owner: %v\n", err)
		return
	}
	if owner.String() == "0x0000000000000000000000000000000000000000" {
		fmt.Println("Domain is not registered.")
	} else {
		fmt.Printf("Owner of %s: %s\n", domain, owner.Hex())
	}
}

// setRecord 对应合约的 addRecord
func setRecord(session dnsrecord.DnsrecordSession, recType uint16, recName string, recValue string) {
	txSendAnswer, err := session.AddRecord(recName, recType, recValue)
	if err != nil {
		log.Printf("Could not set record in contract: %v\n", err)
		return
	}
	fmt.Printf("Record set! Please wait for tx %s to be confirmed.\n", txSendAnswer.Hash().Hex())
}

// showRecord 对应合约的 getRecord
func showRecord(session dnsrecord.DnsrecordSession, recType uint16, recName string) {
	val, err := session.GetRecord(recName, recType)
	if err != nil {
		log.Printf("Could not read record from contract: %v\n", err)
		log.Println(ErrTransactionWait)
		return
	}
	if val == "" {
		fmt.Println("No record found.")
	} else {
		fmt.Printf("Value: %s\n", val)
	}
}

// --- 工具函数 ---

func NewSession(ctx context.Context) (session dnsrecord.DnsrecordSession) {
	return contract.NewSession(ctx, "../config/keystore/"+myenv["KEYSTOREFILE"], myenv["KEYSTOREPASS"])
}

func NewContract(session dnsrecord.DnsrecordSession, client *ethclient.Client) dnsrecord.DnsrecordSession {
	if myenv["CONTRACTADDR"] != "" {
		return session
	}
	session, contractAddress := contract.NewContract(session, client)
	updateEnvFile("CONTRACTADDR", contractAddress)
	return session
}

func LoadContract(session dnsrecord.DnsrecordSession, client *ethclient.Client) dnsrecord.DnsrecordSession {
	return contract.LoadContract(session, client, myenv["CONTRACTADDR"])
}

func readStringStdin() string {
	reader := bufio.NewReader(os.Stdin)
	inputVal, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Invalid input: %v\n", err)
		return ""
	}
	return strings.TrimSpace(inputVal)
}

func parseRecordType(recTypeString string) (uint16, error) {
	switch strings.ToUpper(recTypeString) {
	case "A":
		return 1, nil
	case "AAAA":
		return 28, nil
	case "CNAME":
		return 5, nil
	case "MX":
		return 15, nil
	case "TXT":
		return 16, nil
	default:
		return 0, errors.New("Unsupported record type")
	}
}