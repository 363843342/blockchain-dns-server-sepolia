# blockchain-dns-server-sepolia
stem from ayls/blockchain-dns-server

# Blockchain DNS Server Setup Guide

**Exported:** 4/25/2026 1:12:53  
**Link:** [https://gemini.google.com/app/2c62ddcc1f0dc017](https://gemini.google.com/app/2c62ddcc1f0dc017)

## Response:
Blockchain DNS Server (Sepolia Update)
======================================

This project is a modified version of the 
[blockchain-dns-server](https://github.com/ayls/blockchain-dns-server)
. It updates the deprecated Rinkeby network to **Sepolia**, adds domain ownership verification, implements DNS caching, and includes dual-RPC data validation for enhanced security.

* * *

🚀 Features
-----------

*   **Sepolia Integration:** Fully migrated from Rinkeby.
*   **Ownership Logic:** Smart contract now includes domain owner registration and validation.
*   **Security & Performance:** DNS service features dual RPC data comparison and a built-in DNS cache.
*   **Deployment Tools:** Includes automated scripts for smart contract deployment.

* * *

🛠 Prerequisites
----------------

### 1\. RPC Endpoints

You need at least **two** working Sepolia RPC links for data validation.

*   **Private:** 
    [Infura / MetaMask Developer](https://developer.metamask.io/)
*   **Public:** 
    [Chainlist Sepolia](https://chainlist.org/chain/11155111)

### 2\. Install Dependencies

```
sudo apt update
# Install Go
sudo apt install golang-go -y
# Install Ethereum Tools (abigen and geth)
sudo add-apt-repository ppa:ethereum/ethereum -y
sudo apt update
sudo apt install ethereum solc -y
```

### 3\. Initialize Project

```
cd /path/to/blockchain-dns-server
go mod tidy
go mod download
```

* * *

⚙️ Configuration
----------------

Edit the configuration file at `config/.env`. Ensure you have at least two different RPC providers for the gateways.

```
CONTRACTADDR=""
# RPC addresses (Use different links for PROOF and DATA)
GATEWAY="https://ethereum-sepolia-rpc.publicnode.com"
DATA_GATEWAY="https://ethereum-sepolia-rpc.publicnode.com"
PROOF_GATEWAY="https://ethereum-sepolia-rpc.publicnode.com"
EVENT_GATEWAY="wss://ethereum-sepolia-rpc.publicnode.com"

PRIVKEY=""
KEYSTOREFILE=""
KEYSTOREPASS=""
```

* * *

📖 Step-by-Step Setup
---------------------

### 1\. Generate Ethereum Account

```
cd config/
geth account new --datadir .
```

> **Note:** Save your **password**, **public address**, and the **keystore filename**. Update `KEYSTOREFILE` and `KEYSTOREPASS` in `.env`.

### 2\. Deploy Smart Contract

1.  **Extract Private Key:** Update `extract_key.go` with your keystore path and password, then run:
    ```
    go run extract_key.go
    ```
    Copy the output to the `PRIVKEY` field in `.env`.
2.  **Compile Contract:**
    ```
    cd blockchain-contract/dnsrecord
    solc --abi --bin inet-dns-record.sol -o build
    ```
3.  **Generate Go Bindings:**
    ```
    cd ..
    abigen --abi=dnsrecord/build/InetDnsRecord.abi --bin=dnsrecord/build/InetDnsRecord.bin --pkg=dnsrecord --out=dnsrecord/inet-dns-record.go
    ```
4.  **Deploy:** Ensure your address has Sepolia ETH (get some from a 
    [faucet](https://cloud.google.com/application/web3/faucet/ethereum/sepolia)
    ). Update `deploy.go` with your RPC and Private Key, then:
    ```
    cd dnsrecord/deploy
    go run deploy.go
    ```
    Save the generated **Contract Address** to `CONTRACTADDR` in `.env`.

### 3\. Register Domains & Records

```
cd registrar-client
go run main.go
```

*   **Option 1:** Register a new domain.
*   **Option 2:** Add an IP record (e.g., Type: `A`, Domain: `test.com`, IP: `1.2.3.4`).
*   **Option 3:** Verify the record.

### 4\. Start the DNS Server

```
cd dns-server
go run main.go
```

**Verify Resolution:** Open a new terminal and run:

```
nslookup test.com 127.0.0.1
```
