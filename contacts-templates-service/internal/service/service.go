package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"report/internal/repository"
	"report/pkg/api/report"
	"strconv"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Service struct {
	reportRepo repository.Repository
	report.UnimplementedReportServiceServer
	cache *redis.Client
}

func NewService(rr repository.Repository, c *redis.Client) *Service {
	return &Service{
		reportRepo: rr,
		cache:      c,
	}
}

func getUserId(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, status.Error(codes.Unauthenticated, "metadata is not provided")
	}
	userId := md.Get("x-user-id")
	if len(userId) == 0 {
		return 0, status.Error(codes.Unauthenticated, "x-user-id is not provided")
	}
	id, err := strconv.ParseInt(userId[0], 10, 32)
	if err != nil {
		return 0, status.Error(codes.Unauthenticated, "x-user-id is not provided")
	}
	return id, nil
}

func (s *Service) CreateTemplate(ctx context.Context, req *report.CreateTemplateRequest) (*report.CreateTemplateResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	template := report.Template{Title: req.Title, Message: req.Message}
	templateId, err := s.reportRepo.CreateTemplate(ctx, &template, userId)
	if err != nil {
		return nil, status.Error(codes.Aborted, fmt.Sprintf("failed to create template: %v", err.Error()))
	}
	return &report.CreateTemplateResponse{
		Template: &report.Template{
			Id:      templateId,
			Title:   template.Title,
			Message: template.Message}}, nil
}

func (s *Service) GetTemplate(ctx context.Context, req *report.GetTemplateRequest) (*report.GetTemplateResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	template, err := s.reportRepo.GetTemplate(ctx, req.TemplateId, userId)
	if err != nil {
		return nil, status.Error(codes.Aborted, fmt.Sprintf("failed to get template: %v", err.Error()))
	}
	return &report.GetTemplateResponse{Template: &report.Template{Title: template.Title, Message: template.Message, Id: template.Id}}, nil
}

func (s *Service) UpdateTemplate(ctx context.Context, req *report.UpdateTemplateRequest) (*report.UpdateTemplateResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	template := report.Template{Id: req.TemplateId, Message: req.Message, Title: req.Title}
	t, err := s.reportRepo.UpdateTemplate(ctx, &template, userId)
	if err != nil {
		return nil, status.Error(codes.Aborted, fmt.Sprintf("failed to update template: %v", err.Error()))
	}
	s.cache.Conn().Del(ctx, fmt.Sprintf("template:%d", t.Id))
	s.cache.Conn().Del(ctx, fmt.Sprintf("contacts:%d", t.Id))

	return &report.UpdateTemplateResponse{Template: &report.Template{Id: t.Id, Title: t.Title, Message: t.Message}}, nil
}

func (s *Service) DeleteTemplate(ctx context.Context, req *report.DeleteTemplateRequest) (*report.DeleteTemplateResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	ok, err := s.reportRepo.DeleteTemplate(ctx, req.TemplateId, userId)
	if err != nil {
		return nil, status.Error(codes.Aborted, fmt.Sprintf("failed to delete template: %v", err.Error()))
	}
	s.cache.Conn().Del(ctx, fmt.Sprintf("template:%d", req.TemplateId))
	s.cache.Conn().Del(ctx, fmt.Sprintf("contacts:%d", req.TemplateId))
	return &report.DeleteTemplateResponse{Success: ok}, nil
}

func (s *Service) ListTemplates(ctx context.Context, req *report.ListTemplatesRequest) (*report.ListTemplatesResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	templates, err := s.reportRepo.ListTemplates(ctx, userId)
	if err != nil {
		return nil, status.Error(codes.Aborted, fmt.Sprintf("failed to get list of templates: %v", err.Error()))
	}
	return &report.ListTemplatesResponse{Templates: templates}, nil
}

func (s *Service) CreateContact(ctx context.Context, req *report.CreateContactRequest) (*report.CreateContactResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	contact := report.Contact{Email: req.Email, Name: req.Name}
	id, err := s.reportRepo.CreateContact(ctx, &contact, userId)
	if err != nil {
		return nil, status.Error(codes.Aborted, fmt.Sprintf("failed to create contact: %v", err.Error()))
	}
	return &report.CreateContactResponse{
		Contact: &report.Contact{
			Id:    id,
			Email: contact.Email,
			Name:  contact.Name}}, nil
}

func (s *Service) GetContact(ctx context.Context, req *report.GetContactRequest) (*report.GetContactResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	contact, err := s.reportRepo.GetContact(ctx, req.ContactId, userId)
	if err != nil {
		return nil, status.Error(codes.Aborted, fmt.Sprintf("failed to get contact: %v", err.Error()))
	}
	return &report.GetContactResponse{Contact: contact}, nil
}

func (s *Service) UpdateContact(ctx context.Context, req *report.UpdateContactRequest) (*report.UpdateContactResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	contact := report.Contact{Id: req.ContactId, Email: req.Email, Name: req.Name}
	c, err := s.reportRepo.UpdateContact(ctx, &contact, userId)
	if err != nil {
		return nil, status.Error(codes.Aborted, fmt.Sprintf("failed to update contact: %v", err.Error()))
	}
	return &report.UpdateContactResponse{Contact: c}, nil
}

func (s *Service) DeleteContact(ctx context.Context, req *report.DeleteContactRequest) (*report.DeleteContactResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	ok, err := s.reportRepo.DeleteContact(ctx, req.ContactId, userId)
	if err != nil {
		return nil, status.Error(codes.Aborted, fmt.Sprintf("failed to delete contact: %v", err.Error()))
	}
	return &report.DeleteContactResponse{Success: ok}, nil
}

func (s *Service) GetContactsByTemplate(ctx context.Context, req *report.GetContactsByTemplateRequest) (*report.GetContactsByTemplateResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	contacts, err := s.reportRepo.GetContactsByTemplate(ctx, req.TemplateId, userId)
	if err != nil {
		return nil, status.Error(codes.Aborted, fmt.Sprintf("failed to get list of contacts: %v", err.Error()))
	}
	return &report.GetContactsByTemplateResponse{Contacts: contacts}, nil
}

func (s *Service) AddContactsToTemplate(ctx context.Context, req *report.AddContactsToTemplateRequest) (*report.AddContactsToTemplateResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	err = s.reportRepo.AddContactsToTemplate(ctx, req.TemplateId, req.ContactsId, userId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Aborted, fmt.Sprintf("failed to add contact to template: %v", err.Error()))
	}
	s.cache.Conn().Del(ctx, fmt.Sprintf("template:%d", req.TemplateId))
	s.cache.Conn().Del(ctx, fmt.Sprintf("contacts:%d", req.TemplateId))
	return &report.AddContactsToTemplateResponse{Success: true}, nil
}

func (s *Service) GetContacts(ctx context.Context, req *report.GetContactsRequest) (*report.GetContactsResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	contacts, err := s.reportRepo.GetContacts(ctx, userId)
	if err != nil {
		return nil, status.Error(codes.Aborted, fmt.Sprintf("failed to get list of contacts: %v", err.Error()))
	}
	return &report.GetContactsResponse{Contacts: contacts}, nil
}

func HTTPToGRPCMiddleware(ctx context.Context, req *http.Request) metadata.MD {
	md := metadata.MD{}
	if uid := req.Header.Get("X-User-ID"); uid != "" {
		md.Set("x-user-id", uid)
	}
	return md
}
