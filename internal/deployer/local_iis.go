package deployer

import (
	"context"
    "encoding/pem"
    "crypto/x509"
	"fmt"
	"os"
	"os/exec"
	"runtime"
    "io/ioutil"
	"crypto/rand"
    "strings"
	"software.sslmate.com/src/go-pkcs12"
)

type LocalIISDeployer struct {
	option *DeployerOption
	infos  []string
}

func NewLocalIISDeployer(option *DeployerOption) (Deployer, error) {
	return &LocalIISDeployer{
		option: option,
		infos:  make([]string, 0),
	}, nil
}

func (d *LocalIISDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *LocalIISDeployer) GetInfo() []string {
	return []string{}
}

func (d *LocalIISDeployer) Deploy(ctx context.Context) error {
    // 解码证书
    certBlock, _ := pem.Decode([]byte(d.option.Certificate.Certificate))
	if certBlock == nil || certBlock.Type != "CERTIFICATE" {
        return fmt.Errorf("无效的证书文件")
    }
    cert, err := x509.ParseCertificate(certBlock.Bytes)
    if err != nil {
        return fmt.Errorf("解析证书失败: %w", err)
    }

    // 解码私钥
    keyBlock, _ := pem.Decode([]byte(d.option.Certificate.PrivateKey))
    key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
    if err != nil {
        return fmt.Errorf("解析私钥失败: %w", err)
    }

    // 生成PFX文件
    pfxData, err := pkcs12.Encode(rand.Reader, key, cert, nil, getDeployString(d.option.DeployConfig, "password"))
    if err != nil {
        return fmt.Errorf("生成PFX文件失败: %w", err)
    }

    // 写入PFX文件
	pfxPath := fmt.Sprintf("%s.pfx", d.GetID())
    if err := ioutil.WriteFile(pfxPath, pfxData, 0644); err != nil {
        return fmt.Errorf("写入PFX文件失败: %w", err)
    }
	
	// 获取域名
	domain := strings.Split(getDeployString(d.option.DeployConfig, "domain"),"\r\n");

	for _, value := range domain {
		// 合并 Powershell 脚本
		powershellScript := combinePowershellScript(pfxPath, getDeployString(d.option.DeployConfig, "password"), getDeployString(d.option.DeployConfig, "siteName"), value, getDeployString(d.option.DeployConfig, "ip"), getDeployString(d.option.DeployConfig, "bindingPort"))
		
		// 输出 Powershell 脚本
		powershellScriptPath := fmt.Sprintf("%s.ps1", d.GetID())
		if err := copyFile(powershellScriptPath, powershellScript); err != nil {
			return fmt.Errorf("输出 Powershell 脚本失败: %w", err)
		}

		// 执行 Powershell 脚本
		if err := execCmdLocalIIS(powershellScriptPath); err != nil {
			return err
		}
	}
	
	return nil
}

func execCmdLocalIIS(scriptPath string) error {

	// 判断操作系统
	if runtime.GOOS != "windows" {
		return fmt.Errorf("不支持的操作系统")
	}

	// 执行命令
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", scriptPath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("执行脚本失败: %w", err)
	}

	os.Remove(scriptPath)

	return nil
}

func combinePowershellScript(pfxPath string, pfxPassword string, siteName string, domain string, ip string, bindingPort string) string {
	script := depolyPowershellScript
    script = strings.Replace(script, "$pfxPath", pfxPath, -1)
	script = strings.Replace(script, "$pfxPassword", pfxPassword, -1)
    script = strings.Replace(script, "$siteName", siteName, -1)
    script = strings.Replace(script, "$domain", domain, -1)
    script = strings.Replace(script, "$ip", ip, -1)
    script = strings.Replace(script, "$bindingPort", bindingPort, -1)
    return script
}

const depolyPowershellScript string = `
# 导入证书到本地计算机的个人存储区(My)
$cert = Import-PfxCertificate -FilePath "$pfxPath" -CertStoreLocation Cert:\LocalMachine\My -Password (ConvertTo-SecureString -String "$pfxPassword" -AsPlainText -Force) -Exportable

# 获取 Thumbprint
$thumbprint = $cert.Thumbprint

# 导入 WebAdministration 模块
Import-Module WebAdministration

# 检查是否已存在 HTTPS 绑定
$existingBinding = Get-WebBinding -Name "$siteName" -Protocol "https" -Port $bindingPort -HostHeader "$domain" -ErrorAction SilentlyContinue
if (!$existingBinding) {
    # 添加新的 HTTPS 绑定
	New-WebBinding -Name "$siteName" -Protocol "https" -Port $bindingPort -IPAddress "$ip" -HostHeader "$domain"
}

# 获取绑定对象
$binding = Get-WebBinding -Name "$siteName" -Protocol "https" -Port $bindingPort -IPAddress "$ip" -HostHeader "$domain"

# 绑定 SSL 证书
$binding.AddSslCertificate($thumbprint, "My")

# 删除目录下的证书文件
Remove-Item -Path "$pfxPath" -Force
`