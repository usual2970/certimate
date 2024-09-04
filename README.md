<p align="center">
    <a href="#" target="_blank" rel="noopener">
        <img src="https://i.imgur.com/8cRJqdr.jpeg" alt="Certimate - Your Trusted SSL Automation Partner" />
    </a>
</p>

# Certimate

做个人产品或者在小企业里负责运维的同学，会遇到要管理多个域名的情况，需要给域名申请证书。但是手动申请证书有以下缺点：

1. 😱麻烦：申请证书并部署到服务的流程虽不复杂，但也挺麻烦的，犹其是你有多个域名需要维护的时候。
2. 😭易忘：另外当前免费证书的有效期只有90天，这就要求你定期的操作，增加了工作量的同时，你也很容易忘掉续期，从而导致网站访问不了。

Certimate 就是为了解决上述问题而产生的，它具有以下特点：

1. 支持私有部署：部署方法简单，只需下载二进制文件并执行即可完成安装。
2. 数据安全：由于是私有部署，所有数据均存储在本地，不会保存在服务商的服务器上，确保数据的安全性。
3. 操作方便：通过简单的配置即可轻松申请 SSL 证书，并且在证书即将过期时自动续期，无需人工干预。

Certimate 旨在为用户提供一个安全、简便的 SSL 证书管理解决方案。使用文档请访问[https://docs.certimate.me](https://docs.certimate.me)

- [Certimate](#certimate)
  - [一、安装](#一安装)
    - [1. 二进制文件](#1-二进制文件)
    - [2. Docker 安装](#2-docker-安装)
    - [3. 源代码安装](#3-源代码安装)
  - [二、使用](#二使用)
  - [三、概念](#三概念)
    - [1. 域名](#1-域名)
    - [2. dns 服务商授权信息](#2-dns-服务商授权信息)
    - [3. 部署服务商授权信息](#3-部署服务商授权信息)
  - [四、常见问题](#四常见问题)
  - [五、许可证](#五许可证)



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


## 三、概念

Certimate 的工作流程如下：

* 用户通过 Certimate 管理页面填写申请证书的信息，包括域名、dns 服务商的授权信息、以及要部署到的服务商的授权信息。
* Certimate 向证书场商的 API 发起申请请求，获取 SSL 证书。
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

## 四、常见问题


Q: 提供saas服务吗？

> A: 不提供，目前仅支持self-hosted（私有部署）。

Q: 数据安全？

> A: 由于仅支持私有部署，各种数据都保存在用户的服务器上。另外Certimate源码也开源，二进制包及Docker镜像打包过程全部使用Github actions进行，过程透明可见，可自行审计。

Q: 自动续期证书？

> A: 已经申请的证书会在过期前10天自动续期。每天会检查一次证书是否快要过期，快要过期时会自动重新申请证书并部署到目标服务上。



## 五、许可证

Certimate 采用 MIT 许可证，详情请查看 [LICENSE](LICENSE.md) 文件。
