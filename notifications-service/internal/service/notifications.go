package service

import (
	"context"
	"encoding/json"
	"fmt"
	"notifications/internal/models"
	"notifications/pkg/api/notifications"
	"notifications/pkg/api/report"
	"time"

	"github.com/redis/go-redis/v9"
)

type ReportClient interface {
	GetTemplate(ctx context.Context, templateID int64) (*report.Template, error)
	GetContactsByTemplate(ctx context.Context, templateID int64) ([]*report.Contact, error)
}

type NotificationsService struct {
	notifications.NotificationsServer
	reportClient ReportClient
	redisClient  *redis.Client
	msgChannels  []chan models.KafkaMsg
	batchSize    int
}

func NewNotificationsService(rc ReportClient, redis *redis.Client, channels []chan models.KafkaMsg) *NotificationsService {
	return &NotificationsService{
		reportClient: rc,
		msgChannels:  channels,
		batchSize:    20,
		redisClient:  redis,
	}
}

func (s *NotificationsService) SetupBatchSize(batchSize int) {
	s.batchSize = batchSize
}

func (s *NotificationsService) Send(ctx context.Context, req *notifications.SendRequest) (*notifications.SendResponse, error) {
	templateID := req.TemplateId

	contactsCacheKey := fmt.Sprintf("contacts:%d", templateID)
	contactsResp, err := s.getCachedContacts(ctx, templateID, contactsCacheKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get contacts: %w", err)
	}

	templateCacheKey := fmt.Sprintf("template:%d", templateID)
	templateResp, err := s.getCachedTemplate(ctx, templateID, templateCacheKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	msg := models.KafkaMsg{
		TemplateTitle:   templateResp.Title,
		TemplateMessage: templateResp.Message,
	}

	batches := SplitIntoBatches(contactsResp, s.batchSize)
	for i, batch := range batches {
		msg.ContactEmail = make([]string, len(batch))
		for j, contact := range batch {
			msg.ContactEmail[j] = contact.Email
		}
		s.msgChannels[i%len(s.msgChannels)] <- msg
	}

	return &notifications.SendResponse{
		Status: "sent",
	}, nil
}

func (s *NotificationsService) getCachedTemplate(ctx context.Context, templateID int64, cacheKey string) (*report.Template, error) {
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var template report.Template
		if err := json.Unmarshal([]byte(cached), &template); err == nil {
			return &template, nil
		}
	}
	template, err := s.reportClient.GetTemplate(ctx, templateID)
	if err != nil {
		return nil, err
	}
	jsonData, _ := json.Marshal(template)
	s.redisClient.Set(ctx, cacheKey, jsonData, time.Hour)

	return template, nil
}

func (s *NotificationsService) getCachedContacts(ctx context.Context, templateID int64, cacheKey string) ([]*report.Contact, error) {
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var contacts []*report.Contact
		if err := json.Unmarshal([]byte(cached), &contacts); err == nil {
			return contacts, nil
		}
	}

	contacts, err := s.reportClient.GetContactsByTemplate(ctx, templateID)
	if err != nil {
		return nil, err
	}

	jsonData, _ := json.Marshal(contacts)
	s.redisClient.Set(ctx, cacheKey, jsonData, time.Hour)

	return contacts, nil
}

func SplitIntoBatches(contacts []*report.Contact, batchSize int) [][]*report.Contact {
	var batches [][]*report.Contact
	for i := 0; i < len(contacts); i += batchSize {
		end := i + batchSize
		if end > len(contacts) {
			end = len(contacts)
		}

		batches = append(batches, contacts[i:end])
	}

	return batches
}
