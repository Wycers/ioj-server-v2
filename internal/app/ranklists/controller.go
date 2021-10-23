package ranklists

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller interface {
	GetRankList(c *gin.Context)
}

type controller struct {
	logger  *zap.Logger
	service Service
}

func (ctl *controller) GetRankList(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	ctl.logger.Debug("get ranklist", zap.Uint64("ranklist id", id))

	rankList, err := ctl.service.GetRankList(id)
	if err != nil {
		ctl.logger.Error("get account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	if rankList == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, rankList)
}

//
//func (pc *DefaultController) CreateProgram(c *gin.Context) {
//	session := sessions.GetSession(c)
//	if session == nil {
//		pc.logger.Debug("get principal failed")
//		c.AbortWithStatus(http.StatusUnauthorized)
//		return
//	}
//
//	pc.logger.Debug("create program", zap.Uint64("account id", session.AccountId))
//	request := struct {
//		Definition string `json:"definition" binding:"required,gt=0"`
//	}{}
//
//	if err := c.ShouldBind(&request); err != nil {
//		errs, ok := err.(validator.ValidationErrors)
//		if !ok {
//			c.JSON(http.StatusOK, gin.H{
//				"msg": err.Error(),
//			})
//			return
//		}
//
//		c.JSON(http.StatusBadRequest, gin.H{
//			"msg": errs.Error(),
//		})
//		return
//	}
//	problem, err := pc.service.CreateProgram(request.Definition)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusInternalServerError, &gin.H{
//			"message": err.Error(),
//		})
//		return
//	}
//
//	c.JSON(http.StatusOK, problem)
//}
//
//func (pc *DefaultController) GetProgram(c *gin.Context) {
//	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
//
//	pc.logger.Debug("get program", zap.Uint64("program name", id))
//
//	problem, err := pc.service.GetProgram(id)
//	if err != nil {
//		pc.logger.Error("get account", zap.Error(err))
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"message": err.Error(),
//		})
//		return
//	}
//	if problem == nil {
//		c.AbortWithStatus(http.StatusNotFound)
//		return
//	}
//	c.JSON(http.StatusOK, problem)
//}
//
//func (pc *DefaultController) GetPrograms(c *gin.Context) {
//	pc.logger.Debug("get programs")
//
//	programs, err := pc.service.GetPrograms()
//	if err != nil {
//		pc.logger.Error("get account", zap.Error(err))
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"message": err.Error(),
//		})
//		return
//	}
//	if programs == nil {
//		c.AbortWithStatus(http.StatusNotFound)
//		return
//	}
//	c.JSON(http.StatusOK, programs)
//}

func NewController(logger *zap.Logger, s Service) Controller {
	return &controller{
		logger:  logger,
		service: s,
	}
}
