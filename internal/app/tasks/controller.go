package tasks

import (
	"net/http"

	"github.com/infinity-oj/server-v2/pkg/models"

	"github.com/go-playground/validator/v10"
	"github.com/infinity-oj/server-v2/internal/pkg/sessions"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller interface {
	GetTasks(c *gin.Context)
	GetTask(c *gin.Context)
	UpdateTask(c *gin.Context)
	ReserveTask(c *gin.Context)
}

type DefaultController struct {
	logger  *zap.Logger
	service Service
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
