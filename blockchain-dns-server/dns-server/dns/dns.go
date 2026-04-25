package dns

import (
	"context"
	"fmt"
	"golang.org/x/net/dns/dnsmessage"
	"log"
	"net"
	"sync"
	"time"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type cacheItem struct {
	answers []dnsmessage.Resource
	expiry  time.Time
}

type dnsCache struct {
	mu    sync.Mutex // 使用 Mutex 因为 get 操作现在涉及写权限（更新过期时间）
	items map[string]cacheItem
}

// 核心改进：查询后自动恢复 n 秒寿命
func (c *dnsCache) get(key string) ([]dnsmessage.Resource, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found || time.Now().After(item.expiry) {
		if found {
			delete(c.items, key)
		}
		return nil, false
	}

	// 需求实现：命中后，将有效期重置为当前时间起 300 秒
	item.expiry = time.Now().Add(300 * time.Second)
	c.items[key] = item
	return item.answers, true
}

func (c *dnsCache) set(key string, answers []dnsmessage.Resource, ttl uint32) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = cacheItem{
		answers: answers,
		expiry:  time.Now().Add(time.Duration(ttl) * time.Second),
	}
}

// 核心改进：删除包含特定域名的所有缓存项
func (c *dnsCache) removeByDomain(domain string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 遍历缓存键（格式为 "domain-type"）
	for key := range c.items {
		if len(key) >= len(domain) && key[:len(domain)] == domain {
			delete(c.items, key)
		}
	}
}

type DNSService struct {
	conn      *net.UDPConn
	dnslookup *store // 注意：此处需使用指针以访问 store 的 dataReader
	logger    *log.Logger
	cache     *dnsCache
}

type Packet struct {
	addr    net.UDPAddr
	message dnsmessage.Message
}

// 核心改进：15 秒事件监听任务
// watchEvents 启动基于 WebSocket (WSS) 的实时事件监听
func (s *DNSService) watchEvents() {
	// 1. 配置过滤规则：仅监听当前合约地址的事件
	query := ethereum.FilterQuery{
		Addresses: []common.Address{s.dnslookup.contractAddr},
	}

	// 2. 创建用于接收日志的信道
	logsChan := make(chan types.Log)

	// 3. 使用专门的 eventReader (WSS) 发起订阅
	// 注意：SubscribeFilterLogs 仅在 WSS 连接下有效
	sub, err := s.dnslookup.eventReader.SubscribeFilterLogs(context.Background(), query, logsChan)
	if err != nil {
		s.logger.Printf("❌ WSS 订阅启动失败: %v。 10秒后尝试重连...", err)
		time.AfterFunc(10*time.Second, s.watchEvents)
		return
	}

	// 4. 在独立协程中处理订阅消息
	go func() {
		// 确保退出时释放订阅资源
		defer sub.Unsubscribe()
		s.logger.Println("📡 以太坊 WSS 实时监听已就绪...")

		for {
			select {
			case err := <-sub.Err():
				// 当网络波动或节点断开时，sub.Err() 会收到信号
				s.logger.Printf("⚠️ WSS 连接异常中断: %v。正在重新发起订阅...", err)
				s.watchEvents() // 递归调用以重新订阅
				return

			case vLog := <-logsChan:
				// 5. 过滤掉因链分叉导致的回退日志 (Removed 为 true 表示该记录无效)
				if vLog.Removed {
					s.logger.Printf("检测到分叉回退，忽略该日志 (Tx: %s)", vLog.TxHash.Hex())
					continue
				}

				// 6. 解析日志中的域名信息
				// 该方法复用你 store.go 中的解析逻辑
				domain := s.parseDomainFromLog(vLog)
				
				if domain != "" {
					// 7. 执行缓存清理
					// removeByDomain 会通过 Mutex 锁安全地删除所有相关的 A/CNAME 记录
					s.cache.removeByDomain(domain)
					
					fmt.Printf("Blockchain change: the cache of [%s] has been cleaned.\nTransaction: [%s]\n", domain, vLog.TxHash.Hex())
				}
			}
		}
	}()
}

func (s *DNSService) parseDomainFromLog(vLog types.Log) string {
	// 针对 string 类型事件参数，前 64 字节通常是偏移和长度
	if len(vLog.Data) <= 64 {
		return ""
	}
	length := new(big.Int).SetBytes(vLog.Data[32:64]).Uint64()
	if uint64(len(vLog.Data[64:])) < length {
		return ""
	}
	return string(vLog.Data[64 : 64+length])
}

func (s *DNSService) Listen() {
	addr := &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 53}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		s.logger.Fatalf("端口 53 绑定失败: %v", err)
	}
	s.conn = conn
	defer s.conn.Close()

	fmt.Println("DNS service started successfully, listening to Sepolia ETH....")

	for {
		buf := make([]byte, 512)
		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		var m dnsmessage.Message
		if err := m.Unpack(buf[:n]); err != nil {
			continue
		}
		if len(m.Questions) == 0 {
			continue
		}

		go s.query(Packet{*addr, m})
	}
}

func (s *DNSService) query(p Packet) {
	q := p.message.Questions[0]
	// 缓存键使用 域名+类型 的组合
	key := fmt.Sprintf("%s-%d", q.Name.String(), q.Type)

	if val, ok := s.cache.get(key); ok {
		p.message.Answers = append(p.message.Answers, val...)
		s.sendPacket(p.message, p.addr)
		return
	}

	val, ok := s.dnslookup.getVerifiedRecord(q)
	if ok {
		p.message.Answers = append(p.message.Answers, val...)
		// 初始 TTL 设置为 300
		s.cache.set(key, val, 3600)
	} else {
		p.message.Header.RCode = dnsmessage.RCodeNameError
	}

	s.sendPacket(p.message, p.addr)
}

func (s *DNSService) sendPacket(message dnsmessage.Message, addr net.UDPAddr) {
	packed, err := message.Pack()
	if err != nil {
		return
	}
	_, _ = s.conn.WriteToUDP(packed, &addr)
}

// 在 dnsCache 结构体中添加一个清空方法
func (c *dnsCache) clearAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]cacheItem) // 直接重置 map
}
// 自动调度函数：每天 0 点执行
func (s *DNSService) scheduleMidnightCleanup() {
	now := time.Now()
	
	// 计算下一次 0 点
	// 这里的 Trick 是：设置明天的 0 点 0 分 0 秒
	next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	
	// 计算当前到下一次 0 点的间隔
	delay := next.Sub(now)

	// 使用 AfterFunc，它会在 delay 耗尽时自动在新的协程中执行
	time.AfterFunc(delay, func() {
		s.logger.Println("🕛 触发每日 0 点缓存清空任务")
		s.cache.clearAll()
		
		// 递归调用：再次安排明天的清空任务
		s.scheduleMidnightCleanup()
	})
}


func Start(logger *log.Logger) *DNSService {
	service := &DNSService{
		dnslookup: &store{logger: logger},
		logger:    logger,
		cache:     &dnsCache{items: make(map[string]cacheItem)},
	}
	service.dnslookup.init()
	service.watchEvents() // 启动异步监控协程
	service.scheduleMidnightCleanup()// // 启动 0 点定时清空任务
	go service.Listen()
	return service
}