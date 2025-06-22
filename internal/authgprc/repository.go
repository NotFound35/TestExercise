package authgprc

import (
	"awesomeProject/internal/domain/models" // <-- сюда путь к твоему пакету с User
	"database/sql"
	"errors"
	"github.com/google/uuid"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(user *models.User) (uuid.UUID, error) {
	var id uuid.UUID
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"
	err := r.db.QueryRow(query, user.Email, user.Password).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *AuthRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, email, password FROM users WHERE email = $1"
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}
	return user, nil
}
