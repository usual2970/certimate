<p align="center">
    <a href="#" target="_blank" rel="noopener">
        <img src="https://i.imgur.com/4ff4nUV.jpeg" alt="Certimate - Your Trusted SSL Automation Partner" />
    </a>
</p>

# Certimate

Certimate 是一个开源的 SSL 证书管理工具，具有以下特点：

    1.	支持私有部署：部署方法简单，只需下载二进制文件并执行即可完成安装。
    2.	数据安全：由于是私有部署，所有数据均存储在本地，不会保存在服务商的服务器上，确保数据的安全性。
    3.	操作方便：通过简单的配置即可轻松申请 SSL 证书，并且在证书即将过期时自动续期，无需人工干预。

Certimate 旨在为用户提供一个安全、简便的 SSL 证书管理解决方案。

## 安装使用

你可以直接从[Releases 页](https://github.com/usual2970/certimate/releases)下载预先编译好的二进制文件，解压后执行:

```bash
./certimate serve
```

然后在浏览器中访问 http://127.0.0.1:8090 即可访问 Certimate 管理页面。

默认账号：

```bash
用户名：admin@essay.com
密码：1234567890
```


## 许可证

Certimate 采用 MIT 许可证，详情请查看 [LICENSE](LICENSE.md) 文件。
