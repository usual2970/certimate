package domains

import (
	"context"
	"fmt"
	"strings"
	"time"

	"crypto/rsa"
	"crypto/ecdsa"

	"github.com/pocketbase/pocketbase/models"

	"golang.org/x/exp/slices"

	"github.com/usual2970/certimate/internal/applicant"
	"github.com/usual2970/certimate/internal/deployer"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/utils/app"

	"github.com/usual2970/certimate/internal/pkg/utils/x509"
)

type Phase string

const (
	checkPhase  Phase = "check"
	applyPhase  Phase = "apply"
	deployPhase Phase = "deploy"
)

func deploy(ctx context.Context, record *models.Record) error {
	defer func() {
		if r := recover(); r != nil {
			app.GetApp().Logger().Error("部署失败", "err", r)
		}
	}()
	var certificate *applicant.Certificate

	history := NewHistory(record)
	defer history.commit()

	// ############1.检查域名配置
	history.record(checkPhase, "开始检查", nil)

	currRecord, err := app.GetApp().Dao().FindRecordById("domains", record.Id)
	if err != nil {
		app.GetApp().Logger().Error("获取记录失败", "err", err)
		history.record(checkPhase, "获取域名配置失败", &RecordInfo{Err: err})
		return err
	}
	history.record(checkPhase, "获取记录成功", nil)

	cert := currRecord.GetString("certificate")
	expiredAt := currRecord.GetDateTime("expiredAt").Time()

	// 检查证书是否包含设置的所有域名
	changed := isCertChanged(cert, currRecord)

	if cert != "" && time.Until(expiredAt) > time.Hour*24*10 && currRecord.GetBool("deployed") && !changed {
		app.GetApp().Logger().Info("证书在有效期内")
		history.record(checkPhase, "证书在有效期内且已部署，跳过", &RecordInfo{
			Info: []string{fmt.Sprintf("证书有效期至 %s", expiredAt.Format("2006-01-02"))},
		}, true)

		// 跳过的情况也算成功
		history.setWholeSuccess(true)
		return nil
	}
	history.record(checkPhase, "检查通过", nil, true)

	// ############2.申请证书
	history.record(applyPhase, "开始申请", nil)

	if cert != "" && time.Until(expiredAt) > time.Hour*24 && !changed {
		history.record(applyPhase, "证书在有效期内，跳过", &RecordInfo{
			Info: []string{fmt.Sprintf("证书有效期至 %s", expiredAt.Format("2006-01-02"))},
		})
	} else {
		applicant, err := applicant.Get(currRecord)
		if err != nil {
			history.record(applyPhase, "获取applicant失败", &RecordInfo{Err: err})
			app.GetApp().Logger().Error("获取applicant失败", "err", err)
			return err
		}
		certificate, err = applicant.Apply()
		if err != nil {
			history.record(applyPhase, "申请证书失败", &RecordInfo{Err: err})
			app.GetApp().Logger().Error("申请证书失败", "err", err)
			return err
		}
		history.record(applyPhase, "申请证书成功", &RecordInfo{
			Info: []string{fmt.Sprintf("证书地址: %s", certificate.CertUrl)},
		})
		history.setCert(certificate)
	}

	history.record(applyPhase, "保存证书成功", nil, true)

	// ############3.部署证书
	history.record(deployPhase, "开始部署", nil, false)
	deployers, err := deployer.Gets(currRecord, certificate)
	if err != nil {
		history.record(deployPhase, "获取deployer失败", &RecordInfo{Err: err})
		app.GetApp().Logger().Error("获取deployer失败", "err", err)
		return err
	}

	// 没有部署配置,也算成功
	if len(deployers) == 0 {
		history.record(deployPhase, "没有部署配置", &RecordInfo{Info: []string{"没有部署配置"}})
		history.setWholeSuccess(true)
		return nil
	}

	for _, deployer := range deployers {
		if err = deployer.Deploy(ctx); err != nil {

			app.GetApp().Logger().Error("部署失败", "err", err)
			history.record(deployPhase, "部署失败", &RecordInfo{Err: err, Info: deployer.GetInfos()})
			return err
		}
		history.record(deployPhase, fmt.Sprintf("[%s]-部署成功", deployer.GetID()), &RecordInfo{
			Info: deployer.GetInfos(),
		}, false)

	}

	app.GetApp().Logger().Info("部署成功")
	history.record(deployPhase, "部署成功", nil, true)

	history.setWholeSuccess(true)

	return nil
}

func isCertChanged(certificate string, record *models.Record) bool {
	// 如果证书为空，直接返回false
	if certificate == "" {
		return true
	}

	// 解析证书
	cert, err := x509.ParseCertificateFromPEM(certificate)
	if err != nil {
		app.GetApp().Logger().Error("解析证书失败", "err", err)
		return true
	}

	// 遍历域名列表，检查是否都在证书中，找到第一个不存在证书中域名时提前返回false
	for _, domain := range strings.Split(record.GetString("domain"), ";") {
		if !slices.Contains(cert.DNSNames, domain) && !slices.Contains(cert.DNSNames, "*."+removeLastSubdomain(domain)) {
			return true
		}
	}

	// 解析applyConfig
	applyConfig := &domain.ApplyConfig{}
	record.UnmarshalJSONField("applyConfig", applyConfig)

	
	// 检查证书加密算法是否一致
	switch pubkey := cert.PublicKey.(type) {
	case *rsa.PublicKey:
	  bitSize := pubkey.N.BitLen()
	  switch bitSize {
		case 2048:
		  // RSA2048
		  if applyConfig.KeyAlgorithm != "" && applyConfig.KeyAlgorithm != "RSA2048" { return true }
		case 3072:
		  // RSA3072
		  if applyConfig.KeyAlgorithm != "RSA3072" { return true }
		case 4096:
		  // RSA4096
		  if applyConfig.KeyAlgorithm != "RSA4096" { return true }
		case 8192:
		  // RSA8192
		  if applyConfig.KeyAlgorithm != "RSA8192" { return true }
	  }
	case *ecdsa.PublicKey:
	  bitSize := pubkey.Curve.Params().BitSize
	  switch bitSize {
		case 256:
		  // EC256
		  if applyConfig.KeyAlgorithm != "EC256" { return true }
		case 384:
		  // EC384
		  if applyConfig.KeyAlgorithm != "EC384" { return true }
	  }
	}

	return false
}

func removeLastSubdomain(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) > 1 {
		return strings.Join(parts[1:], ".")
	}
	return domain
}
