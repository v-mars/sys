package role

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/v-mars/frame/db"
	"github.com/v-mars/frame/pkg/convert"
	"github.com/v-mars/frame/pkg/utils"
	"github.com/v-mars/frame/response"
	"github.com/v-mars/sys/app/models/sys"
	"github.com/v-mars/sys/app/sys/user"
	"github.com/v-mars/sys/app/utils/casbin"
	"gorm.io/gorm"
	"strings"
)

type IRole interface {
	GetAll(c *gin.Context)
	Query(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetUserRoles(username string) ([]string, error)
}
var dom = "sys"

type Role struct {
	DB *gorm.DB
}

func NewService(DB *gorm.DB) IRole {
	return Role{DB: DB}
}

// GetAll
// @Tags 角色管理
// @Summary 所有角色
// @Description 角色
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/roles [get]
func (r Role) GetAll(c *gin.Context) {
	var obj []ShowData
	var pageData = response.InitPageData(c, &obj, true)
	o := db.Option{}
	o.Select = "distinct id, name"
	o.Order = "ID DESC"
	o.Scan = true
	err := db.Query(r.DB,&sys.Role{}, o, &pageData)
	if err != nil {
		response.Error(c, err)
		return
	} else {
		response.Success(c, pageData)
	}

}

// Query
// @Tags 角色管理
// @Summary 角色列表
// @Description 角色
// @Produce  json
// @Security ApiKeyAuth
// @Param name query string false "角色名"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/role [get]
func (r Role) Query(c *gin.Context) {
	var obj []sys.Role
	var pageData = response.InitPageData(c, &obj, false)
	type Param struct {
		Name       string `form:"name"` // `form:"name" binding:"required"`
	}
	var param Param
	if err := c.ShouldBindQuery(&param); err!=nil{
		response.ParamFailed(c, err)
		return
	}
	var o db.Option
	o.Where = "name like ?"
	o.Value = append(o.Value, "%"+param.Name+"%")
	o.Select = "distinct id, name, description, by_update, ctime, mtime"
	o.Preloads = []string{"Permissions"}
	o.Order = "ID DESC"
	//o.Scan = true
	err := db.Query(r.DB,&obj, o, &pageData)
	if err != nil {
		response.Error(c, err)
		return
	} else {
		response.Success(c, pageData)
	}
}

// Create
// @Tags 角色管理
// @Summary 创建角色
// @Description 角色
// @Produce  json
// @Security ApiKeyAuth
// @Param payload body  sys.Role true "参数信息"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/role [post]
func (r Role) Create(c *gin.Context) {
	u, err:= user.GetUserValue(c)
	if err!=nil{
		response.Error(c, err)
		return
	}
	var obj PostSchema
	var newRow sys.Role
	if err := c.ShouldBindJSON(&obj); err!=nil{
		response.ParamFailed(c, err)
		return
	}
	newRow.Name = obj.Name
	newRow.Description = obj.Description
	newRow.ByUpdate = u.Nickname
	tx :=r.DB.Begin()
	defer func() {tx.Rollback()}()
	if err:= tx.Model(&sys.Permission{}).Where("id in (?)", obj.Permissions).Select(
		"id,name,path,method").Scan(&newRow.Permissions).Error;err!=nil{
		response.Error(c, err)
		return
	}
	for _,v:=range newRow.Permissions{
		_, err = casbin.Enforcer.AddPolicy(v.Name, dom,v.Path, v.Method)
		if err!=nil{
			response.Error(c, err)
			return
		}
	}
	if err:= db.Create(tx, &newRow,true);err!=nil{
		response.Error(c, err)
		return
	}
	response.CreateSuccess(c, obj)
}

// Update
// @Tags 角色管理
// @Summary 更新角色
// @Description 角色
// @Produce  json
// @Security ApiKeyAuth
// @Param payload body  sys.Role true "参数信息"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/role [put]
func (r Role) Update(c *gin.Context) {
	u, err:= user.GetUserValue(c)
	if err!=nil{
		response.Error(c, err)
		return
	}

	var obj PutSchema
	if err := c.ShouldBindJSON(&obj); err!=nil{
		response.Error(c, err)
		return
	}
	obj.ByUpdate = &u.Nickname

	var MapData map[string]interface{}
	if MapData,err=convert.StructToMap(obj); err!=nil{
		response.Error(c, err)
		return
	}
	tx :=r.DB.Begin()
	defer func() {tx.Rollback()}()
	var ass []db.Association
	if obj.Permissions != nil{
		var permissions []sys.Permission
		if err:= tx.Model(&sys.Permission{}).Where("id in (?)", *obj.Permissions).Select("id,name,path,method").Scan(&permissions).Error;err!=nil{
			response.Error(c, err)
			return
		}
		ass =append(ass, db.Association{Column: "Permissions", Values: &permissions})
		var role sys.Role
		if err:=tx.First(&role, obj.ID).Error;err!=nil{
			response.Error(c, err)
			return
		}

		var newStrArray []string
		for _,v:=range permissions{
			newStrArray=append(newStrArray, fmt.Sprintf("%s:%s:%s:%s",role.Name, dom,v.Path, v.Method))
			_, err = casbin.Enforcer.AddPolicy(role.Name, dom,v.Path, v.Method)
			if err!=nil{
				response.Error(c, err)
				return
			}
		}
		existsList := casbin.Enforcer.GetPermissionsForUser(role.Name, dom)

		//existsList, err := casbin.Enforcer.GetImplicitPermissionsForUser(role.Name, dom)
		//if err!=nil{
		//	response.Error(c, err)
		//	return
		//}
		//fmt.Println("existsList:", existsList)
		for _,v:=range existsList {
			var strTmp = strings.Join(v, ":")
			if len(v)==4 && !utils.InOfStr(strTmp, newStrArray){
				_, err = casbin.Enforcer.RemovePolicy(role.Name, dom,v[2], v[3])
				if err!=nil{
					response.Error(c, err)
					return
				}
			}
		}

	}
	if err:= db.UpdateById(tx, &sys.Role{},obj.ID,MapData,ass, true);err!=nil{
		response.Error(c, err)
		return
	}

	response.UpdateSuccess(c)
}

// Delete
// @Tags 角色管理
// @Summary 删除角色
// @Description 角色
// @Produce  json
// @Security ApiKeyAuth
// //@Param payload body [] true "id list"
// //@Param id path int true "角色id"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/role [delete]
// //@Router /api/v1/sys/role/{id} [delete]
func (r Role) Delete(c *gin.Context) {
	var obj sys.Role
	var data map[string][]int
	if err:= c.ShouldBindJSON(&data);err!=nil{
		response.Error(c, err)
		return
	}
	tx :=r.DB.Begin()
	defer func() {tx.Rollback()}()

	var deleteList []sys.Role
	if err:=tx.Model(&sys.Role{}).Find(&deleteList,"id in (?)", data["rows"]).Error;err!=nil{
		response.Error(c, err)
		return
	}
	for _,v:=range deleteList{
		if _, err := casbin.Enforcer.RemoveFilteredNamedPolicy("p", 0, v.Name, dom);err!=nil{
			response.Error(c, err)
			return}
	}
	for _,_id := range data["rows"]{
		if err:= db.DeleteById(tx, &obj, _id, []string{"Permissions"}, true); err!=nil{
			response.Error(c, err)
			return
		}
	}

	response.DeleteSuccess(c)
}

func (r Role) GetUserRoles(username string) ([]string, error) {
	var roles []string
	var sql = "SELECT DISTINCT role.name FROM role " +
		"LEFT JOIN usergroup_roles ON usergroup_roles.role_id = role.id " +
		"LEFT JOIN usergroup ON usergroup.id = usergroup_roles.usergroup_id " +
		"LEFT JOIN usergroup_users ON usergroup_users.usergroup_id = usergroup.id " +
		"LEFT JOIN user ON user.id = usergroup_users.user_id " +
		"LEFT JOIN user_roles ON user_roles.role_id = role.id " +
		"LEFT JOIN user AS T7 ON user_roles.user_id = T7.id " +
		"WHERE (user.username = ? or T7.username= ?)"
	var o db.Option
	o.Value = append(o.Value, username, username)
	o.Pluck = "role.name"
	if err:= db.Raw(&sys.Role{},sql, o, &roles);err!=nil{
		return nil, err
	}else {
		return roles, nil
	}
}

