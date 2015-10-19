# IPv4

## 更新
### 2015年10月19日
	1.使用官方库（net）作为IPv4Addr作为底层
	2.修复Segament构造函数名称错误的问题
	3.完善文档
	
## 说明
本项目是工具项目，对Go官方库net.IP进行了加工，封装了一些方便的方法供使用（至少对于我个人而言QAQ）。

## 安装
```
go get github.com/YunxiangHuang/IPv4
```

## 使用
只需要在代码头部import即可
```
import "github.com/YunxiangHuang/IPv4"
```

---

## 类
### IPv4Addr - IPv4地址
#### 定义
```
type IPv4Addr net.IP
// type net.IP []byte
```
#### 构造方法
1. `NewIPv4AddrFromString(str string) (IPv4Addr, error)`

	通过字符串构造一个IPv4地址
	
	栗子🌰:
	```
	ip, err := NewIPv4AddrFromString("127.0.0.1")
	// ip为 127.0.0.1
	```
	
2. `NewIPv4AddrFromOther(other IPv4Addr) IPv4Addr`
	
	复制函数，复制另外一个IPv4Addr实例
	
	栗子🌰:
	```
	// ip为 127.0.0.1
	tmp := NewIPv4AddrFromOther(ip)
	// tmp为 127.0.0.1
	```
	
#### 方法
1. `Equal(other IPv4Addr) bool`

	判断IP是否与other相同
	
	栗子🌰:
	```
	// ip1为 1.1.1.1
	// ip2为 1.1.1.1
	// ip3为 1.1.1.2
	ip1.Equal(ip2)		// true
	ip1.Equal(ip3)		// false
	```
	
2. `Greater(other IPv4Addr) bool`

	判断IP是否在other之后，比如：1.1.1.2在1.1.1.1之后。
	
	栗子🌰:
	```
	// ip1为 1.1.1.1
	// ip2为 1.1.1.1
	// ip3为 1.1.1.2
	
	ip1.Greater(ip2)	// false
	ip1.Greater(ip3)	// false
	ip3.Greater(ip1)	// true
	```
	
3. `Less(other IPv4Addr) bool`

	判断IP是否在other之前，比如：1.1.1.1在2.2.2.2之前。
	
	栗子🌰:
	```
	// ip1为 1.1.1.1
	// ip2为 1.1.1.1
	// ip3为 1.1.1.2
	
	ip1.Less(ip2)		// false
	ip1.Less(ip3)		// true
	ip3.Less(ip1)		// false
	```
	
4. `String() string`
	
	格式化输出IPv4，输出栗子：1.1.1.1
	
	栗子🌰:
	```
	// ip为 1.1.1.1
	fmt.Printf(ip.String())
	// 输出结果为 1.1.1.1
	```
	
### IPv4Segament - IPv4地址段
#### 定义
```
type IPv4Segament struct {
	begin IPv4Addr
	end   IPv4Addr
}
```
#### 构造方法
1. `NewIPv4SegamentFromIPandMask(ip, mask IPv4Addr) (IPv4Segament, error)`

	通过IP与Mask来构造
	
	栗子🌰:
	```
	ip, _ := NewIPv4AddrFromString("192.168.199.1")
	mask, _ := NewIPv4AddrFromString("255.255.0.0")
	ips, err := NewIPv4SegamentFromIPandMask(ip, mask)
	// IP段为 192.168.0.0 - 192.168.255.255
	```
	
2. `NewIPv4SegamentFromIPv4Addr(begin, end IPv4Addr) (IPv4Segament, error)`

	通过起始与结束IP来构造
	
	栗子🌰:
	```
	begin, _ := NewIPv4AddrFromString("1.1.1.1")
	end, _ := NewIPv4AddrFromString("2.2.2.2")
	ips, err := NewIPv4SegamentFromIPv4Addr(begin, end)
	// IP段为 1.1.1.1 - 2.2.2.2
	```
	
