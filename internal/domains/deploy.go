package domains

import (
	"certimate/internal/applicant"
	"certimate/internal/deployer"
	"certimate/internal/utils/app"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/pocketbase/pocketbase/models"
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
	if errs := app.GetApp().Dao().ExpandRecord(currRecord, []string{"access", "targetAccess", "group"}, nil); len(errs) > 0 {

		errList := make([]error, 0)
		for name, err := range errs {
			errList = append(errList, fmt.Errorf("展开记录失败,%s: %w", name, err))
		}
		err = errors.Join(errList...)
		app.GetApp().Logger().Error("展开记录失败", "err", err)
		history.record(checkPhase, "获取授权信息失败", &RecordInfo{Err: err})
		return err
	}
	history.record(checkPhase, "获取授权信息成功", nil)

	cert := currRecord.GetString("certificate")
	expiredAt := currRecord.GetDateTime("expiredAt").Time()

	if cert != "" && time.Until(expiredAt) > time.Hour*24*10 && currRecord.GetBool("deployed") {
		app.GetApp().Logger().Info("证书在有效期内")
		history.record(checkPhase, "证书在有效期内且已部署，跳过", &RecordInfo{
			Info: []string{fmt.Sprintf("证书有效期至 %s", expiredAt.Format("2006-01-02"))},
		}, true)
		return err
	}
	history.record(checkPhase, "检查通过", nil, true)

	// ############2.申请证书
	history.record(applyPhase, "开始申请", nil)

	if cert != "" && time.Until(expiredAt) > time.Hour*24 {
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

	for _, deployer := range deployers {
		if err = deployer.Deploy(ctx); err != nil {

			app.GetApp().Logger().Error("部署失败", "err", err)
			history.record(deployPhase, "部署失败", &RecordInfo{Err: err, Info: deployer.GetInfo()})
			return err
		}
		history.record(deployPhase, fmt.Sprintf("[%s]-部署成功", deployer.GetID()), &RecordInfo{
			Info: deployer.GetInfo(),
		}, false)

	}

	app.GetApp().Logger().Info("部署成功")
	history.record(deployPhase, "部署成功", nil, true)

	return nil
}
