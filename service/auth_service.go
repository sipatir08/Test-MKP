package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sipatir08/cinema-ticket-api/models"
)

type TokenIssuer interface {
	Issue(userID int64, role string) (string, error)
}

type PasswordHasher interface {
	Compare(hash, password string) error
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutput struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	User        *models.User `json:"user"`
}

type AuthService struct {
	users  models.UserRepository
	hasher PasswordHasher
	tokens TokenIssuer
}

func NewAuthService(users models.UserRepository, hasher PasswordHasher, tokens TokenIssuer) *AuthService {
	return &AuthService{users: users, hasher: hasher, tokens: tokens}
}

func (s *AuthService) Login(ctx context.Context, in LoginInput) (*LoginOutput, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))
	if email == "" || in.Password == "" {
		return nil, fmt.Errorf("%w: email and password are required", models.ErrValidation)
	}

	user, err := s.users.FindByEmail(ctx, email)
	if errors.Is(err, models.ErrNotFound) {
		return nil, models.ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	if err := s.hasher.Compare(user.PasswordHash, in.Password); err != nil {
		return nil, models.ErrInvalidCredentials
	}

	token, err := s.tokens.Issue(user.ID, user.Role)
	if err != nil {
		return nil, err
	}
	return &LoginOutput{AccessToken: token, TokenType: "Bearer", User: user}, nil
}