# Real IP Middleware

## 更新
### 2015年10月19日
	1.~~使用官方库（net）作为IPv4Addr底层~~(已经换回自己写的IPv4)
	2.修复Segament构造函数名错误的问题
	3.修改IPv4SegamentsMerge函数为IPv4Segaments结构的方法Merge

### 2015年10月16日
	1.增加命令简写（-R, -H, -W)
	2.拆分文件
	3.完善文档
    4.包含YunxiangHuang/ipv4

## 说明
Vulcand的插件，能够通过配置来预配置HTTP头部的内容。

补充说明：
	因为Vulcand在转发的时候会修改X-Forwarded-For，所以本模块只能临时取消-name参数，并将原本应该设置在XFF中的值设置到REALIP中，如需使用请再addheader模块中配置将XFF的值设定成$REALIP。

---
## 安装
使用Vulcand自带的Vbundle安装。假定目录为（/vulcand-bundle）
Tip: 最好新建一个文件夹来做，不要在vulcand目录中
```
vbundle init --middleware="github.com/YunxiangHuang/vulcand-plugin/realip"
```
修改/vulcand-bundle/vctl/main.go
```
// log.Init([]*log.LogConfig{&log.LogConfig{Name:"console"}})
log.NewConsoleLogger(log.Config{log.Console, "console"})
```
在/vulcand-bundle中执行如下语句
```
go build -o vulcand
pushd vctl/; go build -o vctl; popd
```
如果在编译过程中提示找不到 github.com/mailgun/log 或其他内容，请执行以下语句
```
go get -u github.com/mailgun/vulcand
go get -u github.com/mailgun/log
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
|起始与结束IP|`192.168.0.1-192.168.0.5`|

### -name, -N
### 临时取消，待恢复
#### 格式
字符串即可，比如:`REALIP_XFF`（这个是默认值）
#### 说明
增加此选项是为了将设置头部的操作交付给addHeader，realip仅提供配置，由addHeader的配置决定设置的具体内容（包含设置与否，设置在哪个Key等，具体请参看vulcand-plugin/addheader）

以上配置适用于使用etcd、vctl设置，具体格式请查看etcd与vctl的文档，以下是栗子：
```
// vtcl
// --recursive与--whitelist为必选选项
vtcl realip upsert --recursive=on --header="x-forwarded-for" -whitelist 8.8.8.8/24,172.168.199.1 -N realip_test -f f1 --id realip1

// etcd
etcdctl set vulcand/frontends/f1/middlewares/realip1 '
	{
		"Id":"realip1",
		"Priority":1,
		"Type":"realip",
		"Middleware": {
			“Header": "x-forwarded-for",
			"Recursive":"on",
			"Whitelist":"8.8.8.8/24, 172.168.199.1",
			"Name": "realip_test"
		}
	}'
// 注意1，etcd中realip的设置在Middleware部分，之前的部分是vulcand的配置
// 注意2，Middleware部分最好设置成首字母大写，能够避免一些奇怪的问题（比如读取不到配置）
```
