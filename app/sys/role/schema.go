package role

import (
	"github.com/v-mars/sys/app/models"
	"github.com/v-mars/sys/app/models/name"
)


var tbName = name.Role

type ShowData struct {
	ID          uint   `json:"id" form:"id"`
	Name        string `json:"name" form:"name"`
}
func (ShowData) TableName() string {
	return tbName
}

type PostSchema struct {
	Title       string `json:"title" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ByUpdate    string `json:"by_update,-"`
	Apis        []uint `json:"apis,omitempty"`
	ParentID    uint   `json:"parent_id"`
	Sort        int    `json:"sort"`
	Path        models.IntArray  `json:"path"`
}
func (PostSchema) TableName() string {
	return tbName
}

type PutSchema struct {
	ID          uint             `json:"id" binding:"required"`
	Title       *string          `json:"title"`
	Name        *string          `json:"name"`
	Description *string          `json:"description"`
	ByUpdate    *string          `json:"by_update,-"`
	Apis        *[]uint          `json:"apis,omitempty"`
	ParentID    *uint            `json:"parent_id"`
	Sort        *int             `json:"sort"`
	Path        *models.IntArray `json:"path"`
}

func (PutSchema) TableName() string {
	return tbName
}

type DeleteSchema struct {
	Rows []uint `json:"rows"`
}