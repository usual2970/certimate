package domains

import (
	"time"

	"github.com/pocketbase/pocketbase/models"

	"certimate/internal/applicant"
	"certimate/internal/utils/app"
	"certimate/internal/utils/xtime"
)

type historyItem struct {
	Time    string   `json:"time"`
	Message string   `json:"message"`
	Error   string   `json:"error"`
	Info    []string `json:"info"`
}

type RecordInfo struct {
	Err  error    `json:"err"`
	Info []string `json:"info"`
}

type history struct {
	Domain       string                  `json:"domain"`
	Log          map[Phase][]historyItem `json:"log"`
	Phase        Phase                   `json:"phase"`
	PhaseSuccess bool                    `json:"phaseSuccess"`
	DeployedAt   string                  `json:"deployedAt"`
	Cert         *applicant.Certificate  `json:"cert"`
	WholeSuccess bool                    `json:"wholeSuccess"`
}

func NewHistory(record *models.Record) *history {
	return &history{
		Domain:       record.Id,
		DeployedAt:   time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		Log:          make(map[Phase][]historyItem),
		Phase:        checkPhase,
		PhaseSuccess: false,
	}
}

func (a *history) record(phase Phase, msg string, info *RecordInfo, pass ...bool) {
	if info == nil {
		info = &RecordInfo{}
	}
	a.Phase = phase
	if len(pass) > 0 {
		a.PhaseSuccess = pass[0]
	}

	errMsg := ""
	if info.Err != nil {
		errMsg = info.Err.Error()
		a.PhaseSuccess = false
	}

	a.Log[phase] = append(a.Log[phase], historyItem{
		Message: msg,
		Error:   errMsg,
		Info:    info.Info,
		Time:    xtime.BeijingTimeStr(),
	})
}

func (a *history) setCert(cert *applicant.Certificate) {
	a.Cert = cert
}

func (a *history) setWholeSuccess(success bool) {
	a.WholeSuccess = success
}

func (a *history) commit() error {
	collection, err := app.GetApp().Dao().FindCollectionByNameOrId("deployments")
	if err != nil {
		return err
	}

	record := models.NewRecord(collection)

	record.Set("domain", a.Domain)
	record.Set("deployedAt", a.DeployedAt)
	record.Set("log", a.Log)
	record.Set("phase", string(a.Phase))
	record.Set("phaseSuccess", a.PhaseSuccess)
	record.Set("wholeSuccess", a.WholeSuccess)

	if err := app.GetApp().Dao().SaveRecord(record); err != nil {
		return err
	}

	domainRecord, err := app.GetApp().Dao().FindRecordById("domains", a.Domain)
	if err != nil {
		return err
	}

	domainRecord.Set("lastDeployedAt", a.DeployedAt)
	domainRecord.Set("lastDeployment", record.Id)
	domainRecord.Set("rightnow", false)
	if a.Phase == deployPhase && a.PhaseSuccess {
		domainRecord.Set("deployed", true)
	}
	cert := a.Cert
	if cert != nil {
		domainRecord.Set("certUrl", cert.CertUrl)
		domainRecord.Set("certStableUrl", cert.CertStableUrl)
		domainRecord.Set("privateKey", cert.PrivateKey)
		domainRecord.Set("certificate", cert.Certificate)
		domainRecord.Set("issuerCertificate", cert.IssuerCertificate)
		domainRecord.Set("csr", cert.Csr)
		domainRecord.Set("expiredAt", time.Now().Add(time.Hour*24*90))
	}

	if err := app.GetApp().Dao().SaveRecord(domainRecord); err != nil {
		return err
	}

	return nil
}
