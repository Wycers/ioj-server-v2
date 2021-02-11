package api

import (
	"github.com/go-resty/resty/v2"
	"github.com/infinity-oj/server-v2/pkg/models"
)

type ProblemAPI interface {
	//CreateProblem() )
	//GetProblems(c *gin.Context)
	GetProblem(problemName string) (*models.Problem, error)
	//UpdateProblem(c *gin.Context)
}

type problemService struct {
	client *resty.Client
}

func (s problemService) GetProblem(problemName string) (*models.Problem, error) {
	response := &models.Problem{}

	_, err := s.client.R().
		SetPathParams(map[string]string{
			"problemName": problemName,
		}).
		SetResult(response).
		Get("/problem/{problemName}")
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewProblemAPI(client *resty.Client) ProblemAPI {
	return &problemService{
		client: client,
	}
}
