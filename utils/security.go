package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/sipatir08/cinema-ticket-api/models"
)

type Bcrypt struct{}

func (Bcrypt) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

type JWT struct {
	secret []byte
	ttl    time.Duration
}

func NewJWT(secret string, ttl time.Duration) JWT {
	return JWT{secret: []byte(secret), ttl: ttl}
}

type claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func (j JWT) Issue(userID int64, role string) (string, error) {
	now := time.Now()
	c := claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(j.secret)
}

func (j JWT) Verify(token string) (models.AuthClaims, error) {
	var c claims
	_, err := jwt.ParseWithClaims(token, &c, func(*jwt.Token) (any, error) { return j.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), // tolak alg "none" / algoritma palsu
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return models.AuthClaims{}, err
	}
	id, err := strconv.ParseInt(c.Subject, 10, 64)
	if err != nil {
		return models.AuthClaims{}, errors.New("invalid subject")
	}
	return models.AuthClaims{UserID: id, Role: c.Role}, nil
}