package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/infinity-oj/server-v2/internal/pkg/sessions"

	"github.com/gin-gonic/gin"
	"github.com/infinity-oj/server-v2/internal/app/processes/services"
	"go.uber.org/zap"
)

type Controller interface {
	CreateProcess(c *gin.Context)
	GetProcess(c *gin.Context)
}

type DefaultController struct {
	logger  *zap.Logger
	service services.ProcessesService
}

func (pc *DefaultController) CreateProcess(c *gin.Context) {
	session := sessions.GetSession(c)
	if session == nil {
		pc.logger.Debug("get principal failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	pc.logger.Debug("create process", zap.Uint64("account id", session.AccountId))
	request := struct {
		Definition string `json:"definition" binding:"required,gt=0"`
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
	problem, err := pc.service.CreateProcess(request.Definition)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, problem)
}

func (pc *DefaultController) GetProcess(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	pc.logger.Debug("get process", zap.Uint64("problem name", id))

	problem, err := pc.service.GetProcess(id)
	if err != nil {
		pc.logger.Error("get account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	if problem == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, problem)
}

func New(logger *zap.Logger, s services.ProcessesService) Controller {
	return &DefaultController{
		logger:  logger,
		service: s,
	}
}
