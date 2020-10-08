package controllers

import (
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/infinity-oj/server-v2/internal/pkg/sessions"

	"github.com/gin-gonic/gin"
	"github.com/infinity-oj/server-v2/internal/app/volumes/services"
	"go.uber.org/zap"
)

type Controller interface {
	CreateVolume(c *gin.Context)

	CreateFile(c *gin.Context)
	CreateDirectory(c *gin.Context)

	GetFile(c *gin.Context)
	GetDirectory(c *gin.Context)
}

type DefaultController struct {
	logger  *zap.Logger
	service services.Service
}

func (d DefaultController) CreateFile(c *gin.Context) {
	session := sessions.GetSession(c)
	if session == nil {
		d.logger.Debug("get principal failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// single formFile
	formFile, err := c.FormFile("file")
	if err != nil {
		d.logger.Error("create file failed")
		return
	}
	d.logger.Debug("upload file",
		zap.String("filename", formFile.Filename),
	)

	volume := c.Param("name")

	file, err := formFile.Open()
	if err != nil {
		d.logger.Error("create file failed")
		return
	}
	fileData, _ := ioutil.ReadAll(file)
	err = d.service.CreateFile(volume, formFile.Filename, fileData)
	if err != nil {
		d.logger.Error("create file failed")
		return
	}

	c.Status(http.StatusNoContent)
}

func (d DefaultController) CreateDirectory(c *gin.Context) {
	session := sessions.GetSession(c)
	if session == nil {
		d.logger.Debug("get principal failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	request := struct {
		Dirname string `json:"dirname" binding:"required"`
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

	volume := c.Param("name")

	err := d.service.CreateDirectory(volume, request.Dirname)
	if err != nil {
		d.logger.Error("create volume failed")
		return
	}

	c.Status(http.StatusNoContent)
}

func (d DefaultController) CreateVolume(c *gin.Context) {
	session := sessions.GetSession(c)
	if session == nil {
		d.logger.Debug("get principal failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	volume, err := d.service.CreateVolume()
	if err != nil {
		d.logger.Error("create volume failed")
		return
	}

	c.JSON(http.StatusOK, volume)
}

func (d DefaultController) GetFile(c *gin.Context) {
	panic("implement me")
}

func (d DefaultController) GetDirectory(c *gin.Context) {
	panic("implement me")
}

func New(logger *zap.Logger, s services.Service) Controller {
	return &DefaultController{
		logger:  logger,
		service: s,
	}
}