3. `NewIPv4SegamentFromString(str string) (IPv4Segament, error)`
	
	通过字符串来构造
	
	栗子🌰:
	```
	ips1, err := NewIPv4SegamentFromString("1.1.1.1")
	ips2, err := NewIPv4SegamentFromString("1.1.1.1/24")
	ips3, err := NewIPv4SegamentFromString("1.1.1.1-2.2.2.2")
	
	// ips1为：1.1.1.1 - 1.1.1.1
	// ips2为：1.1.1.0 - 1.1.1.255
	// ips3为：1.1.1.1 - 2.2.2.2
	```
	
4. `NewIPv4SegamentFromOther(other IPv4Segament) IPv4Segament`
	
	复制函数，复制另外一个IPv4Segament实例
	
	栗子🌰:
	```
	// ips1为 1.1.1.0 - 1.1.1.255
	ips2 := NewIPv4SegamentFromOther(ips1)
	// ips2为 1.1.1.0 - 1.1.1.255
	```

#### 方法
1. `Equal(other IPv4Segament) bool`
	
	判断是否相同
	
	栗子🌰:
	```
	// ips1 为 1.1.1.1 - 2.2.2.2
	// ips2 为 1.1.1.1 - 2.2.2.2
	// ips3 为 1.1.1.1 - 2.2.3.3
	
	ips1.Equal(ips2)		// true
	ips1.Equal(ips3)		// false
	```

2. `Less(other IPv4Segament) bool`

	判断是否在otherIP段之前，优先比较begin的先后，再比较end的先后

	栗子🌰:
	```
	// ips1 为 1.1.1.1 - 2.2.2.2
	// ips2 为 1.1.1.1 - 2.2.2.2
	// ips3 为 1.1.1.1 - 2.2.3.3
	
	ips1.Less(ips2)			// false
	ips1.Less(ips3)			// true
	ips3.Less(ips1)			// false
	```
	
3. `Greater(other IPv4Segament) bool`

	判断是否在otherIP段之后，优先比较begin的先后，再比较end的先后
	
	栗子🌰:
	```
	// ips1 为 1.1.1.1 - 2.2.2.2
	// ips2 为 1.1.1.1 - 2.2.2.2
	// ips3 为 1.1.1.1 - 2.2.3.3
	
	ips1.Greater(ips2)		// false
	ips1.Greater(ips3)		// false
	ips3.Greater(ips1)		// true
	```
	
4. `IsInclude(ip IPv4Addr) bool `

	判断是否包含ip
	
	栗子🌰:
	```
	// ip1 为 1.1.1.1
	// ip2 为 8.8.8.8
	// ips 为 1.1.1.1 - 2.2.2.2
	
	ips.IsInclude(ip1)		// true
	ips.IsInclude(ip2)		// false
	```
5. `String() string`

	格式化输出IP段

	栗子🌰:
	```
	// ips 为 1.1.1.1 - 2.2.2.2
	fmt.Printf(ips.String())
	// 输出为 1.1.1.1-2.2.2.2
	```
	
### IPv4Segaments - IPv4地址段数组
#### 定义
```
type IPv4Segaments []IPv4Segament
```

#### 构造方法
暂无

#### sort.Interface接口实现
1. `Len() int`

	获取IPv4地址段个数
2. `Less(i, j int) bool`
	
	比较第i与第j个IPv4地址段的先后
3. `Swap(i, j int)`

	交换第i与第j个IPv4地址段

#### 方法
1. `IsInclude(ip IPv4Addr) bool`

	使用二分搜索，在许多已经有序的IPv4地址段中查询是否包含ip
	
---
	
## 其他辅助方法
1. `IPv4IntToMask(mint int) IPv4Addr, error`

	将整型转化为子网掩码，mint的范围[0, 32]

	栗子🌰:
	```
	mask1, err := IPv4IntToMask(24)
	mask2, err := IPv4IntToMask(28)
	
	// mask1 为 255.255.255.0
	// mask2 为 255.255.255.240
	```
	
2. `splitWithoutSpace(str, flag string) []string`

	基于strings.Split方法与strings.TrimSpace方法，在返回结果中，所有的字符串元素将去掉两端的空格

	栗子🌰:
	```
		res := splitWithoutSpace("1.1.1.1 ,    2.2.2.2/24", ",")
		// res 为 ["1.1.1.1", "2.2.2.2/24"]
	```