package sys

import "github.com/v-mars/frame/db"

var PermissionName = "user_group"

type Permission struct {
	db.BaseID
	Name   string `gorm:"type:varchar(128);comment:'名称'" json:"name" form:"name"`
	Path   string `gorm:"type:varchar(128);not null;unique_index:idx_name_code;comment:'路由路径'" json:"path" form:"path"`
	Method string `gorm:"type:varchar(64);not null;unique_index:idx_name_code;comment:'请求方法'" json:"method" form:"method"`
	db.BaseTime
}

func (Permission) TableName() string {
	return PermissionName
}