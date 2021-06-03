package sys

import (
	"github.com/v-mars/frame/db"
)

var UserGroupName = "user_group"

type UserGroup struct {
	db.BaseID
	Name        string `gorm:"type:varchar(128);comment:'名称'" json:"name" form:"name"`
	Description string `gorm:"type:longtext;comment:'描述'" json:"description" form:"description"`
	Users       []User `gorm:"many2many:usergroup_users;association_autoupdate:false;association_autocreate:false;comment:'关联用户'" json:"users" form:"users"`
	Roles       []Role `gorm:"many2many:usergroup_roles;association_autoupdate:false;association_autocreate:false;comment:'关联角色'" json:"roles" form:"roles"`
	OwnerID     *uint  `gorm:"index:owner_id;default:null;comment:'用户组拥有者'" json:"owner_id" form:"owner_id"`
	db.BaseByUpdate
	db.BaseTime
}


func (UserGroup) TableName() string {
	return UserGroupName
}
