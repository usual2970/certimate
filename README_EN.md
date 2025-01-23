[中文](README.md) | [English](README_EN.md)

> [!WARNING]
> The current branch is `next`, which is the development branch for v0.3.x. It is currently unstable and should not be used in production environments.
>
> To access the previous versions, please switch to the `main` branch.

# 🔒Certimate

For individuals managing personal projects or those responsible for IT operations in small businesses who need to manage multiple domain names, applying for certificates manually comes with several drawbacks:

1. 😱Troublesome: Applying for and deploying certificates isn’t difficult, but it can be quite a hassle, especially when managing multiple domains.
2. 😭Easily forgotten: The current free certificate has a validity period of only 90 days, requiring regular renewal operations. This increases the workload and makes it easy to forget, which can result in the website becoming inaccessible.

Certimate was created to solve the above-mentioned issues and has the following features:

1. Simple operation: Automatically apply, deploy, and renew SSL certificates without any manual intervention.
2. Support for self-hosted deployment: The deployment method is simple; you only need to download the binary file and execute it. Both the binary files and Docker images are generated using GitHub Actions, ensuring a transparent process that can be audited independently.
3. Data security: Since it is a self-hosted deployment, all data is stored locally and will not be saved on the service provider’s servers, ensuring the security of the data.

Related articles:

