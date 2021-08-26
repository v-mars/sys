package api

import (
	"github.com/v-mars/sys/app/models/name"
)

var tbName = name.Api

type ShowGroupData struct {
	Group       string `json:"group" form:"group"`
}

func (ShowGroupData) TableName() string {
	return tbName
}


type ShowData struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" form:"name"`
	Url         string `json:"url"`
	Method      string `json:"Method"`
}
func (ShowData) TableName() string {
	return tbName
}

type PostSchema struct {
	Name     string `json:"name" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Group    string `json:"group"`
	Disabled bool   `json:"disabled"`
	Path     string `json:"path"`
	Method   string `json:"Method"`
	//ByUpdate    string `json:"by_update,-"`
}
func (PostSchema) TableName() string {
	return tbName
}

type PutSchema struct {
	Name     *string `json:"name"`
	Title    *string `json:"title"`
	Group    *string `json:"group"`
	Disabled *bool   `json:"disabled"`
	Path     *string `json:"path"`
	Method   *string `json:"Method"`
	//ByUpdate    *string `json:"by_update,-"`
}

func (PutSchema) TableName() string {
	return tbName
}

type DeleteSchema struct {
	Rows []uint `json:"rows"`
}