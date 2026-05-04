package service

import (
	"context"
	"errors"
	"report/internal/repository"
	"report/pkg/api/report"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestGetUserId(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		expected    int64
		expectedErr error
	}{
		{
			name: "Success",
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.Pairs("x-user-id", "123")),
			expected:    123,
			expectedErr: nil,
		},
		{
			name:        "No metadata",
			ctx:         context.Background(),
			expected:    0,
			expectedErr: status.Error(codes.Unauthenticated, "metadata is not provided"),
		},
		{
			name: "No user id",
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.Pairs("other-header", "value")),
			expected:    0,
			expectedErr: status.Error(codes.Unauthenticated, "x-user-id is not provided"),
		},
		{
			name: "Invalid user id",
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.Pairs("x-user-id", "invalid")),
			expected:    0,
			expectedErr: status.Error(codes.Unauthenticated, "x-user-id is not provided"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getUserId(tt.ctx)
			assert.Equal(t, tt.expected, got)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestCreateTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	tests := []struct {
		name        string
		mock        func()
		ctx         context.Context
		req         *report.CreateTemplateRequest
		expected    *report.CreateTemplateResponse
		expectedErr error
	}{
		{
			name: "Success",
			mock: func() {
				mockRepo.EXPECT().CreateTemplate(gomock.Any(), &report.Template{
					Title:   "Test",
					Message: "Message",
				}, int64(123)).Return(int64(1), nil)
			},
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.Pairs("x-user-id", "123")),
			req: &report.CreateTemplateRequest{
				Title:   "Test",
				Message: "Message",
			},
			expected: &report.CreateTemplateResponse{
				Template: &report.Template{
					Id:      1,
					Title:   "Test",
					Message: "Message",
				},
			},
			expectedErr: nil,
		},
		{
			name: "Unauthenticated",
			mock: func() {},
			ctx:  context.Background(),
			req: &report.CreateTemplateRequest{
				Title:   "Test",
				Message: "Message",
			},
			expected:    nil,
			expectedErr: status.Error(codes.Unauthenticated, "metadata is not provided"),
		},
		{
			name: "Repository error",
			mock: func() {
				mockRepo.EXPECT().CreateTemplate(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int64(0), errors.New("db error"))
			},
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.Pairs("x-user-id", "123")),
			req: &report.CreateTemplateRequest{
				Title:   "Test",
				Message: "Message",
			},
			expected:    nil,
			expectedErr: status.Error(codes.Aborted, "failed to create template: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			resp, err := service.CreateTemplate(tt.ctx, tt.req)
			assert.Equal(t, tt.expected, resp)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestGetTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	tests := []struct {
		name        string
		mock        func()
		ctx         context.Context
		req         *report.GetTemplateRequest
		expected    *report.GetTemplateResponse
		expectedErr error
	}{
		{
			name: "Success",
			mock: func() {
				mockRepo.EXPECT().GetTemplate(gomock.Any(), int64(1), int64(123)).
					Return(&report.Template{
						Id:      1,
						Title:   "Test",
						Message: "Message",
					}, nil)
			},
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.Pairs("x-user-id", "123")),
			req: &report.GetTemplateRequest{
				TemplateId: 1,
			},
			expected: &report.GetTemplateResponse{
				Template: &report.Template{
					Id:      1,
					Title:   "Test",
					Message: "Message",
				},
			},
			expectedErr: nil,
		},
		{
			name: "Not found",
			mock: func() {
				mockRepo.EXPECT().GetTemplate(gomock.Any(), int64(1), int64(123)).
					Return(nil, repository.ErrNotFound)
			},
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.Pairs("x-user-id", "123")),
			req: &report.GetTemplateRequest{
				TemplateId: 1,
			},
			expected:    nil,
			expectedErr: status.Error(codes.Aborted, "failed to get template: not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			resp, err := service.GetTemplate(tt.ctx, tt.req)
			assert.Equal(t, tt.expected, resp)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestUpdateTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	tests := []struct {
		name        string
		mock        func()
		ctx         context.Context
		req         *report.UpdateTemplateRequest
		expected    *report.UpdateTemplateResponse
		expectedErr error
	}{
		{
			name: "Success",
			mock: func() {
				mockRepo.EXPECT().UpdateTemplate(gomock.Any(), &report.Template{
					Id:      1,
					Title:   "Updated",
					Message: "Updated message",
				}, int64(123)).
					Return(&report.Template{
						Id:      1,
						Title:   "Updated",
						Message: "Updated message",
					}, nil)
			},
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.Pairs("x-user-id", "123")),
			req: &report.UpdateTemplateRequest{
				TemplateId: 1,
				Title:      "Updated",
				Message:    "Updated message",
			},
			expected: &report.UpdateTemplateResponse{
				Template: &report.Template{
					Id:      1,
					Title:   "Updated",
					Message: "Updated message",
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			resp, err := service.UpdateTemplate(tt.ctx, tt.req)
			assert.Equal(t, tt.expected, resp)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestDeleteTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	tests := []struct {
		name        string
		mock        func()
		ctx         context.Context
		req         *report.DeleteTemplateRequest
		expected    *report.DeleteTemplateResponse
		expectedErr error
	}{
		{
			name: "Success",
			mock: func() {
				mockRepo.EXPECT().DeleteTemplate(gomock.Any(), int64(1), int64(123)).
					Return(true, nil)
			},
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.Pairs("x-user-id", "123")),
			req: &report.DeleteTemplateRequest{
				TemplateId: 1,
			},
			expected: &report.DeleteTemplateResponse{
				Success: true,
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			resp, err := service.DeleteTemplate(tt.ctx, tt.req)
			assert.Equal(t, tt.expected, resp)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestListTemplates(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	tests := []struct {
		name        string
		mock        func()
		ctx         context.Context
		req         *report.ListTemplatesRequest
		expected    *report.ListTemplatesResponse
		expectedErr error
	}{
		{
			name: "Success",
			mock: func() {
				mockRepo.EXPECT().ListTemplates(gomock.Any(), int64(123)).
					Return([]*report.Template{
						{Id: 1, Title: "Test 1", Message: "Message 1"},
						{Id: 2, Title: "Test 2", Message: "Message 2"},
					}, nil)
			},
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.Pairs("x-user-id", "123")),
			req: &report.ListTemplatesRequest{},
			expected: &report.ListTemplatesResponse{
				Templates: []*report.Template{
					{Id: 1, Title: "Test 1", Message: "Message 1"},
					{Id: 2, Title: "Test 2", Message: "Message 2"},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			resp, err := service.ListTemplates(tt.ctx, tt.req)
			assert.Equal(t, tt.expected, resp)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestCreateContact(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	tests := []struct {
		name        string
		mock        func()
		ctx         context.Context
		req         *report.CreateContactRequest
		expected    *report.CreateContactResponse
		expectedErr error
	}{
		{
			name: "Success",
			mock: func() {
				mockRepo.EXPECT().CreateContact(gomock.Any(), &report.Contact{
					Name:  "John",
					Email: "john@example.com",
				}, int64(123)).Return(int64(1), nil)
			},
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.Pairs("x-user-id", "123")),
			req: &report.CreateContactRequest{
				Name:  "John",
				Email: "john@example.com",
			},
			expected: &report.CreateContactResponse{
				Contact: &report.Contact{
					Id:    1,
					Name:  "John",
					Email: "john@example.com",
				},
			},
			expectedErr: nil,
		},
		{
			name: "Contact already exists",
			mock: func() {
				mockRepo.EXPECT().CreateContact(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int64(0), repository.ErrContactAlreadyExists)
			},
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.Pairs("x-user-id", "123")),
			req: &report.CreateContactRequest{
				Name:  "John",
				Email: "john@example.com",
			},
			expected:    nil,
			expectedErr: status.Error(codes.Aborted, "failed to create contact: contact already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			resp, err := service.CreateContact(tt.ctx, tt.req)
			assert.Equal(t, tt.expected, resp)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestAddContactsToTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	tests := []struct {
		name        string
		mock        func()
		ctx         context.Context
		req         *report.AddContactsToTemplateRequest
		expected    *report.AddContactsToTemplateResponse
		expectedErr error
	}{
		{
			name: "Success",
			mock: func() {
				mockRepo.EXPECT().AddContactsToTemplate(gomock.Any(), int64(1), []int64{1, 2}, int64(123)).
					Return(nil)
			},
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.Pairs("x-user-id", "123")),
			req: &report.AddContactsToTemplateRequest{
				TemplateId: 1,
				ContactsId: []int64{1, 2},
			},
			expected: &report.AddContactsToTemplateResponse{
				Success: true,
			},
			expectedErr: nil,
		},
		{
			name: "Not found",
			mock: func() {
				mockRepo.EXPECT().AddContactsToTemplate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(repository.ErrNotFound)
			},
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.Pairs("x-user-id", "123")),
			req: &report.AddContactsToTemplateRequest{
				TemplateId: 1,
				ContactsId: []int64{1, 2},
			},
			expected:    nil,
			expectedErr: status.Error(codes.NotFound, "not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			resp, err := service.AddContactsToTemplate(tt.ctx, tt.req)
			assert.Equal(t, tt.expected, resp)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
