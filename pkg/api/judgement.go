package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/infinity-oj/server-v2/pkg/models"
)

type JudgementAPI interface {
	Create(problemId, volume string) (*models.Judgement, error)
	QueryJudgements() ([]*models.Judgement, error)
	QueryJudgement(judgementId string) (*models.Judgement, error)
	CancelJudgement(judgementId string) (*models.Judgement, error)
}

type judgementService struct {
	client *resty.Client
}

func (s *judgementService) CancelJudgement(judgementId string) (*models.Judgement, error) {
	response := &models.Judgement{}

	url := fmt.Sprintf("/judgement/%s/cancel", judgementId)

	_, err := s.client.R().
		SetResult(response).
		Post(url)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *judgementService) QueryJudgement(judgementId string) (*models.Judgement, error) {

	response := &models.Judgement{}

	_, err := s.client.R().
		SetResult(response).
		Get("/judgement/" + judgementId)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *judgementService) QueryJudgements() ([]*models.Judgement, error) {
	var response []*models.Judgement

	_, err := s.client.R().
		SetResult(&response).
		Get("/judgement")
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *judgementService) Create(problemId, volume string) (*models.Judgement, error) {

	request := map[string]interface{}{
		"problemId": problemId,
		"volume":    volume,
	}

	response := &models.Judgement{}

	_, err := s.client.R().
		SetBody(request).
		SetResult(response).
		Post("/judgement")
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewJudgementAPI(client *resty.Client) JudgementAPI {
	return &judgementService{
		client: client,
	}
}
