package dept

import (
	"github.com/gin-gonic/gin"
	"github.com/v-mars/frame/db"
	"github.com/v-mars/frame/pkg/convert"
	"github.com/v-mars/frame/pkg/utils"
	"github.com/v-mars/frame/response"
	"github.com/v-mars/sys/app/models/sys"
	"gorm.io/gorm"
	"strconv"
)

type IDept interface {
	GetAll(c *gin.Context)
	Query(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type Dept struct {
	DB *gorm.DB
}

func NewService(DB *gorm.DB) IDept {
	return Dept{DB: DB}
}

// GetAll
// @Tags 部门管理
// @Summary 所有部门
// @Description 部门
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/depts [get]
func (r Dept) GetAll(c *gin.Context) {
	var api sys.Dept
	tree, err := api.GetListTree(r.DB)
	if err != nil {
		response.Error(c, err)
		return
	} else {
		response.Success(c, tree)
	}

}

// Query
// @Tags 部门管理
// @Summary 部门列表
// @Description 部门
// @Produce  json
// @Security ApiKeyAuth
// @Param name query string false "部门名"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/dept [get]
func (r Dept) Query(c *gin.Context) {
	var obj []sys.Dept
	var pageData = response.InitPageData(c, &obj, false)
	type Param struct {
		Name       string `form:"name"` // `form:"name" binding:"required"`
		Code       string `form:"code"` // `form:"name" binding:"required"`
	}
	var param Param
	if err := c.ShouldBindQuery(&param); err!=nil{
		response.ParamFailed(c, err)
		return
	}
	var o db.Option
	o.Where = "(parent_id = 0 or parent_id is null) and (name like ? or code like ?)"
	o.Value = append(o.Value, "%"+param.Name+"%","%"+param.Code+"%")
	o.Order = "sort asc"

	var api sys.Dept
	tree, err := api.GetMapTree(r.DB)
	if err!=nil{
		response.Error(c, err)
		return
	}

	var data = make([]map[string]interface{},0)
	err = db.Query(r.DB,&obj, o, &pageData)
	if err != nil {
		response.Error(c, err)
		return
	} else {
		for _,v := range obj{
			if tree[v.ID] != nil{
				data = append(data, tree[v.ID])
			}else {
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
// @Tags 部门管理
// @Summary 创建部门
// @Description 部门
// @Produce  json
// @Security ApiKeyAuth
// @Param payload body  sys.Role true "参数信息"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/dept [post]
func (r Dept) Create(c *gin.Context) {
	//u, err:= user.GetUserValue(c)
	//if err!=nil{
	//	response.Error(c, err)
	//	return
	//}
	//var obj PostSchema
	var newRow sys.Dept
	if err := c.ShouldBindJSON(&newRow); err!=nil{
		response.ParamFailed(c, err)
		return
	}
	//newRow.Name = obj.Name
	tx :=r.DB.Begin()
	defer func() {tx.Rollback()}()
	if err:= db.Create(tx, &newRow,true);err!=nil{
		response.Error(c, err)
		return
	}
	response.CreateSuccess(c, newRow)
}

// Update
// @Tags 部门管理
// @Summary 更新部门
// @Description 部门
// @Produce  json
// @Security ApiKeyAuth
// @Param payload body  sys.Role true "参数信息"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/dept [put]
func (r Dept) Update(c *gin.Context) {
	_id := c.Params.ByName("id")
	var err error
	var obj PutSchema
	if err = c.ShouldBindJSON(&obj); err!=nil{
		response.Error(c, err)
		return
	}
	if _id == strconv.Itoa(int(*obj.ParentID)) {
		*obj.ParentID=0
	}

	var MapData map[string]interface{}
	if MapData,err=convert.StructToMap(obj); err!=nil{
		response.Error(c, err)
		return
	}
	if obj.Path!=nil{
		MapData["path"] = obj.Path
	}
	tx :=r.DB.Begin()
	defer func() {tx.Rollback()}()
	var ass []db.Association
	if err= db.UpdateById(tx, &sys.Dept{},_id,MapData,ass, true);err!=nil{
		response.Error(c, err)
		return
	}

	response.UpdateSuccess(c)
}

// Delete
// @Tags 部门管理
// @Summary 删除部门
// @Description 部门
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "部门id"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/dept/{id} [delete]
func (r Dept) Delete(c *gin.Context) {
	_id := c.Params.ByName("id")
	tx :=r.DB.Begin()
	defer func() {tx.Rollback()}()
	if err:= db.DeleteById(tx, &sys.Dept{}, _id, []string{}, true); err!=nil{
		response.Error(c, err)
		return
	}
	response.DeleteSuccess(c)
}

