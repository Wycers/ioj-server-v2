package services

import (
	"errors"
	"fmt"

	"go.uber.org/zap"

	problemRepository "github.com/infinity-oj/server-v2/internal/app/problems/repositories"
	"github.com/infinity-oj/server-v2/internal/app/submissions/repositories"
	"github.com/infinity-oj/server-v2/internal/pkg/models"
)

type SubmissionsService interface {
	Create(submitterID uint64, problemName string, userSpace string) (s *models.Submission, err error)
	GetSubmission(submissionId string) (s *models.Submission, err error)
	GetSubmissions(problemId string, page, pageSize int) (res []*models.Submission, err error)

	// Following are need to rearrange
	DeliverJudgement(element *repositories.Process) error
	DispatchJudgement(submissionId string) error
	ReturnJudgement(judgementId string, outputs [][]byte) error
}

type DefaultSubmissionService struct {
	logger               *zap.Logger
	SubmissionRepository repositories.Repository
	ProblemRepository    problemRepository.Repository

	processMap map[string]string
	idMap      map[string]int
}

func (d DefaultSubmissionService) GetSubmissions(problemId string, page, pageSize int) (res []*models.Submission, err error) {
	offset := (page - 1) * pageSize
	res, err = d.SubmissionRepository.GetSubmissions(offset, pageSize, problemId)
	if err != nil {
		d.logger.Error("get submissions", zap.String("problem id", problemId), zap.Error(err))
		return nil, err
	}
	return
}

func (d DefaultSubmissionService) GetSubmission(submissionId string) (s *models.Submission, err error) {
	s, err = d.SubmissionRepository.GetSubmission(submissionId)
	if err != nil {
		d.logger.Error("get submission", zap.String("submission id", submissionId), zap.Error(err))
		return nil, err
	}
	return
}

func (d DefaultSubmissionService) Create(submitterID uint64, problemName, userSpace string) (s *models.Submission, err error) {
	d.logger.Debug("create submission",
		zap.Uint64("submitter Id", submitterID),
		zap.String("problem name", problemName),
		zap.String("user space", userSpace),
	)
	problem, err := d.ProblemRepository.GetProblemByName(problemName)
	if err != nil {
		d.logger.Error("create submission", zap.Error(err))
		return nil, err
	}
	if problem == nil {
		d.logger.Error("create submission: unknown problem")
		return nil, errors.New("unknown problem")
	}
	d.logger.Debug("create submission",
		zap.Uint64("submitter Id", submitterID),
		zap.Uint64("problem id", problem.ID),
		zap.String("user space", userSpace),
	)
	s, err = d.SubmissionRepository.Create(submitterID, problem.ID, userSpace)
	if err != nil {
		d.logger.Error("create submission", zap.Error(err))
		return nil, err
	}
	return
}

func (d DefaultSubmissionService) DispatchJudgement(submissionId string) error {
	submission, err := d.SubmissionRepository.GetSubmission(submissionId)
	d.logger.Info("dispatch judgement")
	if err != nil {
		d.logger.Error("error:", zap.Error(err), zap.String("submissionId", submissionId))
		return err
	}

	if submission == nil {
		d.logger.Error("unknown submission")
		return errors.New("unknown submission")
	}

	submissionElement := d.SubmissionRepository.CreateProcess(submission)

	err = d.DeliverJudgement(submissionElement)
	if err != nil {
		d.logger.Error("dispatch judgement error", zap.Error(err))
		return err
	}
	return nil
}

func (d DefaultSubmissionService) DeliverJudgement(element *repositories.Process) error {

	upstreams := element.FindUpstreams()

	for _, upstream := range upstreams {
		upstreamType := upstream.Type

		fmt.Println(upstream.Properties)

		d.logger.Info("create judgement",
			zap.String("process Id", element.ProcessId),
			zap.Int("block Id", upstream.Id),
			zap.String("judgement type", upstreamType),
		)

		//judgementId, err := d.JudgementService.Create(
		//	context.TODO(),
		//	upstreamType,
		//	upstream.Properties,
		//	upstream.Inputs,
		//)
		//
		//if err != nil {
		//	d.logger.Error("create judgement error", zap.Error(err))
		//	return err
		//}
		//d.logger.Info("create judgement success", zap.String("judgement id", judgementId))
		//
		//d.idMap[judgementId] = upstream.Id
		//d.processMap[judgementId] = element.ProcessId
	}

	return nil

}

func (d DefaultSubmissionService) ReturnJudgement(judgementId string, outputs [][]byte) error {
	d.logger.Info("return judgement",
		zap.String("judgement id", judgementId),
	)

	blockId, ok := d.idMap[judgementId]
	if !ok {
		err := errors.New("unknown judgement id")
		d.logger.Error("unknown judgement id", zap.String("judgement id", judgementId))
		return err
	}
	processId, ok := d.processMap[judgementId]
	if !ok {
		err := errors.New("unknown judgement id")
		d.logger.Error("unknown judgement id", zap.String("judgement id", judgementId))
		return err
	}

	submissionElement := d.SubmissionRepository.FetchProcess(processId)
	if submissionElement == nil {
		err := errors.New("internal error: unknown submission")
		d.logger.Error("return judgement failed", zap.Error(err))
		return err
	}
	err := submissionElement.SetOutputs(blockId, outputs)
	if err != nil {
		d.logger.Error("return judgement failed", zap.Error(err))
	}

	err = d.DeliverJudgement(submissionElement)
	return err
}

func NewSubmissionService(
	logger *zap.Logger,
	Repository repositories.Repository,
	pRepository problemRepository.Repository,
) SubmissionsService {
	return &DefaultSubmissionService{
		logger:               logger.With(zap.String("type", "DefaultSubmissionService")),
		SubmissionRepository: Repository,
		ProblemRepository:    pRepository,
		processMap:           make(map[string]string),
		idMap:                make(map[string]int),
	}
}
