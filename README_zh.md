Blockchain DNS Server (Sepolia 升级版)
===================================

[🌐 English Document](./README.md)

本项目基于 
[blockchain-dns-server](https://github.com/ayls/blockchain-dns-server)
 修改而来 。主要将原项目中已弃用的 Rinkeby 网络迁移至 **Sepolia** 测试网，并增强了安全性与功能性 。

🚀 项目改进
-------

*   **网络迁移**：将过时的 Rinkeby 区块链网络更换为 Sepolia 。
*   **权属验证**：合约新增了域名所有者注册及验证功能 。
*   **双重校验**：DNS 服务程序添加了两次 RPC 数据对比验证，确保解析结果的一致性 。
*   **性能优化**：引入了 DNS 缓存功能以提升解析速度 。
*   **便捷部署**：新增了手动部署合约的程序脚本 。

* * *

🛠 环境准备
-------

### 1\. 获取 RPC 节点

项目运行至少需要两个可用的 Sepolia RPC 链接 ：

*   **私人节点**：前往 
    [MetaMask Developer (Infura)](https://developer.metamask.io/key/active-endpoints)
     注册获取 。
*   **公共节点**：可从 
    [Chainlist](https://chainlist.org/chain/11155111)
     获取 。

### 2\. 安装必要运行库

```
sudo apt update
# 安装 Go 语言
sudo apt install golang-go -y
# 安装以太坊工具包 (需要 abigen 和 geth)
sudo add-apt-repository ppa:ethereum/ethereum -y
sudo apt update
sudo apt install ethereum solc -y
```

### 3\. 初始化项目

```
cd /<项目路径>/blockchain-dns-server
go mod tidy
go mod download
```

* * *

⚙️ 配置文件说明
---------

编辑 `config/.env` 文件，确保 `GATEWAY` 与 `PROOF_GATEWAY` 使用不同的 RPC 链接以启用验证功能 。

```
CONTRACTADDR=""
# RPC 地址
GATEWAY="https://ethereum-sepolia-rpc.publicnode.com"
DATA_GATEWAY="https://ethereum-sepolia-rpc.publicnode.com"
# PROOF_GATEWAY 与 DATA_GATEWAY 需要使用不同的 Sepolia RPC 链接，不然怎么叫验证呢？
PROOF_GATEWAY="https://ethereum-sepolia-rpc.publicnode.com" 
EVENT_GATEWAY="wss://ethereum-sepolia-rpc.publicnode.com"
PRIVKEY=""
KEYSTOREFILE=""
KEYSTOREPASS=""
```

* * *

📖 操作步骤
-------

### 1\. 生成以太坊钱包

```
cd config/
geth account new --datadir .
```

> **注意**：请务必记住输入的**密码**、生成的**公共地址**（Public address）以及**密钥文件路径** 。随后将文件名及密码填入 `.env` 的 `KEYSTOREFILE` 和 `KEYSTOREPASS` 中 。

### 2\. 部署智能合约

1.  **获取私钥**：在项目根目录/blockchain-dns-server，用文本编辑器打开并修改 `extract_key.go` 中的路径 `keyPath := "config/keystore/UTC--xxx`和密码`password := "xxx"`，运行 `go run extract_key.go` 获取私钥明文，并填入 `.env` 的 `PRIVKEY` 。
2.  **编译合约**：
    ```
    cd blockchain-contract/dnsrecord
    solc --abi --bin inet-dns-record.sol -o build
    ```
3.  **生成 Go 代理代码**：
    ```
    cd ..
    abigen --abi=dnsrecord/build/InetDnsRecord.abi --bin=dnsrecord/build/InetDnsRecord.bin --pkg=dnsrecord --out=dnsrecord/inet-dns-record.go
    ```
4.  **执行部署**：确保钱包内有 Sepolia 测试币（可从 
    [Google Faucet](https://cloud.google.com/application/web3/faucet/ethereum/sepolia)
     领取）。在/blockchain-dns-server/blockchain-contract/dnsrecord/deploy打开并编辑deploy.go，将里面`rpc :`和`privateKeyHex :`改成你自己的RPC地址和私钥，然后运行 ：
    ```
    cd dnsrecord/deploy
    go run deploy.go
    ```
    将生成的 `Contract address` 填入 `.env` 的 `CONTRACTADDR` 。

    部署成功合约以后，可以把.env中的明文的私钥`PRIVKEY`删除

### 3\. 注册域名与添加记录

```
cd registrar-client
go run main.go
```

*   输入 **1**：注册新域名。
*   输入 **2**：为域名添加记录（如：类型 `A`, 域名 `test.com`, IP `1.2.3.4`）。
*   输入 **3**：检查记录是否生效。
*   ctrl+c 退出

### 4\. 开启 DNS 服务器

```
cd dns-server
go run main.go
```

**测试解析**： 保持当前ssh终端的同时打开另一个终端窗口执行：

```
nslookup test.com 127.0.0.1
```
