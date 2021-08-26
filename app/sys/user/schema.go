package user

import (
	"github.com/v-mars/frame/db"
	"github.com/v-mars/sys/app/models"
	"github.com/v-mars/sys/app/models/name"
)

var tbName = name.User

type AllUser struct {
	ID        uint     `json:"id"`
	Nickname  string   `json:"nickname"`
	Username  string   `json:"username"`
}

func (AllUser) TableName() string {
	return tbName
}

type InfoUser struct {
	ID        uint     `json:"id"`
	Name      string   `json:"name"`
	Nickname  string   `json:"nickname"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
}

func (InfoUser) TableName() string {
	return tbName
}

type Role struct {
	ID        uint     `json:"id"`
	Name      string   `json:"name"`
	Title     string   `json:"title"`
}

type User struct {
	ID         uint         `json:"id"`
	Nickname   string       `json:"nickname"`
	Username   string       `json:"username"`
	Password   string       `json:"password"`
	Email      string       `json:"email"`
	Phone      string       `json:"phone"`
	UserType   string       `json:"user_type"`
	Path       models.IntNestArray `json:"path"`
	Status     bool         `json:"status"`
	ByUpdate   string       `json:"by_update"`
	Roles      []Role       `gorm:"many2many:user_roles;" json:"roles"`
	Ctime      db.JSONTime `json:"ctime" form:"ctime"`
	Mtime      db.JSONTime `json:"mtime" form:"mtime"`
}

func (User) TableName() string {
	return tbName
}

type PostSchema struct {
	Nickname    string  `json:"nickname" binding:"required"`
	Username    string  `json:"username" binding:"required"`
	Password    string  `json:"password" binding:"required"`
	Email       string  `json:"email" binding:"required"`
	Path        models.IntNestArray `json:"path"`
	Phone       string  `json:"phone"`
	Roles       []uint  `json:"roles" form:"roles"`
}

func (PostSchema) TableName() string {
	return tbName
}

type PutSchema struct {
	//ID       uint   `json:"id" binding:"required"`
	Nickname *string `json:"nickname,omitempty"`
	Username *string `json:"username,omitempty"`
	//Password string `json:"password" binding:"required"`
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Path     *models.IntNestArray `json:"path"`
	Status   *bool   `json:"status,omitempty"`
	Roles    *[]uint `json:"roles,omitempty"`
	ByUpdate string `json:"by_update,-"`
}

func (PutSchema) TableName() string {
	return tbName
}

type RestPass struct {
	ID       uint   `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePass struct {
	OldPassword string `json:"old_password" binding:"required"`
	Password    string `json:"password" binding:"required"`
	RePassword  string `json:"re_password" binding:"required"`
}

type ChangeProfile struct {
	Nickname *string `json:"nickname,omitempty"`
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	ByUpdate string `json:"by_update,-"`
}

func (ChangeProfile) TableName() string {
	return tbName
}