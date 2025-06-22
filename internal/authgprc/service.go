package authgprc

import (
	"errors"
	"github.com/google/uuid"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"awesomeProject/internal/domain/models" // импорт модели User
)

type AuthService struct {
	repo      *AuthRepository
	jwtSecret string
}

func NewAuthService(repo *AuthRepository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

func (s *AuthService) Register(email, password string) (uuid.UUID, error) {
	// Проверка — не существует ли юзер
	existingUser, _ := s.repo.GetUserByEmail(email)
	if existingUser != nil {
		return uuid.UUID{}, errors.New("user already exists")
	}

	// Хеш пароля
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.UUID{}, err
	}

	user := &models.User{
		Email:    email,
		Password: string(hash),
	}

	return s.repo.CreateUser(user)
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Генерация токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
