package dept

import (
	"github.com/v-mars/sys/app/models"
	"github.com/v-mars/sys/app/models/name"
)

var tbName = name.Dept

type ShowData struct {
	ID          uint   `json:"id" form:"id"`
	Name        string `json:"name" form:"name"`
	Code        string `json:"code" form:"code"`
}
func (ShowData) TableName() string {
	return tbName
}

type PostSchema struct {
	Name     string `json:"name" binding:"required"`
	Code     string `json:"code" binding:"required"`
	ParentID uint   `json:"parent_id"`
	Sort     int    `json:"sort"`
	Leader   string `json:"leader"`
	Status   bool   `json:"status"`
	Path     models.IntArray  `json:"path"`
}
func (PostSchema) TableName() string {
	return tbName
}

type PutSchema struct {
	//ID          uint   `json:"id" binding:"required"`
	Name     *string `json:"name"`
	Code     *string `json:"code"`
	ParentID *uint   `json:"parent_id"`
	Sort     *int    `json:"sort"`
	Leader   *string `json:"leader"`
	Status   *bool   `json:"status"`
	Path     *models.IntArray  `json:"path"`
}

func (PutSchema) TableName() string {
	return tbName
}

type DeleteSchema struct {
	Rows []uint `json:"rows"`
}