package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/v-mars/frame/db"
	"github.com/v-mars/frame/pkg/convert"
	"github.com/v-mars/frame/pkg/utils"
	"github.com/v-mars/frame/response"
	"github.com/v-mars/sys/app"
	"github.com/v-mars/sys/app/models/sys"
	"github.com/v-mars/sys/app/utils/casbin"
	"gorm.io/gorm"
	"strconv"
)

type IUser interface {
	GetAll(c *gin.Context)
	ResetPassword(c *gin.Context)
	ChangePassword(c *gin.Context)
	ChangeProfile(c *gin.Context)
	app.CommonInterfaces
}
var dom = "sys"
type SUser struct {
	DB *gorm.DB
}

var Model = &sys.User{}

func NewService(DB *gorm.DB) IUser {
	return SUser{DB: DB}
}

// Get
// @Tags 用户管理
// @Summary 用户详细
// @Description 用户
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "用户id"
// @Security ApiKeyAuth
// @Param username query string false "用户名"
// @Param nickname query string false "用户显示名"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 200 object response.Data {"code": 4000, "status": "error", "message": "error"}
// @Failure 200 object response.Data {"code": 5000, "status": "error", "message": "error"}
// @Router /api/v1/sys/user/{id} [get]
func (r SUser) Get(c *gin.Context) {
	id := c.Params.ByName("id")
	o := r.Option()
	o.Where = "user.id = ?"
	o.Value = append(o.Value, id)
	o.Preloads = []string{"Roles"}
	o.First = true
	o.NullError = true
	var obj User
	err := db.Get(r.DB,&User{}, o, &obj)
	if err!= nil {
		response.Error(c, err)
	} else {
		response.Success(c,&obj)
	}
}

// GetAll
// @Tags 用户管理
// @Summary 所有用户列表
// @Description 用户
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 200 object response.Data {"code": 4000, "status": "error", "message": "error"}
// @Failure 200 object response.Data {"code": 5000, "status": "error", "message": "error"}
// @Router /api/v1/sys/users [get]
func (r SUser) GetAll(c *gin.Context) {
	var obj []AllUser
	var pageData = response.InitPageData(c, &obj, true)
	o := db.Option{}
	o.Where = "status = ?"
	o.Value=append(o.Value, true)
	o.Select = "distinct user.id, nickname, username"
	o.Order = "ID DESC"
	o.Scan = true
	err := db.Query(r.DB,&sys.User{}, o, &pageData)
	if err != nil {
		response.Error(c, err)
		return
	} else {
		response.Success(c, pageData)
	}

}

// Query
// @Tags 用户管理
// @Summary 用户列表
// @Description 用户
// @Produce  json
// @Security ApiKeyAuth
// @Param page query int false "page"
// @Param pageSize query int false "pageSize"
// @Param username query string false "用户名"
// @Param nickname query string false "用户显示名"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/user [get]
func (r SUser) Query(c *gin.Context) {
	var obj []User
	var pageData = response.InitPageData(c, &obj, false)
	type Param struct {
		Username       string `form:"username"` // `form:"name" binding:"required"`
		Nickname       string `form:"nickname"` // `form:"name" binding:"required"`
		Roles          []string `form:"roles[]"`
	}
	var param Param
	if err := c.ShouldBindQuery(&param);err!=nil{
		response.Error(c, err)
		return
	}
	//fmt.Println("Query:", db.DB.)
	o := r.Option()
	o.Where = "username like ? and nickname like ?"
	o.Value = append(o.Value, "%"+param.Username+"%", "%"+param.Nickname+"%")
	if len(param.Roles)>0{
		o.Where = o.Where + " and role.id in (?)"
		o.Value = append(o.Value, param.Roles)
	}
	o.Order = "user.ID DESC"
	o.Group = "user.id"
	o.Preloads = []string{"Roles"}
	//o.Scan = true
	err := db.Query(r.DB, &User{}, o, &pageData)
	if err != nil {
		response.Error(c, err)
		return
	} else {
		response.Success(c, pageData)
	}
}

// Create
// @Tags 用户管理
// @Summary 创建用户
// @Description 用户
// @Produce  json
// @Security ApiKeyAuth
// @Param payload body  PostSchema true "参数信息"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/user [post]
func (r SUser) Create(c *gin.Context) {
	u, err:= GetUserValue(c)
	if err!=nil{
		response.Error(c, err)
		return
	}
	var obj PostSchema
	if err:= c.ShouldBindJSON(&obj);err!=nil{
		response.Error(c, err)
		return
	}
	var newRow sys.User
	newRow.Nickname = obj.Nickname
	newRow.Username = obj.Username
	newRow.Password = utils.SHA256HashString(obj.Password)
	newRow.Email = obj.Email
	newRow.Phone = obj.Phone
	tx :=r.DB.Begin()
	defer func() {tx.Rollback()}()
	if err:= tx.Model(&sys.Role{}).Where("id in (?)", obj.Roles).Select("id,name").Scan(&newRow.Roles).Error;err!=nil{
		response.Error(c, err)
		return
	}
	for _,v:=range newRow.Roles{
		_, err = casbin.Enforcer.AddGroupingPolicy(newRow.Username,v.Name, dom) // user role dom
		if err!=nil{
			response.Error(c, err)
			return
		}
	}
	//fmt.Println(user, reflect.TypeOf(user))
	newRow.ByUpdate = u.Nickname
	if err:= db.Create(tx, &newRow, true);err!=nil{
		response.Error(c, err)
		return
	}
	response.CreateSuccess(c, newRow)
}

