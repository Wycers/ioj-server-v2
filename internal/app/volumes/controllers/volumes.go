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

	GetVolume(c *gin.Context)

	DownloadDirectory(c *gin.Context)
	GetFile(c *gin.Context)
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
		d.logger.Error("create file failed", zap.Error(err))
		return
	}
	d.logger.Debug("upload file",
		zap.String("filename", formFile.Filename),
	)

	volumeName := c.Param("name")

	file, err := formFile.Open()
	if err != nil {
		d.logger.Error("create file failed", zap.Error(err))
		return
	}
	fileData, _ := ioutil.ReadAll(file)
	volume, err := d.service.CreateFile(volumeName, formFile.Filename, fileData)
	if err != nil {
		d.logger.Error("create file failed", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, volume)
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

	volumeName := c.Param("name")

	volume, err := d.service.CreateDirectory(volumeName, request.Dirname)
	if err != nil {
		d.logger.Error("create volume failed")
		return
	}

	c.JSON(http.StatusOK, volume)
}

func (d DefaultController) CreateVolume(c *gin.Context) {
	session := sessions.GetSession(c)
	if session == nil {
		d.logger.Debug("get principal failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	volume, err := d.service.CreateVolume(session.AccountId)
	if err != nil {
		d.logger.Error("create volume failed")
		return
	}

	c.JSON(http.StatusOK, volume)
}

func (d DefaultController) GetFile(c *gin.Context) {
	panic("implement me")
}

func (d DefaultController) GetVolume(c *gin.Context) {
	session := sessions.GetSession(c)
	if session == nil {
		d.logger.Debug("get principal failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	volumeName := c.Param("name")

	volume, err := d.service.GetVolume(volumeName)
	if err != nil {
		d.logger.Error("create volume failed")
		return
	}

	c.JSON(http.StatusOK, volume)
}

func (d DefaultController) DownloadDirectory(c *gin.Context) {

	request := struct {
		Dirname string `form:"dirname" binding:"required"`
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

	volume := c.Param("name")

	file, err := d.service.DownloadDirectory(volume, request.Dirname)
	if err != nil {
		d.logger.Error("Download directory", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.File(file.Name())
}

func New(logger *zap.Logger, s services.Service) Controller {
	return &DefaultController{
		logger:  logger,
		service: s,
	}
}
