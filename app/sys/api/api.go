package api

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/v-mars/frame/db"
	"github.com/v-mars/frame/pkg/convert"
	"github.com/v-mars/frame/pkg/utils"
	"github.com/v-mars/frame/response"
	"github.com/v-mars/sys/app"
	"github.com/v-mars/sys/app/models/name"
	"github.com/v-mars/sys/app/models/sys"
	"github.com/v-mars/sys/app/utils/casbin"
	"gorm.io/gorm"
	"strings"
)

type IApi interface {
	app.CommonInterfaces
	GetAllPerm(c *gin.Context)
	GetAllGroup(c *gin.Context)
	UpdateApi(conditions []sys.Api) error
}

var dom = "sys"

type Api struct {
	DB *gorm.DB
}

func NewService(DB *gorm.DB) IApi {
	return Api{DB: DB}
}

// GetAllPerm
// @Tags Api接口管理
// @Summary Api接口详细
// @Description Api
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "ID"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/apis [get]
func (r Api) GetAllPerm(c *gin.Context) {
	var api sys.Api
	tree, err := api.GetListTree(r.DB)
	if err != nil {
		response.Error(c, err)
		return
	} else {
		response.Success(c, tree)
	}
}

// GetAllGroup
// @Tags Api接口管理
// @Summary Api接口Group
// @Description Api
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/api-group [get]
func (r Api) GetAllGroup(c *gin.Context) {
	var obj []ShowGroupData
	var o = r.Option()
	o.Select = "DISTINCT api.group"
	o.Where = "api.group != null or api.group!=''"
	o.Scan = true
	tx :=r.DB.Begin()
	defer func() {
		tx.Rollback()
	}()
	err := db.Get(tx,&sys.Api{}, o, &obj)
	if err != nil {
		response.Error(c, err)
		return
	} else {
		var pageData = response.PageDataList{Page: 1,PageSize:0,List:&obj,Total: int64(len(obj))}
		response.Success(c, pageData)
	}
}

// Get
// @Tags Api接口管理
// @Summary Api接口详细
// @Description Api
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "ID"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/api/{id} [get]
func (r Api) Get(c *gin.Context) {
	_id := c.Params.ByName("id")
	var obj sys.Api
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
// @Tags Api接口管理
// @Summary Api接口列表
// @Description Api
// @Produce  json
// @Security ApiKeyAuth
// @Param name query string false "Api接口名"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/api [get]
func (r Api) Query(c *gin.Context) {
	var obj []sys.Api
	var pageData = response.InitPageData(c, &obj, false)
	type Param struct {
		Title   string `form:"title"`   // `form:"name" binding:"required"`
		Name   string `form:"name"`   // `form:"name" binding:"required"`
		Method string `form:"method"`
		Path   string `form:"path"`
		Groups    []string `form:"groups[]"`
	}
	var param Param
	if err := c.ShouldBindQuery(&param); err!=nil{
		response.ParamFailed(c, err)
		return
	}
	var o = r.Option()
	o.Where = "name like ? and method like ? and path like ?"
	o.Value = append(o.Value, "%"+param.Name+"%", "%"+param.Method+"%", "%"+param.Path+"%")
	if len(param.Groups)>0{
		o.Where = o.Where + " and api.group in (?)"
		o.Value = append(o.Value, param.Groups)
	}
	o.Order = "name asc,path asc"
	tx :=r.DB.Begin()
	defer func() {tx.Rollback()}()
	if err := db.Query(tx,&sys.Api{}, o, &pageData);err != nil {
		response.Error(c, err)
		return
	} else {
		response.Success(c, pageData)
	}
}

// Create
// @Tags Api接口管理
// @Summary 创建Api接口
// @Description Api
// @Produce  json
// @Security ApiKeyAuth
// @Param payload body PostSchema true "参数信息"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/api [post]
func (r Api) Create(c *gin.Context) {
	var obj PostSchema
	if err := c.ShouldBindJSON(&obj); err!=nil{
		response.ParamFailed(c, err)
		return
	}
	tx :=r.DB.Begin()
	defer func() {
		tx.Rollback()
	}()
	var newRow = sys.Api{}
	err := utils.AnyToAny(obj, &newRow)
	if err != nil {
		response.Error(c, err)
		return
	}
	if err= db.Create(tx, &newRow,true);err!=nil{
		response.Error(c, err)
		return
	}
	response.CreateSuccess(c, obj)
}

// Update
// @Tags Api接口管理
// @Summary 更新Api接口
// @Description Api
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "ID"
// @Param payload body PutSchema true "参数信息"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/api/{id} [put]
func (r Api) Update(c *gin.Context) {
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
	_id := c.Params.ByName("id")
	var MapData map[string]interface{}
	if MapData,err=convert.StructToMap(obj); err!=nil{
		response.Error(c, err)
		return
	}
	tx :=r.DB.Begin()
	defer func() {tx.Rollback()}()
	var res sys.Api
	if err= db.UpdateById(tx, &res, _id,MapData,nil, true);err!=nil{
		response.Error(c, err)
		return
	}
	response.UpdateSuccess(c, res)
}

// Delete
// @Tags Api接口管理
// @Summary 删除Api接口
// @Description Api
// @Produce  json
// @Security ApiKeyAuth
// //@Param id path int true "ID"
// @Param payload body DeleteSchema true "参数信息: {rows:[1,2]}"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/api [delete]
func (r Api) Delete(c *gin.Context) {
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
		if err:= db.DeleteById(tx, &sys.Api{}, _id, []string{}, true); err!=nil{
			response.Error(c, err)
			return
		}
	}

	response.DeleteSuccess(c)
}


func (r Api) Option() db.Option {
	var o db.Option
	//o.Select = "distinct id,name,path,method,ctime,mtime"
	return o
}

func (r Api) UpdateApi(conditions []sys.Api) error {
	tx :=r.DB.Begin()
	defer func() {tx.Rollback()}()
	var obj []sys.Api
	var pageData = response.PageDataList{Page: 1,PageSize:0,List:&obj}
	var o = r.Option()
	o.All = true
	var cdsStrList []string
	for _,v:=range conditions{
		cdsStrList = append(cdsStrList, fmt.Sprintf("path='%s' and method='%s'", v.Path, v.Method))
	}
	o.Where = strings.Join(cdsStrList, " or ")
	err := db.Query(tx, &sys.Api{},o, &pageData)
	if err!=nil{
		return err
	}

	// 过滤新增的列表
	var newRows []sys.Api
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
	var deleteList []sys.Api
	var deleteIdList []uint
	if err:=tx.Find(&deleteList, "id not in (?)", existIdList).Error;err!=nil{return err}

	for _,v:=range deleteList{
		deleteIdList = append(deleteIdList, v.ID)
		if _, err := casbin.Enforcer.RemoveFilteredNamedPolicy("p", 1, dom, v.Path, v.Method);err!=nil{return err}
	}
	if err:=tx.Exec("DELETE FROM role_api WHERE api_id in (?)", deleteIdList).Error;err!=nil{return err}
	if err:=tx.Delete(&sys.Api{}, "id in (?)", deleteIdList).Error;err!=nil{return err}
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
func BatchSave(db *gorm.DB, perms []sys.Api) error {
	var buffer bytes.Buffer
	sql := fmt.Sprintf("insert into `%s` (`name`,`path`,`method`) values", name.Api)
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