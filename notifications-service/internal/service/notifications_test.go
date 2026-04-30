package service_test

import (
	"context"
	"encoding/json"
	"notifications/internal/service"
	"testing"
	"time"

	mock_client "notifications/internal/client/mock"
	"notifications/internal/models"
	"notifications/pkg/api/notifications"
	"notifications/pkg/api/report"
	"notifications/pkg/logger"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupTestRedis() *redis.Client {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	return redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
}

func TestNotificationsService_Send(t *testing.T) {
	ctx, _ := logger.New(context.Background())
	redisClient := setupTestRedis()

	t.Run("success with single batch", func(t *testing.T) {
		mockClient := &mock_client.MockReportClient{
			GetTemplateFunc: func(ctx context.Context, templateID int64) (*report.Template, error) {
				return &report.Template{
					Id:      1,
					Title:   "Test Title",
					Message: "Test Message",
				}, nil
			},
			GetContactsByTemplateFunc: func(ctx context.Context, templateID int64) ([]*report.Contact, error) {
				return []*report.Contact{
					{Email: "test1@example.com"},
					{Email: "test2@example.com"},
				}, nil
			},
		}

		msgChannels := make([]chan models.KafkaMsg, 1)
		msgChannels[0] = make(chan models.KafkaMsg, 1)

		svc := service.NewNotificationsService(mockClient, redisClient, msgChannels)
		svc.SetupBatchSize(10)

		resp, err := svc.Send(ctx, &notifications.SendRequest{TemplateId: 1})
		assert.NoError(t, err)
		assert.Equal(t, "sent", resp.Status)

		msg := <-msgChannels[0]
		assert.Equal(t, "Test Title", msg.TemplateTitle)
		assert.Equal(t, "Test Message", msg.TemplateMessage)
		assert.Equal(t, []string{"test1@example.com", "test2@example.com"}, msg.ContactEmail)
		cachedTemplate, err := redisClient.Get(ctx, "template:1").Result()
		assert.NoError(t, err)
		var template report.Template
		assert.NoError(t, json.Unmarshal([]byte(cachedTemplate), &template))
		assert.Equal(t, "Test Title", template.Title)

		cachedContacts, err := redisClient.Get(ctx, "contacts:1").Result()
		assert.NoError(t, err)
		var contacts []*report.Contact
		assert.NoError(t, json.Unmarshal([]byte(cachedContacts), &contacts))
		assert.Len(t, contacts, 2)
	})

	t.Run("success with cached data", func(t *testing.T) {
		template := &report.Template{
			Id:      2,
			Title:   "Cached Title",
			Message: "Cached Message",
		}
		templateData, _ := json.Marshal(template)
		redisClient.Set(ctx, "template:2", templateData, time.Hour)

		contacts := []*report.Contact{
			{Email: "cached1@example.com"},
			{Email: "cached2@example.com"},
		}
		contactsData, _ := json.Marshal(contacts)
		redisClient.Set(ctx, "contacts:2", contactsData, time.Hour)

		mockClient := &mock_client.MockReportClient{
			GetTemplateFunc: func(ctx context.Context, templateID int64) (*report.Template, error) {
				assert.Fail(t, "Should use cached template")
				return nil, nil
			},
			GetContactsByTemplateFunc: func(ctx context.Context, templateID int64) ([]*report.Contact, error) {
				assert.Fail(t, "Should use cached contacts")
				return nil, nil
			},
		}

		msgChannels := make([]chan models.KafkaMsg, 1)
		msgChannels[0] = make(chan models.KafkaMsg, 1)

		svc := service.NewNotificationsService(mockClient, redisClient, msgChannels)
		svc.SetupBatchSize(10)

		resp, err := svc.Send(ctx, &notifications.SendRequest{TemplateId: 2})
		assert.NoError(t, err)
		assert.Equal(t, "sent", resp.Status)

		msg := <-msgChannels[0]
		assert.Equal(t, "Cached Title", msg.TemplateTitle)
		assert.Equal(t, []string{"cached1@example.com", "cached2@example.com"}, msg.ContactEmail)
	})

	t.Run("success with multiple batches", func(t *testing.T) {
		mockClient := &mock_client.MockReportClient{
			GetTemplateFunc: func(ctx context.Context, templateID int64) (*report.Template, error) {
				return &report.Template{
					Id:      3,
					Title:   "Test Title",
					Message: "Test Message",
				}, nil
			},
			GetContactsByTemplateFunc: func(ctx context.Context, templateID int64) ([]*report.Contact, error) {
				return []*report.Contact{
					{Email: "test1@example.com"},
					{Email: "test2@example.com"},
					{Email: "test3@example.com"},
					{Email: "test4@example.com"},
					{Email: "test5@example.com"},
				}, nil
			},
		}

		msgChannels := make([]chan models.KafkaMsg, 2)
		msgChannels[0] = make(chan models.KafkaMsg, 3)
		msgChannels[1] = make(chan models.KafkaMsg, 3)

		svc := service.NewNotificationsService(mockClient, redisClient, msgChannels)
		svc.SetupBatchSize(2)

		resp, err := svc.Send(ctx, &notifications.SendRequest{TemplateId: 3})
		assert.NoError(t, err)
		assert.Equal(t, "sent", resp.Status)

		var msgs []models.KafkaMsg
		for i := 0; i < 3; i++ {
			select {
			case msg := <-msgChannels[0]:
				msgs = append(msgs, msg)
			case msg := <-msgChannels[1]:
				msgs = append(msgs, msg)
			case <-time.After(100 * time.Millisecond):
				break
			}
		}

		assert.Len(t, msgs, 3)
	})
}
