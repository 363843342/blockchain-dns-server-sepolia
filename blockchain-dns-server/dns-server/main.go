package main

import (
	"ayls/blockchain-dns-server/dns-server/dns"
	"fmt"
	"github.com/natefinch/lumberjack"
	"log"
	"os"
)

var logger *log.Logger

const (
	logLoc = "./dns-server.log"
)

func main() {
	e, err := os.OpenFile(logLoc, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		os.Exit(1)
	}
	logger = log.New(e, "", log.Ldate|log.Ltime)
	logger.SetOutput(&lumberjack.Logger{
		Filename:   logLoc,
		MaxSize:    1,  // megabytes after which new file is created
		MaxBackups: 3,  // number of backups
		MaxAge:     28, //days
	})

	// 修复：不要使用 logger.Fatal 包装 Start 函数
	// Start 内部会启动 goroutine 监听 UDP
	_ = dns.Start(logger)

	// 阻塞主线程，防止程序退出
	s := make(chan struct{})
	fmt.Println("DNS service started successfully and entered listening....")
	<-s 
}