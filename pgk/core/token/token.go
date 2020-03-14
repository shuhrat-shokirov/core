package token

import (
	"github.com/shuhrat-shokirov/jwt/pkg/cmd"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"

)

type Service struct {
	secret []byte
}

type Payload struct {
	Id    int64    `json:"id"`
	Exp   int64    `json:"exp"`
	Roles []string `json:"roles"`
}

type RequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseDTO struct {
	Token string `json:"token"`
}

var ErrInvalidLogin = errors.New("invalid password")
var ErrInvalidPassword = errors.New("invalid password")

func NewService(secret jwt.Secret) *Service {
	return &Service{secret: secret}
}

func (s *Service) Generate(context context.Context, request *RequestDTO) (response ResponseDTO, err error) {
	// TODO: Go to DB & get user by login
	hash, err := bcrypt.GenerateFromPassword([]byte("hash"), bcrypt.DefaultCost)
	log.Print(string(hash))

	err = bcrypt.CompareHashAndPassword([]byte("$2a$10$Rl.2ep5Jnq7Spj4oS9POneJiO/YhkCUzfzRVxvDp8tolE8GlXYYEi"), []byte(request.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		err = ErrInvalidPassword
		return
	}

	response.Token, err = jwt.Encode(Payload{
		Id:  1,
		Exp: time.Now().Add(time.Hour).Unix(),
		Roles: []string{"ROLE_USER"},
	}, s.secret)
	return
}
