package sys

import (
	"github.com/v-mars/frame/db"
	"github.com/v-mars/sys/app/models"
	"github.com/v-mars/sys/app/models/name"
	_ "gorm.io/driver/mysql"
)

// User
// (1, '本地'),
// (2, 'LDAP'),
// (3, 'SSO'),
type User struct {
	db.BaseID
	Nickname   string      `gorm:"type:varchar(128);index:idx_name_code;comment:'显示名'" json:"nickname" form:"nickname"`
	Username   string      `gorm:"type:varchar(128);unique_index:idx_name_code;comment:'用户名'" json:"username" form:"username"` // `unique_index` also works
	Password   string      `gorm:"type:varchar(128);comment:'密码'" json:"password" form:"password"`
	Email      string      `gorm:"type:varchar(156);comment:'邮箱'" json:"email" form:"email"`
	Phone      string      `gorm:"type:varchar(32);comment:'电话'" json:"phone" form:"phone"`
	UserType   string      `gorm:"type:varchar(16);default:'local';comment:'用户类型:local,ldap,sso'" json:"user_type" form:"user_type"`
	Path       models.IntNestArray `json:"path" gorm:"type:text;comment:'节点路径'"`
	//UserTypeID uint        `gorm:"column:user_type_id;type:varchar(156);default:1" json:"user_type_id" form:"user_type_id"` // 1 本地，2 ldap
	Roles      []Role `gorm:"many2many:user_roles;association_autoupdate:false;association_autocreate:false;constraint:OnDelete:CASCADE" json:"roles"` // Many-To-Many
	//AccessSecret string    `gorm:"type:varchar(156);" json:"access_secret" form:"access_secret"`
	Status     bool        `gorm:"type:varchar(32);default:true;comment:'状态'" json:"status" form:"status"`
	db.BaseByUpdate
	db.BaseTime
}
func (User) TableName() string {
	return name.User
}

type Usergroup struct {
	db.BaseID
	Name        string `gorm:"type:varchar(128);comment:'名称'" json:"name" form:"name"`
	Description string `gorm:"type:longtext;comment:'描述'" json:"description" form:"description"`
	Users       []User `gorm:"many2many:usergroup_users;association_autoupdate:false;association_autocreate:false;comment:'关联用户';constraint:OnDelete:CASCADE" json:"users" form:"users"`
	Roles       []Role `gorm:"many2many:usergroup_roles;association_autoupdate:false;association_autocreate:false;comment:'关联角色';constraint:OnDelete:CASCADE" json:"roles" form:"roles"`
	OwnerID     *uint  `gorm:"index:owner_id;default:null;comment:'用户组拥有者'" json:"owner_id" form:"owner_id"`
	db.BaseByUpdate
	db.BaseTime
}
func (Usergroup) TableName() string {
	return name.UserGroup
}

type Portal struct {
	db.BaseID
	Name        string `gorm:"type:varchar(128);not null;unique_index:idx_name_code;comment:'名称'" json:"name" form:"name"`
	Description string `gorm:"type:longtext;comment:'描述'" json:"description" form:"description"`
	Url         string `gorm:"type:varchar(256);not null;comment:'URL'" json:"url" form:"url"`
	Icon        string `gorm:"type:varchar(256);comment:'图标'" json:"icon" form:"icon"`
	Type        string `gorm:"type:varchar(128);default:'default';comment:'类型'" json:"type" form:"type"`
	Favors      []User `gorm:"many2many:portal_favors;association_autoupdate:false;association_autocreate:false;comment:'portal收藏'" json:"favors"`
	db.BaseByUpdate
	db.BaseTime
}
func (Portal) TableName() string {
	return name.Portal
}


type Property struct {
	db.BaseID
	Name     string   `gorm:"type:varchar(128);unique_index:idx_name_code;comment:'名称'" json:"name" form:"name"`
	K        string   `gorm:"column:k;type:varchar(128);comment:'key'" json:"k" form:"k"`
	V        string   `gorm:"column:v;type:varchar(256);comment:'value'" json:"v" form:"v"`
	db.BaseByUpdate
	db.BaseTime
}

type TreeV1 struct {
	db.BaseID
	Name     string  `gorm:"unique_index;type:varchar(128);comment:'节点名'" json:"name" form:"name"`
	ParentID *uint   `gorm:"column:parent_id;index:parent_id;comment:'父节点'" json:"parent_id" form:"parent_id"`
	Mark     string  `gorm:"type:varchar(128);comment:'标识'" json:"mark" form:"mark"`
}

type TreeNode struct {
	db.BaseID
	ParentID uint   `gorm:"index:parent_id;comment:'父节点'" json:"parent_id" form:"parent_id"`
	Lft      int    `gorm:"comment:'left'" json:"lft" form:"lft"`
	Rgt      int    `gorm:"comment:'Right'" json:"rgt" form:"rgt"`
	Name     string `gorm:"unique_index;type:varchar(128);comment:'节点名'" json:"name" form:"name"`
	Mark     string `gorm:"index:mark;type:varchar(128);not null;comment:'标识'" json:"mark" form:"mark"`
}




// 以下sql中需要传的值全用???表示
// 根据节点id获取此节点所有子孙节点
// select * from tree_table where
// left > (select left from tree_table where id=???) and
// right < (select right from tree_table where id=???)

// 根据节点id获取此节点的所有子孙节点(包含自己)
// select * from tree_table where
// left >= (select left from tree_table where id=???) and
// right <= (select right from tree_table where id=???)

// 根据节点id获取此节点的所有上级节点
// select * from tree_table where
// left < (select left from tree_table where id=???) and
// right > (select right from tree_table where id=???)

// 根据节点id获取此节点的所有上级节点(包括自己)
// select * from tree_table where
// left <= (select left from tree_table where id=???) and
// right >= (select right from tree_table where id=???)