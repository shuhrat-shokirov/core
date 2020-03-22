package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Url string

type TokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type NewUser struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

// errors are part API
var ErrUnknown = errors.New("unknown error")
var ErrResponse = errors.New("response error")
var ErrAddNewUser = errors.New("allready is exist user by this login")

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

func (e *ErrorResponse) Error() string {
	return strings.Join(e.Errors, ", ")
}

// for errors.Is
func (e *ErrorResponse) Unwrap() error {
	return ErrResponse
}

type Client struct {
	url Url
}

func NewClient(url Url) *Client {
	return &Client{url: url}
}

type Rooms struct {
	Id     int64
	Status bool
	TimeInFour int
	TimeInMinutes int
	TimeOutFour int
	TimeOutMinutes int
	FileName string
}

func (c *Client) Login(ctx context.Context, login string, password string) (token string, err error) {
	requestData := TokenRequest{
		Username: login,
		Password: password,
	}
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("can't encode requestBody %v: %w", requestData, err)
	}
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/tokens", c.url),
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return "", fmt.Errorf("can't create request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("can't send request: %w", err)
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("can't parse response: %w", err)
	}

	switch response.StatusCode {
	case 200:
		var responseData *TokenResponse
		err = json.Unmarshal(responseBody, &responseData)
		if err != nil {
			return "", fmt.Errorf("can't decode response: %w", err)
		}
		return responseData.Token, nil
	case 400:
		var responseData *ErrorResponse
		err = json.Unmarshal(responseBody, &responseData)
		if err != nil {
			return "", fmt.Errorf("can't decode response: %w", err)
		}
		return "", responseData
	default:
		return "", ErrUnknown
	}

}

func (c *Client) Register(ctx context.Context, name, login, password string) (err error) {
	requestData := NewUser{
		Name:     name,
		Login:    login,
		Password: password,
	}
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("can't encode requestBody %v: %w", requestData, err)
	}
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/newUser", c.url),
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return fmt.Errorf("can't create request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return fmt.Errorf("can't send request: %w", err)
	}
	defer response.Body.Close()
	switch response.StatusCode {
	case 200:
		return nil
	case 400:
		return ErrAddNewUser
	default:
		return ErrUnknown
	}
}

func (c *Client) HomePage(ctx context.Context) ([]byte, error){
	response, err := http.Get("http://localhost:9999/api/rooms/list")
	if err != nil {
		return nil,errors.New("can't response")
	}
	response.Header.Set("Content-Type", "application/json")
	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("Can't read all")
	}
	return all, nil
}
