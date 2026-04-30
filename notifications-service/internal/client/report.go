package report

import (
	"context"
	"notifications/pkg/api/report"
	"notifications/pkg/logger"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	conn          *grpc.ClientConn
	serviceClient report.ReportServiceClient
	timeout       time.Duration
}

func NewClient(serverAddr string, opts ...ClientOption) (*Client, error) {
	cfg := &clientConfig{
		timeout: 5 * time.Second,
		dialOptions: []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	conn, err := grpc.Dial(serverAddr, cfg.dialOptions...)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:          conn,
		serviceClient: report.NewReportServiceClient(conn),
		timeout:       cfg.timeout,
	}, nil
}

type ClientOption func(*clientConfig)

type clientConfig struct {
	timeout     time.Duration
	dialOptions []grpc.DialOption
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) getLoggerFromContext(ctx context.Context) logger.Logger {
	ctx, _ = logger.New(ctx)
	return *logger.GetLoggerFromContext(ctx)
}

func (c *Client) enrichContext(ctx context.Context) context.Context {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}

func (c *Client) GetTemplate(ctx context.Context, templateID int64) (*report.Template, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	log := c.getLoggerFromContext(ctx)
	ctx = c.enrichContext(ctx)
	log.Info(ctx, "getting template", zap.Int64("template_id", templateID))

	resp, err := c.serviceClient.GetTemplate(ctx, &report.GetTemplateRequest{
		TemplateId: templateID,
	})
	if err != nil {
		log.Error(ctx, "failed to get template",
			zap.Int64("template_id", templateID),
			zap.Error(err))
		return nil, err
	}

	return resp.Template, nil
}

func (c *Client) GetContactsByTemplate(ctx context.Context, templateID int64) ([]*report.Contact, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	log := c.getLoggerFromContext(ctx)
	ctx = c.enrichContext(ctx)

	log.Info(ctx, "getting contacts by template",
		zap.Int64("template_id", templateID))

	resp, err := c.serviceClient.GetContactsByTemplate(ctx, &report.GetContactsByTemplateRequest{
		TemplateId: templateID,
	})
	if err != nil {
		log.Error(ctx, "failed to get contacts by template",
			zap.Int64("template_id", templateID),
			zap.Error(err))
		return nil, err
	}

	log.Info(ctx, "successfully got contacts",
		zap.Int("contacts_count", len(resp.Contacts)))

	return resp.Contacts, nil
}
