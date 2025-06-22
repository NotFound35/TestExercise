package authgprc

import (
	authpb "awesomeProject/internal/authgprc/proto"
	"awesomeProject/internal/domain/models"
	"awesomeProject/internal/jwt"
	"awesomeProject/internal/userservice"
	"context"
	"errors"
)

type AuthServer struct {
	authpb.UnimplementedAuthServer
	UserService *userservice.UserService // подключаем бизнес-логику
}

func NewAuthServer(userService *userservice.UserService) *AuthServer {
	return &AuthServer{
		UserService: userService,
	}
}

func (s *AuthServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	user := models.User{
		Email:    req.Email,
		Password: req.Password, // ⚠️ Здесь должна быть хешировка, если ее нет — добавим
	}

	err := s.UserService.SaveUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	return &authpb.RegisterResponse{UserId: 1}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	user, err := s.UserService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !userservice.CheckPassword(user.Password, req.Password) { // функция сравнения
		return nil, errors.New("invalid credentials")
	}

	token, err := jwt.GenerateToken(user.ID) // реальная генерация токена
	if err != nil {
		return nil, err
	}

	return &authpb.LoginResponse{Token: token}, nil
}
