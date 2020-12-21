package api

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"net/http"

	"github.com/infinity-oj/server-v2/pkg/models"
)

type AccountAPI interface {
	Create(username, password, email string) (*models.Account, error)
	Login(username, password string) error
	ResetCredential(username, oldPassword, newPassword string) error
	Test() (*models.Account, error)
}

type accountAPI struct {
	client *resty.Client
}

func (s *accountAPI) ResetCredential(username, oldPassword, newPassword string) error {

	request := map[string]interface{}{
		"username":    username,
		"oldPassword": oldPassword,
		"newPassword": newPassword,
	}

	resp, err := s.client.R().
		SetBody(request).
		Post(fmt.Sprintf("/account/%s/credential/application", username))
	if err != nil {
		return err
	}
	if resp.StatusCode() == http.StatusNoContent {
		return nil
	} else {
		return errors.New("reset password failed")
	}
}

func (s *accountAPI) Create(username, password, email string) (*models.Account, error) {

	account := &models.Account{}

	request := map[string]interface{}{
		"username": username,
		"password": password,
		"email":    email,
	}

	_, err := s.client.R().
		SetBody(request).
		SetResult(account).
		Post("/account/application")
	if err != nil {
		return nil, err
	}

	// Explore response object
	//fmt.Println("Response Info:")
	//fmt.Println("  ", resp.Request.URL)
	//fmt.Println("  Error      :", err)
	//fmt.Println("  Status Code:", resp.StatusCode())
	//fmt.Println("  Status     :", resp.Status())
	//fmt.Println("  Proto      :", resp.Proto())
	//fmt.Println("  Time       :", resp.Time())
	//fmt.Println("  Received At:", resp.ReceivedAt())
	//fmt.Println("  Body       :\n", resp)
	//fmt.Println()

	return account, nil
}

func (s *accountAPI) Login(username, password string) error {

	request := map[string]interface{}{
		"username": username,
		"password": password,
	}

	_, err := s.client.R().
		SetBody(request).
		Post("/session/principal")
	if err != nil {
		return err
	}
	return err
}

func (s *accountAPI) Test() (*models.Account, error) {

	account := &models.Account{}

	_, err := s.client.R().
		SetResult(account).
		Get("/session/principal")
	if err != nil {
		return nil, err
	}

	return account, err
}

func NewAccountAPI(client *resty.Client) AccountAPI {
	return &accountAPI{
		client: client,
	}
}
