package dns

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"

	"ayls/blockchain-dns-server/blockchain-contract"
	"ayls/blockchain-dns-server/blockchain-contract/dnsrecord"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"golang.org/x/net/dns/dnsmessage"
)

// --- 结构体定义 ---

type StorageProof struct {
	Key   string       `json:"key"`
	Value *hexutil.Big `json:"value"` // 自动解析十六进制大整数
	Proof []string     `json:"proof"`
}

type AccountResult struct {
	AccountProof []string       `json:"accountProof"`
	StorageProof []StorageProof `json:"storageProof"`
}

type store struct {
	logger       *log.Logger
	dataReader   *ethclient.Client  // 负责常规合约方法调用 (Data Gateway)
	proofReader  *ethclient.Client  // 负责 eth_getProof 校验 (Proof Gateway)
	eventReader  *ethclient.Client  // 负责 WebSocket 订阅 (WSS Gateway) -> 新增
	session      dnsrecord.DnsrecordSession
	contractAddr common.Address
}

var (
	errTypeNotSupported = errors.New("dns type not supported")
	errIPInvalid        = errors.New("invalid IP address format")
)

const envLoc = "../config/.env"

// --- 核心校验逻辑 ---

// calculateSlot 计算 Solidity 映射项在 EVM 中的存储位置
// 逻辑: keccak256(h(k) + p) 其中 p 是 mapping 的槽位下标
func calculateSlot(name string, dnsType uint16, baseSlot int) common.Hash {
	// 1. 第一层映射: dnsMapping[domain]
	nameBytes := []byte(name)
	slotBase := common.BigToHash(big.NewInt(int64(baseSlot))) 
	h1 := crypto.Keccak256Hash(append(nameBytes, slotBase.Bytes()...))

	// 2. 第二层映射: [recType]
	typePadded := common.BigToHash(big.NewInt(int64(dnsType)))
	finalSlot := crypto.Keccak256Hash(append(typePadded.Bytes(), h1.Bytes()...))
	
	return finalSlot
}

// decodeSolidityShortString 解码 Solidity 短字符串存储格式
// 逻辑: 槽位末尾存 length*2，前部存数据。仅支持 < 31 字节的记录值
func decodeSolidityShortString(val *big.Int) string {
	if val == nil { return "" }
	
	buf := make([]byte, 32)
	valBytes := val.Bytes()
	// 将 big.Int 填充回 32 字节数组
	copy(buf[32-len(valBytes):], valBytes)

	lastByte := buf[31]
	
	// 如果最后一位是 0，说明是短字符串
	if lastByte&1 == 0 {
		length := int(lastByte) / 2
		if length <= 0 { return "" }
		if length > 31 { length = 31 } 
		return string(buf[:length])
	}
	
	return "LONG_STRING_NOT_SUPPORTED"
}

// getVerifiedRecord 执行带有 Merkle Proof 校验的查询
func (s *store) getVerifiedRecord(q dnsmessage.Question) ([]dnsmessage.Resource, bool) {
	rawName := q.Name.String()
	recName := rawName
	// 移除域名末尾的点
	if len(rawName) > 0 && rawName[len(rawName)-1] == '.' {
		recName = rawName[:len(rawName)-1]
	}

	// 步骤 1: 通过数据节点调用合约只读方法
	recValue, err := s.session.GetRecord(recName, uint16(q.Type))
	if err != nil || recValue == "" {
		return nil, false
	}

	// 步骤 2: 计算对应的存储槽位并获取 Merkle 证明 (eth_getProof)
	slot := calculateSlot(recName, uint16(q.Type), 1) // 假设映射在合约中的槽位下标为 1
	var result AccountResult
	
	err = s.proofReader.Client().Call(&result, "eth_getProof", s.contractAddr, []string{slot.Hex()}, "latest")
	if err != nil {
		s.logger.Printf("证明获取失败: %v", err)
		return nil, false
	}

	// 步骤 3: 基础验证
	if len(result.StorageProof) == 0 || len(result.StorageProof[0].Proof) == 0 {
		s.logger.Printf("验证失败: 节点未返回 Merkle 证明数据")
		return nil, false
	}

	// 步骤 4: 深度比较合约返回值与区块链状态根数据
	rawChainVal := result.StorageProof[0].Value.ToInt()
	chainValueStr := decodeSolidityShortString(rawChainVal)

	if chainValueStr != recValue {
		s.logger.Printf("❗警告：检测到数据篡改！")
		s.logger.Printf("节点返回数据: %s", recValue)
		s.logger.Printf("链上真实数据: %s", chainValueStr)
		return nil, false 
	}

	fmt.Printf("✅ Verification is consistent: [%s] -> %s \n", recName, chainValueStr)

	// 步骤 5: 转换为 DNS 资源记录
	resource, err := s.toResource(uint16(q.Type), rawName, recValue)
	if err != nil {
		return nil, false
	}

	return []dnsmessage.Resource{resource}, true
}

