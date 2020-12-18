package api

import (
	"github.com/go-resty/resty/v2"
	"github.com/infinity-oj/server-v2/pkg/models"
)

type JudgementService interface {
	Create(problemId, volume string) (*models.Judgement, error)
	Query() ([]*models.Judgement, error)
}

type judgementService struct {
	client *resty.Client
}

func (s *judgementService) Query() ([]*models.Judgement, error) {
	var response []*models.Judgement

	_, err := s.client.R().
		SetResult(&response).
		Get("/judgement")
	if err != nil {
		return nil, err
	}

	return response, err
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
	return response, nil
}

func NewJudgementService(client *resty.Client) JudgementService {
	return &judgementService{
		client: client,
	}
}
