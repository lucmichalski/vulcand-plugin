# Real IP Middleware

## 更新
### 2015年10月16日
	1.增加命令简写（-R, -H, -W)
	2.拆分文件
	3.完善文档
    4.包含YunxiangHuang/ipv4

## 说明
Vulcand的插件，能够通过配置来自动设置Remote_Addr的内容。

---
## 安装
使用Vulcand自带的Vbundle安装。假定目录为（/vulcand-bundle）
```
Vbundle --middleware="github.com/YunxiangHuang/vulcand-plugin/realip
```
修改/vulcand-bundle/vctl/main.go
```
// log.Init([]*log.LogConfig{&log.LogConfig{Name:"console"}})
log.NewConsoleLogger(log.Config{log.Console, "console"})
```
在/vulcand-bundle中执行如下语句
```
go build -o vulcand
push vctl/; go build -o vctl; popd
```
之后回在/vulcand-bundle中生成vulcand，在/vulcand-bundle/vctl中生成vctl
两个二进制文件需要搭配使用。

---

## 配置

### -recursive, -R
#### 可选值
`ON|OFF`
#### 说明
作用请参看下面的 -header的说明

### -header, -H
#### 可选值
`X-Forwarded-For|Remote_Addr(默认)`
#### 说明
当IP不在白名单中时，用于替换HTTP头部Remote_Addr的内容。

当此项设置为X-Forwarded-For：

	recursive为ON：从HTTP头部取得X-Forwarded-For，并从最后一个IP起，将第一个不在白名单中的IP设置为Remote_Addr的值
	recursive为OFF: 直接设置X-Forwarded-For中最后一个IP为Remote_Addr的值
当此项设置为Remote_Addr：
	会设置IP层的IP到HTTP头部的Remote_Addr

### -whitelist, -W
#### 格式
|格式|样例|
|---|---|
|单一IP地址|`192.168.0.1`|
|带子网掩码|`192.168.0.1/24`|
|起始与结束IP~~|~~`192.168.0.1-192.168.0.5`|

以上配置适用于使用etcd、vctl设置，具体格式请查看etcd与vctl的文档，以下是栗子：
```
// vtcl
// --recursive与--whitelist为必选选项
vtcl realip upsert --recursive=on --header="x-forwarded-for" -whitelist 8.8.8.8/24,172.168.199.1 -f f1 --id realip1

// etcd
etcdctl set vulcand/frontends/f1/middlwares/realip1 '
	{
		"Id":"realip1",
		"Priority":1,
		"Type":"realip",
		"Middleware": {
			“header": "x-forwarded-for",
			"recursive":"on",
			"whitelist":"8.8.8.8/24, 172.168.199.1"
		}
	}'
// 注意，etcd中realip的设置在Middleware部分，之前的部分是vulcand的配置
```
