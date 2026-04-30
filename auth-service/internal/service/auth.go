package service

import (
	"auth/internal/jwtToken"
	"auth/internal/models"
	"auth/internal/repository"
	. "auth/pkg/api/auth"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	AuthServer
	UsersRepo repository.Repository
}

func NewAuthService(repo repository.Repository) *AuthService {
	return &AuthService{
		UsersRepo: repo,
	}
}

func (s *AuthService) Register(ctx context.Context, request *RegisterRequest) (*RegisterResponse, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	user, err := s.UsersRepo.Register(ctx, &models.RegisterRequest{
		Email:    request.Email,
		Username: request.Username,
		Password: string(passwordHash),
	})
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "User already exists")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	token, err := jwtToken.New(*user, 24*time.Hour)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	registerResponse := &RegisterResponse{
		Token: token,
	}
	return registerResponse, nil
}

func (s *AuthService) Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {
	user, err := s.UsersRepo.Login(ctx, &models.LoginRequest{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		if errors.Is(err, repository.ErrInvalidPassword) {
			return nil, status.Error(codes.Unauthenticated, "Invalid password")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	token, err := jwtToken.New(*user, 24*time.Hour)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	loginResponse := &LoginResponse{
		Token: token,
	}
	return loginResponse, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, request *ValidateRequest) (*ValidateResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}
	tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")
	token, err := jwtToken.Verify(tokenString)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid claims")
	}
	uid, ok := claims["uid"].(float64)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid uid")
	}

	newMD := metadata.Pairs("X-User-ID", fmt.Sprintf("%v", uid))

	if err := grpc.SetHeader(ctx, newMD); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &ValidateResponse{Valid: token.Valid}, nil
}
