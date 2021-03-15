package problems

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/infinity-oj/server-v2/internal/pkg/sessions"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller interface {
	CreateProblem(c *gin.Context)
	GetProblems(c *gin.Context)
	GetProblem(c *gin.Context)
	UpdateProblem(c *gin.Context)
}

type DefaultController struct {
	logger  *zap.Logger
	service Service
}

func (pc *DefaultController) CreateProblem(c *gin.Context) {
	session := sessions.GetSession(c)
	if session == nil {
		pc.logger.Debug("get principal failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	pc.logger.Debug("create problem", zap.Uint64("account id", session.AccountId))
	request := struct {
		Name  string `json:"name" binding:"required,gt=0"`
		Title string `json:"title" binding:"required,gt=0"`
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
	problem, err := pc.service.CreateProblem(request.Name, request.Title)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, problem)
}

func (pc *DefaultController) GetProblems(c *gin.Context) {
	request := struct {
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

	pc.logger.Debug("get problems",
		zap.Int("page", request.Page),
		zap.Int("pageSize", request.PageSize),
	)

	problem, err := pc.service.GetProblems(request.Page, request.PageSize)
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

func (pc *DefaultController) GetProblem(c *gin.Context) {
	name := c.Param("name")

	pc.logger.Debug("get problem", zap.String("problem name", name))

	problem, err := pc.service.GetProblemByName(name)
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

func (pc *DefaultController) UpdateProblem(c *gin.Context) {
	name := c.Param("name")
	problem, err := pc.service.GetProblemByName(name)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if problem == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	session := sessions.GetSession(c)
	if session == nil {
		pc.logger.Debug("get principal failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	pc.logger.Debug("update problem", zap.Uint64("account id", session.AccountId))

	request := struct {
		Title string `json:"title" binding:"required,gt=0"`

		PublicVolume  string `json:"publicVolume" binding:"required,gt=0"`
		PrivateVolume string `json:"privateVolume" binding:"required,gt=0"`
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

	problem, err = pc.service.UpdateProblem(problem, name, request.Title, request.PublicVolume, request.PrivateVolume)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, problem)
}

func NewController(logger *zap.Logger, s Service) Controller {
	return &DefaultController{
		logger:  logger,
		service: s,
	}
}