// --- 初始化逻辑 ---

func (s *store) init() {
	_ = godotenv.Load(envLoc)

	dataGateway := os.Getenv("DATA_GATEWAY")
	proofGateway := os.Getenv("PROOF_GATEWAY")
	eventGateway := os.Getenv("EVENT_GATEWAY") // 获取 WSS 地址
	contractAddrStr := os.Getenv("CONTRACTADDR")

	if dataGateway == "" || proofGateway == "" {
		s.logger.Fatal("配置文件缺失网关地址")
	}

	// 初始化双客户端
	dClient, err := ethclient.Dial(dataGateway)
	if err != nil {
		s.logger.Fatalf("数据节点连接失败: %v", err)
	}
	s.dataReader = dClient

	pClient, err := ethclient.Dial(proofGateway)
	if err != nil {
		s.logger.Fatalf("证明节点连接失败: %v", err)
	}
	s.proofReader = pClient
	
	eClient, err := ethclient.Dial(eventGateway)
	if err != nil {
		s.logger.Fatalf("事件监听节点 (WSS) 连接失败: %v", err)
	}
	s.eventReader = eClient

	if contractAddrStr != "" {
		s.contractAddr = common.HexToAddress(contractAddrStr)
	}

	ctx := context.Background()
	s.newSession(ctx) 
	s.loadContract(contractAddrStr)
	
	fmt.Printf("Connection ready:\n - Data gateway: %s\n - Proof gateway: %s\n", dataGateway, proofGateway)
}

func (s *store) newSession(ctx context.Context) {
	keystorePath := "../config/keystore/" + os.Getenv("KEYSTOREFILE")
	pass := os.Getenv("KEYSTOREPASS")
	s.session = contract.NewSession(ctx, keystorePath, pass) 
}

func (s *store) loadContract(addr string) {
	s.session = contract.LoadContract(s.session, s.dataReader, addr)
}

// toResource 构建标准 DNS 响应体
func (s *store) toResource(recType uint16, recName string, recValue string) (dnsmessage.Resource, error) {
	rName, _ := dnsmessage.NewName(recName)
	var rType dnsmessage.Type
	var rBody dnsmessage.ResourceBody

	switch dnsmessage.Type(recType) {
	case dnsmessage.TypeA:
		rType = dnsmessage.TypeA
		ip := net.ParseIP(recValue).To4()
		if ip == nil { return dnsmessage.Resource{}, errIPInvalid }
		rBody = &dnsmessage.AResource{A: [4]byte{ip[0], ip[1], ip[2], ip[3]}}
	case dnsmessage.TypeCNAME:
		rType = dnsmessage.TypeCNAME
		cname, _ := dnsmessage.NewName(recValue)
		rBody = &dnsmessage.CNAMEResource{CNAME: cname}
	default:
		return dnsmessage.Resource{}, errTypeNotSupported
	}

	return dnsmessage.Resource{
		Header: dnsmessage.ResourceHeader{
			Name:  rName,
			Type:  rType,
			Class: dnsmessage.ClassINET,
			TTL:   300, // 此 TTL 供客户端参考，后端缓存通过 dns.go 独立管理
		},
		Body: rBody,
	}, nil
}