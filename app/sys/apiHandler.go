package sys

import (
	"github.com/gin-gonic/gin"
	"github.com/v-mars/sys/app/sys/api"
	//"github.com/v-mars/sys/app/sys/auth"
	"github.com/v-mars/sys/app/sys/dept"
	"github.com/v-mars/sys/app/sys/portal"
	"github.com/v-mars/sys/app/sys/property"
	"github.com/v-mars/sys/app/sys/role"
	"github.com/v-mars/sys/app/sys/server"
	"github.com/v-mars/sys/app/sys/tree_node"
	"github.com/v-mars/sys/app/sys/user"
	"github.com/v-mars/sys/app/sys/usergroup"
	"gorm.io/gorm"
)

func ApiHandler(rg *gin.RouterGroup, DB *gorm.DB)  {
	//sysG := rg.Group("/sys")
	//app.POST("login", auth.LoginAuth)
	//app.POST("refresh-token", auth.RefreshToken)
	//app.GET("/user-info", user.GetUserInfo)
	//sys := app.Group("/sys")
	{
		var userInter = user.NewService(DB)
		rg.GET("/user/:id", userInter.Get)
		rg.GET("/users", userInter.GetAll)
		rg.GET("/user", userInter.Query)
		rg.POST("/user", userInter.Create)
		rg.PUT("/user-pass-reset/:id", userInter.ResetPassword)
		rg.PUT("/user-pass/:id", userInter.ChangePassword)
		rg.PUT("/user-profile/:id", userInter.ChangeProfile)
		rg.PUT("/user/:id", userInter.Update)
		rg.DELETE("/user", userInter.Delete)

		var userGroup = usergroup.NewService(DB)
		rg.GET("/usergroups", userGroup.GetAll)
		rg.GET("/usergroup", userGroup.Query)
		rg.POST("/usergroup", userGroup.Create)
		rg.PUT("/usergroup", userGroup.Update)
		rg.DELETE("/usergroup", userGroup.Delete)

		var roleInter = role.NewService(DB)
		rg.GET("/roles", roleInter.GetAll)
		rg.GET("/role", roleInter.Query)
		rg.POST("/role", roleInter.Create)
		rg.PUT("/role", roleInter.Update)
		rg.DELETE("/role", roleInter.Delete)

		var iApi = api.NewService(DB)
		rg.GET("/apis", iApi.GetAllPerm)
		rg.GET("/api-group", iApi.GetAllGroup)
		rg.GET("/api", iApi.Query)
		rg.POST("/api", iApi.Create)
		rg.PUT("/api/:id", iApi.Update)
		rg.DELETE("/api", iApi.Delete)

		var deptInter = dept.NewService(DB)
		rg.GET("/depts", deptInter.GetAll)
		rg.GET("/dept", deptInter.Query)
		rg.POST("/dept", deptInter.Create)
		rg.PUT("/dept/:id", deptInter.Update)
		rg.DELETE("/dept/:id", deptInter.Delete)

		var portalM = portal.NewService(DB)
		rg.GET("/portal-type", portalM.GetAllType)
		rg.GET("/portals", portalM.GetAll)
		rg.GET("/portal", portalM.Query)
		rg.POST("/portal", portalM.Create)
		rg.PUT("/portal/:id", portalM.Update)
		rg.DELETE("/portal", portalM.Delete)
		rg.POST("/portal/favor/:id", portalM.FavorCreate)
		rg.DELETE("/portal/favor/:id", portalM.FavorDelete)

		var tree = tree_node.NewService(DB)
		rg.GET("/tree_node_mark", tree.GetAllMark)
		rg.GET("/tree_node", tree.Query)
		rg.POST("/tree_node", tree.Create)
		rg.PUT("/tree_node/:id/rename", tree.Rename)
		rg.PUT("/tree_node", tree.Update)
		rg.DELETE("/tree_node/:id/:mark", tree.Delete)

		var propertySVC = property.NewService(DB)

		rg.GET("/property", propertySVC.Query)
		rg.GET("/ldap", propertySVC.QueryLDAP)
		rg.POST("/ldap", propertySVC.CreateOrUpdateLDAP)
		rg.POST("/check-ldap-connect", propertySVC.CheckLDAPConnect)
		rg.GET("/email", propertySVC.QueryEmail)
		rg.POST("/email", propertySVC.CreateOrUpdateEmail)
		rg.GET("/aliyun", propertySVC.QueryAliYun)
		rg.POST("/aliyun", propertySVC.CreateOrUpdateAliYun)
		rg.GET("/gitlab", propertySVC.QueryGitlab)
		rg.POST("/gitlab", propertySVC.CreateOrUpdateGitlab)
		rg.GET("/jenkins", propertySVC.QueryJenkins)
		rg.POST("/jenkins", propertySVC.CreateOrUpdateJenkins)

		var state = server.NewService(DB)
		rg.GET("/state", state.Get)
	}
}
