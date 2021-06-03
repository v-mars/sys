package portal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/v-mars/frame/db"
	"github.com/v-mars/frame/pkg/convert"
	"github.com/v-mars/frame/response"
	"github.com/v-mars/sys/app/models/sys"
	"github.com/v-mars/sys/app/sys/user"
	"gorm.io/gorm"
	"strings"
)

type IPortal interface {
	Get(c *gin.Context)
	Query(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Option(userId interface{}) db.Option
	GetAll(c *gin.Context)
	GetAllType(c *gin.Context)
	FavorCreate(c *gin.Context)
	FavorDelete(c *gin.Context)
}

type Portal struct {
	DB *gorm.DB
}

func NewService(DB *gorm.DB) IPortal {
	return Portal{DB: DB}
}

// Get
// @Tags Portal管理
// @Summary Portal详细
// @Description Portal
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "ID"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/portal/{id} [get]
func (r Portal) Get(c *gin.Context) {
	u, err:= user.GetUserValue(c)
	if err!=nil{
		response.Error(c, err)
		return
	}

	_id := c.Params.ByName("id")
	var obj ShowData
	o := r.Option(u.ID)
	o.Where = "id = ?"
	o.Value = append(o.Value, _id)
	o.First = true
	o.NullError = true
	err = db.Get(r.DB,&obj, o, &obj)
	if err!= nil {
		response.Error(c, err)
		return
	} else {
		response.Success(c,&obj)
		return
	}
}

// GetAllType
// @Tags Portal管理
// @Summary Portal type列表
// @Description Portal
// @Produce  json
// @Security ApiKeyAuth
// @Param name query string false "Portal名"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/portal-type [get]
func (r Portal) GetAllType(c *gin.Context) {
	var obj []ShowPortalType
	var o = r.Option(nil)
	o.Select = "DISTINCT type"
	o.Joins = ""
	o.Order = "type"
	o.Scan = true
	tx :=r.DB.Begin()
	defer func() {
		tx.Rollback()
	}()
	err := db.Get(tx,&sys.Portal{}, o, &obj)
	if err != nil {
		response.Error(c, err)
		return
	} else {
		var pageData = response.PageDataList{PageNumber: 1,PageSize:0,List:&obj,Total: int64(len(obj))}
		response.Success(c, pageData)
	}
}

// GetAll
// @Tags Portal管理
// @Summary Portal all列表
// @Description Portal
// @Produce  json
// @Security ApiKeyAuth
// @Param name query string false "Portal名"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/portal-all [get]
func (r Portal) GetAll(c *gin.Context) {
	u, err:= user.GetUserValue(c)
	if err!=nil{
		response.Error(c, err)
		return
	}
	var obj []ShowData
	var pageData = response.InitPageData(c, &obj, true)
	type Param struct {
		Name     string `form:"name"`
		Types    []string `form:"types[]"`
	}
	var param Param
	if err := c.ShouldBindQuery(&param); err!=nil{
		response.ParamFailed(c, err)
		return
	}
	var o = r.Option(u.ID)
	o.Where = "name like ?"
	o.Value = append(o.Value, "%"+param.Name+"%")
	if len(param.Types)>0{
		o.Where = o.Where + " and type in (?)"
		o.Value = append(o.Value, param.Types)
	}
	o.Order = "type,id ASC"
	o.Scan = true
	tx :=r.DB.Begin()
	defer func() {
		tx.Rollback()
	}()
	err = db.Query(tx,&sys.Portal{}, o, &pageData)
	if err != nil {
		response.Error(c, err)
		return
	} else {
		type Res struct {
			Favor       []ShowData            `json:"favor"`
			Default     []ShowData            `json:"default"`
			TypeMapData map[string][]ShowData `json:"type_map_data"`
			//All         []ShowData            `json:"all"`
		}
		tmp:=map[string][]ShowData{}
		Result := Res{TypeMapData: tmp}
		for _,v:=range obj{
			if v.IsFavor>0{
				Result.Favor = append(Result.Favor, v)
			}
			if len(v.Type)==0 || strings.ToLower(v.Type) == "default"{
				Result.Default = append(Result.Default, v)
			} else {
				Result.TypeMapData[v.Type] = append(Result.TypeMapData[v.Type], v)
			}
		}
		response.Success(c, Result)
	}
}

// Query
// @Tags Portal管理
// @Summary Portal列表
// @Description Portal
// @Produce  json
// @Security ApiKeyAuth
// @Param name query string false "Portal名"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/portal [get]
func (r Portal) Query(c *gin.Context) {
	u, err:= user.GetUserValue(c)
	if err!=nil{
		response.Error(c, err)
		return
	}
	var obj []ShowData
	var pageData = response.InitPageData(c, &obj, false)
	type Param struct {
		Name     string `form:"name"`
		Types    []string `form:"types[]"`
	}
	var param Param
	if err := c.ShouldBindQuery(&param); err!=nil{
		response.ParamFailed(c, err)
		return
	}
	var o = r.Option(u.ID)
	o.Where = "name like ?"
	o.Value = append(o.Value, "%"+param.Name+"%")
	if len(param.Types)>0{
		o.Where = o.Where + " and type in (?)"
		o.Value = append(o.Value, param.Types)
	}
	o.Order = "ID DESC"
	o.Scan = true
	tx :=r.DB.Begin()
	defer func() {
		tx.Rollback()
	}()
	err = db.Query(tx,&sys.Portal{}, o, &pageData)
	if err != nil {
		response.Error(c, err)
		return
	} else {
		response.Success(c, pageData)
	}
}

// Create
// @Tags Portal管理
// @Summary 创建Portal
// @Description Portal
// @Produce  json
// @Security ApiKeyAuth
// @Param payload body PostSchema true "参数信息"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/portal [post]
func (r Portal) Create(c *gin.Context) {
	u, err:= user.GetUserValue(c)
	if err!=nil{
		response.Error(c, err)
		return
	}
	var obj PostSchema
	if err := c.ShouldBindJSON(&obj); err!=nil{
		response.ParamFailed(c, err)
		return
	}
	obj.ByUpdate = u.Nickname
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
// @Tags Portal管理
// @Summary 更新Portal
// @Description Portal
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "ID"
// @Param payload body PutSchema true "参数信息"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/portal/{id} [put]
func (r Portal) Update(c *gin.Context) {
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
	obj.ByUpdate = u.Nickname
	_id := c.Params.ByName("id")
	var MapData map[string]interface{}
	if MapData,err=convert.StructToMap(obj); err!=nil{
		response.Error(c, err)
		return
	}
	//fmt.Println("mapData", MapData)
	tx :=r.DB.Begin()
	defer func() {
		tx.Rollback()
	}()
	var res sys.Portal
	if err:= db.UpdateById(tx, &res, _id,MapData,nil, true);err!=nil{
		response.Error(c, err)
		return
	}
	response.UpdateSuccess(c, res)
}

// Delete
// @Tags Portal管理
// @Summary 删除Portal
// @Description Portal
// @Produce  json
// @Security ApiKeyAuth
// //@Param id path int true "ID"
// @Param payload body DeleteSchema true "参数信息: {rows:[1,2]}"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/portal [delete]
func (r Portal) Delete(c *gin.Context) {
	//_id := c.Params.ByName("id")
	var data map[string][]int
	if err:= c.ShouldBindJSON(&data);err!=nil{
		response.Error(c, err)
		return
	}
	tx :=r.DB.Begin()
	defer func() {
		tx.Rollback()
	}()
	for _,_id := range data["rows"]{
		if err:= db.DeleteById(tx, &sys.Portal{}, _id, []string{}, true); err!=nil{
			response.Error(c, err)
			return
		}
	}

	response.DeleteSuccess(c)
}

// FavorCreate
// @Tags Portal管理
// @Summary Portal收藏
// @Description Portal
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/portal/favor/:id [post]
func (r Portal) FavorCreate(c *gin.Context) {
	u, err:= user.GetUserValue(c)
	if err!=nil{
		response.Error(c, err)
		return
	}
	if u.ID < 1{
		response.Error(c, fmt.Errorf("用户信息不存在，请联系管理员。"))
		return
	}
	_id := c.Params.ByName("id")
	tx :=r.DB.Begin()
	defer func() {tx.Rollback()}()
	var Po = sys.Portal{}
	err = db.GetById(tx,&sys.Portal{},_id,&Po,true)
	if err!=nil{
		response.Error(c, err)
		return
	}
	if err:=tx.Model(&Po).Association("Favors").Append(&sys.User{BaseID: db.BaseID{ID: u.ID}});err!=nil{
		response.Error(c, err)
		return
	}

	if err:=tx.Commit().Error;err!=nil{
		tx.Rollback()
		response.Error(c, err)
		return
	}

	response.SuccessMsg(c, "收藏成功",Po)
}

// FavorDelete
// @Tags Portal管理
// @Summary Portal收藏取消
// @Description Portal
// @Produce  json
// @Security ApiKeyAuth
// //@Param id path int true "ID"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/portal/favor/:id [delete]
func (r Portal) FavorDelete(c *gin.Context) {
	u, err:= user.GetUserValue(c)
	if err!=nil{
		response.Error(c, err)
		return
	}
	if u.ID < 1{
		response.Error(c, fmt.Errorf("用户信息不存在，请联系管理员。"))
		return
	}
	portalId := c.Params.ByName("id")
	tx :=r.DB.Begin()
	defer func() {
		tx.Rollback()
	}()
	var Po = sys.Portal{}
	err = db.GetById(tx, &sys.Portal{},portalId,&Po, true)
	if err!=nil{
		response.Error(c, err)
		return
	}
	if err:=tx.Model(&Po).Association("Favors").Delete(&sys.User{BaseID: db.BaseID{ID: u.ID}});err!=nil{
		response.Error(c, err)
		return
	}
	if err:=tx.Commit().Error;err!=nil{
		tx.Rollback()
		response.Error(c, err)
		return
	}
	response.SuccessMsg(c, "已取消收藏",Po)
}


func (r Portal) Option(userId interface{}) db.Option {
	var o db.Option
	o.Select = "distinct id,name,description,url,icon,type,portal_favors.user_id as is_favor,by_update,ctime,mtime"
	o.Joins = fmt.Sprintf(
		"left join `portal_favors` on `portal_favors`.user_id=%d and `portal_favors`.portal_id=`portal`.id", userId)
	return o
}
