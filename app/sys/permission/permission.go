package permission

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/v-mars/frame/db"
	"github.com/v-mars/frame/pkg/convert"
	"github.com/v-mars/frame/response"
	"github.com/v-mars/sys/app"
	"github.com/v-mars/sys/app/models/sys"
	"github.com/v-mars/sys/app/utils/casbin"
	"gorm.io/gorm"
	"strings"
)

type IPermission interface {
	app.CommonInterfaces
	UpdatePermission(conditions []sys.Permission) error
}
var dom = "sys"
type Permission struct {
	DB *gorm.DB
}

func NewService(DB *gorm.DB) IPermission {
	return Permission{DB: DB}
}

// Get
// @Tags 权限管理
// @Summary 权限详细
// @Description Permission
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "ID"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/permission/{id} [get]
func (r Permission) Get(c *gin.Context) {
	_id := c.Params.ByName("id")
	var obj sys.Permission
	o := r.Option()
	o.Where = "id = ?"
	o.Value = append(o.Value, _id)
	o.First = true
	o.NullError = true
	err := db.Get(r.DB,&obj, o, &obj)
	if err!= nil {
		response.Error(c, err)
		return
	} else {
		response.Success(c,&obj)
		return
	}
}

// Query
// @Tags 权限管理
// @Summary 权限列表
// @Description Permission
// @Produce  json
// @Security ApiKeyAuth
// @Param name query string false "权限名"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/permission [get]
func (r Permission) Query(c *gin.Context) {
	var obj []sys.Permission
	var pageData = response.InitPageData(c, &obj, false)
	type Param struct {
		Name   string `form:"name"`   // `form:"name" binding:"required"`
		Method string `form:"method"` // `form:"name" binding:"required"`
		Path   string `form:"path"`    // `form:"name" binding:"required"`
	}
	var param Param
	if err := c.ShouldBindQuery(&param); err!=nil{
		response.ParamFailed(c, err)
		return
	}
	var o = r.Option()
	o.Where = "name like ? and method like ? and path like ?"
	o.Value = append(o.Value, "%"+param.Name+"%", "%"+param.Method+"%", "%"+param.Path+"%")
	o.Order = "ID DESC"
	o.Scan = true
	tx :=r.DB.Begin()
	defer func() {tx.Rollback()}()
	err := db.Query(tx,&sys.Permission{}, o, &pageData)
	if err != nil {
		response.Error(c, err)
		return
	} else {
		response.Success(c, pageData)
	}
}

// Create
// @Tags 权限管理
// @Summary 创建权限
// @Description Permission
// @Produce  json
// @Security ApiKeyAuth
// @Param payload body PostSchema true "参数信息"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/permission [post]
func (r Permission) Create(c *gin.Context) {
	//var err error
	//u, err:= user.GetUserValue(c)
	//if err!=nil{
	//	response.Error(c, err)
	//	return
	//}
	var obj PostSchema
	if err := c.ShouldBindQuery(&obj); err!=nil{
		response.ParamFailed(c, err)
		return
	}
	//obj.ByUpdate = u.Nickname
	tx :=r.DB.Begin()
	defer func() {
		tx.Rollback()
	}()

	if err:= db.Create(tx, &obj,true);err!=nil{
		response.Error(c, err)
		return
	}
	response.CreateSuccess(c, obj)
}

// Update
// @Tags 权限管理
// @Summary 更新权限
// @Description Permission
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "ID"
// @Param payload body PutSchema true "参数信息"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/permission/{id} [put]
func (r Permission) Update(c *gin.Context) {
	var err error
	//u, err:= user.GetUserValue(c)
	//if err!=nil{
	//	response.Error(c, err)
	//	return
	//}
	var obj PutSchema
	if err := c.ShouldBindJSON(&obj); err!=nil{
		response.Error(c, err)
		return
	}
	//obj.ByUpdate = &u.Nickname
	_id := c.Params.ByName("id")
	var MapData map[string]interface{}
	if MapData,err=convert.StructToMap(obj); err!=nil{
		response.Error(c, err)
		return
	}
	//fmt.Println("mapData", MapData)
	tx :=r.DB.Begin()
	defer func() {tx.Rollback()}()
	var res sys.Permission
	if err:= db.UpdateById(tx, &res, _id,MapData,nil, true);err!=nil{
		response.Error(c, err)
		return
	}
	response.UpdateSuccess(c, res)
}

