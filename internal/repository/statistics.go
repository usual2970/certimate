package repository

import (
	"context"

	"github.com/certimate-go/certimate/internal/app"
	"github.com/certimate-go/certimate/internal/domain"
)

type StatisticsRepository struct{}

func NewStatisticsRepository() *StatisticsRepository {
	return &StatisticsRepository{}
}

func (r *StatisticsRepository) Get(ctx context.Context) (*domain.Statistics, error) {
	rs := &domain.Statistics{}

	// 所有证书
	certTotal := struct {
		Total int `db:"total"`
	}{}
	if err := app.GetDB().
		NewQuery("SELECT COUNT(*) AS total FROM certificate WHERE deleted = ''").
		One(&certTotal); err != nil {
		return nil, err
	}
	rs.CertificateTotal = certTotal.Total

	// 即将过期证书
	certExpireSoonTotal := struct {
		Total int `db:"total"`
	}{}
	if err := app.GetDB().
		NewQuery("SELECT COUNT(*) AS total FROM certificate WHERE expireAt > DATETIME('now') and expireAt < DATETIME('now', '+20 days') AND deleted = ''").
		One(&certExpireSoonTotal); err != nil {
		return nil, err
	}
	rs.CertificateExpireSoon = certExpireSoonTotal.Total

	// 已过期证书
	certExpiredTotal := struct {
		Total int `db:"total"`
	}{}
	if err := app.GetDB().
		NewQuery("SELECT COUNT(*) AS total FROM certificate WHERE expireAt < DATETIME('now') AND deleted = ''").
		One(&certExpiredTotal); err != nil {
		return nil, err
	}
	rs.CertificateExpired = certExpiredTotal.Total

	// 所有工作流
	workflowTotal := struct {
		Total int `db:"total"`
	}{}
	if err := app.GetDB().
		NewQuery("SELECT COUNT(*) AS total FROM workflow").
		One(&workflowTotal); err != nil {
		return nil, err
	}
	rs.WorkflowTotal = workflowTotal.Total

	// 已启用工作流
	workflowEnabledTotal := struct {
		Total int `db:"total"`
	}{}
	if err := app.GetDB().
		NewQuery("SELECT COUNT(*) AS total FROM workflow WHERE enabled IS TRUE").
		One(&workflowEnabledTotal); err != nil {
		return nil, err
	}
	rs.WorkflowEnabled = workflowEnabledTotal.Total
	rs.WorkflowDisabled = workflowTotal.Total - workflowEnabledTotal.Total

	return rs, nil
}
