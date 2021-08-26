package router

import (
	"github.com/gin-gonic/gin"
	"github.com/v-mars/frame/db"
	"github.com/v-mars/frame/pkg/logger"
	"github.com/v-mars/frame/pkg/utils"
	"github.com/v-mars/sys/app/sys"
	"net"
)

// RegisterRouter 注册/api路由
func RegisterRouter(engine *gin.Engine) error {
	v1 := engine.Group("/api/v1")
	sys.Api(v1, db.DB)

	engine.GET("/ping", Ping)
	engine.GET("/gin/routes", Query)

	return nil
}

var adminRoles = []string{"administrators", "admin", "root"}
func callHandlerPermission(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles = utils.Union(roles, adminRoles)
		//fmt.Println("roles:", roles)
		userRoles := []string{"bb", "cc", "devops"}
		intersects := utils.Intersect(roles, userRoles)
		//fmt.Println("intersect:", intersects, len(intersects))

		/**/
		if len(intersects) > 0 {
			c.Next()
			return
		}else {
			c.JSON(200, gin.H{
				"status": "error",
				"message": "没有权限",
				"code": 401,
			})
			c.Abort()
			return
		}

	}
}


func Ping(c *gin.Context)  {
	logger.Infof("from client %s ping", c.ClientIP())
	ips,_:= ServerIPv4s()
	c.JSON(200, gin.H{
		"app": "sys",
		"message": "pong",
		"status": "success",
		"clientIP": c.ClientIP(),
		"serverIPs": ips,
	})
}

// ServerIPv4s LocalIPs return all non-loopback IPv4 addresses
func ServerIPv4s() ([]string, error) {
	var ips []string
	adders, err := net.InterfaceAddrs()

	if err != nil {
		return ips, err
	}

	for _, a := range adders {
		if inet, ok := a.(*net.IPNet); ok && !inet.IP.IsLoopback() && inet.IP.To4() != nil {
			ips = append(ips, inet.IP.String())
		}
	}
	return ips, nil
}