- [Why Certimate?](https://docs.certimate.me/blog/why-certimate)
- [Introduction to Domain Variables and Deployment Authorization Groups](https://docs.certimate.me/blog/multi-deployer)

Certimate aims to provide users with a secure and user-friendly SSL certificate management solution. For usage documentation, please visit [https://docs.certimate.me](https://docs.certimate.me).

## Installation

Installing Certimate is very simple, you can choose one of the following methods for installation:

### 1. Binary File

You can download the precompiled binary files directly from the [Releases page](https://github.com/usual2970/certimate/releases), and after extracting them, execute:

```bash
./certimate serve
```

Or run the following command to automatically add a certificate to Certimate itself.

```bash
./certimate serve yourDomain
```

> [!NOTE]
> When executing the binary file on macOS, you may see a prompt saying: “Cannot open ‘certimate’ because Apple cannot check it for malicious software.” You can go to System Preferences > Security & Privacy > General, then click “Allow Anyway,” and try executing the binary file again.

### 2. Docker Installation

```bash

mkdir -p ~/.certimate && cd ~/.certimate && curl -O https://raw.githubusercontent.com/usual2970/certimate/refs/heads/main/docker/docker-compose.yml && docker compose up -d

```

### 3. Source Code Installation

```bash
git clone EMAIL:usual2970/certimate.git
cd certimate
make local.run
```

## Usage

After completing the installation steps above, you can access the Certimate management page by visiting <http://127.0.0.1:8090> in your browser.

```bash
username：admin@certimate.fun
password：1234567890
```

### Watch the video

[![Watch the video](https://i.imgur.com/4DAUKEE.gif)](https://www.youtube.com/watch?v=am_yzdfyNOE)

## List of Supported Providers

### Request certificates

The following DNS providers are supported:

<details>

<summary>[Fold/Unfold to view ...]</summary>

| Provider                                                     | Remarks                               |
| :----------------------------------------------------------- | :------------------------------------ |
| [Alibaba Cloud](https://www.alibabacloud.com/)               |                                       |
| [Tencent Cloud](https://www.tencentcloud.com/)               |                                       |
| [Huawei Cloud](https://www.huaweicloud.com/)                 |                                       |
| [Volcengine](https://www.volcengine.com/)                    |                                       |
| [AWS Route53](https://aws.amazon.com/route53/)               |                                       |
| [Azure DNS](https://azure.microsoft.com/)                    |                                       |
| [CloudFlare](https://www.cloudflare.com/)                    |                                       |
| [ClouDNS](https://www.cloudns.net//)                         |                                       |
| [GNAME](https://www.gname.com/)                              |                                       |
| [GoDaddy](https://www.godaddy.com/)                          |                                       |
| [Name.com](https://www.name.com/)                            |                                       |
| [NameSilo](https://www.namesilo.com/)                        |                                       |
| [IBM NS1 Connect](https://www.ibm.com/products/ns1-connect/) |                                       |
| [Rain Yun](https://www.rainyun.com/)                         |                                       |
| [West.cn](https://www.west.cn/)                              |                                       |
| [PowerDNS](https://www.powerdns.com/)                        |                                       |
| ACME Proxy HTTP Request                                      | Supports managing DNS by HTTP request |

</details>

### Deploy certificates

The following hosting providers are supported:

<details>

<summary>[Fold/Unfold to view ...]</summary>

| Provider                                        | Remarks                                                                          |
| :---------------------------------------------- | :------------------------------------------------------------------------------- |
| Local                                           | Supports deployment to local servers                                             |
| SSH                                             | Supports deployment to remote servers (via SSH+SFTP/SCP)                         |
| Webhook                                         | Supports deployment to Webhook                                                   |
| [Kubernetes](https://kubernetes.io/)            | Supports deployment to Kubernetes Secret                                         |
| [Alibaba Cloud](https://www.alibabacloud.com/)  | Supports deployment to Alibaba Cloud OSS, CDN, DCDN, SLB(CLB/ALB/NLB), WAF, Live |
| [Tencent Cloud](https://www.tencentcloud.com/)  | Supports deployment to Tencent Cloud COS, CDN, ECDN, EdgeOne, CLB, CSS           |
| [Baidu AI Cloud](https://intl.cloud.baidu.com/) | Supports deployment to Baidu AI CLoud CDN                                        |
| [Huawei Cloud](https://www.huaweicloud.com/)    | Supports deployment to Huawei Cloud CDN, ELB                                     |
| [Volcengine](https://www.volcengine.com/)       | Supports deployment to Volcengine TOS, CDN, DCDN, CLB, Live                      |
| [Qiniu Cloud](https://www.qiniu.com/)           | Supports deployment to Qiniu Cloud CDN, Pili                                     |
| [Doge Cloud](https://www.dogecloud.com/)        | Supports deployment to Doge Cloud CDN                                            |
| [BytePlus](https://www.byteplus.com/)           | Supports deployment to BytePlus CDN                                              |
| [UCloud](https://www.ucloud-global.com/)        | Supports deployment to UCloud US3, UCDN                                          |
| [AWS](https://aws.amazon.com/)                  | Supports deployment to AWS CloudFront                                            |
| [Edgio](https://edg.io/)                        | Supports deployment to Edgio Applications                                        |

</details>

## Concepts

The workflow of Certimate is as follows:

- Users fill in the certificate application information on the Certimate management page, including domain name, authorization information for the DNS provider, and authorization information for the service provider to deploy to.
- Certimate sends a request to the certificate vendor's API to apply for an SSL certificate.
- Certimate stores the certificate information, including the certificate content, private key, validity period, etc., and automatically renews the certificate when it is about to expire.
- Certimate sends a deployment request to the service provider's API to deploy the certificate to the service provider's servers.

This involves authorization information for the domain, DNS provider, and deployment service provider.

### 1. Domain

It involves the domain name for which the certificate is being requested.

### 2. Authorization Information for the DNS Provider

To apply for a certificate for a domain, you need to prove that the domain belongs to you. Therefore, when manually applying for a certificate, you typically need to add a TXT record to the DNS records in the domain provider's control panel.

Certimate will automatically add a TXT record for you; you only need to fill in the authorization information for your DNS provider in the Certimate backend.

For example, if you purchased the domain from Alibaba Cloud, the authorization information would be as follows:

```bash
accessKeyId: your-access-key-id
accessKeySecret: your-access-key-secret
```

If you purchased the domain from Tencent Cloud, the authorization information would be as follows:

```bash
secretId: your-secret-id
secretKey: your-secret-key
```

Notes: This authorization information requires relevant administration permissions for accessing the DNS services. Please refer to the documentations of each service provider for the specific permissions list.

### 3. Authorization Information for the Deployment Service Provider

After Certimate applies for the certificate, it will automatically deploy the certificate to your specified target, such as Alibaba Cloud CDN. At this point, you need to fill in the authorization information for Alibaba Cloud. Certimate will use the authorization information and domain name you provided to locate the corresponding CDN service and deploy the certificate to that service.

The authorization information for the deployment service provider is the same as that for the DNS provider, with the distinction that the DNS provider's authorization information is used to prove that the domain belongs to you, while the deployment service provider's authorization information is used to provide authorization for the certificate deployment.

Notes: This authorization information requires relevant administration permissions to access the target deployment services. Please refer to the documentations of each service provider for the specific permissions list.

## FAQ

Q: Do you provide SaaS services?

> A: No, we do not provide that. Currently, we only support self-hosted.

Q: Data Security?

> A: Since only self-hosted is supported, all data is stored on the user’s server. Additionally, the source code of Certimate is open-source, and the packaging process for binary files and Docker images is entirely done using GitHub Actions. This process is transparent and visible, allowing for independent auditing.

Q: Automatic Certificate Renewal?

> A: Certificates that have already been issued will be automatically renewed **10 days before expiration**. The system checks once a day to see if any certificates are nearing expiration, and if so, it will automatically reapply for the certificate and deploy it to the target service.

## Contributing

Certimate is a free and open-source project, licensed under the [MIT License](LICENSE.md). You can use it for anything you want, even offering it as a paid service to users.

You can support the development of Certimate in the following ways:

- **Submit Code**: If you find a bug or have new feature requests, and you have relevant experience, [you can submit code to us](CONTRIBUTING_EN.md).
- **Submit an Issue**: For feature suggestions or bugs, you can [submit an issue](https://github.com/usual2970/certimate/issues) to us.

Support for more service providers, UI enhancements, bug fixes, and documentation improvements are all welcome. We encourage everyone to submit pull requests (PRs).

## Disclaimer

This software is provided under the MIT License and distributed “as-is” without any warranty of any kind. The authors and contributors are not responsible for any damages or losses resulting from the use or inability to use this software, including but not limited to data loss, business interruption, or any other potential harm.

No Warranties: This software comes without any express or implied warranties, including but not limited to implied warranties of merchantability, fitness for a particular purpose, and non-infringement.

User Responsibility: By using this software, you agree to take full responsibility for any outcomes resulting from its use.

## Join the Community

- [Telegram-a new era of messaging](https://t.me/+ZXphsppxUg41YmVl)
- Wechat Group

<img src="https://i.imgur.com/zSHEoIm.png" width="400"/>

## Star History

[![Stargazers over time](https://starchart.cc/usual2970/certimate.svg?variant=adaptive)](https://starchart.cc/usual2970/certimate)
