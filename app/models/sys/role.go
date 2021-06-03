package sys

import "github.com/v-mars/frame/db"

var RoleName = "role"

type Role struct {
	db.BaseID
	Name        string       `gorm:"type:varchar(128);unique_index:idx_name_code;comment:'名称'" json:"name" form:"name"`
	Description string       `gorm:"type:longtext;comment:'描述'" json:"description" form:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions;association_autoupdate:false;association_autocreate:false" json:"permissions"`
	db.BaseByUpdate
	db.BaseTime
}

func (Role) TableName() string {
	return RoleName
}
