// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract InetDnsRecord {
    // 映射：域名 => 拥有者地址
    mapping(string => address) private domainOwners;
    
    // 映射：域名 => (记录类型 => 记录值)
    mapping(string => mapping(uint16 => string)) private dnsMapping;

    // 定义事件，用于通知 Go 后端清理缓存
    // 注意：domain 不建议使用 indexed，以便后端能更简单地通过 data 解析出原始字符串
    event RecordUpdated(string domain);

    // 域名注册
    function registerDomain(string calldata domain) external {
        require(domainOwners[domain] == address(0), "Domain already registered");
        domainOwners[domain] = msg.sender;

    }

    // 域名注销
    function unregisterDomain(string calldata domain) external {
        require(domainOwners[domain] == msg.sender, "Not owner");
        
        delete domainOwners[domain];
        // 清理常用 DNS 记录类型
        delete dnsMapping[domain][1];  // A
        delete dnsMapping[domain][28]; // AAAA
        delete dnsMapping[domain][5];  // CNAME
        delete dnsMapping[domain][16]; // TXT
        delete dnsMapping[domain][15]; // MX

        // 通知缓存失效
        emit RecordUpdated(domain);
    }

    // 修改/添加记录
    function addRecord(string calldata key, uint16 recType, string calldata recValue) external {
        require(domainOwners[key] == msg.sender, "You are not the owner");
        dnsMapping[key][recType] = recValue;
        
        // 通知缓存失效
        emit RecordUpdated(key);
    }

    // 查询功能
    function getRecord(string calldata key, uint16 recType) external view returns (string memory) {
        return dnsMapping[key][recType];
    }
    
    function getOwner(string calldata domain) external view returns (address) {
        return domainOwners[domain];
    }
}