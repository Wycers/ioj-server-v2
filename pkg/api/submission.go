package api

import (
	"github.com/go-resty/resty/v2"
	"github.com/infinity-oj/server-v2/pkg/models"
)

type SubmissionService interface {
	Create(problemId, volume string) (*models.Submission, error)
}

type service struct {
	client *resty.Client
}

func (s *service) Create(problemId, volume string) (*models.Submission, error) {

	request := map[string]interface{}{
		"problemId": problemId,
		"volume":    volume,
	}

	response := &struct {
		Submission *models.Submission `json:"submission"`
		Judgement  *models.Judgement  `json:"judgement"`
	}{}

	_, err := s.client.R().
		SetBody(request).
		SetResult(response).
		Post("/submission")
	if err != nil {
		return nil, err
	}

	//Explore response object
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

	return response.Submission, nil
}

func NewSubmissionService(client *resty.Client) SubmissionService {
	return &service{
		client: client,
	}
}
