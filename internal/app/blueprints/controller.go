package blueprints

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/infinity-oj/server-v2/internal/pkg/sessions"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller interface {
	CreateBlueprint(c *gin.Context)
	GetJudgementPrerequisites(c *gin.Context)
	GetBlueprint(c *gin.Context)
	GetBlueprints(c *gin.Context)
}

type DefaultController struct {
	logger  *zap.Logger
	service Service
}

func (pc *DefaultController) GetBlueprints(c *gin.Context) {
	pc.logger.Debug("get blueprints")

	programs, err := pc.service.GetBlueprints()
	if err != nil {
		pc.logger.Error("get account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	if programs == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, programs)
}

func (pc *DefaultController) GetJudgementPrerequisites(c *gin.Context) {
	c.JSON(200, &gin.H{
		"upload": "*.cpp,*.c,*.py,*.zip",
	})
}

func (pc *DefaultController) CreateBlueprint(c *gin.Context) {
	session := sessions.GetSession(c)
	if session == nil {
		pc.logger.Debug("get principal failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	pc.logger.Debug("create blueprint", zap.Uint64("account id", session.AccountId))
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

	problem, err := pc.service.CreateBlueprint(request.Definition)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, problem)
}

func (pc *DefaultController) GetBlueprint(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	pc.logger.Debug("get blueprint", zap.Uint64("problem name", id))

	problem, err := pc.service.GetBlueprint(id)
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

func NewController(logger *zap.Logger, s Service) Controller {
	return &DefaultController{
		logger:  logger,
		service: s,
	}
}
