一个通过Cloudflare开放API接口实现的动态域名服务。

---
[English](./README.md) | 中文

# 什么是 Cloudflare 动态域名?

DNS域名解析本质上是静态的，它不能很好地处理动态IP地址。例如家庭宽带基本上都是分配的动态公网IP，也就是说过一段时间，或者重启光猫之后，公网IP就变了。

Cloudflare 提供了API接口，可以通过接口的方式去管理您的DNS域名解析。

使用该 Cloudflare 动态域名解析服务，必须在你的网络环境里运行该程序。

该程序主要做两件事：获取你当前网络的公网IP，然后自动绑定更新到域名上。

原理如下图所示：

![](images/ddns.png)

# 前提条件
- 拥有一个域名
- 拥有一个 Cloudflare 账号，免费账号即可
- [域名已托管至Cloudflare](https://www.google.com.hk/search?q=cloudflare+%E5%9F%9F%E5%90%8D%E6%89%98%E7%AE%A1)
- 域名下添加了解析记录
- 获取 Cloudflare API token

# 使用

``` bash
./ddns-cloudflare run --conf /yourfolder/config.json
```

建议创建定时任务来执行，例如crontab等。
```
*/15 * * * * ~/ddns-cloudflare run > ~/cf_ddns.log
```