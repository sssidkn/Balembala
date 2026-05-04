package service

import (
	"auth/internal/models"
	"auth/internal/repository"
	"auth/pkg/api/auth"
	"auth/pkg/logger"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestAuthService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, _ := logger.New(context.Background())

	mockRepo := repository.NewMockRepository(ctrl)
	authService := NewAuthService(mockRepo)

	tests := []struct {
		name          string
		request       *auth.RegisterRequest
		mockSetup     func()
		expectedError error
	}{
		{
			name: "successful registration",
			request: &auth.RegisterRequest{
				Email:    "email1@mail.com",
				Username: "user1",
				Password: "password123",
			},
			mockSetup: func() {
				mockRepo.EXPECT().Register(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
						assert.Equal(t, "email1@mail.com", req.Email)
						assert.Equal(t, "user1", req.Username)
						assert.NotEmpty(t, req.Password)
						assert.NotEqual(t, "password123", req.Password)

						return &models.User{
							UserID:   123,
							Email:    "email1@mail.com",
							Username: "user1",
						}, nil
					})
			},
			expectedError: nil,
		},
		{
			name: "user already exists",
			request: &auth.RegisterRequest{
				Email:    "exists@mail.com",
				Username: "user2",
				Password: "password123",
			},
			mockSetup: func() {
				mockRepo.EXPECT().Register(gomock.Any(), gomock.Any()).
					Return(nil, repository.ErrUserAlreadyExists)
			},
			expectedError: status.Error(codes.AlreadyExists, "User already exists"),
		},
		{
			name: "repository error",
			request: &auth.RegisterRequest{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "password123",
			},
			mockSetup: func() {
				mockRepo.EXPECT().Register(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("db error"))
			},
			expectedError: status.Error(codes.Internal, "db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			response, err := authService.Register(ctx, tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotEmpty(t, response.Token)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, _ := logger.New(context.Background())

	mockRepo := repository.NewMockRepository(ctrl)
	authService := NewAuthService(mockRepo)

	tests := []struct {
		name          string
		request       *auth.LoginRequest
		mockSetup     func()
		expectedError error
	}{
		{
			name: "successful login",
			request: &auth.LoginRequest{
				Email:    "email1@mail.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockRepo.EXPECT().Login(gomock.Any(), &models.LoginRequest{
					Email:    "email1@mail.com",
					Password: "password123",
				}).Return(&models.User{
					UserID:   123,
					Email:    "email1@mail.com",
					Username: "user1",
				}, nil)
			},
			expectedError: nil,
		},
		{
			name: "user not found",
			request: &auth.LoginRequest{
				Email:    "notfound@mail.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockRepo.EXPECT().Login(gomock.Any(), gomock.Any()).
					Return(nil, repository.ErrUserNotFound)
			},
			expectedError: status.Error(codes.NotFound, "User not found"),
		},
		{
			name: "invalid password",
			request: &auth.LoginRequest{
				Email:    "email1@mail.com",
				Password: "wrongpass",
			},
			mockSetup: func() {
				mockRepo.EXPECT().Login(gomock.Any(), gomock.Any()).
					Return(nil, repository.ErrInvalidPassword)
			},
			expectedError: status.Error(codes.Unauthenticated, "Invalid password"),
		},
		{
			name: "repository error",
			request: &auth.LoginRequest{
				Email:    "email1@mail.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockRepo.EXPECT().Login(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("db error"))
			},
			expectedError: status.Error(codes.Internal, "db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			response, err := authService.Login(ctx, tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotEmpty(t, response.Token)
			}
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, _ := logger.New(context.Background())

	mockRepo := repository.NewMockRepository(ctrl)
	authService := NewAuthService(mockRepo)

	tests := []struct {
		name          string
		token         string
		ctx           context.Context
		expectedError error
		expectedValid bool
	}{
		{
			name:          "missing metadata",
			token:         "",
			ctx:           ctx,
			expectedError: status.Error(codes.Unauthenticated, "missing metadata"),
			expectedValid: false,
		},
		{
			name:  "missing token",
			token: "",
			ctx: metadata.NewIncomingContext(ctx, metadata.Pairs(
				"other-header", "value",
			)),
			expectedError: status.Error(codes.Unauthenticated, "missing token"),
			expectedValid: false,
		},
		{
			name:  "invalid token",
			token: "invalid.token.here",
			ctx: metadata.NewIncomingContext(ctx, metadata.Pairs(
				"authorization", "Bearer invalid.token.here",
			)),
			expectedError: status.Error(codes.Unauthenticated, "invalid token"),
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := authService.ValidateToken(tt.ctx, &auth.ValidateRequest{})

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tt.expectedValid, response.Valid)
			}
		})
	}
}