// Delete
// @Tags 权限管理
// @Summary 删除权限
// @Description Permission
// @Produce  json
// @Security ApiKeyAuth
// //@Param id path int true "ID"
// @Param payload body DeleteSchema true "参数信息: {rows:[1,2]}"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/permission [delete]
func (r Permission) Delete(c *gin.Context) {
	//_id := c.Params.ByName("id")
	var data map[string][]int
	if err:= c.ShouldBindJSON(&data);err!=nil{
		response.Error(c, err)
		return
	}
	//rows["rows"]
	tx :=r.DB.Begin()
	defer func() {
		tx.Rollback()
	}()
	for _,_id := range data["rows"]{
		if err:= db.DeleteById(tx, &sys.Permission{}, _id, []string{}, true); err!=nil{
			response.Error(c, err)
			return
		}
	}

	response.DeleteSuccess(c)
}


func (r Permission) Option() db.Option {
	var o db.Option
	o.Select = "distinct id,name,path,method,ctime,mtime"
	return o
}

func (r Permission) UpdatePermission(conditions []sys.Permission) error {
	tx :=r.DB.Begin()
	defer func() {tx.Rollback()}()
	var obj []sys.Permission
	var pageData = response.PageDataList{PageNumber: 1,PageSize:0,List:&obj}
	var o = r.Option()
	o.All = true
	var cdsStrList []string
	for _,v:=range conditions{
		cdsStrList = append(cdsStrList, fmt.Sprintf("path='%s' and method='%s'", v.Path, v.Method))
	}
	o.Where = strings.Join(cdsStrList, " or ")
	err := db.Query(tx, &sys.Permission{},o, &pageData)
	if err!=nil{
		return err
	}

	// 过滤新增的列表
	var newRows []sys.Permission
	var existList []string
	var existIdList []uint
	for _,v:=range obj{
		existIdList = append(existIdList, v.ID)
		existList = append(existList, fmt.Sprintf("path='%s' and method='%s'", v.Path, v.Method))
	}

	for _,v:=range conditions{
		var str1 = fmt.Sprintf("path='%s' and method='%s'", v.Path, v.Method)
		if !CheckExist(str1, existList){
			newRows = append(newRows, v)
		}
	}

	// 过滤出需要删除的列表
	var deleteList []sys.Permission
	var deleteIdList []uint
	if err:=tx.Find(&deleteList, "id not in (?)", existIdList).Error;err!=nil{return err}

	for _,v:=range deleteList{
		deleteIdList = append(deleteIdList, v.ID)
		if _, err := casbin.Enforcer.RemoveFilteredNamedPolicy("p", 1, dom, v.Path, v.Method);err!=nil{return err}
	}
	if err:=tx.Exec("DELETE FROM role_permissions WHERE permission_id in (?)", deleteIdList).Error;err!=nil{return err}
	if err:=tx.Delete(&sys.Permission{}, "id in (?)", deleteIdList).Error;err!=nil{return err}
	if len(newRows)>0 {
		err=tx.CreateInBatches(newRows, 100).Error
		//err = BatchSave(tx, newRows)
		if err!=nil{return err}
	}

	if err:=tx.Commit().Error; err!=nil{
		tx.Rollback()
		return err
	}
	return nil
}

// BatchSave 批量插入数据
func BatchSave(db *gorm.DB, perms []sys.Permission) error {
	var buffer bytes.Buffer
	sql := fmt.Sprintf("insert into `%s` (`name`,`path`,`method`) values", tbName)
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	for i, e := range perms {
		if i == len(perms)-1 {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s');", e.Name,e.Path, e.Method))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s'),", e.Name,e.Path, e.Method))
		}
	}
	return db.Exec(buffer.String()).Error
}

func CheckExist(str string, existList []string) bool {
	for _,v:=range existList{
		if str == v{
			return true
		}
	}
	return false
}