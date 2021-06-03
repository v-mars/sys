package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/v-mars/frame/db"
	"github.com/v-mars/frame/pkg/convert"
	goJWT "github.com/v-mars/frame/pkg/jwt"
	"github.com/v-mars/frame/pkg/logger"
	"github.com/v-mars/frame/pkg/utils"
	"github.com/v-mars/frame/response"
	"github.com/v-mars/sys/app/config"
	"github.com/v-mars/sys/app/models/sys"
	"github.com/v-mars/sys/app/sys/property"
	"github.com/v-mars/sys/app/sys/role"
	"github.com/v-mars/sys/app/sys/user"
	"github.com/v-mars/sys/app/utils/ldap"
	"gorm.io/gorm"
	"time"
)

// LoginAuth
// @Tags 登录认证
// @Summary 登录
// @Description 登录
// @Produce  json
// @Param payload body  Params true "参数信息"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/login [post]
func LoginAuth(c *gin.Context)  {
	var obj Params
	if err:= c.ShouldBindJSON(&obj);err!=nil{
		response.Error(c, err)
		return
	}

	var localAuth int64
	if err:= db.DB.Model(&sys.User{}).Where("username = ? and password = ? and user_type = ?",
		obj.Username, utils.SHA256HashString(obj.Password), "local").Count(&localAuth).Error;err!=nil {
		//fmt.Println("gorm.IsRecordNotFoundError(err) localAuth", gorm.IsRecordNotFoundError(err))
		response.Error(c, err)
		return
	}

	if localAuth > 0 {
		data, err := returnResult(db.DB, obj.Username, "local",ldap.Result{},config.Configs.JWT.Age)
		if err!=nil{response.Error(c,err);return}
		response.Success(c, data)
		return
	}

	//var ldapApi = ldap.NewLDAP()
	var ldapApi, err = property.GetLdapApiFromDB(db.DB)
	if err!=nil{
		logger.Errorf("ldap get ldap(from db config) api  err: %s",err)
	} else {
		r, ldapAuthErr := ldapApi.Authentication(obj.Username, obj.Password)

		if ldapAuthErr ==  nil {
			data, err := returnResult(db.DB, obj.Username, "ldap",r,config.Configs.JWT.Age)
			if err!=nil{response.Error(c,err);return}
			response.Success(c, data)
			return
		} else {
			logger.Errorf("ldap auth err: %s",ldapAuthErr)
		}
	}

	logger.Errorf(fmt.Sprintf("[%s]认证失败，用户名或密码不对,请重新输入", obj.Username))
	response.ErrorNoStack(c, fmt.Errorf("[%s]认证失败，用户名或密码不对,请重新输入！", obj.Username))
	return

}

// RefreshToken
// @Tags 登录认证
// @Summary 刷新Token
// @Description 刷新Token
// @Produce  json
// @Param payload body  RefreshParams true "参数信息"
// @Success 200 object response.Data {"code": 2000, "status": "ok", "message": "success", "data": ""}
// @Failure 400 object response.Data {"code": 4001, "status": "error", "message": "error", "data": ""}
// @Router /api/v1/refresh-token [post]
func RefreshToken(c *gin.Context)  {
	var obj RefreshParams
	if err:= c.ShouldBindJSON(&obj);err!=nil{
		response.Error(c, err)
		return
	}
	tokenApi := goJWT.New()
	tokenApi.SetKey([]byte(config.Configs.JWT.RefreshKey))
	claims, err := tokenApi.ParseToken(obj.Token,config.Configs.JWT.RefreshKey)
	if err !=nil{
		response.Error(c, err)
		return
	}


	tokenApi.SetClaims(claims)
	data,err := TokenMethod(tokenApi, claims, config.Configs.JWT.Age)
	if err != nil {
		response.Error(c, err, map[string]interface{}{})
		return
	}
	data["roles"] = claims["roles"]
	logger.Infof("%s(%s) token刷新成功!", data["nickname"],data["username"])
	response.Success(c, data)
	return
}

func TokenMethod(tokenApi goJWT.JWTAuth, cla jwt.MapClaims, expired int) (map[string]interface{}, error) {
	tokenApi.Options.TokenType = config.Configs.JWT.TokenType //"Bearer"
	tokenApi.Options.SigningKey = []byte(config.Configs.JWT.TokenKey)
	tokenApi.Options.Expired = expired
	tokenApi.Options.Claims = cla
	tokenString, err := tokenApi.GenerateToken()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	expiresAt := time.Now().Add(time.Duration(tokenApi.Options.Expired) * time.Second).Unix()
	expireTime := now.Add(time.Second * time.Duration(tokenApi.Options.Expired)).Format("2006-01-02 15:04:05")
	tokenApi.Options.SigningKey = []byte(config.Configs.JWT.RefreshKey)
	tokenApi.Options.Expired = config.Configs.JWT.Age + 3600
	RefreshTokenString, err := tokenApi.GenerateToken()
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{
		"token":         tokenString.GetAccessToken(),
		"refresh_token": RefreshTokenString.GetAccessToken(),
		"expire_time":   expireTime,
		"expires_at":    expiresAt,
		"user":          map[string]interface{}{},
		"user_type":     "local",
		"token_type":    tokenString.GetTokenType(),   // tokenString.GetTokenType()
		"nickname":      cla["nickname"],
		"username":      cla["username"],
		"email":         cla["email"],
	}
	return data, nil
}

func returnResult(DB *gorm.DB, username string, userType string,r ldap.Result, expired int) (map[string]interface{}, error) {
	type User struct {
		ID        uint   `json:"id"`
		Nickname  string `json:"nickname"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		Status    bool   `json:"status"`
	}
	var userObj sys.User
	if err:= DB.Table("user").Where("username = ? and user_type = ?", username,
		userType).Select("id, nickname, username, email, status").First(&userObj).Error;
		errors.Is(err, gorm.ErrRecordNotFound) && userType=="ldap"{
		var roles []sys.Role
		DB.Model(&sys.Role{}).Where("name='normal'").Find(&roles)
		userObj = sys.User{Username: r.Username,Nickname: r.Username,Email: r.Email,Phone: r.Phone,Status: true,
			UserType: userType,Roles: roles, BaseByUpdate: db.BaseByUpdate{ByUpdate: "ldap"}}
		if len(r.DisplayName)>0{userObj.Nickname = r.DisplayName}
		if len(r.Email)==0{userObj.Email = fmt.Sprintf("%s@example.com",username)}
		if err=user.AddUser(db.DB,&userObj); err != nil {
			return nil, err
		}
	} else if err!=nil {
		return nil, err
	}
	if !userObj.Status{
		return nil, fmt.Errorf("用户名%s已经被禁用", username)
	}

	var iRole role.IRole = &role.Role{}
	roleList, err:= iRole.GetUserRoles(username)
	if err != nil {return nil, err}
	var u = user.InfoUser{ID: userObj.ID,Name:userObj.Nickname,Nickname:userObj.Nickname,Username:userObj.Username,
		Email:userObj.Email,Roles:roleList}
	var cla map[string]interface{}
	if err := convert.StructToMapOut(u, &cla); err!=nil{
		return nil, fmt.Errorf("user info parse claims err: %s", err)
	}
	tokenApi := goJWT.New()
	tokenApi.SetKey([]byte(config.Configs.JWT.TokenKey))
	tokenApi.SetClaims(cla)
	data,err := TokenMethod(tokenApi, cla, expired)
	if err != nil {return nil, err}
	data["roles"] = roleList
	logger.Infof("%s(%s) 登陆成功!", userObj.Nickname,userObj.Username)
	return data, nil
}