// Update
// @Tags 用户管理
// @Summary 更新用户
// @Description 用户
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "ID"
// @Param payload body  PutSchema true "参数信息"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/user/{id} [put]
func (r SUser) Update(c *gin.Context) {
	u, err:= GetUserValue(c)
	if err!=nil{
		response.Error(c, err)
		return
	}
	//nickname, _:=c.Get("nickname")
	var obj PutSchema
	if err := c.ShouldBindJSON(&obj); err!=nil{
		response.ParamFailed(c, err)
		return
	}
	obj.ByUpdate = u.Nickname
	_id := c.Params.ByName("id")
	var MapData map[string]interface{}
	if MapData,err=convert.StructToMap(obj); err!=nil{
		response.Error(c, err)
		return
	}
	tx :=r.DB.Begin()
	defer func() {tx.Rollback()}()
	var ass []db.Association
	if obj.Roles != nil{
		var roles []sys.Role
		if err:= tx.Model(&sys.Role{}).Where("id in (?)", *obj.Roles).Select("id,name").Scan(&roles).Error;err!=nil{
			response.Error(c, err)
			return
		}
		ass =append(ass, db.Association{Column: "Roles", Values: &roles})
		var userObj sys.User
		if err:=tx.First(&userObj, _id).Error;err!=nil{
			response.Error(c, err)
			return
		}

		var newStrArray []string
		for _,v:=range roles{
			newStrArray=append(newStrArray, v.Name)
			_, err = casbin.Enforcer.AddGroupingPolicy(userObj.Username,v.Name, dom) // user role dom
			if err!=nil{
				response.Error(c, err)
				return
			}
		}
		existsList, err := casbin.Enforcer.GetRolesForUser(userObj.Username, dom)
		if err!=nil{
			response.Error(c, err)
			return
		}
		//fmt.Println("roles:", existsList)
		diff := utils.Difference(existsList, newStrArray)
		//fmt.Println("diff:", diff)
		for _,v :=range diff{
			_, err = casbin.Enforcer.RemoveGroupingPolicy(userObj.Username, v, dom)
			if err!=nil{
				response.Error(c, err)
				return
			}
		}

	}
	if obj.Path!=nil{
		MapData["path"] = obj.Path
	}
	var user sys.User
	if err:= db.UpdateById(tx, &user,_id,MapData,ass, true);err!=nil{
		response.Error(c, err)
		return
	}

	response.UpdateSuccess(c, user)
}

// Delete
// @Tags 用户管理
// @Summary 删除用户
// @Description 用户
// @Produce  json
// @Security ApiKeyAuth
// //@Param payload body [] true "用户id list"
// //@Param id path int true "用户id"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/user [delete]
// //@Router /api/v1/sys/user/{id} [delete]
func (r SUser) Delete(c *gin.Context) {
	//var obj sys.User
	var data map[string][]int
	if err := c.ShouldBindJSON(&data); err!=nil{
		response.Error(c, err)
		return
	}
	tx :=r.DB.Begin()
	defer func() {tx.Rollback()}()
	var us []sys.User
	if err:=tx.Model(&sys.User{}).Find(&us,"id in (?)", data["rows"]).Error;err!=nil{
		response.Error(c, err)
		return
	}
	for _,v:=range us{
		if _, err:= casbin.Enforcer.RemoveFilteredNamedGroupingPolicy("g", 0, v.Username, "", dom);err!=nil{
			response.Error(c, err)
			return
		}
	}
	//time.Sleep(time.Second*30)
	for _, i := range data["rows"]{
		//var o = model.Option{Where: "id = ?", Value: []interface{}{i}}
		if err:= db.DeleteById(tx, &sys.User{}, i, []string{"Roles"}, true); err!=nil{
			response.Error(c, err)
			return
		}
	}

	response.DeleteSuccess(c)
}

