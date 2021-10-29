package processes

import (
	"net/http"

	"github.com/infinity-oj/server-v2/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Controller interface {
	GetProcesses(c *gin.Context)
	GetProcess(c *gin.Context)
	UpdateProcess(c *gin.Context)
	ReserveProcess(c *gin.Context)
}

type DefaultController struct {
	logger  *zap.Logger
	service Service
}

func (d *DefaultController) GetProcesses(c *gin.Context) {
	d.logger.Debug("get processes")
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

	d.logger.Debug("get processes", zap.String("page", request.Type))

	processes, err := d.service.GetProcesses(request.Type)
	if err != nil {
		d.logger.Error("get processes", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, processes)
}

func (d *DefaultController) GetProcess(c *gin.Context) {
	processId := c.Param("processId")

	d.logger.Debug("get process", zap.String("process processId", processId))

	process, err := d.service.GetProcess(processId)
	if err != nil {
		d.logger.Error("get process", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	if process == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, process)
}

func (d *DefaultController) UpdateProcess(c *gin.Context) {
	//session := sessions.GetSession(c)
	//if session == nil {
	//	d.logger.Debug("get principal failed")
	//	//c.AbortWithStatus(http.StatusUnauthorized)
	//	//return
	//}

	processId := c.Param("processId")

	d.logger.Debug("update process",
		//zap.Uint64("account id", session.AccountId),
		zap.String("process id", processId),
	)

	request := struct {
		Token   string       `json:"token" binding:"required"`
		Warning string       `json:"warning" binding:""`
		Error   string       `json:"error" binding:""`
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
	d.logger.Debug("update process",
		zap.String("token", request.Token),
		zap.String("warning", request.Warning),
		zap.String("error", request.Error),
	)

	process, err := d.service.UpdateProcess(processId, request.Warning, request.Error, &request.Outputs)
	if err != nil {
		d.logger.Error("update process", zap.Error(err))

		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(200, process)
}

func (d *DefaultController) ReserveProcess(c *gin.Context) {
	processId := c.Param("processId")

	d.logger.Debug("reserve process", zap.String("process processId", processId))

	token, locked, err := d.service.ReserveProcess(processId)
	if !locked {
		if err != nil {
			d.logger.Error("reserve process failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		d.logger.Error("reserve process failed: locked before")
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
