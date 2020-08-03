package repositories

import (
	"container/list"
	"sync"

	"github.com/infinity-oj/server-v2/internal/pkg/nodeEngine"

	"github.com/google/uuid"

	"github.com/infinity-oj/server-v2/internal/pkg/models"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Repository interface {
	GetSubmission(submissionId string) (*models.Submission, error)
	GetSubmissions(offset, limit int, problemId string) ([]*models.Submission, error)
	Create(submitterID, problemId uint64, userSpace string) (s *models.Submission, err error)
	Update(s *models.Submission) error

	CreateProcess(s *models.Submission) *Process
	FetchProcess(processId string) *Process
}

type DefaultRepository struct {
	logger *zap.Logger
	db     *gorm.DB
	queue  *list.List
	mutex  *sync.Mutex
}

func (m DefaultRepository) GetSubmissions(offset, limit int, problemId string) (res []*models.Submission, err error) {
	if err = m.db.Table("submissions").Where("problem_id = ?", problemId).
		Offset(offset).Limit(limit).
		Find(&res).Error; err != nil {
		return nil, err
	}
	return
}

type Process struct {
	ProcessId string

	obj    *models.Submission
	graph  *nodeEngine.Graph
	result map[int][]byte
}

type JudgementElement struct {
	Id         int
	Type       string
	Properties map[string]string
	Inputs     [][]byte
}

func (se Process) FindUpstreams() (res []*JudgementElement) {

	ids := se.graph.Run()

	for _, block := range ids {
		var inputs [][]byte
		for _, linkId := range block.Inputs {
			if data, ok := se.result[linkId]; ok {
				inputs = append(inputs, data)
			} else {
				panic("something wrong")
			}
		}

		res = append(res, &JudgementElement{
			Id:         block.Id,
			Type:       block.Type,
			Properties: block.Properties,
			Inputs:     inputs,
		})
	}

	return
}

func (se Process) SetOutputs(blockId int, outputs [][]byte) error {

	block := se.graph.FindBlockById(blockId)

	if len(block.Output) != len(outputs) {
		return errors.New("output slots mismatch")
	}

	for index, output := range outputs {
		links := se.graph.FindLinkBySourcePort(blockId, index)
		for _, link := range links {
			se.result[link.Id] = output
		}
	}

	block.Done()
	return nil
}

func (m DefaultRepository) CreateProcess(s *models.Submission) *Process {

	graph := nodeEngine.NewGraphByFile("easyGraph.json")

	result := make(map[int][]byte)

	element := &Process{
		ProcessId: uuid.New().String(),
		obj:       s,
		graph:     graph,
		result:    result,
	}

	m.queue.PushBack(element)

	return element
}

func (m DefaultRepository) FetchProcess(processId string) *Process {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for se := m.queue.Front(); se != nil; se = se.Next() {
		element, ok := se.Value.(*Process)
		if !ok {
			continue
		}

		if element.ProcessId == processId {
			return element
		}
	}

	return nil
}

func (m DefaultRepository) GetSubmission(submissionId string) (*models.Submission, error) {
	submission := &models.Submission{}
	err := m.db.Where("submission_id = ?", submissionId).First(&submission).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return submission, nil
}

func (m DefaultRepository) Create(submitterId, problemId uint64, userSpace string) (s *models.Submission, err error) {
	s = &models.Submission{
		SubmissionId: uuid.New().String(),
		SubmitterId:  submitterId,
		ProblemId:    problemId,
		UserSpace:    userSpace,
	}

	if err = m.db.Create(s).Error; err != nil {
		return nil, errors.Wrapf(err,
			" create submission with username: %d, problemID: %s, userSpace: %s",
			submitterId, problemId, userSpace,
		)
	}

	return
}

func (m DefaultRepository) GetUpstreams() (res []*JudgementElement) {

	for se := m.queue.Front(); se != nil; se = se.Next() {
		judgement, ok := se.Value.(*Process)

		if !ok {
			continue
		}

		res = append(res, judgement.FindUpstreams()...)
	}

	return
}

func (m DefaultRepository) Update(s *models.Submission) (err error) {
	err = m.db.Save(s).Error
	return
}

func NewMysqlSubmissionsRepository(logger *zap.Logger, db *gorm.DB) Repository {
	return &DefaultRepository{
		logger: logger.With(zap.String("type", "Repository")),
		db:     db,
		queue:  list.New(),
		mutex:  &sync.Mutex{},
	}
}
