package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/infinity-oj/server-v2/internal/app/accounts/services"
	"github.com/infinity-oj/server-v2/internal/pkg/sessions"
	"go.uber.org/zap"
)

type Controller interface {
	CreateAccount(c *gin.Context)
	GetAccount(c *gin.Context)
	UpdateAccount(c *gin.Context)
	UpdateAccountCredential(c *gin.Context)

	CreatePrincipal(c *gin.Context)
	GetPrincipal(c *gin.Context)
	DeletePrincipal(c *gin.Context)

	GetRole(c *gin.Context)
}

type DefaultController struct {
	logger  *zap.Logger
	service services.Service
}

func New(logger *zap.Logger, s services.Service) Controller {
	return &DefaultController{
		logger:  logger,
		service: s,
	}
}

func (d DefaultController) CreateAccount(c *gin.Context) {
	c.AbortWithStatus(http.StatusForbidden)
	return

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

func (d DefaultController) UpdateAccountCredential(c *gin.Context) {
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
		OldPassword string `json:"oldPassword"  binding:"required,gte=6"`
		NewPassword string `json:"newPassword"  binding:"required,gte=6"`
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
	res, err := d.service.UpdateCredential(account.Name, request.OldPassword, request.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	if res  {
		c.Status(http.StatusNoContent)
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
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

func (d DefaultController) CreatePrincipal(c *gin.Context) {
	session := sessions.GetSession(c)
	if session != nil {
		session.Clear(c)
	}
	request := struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
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

	isValid, err := d.service.VerifyCredential(request.Username, request.Password)
	if err != nil {
		d.logger.Error("verify credential", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if !isValid {
		d.logger.Debug("verify credential", zap.String("username", request.Username))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	account, err := d.service.GetAccount(request.Username)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	session = sessions.New()
	session.AccountId = account.ID

	roles, err := d.service.GetRoleById(account.ID)
	if err != nil {
		d.logger.Error("ger roles failed", zap.Uint64("accountId", account.ID))
	}
	if roles != nil {
		session.Roles = []string{}
		for _, v := range roles {
			session.Roles = append(session.Roles, v.Name)
		}
	}

	err = session.Save(c)
	if err != nil {
		d.logger.Error("verify credential, save session", zap.String("username", request.Username))
	}

	c.Status(http.StatusNoContent)
}

func (d DefaultController) GetPrincipal(c *gin.Context) {
	session := sessions.GetSession(c)
	if session == nil {
		d.logger.Debug("get principal failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	account, err := d.service.GetAccountById(session.AccountId)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if account == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, account)
}

func (d DefaultController) DeletePrincipal(c *gin.Context) {
	session := sessions.GetSession(c)
	if session != nil {
		session.Clear(c)
	}
	c.Status(http.StatusNoContent)
}

func (d DefaultController) DeleteAccount(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (d DefaultController) GetRole(c *gin.Context) {
	session := sessions.GetSession(c)
	if session == nil {
		d.logger.Debug("get principal failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"roles": session.Roles,
	})

}
