package main

import (
    "fmt"
    "io/ioutil"
    "github.com/ethereum/go-ethereum/accounts/keystore"
)

func main() {
    // 这里的路径换成你实际的 Keystore 文件名
    keyPath := "config/keystore/UTC--2026-02-14T03-02-42.948677670Z--3e4e0b52c715bf4cbb47a5a214ba6e53bc4dece5"
    password := "169" // 你创建账户时设置的密码

    keyjson, err := ioutil.ReadFile(keyPath)
    if err != nil {
        panic(err)
    }

    key, err := keystore.DecryptKey(keyjson, password)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Private Key: %x\n", key.PrivateKey.D)
}