

# 🔒Certimate

做个人产品或在小企业负责运维的同学，需要管理多个域名，要给域名申请证书。但手动申请证书有以下缺点：

1. 😱麻烦：申请、部署证书虽不困难，但也挺麻烦的，尤其是维护多个域名的时候。
2. 😭易忘：当前免费证书有效期仅90天，这就要求定期操作，增加工作量的同时，也很容易忘掉，导致网站无法访问。

Certimate 就是为了解决上述问题而产生的，它具有以下特点：

1. 操作简单：自动申请、部署、续期 SSL 证书，全程无需人工干预。
2. 支持私有部署：部署方法简单，只需下载二进制文件执行即可。二进制文件、docker 镜像全部用 github actions 生成，过程透明，可自行审计。
3. 数据安全：由于是私有部署，所有数据均存储在本地，不会保存在服务商的服务器，确保数据的安全性。

相关文章：

* [Why Certimate?](https://docs.certimate.me/blog/why-certimate)
* [域名变量及部署授权组介绍](https://docs.certimate.me/blog/multi-deployer)


Certimate 旨在为用户提供一个安全、简便的 SSL 证书管理解决方案。使用文档请访问[https://docs.certimate.me](https://docs.certimate.me)

- [🔒Certimate](#certimate)
  - [一、安装](#一安装)
    - [1. 二进制文件](#1-二进制文件)
    - [2. Docker 安装](#2-docker-安装)
    - [3. 源代码安装](#3-源代码安装)
  - [二、使用](#二使用)
  - [三、支持的服务商列表](#三支持的服务商列表)
  - [四、系统截图](#四系统截图)
  - [五、概念](#五概念)
    - [1. 域名](#1-域名)
    - [2. dns 服务商授权信息](#2-dns-服务商授权信息)
    - [3. 部署服务商授权信息](#3-部署服务商授权信息)
  - [六、常见问题](#六常见问题)
  - [七、贡献](#七贡献)
  - [八、加入社区](#八加入社区)



## 一、安装

安装 Certimate 非常简单，你可以选择以下方式之一进行安装：

### 1. 二进制文件

你可以直接从[Releases 页](https://github.com/usual2970/certimate/releases)下载预先编译好的二进制文件，解压后执行:

```bash
./certimate serve
```

> [!NOTE]
> MacOS 在执行二进制文件时会提示：无法打开“certimate”，因为Apple无法检查其是否包含恶意软件。可在系统设置> 隐私与安全性> 安全性 中点击 "仍然允许"，然后再次尝试执行二进制文件。


### 2. Docker 安装

```bash

git clone git@github.com:usual2970/certimate.git && cd certimate/docker && docker compose up -d

```

### 3. 源代码安装

```bash
git clone EMAIL:usual2970/certimate.git
cd certimate
go run main.go serve
```


## 二、使用

执行完上述安装操作后，在浏览器中访问 `http://127.0.0.1:8090` 即可访问 Certimate 管理页面。

```bash
用户名：admin@certimate.fun
密码：1234567890
```

![usage.gif](https://i.imgur.com/zpCoLVM.gif)

## 三、支持的服务商列表

| 服务商 | 是否域名服务商 | 是否部署服务 | 备注 |
|------|------|-----|------|
| 阿里云| 是 | 是 | 支持阿里云注册的域名，支持部署到阿里云 CDN,OSS |
| 腾讯云| 是 | 是 | 支持腾讯云注册的域名，支持部署到腾讯云 CDN |
| 七牛云| 否 | 是 | 七牛云没有注册域名服务，支持部署到七牛云 CDN |
|CloudFlare| 是 | 否 | 支持 CloudFlare 注册的域名，CloudFlare 服务自带SSL证书 |
|SSH| 否 | 是 | 支持部署到 SSH 服务器 |
|WEBHOOK| 否 | 是 | 支持回调到 WEBHOOK |




## 四、系统截图

![login](https://i.imgur.com/SYjjbql.jpeg)

![dashboard](https://i.imgur.com/WMVbBId.jpeg)

![domains](https://i.imgur.com/8wit3ZA.jpeg)

![accesses](https://i.imgur.com/EWtOoJ0.jpeg)

![history](https://i.imgur.com/aaPtSW7.jpeg)


## 五、概念

Certimate 的工作流程如下：

* 用户通过 Certimate 管理页面填写申请证书的信息，包括域名、dns 服务商的授权信息、以及要部署到的服务商的授权信息。
* Certimate 向证书厂商的 API 发起申请请求，获取 SSL 证书。
* Certimate 存储证书信息，包括证书内容、私钥、证书有效期等，并在证书即将过期时自动续期。
* Certimate 向服务商的 API 发起部署请求，将证书部署到服务商的服务器上。

这就涉及域名、dns 服务商的授权信息、部署服务商的授权信息等。

### 1. 域名

就是要申请证书的域名。

### 2. dns 服务商授权信息

给域名申请证书需要证明域名是你的，所以我们手动申请证书的时候一般需要在域名服务商的控制台解析记录中添加一个 TXT 记录。

Certimate 会自动添加一个 TXT 记录，你只需要在 Certimate 后台中填写你的域名服务商的授权信息即可。

比如你在阿里云购买的域名，授权信息如下：

```bash
accessKeyId: xxx
accessKeySecret: TOKEN
```

在腾讯云购买的域名，授权信息如下：

```bash
secretId: xxx
secretKey: TOKEN
```

### 3. 部署服务商授权信息

Certimate 申请证书后，会自动将证书部署到你指定的目标上，比如阿里云 CDN 这时你需要填写阿里云的授权信息。Certimate 会根据你填写的授权信息及域名找到对应的 CDN 服务,并将证书部署到对应的 CDN 服务上。

部署服务商授权信息和 dns 服务商授权信息一致，区别在于 dns 服务商授权信息用于证明域名是你的，部署服务商授权信息用于提供证书部署的授权信息。

## 六、常见问题


Q: 提供saas服务吗？

> A: 不提供，目前仅支持self-hosted（私有部署）。

Q: 数据安全？

> A: 由于仅支持私有部署，各种数据都保存在用户的服务器上。另外Certimate源码也开源，二进制包及Docker镜像打包过程全部使用Github actions进行，过程透明可见，可自行审计。

Q: 自动续期证书？

> A: 已经申请的证书会在过期前10天自动续期。每天会检查一次证书是否快要过期，快要过期时会自动重新申请证书并部署到目标服务上。



## 七、贡献

Certimate 是一个免费且开源的项目，采用 [MIT 开源协议](LICENSE.md)。你可以使用它做任何你想做的事，甚至把它当作一个付费服务提供给用户。

你可以通过以下方式来支持 Certimate 的开发：

* 提交代码：如果你发现了 bug 或有新的功能需求，而你又有相关经验，可以提交代码给我们。
* 提交 issue：功能建议或者 bug 可以[提交 issue](https://github.com/usual2970/certimate/issues) 给我们。

支持更多服务商、UI 的优化改进、BUG 修复、文档完善等，欢迎大家提交 PR。

## 八、加入社区

* [Telegram-a new era of messaging](https://t.me/+ZXphsppxUg41YmVl) 

* 微信群聊

<img src="https://i.imgur.com/lJUfTeD.png" width="400"/>
