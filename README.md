[中文](README.md) | [English](README_EN.md)

# 🔒Certimate

做个人产品或在小企业负责运维的同学，需要管理多个域名，要给域名申请证书。但手动申请证书有以下缺点：

1. 😱 麻烦：申请、部署证书虽不困难，但也挺麻烦的，尤其是维护多个域名的时候。
2. 😭 易忘：当前免费证书有效期仅 90 天，这就要求定期操作，增加工作量的同时，也很容易忘掉，导致网站无法访问。

Certimate 就是为了解决上述问题而产生的，它具有以下特点：

1. 操作简单：自动申请、部署、续期 SSL 证书，全程无需人工干预。
2. 支持私有部署：部署方法简单，只需下载二进制文件执行即可。二进制文件、Docker 镜像全部用 Github Actions 生成，过程透明，可自行审计。
3. 数据安全：由于是私有部署，所有数据均存储在本地，不会保存在服务商的服务器，确保数据的安全性。

相关文章：

- [⚠️⚠️⚠️V0.2.0-第一个不向后兼容的版本](https://docs.certimate.me/blog/v0.2.0)
- [Why Certimate?](https://docs.certimate.me/blog/why-certimate)
- [域名变量及部署授权组介绍](https://docs.certimate.me/blog/multi-deployer)

Certimate 旨在为用户提供一个安全、简便的 SSL 证书管理解决方案。使用文档请访问 [https://docs.certimate.me](https://docs.certimate.me)

## 一、安装

安装 Certimate 非常简单，你可以选择以下方式之一进行安装：

### 1. 二进制文件

你可以直接从[Releases 页](https://github.com/usual2970/certimate/releases)下载预先编译好的二进制文件，解压后执行:

```bash
./certimate serve
```

或运行以下命令自动给 Certimate 自身添加证书

```bash
./certimate serve 你的域名
```

> [!NOTE]
> MacOS 在执行二进制文件时会提示：无法打开“Certimate”，因为 Apple 无法检查其是否包含恶意软件。可在“系统设置 > 隐私与安全性 > 安全性”中点击“仍然允许”，然后再次尝试执行二进制文件。

### 2. Docker 安装

```bash

mkdir -p ~/.certimate && cd ~/.certimate && curl -O https://raw.githubusercontent.com/usual2970/certimate/refs/heads/main/docker/docker-compose.yml && docker compose up -d

```

### 3. 源代码安装

```bash
git clone EMAIL:usual2970/certimate.git
cd certimate
make local.run
```

## 二、使用

执行完上述安装操作后，在浏览器中访问 `http://127.0.0.1:8090` 即可访问 Certimate 管理页面。

```bash
用户名：admin@certimate.fun
密码：1234567890
```

![usage.gif](https://i.imgur.com/zpCoLVM.gif)

## 三、支持的服务商列表

|   服务商   | 支持申请证书 | 支持部署证书 | 备注                                                              |
| :--------: | :----------: | :----------: | ----------------------------------------------------------------- |
|   阿里云   |      √       |      √       | 可签发在阿里云注册的域名；可部署到阿里云 OSS、CDN、SLB            |
|   腾讯云   |      √       |      √       | 可签发在腾讯云注册的域名；可部署到腾讯云 COS、CDN、ECDN、CLB、TEO |
| 百度智能云 |              |      √       | 可部署到百度智能云 CDN                                            |
|   华为云   |      √       |      √       | 可签发在华为云注册的域名；可部署到华为云 CDN、ELB                 |
|   七牛云   |              |      √       | 可部署到七牛云 CDN                                                |
|   多吉云   |              |      √       | 可部署到多吉云 CDN                                                |
|  火山引擎  |       √       |      √       | 可签发在火山引擎注册的域名；可部署到火山引擎 Live                     |
|    AWS     |      √       |              | 可签发在 AWS Route53 托管的域名                                   |
| CloudFlare |      √       |              | 可签发在 CloudFlare 注册的域名；CloudFlare 服务自带 SSL 证书      |
|  GoDaddy   |      √       |              | 可签发在 GoDaddy 注册的域名                                       |
|  Namesilo  |      √       |              | 可签发在 Namesilo 注册的域名                                      |
|  PowerDNS  |      √       |              | 可签发在 PowerDNS 托管的域名                                      |
| HTTP 请求  |      √       |              | 可签发允许通过 HTTP 请求修改 DNS 的域名                           |
|  本地部署  |              |      √       | 可部署到本地服务器                                                |
|    SSH     |              |      √       | 可部署到 SSH 服务器                                               |
|  Webhook   |              |      √       | 可部署时回调到 Webhook                                            |
| Kubernetes |              |      √       | 可部署到 Kubernetes Secret                                        |

## 四、系统截图

<div align="center">
<img src="https://i.imgur.com/SYjjbql.jpeg" title="Login page" width="95%"/>
<img src="https://i.imgur.com/WMVbBId.jpeg" title="Dashboard page" width="47%"/>
<img src="https://i.imgur.com/8wit3ZA.jpeg" title="Domains page" width="47%"/>
<img src="https://i.imgur.com/EWtOoJ0.jpeg" title="Accesses page" width="47%"/>
<img src="https://i.imgur.com/aaPtSW7.jpeg" title="History page" width="47%"/>
</div>

## 五、概念

Certimate 的工作流程如下：

- 用户通过 Certimate 管理页面填写申请证书的信息，包括域名、DNS 服务商的授权信息、以及要部署到的服务商的授权信息。
- Certimate 向证书厂商的 API 发起申请请求，获取 SSL 证书。
- Certimate 存储证书信息，包括证书内容、私钥、证书有效期等，并在证书即将过期时自动续期。
- Certimate 向服务商的 API 发起部署请求，将证书部署到服务商的服务器上。

这就涉及域名、DNS 服务商的授权信息、部署服务商的授权信息等。

### 1. 域名

就是要申请证书的域名。

### 2. DNS 服务商授权信息

给域名申请证书需要证明域名是你的，所以我们手动申请证书的时候一般需要在域名服务商的控制台解析记录中添加一个 TXT 域名解析记录。

Certimate 会自动添加一个 TXT 域名解析记录，你只需要在 Certimate 后台中填写你的域名服务商的授权信息即可。

比如你在阿里云购买的域名，授权信息如下：

```bash
accessKeyId: your-access-key-id
accessKeySecret: your-access-key-secret
```

在腾讯云购买的域名，授权信息如下：

```bash
secretId: your-secret-id
secretKey: your-secret-key
```

注意，此授权信息需具有访问域名及 DNS 解析的管理权限，具体的权限清单请参阅各服务商自己的技术文档。

### 3. 部署服务商授权信息

Certimate 申请证书后，会自动将证书部署到你指定的目标上，比如阿里云 CDN，Certimate 会根据你填写的授权信息及域名找到对应的 CDN 服务，并将证书部署到对应的 CDN 服务上。

部署服务商授权信息和 DNS 服务商授权信息基本一致，区别在于 DNS 服务商授权信息用于证明域名是你的，部署服务商授权信息用于提供证书部署的授权信息。

注意，此授权信息需具有访问部署目标服务的相关管理权限，具体的权限清单请参阅各服务商自己的技术文档。

## 六、常见问题

Q: 提供 SaaS 服务吗？

> A: 不提供，目前仅支持 self-hosted（私有部署）。

Q: 数据安全？

> A: 由于仅支持私有部署，各种数据都保存在用户的服务器上。另外 Certimate 源码也开源，二进制包及 Docker 镜像打包过程全部使用 Github Actions 进行，过程透明可见，可自行审计。

Q: 自动续期证书？

> A: 已经申请的证书会在**过期前 10 天**自动续期。每天会检查一次证书是否快要过期，快要过期时会自动重新申请证书并部署到目标服务上。

## 七、贡献

Certimate 是一个免费且开源的项目，采用 [MIT 开源协议](LICENSE.md)。你可以使用它做任何你想做的事，甚至把它当作一个付费服务提供给用户。

你可以通过以下方式来支持 Certimate 的开发：

- 提交代码：如果你发现了 Bug 或有新的功能需求，而你又有相关经验，可以[提交代码](CONTRIBUTING.md)给我们。
- 提交 Issue：功能建议或者 Bug 可以[提交 Issue](https://github.com/usual2970/certimate/issues) 给我们。

支持更多服务商、UI 的优化改进、Bug 修复、文档完善等，欢迎大家提交 PR。

## 八、免责声明

本软件依据 MIT 许可证（MIT License）发布，免费提供，旨在“按现状”供用户使用。作者及贡献者不对使用本软件所产生的任何直接或间接后果承担责任，包括但不限于性能下降、数据丢失、服务中断、或任何其他类型的损害。

无任何保证：本软件不提供任何明示或暗示的保证，包括但不限于对特定用途的适用性、无侵权性、商用性及可靠性的保证。

用户责任：使用本软件即表示您理解并同意承担由此产生的一切风险及责任。

## 九、加入社区

- [Telegram-a new era of messaging](https://t.me/+ZXphsppxUg41YmVl)
- 微信群聊（超 200 人需邀请入群，可先加作者好友）

<img src="https://i.imgur.com/8xwsLTA.png" width="400"/>

## 十、Star 趋势图

[![Stargazers over time](https://starchart.cc/usual2970/certimate.svg?variant=adaptive)](https://starchart.cc/usual2970/certimate)

