package usergroup

import (
	"github.com/v-mars/frame/db"
	"github.com/v-mars/sys/app/models/name"
)

var tbName = name.UserGroup

type Role struct {
	ID        uint  	    `json:"id"`
	Name      string        `json:"name"`
	Title     string        `json:"title"`
}

type User struct {
	ID        uint     `json:"id"`
	Nickname  string   `json:"nickname"`
	Username  string   `json:"username"`
}

type Usergroup struct {
	ID          uint         `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Users       []User       `gorm:"many2many:usergroup_users;" json:"users"`
	Roles       []Role       `gorm:"many2many:usergroup_roles;" json:"roles"`
	ByUpdate    string       `json:"by_update"`
	OwnerID     *uint        `json:"owner_id"`
	OwnerName   string       `json:"owner_name"`
	Ctime       db.JSONTime `json:"ctime"`
	Mtime       db.JSONTime `json:"mtime"`
}

func (Usergroup) TableName() string {
	return tbName
}

type All struct {
	ID          uint         `json:"id"`
	Name        string       `json:"name"`
}

func (All) TableName() string {
	return tbName
}

type PostSchema struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	//OwnerID     *uint   `json:"owner_id"`
	Users       []uint  `json:"users" form:"users"`
	Roles       []uint  `json:"roles" form:"roles"`
}

func (PostSchema) TableName() string {
	return tbName
}

type PutSchema struct {
	ID          uint    `json:"id" binding:"required"`
	Name        *string `json:"name" binding:"required"`
	Description *string `json:"description"`
	Email       *string `json:"email,omitempty"`
	OwnerID     *uint   `json:"owner_id"`
	Users       *[]uint `json:"users" form:"users"`
	Roles       *[]uint `json:"roles,omitempty"`
	ByUpdate    string  `json:"by_update,-"`
}

func (PutSchema) TableName() string {
	return tbName
}