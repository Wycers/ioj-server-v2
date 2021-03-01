package submissions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/infinity-oj/server-v2/internal/pkg/sessions"
	"go.uber.org/zap"
)

type Controller interface {
	CreateSubmission(c *gin.Context)
	GetSubmissions(c *gin.Context)
	GetSubmission(c *gin.Context)
}

type controller struct {
	logger  *zap.Logger
	service Service
}

func (d *controller) CreateSubmission(c *gin.Context) {
	d.logger.Debug("create submission")
	session := sessions.GetSession(c)
	if session == nil {
		d.logger.Debug("get principal failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	d.logger.Debug("create submission", zap.Uint64("account id", session.AccountId))

	request := struct {
		ProblemId string `json:"problemId" binding:"required"`
		UserSpace string `json:"volume" binding:"required"`
	}{}

	if err := c.ShouldBind(&request); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"msg": errs.Error(),
		})
		return
	}
	d.logger.Debug("create submission",
		zap.String("problem id", request.ProblemId),
		zap.String("user space", request.UserSpace),
	)

	code, submission, judgement, err := d.service.Create(session.AccountId, request.ProblemId, request.UserSpace)
	if err != nil {
		d.logger.Error("create submission", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	if code != http.StatusOK {
		c.AbortWithStatus(code)
		return
	}
	c.JSON(200, &gin.H{
		"submission": submission,
		"judgement":  judgement,
	})
}

func (d *controller) GetSubmissions(c *gin.Context) {
	request := struct {
		ProblemId int `form:"problemId"`

		Page     int `form:"page" binding:"required,gt=0"`
		PageSize int `form:"pageSize" binding:"required,gt=0,lte=15"`

		Sort string `json:"sort"`
	}{}

	if err := c.ShouldBindQuery(&request); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"msg": errs.Error(),
		})
		return
	}

}

func (d *controller) GetSubmission(c *gin.Context) {
	submissionId := c.Param("submissionId")
	submission, err := d.service.GetSubmission(submissionId)
	if err != nil {
		d.logger.Error("get submission", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}
	if submission == nil {
		d.logger.Debug("get submission missing")
		c.Status(http.StatusNotFound)
		return
	}

	session := sessions.GetSession(c)
	if session == nil {
		d.logger.Debug("get principal failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if submission.SubmitterId != session.AccountId {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, submission)
}

func NewController(logger *zap.Logger, s Service) Controller {
	return &controller{
		logger:  logger,
		service: s,
	}
}
