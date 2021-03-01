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
	GetJudgement(c *gin.Context)
	CancelJudgement(c *gin.Context)

	GetTasks(c *gin.Context)
	GetTask(c *gin.Context)
	UpdateTask(c *gin.Context)
	ReserveTask(c *gin.Context)
}

type DefaultController struct {
	logger  *zap.Logger
	service Service
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
		d.logger.Error("create judgement", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(200, judgement)
}

func (d *DefaultController) CreateJudgement(c *gin.Context) {
	d.logger.Debug("create judgement")
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
		zap.Uint64("processes id", request.ProcessId),
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

func (d *DefaultController) GetTasks(c *gin.Context) {
	request := struct {
		Type string `form:"type" binding:"required,gt=0"`
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

	d.logger.Debug("get tasks",
		zap.String("page", request.Type),
	)

	tasks, err := d.service.GetTasks(request.Type)
	if err != nil {
		d.logger.Error("get tasks", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (d *DefaultController) GetTask(c *gin.Context) {
	taskId := c.Param("taskId")

	d.logger.Debug("get task", zap.String("task taskId", taskId))

	task, err := d.service.GetTask(taskId)
	if err != nil {
		d.logger.Error("get task", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	if task == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, task)
}

func (d *DefaultController) UpdateTask(c *gin.Context) {
	session := sessions.GetSession(c)
	if session == nil {
		d.logger.Debug("get principal failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	taskId := c.Param("taskId")

	d.logger.Debug("update task",
		zap.Uint64("account id", session.AccountId),
		zap.String("task id", taskId),
	)

	request := struct {
		Token   string `json:"token" binding:"required"`
		Warning string `json:"warning" binding:""`
		Error   string `json:"error" binding:""`

		Outputs models.Slots `json:"outputs" binding:"required"`
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
	d.logger.Debug("update task",
		zap.String("token", request.Token),
		zap.String("warning", request.Warning),
		zap.String("error", request.Error),
	)

	task, err := d.service.UpdateTask(taskId, request.Warning, request.Error, &request.Outputs)
	if err != nil {
		d.logger.Error("update task", zap.Error(err))

		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(200, task)
}

func (d *DefaultController) ReserveTask(c *gin.Context) {
	taskId := c.Param("taskId")

	d.logger.Debug("reserve task", zap.String("task taskId", taskId))

	token, locked, err := d.service.ReserveTask(taskId)
	if !locked {
		if err != nil {
			d.logger.Error("reserve task failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		d.logger.Error("reserve task failed: locked before")
		c.Status(http.StatusPreconditionFailed)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func NewController(logger *zap.Logger, s Service) Controller {
	return &DefaultController{
		logger:  logger,
		service: s,
	}
}
