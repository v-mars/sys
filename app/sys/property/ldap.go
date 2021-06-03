package property

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/v-mars/frame/pkg/convert"
	"github.com/v-mars/frame/pkg/logger"
	"github.com/v-mars/frame/pkg/utils"
	"github.com/v-mars/frame/response"
	utilsApi "github.com/v-mars/sys/app/utils"
	"github.com/v-mars/sys/app/utils/ldap"
	"gorm.io/gorm"
)

type LDAP struct {
	Enable    string `json:"enable"`	// "true" or "false"
	Addr      string `json:"addr"`
	TLS       string `json:"tls"`
	ManagerDN string `json:"username"`
	Password  string `json:"password"`
	BaseDN    string `json:"base_dn"`
	Filter    string `json:"filter"`
}

type LDAPCheckUser struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func (r Property) QueryLDAP(c *gin.Context) {
	data,err := GetProperty(r.DB,"ldap")
	if err != nil{
		response.Error(c, err)
		return
	}
	if data["password"] != "" {
		data["password"] = utils.Base64Dec(convert.ToString(data["password"]))
	}
	response.Success(c, map[string]interface{}{"result": data})
	return
}

func (r Property) CreateOrUpdateLDAP(c *gin.Context) {
	type Param struct {
		//CheckEmail  string                 `json:"check_email"`
		Row        map[string]interface{} `json:"row" binding:"required"`
	}
	var obj Param
	if err := c.ShouldBindJSON(&obj); err!=nil{
		response.ParamFailed(c, err)
		return
	}
	obj.Row["password"] = utils.Base64Enc([]byte(convert.ToString(obj.Row["password"])))
	if err:= CreateOrUpdate(r.DB,"ldap", obj.Row); err!=nil{
		response.Error(c, err)
		return
	}
	response.SuccessMsg(c, "保存成功！", "")
}

func (r Property) CheckLDAPConnect(c *gin.Context) {
	type Param struct {
		User LDAPCheckUser `json:"user" binding:"required"`
		Ldap LDAP          `json:"ldap" binding:"required"`
	}
	var obj Param
	if err := c.ShouldBindJSON(&obj); err!=nil{
		response.ParamFailed(c, err)
		return
	}
	var ldapApi = ldap.NewLDAP()
	ldapApi.Logger = logger.Logger
	ldapApi.Option.Addr = obj.Ldap.Addr
	ldapApi.Option.AuthFilter = obj.Ldap.Filter
	ldapApi.Option.BaseDN = obj.Ldap.BaseDN
	ldapApi.Option.Username = obj.Ldap.ManagerDN
	ldapApi.Option.Password = obj.Ldap.Password
	ldapApi.Option.Tls = obj.Ldap.TLS == "true"
	authentication, err := ldapApi.Authentication(obj.User.Username,obj.User.Password)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessMsg(c, "LDAP连接测试成功！", authentication)
}

func GetLdapApiFromDB(DB *gorm.DB) (*ldap.LDAP,error) {
	data,err := GetProperty(DB,"ldap")
	if err != nil{
		return nil,fmt.Errorf("get ldap data err: %s", err)
	}
	if data["password"] != "" {
		data["password"] = utils.Base64Dec(convert.ToString(data["password"]))
	}
	var obj LDAP
	err = utilsApi.StructToStruct(data, &obj)
	if err != nil {
		return nil,fmt.Errorf("map convert ldap struct err: %s", err)
	}
	if obj.Enable != "true" || len(obj.Addr) == 0 || len(obj.ManagerDN) == 0 {
		return nil, fmt.Errorf("ldap data is none")
	}
	var ldapApi = ldap.NewLDAP()
	ldapApi.Logger = logger.Logger
	ldapApi.Option.Addr = obj.Addr
	ldapApi.Option.AuthFilter = obj.Filter
	ldapApi.Option.BaseDN = obj.BaseDN
	ldapApi.Option.Username = obj.ManagerDN
	ldapApi.Option.Password = obj.Password
	ldapApi.Option.Tls = obj.TLS == "true"
	return ldapApi,nil
}