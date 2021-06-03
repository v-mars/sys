package sys

import (
	"github.com/gin-gonic/gin"
	"github.com/v-mars/sys/app/sys/auth"
	"github.com/v-mars/sys/app/sys/permission"
	"github.com/v-mars/sys/app/sys/portal"
	"github.com/v-mars/sys/app/sys/property"
	"github.com/v-mars/sys/app/sys/role"
	"github.com/v-mars/sys/app/sys/server"
	"github.com/v-mars/sys/app/sys/tree_node"
	"github.com/v-mars/sys/app/sys/user"
	"github.com/v-mars/sys/app/sys/usergroup"
	"gorm.io/gorm"
)

func Api(app *gin.RouterGroup, d *gorm.DB)  {
	app.POST("login", auth.LoginAuth)
	app.POST("refresh-token", auth.RefreshToken)
	app.GET("/user-info", user.GetUserInfo)
	sys := app.Group("/sys")
	{
		var iUser = user.NewService(d)
		sys.GET("/iUser/:id", iUser.Get)
		sys.GET("/users", iUser.GetAll)
		sys.GET("/iUser", iUser.Query)
		sys.POST("/iUser", iUser.Create)
		sys.PUT("/iUser/:id", iUser.Update)
		sys.DELETE("/iUser", iUser.Delete)

		var userGroup = usergroup.NewService(d)
		sys.GET("/usergroups", userGroup.GetAll)
		sys.GET("/usergroup", userGroup.Query)
		sys.POST("/usergroup", userGroup.Create)
		sys.PUT("/usergroup", userGroup.Update)
		sys.DELETE("/usergroup", userGroup.Delete)

		var iRole = role.NewService(d)
		sys.GET("/roles", iRole.GetAll)
		sys.GET("/iRole", iRole.Query)
		sys.POST("/iRole", iRole.Create)
		sys.PUT("/iRole", iRole.Update)
		sys.DELETE("/iRole", iRole.Delete)

		var permissionM = permission.NewService(d)
		sys.GET("/permissions", permissionM.Query)
		sys.GET("/permission", permissionM.Query)
		sys.POST("/permission", permissionM.Create)
		sys.PUT("/permission/:id", permissionM.Update)
		sys.DELETE("/permission", permissionM.Delete)

		var portalM = portal.NewService(d)
		sys.GET("/portal-type", portalM.GetAllType)
		sys.GET("/portals", portalM.GetAll)
		sys.GET("/portal", portalM.Query)
		sys.POST("/portal", portalM.Create)
		sys.PUT("/portal/:id", portalM.Update)
		sys.DELETE("/portal", portalM.Delete)
		sys.POST("/portal/favor/:id", portalM.FavorCreate)
		sys.DELETE("/portal/favor/:id", portalM.FavorDelete)

		var tree = tree_node.NewService(d)
		sys.GET("/tree_node_mark", tree.GetAllMark)
		sys.GET("/tree_node", tree.Query)
		sys.POST("/tree_node", tree.Create)
		sys.PUT("/tree_node/:id/rename", tree.Rename)
		sys.PUT("/tree_node", tree.Update)
		sys.DELETE("/tree_node/:id/:mark", tree.Delete)

		var iProperty = property.NewService(d)

		sys.GET("/iProperty", iProperty.Query)
		sys.GET("/ldap", iProperty.QueryLDAP)
		sys.POST("/ldap", iProperty.CreateOrUpdateLDAP)
		sys.POST("/check-ldap-connect", iProperty.CheckLDAPConnect)
		sys.GET("/email", iProperty.QueryEmail)
		sys.POST("/email", iProperty.CreateOrUpdateEmail)
		sys.GET("/aliyun", iProperty.QueryAliYun)
		sys.POST("/aliyun", iProperty.CreateOrUpdateAliYun)
		sys.GET("/gitlab", iProperty.QueryGitlab)
		sys.POST("/gitlab", iProperty.CreateOrUpdateGitlab)
		sys.GET("/jenkins", iProperty.QueryJenkins)
		sys.POST("/jenkins", iProperty.CreateOrUpdateJenkins)

		var state = server.NewService(d)
		sys.GET("/state", state.Get)
	}
}
