# Add Header

## 说明
通过简单的配置来设置HTTP头部的内容，目前仅设置request的头部，response以后可能会支持。

--
## 安装
使用Vulcand自带的Vbundle安装。假定目录为（/vulcand-bundle)
Tip: 最好新建一个文件夹来坐，不要再Vulcand的目录中

```
vbundle init --middleware="github.com/YunxiangHuang/vulcand-plugin/addheader"
```
修改 /vulcand-bundle/vctl/main.go
```
// log.Init([]*log.LogConfig{&log.LogConfig{Name:"console"}})
log.NewConsoleLogger(log.Config{log.Console, "console"})
```

在 /vulcand-bundle 中执行如下语句：
```
go build -o vulcand
pushd vctl/; go build -o vctl; popd
```
Tip: 如果在编译过程中提示找不到 github.com/mailgun/log 或者其他内容，请执行以下语句：
```
go get -u github.com/mailgun/vulcand
go get -u github.com/mailgun/log
```
如果还是提示缺少，请参考上面自行`go get`

---

## 配置
### -setproxyheader, -S
#### 格式
|格式|栗子🌰|备注|
|---|-----|----|
|单一Header，设置字符串|Test:test|设置Test的值为"test"|
|单一Header，设置变量*|Test:$X-Forwarded-For|设置Test的值为头部中X-Forwarded-For的值|
|多个Header|test1:test1, test2:test2||

	*：此处变量是指HTTP头部中的Key，例如：X-Forwarded-For

以上配置适用于etcd、vctl设置，具体格式请查看etcd与vctl的文档，一下是栗子🌰:
```
// vtcl
vctl addheader upsert -S='X-Forwarded-For:$REALIP_XFF, test:ohaha' -f f1 --id addheader1

etcdctl set vulcand/frontends/f1/middlewares/addheader1 '
	{
		"Id": "addheader1",
		"Priority": 1,
		"Type": "addheader",
		"Middleware": {
			"setproxyheader": "X-Forwarded-For:$REALIP_XFF, test:ohaha"
		}
	}
'

```

注意:因为部分shell将$认为是变量，因此在设置setproxyheader参数时请尽量用`''`进行包裹