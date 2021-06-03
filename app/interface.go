package app

import (
	"github.com/gin-gonic/gin"
	"github.com/v-mars/frame/db"
)

type CommonInterfaces interface {
	Get(c *gin.Context)
	Query(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Option() db.Option
}