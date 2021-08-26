package usergroup

import (
	"github.com/gin-gonic/gin"
	"github.com/v-mars/frame/db"
	"github.com/v-mars/frame/pkg/convert"
	"github.com/v-mars/frame/response"
	"github.com/v-mars/sys/app/models/sys"
	"github.com/v-mars/sys/app/sys/user"
	"gorm.io/gorm"
)


type IUserGroup interface {
	Query(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetAll(c *gin.Context)
	Option() db.Option
}
type UserGroup struct {
	DB *gorm.DB
}

func NewService(DB *gorm.DB) IUserGroup {
	return UserGroup{DB: DB}
}

// GetAll
// @Tags 用户组管理
// @Summary 所有用户组
// @Description 用户组
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/usergroups [get]
func (r UserGroup) GetAll(c *gin.Context) {

	var obj []All
	var pageData = response.PageDataList{Page: 1,PageSize:0,List:&obj}
	o := db.Option{}
	o.Select = "distinct usergroup.id, name"
	o.Order = "ID DESC"
	o.Scan = true
	err := db.Query(r.DB,&sys.Usergroup{}, o, &pageData)
	if err != nil {
		response.Error(c, err)
		return
	} else {
		response.Success(c, pageData)
	}

}


// @Tags 用户组管理
// @Summary 用户组列表
// @Description 用户组
// @Produce  json
// @Security ApiKeyAuth
// @Param page query int false "page"
// @Param pageSize query int false "pageSize"
// @Param name query string false "用户组显示名"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/usergroup [get]
func (r UserGroup) Query(c *gin.Context) {
	var obj []Usergroup
	var pageData = response.InitPageData(c, &obj, false)
	type Param struct {
		Name       string `form:"name"` // `form:"name" binding:"required"`
	}
	var param Param

	if err := c.ShouldBindQuery(&param); err!=nil{
		response.ParamFailed(c, err)
		return
	}
	o := r.Option()
	o.Where = "name like ?"
	o.Value = append(o.Value, "%"+param.Name+"%")
	//o.Select = "distinct usergroup.id, name, description, usergroup.ctime, usergroup.mtime, " +
	//	"usergroup.by_update, user.nickname as owner_name, user.id as owner_id"
	o.Order = "ID DESC"
	o.Preloads = []string{"Users", "Roles"}
	//o.Scan = true
	err := db.Query(r.DB,&Usergroup{}, o, &pageData)
	if err != nil {
		response.Error(c, err)
	} else {
		response.Success(c, pageData)
	}
}

// Create
// @Tags 用户组管理
// @Summary 创建用户组
// @Description 用户组
// @Produce  json
// @Security ApiKeyAuth
// @Param payload body  sys.User true "参数信息"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/usergroup [post]
func (r UserGroup) Create(c *gin.Context) {
	u, err:= user.GetUserValue(c)
	if err!=nil{
		response.Error(c, err)
		return
	}
	var obj PostSchema
	if err := c.ShouldBindJSON(&obj); err != nil{
		response.Error(c, err)
		return
	}
	var newRow sys.Usergroup
	newRow.Name = obj.Name
	newRow.OwnerID = &u.ID
	tx :=r.DB.Begin()
	defer func() {tx.Rollback()}()
	if err:= tx.Where("id in (?)", obj.Users).Find(&newRow.Users).Error; err!= nil{
		response.Error(c, err)
		return
	}
	if err:= tx.Where("id in (?)", obj.Roles).Find(&newRow.Roles).Error; err!= nil{
		response.Error(c, err)
		return
	}
	newRow.ByUpdate = u.Nickname

	if err:= db.Create(tx, &newRow,true);err!=nil{
		response.Error(c, err)
		return
	}
	response.CreateSuccess(c)
}

// Update
// @Tags 用户组管理
// @Summary 更新用户组
// @Description 用户组
// @Produce  json
// @Security ApiKeyAuth
// @Param payload body  sys.User true "参数信息"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/usergroup [put]
func (r UserGroup) Update(c *gin.Context) {
	u, err:= user.GetUserValue(c)
	if err!=nil{
		response.Error(c, err)
		return
	}
	var obj PutSchema
	if err := c.ShouldBindJSON(&obj); err != nil{
		response.Error(c, err)
		return
	}
	obj.ByUpdate = u.Nickname
	var MapData map[string]interface{}
	if MapData,err=convert.StructToMap(obj); err!=nil{
		response.Error(c, err)
		return
	}
	//fmt.Println(MapData)
	//panic("dsfadsfadsf")
	tx :=r.DB.Begin()
	defer func() {tx.Rollback()}()
	var ass []db.Association
	if obj.Roles!=nil{
		var roles []sys.Role
		if err:= tx.Model(&sys.Role{}).Where("id in (?)", *obj.Roles).Select("id").Scan(&roles).Error;err!=nil{
			response.Error(c, err)
			return
		}
		ass =append(ass, db.Association{Column: "Roles", Values: &roles})
	}

	if obj.Users!=nil{
		var users []sys.User
		if err:= tx.Where("id in (?)", *obj.Users).Find(&users).Error;err!=nil{
			response.Error(c, err)
			return
		}
		ass =append(ass, db.Association{Column: "Users", Values: &users})
	}

	if err:= db.UpdateById(tx, &sys.Usergroup{},obj.ID,MapData,ass, true);err!=nil{
		response.Error(c, err)
		return
	}
	response.UpdateSuccess(c)
}

// Delete
// @Tags 用户组管理
// @Summary 删除用户组
// @Description 用户组
// @Produce  json
// @Security ApiKeyAuth
// //@Param payload body [] true "id list"
// //@Param id path int true "用户组id"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/usergroup [delete]
// //@Router /api/v1/sys/usergroup/{id} [delete]
func (r UserGroup) Delete(c *gin.Context) {
	var obj sys.Usergroup
	var data map[string][]int
	if err := c.ShouldBindJSON(&data); err!=nil{
		response.Error(c, err)
		return
	}

	tx :=r.DB.Begin()
	defer func() {tx.Rollback()}()
	//time.Sleep(time.Second*30)
	for _, i := range data["rows"]{
		//var o = model.Option{Where: "id = ?", Value: []interface{}{i}}
		if err:= db.DeleteById(tx, &obj, i, []string{"Roles", "Users"}, true); err!=nil{
			response.Error(c, err)
			return
		}
	}

	response.DeleteSuccess(c)
}

func (r UserGroup) Option() db.Option {
	var o db.Option
	o.Select = "distinct usergroup.id, name, description, usergroup.ctime, usergroup.mtime, " +
		"usergroup.by_update, user.nickname as owner_name, user.id as owner_id"
	o.Joins = "left join `user` on `user`.id=`usergroup`.owner_id "
	return o
}