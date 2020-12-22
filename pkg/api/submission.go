package api

import (
	"github.com/go-resty/resty/v2"
	"github.com/infinity-oj/server-v2/pkg/models"
)

type SubmissionAPI interface {
	Create(problemId, volume string) (int, *models.Submission, error)
}

type service struct {
	client *resty.Client
}

func (s *service) Create(problemId, volume string) (int, *models.Submission, error) {

	request := map[string]interface{}{
		"problemId": problemId,
		"volume":    volume,
	}

	response := &struct {
		Submission *models.Submission `json:"submission"`
		Judgement  *models.Judgement  `json:"judgement"`
	}{}

	resp, err := s.client.R().
		SetBody(request).
		SetResult(response).
		Post("/submission")
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode(), response.Submission, nil
}

func NewSubmissionAPI(client *resty.Client) SubmissionAPI {
	return &service{
		client: client,
	}
}
