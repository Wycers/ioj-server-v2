package controllers

import (
	"io/ioutil"
	"net/http"
	"strings"

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

	DeleteFile(c *gin.Context)

	GetVolume(c *gin.Context)
	DownloadDirectory(c *gin.Context)
	DownloadFile(c *gin.Context)
}
type DefaultController struct {
	logger  *zap.Logger
	service services.Service
}

func (d DefaultController) DownloadDirectory(c *gin.Context) {
	request := struct {
		Dirname string `form:"dirname"`
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

	if request.Dirname == "" {
		request.Dirname = "/"
	}

	volume := c.Param("name")

	file, err := d.service.GetDirectory(volume, "/")
	if err != nil {
		d.logger.Error("Download directory", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.File(file.Name())
}

func (d DefaultController) CreateFile(c *gin.Context) {
	session := sessions.GetSession(c)
	if session == nil {
		d.logger.Debug("get principal failed")
		//c.AbortWithStatus(http.StatusUnauthorized)
		//return
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
	volume, err := d.service.CreateFile(volumeName, "/", formFile.Filename, fileData)
	if err != nil {
		d.logger.Error("create file failed", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, volume)
}

func (d DefaultController) DeleteFile(c *gin.Context) {
	session := sessions.GetSession(c)
	if session == nil {
		d.logger.Debug("get principal failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	request := struct {
		Filename string `form:"filename" binding:"required,gt=0"`
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

	d.logger.Debug("upload file",
		zap.String("filename", request.Filename),
	)

	volumeName := c.Param("name")

	volume, err := d.service.RemoveFile(volumeName, "/", request.Filename)
	if err != nil {
		d.logger.Error("remove file failed", zap.Error(err))
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
	createBy := uint64(0)
	if session != nil {
		//d.logger.Debug("get principal failed")
		//c.AbortWithStatus(http.StatusUnauthorized)
		//return
		createBy = session.AccountId
	}

	volume, err := d.service.CreateVolume(createBy)
	if err != nil {
		d.logger.Error("create volume failed")
		return
	}

	c.JSON(http.StatusOK, volume)
}

func (d DefaultController) DownloadFile(c *gin.Context) {
	request := struct {
		Filename string `form:"filename" binding:"required,gt=0"`
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
	filename := strings.ReplaceAll(request.Filename, "%2f", "/")
	filename = strings.ReplaceAll(request.Filename, "%2F", "/")

	file, err := d.service.GetFile(volume, filename)
	if err != nil {
		if err.Error() == "not found" {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			d.logger.Error("Download File", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	c.File(file.Name())
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

func New(logger *zap.Logger, s services.Service) Controller {
	return &DefaultController{
		logger:  logger,
		service: s,
	}
}
