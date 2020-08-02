package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/infinity-oj/server-v2/internal/app/accounts/services"
	"go.uber.org/zap"
	"net/http"
)

type Controller interface {
	CreateAccount(c *gin.Context)
	GetAccount(c *gin.Context)
	UpdateAccount(c *gin.Context)
}

type DefaultController struct {
	logger  *zap.Logger
	service services.Service
}

func (d DefaultController) CreateAccount(c *gin.Context) {
	request := struct {
		Username string `json:"username" binding:"required,gte=6"`
		Password string `json:"password" binding:"required,gte=6"`
		Email    string `json:"email" binding:"required,email"`
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

	account, err := d.service.CreateAccount(request.Username, request.Password, request.Email)
	if err != nil {
		d.logger.Error("create account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, account)
}

func (d DefaultController) GetAccount(c *gin.Context) {
	name := c.Param("name")

	d.logger.Debug("get account", zap.String("account name", name))

	account, err := d.service.GetAccount(name)
	if err != nil {
		d.logger.Error("get account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	if account == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, account)
}

func (d DefaultController) UpdateAccount(c *gin.Context) {
	name := c.Param("name")

	account, err := d.service.GetAccount(name)
	if err != nil {
		d.logger.Error("get account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	if account == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	request := struct {
		Nickname string `json:"nickname" binding:"required"`
		Locale   string `json:"locale" binding:"required,gt=0"`
		Email    string `json:"email" binding:"required,email"`
		Gender   string `json:"gender" binding:"required"`
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

	d.logger.Debug("update account", zap.String("account name", name))
	account, err = d.service.UpdateAccount(
		account,
		request.Nickname,
		request.Email,
		request.Gender,
		request.Locale,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, account)
}

func New(logger *zap.Logger, s services.Service) Controller {
	return &DefaultController{
		logger:  logger,
		service: s,
	}
}
