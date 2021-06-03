package portal

import (
	"github.com/v-mars/frame/db"
	"github.com/v-mars/sys/app/models/name"
)

var tbName = name.Portal

type ShowData struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" form:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Icon        string `json:"icon"`
	Type        string `json:"type"`
	IsFavor     uint   `json:"is_favor"`
	db.BaseByUpdate
	db.BaseTime
}

func (ShowData) TableName() string {
	return tbName
}

type ShowPortalType struct {
	Type        string `json:"type"`
}

func (ShowPortalType) TableName() string {
	return tbName
}

type PostSchema struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Icon        string `json:"icon"`
	Type        string `json:"type"`
	ByUpdate    string `json:"by_update,-"`
}

func (PostSchema) TableName() string {
	return tbName
}

type PutSchema struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Url         *string `json:"url"`
	Icon        *string `json:"icon"`
	Type        *string `json:"type"`
	ByUpdate    string `json:"by_update,-"`
}

func (PutSchema) TableName() string {
	return tbName
}

type DeleteSchema struct {
	Rows []uint `json:"rows"`
}