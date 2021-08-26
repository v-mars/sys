package role

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/v-mars/frame/pkg/convert"
	"github.com/v-mars/frame/pkg/logger"
	"github.com/v-mars/frame/pkg/utils"
	"github.com/v-mars/frame/response"
	//"github.com/v-mars/sys/app/models"
	"github.com/v-mars/frame/db"
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
	//var obj []ShowData
	//var pageData = response.InitPageData(c, &obj, true)
	//o := models.Option{}
	//o.Select = "distinct id, name"
	//o.Order = "ID DESC"
	//o.Scan = true
	//err := models.Query(r.DB,&tbs.Role{}, o, &pageData)

	var api sys.Role
	tree, err := api.GetListTree(r.DB)
	if err != nil {
		response.Error(c, err)
		return
	} else {
		response.Success(c, tree)
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
	o.Where = "(parent_id = 0 or parent_id is null) and name like ?"
	o.Value = append(o.Value, "%"+param.Name+"%")
	o.Preloads = []string{"Apis"}
	o.Order = "sort asc"

	var api sys.Role
	tree, err := api.GetMapTree(r.DB)
	if err!=nil{
		response.Error(c, err)
		return
	}

	var data = make([]map[string]interface{},0)
	if err = db.Query(r.DB,&obj, o, &pageData);err != nil {
		response.Error(c, err)
		return
	} else {
		for _,v := range obj{
			tree[v.ID]["apis"] = v.Apis
			if tree[v.ID] != nil{
				data = append(data, tree[v.ID])
			} else {
				var t = map[string]interface{}{}
				err = utils.AnyToAny(tree[v.ID], &t)
				if err != nil {
					response.Error(c, err)
					return
				}
				data = append(data, t)
			}
		}
		pageData.List = data
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
	if err = c.ShouldBindJSON(&obj); err!=nil{
		response.ParamFailed(c, err)
		return
	}
	newRow.Name = obj.Name
	newRow.Description = obj.Description
	newRow.ByUpdate = u.Nickname
	newRow.ParentID = obj.ParentID
	newRow.Path = obj.Path
	tx :=r.DB.Begin()
	defer func() {tx.Rollback()}()
	if err= tx.Model(&sys.Api{}).Where("id in (?)", obj.Apis).Select(
		"id,name,path,method").Scan(&newRow.Apis).Error;err!=nil{
		response.Error(c, err)
		return
	}
	for _,v:=range newRow.Apis{
		_, err = casbin.Enforcer.AddPolicy(v.Name, dom,v.Path, v.Method)
		if err!=nil{
			response.Error(c, err)
			return
		}
	}
	if err= db.Create(tx, &newRow,true);err!=nil{
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
	if err = c.ShouldBindJSON(&obj); err!=nil{
		response.Error(c, err)
		return
	}
	obj.ByUpdate = &u.Nickname
	if obj.ID == *obj.ParentID {
		*obj.ParentID=0
	}
	var MapData map[string]interface{}
	if MapData,err=convert.StructToMap(obj); err!=nil{
		response.Error(c, err)
		return
	}
	tx :=r.DB.Begin()
	defer func() {tx.Rollback()}()
	var ass []db.Association
	var role sys.Role
	if err=tx.First(&role, obj.ID).Error;err!=nil{
		response.Error(c, err)
		return
	}
	if obj.Path!=nil{
		if len(*obj.Path)>9{
			response.Error(c, fmt.Errorf("角色层级深度不能超过10级"))
			return
		}
		MapData["path"] = obj.Path
	}

	if obj.Apis != nil{
		var apis []sys.Api
		if err= tx.Model(&sys.Api{}).Where("id in (?)", *obj.Apis).Select("id,name,path,method").Scan(&apis).Error;err!=nil{
			response.Error(c, err)
			return
		}
		ass =append(ass, db.Association{Column: "Apis", Values: &apis})

		var newStrArray []string
		for _,v:=range apis{
			newStrArray=append(newStrArray, fmt.Sprintf("%s:%s:%s:%s",role.Name, dom,v.Path, v.Method))
			_, err = casbin.Enforcer.AddPolicy(role.Name, dom,v.Path, v.Method)
			if err!=nil{
				response.Error(c, err)
				return
			}
		}

		existsList := casbin.Enforcer.GetPermissionsForUser(role.Name, dom)
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

	if err= db.UpdateById(tx, &sys.Role{},obj.ID,MapData,ass, true);err!=nil{
		response.Error(c, err)
		return
	}

	if obj.ParentID!=nil && *obj.ParentID != 0 {
		var parentRole sys.Role
		if err=r.DB.First(&parentRole, *obj.ParentID).Error;err!=nil{
			response.Error(c, err)
			return
		}
		if role.ParentID != 0 && role.ParentID != *obj.ParentID {
			_, err = casbin.Enforcer.RemoveFilteredNamedGroupingPolicy("g", 1, role.Name, parentRole.Name, dom)
			if err != nil {
				response.Error(c, err)
				return
			}
		}
		_, err = casbin.Enforcer.AddGroupingPolicy(parentRole.Name, role.Name, dom)
		if err != nil {
			response.Error(c, err)
			return
		}
	}else {
		if role.ParentID != 0 {
			_, err = casbin.Enforcer.RemoveFilteredNamedGroupingPolicy("g", 1, role.Name, "", dom)
			if err != nil {
				response.Error(c, err)
				return
			}
		}
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
	var childList []sys.Role
	if err:=tx.Model(&sys.Role{}).Find(&deleteList,"id in (?)", data["rows"]).Error;err!=nil{
		response.Error(c, err)
		return
	}
	if err:=tx.Model(&sys.Role{}).Find(&childList,"parent_id in (?)", data["rows"]).Error;err!=nil{
		response.Error(c, err)
		return
	}
	if len(childList)>0{
		response.Error(c, fmt.Errorf("角色包含子角色不能删除"))
		return
	}


	for _,v:=range deleteList{
		if err:=tx.Model(&sys.Role{}).Find(&deleteList,"id in (?)", data["rows"]).Error;err!=nil{
			response.Error(c, err)
			return
		}
		if _, err := casbin.Enforcer.RemoveFilteredNamedPolicy("p", 0, v.Name, dom);err!=nil{
			response.Error(c, err)
			return}
		if _, err := casbin.Enforcer.RemoveFilteredNamedGroupingPolicy("g", 0, v.Name, "", dom);err!=nil{
			response.Error(c, err)
			return
		}
		if _, err := casbin.Enforcer.RemoveFilteredNamedGroupingPolicy("g",1, v.Name, dom);err!=nil{
			response.Error(c, err)
			return
		}
	}
	logger.Logger.Printf("角色关联api删除成功")
	for _,_id := range data["rows"]{
		if err:= db.DeleteById(tx, &obj, _id, []string{"Apis"}, true); err!=nil{
			response.Error(c, err)
			return
		}
	}
	logger.Logger.Printf("角色删除成功")
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

