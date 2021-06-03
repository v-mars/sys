package server

import (
	"github.com/gin-gonic/gin"
	"github.com/v-mars/frame/response"
	"gorm.io/gorm"
)

type IState interface {
	Get(c *gin.Context)
}
type SState struct {
	DB *gorm.DB
}

func NewService(DB *gorm.DB) IState {
	return SState{DB: DB}
}

func (r SState)Get(c *gin.Context)  {
	data, err := GetServerInfo()
	if err!=nil{
		response.Error(c, err)
		return
	}
	response.Success(c, data)
}