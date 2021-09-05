package judgements

import (
	"net/http"

	"github.com/infinity-oj/server-v2/pkg/models"

	"github.com/go-playground/validator/v10"
	"github.com/infinity-oj/server-v2/internal/pkg/sessions"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller interface {
	CreateJudgement(c *gin.Context)
	GetJudgements(c *gin.Context)
	GetJudgementPrerequisites(c *gin.Context)
	GetJudgement(c *gin.Context)
	CancelJudgement(c *gin.Context)
}

type DefaultController struct {
	logger  *zap.Logger
	service Service
}

func (d *DefaultController) GetJudgementPrerequisites(c *gin.Context) {
	c.JSON(200, &gin.H{
		"upload": "*.cpp,*.c,*.py,*.zip",
	})
}

func (d *DefaultController) CancelJudgement(c *gin.Context) {
	d.logger.Debug("cancel judgement")
	session := sessions.GetSession(c)
	if session == nil {
		d.logger.Debug("get principal failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	judgementId := c.Param("judgementId")

	d.logger.Debug("cancel judgement",
		zap.Uint64("account id", session.AccountId),
		zap.String("judgement id", judgementId),
	)

	judgement, err := d.service.GetJudgement(judgementId)
	if err != nil {
		d.logger.Error("cancel judgement", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	if judgement == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if judgement.Status != models.Pending {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	judgement, err = d.service.UpdateJudgement(judgementId, models.Canceled, -1, "User cancel")
	if err != nil {
		d.logger.Error("cancel judgement", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(200, judgement)
}

func (d *DefaultController) CreateJudgement(c *gin.Context) {
	d.logger.Debug("create judgement by process")
	session := sessions.GetSession(c)
	if session == nil {
		d.logger.Debug("get principal failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	d.logger.Debug("create judgement", zap.Uint64("account id", session.AccountId))

	request := struct {
		ProcessId    uint64 `json:"processId" binding:"required"`
		SubmissionId uint64 `json:"submissionId" binding:"required"`
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

	d.logger.Debug("create judgement",
		zap.Uint64("process id", request.ProcessId),
		zap.Uint64("submission id", request.SubmissionId),
	)

	code, judgement, err := d.service.CreateJudgement(session.AccountId, request.ProcessId, request.SubmissionId)
	if err != nil {
		d.logger.Error("create judgement", zap.Error(err))
		c.JSON(code, gin.H{
			"msg": err.Error(),
		})
		return
	}
	d.logger.Debug("create judgement",
		zap.String("new judgement id", judgement.JudgementId),
	)
	c.JSON(http.StatusOK, judgement)
}

func (d *DefaultController) GetJudgements(c *gin.Context) {
	d.logger.Debug("get judgements")
	session := sessions.GetSession(c)
	if session == nil {
		d.logger.Debug("get principal failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	d.logger.Debug("create judgement", zap.Uint64("account id", session.AccountId))

	judgements, err := d.service.GetJudgements(session.AccountId)
	if err != nil {
		d.logger.Error("create judgement", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(200, judgements)
}

func (d *DefaultController) GetJudgement(c *gin.Context) {
	d.logger.Debug("get judgement")
	session := sessions.GetSession(c)
	if session == nil {
		d.logger.Debug("get principal failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	judgementId := c.Param("judgementId")
	d.logger.Debug("get judgement", zap.String("judgement id", judgementId))

	judgement, err := d.service.GetJudgement(judgementId)
	if err != nil {
		d.logger.Error("create judgement", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(200, judgement)
}

func NewController(logger *zap.Logger, s Service) Controller {
	return &DefaultController{
		logger:  logger,
		service: s,
	}
}