// ResetPassword
// @Tags 用户管理
// @Summary 更新用户密码
// @Description 用户
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "ID"
// @Param payload body  RestPass true "参数信息"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/user-pass-reset/{id} [put]
func (r SUser) ResetPassword (c *gin.Context)  {
	var obj RestPass
	if err := c.ShouldBindJSON(&obj); err!=nil{
		response.Error(c, err)
		return
	}
	encPassword := utils.SHA256HashString(obj.Password)
	if err:= r.DB.Model(&sys.User{}).Where("id = ?", obj.ID).Updates(
		map[string]interface{}{"password": encPassword}).Error;err!=nil{
		response.Error(c, err)
		return
	}
	response.SuccessMsg(c, "密码更新成功", map[string]string{})
}

// ChangePassword
// @Tags 用户管理
// @Summary 用户密码修改
// @Description 用户
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "ID"
// @Param payload body  RestPass true "参数信息"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/user-pass/{id} [put]
func (r SUser) ChangePassword (c *gin.Context)  {
	u, err:= GetUserValue(c)
	if err!=nil{
		response.Error(c, err)
		return
	}
	_id := c.Params.ByName("id")
	if _id!= strconv.Itoa(int(u.ID)){
		response.Error(c, fmt.Errorf("当前用户%s[%d]非法修改用户密码",u.Username,u.ID))
		return
	}
	var obj ChangePass
	if err = c.ShouldBindJSON(&obj); err!=nil{
		response.Error(c, err)
		return
	}
	auth := Auth(r.DB,u.Username, obj.OldPassword)
	if !auth {
		response.Error(c, fmt.Errorf("旧密码不正确"))
		return
	}

	if obj.Password!=obj.RePassword {
		response.Error(c, fmt.Errorf("密码不一致"))
		return
	}
	encPassword := utils.SHA256HashString(obj.Password)
	if err= r.DB.Model(&sys.User{}).Where("id = ?", _id).Updates(
		map[string]interface{}{"password": encPassword}).Error;err!=nil{
		response.Error(c, err)
		return
	}
	response.SuccessMsg(c, "密码更新成功", map[string]string{})
}

// ChangeProfile
// @Tags 用户管理
// @Summary 更新用户信息
// @Description 用户
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "ID"
// @Param payload body  RestPass true "参数信息"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/sys/user-profile/{id} [put]
func (r SUser) ChangeProfile (c *gin.Context)  {
	u, err:= GetUserValue(c)
	if err!=nil{
		response.Error(c, err)
		return
	}
	_id := c.Params.ByName("id")
	if _id != strconv.Itoa(int(u.ID)){
		response.Error(c, fmt.Errorf("当前用户%s[%d]非法修改用户信息",u.Username,u.ID))
		return
	}
	var obj ChangeProfile
	if err = c.ShouldBindJSON(&obj); err!=nil{
		response.Error(c, err)
		return
	}
	obj.ByUpdate = u.Nickname
	var MapData map[string]interface{}
	if MapData,err=convert.StructToMap(obj); err!=nil{
		response.Error(c, err)
		return
	}
	if err= r.DB.Model(&sys.User{}).Where("id = ?", _id).Updates(MapData).Error;err!=nil{
		response.Error(c, err)
		return
	}
	response.SuccessMsg(c, "更新成功", map[string]string{})
}

func (r SUser) Option() db.Option {
	var o db.Option
	o.Select = "distinct user.id as id, user.nickname, user.username, phone, " +
		"email, user_type,user.path, status,user.by_update,user.ctime, user.mtime"
	o.Joins = "left join user_roles on user_roles.user_id = user.id left join role on user_roles.role_id = role.id "
	return o
}


func Auth(DB *gorm.DB,username, password string) bool {
	var count int64
	if err:= DB.Model(&sys.User{}).Where("username = ? and password = ? and user_type = ? and status = ?",
		username, utils.SHA256HashString(password), "local", true).Count(&count).Error;err!=nil {
		return false
	}
	if count == 0 {
		return false
	} else {
		return true
	}

}

func GetUserInfo(c *gin.Context)  {
	u, err:= GetUserValue(c)
	if err!=nil{
		response.Error(c, err)
		return
	}
	if u.Name == "" && u.ID == 0 {
		response.Error(c, fmt.Errorf("user info is null"))
	}
	response.Success(c, u)
	return
}

func GetUserValue(c *gin.Context) (InfoUser, error) {
	userInfo := c.Value("userInfo")
	var u InfoUser
	err := utils.AnyToAny(userInfo, &u)
	if err != nil {
		return InfoUser{}, err
	}
	if u.Name == "" && u.ID == 0 {
		return u, fmt.Errorf("user info is null")
	}
	return u, nil
}

func AddUser(DB *gorm.DB,user *sys.User) error {
	user.Password = utils.SHA256HashString(user.Password)
	tx :=DB.Begin()
	defer func() {tx.Rollback()}()
	for _,v:=range user.Roles{
		_, err := casbin.Enforcer.AddGroupingPolicy(user.Username,v.Name, dom) // user role dom
		if err!=nil{
			return err
		}
	}
	if err:= db.Create(tx, &user, true);err!=nil{
		return err
	}
	return nil
}