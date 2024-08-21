package domains

import (
	"certimate/internal/utils/app"
	"certimate/internal/utils/xtime"
	"time"

	"github.com/pocketbase/pocketbase/models"
)

type historyItem struct {
	Time    string `json:"time"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type history struct {
	Domain       string                  `json:"domain"`
	Log          map[Phase][]historyItem `json:"log"`
	Phase        Phase                   `json:"phase"`
	PhaseSuccess bool                    `json:"phaseSuccess"`
	DeployedAt   string                  `json:"deployedAt"`
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

func (a *history) record(phase Phase, msg string, err error, pass ...bool) {
	a.Phase = phase
	if len(pass) > 0 && pass[0] {
		a.PhaseSuccess = true
	}

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
		a.PhaseSuccess = false
	}

	a.Log[phase] = append(a.Log[phase], historyItem{
		Message: msg,
		Error:   errMsg,
		Time:    xtime.BeijingTimeStr(),
	})

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

	if err := app.GetApp().Dao().SaveRecord(domainRecord); err != nil {
		return err
	}

	return nil
}
