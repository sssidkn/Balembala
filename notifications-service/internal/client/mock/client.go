package mock_client

import (
	"context"
	"notifications/pkg/api/report"
)

type MockReportClient struct {
	GetContactsByTemplateFunc func(ctx context.Context, templateID int64) ([]*report.Contact, error)
	GetTemplateFunc           func(ctx context.Context, templateID int64) (*report.Template, error)
}

func (m *MockReportClient) GetContactsByTemplate(ctx context.Context, templateID int64) ([]*report.Contact, error) {
	return m.GetContactsByTemplateFunc(ctx, templateID)
}

func (m *MockReportClient) GetTemplate(ctx context.Context, templateID int64) (*report.Template, error) {
	return m.GetTemplateFunc(ctx, templateID)
}
